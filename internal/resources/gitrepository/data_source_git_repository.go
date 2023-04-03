/*
Copyright © 2023 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: MPL-2.0
*/

package gitrepository

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"

	"github.com/vmware/terraform-provider-tanzu-mission-control/internal/authctx"
	clienterrors "github.com/vmware/terraform-provider-tanzu-mission-control/internal/client/errors"
	"github.com/vmware/terraform-provider-tanzu-mission-control/internal/helper"
	gitrepositoryclustermodel "github.com/vmware/terraform-provider-tanzu-mission-control/internal/models/gitrepository/cluster"
	gitrepositoryclustergroupmodel "github.com/vmware/terraform-provider-tanzu-mission-control/internal/models/gitrepository/clustergroup"
	objectmetamodel "github.com/vmware/terraform-provider-tanzu-mission-control/internal/models/objectmeta"
	"github.com/vmware/terraform-provider-tanzu-mission-control/internal/resources/common"
	commonscope "github.com/vmware/terraform-provider-tanzu-mission-control/internal/resources/common/scope"
	"github.com/vmware/terraform-provider-tanzu-mission-control/internal/resources/gitrepository/scope"
	"github.com/vmware/terraform-provider-tanzu-mission-control/internal/resources/gitrepository/spec"
	"github.com/vmware/terraform-provider-tanzu-mission-control/internal/resources/gitrepository/status"
)

func DataSourceGitRepository() *schema.Resource {
	return &schema.Resource{
		Schema: gitRepositorySchema,
		ReadContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
			return dataSourceGitRepositoryRead(helper.GetContextWithCaller(ctx, helper.DataRead), d, m)
		},
	}
}

func dataSourceGitRepositoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	config := m.(authctx.TanzuContext)

	gitRepositoryName, ok := d.Get(nameKey).(string)
	if !ok {
		return diag.Errorf("unable to read git repository name")
	}

	gitRepositoryNamespaceName, ok := d.Get(namespaceNameKey).(string)
	if !ok {
		return diag.Errorf("unable to read git repository namespace name")
	}

	scopedFullnameData := scope.ConstructScope(d, gitRepositoryName, gitRepositoryNamespaceName)

	if scopedFullnameData == nil {
		return diag.Errorf("Unable to get Tanzu Mission Control git repository entry; Scope full name is empty")
	}

	err := enableContinuousDelivery(&config, scopedFullnameData, common.ConstructMeta(d))
	if err != nil {
		return diag.FromErr(errors.Wrapf(err, "Unable to get Tanzu Mission Control git repository entry, name : %s", gitRepositoryName))
	}

	UID, meta, atomicSpec, clusterScopeStatus, clusterGroupScopeStatus, err := retrieveGitRepositoryUIDMetaAndSpecFromServer(config, scopedFullnameData, d)
	if err != nil {
		if clienterrors.IsNotFoundError(err) && !helper.IsDataRead(ctx) {
			_ = schema.RemoveFromState(d, m)
			return
		}

		return diag.FromErr(err)
	}

	// always run
	d.SetId(UID)

	if err := d.Set(common.MetaKey, common.FlattenMeta(meta)); err != nil {
		return diag.FromErr(err)
	}

	var (
		flattenedSpec   []interface{}
		flattenedStatus interface{}
	)

	switch scopedFullnameData.Scope {
	case commonscope.ClusterScope:
		flattenedSpec = spec.FlattenSpecForClusterScope(atomicSpec)
		flattenedStatus = status.FlattenStatusForClusterScope(clusterScopeStatus)
	case commonscope.ClusterGroupScope:
		clusterGroupScopeSpec := &gitrepositoryclustergroupmodel.VmwareTanzuManageV1alpha1ClustergroupNamespaceFluxcdGitrepositorySpec{
			AtomicSpec: atomicSpec,
		}
		flattenedSpec = spec.FlattenSpecForClusterGroupScope(clusterGroupScopeSpec)
		flattenedStatus = status.FlattenStatusForClusterGroupScope(clusterGroupScopeStatus)
	}

	if err := d.Set(spec.SpecKey, flattenedSpec); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(statusKey, flattenedStatus); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// nolint: gocognit
