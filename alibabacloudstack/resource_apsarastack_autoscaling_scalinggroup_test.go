package alibabacloudstack

import (
	"fmt"
	"log"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_ess_scalinggroup", &resource.Sweeper{
		Name: "alibabacloudstack_ess_scalinggroup",
		F:    testSweepEssGroups,
	})
}

func testSweepEssGroups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alibabacloudstack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var groups []ess.ScalingGroup
	req := ess.CreateDescribeScalingGroupsRequest()

	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeScalingGroups(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Scaling groups: %s", err)
		}
		resp, _ := raw.(*ess.DescribeScalingGroupsResponse)
		if resp == nil || len(resp.ScalingGroups.ScalingGroup) < 1 {
			break
		}
		groups = append(groups, resp.ScalingGroups.ScalingGroup...)

		if len(resp.ScalingGroups.ScalingGroup) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	sweeped := false
	for _, v := range groups {
		name := v.ScalingGroupName
		id := v.ScalingGroupId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Scaling Group: %s (%s)", name, id)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting Scaling Group: %s (%s)", name, id)
		req := ess.CreateDeleteScalingGroupRequest()
		req.ScalingGroupId = id
		req.ForceDelete = requests.NewBoolean(true)
		_, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DeleteScalingGroup(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Scaling Group (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		time.Sleep(2 * time.Minute)
	}
	return nil
}

func TestAccAlibabacloudStackEssScalingGroup_basic(t *testing.T) {
	rand := getAccTestRandInt(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alibabacloudstack_ess_scaling_group.default"

	name := fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScalingGroup)

	basicMap := map[string]string{
		"min_size":           "1",
		"max_size":           "4",
		"default_cooldown":   "20",
		"vswitch_ids.#":      "1",
		"removal_policies.#": "2",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           1,
					"max_size":           4,
					"scaling_group_name": name,
					"default_cooldown":   20,
					"vswitch_ids":        []string{"${alibabacloudstack_vpc_vswitch.default.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_name": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_size": 5,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_size": "5",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"removal_policies": []string{"OldestInstance"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"removal_policies.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_cooldown": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_cooldown": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_size": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_ids": []string{"${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.update.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_ids": []string{"${alibabacloudstack_vpc_vswitch.update.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_ids.#": "1",
					}),
				),
			},
		},
	})

}

func TestAccAlibabacloudStackEssScalingGroup_slb(t *testing.T) {
	rand := getAccTestRandInt(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alibabacloudstack_ess_scaling_group.default"

	name := fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScalingGroupWithSlb)

	basicMap := map[string]string{
		"min_size":           "1",
		"max_size":           "4",
		"default_cooldown":   "20",
		"scaling_group_name": name,
		"vswitch_ids.#":      "1",
		"removal_policies.#": "2",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           1,
					"max_size":           4,
					"scaling_group_name": name,
					"default_cooldown":   20,
					"vswitch_ids":        []string{"${alibabacloudstack_vpc_vswitch.default.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"loadbalancer_ids":   []string{"${alibabacloudstack_slb.default.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_name": name,
						"loadbalancer_ids.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"loadbalancer_ids": []string{"${alibabacloudstack_slb.default.0.id}", "${alibabacloudstack_slb.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"loadbalancer_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"loadbalancer_ids": []string{"${alibabacloudstack_slb.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"loadbalancer_ids.#": "1",
					}),
				),
			},
		},
	})

}

func TestAccAlibabacloudStackEssScalingGroup_rds(t *testing.T) {
	rand := getAccTestRandInt(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alibabacloudstack_ess_scaling_group.default"

	name := fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScalingGroupWithRds)

	basicMap := map[string]string{
		"min_size":           "1",
		"max_size":           "4",
		"default_cooldown":   "20",
		"vswitch_ids.#":      "1",
		"removal_policies.#": "2",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"min_size":           1,
					"max_size":           4,
					"scaling_group_name": name,
					"default_cooldown":   20,
					"vswitch_ids":        []string{"${alibabacloudstack_vpc_vswitch.default.id}"},
					"removal_policies":   []string{"OldestInstance", "NewestInstance"},
					"db_instance_ids":   []string{"${alibabacloudstack_db_instance.default.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_name": name,
						"db_instance_ids.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_ids": []string{"${alibabacloudstack_db_instance.default.0.id}", "${alibabacloudstack_db_instance.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_ids": []string{"${alibabacloudstack_db_instance.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_ids.#": "1",
					}),
				),
			},
		},
	})

}

func testAccCheckEssScalingGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_ess_scaling_group" {
			continue
		}

		if _, err := essService.DescribeEssScalingGroup(rs.Primary.ID); err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
		return errmsgs.WrapError(fmt.Errorf("Scaling group %s still exists.", rs.Primary.ID))
	}

	return nil
}

func testAccEssScalingGroup(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
	%s
	
	resource "alibabacloudstack_vpc_vswitch" "update" {
	  vswitch_name = "${var.name}_vsw2"
	  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	  cidr_block = "172.16.2.0/24"
	  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	  lifecycle {
	      ignore_changes = [
	        tags
	      ]
	  }
	}
	`, name, VSwitchCommonTestCase)
}
func testAccEssScalingGroupWithSlb(name string) string {
	return fmt.Sprintf(`
		%s
	resource "alibabacloudstack_slb" "default" {
	  count=2
	  name = "${var.name}"
	  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	}

	resource "alibabacloudstack_slb_listener" "default" {
	  count = 2
	  load_balancer_id = "${element(alibabacloudstack_slb.default.*.id, count.index)}"
	  backend_port = "22"
	  frontend_port = "22"
	  protocol = "http"
	  bandwidth = "10"
	  health_check_type = "http"
	  health_check ="off"
	  sticky_session ="off"
	}

`, testAccEssScalingGroup(name))
}

func testAccEssScalingGroupWithRds(name string) string {
	return fmt.Sprintf(`
%s
data "alibabacloudstack_rds_instance_types" "default" {
  engine               = "MySQL"
  sorted_by            = "CPU"
}

resource "alibabacloudstack_db_instance" "default" {
  count = 2
  engine               = data.alibabacloudstack_rds_instance_types.default.instance_types.0.engine
  engine_version       = data.alibabacloudstack_rds_instance_types.default.instance_types.0.engine_version
  instance_type        = data.alibabacloudstack_rds_instance_types.default.instance_types.0.id
  instance_storage     = data.alibabacloudstack_rds_instance_types.default.instance_types.0.storage_min
  instance_name        = "${var.name}"
  vswitch_id           = "${alibabacloudstack_vpc_vswitch.default.id}"
  monitoring_period    = "60"
  storage_type         = data.alibabacloudstack_rds_instance_types.default.instance_types.0.storage_type
}
`, testAccEssScalingGroup(name))
}
