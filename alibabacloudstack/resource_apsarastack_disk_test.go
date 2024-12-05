package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_disk", &resource.Sweeper{
		Name: "alibabacloudstack_disk",
		F:    testSweepDisks,
	})
}

func testSweepDisks(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
	}

	var disks []ecs.Disk
	req := ecs.CreateDescribeDisksRequest()
	req.RegionId = client.RegionId
	req.Headers = map[string]string{"RegionId": client.RegionId}
	req.QueryParams = map[string]string{ "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	req.QueryParams["Department"] = client.Department
	req.QueryParams["ResourceGroup"] = client.ResourceGroup
	if strings.ToLower(client.Config.Protocol) == "https" {
		req.Scheme = "https"
	} else {
		req.Scheme = "http"
	}
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeDisks(req)
		})
		if err != nil {
			return WrapError(err)
		}
		resp, _ := raw.(*ecs.DescribeDisksResponse)
		if resp == nil || len(resp.Disks.Disk) < 1 {
			break
		}
		disks = append(disks, resp.Disks.Disk...)

		if len(resp.Disks.Disk) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	for _, v := range disks {
		name := v.DiskName
		id := v.DiskId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Disk: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Disk: %s (%s)", name, id)
		req := ecs.CreateDeleteDiskRequest()
		req.Headers = map[string]string{"RegionId": client.RegionId}
		req.QueryParams = map[string]string{ "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		req.QueryParams["Department"] = client.Department
		req.QueryParams["ResourceGroup"] = client.ResourceGroup
		req.DiskId = id
		if strings.ToLower(client.Config.Protocol) == "https" {
			req.Scheme = "https"
		} else {
			req.Scheme = "http"
		}
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteDisk(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Disk (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func testAccCheckDiskDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_disk" {
			continue
		}

		// Try to find the Disk
		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		ecsService := EcsService{client}

		_, err := ecsService.DescribeDisk(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

func TestAccAlibabacloudStackDisk_basic(t *testing.T) {
	var v ecs.Disk
	resourceId := "alibabacloudstack_disk.default"
	serverFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serverFunc)
	ra := resourceAttrInit(resourceId, testAccCheckResourceDiskBasicMap)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alibabacloudstack_disk.default",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDiskConfig_basic(),
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
				Config: testAccDiskConfig_size(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size": "70",
					}),
				),
			},
			{
				Config: testAccDiskConfig_name(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccDiskConfig",
					}),
				),
			},
			{
				Config: testAccDiskConfig_description(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testAccDiskConfig_description",
					}),
				),
			},
			{
				Config: testAccDiskConfig_delete_auto_snapshot(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_auto_snapshot": "true",
					}),
				),
			},
			{
				Config: testAccDiskConfig_delete_with_instance(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_with_instance": "true",
					}),
				),
			},
			{
				Config: testAccDiskConfig_enable_auto_snapshot(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_auto_snapshot": "true",
					}),
				),
			},
			{
				Config: testAccDiskConfig_all(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":               "3",
						"tags.name1":           "name1",
						"tags.name2":           "name2",
						"tags.name3":           "name3",
						"name":                 "tf-testAccDiskConfig_all",
						"description":          "nothing",
						"delete_auto_snapshot": "false",
						"delete_with_instance": "false",
						"enable_auto_snapshot": "false",
					}),
				),
			},
		},
	})

}

func TestAccAlibabacloudStackDisk_multi(t *testing.T) {
	var v ecs.Disk
	resourceId := "alibabacloudstack_disk.default.4"
	serverFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serverFunc)
	ra := resourceAttrInit(resourceId, testAccCheckResourceDiskBasicMap)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alibabacloudstack_disk.default.4",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDiskConfig_multi(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        "tf-testAccDiskConfig_multi",
						"description": "nothing",
					}),
				),
			},
		},
	})

}

func testAccDiskConfig_basic() string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alibabacloudstack_disk" "default" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  	size = "50"
}
`)
}

func testAccDiskConfig_size() string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}


resource "alibabacloudstack_disk" "default" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  	size = "70"
}
`)
}

func testAccDiskConfig_name() string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}


variable "name" {
	default = "tf-testAccDiskConfig"
}

resource "alibabacloudstack_disk" "default" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  	size = "70"
	name = "${var.name}"
}
`)
}

func testAccDiskConfig_description() string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}


variable "name" {
	default = "tf-testAccDiskConfig"
}

resource "alibabacloudstack_disk" "default" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  	size = "70"
	name = "${var.name}"
	description = "${var.name}_description"
}
`)
}

func testAccDiskConfig_delete_auto_snapshot() string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}


variable "name" {
	default = "tf-testAccDiskConfig"
}



resource "alibabacloudstack_disk" "default" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  	size = "70"
	name = "${var.name}"
	description = "${var.name}_description"
	category = "cloud_efficiency"
	encrypted = "false"
	delete_auto_snapshot = "true"
}
`)
}

func testAccDiskConfig_delete_with_instance() string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}


variable "name" {
	default = "tf-testAccDiskConfig"
}

resource "alibabacloudstack_disk" "default" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  	size = "70"
	name = "${var.name}"
	description = "${var.name}_description"
	category = "cloud_efficiency"
	encrypted = "false"
	tags = {
		name1 = "name1"
		Name2 = "Name2"
		name3 = "name3"
			}
	delete_auto_snapshot = "true"
	delete_with_instance = "true"
}
`)
}

func testAccDiskConfig_enable_auto_snapshot() string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}


variable "name" {
	default = "tf-testAccDiskConfig"
}

resource "alibabacloudstack_disk" "default" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  	size = "70"
	name = "${var.name}"
	description = "${var.name}_description"
	category = "cloud_efficiency"
	encrypted = "false"
	tags = {
		name1 = "name1"
		Name2 = "Name2"
		name3 = "name3"
			}
	delete_auto_snapshot = "true"
	delete_with_instance = "true"
	enable_auto_snapshot = "true"
}
`)
}

func testAccDiskConfig_all() string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}


variable "name" {
	default = "tf-testAccDiskConfig_all"
}

resource "alibabacloudstack_disk" "default" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  	size = "70"
	name = "${var.name}"
	description = "nothing"
	category = "cloud_efficiency"
	encrypted = "false"
	tags = {
		name1 = "name1"
		name2 = "name2"
		name3 = "name3"
			}
	delete_auto_snapshot = "false"
	delete_with_instance = "false"
	enable_auto_snapshot = "false"
}
`)
}

func testAccDiskConfig_multi() string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}

variable "name" {
	default = "tf-testAccDiskConfig"
}

resource "alibabacloudstack_disk" "default" {
	count = "5"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  	size = "50"
	name = "${var.name}_multi"
	description = "nothing"
	category = "cloud_efficiency"
	encrypted = "false"
}
`)
}

var testAccCheckResourceDiskBasicMap = map[string]string{
	"availability_zone":    CHECKSET,
	"size":                 "50",
	"name":                 "",
	"description":          "",
	"category":             "cloud_efficiency",
	"snapshot_id":          "",
	"encrypted":            "false",
	"tags":                 NOSET,
	"status":               string(Available),
	"delete_auto_snapshot": "false",
	"delete_with_instance": "false",
	"enable_auto_snapshot": "false",
}
