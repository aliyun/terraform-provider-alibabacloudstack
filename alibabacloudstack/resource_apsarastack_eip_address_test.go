package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEipAddress0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_eip_address.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEipAddressCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribeeipaddressesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEipAddressBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test",

					"isp": "BGP",

					"address_name": "rdktest",

					"netmode": "public",

					"bandwidth": "1",

					"payment_type": "PayAsYouGo",

					"region_id": "cn-hangzhou",

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::defaultubT769.ResourceGroupId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"isp": "BGP",

						"address_name": "rdktest",

						"netmode": "public",

						"bandwidth": "1",

						"payment_type": "PayAsYouGo",

						"region_id": "cn-hangzhou",

						"resource_group_id": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test",

					"address_name": "rdktest",

					"bandwidth": "1000",

					"region_id": "cn-hangzhou",

					"instance_type": "Nat",

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::defaultubT769.ResourceGroupId)}}",

					"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

					"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

					"high_definition_monitor_log_status": "ON",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"address_name": "rdktest",

						"bandwidth": "1000",

						"region_id": "cn-hangzhou",

						"instance_type": "Nat",

						"resource_group_id": CHECKSET,

						"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

						"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

						"high_definition_monitor_log_status": "ON",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test",

					"isp": "BGP",

					"address_name": "rdktest",

					"netmode": "public",

					"bandwidth": "1",

					"payment_type": "PayAsYouGo",

					"region_id": "cn-hangzhou",

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::defaultubT769.ResourceGroupId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"isp": "BGP",

						"address_name": "rdktest",

						"netmode": "public",

						"bandwidth": "1",

						"payment_type": "PayAsYouGo",

						"region_id": "cn-hangzhou",

						"resource_group_id": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test",

					"address_name": "rdktest",

					"bandwidth": "1000",

					"region_id": "cn-hangzhou",

					"instance_type": "Nat",

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::defaultubT769.ResourceGroupId)}}",

					"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

					"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

					"high_definition_monitor_log_status": "ON",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"address_name": "rdktest",

						"bandwidth": "1000",

						"region_id": "cn-hangzhou",

						"instance_type": "Nat",

						"resource_group_id": CHECKSET,

						"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

						"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

						"high_definition_monitor_log_status": "ON",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test",

					"address_name": "rdktest",

					"bandwidth": "1000",

					"region_id": "cn-hangzhou",

					"instance_type": "Nat",

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::defaultubT769.ResourceGroupId)}}",

					"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

					"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

					"high_definition_monitor_log_status": "ON",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"address_name": "rdktest",

						"bandwidth": "1000",

						"region_id": "cn-hangzhou",

						"instance_type": "Nat",

						"resource_group_id": CHECKSET,

						"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

						"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

						"high_definition_monitor_log_status": "ON",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test",

					"address_name": "rdktest",

					"bandwidth": "1000",

					"region_id": "cn-hangzhou",

					"instance_type": "Nat",

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::defaultubT769.ResourceGroupId)}}",

					"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

					"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

					"high_definition_monitor_log_status": "ON",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"address_name": "rdktest",

						"bandwidth": "1000",

						"region_id": "cn-hangzhou",

						"instance_type": "Nat",

						"resource_group_id": CHECKSET,

						"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

						"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

						"high_definition_monitor_log_status": "ON",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test",

					"address_name": "rdktest",

					"bandwidth": "1000",

					"region_id": "cn-hangzhou",

					"instance_type": "Nat",

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::defaultubT769.ResourceGroupId)}}",

					"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

					"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

					"high_definition_monitor_log_status": "ON",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"address_name": "rdktest",

						"bandwidth": "1000",

						"region_id": "cn-hangzhou",

						"instance_type": "Nat",

						"resource_group_id": CHECKSET,

						"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

						"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

						"high_definition_monitor_log_status": "ON",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test",

					"address_name": "rdktest",

					"bandwidth": "1000",

					"region_id": "cn-hangzhou",

					"instance_type": "Nat",

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::defaultubT769.ResourceGroupId)}}",

					"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

					"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

					"high_definition_monitor_log_status": "ON",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"address_name": "rdktest",

						"bandwidth": "1000",

						"region_id": "cn-hangzhou",

						"instance_type": "Nat",

						"resource_group_id": CHECKSET,

						"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

						"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

						"high_definition_monitor_log_status": "ON",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test",

					"address_name": "rdktest",

					"bandwidth": "1000",

					"region_id": "cn-hangzhou",

					"instance_type": "Nat",

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::defaultubT769.ResourceGroupId)}}",

					"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

					"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

					"high_definition_monitor_log_status": "ON",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"address_name": "rdktest",

						"bandwidth": "1000",

						"region_id": "cn-hangzhou",

						"instance_type": "Nat",

						"resource_group_id": CHECKSET,

						"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

						"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

						"high_definition_monitor_log_status": "ON",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test",

					"address_name": "rdktest",

					"bandwidth": "1000",

					"region_id": "cn-hangzhou",

					"instance_type": "Nat",

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::defaultubT769.ResourceGroupId)}}",

					"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

					"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

					"high_definition_monitor_log_status": "ON",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"address_name": "rdktest",

						"bandwidth": "1000",

						"region_id": "cn-hangzhou",

						"instance_type": "Nat",

						"resource_group_id": CHECKSET,

						"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

						"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

						"high_definition_monitor_log_status": "ON",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test",

					"address_name": "rdktest",

					"bandwidth": "1000",

					"region_id": "cn-hangzhou",

					"instance_type": "Nat",

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::defaultubT769.ResourceGroupId)}}",

					"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

					"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

					"high_definition_monitor_log_status": "ON",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"address_name": "rdktest",

						"bandwidth": "1000",

						"region_id": "cn-hangzhou",

						"instance_type": "Nat",

						"resource_group_id": CHECKSET,

						"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

						"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

						"high_definition_monitor_log_status": "ON",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test",

					"address_name": "rdktest",

					"bandwidth": "1000",

					"region_id": "cn-hangzhou",

					"instance_type": "Nat",

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::defaultubT769.ResourceGroupId)}}",

					"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

					"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

					"high_definition_monitor_log_status": "ON",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"address_name": "rdktest",

						"bandwidth": "1000",

						"region_id": "cn-hangzhou",

						"instance_type": "Nat",

						"resource_group_id": CHECKSET,

						"log_store": "${{ref(resource, SLS::LogStore::2.0.0::defaultStore.LogstoreName)}}",

						"log_project": "${{ref(resource, SLS::Project::2.0.0::defaultProject.ProjectName)}}",

						"high_definition_monitor_log_status": "ON",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccEipAddressCheckmap = map[string]string{

	"resource_group_id": CHECKSET,

	"high_definition_monitor_log_status": CHECKSET,

	"second_limited": CHECKSET,

	"log_project": CHECKSET,

	"reservation_order_type": CHECKSET,

	"segment_instance_id": CHECKSET,

	"expired_time": CHECKSET,

	"bandwidth_package_id": CHECKSET,

	"reservation_active_time": CHECKSET,

	"reservation_bandwidth": CHECKSET,

	"tags": CHECKSET,

	"status": CHECKSET,

	"instance_id": CHECKSET,

	"log_store": CHECKSET,

	"deletion_protection": CHECKSET,

	"bandwidth_package_type": CHECKSET,

	"bandwidth_package_bandwidth": CHECKSET,

	"internet_charge_type": CHECKSET,

	"address_name": CHECKSET,

	"netmode": CHECKSET,

	"security_protection_types": CHECKSET,

	"description": CHECKSET,

	"allocation_id": CHECKSET,

	"bandwidth": CHECKSET,

	"payment_type": CHECKSET,

	"instance_type": CHECKSET,

	"create_time": CHECKSET,

	"isp": CHECKSET,

	"mode": CHECKSET,

	"has_reservation_data": CHECKSET,

	"operation_locks": CHECKSET,

	"eip_bandwidth": CHECKSET,

	"ip_address": CHECKSET,

	"region_id": CHECKSET,
}

func AlibabacloudTestAccEipAddressBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
