package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEssScalingConfigurationUpdate(t *testing.T) {
	rand := getAccTestRandInt(1000, 999999)
	var v ess.ScalingConfigurationInDescribeScalingConfigurations
	resourceId := "alibabacloudstack_ess_scaling_configuration.default"
	basicMap := map[string]string{
		"scaling_group_id":     CHECKSET,
		"instance_type":        CHECKSET,
		"security_group_ids.#": "1",
		"image_id":             CHECKSET,
		"override":             "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScCon-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":   "${alibabacloudstack_ess_scaling_group.default.id}",
					"image_id":           "${data.alibabacloudstack_images.default.images.0.id}",
					"instance_type":      "ecs.n4.large",
					"security_group_ids": []string{"${alibabacloudstack_ecs_securitygroup.default.id}"},
					"deployment_set_id":  "${alibabacloudstack_ecs_deployment_set.default.id}",
					"force_delete":       "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName: resourceId,
				ImportState:  true,

				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "security_group_id", "password", "kms_encrypted_password", "kms_encryption_context"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"active": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"active": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_configuration_name": fmt.Sprintf("tf-testAccEssScCon-%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_configuration_name": fmt.Sprintf("tf-testAccEssScCon-%d", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_category": "cloud_ssd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_category": "cloud_ssd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disk": []map[string]string{{
						"size":                 "20",
						"category":             "cloud_ssd",
						"delete_with_instance": "false",
						"encrypted":            "true",
					},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disk.#":                      "1",
						"data_disk.0.size":                 "20",
						"data_disk.0.category":             "cloud_ssd",
						"data_disk.0.delete_with_instance": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_data": `#!/bin/bash\necho \"hello\"\n`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_data": "#!/bin/bash\necho \"hello\"\n",
					}),
				),
				//ExpectNonEmptyPlan: true,
			},

			//			{
			//				Config: testAccConfig(map[string]interface{}{
			//					"key_name": "${alibabacloudstack_key_pair.default.id}",
			//				}),
			//				Check: resource.ComposeTestCheckFunc(
			//					testAccCheck(map[string]string{
			//						"key_name": CHECKSET,
			//					}),
			//				),
			//			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"name": "tf-test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.name": "tf-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackEssScalingConfigurationMulti(t *testing.T) {
	rand := getAccTestRandInt(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alibabacloudstack_ess_scaling_configuration.default.0"
	basicMap := map[string]string{
		"scaling_group_id":     CHECKSET,
		"instance_type":        CHECKSET,
		"security_group_ids.#": "1",
		"image_id":             CHECKSET,
		"override":             "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScCon-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":              "1",
					"scaling_group_id":   "${alibabacloudstack_ess_scaling_group.default.id}",
					"image_id":           "${data.alibabacloudstack_images.default.images.0.id}",
					"instance_type":      "ecs.n4.large",
					"security_group_ids": []string{"${alibabacloudstack_ecs_securitygroup.default.id}"},
					"deployment_set_id":  "${alibabacloudstack_ecs_deployment_set.default.id}",
					"force_delete":       "true",
					"data_disk": []map[string]string{{
						"size":                 "20",
						"category":             "cloud_ssd",
						"delete_with_instance": "false",
						"encrypted":            "true",
						"kms_key_id":           "149ca9b2-564d-42f7-ab60-abfd15a91503",
					},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceEssScalingConfigurationConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "%s"
	}

	resource "alibabacloudstack_ecs_deployment_set" "default" {
		strategy            = "Availability"
		domain              = "Default"
		granularity         = "Host"
		deployment_set_name = "example_value"
		description         = "example_value"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}"]
	}`, ECSInstanceCommonTestCase, name)
}
