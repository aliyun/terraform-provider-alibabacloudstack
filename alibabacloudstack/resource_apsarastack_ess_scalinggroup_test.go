package alibabacloudstack

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"log"
	"testing"

	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
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

	basicMap := map[string]string{
		"min_size":           "1",
		"max_size":           "4",
		"default_cooldown":   "20",
		"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroup-%d", rand),
		"vswitch_ids.#":      "2",
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
				Config: testAccEssScalingGroup(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEssScalingGroupUpdateMaxSize(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_size": "5",
					}),
				),
			},

			{
				Config: testAccEssScalingGroupUpdateScalingGroupName(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					}),
				),
			},
			{
				Config: testAccEssScalingGroupUpdateRemovalPolicies(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"removal_policies.#": "1",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupUpdateDefaultCooldown(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_cooldown": "200",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupUpdateMinSize(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_size": "2",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupModifyVSwitchIds(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccEssScalingGroup(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
				),
			},
		},
	})

}

func TestAccAlibabacloudStackEssScalingGroup_vpc(t *testing.T) {
	rand := getAccTestRandInt(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alibabacloudstack_ess_scaling_group.default"

	basicMap := map[string]string{
		"min_size":           "1",
		"max_size":           "1",
		"default_cooldown":   "20",
		"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroup_vpc-%d", rand),
		"vswitch_ids.#":      "2",
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
				Config: testAccEssScalingGroupVpc(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
				//ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEssScalingGroupVpcUpdateMaxSize(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_size": "2",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupVpcUpdateScalingGroupName(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					}),
				),
			},
			{
				Config: testAccEssScalingGroupVpcUpdateRemovalPolicies(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"removal_policies.#": "1",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupVpcUpdateDefaultCooldown(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_cooldown": "200",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupVpcUpdateMinSize(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_size": "2",
					}),
				),
			},
			{
				Config: testAccEssScalingGroupVpc(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
				),
			},
		},
	})

}

