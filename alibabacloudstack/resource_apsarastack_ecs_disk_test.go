package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackEcsDisk0(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"description": "加密测试",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"performance_level": "PL0",

					"disk_name": "单测",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"description": "加密测试",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"performance_level": "PL0",

						"disk_name": "单测",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"description": "加密测试",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"performance_level": "PL0",

					"disk_name": "单测",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"description": "加密测试",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"performance_level": "PL0",

						"disk_name": "单测",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlibabacloudTestAccEcsDiskCheckmap = map[string]string{

	"encrypted": CHECKSET,

	"size": CHECKSET,

	"delete_auto_snapshot": CHECKSET,

	"tags": CHECKSET,

	"status": CHECKSET,

	"kms_key_id": CHECKSET,

	"delete_with_instance": CHECKSET,

	"category": CHECKSET,

	"description": CHECKSET,

	"enable_auto_snapshot": CHECKSET,

	"disk_name": CHECKSET,

	"snapshot_id": CHECKSET,
}

func AlibabacloudTestAccEcsDiskBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}




`, name)
}
func TestAccAlibabacloudStackEcsDisk1(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"performance_level": "PL1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"performance_level": "PL1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "11111",

					"disk_name": "tf-testacccn-hangzhouecsdisk69920",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "11111",

						"disk_name": "tf-testacccn-hangzhouecsdisk69920",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"performance_level": "PL1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"performance_level": "PL1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "11111",

					"disk_name": "tf-testacccn-hangzhouecsdisk69920",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "11111",

						"disk_name": "tf-testacccn-hangzhouecsdisk69920",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "11111",

					"disk_name": "tf-testacccn-hangzhouecsdisk69920",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "11111",

						"disk_name": "tf-testacccn-hangzhouecsdisk69920",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "11111",

					"disk_name": "tf-testacccn-hangzhouecsdisk69920",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "11111",

						"disk_name": "tf-testacccn-hangzhouecsdisk69920",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "11111",

					"disk_name": "tf-testacccn-hangzhouecsdisk69920",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "11111",

						"disk_name": "tf-testacccn-hangzhouecsdisk69920",
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
func TestAccAlibabacloudStackEcsDisk2(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"disk_name": "tf-testacccn-hangzhouecsdisk69920",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"disk_name": "tf-testacccn-hangzhouecsdisk69920",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "tf-testAcc-1yGEO",

					"disk_name": "tf-testacccn-hangzhouecsdisk69920",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "tf-testAcc-1yGEO",

						"disk_name": "tf-testacccn-hangzhouecsdisk69920",
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
func TestAccAlibabacloudStackEcsDisk3(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"disk_name": "tf-testacccn-hangzhouecsdisk69920",

					"storage_cluster_id": "1",

					"kms_key_id": "1",

					"payment_type": "Subscription",

					"description": "1",

					"instance_id": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"disk_name": "tf-testacccn-hangzhouecsdisk69920",

						"storage_cluster_id": "1",

						"kms_key_id": "1",

						"payment_type": "Subscription",

						"description": "1",

						"instance_id": "1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "tf-testAcc-1yGEO",

					"disk_name": "tf-testacccn-hangzhouecsdisk69920",

					"payment_type": "PayAsYouGo",

					"category": "cloud_efficiency",

					"instance_id": "2",

					"performance_level": "2",

					"disk_id": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "tf-testAcc-1yGEO",

						"disk_name": "tf-testacccn-hangzhouecsdisk69920",

						"payment_type": "PayAsYouGo",

						"category": "cloud_efficiency",

						"instance_id": "2",

						"performance_level": "2",

						"disk_id": "2",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_ssd",

					"payment_type": "PayAsYouGo",

					"description": "3",

					"instance_id": "3",

					"performance_level": "3",

					"disk_name": "3",

					"disk_id": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_ssd",

						"payment_type": "PayAsYouGo",

						"description": "3",

						"instance_id": "3",

						"performance_level": "3",

						"disk_name": "3",

						"disk_id": "1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"disk_name": "tf-testacccn-hangzhouecsdisk69920",

					"storage_cluster_id": "1",

					"kms_key_id": "1",

					"payment_type": "Subscription",

					"description": "1",

					"instance_id": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"disk_name": "tf-testacccn-hangzhouecsdisk69920",

						"storage_cluster_id": "1",

						"kms_key_id": "1",

						"payment_type": "Subscription",

						"description": "1",

						"instance_id": "1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "tf-testAcc-1yGEO",

					"disk_name": "tf-testacccn-hangzhouecsdisk69920",

					"payment_type": "PayAsYouGo",

					"category": "cloud_efficiency",

					"instance_id": "2",

					"performance_level": "2",

					"disk_id": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "tf-testAcc-1yGEO",

						"disk_name": "tf-testacccn-hangzhouecsdisk69920",

						"payment_type": "PayAsYouGo",

						"category": "cloud_efficiency",

						"instance_id": "2",

						"performance_level": "2",

						"disk_id": "2",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "tf-testAcc-1yGEO",

					"disk_name": "tf-testacccn-hangzhouecsdisk69920",

					"payment_type": "PayAsYouGo",

					"category": "cloud_efficiency",

					"instance_id": "2",

					"performance_level": "2",

					"disk_id": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "tf-testAcc-1yGEO",

						"disk_name": "tf-testacccn-hangzhouecsdisk69920",

						"payment_type": "PayAsYouGo",

						"category": "cloud_efficiency",

						"instance_id": "2",

						"performance_level": "2",

						"disk_id": "2",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "tf-testAcc-1yGEO",

					"disk_name": "tf-testacccn-hangzhouecsdisk69920",

					"payment_type": "PayAsYouGo",

					"category": "cloud_efficiency",

					"instance_id": "2",

					"performance_level": "2",

					"disk_id": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "tf-testAcc-1yGEO",

						"disk_name": "tf-testacccn-hangzhouecsdisk69920",

						"payment_type": "PayAsYouGo",

						"category": "cloud_efficiency",

						"instance_id": "2",

						"performance_level": "2",

						"disk_id": "2",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "tf-testAcc-1yGEO",

					"disk_name": "tf-testacccn-hangzhouecsdisk69920",

					"payment_type": "PayAsYouGo",

					"category": "cloud_efficiency",

					"instance_id": "2",

					"performance_level": "2",

					"disk_id": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "tf-testAcc-1yGEO",

						"disk_name": "tf-testacccn-hangzhouecsdisk69920",

						"payment_type": "PayAsYouGo",

						"category": "cloud_efficiency",

						"instance_id": "2",

						"performance_level": "2",

						"disk_id": "2",
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
func TestAccAlibabacloudStackEcsDisk4(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"description": "test",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"disk_name": "testdisk",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"description": "test",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"disk_name": "testdisk",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "eeeeee",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "eeeeee",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"description": "test",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"disk_name": "testdisk",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"description": "test",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"disk_name": "testdisk",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "eeeeee",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "eeeeee",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "eeeeee",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "eeeeee",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "eeeeee",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "eeeeee",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "eeeeee",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "eeeeee",
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
func TestAccAlibabacloudStackEcsDisk5(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_auto",

					"description": "挂盘测试",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"disk_name": "挂载测试",

					"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_auto",

						"description": "挂盘测试",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"disk_name": "挂载测试",

						"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_auto",

					"description": "挂盘测试",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"disk_name": "挂载测试",

					"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_auto",

						"description": "挂盘测试",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"disk_name": "挂载测试",

						"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
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
func TestAccAlibabacloudStackEcsDisk6(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"description": "test",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"disk_name": "testdisk",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"description": "test",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"disk_name": "testdisk",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "eeeeee",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "eeeeee",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"description": "test",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"disk_name": "testdisk",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"description": "test",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"disk_name": "testdisk",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "eeeeee",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "eeeeee",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "eeeeee",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "eeeeee",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "eeeeee",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "eeeeee",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "eeeeee",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "eeeeee",
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
func TestAccAlibabacloudStackEcsDisk7(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_auto",

					"description": "挂盘测试",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"disk_name": "testDiskName",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_auto",

						"description": "挂盘测试",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"disk_name": "testDiskName",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_auto",

					"description": "挂盘测试",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"disk_name": "testDiskName",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_auto",

						"description": "挂盘测试",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"disk_name": "testDiskName",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",
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
func TestAccAlibabacloudStackEcsDisk8(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"description": "加密测试",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"disk_name": "zone测试",

					"multi_attach": "Disabled",

					"performance_level": "PL2",

					"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.47.pre::key.KeyId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"description": "加密测试",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"disk_name": "zone测试",

						"multi_attach": "Disabled",

						"performance_level": "PL2",

						"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.47.pre::key.KeyId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "zone测试",

					"category": "cloud_auto",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "zone测试",

						"category": "cloud_auto",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"description": "加密测试",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"disk_name": "zone测试",

					"multi_attach": "Disabled",

					"performance_level": "PL2",

					"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.47.pre::key.KeyId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"description": "加密测试",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"disk_name": "zone测试",

						"multi_attach": "Disabled",

						"performance_level": "PL2",

						"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.47.pre::key.KeyId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
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
func TestAccAlibabacloudStackEcsDisk9(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"description": "create",

					"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.3.pre::key.KeyId)}}",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"performance_level": "PL0",

					"disk_name": "createname",

					"multi_attach": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"description": "create",

						"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.3.pre::key.KeyId)}}",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"performance_level": "PL0",

						"disk_name": "createname",

						"multi_attach": "Enabled",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"description": "update",

					"performance_level": "PL1",

					"disk_name": "updatename",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"description": "update",

						"performance_level": "PL1",

						"disk_name": "updatename",

						"payment_type": "PayAsYouGo",
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
func TestAccAlibabacloudStackEcsDisk10(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"disk_name": "CCCC2",

					"instance_id": "alibabacloudstack_instance.ecs.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"disk_name": "CCCC2",

						"instance_id": "alibabacloudstack_instance.ecs.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "bbbb2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "bbbb2",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"disk_name": "CCCC2",

					"instance_id": "alibabacloudstack_instance.ecs.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"disk_name": "CCCC2",

						"instance_id": "alibabacloudstack_instance.ecs.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "bbbb2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "bbbb2",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "bbbb2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "bbbb2",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "bbbb2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "bbbb2",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "bbbb2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "bbbb2",
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
func TestAccAlibabacloudStackEcsDisk11(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"description": "create",

					"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.3.pre::key.KeyId)}}",

					"performance_level": "PL0",

					"disk_name": "createname",

					"instance_id": "alibabacloudstack_instance.ecs.id",

					"payment_type": "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"description": "create",

						"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.3.pre::key.KeyId)}}",

						"performance_level": "PL0",

						"disk_name": "createname",

						"instance_id": "alibabacloudstack_instance.ecs.id",

						"payment_type": "Subscription",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "update",

					"disk_name": "updatename",

					"instance_id": "alibabacloudstack_instance.ecs.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "update",

						"disk_name": "updatename",

						"instance_id": "alibabacloudstack_instance.ecs.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"description": "create",

					"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.3.pre::key.KeyId)}}",

					"performance_level": "PL0",

					"disk_name": "createname",

					"instance_id": "alibabacloudstack_instance.ecs.id",

					"payment_type": "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"description": "create",

						"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.3.pre::key.KeyId)}}",

						"performance_level": "PL0",

						"disk_name": "createname",

						"instance_id": "alibabacloudstack_instance.ecs.id",

						"payment_type": "Subscription",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "update",

					"disk_name": "updatename",

					"instance_id": "alibabacloudstack_instance.ecs.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "update",

						"disk_name": "updatename",

						"instance_id": "alibabacloudstack_instance.ecs.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "update",

					"disk_name": "updatename",

					"instance_id": "alibabacloudstack_instance.ecs.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "update",

						"disk_name": "updatename",

						"instance_id": "alibabacloudstack_instance.ecs.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "update",

					"disk_name": "updatename",

					"instance_id": "alibabacloudstack_instance.ecs.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "update",

						"disk_name": "updatename",

						"instance_id": "alibabacloudstack_instance.ecs.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "update",

					"disk_name": "updatename",

					"instance_id": "alibabacloudstack_instance.ecs.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "update",

						"disk_name": "updatename",

						"instance_id": "alibabacloudstack_instance.ecs.id",
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
func TestAccAlibabacloudStackEcsDisk12(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"description": "test",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"disk_name": "testdisk",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"description": "test",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"disk_name": "testdisk",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "eeeeee",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "eeeeee",
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
func TestAccAlibabacloudStackEcsDisk13(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_auto",

					"description": "挂盘测试",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"disk_name": "挂载测试",

					"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_auto",

						"description": "挂盘测试",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"disk_name": "挂载测试",

						"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
func TestAccAlibabacloudStackEcsDisk14(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"description": "加密测试",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"disk_name": "zone测试",

					"multi_attach": "Disabled",

					"performance_level": "PL2",

					"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.47.pre::key.KeyId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"description": "加密测试",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"disk_name": "zone测试",

						"multi_attach": "Disabled",

						"performance_level": "PL2",

						"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.47.pre::key.KeyId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "zone测试",

					"category": "cloud_auto",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "zone测试",

						"category": "cloud_auto",
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
func TestAccAlibabacloudStackEcsDisk15(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"description": "加密测试",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"disk_name": "zone测试",

					"multi_attach": "Disabled",

					"performance_level": "PL2",

					"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.47.pre::key.KeyId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"description": "加密测试",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"disk_name": "zone测试",

						"multi_attach": "Disabled",

						"performance_level": "PL2",

						"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.47.pre::key.KeyId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "zone测试",

					"category": "cloud_auto",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "zone测试",

						"category": "cloud_auto",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"description": "加密测试",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"disk_name": "zone测试",

					"multi_attach": "Disabled",

					"performance_level": "PL2",

					"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.47.pre::key.KeyId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"description": "加密测试",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"disk_name": "zone测试",

						"multi_attach": "Disabled",

						"performance_level": "PL2",

						"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.47.pre::key.KeyId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
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
func TestAccAlibabacloudStackEcsDisk16(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_auto",

					"description": "挂盘测试",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"disk_name": "挂载测试",

					"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_auto",

						"description": "挂盘测试",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"disk_name": "挂载测试",

						"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_auto",

					"description": "挂盘测试",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"disk_name": "挂载测试",

					"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_auto",

						"description": "挂盘测试",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"disk_name": "挂载测试",

						"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
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
func TestAccAlibabacloudStackEcsDisk17(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_auto",

					"description": "挂盘测试",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"disk_name": "挂载测试",

					"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_auto",

						"description": "挂盘测试",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"disk_name": "挂载测试",

						"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"payment_type": "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"payment_type": "Subscription",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_auto",

					"description": "挂盘测试",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"disk_name": "挂载测试",

					"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_auto",

						"description": "挂盘测试",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"disk_name": "挂载测试",

						"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "PayAsYouGo",
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
func TestAccAlibabacloudStackEcsDisk18(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"description": "加密测试",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"disk_name": "zone测试",

					"multi_attach": "Disabled",

					"performance_level": "PL2",

					"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.47.pre::key.KeyId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"description": "加密测试",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"disk_name": "zone测试",

						"multi_attach": "Disabled",

						"performance_level": "PL2",

						"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.47.pre::key.KeyId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"disk_name": "zone测试",

					"category": "cloud_auto",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"disk_name": "zone测试",

						"category": "cloud_auto",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_essd",

					"description": "加密测试",

					"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

					"disk_name": "zone测试",

					"multi_attach": "Disabled",

					"performance_level": "PL2",

					"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.47.pre::key.KeyId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_essd",

						"description": "加密测试",

						"zone_id": "data.alibabacloudstack_zones.default.zones.0.id",

						"disk_name": "zone测试",

						"multi_attach": "Disabled",

						"performance_level": "PL2",

						"kms_key_id": "${{ref(resource, KMS::Key::5.0.0.47.pre::key.KeyId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "zoneIdUpdate",

					"performance_level": "PL1",

					"disk_name": "zone测试update",

					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "zoneIdUpdate",

						"performance_level": "PL1",

						"disk_name": "zone测试update",

						"category": "cloud_essd",
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
func TestAccAlibabacloudStackEcsDisk19(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_auto",

					"description": "挂盘测试",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"disk_name": "挂载测试",

					"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_auto",

						"description": "挂盘测试",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"disk_name": "挂载测试",

						"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",

					"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",

						"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_auto",

					"description": "挂盘测试",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"disk_name": "挂载测试",

					"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_auto",

						"description": "挂盘测试",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"disk_name": "挂载测试",

						"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",

					"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",

						"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",

					"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",

						"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",

					"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",

						"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",

					"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",

						"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",

					"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",

						"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",

					"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",

						"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",
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
func TestAccAlibabacloudStackEcsDisk20(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDiskCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedisksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdisk%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDiskBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_auto",

					"description": "挂盘测试",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"disk_name": "挂载测试",

					"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_auto",

						"description": "挂盘测试",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"disk_name": "挂载测试",

						"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",

					"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",

					"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",

						"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",

						"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"category": "cloud_auto",

					"description": "挂盘测试",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"disk_name": "挂载测试",

					"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"category": "cloud_auto",

						"description": "挂盘测试",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"disk_name": "挂载测试",

						"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",

					"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",

					"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",

						"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",

						"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",

					"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",

					"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",

						"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",

						"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",

					"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",

					"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",

						"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",

						"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",

					"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",

					"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",

						"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",

						"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",

					"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",

					"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",

						"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",

						"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",

					"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",

					"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",

						"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",

						"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "挂盘",

					"disk_name": "updateDiskName",

					"category": "cloud_auto",

					"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

					"payment_type": "Subscription",

					"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",

					"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "挂盘",

						"disk_name": "updateDiskName",

						"category": "cloud_auto",

						"instance_id": "${{ref(resource, ECS::Instance::5.1.0::instance.InstanceId)}}",

						"payment_type": "Subscription",

						"image_id": "${{ref(resource, ECS::Image::5.0.0.28.pre::createImage.ImageId)}}",

						"snapshot_id": "${{ref(resource, ECS::Snapshot::4.0.0::createSnapshot.SnapshotId)}}",
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

