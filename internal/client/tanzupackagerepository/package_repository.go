/*
Copyright © 2023 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: MPL-2.0
*/

package tanzupackagerepositoryclient

import (
	"net/url"

	"github.com/vmware/terraform-provider-tanzu-mission-control/internal/client/transport"
	"github.com/vmware/terraform-provider-tanzu-mission-control/internal/helper"
	pkgrepository "github.com/vmware/terraform-provider-tanzu-mission-control/internal/models/tanzupackagerepository"
)

const (
	apiVersionAndGroup                 = "v1alpha1/clusters"
	apiKind                            = "tanzupackage/repositories"
	namespaces                         = "namespaces"
	queryParamKeyManagementClusterName = "fullName.managementClusterName"
	queryParamKeyProvisionerName       = "fullName.provisionerName"
)

// New creates a new repository resource service API client.
func New(transport *transport.Client) ClientService {
	return &Client{Client: transport}
}

/*
Client for repository resource service API.
*/
type Client struct {
	*transport.Client
}

type ClientService interface {
	RepositoryResourceServiceCreate(request *pkgrepository.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageRepositoryRepositoryRequest) (*pkgrepository.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageRepositoryResponse, error)

	RepositoryResourceServiceDelete(fn *pkgrepository.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageRepositoryFullName) error

	RepositoryResourceServiceGet(fn *pkgrepository.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageRepositoryFullName) (*pkgrepository.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageRepositoryGetResponse, error)

	RepositoryResourceServiceUpdate(request *pkgrepository.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageRepositoryRepositoryRequest) (*pkgrepository.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageRepositoryResponse, error)
}

/*
RepositoryResourceServiceCreate creates a repository.
*/
func (c *Client) RepositoryResourceServiceCreate(request *pkgrepository.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageRepositoryRepositoryRequest) (*pkgrepository.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageRepositoryResponse, error) {
	// "/v1alpha1/clusters/{repository.fullName.clusterName}/namespaces/{repository.fullName.namespaceName}/tanzupackage/repositories"
	requestURL := helper.ConstructRequestURL(apiVersionAndGroup, request.Repository.FullName.ClusterName, namespaces, request.Repository.FullName.NamespaceName, apiKind).String()
	pkgRepositoryResponse := &pkgrepository.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageRepositoryResponse{}
	err := c.Create(requestURL, request, pkgRepositoryResponse)

	return pkgRepositoryResponse, err
}

/*
RepositoryResourceServiceDelete deletes a repository.
*/
func (c *Client) RepositoryResourceServiceDelete(fn *pkgrepository.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageRepositoryFullName) error {
	queryParams := url.Values{}

	if fn.ManagementClusterName != "" {
		queryParams[queryParamKeyManagementClusterName] = []string{fn.ManagementClusterName}
	}

	if fn.ProvisionerName != "" {
		queryParams[queryParamKeyProvisionerName] = []string{fn.ProvisionerName}
	}

	// "/v1alpha1/clusters/{fullName.clusterName}/namespaces/{fullName.namespaceName}/tanzupackage/repositories/{fullName.name}"
	requestURL := helper.ConstructRequestURL(apiVersionAndGroup, fn.ClusterName, namespaces, fn.NamespaceName, apiKind, fn.Name).AppendQueryParams(queryParams).String()

	return c.Delete(requestURL)
}

/*
RepositoryResourceServiceGet gets a repository.
*/
func (c *Client) RepositoryResourceServiceGet(fn *pkgrepository.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageRepositoryFullName) (*pkgrepository.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageRepositoryGetResponse, error) {
	queryParams := url.Values{}

	if fn.ManagementClusterName != "" {
		queryParams[queryParamKeyManagementClusterName] = []string{fn.ManagementClusterName}
	}

	if fn.ProvisionerName != "" {
		queryParams[queryParamKeyProvisionerName] = []string{fn.ProvisionerName}
	}

	// "/v1alpha1/clusters/{fullName.clusterName}/namespaces/{fullName.namespaceName}/tanzupackage/repositories/{fullName.name}"
	requestURL := helper.ConstructRequestURL(apiVersionAndGroup, fn.ClusterName, namespaces, fn.NamespaceName, apiKind, fn.Name).AppendQueryParams(queryParams).String()
	pkgRepositoryResponse := &pkgrepository.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageRepositoryGetResponse{}
	err := c.Get(requestURL, pkgRepositoryResponse)

	return pkgRepositoryResponse, err
}

/*
RepositoryResourceServiceUpdate updates overwrite a repository.
*/
func (c *Client) RepositoryResourceServiceUpdate(request *pkgrepository.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageRepositoryRepositoryRequest) (*pkgrepository.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageRepositoryResponse, error) {
	// "/v1alpha1/clusters/{repository.fullName.clusterName}/namespaces/{repository.fullName.namespaceName}/tanzupackage/repositories/{repository.fullName.name}"
	requestURL := helper.ConstructRequestURL(apiVersionAndGroup, request.Repository.FullName.ClusterName, namespaces, request.Repository.FullName.NamespaceName, apiKind, request.Repository.FullName.Name).String()
	pkgRepositoryResponse := &pkgrepository.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageRepositoryResponse{}
	err := c.Update(requestURL, request, pkgRepositoryResponse)

	return pkgRepositoryResponse, err
}
