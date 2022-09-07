package alibabacloudstack

import (
	"fmt"
	"log"
	"testing"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_security_group", &resource.Sweeper{
		Name: "alibabacloudstack_security_group",
		F:    testSweepSecurityGroups,
		//When implemented, these should be removed firstly
		Dependencies: []string{
			"alibabacloudstack_instance",
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
	req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
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
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return WrapError(Error("Error SecurityGroup still exist"))
	}
	return nil
}

func TestAccAlibabacloudStackSecurityGroupBasic(t *testing.T) {
	var v ecs.DescribeSecurityGroupAttributeResponse
	resourceId := "alibabacloudstack_security_group.default"
	ra := resourceAttrInit(resourceId, testAccCheckSecurityBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSecurityGroupConfigBasic(),
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
				Config: testAccCheckSecurityGroupConfigInnerAccess(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"inner_access_policy": "Accept",
					}),
				),
			},
			{
				Config: testAccCheckSecurityGroupConfigName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccCheckSecurityGroupName_change",
					}),
				),
			},

			{
				Config: testAccCheckSecurityGroupConfigDescribe(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testAccCheckSecurityGroupName_describe_change",
					}),
				),
			},
			//{
			//	Config: testAccCheckSecurityGroupConfigTags(),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"tags.%":    "1",
			//			"tags.Test": REMOVEKEY,
			//		}),
			//	),
			//},

			{
				Config: testAccCheckSecurityGroupConfigAll(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(testAccCheckSecurityBasicMap),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackSecurityGroupMulti(t *testing.T) {
	var v ecs.DescribeSecurityGroupAttributeResponse
	resourceId := "alibabacloudstack_security_group.default.9"
	ra := resourceAttrInit(resourceId, testAccCheckSecurityBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSecurityGroupConfigMulti(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccCheckSecurityGroupConfigBasic() string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccCheckSecurityGroupName"
}


resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}_vpc"
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_security_group" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  name = "${var.name}"
  description = "${var.name}_describe"
  
}
`)
}

func testAccCheckSecurityGroupConfigInnerAccess() string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccCheckSecurityGroupName"
}


resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}_vpc"
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_security_group" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  inner_access_policy = "Accept"
  name = "${var.name}"
  description = "${var.name}_describe"
  
}`)
}

func testAccCheckSecurityGroupConfigName() string {
	return fmt.Sprintf(`

variable "name" {
  default = "tf-testAccCheckSecurityGroupName"
}


resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}_vpc"
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_security_group" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  name = "${var.name}_change"
  description = "${var.name}_describe"
  
}`)
}

func testAccCheckSecurityGroupConfigDescribe() string {
	return fmt.Sprintf(`

variable "name" {
  default = "tf-testAccCheckSecurityGroupName"
}


resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}_vpc"
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_security_group" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  name = "${var.name}_change"
  description = "${var.name}_describe_change"
}`)
}
func testAccCheckSecurityGroupConfigTags() string {
	return fmt.Sprintf(`

variable "name" {
  default = "tf-testAccCheckSecurityGroupName"
}


resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}_vpc"
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_security_group" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  name = "${var.name}_change"
  description = "${var.name}_describe_change"

}`)
}

func testAccCheckSecurityGroupConfigAll() string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccCheckSecurityGroupName"
}


resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}_vpc"
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_security_group" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  inner_access_policy = "Accept"
  name = "${var.name}"
  description = "${var.name}_describe"
}`)
}

func testAccCheckSecurityGroupConfigMulti() string {
	return fmt.Sprintf(`

variable "name" {
  default = "tf-testAccCheckSecurityGroupName"
}


resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}_vpc"
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_security_group" "default" {
  count = 10
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  name = "${var.name}"
  description = "${var.name}_describe"
}`)
}

var testAccCheckSecurityBasicMap = map[string]string{
	"vpc_id":              CHECKSET,
	"inner_access_policy": "Accept",
	"name":                "tf-testAccCheckSecurityGroupName",
	"description":         "tf-testAccCheckSecurityGroupName_describe",
	//"tags.%":              "2",
	//"tags.foo":            "foo",
	//"tags.Test":           "Test",
}