func TestAccAlibabacloudStackEssScalingGroup_slb(t *testing.T) {
	var v ess.ScalingGroup
	var slb *slb.DescribeLoadBalancerAttributeResponse
	rand := getAccTestRandInt(10000, 999999)
	resourceId := "alibabacloudstack_ess_scaling_group.default"

	basicMap := map[string]string{
		"min_size":           "1",
		"max_size":           "1",
		"default_cooldown":   "300",
		"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroup_slb-%d", rand),
		"vswitch_ids.#":      "1",
		"removal_policies.#": "2",
		"loadbalancer_ids.#": "0",
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rcSlb0 := resourceCheckInit("alibabacloudstack_slb.default.0", &slb, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rcSlb1 := resourceCheckInit("alibabacloudstack_slb.default.1", &slb, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
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
				Config: testAccEssScalingGroupSlbempty(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
				//ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccEssScalingGroupSlb(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					rcSlb1.checkResourceExists(),
					testAccCheck(map[string]string{
						"loadbalancer_ids.#": "2",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccEssScalingGroupSlbDetach(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					testAccCheck(map[string]string{
						"loadbalancer_ids.#": "1",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccEssScalingGroupSlbUpdateMaxSize(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					rcSlb1.checkResourceExists(),
					testAccCheck(map[string]string{
						"max_size":           "2",
						"loadbalancer_ids.#": "2",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccEssScalingGroupSlbUpdateScalingGroupName(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					rcSlb1.checkResourceExists(),
					testAccCheck(map[string]string{
						"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroupUpdate-%d", rand),
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccEssScalingGroupSlbUpdateRemovalPolicies(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					rcSlb1.checkResourceExists(),
					testAccCheck(map[string]string{
						"removal_policies.#": "1",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccEssScalingGroupSlbUpdateDefaultCooldown(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					rcSlb1.checkResourceExists(),
					testAccCheck(map[string]string{
						"default_cooldown": "200",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccEssScalingGroupSlbUpdateMinSize(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					rcSlb0.checkResourceExists(),
					rcSlb1.checkResourceExists(),
					testAccCheck(map[string]string{
						"min_size": "2",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccEssScalingGroupSlbempty(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"loadbalancer_ids.#": "0",
						"min_size":           "1",
						"max_size":           "1",
						"default_cooldown":   "300",
						"removal_policies.#": "2",
						"scaling_group_name": fmt.Sprintf("tf-testAccEssScalingGroup_slb-%d", rand),
					}),
				),
				//ExpectNonEmptyPlan: true,
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

func testAccEssScalingGroup(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup-%d"
	}
	
	resource "alibabacloudstack_vswitch" "default2" {
		  vpc_id = "${alibabacloudstack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 4
		scaling_group_name = "${var.name}"
		default_cooldown = 20
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}", "${alibabacloudstack_vswitch.default2.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}`, common, rand)
}

func testAccEssScalingGroupUpdateMaxSize(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup-%d"
	}
	
	resource "alibabacloudstack_vswitch" "default2" {
		  vpc_id = "${alibabacloudstack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 5
		scaling_group_name = "${var.name}"
		default_cooldown = 20
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}", "${alibabacloudstack_vswitch.default2.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}`, common, rand)
}

func testAccEssScalingGroupUpdateScalingGroupName(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alibabacloudstack_vswitch" "default2" {
		  vpc_id = "${alibabacloudstack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 5
		scaling_group_name = "${var.name}"
		default_cooldown = 20
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}", "${alibabacloudstack_vswitch.default2.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}`, common, rand)
}

func testAccEssScalingGroupUpdateRemovalPolicies(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alibabacloudstack_vswitch" "default2" {
		  vpc_id = "${alibabacloudstack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 5
		scaling_group_name = "${var.name}"
		default_cooldown = 20
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}", "${alibabacloudstack_vswitch.default2.id}"]
		removal_policies = ["OldestInstance"]
	}`, common, rand)
}

func testAccEssScalingGroupUpdateDefaultCooldown(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alibabacloudstack_vswitch" "default2" {
		  vpc_id = "${alibabacloudstack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 5
		scaling_group_name = "${var.name}"
		default_cooldown = 200
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}", "${alibabacloudstack_vswitch.default2.id}"]
		removal_policies = ["OldestInstance"]
	}`, common, rand)
}

func testAccEssScalingGroupUpdateMinSize(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alibabacloudstack_vswitch" "default2" {
		  vpc_id = "${alibabacloudstack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 2
		max_size = 5
		scaling_group_name = "${var.name}"
		default_cooldown = 200
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}", "${alibabacloudstack_vswitch.default2.id}"]
		removal_policies = ["OldestInstance"]
	}`, common, rand)
}
func testAccEssScalingGroupVpc(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup_vpc-%d"
	}
	
	resource "alibabacloudstack_vswitch" "default2" {
		  vpc_id = "${alibabacloudstack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		default_cooldown = 20
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}", "${alibabacloudstack_vswitch.default2.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}`, common, rand)
}

func testAccEssScalingGroupVpcUpdateMaxSize(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup_vpc-%d"
	}
	
	resource "alibabacloudstack_vswitch" "default2" {
		  vpc_id = "${alibabacloudstack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 2
		scaling_group_name = "${var.name}"
		default_cooldown = 20
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}", "${alibabacloudstack_vswitch.default2.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}`, common, rand)
}

func testAccEssScalingGroupVpcUpdateScalingGroupName(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alibabacloudstack_vswitch" "default2" {
		  vpc_id = "${alibabacloudstack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 2
		scaling_group_name = "${var.name}"
		default_cooldown = 20
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}", "${alibabacloudstack_vswitch.default2.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}`, common, rand)
}

func testAccEssScalingGroupVpcUpdateRemovalPolicies(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alibabacloudstack_vswitch" "default2" {
		  vpc_id = "${alibabacloudstack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 2
		scaling_group_name = "${var.name}"
		default_cooldown = 20
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}", "${alibabacloudstack_vswitch.default2.id}"]
		removal_policies = ["OldestInstance"]
	}`, common, rand)
}

func testAccEssScalingGroupVpcUpdateDefaultCooldown(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alibabacloudstack_vswitch" "default2" {
		  vpc_id = "${alibabacloudstack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 2
		scaling_group_name = "${var.name}"
		default_cooldown = 200
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}", "${alibabacloudstack_vswitch.default2.id}"]
		removal_policies = ["OldestInstance"]
	}`, common, rand)
}

func testAccEssScalingGroupVpcUpdateMinSize(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alibabacloudstack_vswitch" "default2" {
		  vpc_id = "${alibabacloudstack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 2
		max_size = 2
		scaling_group_name = "${var.name}"
		default_cooldown = 200
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}", "${alibabacloudstack_vswitch.default2.id}"]
		removal_policies = ["OldestInstance"]
	}`, common, rand)
}

func testAccEssScalingGroupSlb(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup_slb-%d"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
	  min_size = "1"
	  max_size = "1"
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance", "NewestInstance"]
	  vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
	  loadbalancer_ids = ["${alibabacloudstack_slb.default.0.id}","${alibabacloudstack_slb.default.1.id}"]
	  depends_on = ["alibabacloudstack_slb_listener.default"]
	}

	resource "alibabacloudstack_slb" "default" {
	  count=2
	  name = "${var.name}"
	  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
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
	`, common, rand)
}

func testAccEssScalingGroupSlbDetach(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup_slb-%d"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
	  min_size = "1"
	  max_size = "1"
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance", "NewestInstance"]
	  vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
	  loadbalancer_ids = ["${alibabacloudstack_slb.default.0.id}"]
	  depends_on = ["alibabacloudstack_slb_listener.default"]
	}

	resource "alibabacloudstack_slb" "default" {
	  count=2
	  name = "${var.name}"
	  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
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
	`, common, rand)
}

func testAccEssScalingGroupSlbempty(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup_slb-%d"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
	  min_size = "1"
	  max_size = "1"
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance", "NewestInstance"]
	  vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
	  loadbalancer_ids = []
	}`, common, rand)
}

func testAccEssScalingGroupSlbUpdateMaxSize(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroup_slb-%d"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
	  min_size = "1"
	  max_size = "2"
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance", "NewestInstance"]
	  vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
	  loadbalancer_ids = ["${alibabacloudstack_slb.default.0.id}","${alibabacloudstack_slb.default.1.id}"]
	  depends_on = ["alibabacloudstack_slb_listener.default"]
	}

	resource "alibabacloudstack_slb" "default" {
	  count=2
	  name = "${var.name}"
	  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
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
	`, common, rand)
}

func testAccEssScalingGroupSlbUpdateScalingGroupName(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
	  min_size = "1"
	  max_size = "2"
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance", "NewestInstance"]
	  vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
	  loadbalancer_ids = ["${alibabacloudstack_slb.default.0.id}","${alibabacloudstack_slb.default.1.id}"]
	  depends_on = ["alibabacloudstack_slb_listener.default"]
	}

	resource "alibabacloudstack_slb" "default" {
	  count=2
	  name = "${var.name}"
	  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	}

	resource "alibabacloudstack_slb_listener" "default" {
	  count = 2
	  load_balancer_id = "${element(alibabacloudstack_slb.default.*.id, count.index)}"
	  backend_port = "22"
	  frontend_port = "22"
	  protocol = "http"
	  bandwidth = "10"
	  health_check_type = "http"
       sticky_session="off"
	  health_check= "off"
	}
	`, common, rand)
}

func testAccEssScalingGroupSlbUpdateRemovalPolicies(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
	  min_size = "1"
	  max_size = "2"
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance"]
	  vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
	  loadbalancer_ids = ["${alibabacloudstack_slb.default.0.id}","${alibabacloudstack_slb.default.1.id}"]
	  depends_on = ["alibabacloudstack_slb_listener.default"]
	}

	resource "alibabacloudstack_slb" "default" {
	  count=2
	  name = "${var.name}"
	  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
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
	`, common, rand)
}

func testAccEssScalingGroupSlbUpdateDefaultCooldown(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
	  min_size = "1"
	  max_size = "2"
      default_cooldown = 200
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance"]
	  vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
	  loadbalancer_ids = ["${alibabacloudstack_slb.default.0.id}","${alibabacloudstack_slb.default.1.id}"]
	  depends_on = ["alibabacloudstack_slb_listener.default"]
	}

	resource "alibabacloudstack_slb" "default" {
	  count=2
	  name = "${var.name}"
	  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
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
	`, common, rand)
}

func testAccEssScalingGroupSlbUpdateMinSize(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
	  min_size = "2"
	  max_size = "2"
      default_cooldown = 200
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance"]
	  vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
	  loadbalancer_ids = ["${alibabacloudstack_slb.default.0.id}","${alibabacloudstack_slb.default.1.id}"]
	  depends_on = ["alibabacloudstack_slb_listener.default"]
	}

	resource "alibabacloudstack_slb" "default" {
	  count=2
	  name = "${var.name}"
	  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
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
	`, common, rand)
}

func testAccEssScalingGroupModifyVSwitchIds(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingGroupUpdate-%d"
	}
	
	resource "alibabacloudstack_vswitch" "default2" {
		  vpc_id = "${alibabacloudstack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		  name = "${var.name}-bar"
	}

	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 2
		max_size = 5
		scaling_group_name = "${var.name}"
		default_cooldown = 200
		vswitch_ids = ["${alibabacloudstack_vswitch.default2.id}"]
		removal_policies = ["OldestInstance"]
	}`, common, rand)
}
