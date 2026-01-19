package alibabacloudstack

import (
	"fmt"
	"log"
	"testing"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_security_group", &resource.Sweeper{
		Name: "alibabacloudstack_security_group",
		F:    testSweepSecurityGroups,
		//When implemented, these should be removed firstly
		Dependencies: []string{
			"alibabacloudstack_ecs_instance",
			"alibabacloudstack_network_interface",
			"alibabacloudstack_yundun_bastionhost_instance",
		},
	})
}

func testSweepSecurityGroups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var groups []ecs.SecurityGroup
	req := ecs.CreateDescribeSecurityGroupsRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		req.Scheme = "https"
	} else {
		req.Scheme = "http"
	}
	req.Headers = map[string]string{"RegionId": client.RegionId}
	req.QueryParams = map[string]string{"Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeSecurityGroups(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Security Groups: %s", err)
		}
		resp, _ := raw.(*ecs.DescribeSecurityGroupsResponse)
		if resp == nil || len(resp.SecurityGroups.SecurityGroup) < 1 {
			break
		}
		groups = append(groups, resp.SecurityGroups.SecurityGroup...)

		if len(resp.SecurityGroups.SecurityGroup) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	vpcService := VpcService{client}
	ecsService := EcsService{client}
	for _, v := range groups {
		name := v.SecurityGroupName
		id := v.SecurityGroupId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a Security Group created by other service, it should be fetched by vpc name and deleted.
		if skip {
			if need, err := vpcService.needSweepVpc(v.VpcId, ""); err == nil {
				skip = !need
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Security Group: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Security Group: %s (%s)", name, id)
		if err := ecsService.sweepSecurityGroup(id); err != nil {
			log.Printf("[ERROR] Failed to delete Security Group (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func testAccCheckSecurityGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_security_group" {
			continue
		}

		_, err := ecsService.DescribeSecurityGroup(rs.Primary.ID)

		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return err
		}
		return errmsgs.WrapError(errmsgs.Error("Error SecurityGroup still exist"))
	}
	return nil
}

func TestAccAlibabacloudStackEcsSecurityGroupBasic(t *testing.T) {
	var v ecs.DescribeSecurityGroupAttributeResponse
	resourceId := "alibabacloudstack_security_group.default"
	ra := resourceAttrInit(resourceId, testAccCheckSecurityBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc_sg_%d", rand)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSecurityGroupConfigBasic(name),
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
				Config: testAccCheckSecurityGroupConfigName(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                fmt.Sprintf("%s_change", name),
						"inner_access_policy": "Drop",
						"description":         fmt.Sprintf("%s_change", name),
					}),
				),
			},

			{
				Config: testAccCheckSecurityGroupConfigAll(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(testAccCheckSecurityBasicMap),
				),
			},
			{
				Config: testAccCheckSecurityGroupConfigTags(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccCheckSecurityGroupConfigTagsUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccCheckSecurityGroupConfigAll(name),
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

func TestAccAlibabacloudStackEcsSecurityGroupMulti(t *testing.T) {
	var v ecs.DescribeSecurityGroupAttributeResponse
	resourceId := "alibabacloudstack_security_group.default.2"
	ra := resourceAttrInit(resourceId, testAccCheckSecurityBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc_sg_%d", rand)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSecurityGroupConfigMulti(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccCheckSecurityGroupConfigBasic(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

resource "alibabacloudstack_security_group" "default" {
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  name = "${var.name}"
  description = "${var.name}_describe"
  type = "normal"
  
}
`, name, VpcCommonTestCase)
}

func testAccCheckSecurityGroupConfigName(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}


%s

resource "alibabacloudstack_security_group" "default" {
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  name = "${var.name}_change"
  description = "${var.name}_change"
  type = "normal"
  inner_access_policy = "Drop"
  
}`, name, VpcCommonTestCase)
}

func testAccCheckSecurityGroupConfigAll(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

resource "alibabacloudstack_security_group" "default" {
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  inner_access_policy = "Drop"
  name = "${var.name}"
  description = "${var.name}_describe"
  type = "normal"
}`, name, VpcCommonTestCase)
}

func testAccCheckSecurityGroupConfigMulti(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}


%s

resource "alibabacloudstack_security_group" "default" {
  count = 3
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  name = "${var.name}"
  description = "${var.name}_describe"
  type = "normal"
}`, name, VpcCommonTestCase)
}

func testAccCheckSecurityGroupConfigTags(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

%s

resource "alibabacloudstack_security_group" "default" {
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  inner_access_policy = "Drop"
  name = "${var.name}"
  description = "${var.name}_describe"
  type = "normal"
  tags = {
    "Created" = "TF"
    "For" = "Test"
  }
}`, name, VpcCommonTestCase)
}

func testAccCheckSecurityGroupConfigTagsUpdate(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

%s

resource "alibabacloudstack_security_group" "default" {
   vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  inner_access_policy = "Drop"
  name = "${var.name}"
  description = "${var.name}_describe"
  type = "normal"
  tags = {
    "Created" = "TF-update"
    "For" = "Test-update"
  }
}`, name, VpcCommonTestCase)
}

var testAccCheckSecurityBasicMap = map[string]string{
	"vpc_id":              CHECKSET,
	"inner_access_policy": CHECKSET,
	"name":                CHECKSET,
	"description":         CHECKSET,
	"type":                "normal",
	//"tags.%":              "2",
	//"tags.foo":            "foo",
	//"tags.Test":           "Test",
}
