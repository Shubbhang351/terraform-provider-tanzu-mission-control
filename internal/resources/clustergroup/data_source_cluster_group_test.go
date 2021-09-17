/*
Copyright © 2021 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: MPL-2.0
*/

package clustergroup

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	testhelper "gitlab.eng.vmware.com/olympus/terraform-provider-tanzu/internal/resources/testing"
)

func TestAcceptanceForClusterGroupDataSource(t *testing.T) {
	var provider = initTestProvider(t)

	clusterGroup := acctest.RandomWithPrefix("tf-cg-test")
	dataSourceName := "data.tmc_cluster_group.test_data_cluster_group"
	resourceName := "tmc_cluster_group.test_cluster_group"

	resource.Test(t, resource.TestCase{
		PreCheck:          testhelper.TestPreCheck(t),
		ProviderFactories: testhelper.GetTestProviderFactories(provider),
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: getTestDataSourceClusterGroupConfigValue(clusterGroup, testhelper.MetaTemplate),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceAttributes(dataSourceName, resourceName),
				),
			},
		},
	},
	)
	t.Log("cluster group data source acceptance test complete!")
}

func getTestDataSourceClusterGroupConfigValue(clusterGroupName, meta string) string {
	return fmt.Sprintf(`
resource "tmc_cluster_group" "test_cluster_group" {
  name = "%s"
  %s
}

data "tmc_cluster_group" "test_data_cluster_group" {
  name = tmc_cluster_group.test_cluster_group.name
}
`, clusterGroupName, meta)
}

func checkDataSourceAttributes(dataSourceName, resourceName string) resource.TestCheckFunc {
	var check = []resource.TestCheckFunc{
		verifyClusterGroupDataSource(dataSourceName),
		resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
		resource.TestCheckResourceAttrSet(dataSourceName, "id"),
	}

	check = append(check, testhelper.MetaDataSourceAttributeCheck(dataSourceName, resourceName)...)

	return resource.ComposeTestCheckFunc(check...)
}

func verifyClusterGroupDataSource(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module does not have cluster group resource %s", name)
		}

		return nil
	}
}
