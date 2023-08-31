/*
Copyright © 2023 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: MPL-2.0
*/

package spec

import (
	"strconv"

	valid "github.com/asaskevich/govalidator"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	packageinstallmodel "github.com/vmware/terraform-provider-tanzu-mission-control/internal/models/tanzupackageinstall"
	common "github.com/vmware/terraform-provider-tanzu-mission-control/internal/resources/common"
)

func ConstructSpecForClusterScope(d *schema.ResourceData) (spec *packageinstallmodel.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageInstallSpec) {
	value, ok := d.GetOk(SpecKey)
	if !ok {
		return spec
	}

	data, _ := value.([]interface{})

	if len(data) == 0 || data[0] == nil {
		return spec
	}

	specData := data[0].(map[string]interface{})

	spec = &packageinstallmodel.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageInstallSpec{
		PackageRef: &packageinstallmodel.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageMetadataPackagePackageRef{
			VersionSelection: &packageinstallmodel.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageMetadataPackageVersionSelection{},
		},
	}

	if data, ok := specData[PackageRefKey]; ok {
		if v1, ok := data.([]interface{}); ok && len(v1) != 0 {
			specData := v1[0].(map[string]interface{})

			var metadataName string

			if v, ok := specData[PackageMetadataNameKey]; ok {
				metadataName = v.(string)
			}

			spec.PackageRef.PackageMetadataName = metadataName

			if versionSelectionData, ok := specData[VersionSelectionKey]; ok {
				if v1, ok := versionSelectionData.([]interface{}); ok && len(v1) != 0 {
					specData := v1[0].(map[string]interface{})

					var constraints string

					if v, ok := specData[ConstraintsKey]; ok {
						constraints = v.(string)
					}

					spec.PackageRef.VersionSelection.Constraints = constraints
				}
			}
		}
	}

	spec.RoleBindingScope = packageinstallmodel.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageInstallRoleBindingScopeCLUSTER.Pointer()

	if v, ok := specData[InlineValuesKey]; ok {
		if v1, ok := v.(map[string]interface{}); ok {
			for key, value := range v1 {
				switch {
				case valid.IsInt(value.(string)):
					number, err := strconv.ParseUint(value.(string), 10, 32)
					if err != nil {
						v1[key] = value.(string)
						break
					}

					finalIntNum := int(number) // Convert uint64 To int
					v1[key] = finalIntNum
				case valid.IsFloat(value.(string)):
					floatNum, err := strconv.ParseFloat(value.(string), 64)
					if err != nil {
						v1[key] = value.(string)
						break
					}

					v1[key] = floatNum
				default:
					v1[key] = value.(string)
				}
			}

			spec.InlineValues = v1
		}
	}

	return spec
}

func FlattenSpecForClusterScope(spec *packageinstallmodel.VmwareTanzuManageV1alpha1ClusterNamespaceTanzupackageInstallSpec) (data []interface{}) {
	if spec == nil {
		return data
	}

	flattenSpecData := make(map[string]interface{})

	pkgRefSpec := make(map[string]interface{})

	pkgMetadataName := spec.PackageRef.PackageMetadataName
	constraints := spec.PackageRef.VersionSelection.Constraints

	versionSelectionSpec := []interface{}{
		map[string]interface{}{
			ConstraintsKey: constraints,
		},
	}

	pkgRefSpec[PackageMetadataNameKey] = pkgMetadataName
	pkgRefSpec[VersionSelectionKey] = versionSelectionSpec

	if v1, ok := spec.InlineValues.(map[string]interface{}); ok {
		inline := common.GetTypeStringMapData(v1)
		flattenSpecData[InlineValuesKey] = inline
	} else {
		flattenSpecData[InlineValuesKey] = spec.InlineValues
	}

	flattenSpecData[RoleBindingScopeKey] = string(*spec.RoleBindingScope)

	flattenSpecData[PackageRefKey] = []interface{}{pkgRefSpec}

	return []interface{}{flattenSpecData}
}