func retrieveGitRepositoryUIDMetaAndSpecFromServer(config authctx.TanzuContext, scopedFullnameData *scope.ScopedFullname, d *schema.ResourceData) (
	string,
	*objectmetamodel.VmwareTanzuCoreV1alpha1ObjectMeta,
	*gitrepositoryclustermodel.VmwareTanzuManageV1alpha1ClusterNamespaceFluxcdGitrepositorySpec,
	*gitrepositoryclustermodel.VmwareTanzuManageV1alpha1ClusterNamespaceFluxcdGitrepositoryStatus,
	*gitrepositoryclustergroupmodel.VmwareTanzuManageV1alpha1ClustergroupNamespaceFluxcdGitrepositoryStatus,
	error) {
	var (
		UID                     string
		meta                    *objectmetamodel.VmwareTanzuCoreV1alpha1ObjectMeta
		spec                    *gitrepositoryclustermodel.VmwareTanzuManageV1alpha1ClusterNamespaceFluxcdGitrepositorySpec
		clusterScopeStatus      *gitrepositoryclustermodel.VmwareTanzuManageV1alpha1ClusterNamespaceFluxcdGitrepositoryStatus
		clusterGroupScopeStatus *gitrepositoryclustergroupmodel.VmwareTanzuManageV1alpha1ClustergroupNamespaceFluxcdGitrepositoryStatus
	)
	// nolint: dupl
	switch scopedFullnameData.Scope {
	case commonscope.ClusterScope:
		if scopedFullnameData.FullnameCluster != nil {
			resp, err := config.TMCConnection.ClusterGitRepositoryResourceService.VmwareTanzuManageV1alpha1ClusterFluxcdGitrepositoryResourceServiceGet(scopedFullnameData.FullnameCluster)
			if err != nil {
				if clienterrors.IsNotFoundError(err) {
					d.SetId("")
					return "", nil, nil, nil, nil, nil
				}

				return "", nil, nil, nil, nil, errors.Wrapf(err, "Unable to get Tanzu Mission Control cluster git repository entry, name : %s", scopedFullnameData.FullnameCluster.Name)
			}

			scopedFullnameData = &scope.ScopedFullname{
				Scope:           commonscope.ClusterScope,
				FullnameCluster: resp.GitRepository.FullName,
			}

			fullName, name, namespace := scope.FlattenScope(scopedFullnameData)

			if err := d.Set(nameKey, name); err != nil {
				return "", nil, nil, nil, nil, err
			}

			if err := d.Set(namespaceNameKey, namespace); err != nil {
				return "", nil, nil, nil, nil, err
			}

			if err := d.Set(commonscope.ScopeKey, fullName); err != nil {
				return "", nil, nil, nil, nil, err
			}

			UID = resp.GitRepository.Meta.UID
			meta = resp.GitRepository.Meta
			spec = resp.GitRepository.Spec
			clusterScopeStatus = resp.GitRepository.Status
		}
	case commonscope.ClusterGroupScope:
		if scopedFullnameData.FullnameClusterGroup != nil {
			resp, err := config.TMCConnection.ClusterGroupGitRepositoryResourceService.VmwareTanzuManageV1alpha1ClustergroupFluxcdGitrepositoryResourceServiceGet(scopedFullnameData.FullnameClusterGroup)
			if err != nil {
				if clienterrors.IsNotFoundError(err) {
					d.SetId("")
					return "", nil, nil, nil, nil, nil
				}

				return "", nil, nil, nil, nil, errors.Wrapf(err, "Unable to get Tanzu Mission Control cluster group git repository entry, name : %s", scopedFullnameData.FullnameClusterGroup.Name)
			}

			scopedFullnameData = &scope.ScopedFullname{
				Scope:                commonscope.ClusterGroupScope,
				FullnameClusterGroup: resp.GitRepository.FullName,
			}

			fullName, name, namespace := scope.FlattenScope(scopedFullnameData)

			if err := d.Set(nameKey, name); err != nil {
				return "", nil, nil, nil, nil, err
			}

			if err := d.Set(namespaceNameKey, namespace); err != nil {
				return "", nil, nil, nil, nil, err
			}

			if err := d.Set(commonscope.ScopeKey, fullName); err != nil {
				return "", nil, nil, nil, nil, err
			}

			UID = resp.GitRepository.Meta.UID
			meta = resp.GitRepository.Meta
			spec = resp.GitRepository.Spec.AtomicSpec
			clusterGroupScopeStatus = resp.GitRepository.Status
		}
	case commonscope.UnknownScope:
		return "", nil, nil, nil, nil, errors.Errorf("no valid scope type block found: minimum one valid scope type block is required among: %v. Please check the schema.", strings.Join(scope.ScopesAllowed[:], `, `))
	}

	return UID, meta, spec, clusterScopeStatus, clusterGroupScopeStatus, nil
}
