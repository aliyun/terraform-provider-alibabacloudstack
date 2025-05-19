package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
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
	req.QueryParams = map[string]string{"Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
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
			return errmsgs.WrapError(err)
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
		req.QueryParams = map[string]string{"Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
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
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
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
	ResourceTest(t, resource.TestCase{
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
						"disk_name": "tf-testAccDiskConfig",
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
						"disk_name":            "tf-testAccDiskConfig_all",
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
	resourceId := "alibabacloudstack_disk.default.2"
	serverFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serverFunc)
	ra := resourceAttrInit(resourceId, testAccCheckResourceDiskBasicMap)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alibabacloudstack_disk.default.2",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDiskConfig_multi(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_name":   "tf-testAccDiskConfig_multi",
						"description": "nothing",
					}),
				),
			},
		},
	})

}

func TestAccAlibabacloudStackDisk_Encrypted(t *testing.T) {
	var v ecs.Disk
	resourceId := "alibabacloudstack_disk.default"
	serverFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serverFunc)
	ra := resourceAttrInit(resourceId, map[string]string{})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alibabacloudstack_disk.default",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDiskConfig_encrypted(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      "testAccDiskConfig_encrypted",
						"encrypted": "true",
					}),
				),
			},
		},
	})
}

func testAccDiskConfig_basic() string {
	return fmt.Sprintf(`
%s

resource "alibabacloudstack_disk" "default" {
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  	size = "50"
}
`, DataAlibabacloudstackVswitchZones)
}

func testAccDiskConfig_encrypted() string {
	return fmt.Sprintf(`
%s

resource "alibabacloudstack_disk" "default" {
    name = "testAccDiskConfig_encrypted"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  	size = "50"
	category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
	encrypted = true
	kms_key_id = "${alibabacloudstack_kms_key.key.id}"
}

%s
`, DataZoneCommonTestCase, KeyCommonTestCase)
}

func testAccDiskConfig_size() string {
	return fmt.Sprintf(`
%s


resource "alibabacloudstack_disk" "default" {
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  	size = "70"
}
`, DataAlibabacloudstackVswitchZones)
}

func testAccDiskConfig_name() string {
	return fmt.Sprintf(`
%s


variable "name" {
	default = "tf-testAccDiskConfig"
}

resource "alibabacloudstack_disk" "default" {
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  	size = "70"
	disk_name = "${var.name}"
}
`, DataAlibabacloudstackVswitchZones)
}

func testAccDiskConfig_description() string {
	return fmt.Sprintf(`
%s


variable "name" {
	default = "tf-testAccDiskConfig"
}

resource "alibabacloudstack_disk" "default" {
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  	size = "70"
	disk_name = "${var.name}"
	description = "${var.name}_description"
}
`, DataAlibabacloudstackVswitchZones)
}

func testAccDiskConfig_delete_auto_snapshot() string {
	return fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccDiskConfig"
}



resource "alibabacloudstack_disk" "default" {
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  	size = "70"
	disk_name = "${var.name}"
	description = "${var.name}_description"
	encrypted = "false"
	delete_auto_snapshot = "true"
}
`, DataAlibabacloudstackVswitchZones)
}

func testAccDiskConfig_delete_with_instance() string {
	return fmt.Sprintf(`
%s


variable "name" {
	default = "tf-testAccDiskConfig"
}

resource "alibabacloudstack_disk" "default" {
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  	size = "70"
	disk_name = "${var.name}"
	description = "${var.name}_description"
	encrypted = "false"
	delete_auto_snapshot = "true"
	delete_with_instance = "true"
}
`, DataAlibabacloudstackVswitchZones)
}

func testAccDiskConfig_enable_auto_snapshot() string {
	return fmt.Sprintf(`
%s


variable "name" {
	default = "tf-testAccDiskConfig"
}

resource "alibabacloudstack_disk" "default" {
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  	size = "70"
	disk_name = "${var.name}"
	description = "${var.name}_description"
	encrypted = "false"
	delete_auto_snapshot = "true"
	delete_with_instance = "true"
	enable_auto_snapshot = "true"
}
`, DataAlibabacloudstackVswitchZones)
}

func testAccDiskConfig_all() string {
	return fmt.Sprintf(`
%s


variable "name" {
	default = "tf-testAccDiskConfig_all"
}

resource "alibabacloudstack_disk" "default" {
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  	size = "70"
	disk_name = "${var.name}"
	description = "nothing"
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
`, DataAlibabacloudstackVswitchZones)
}

func testAccDiskConfig_multi() string {
	return fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccDiskConfig"
}

resource "alibabacloudstack_disk" "default" {
	count = "3"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  	size = "50"
	disk_name = "${var.name}_multi"
	description = "nothing"
	category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
	encrypted = "false"
}
`, DataAlibabacloudstackVswitchZones)
}

var testAccCheckResourceDiskBasicMap = map[string]string{
	"zone_id":              CHECKSET,
	"size":                 "50",
	"disk_name":            "",
	"description":          "",
	"category":             CHECKSET,
	"snapshot_id":          "",
	"encrypted":            "false",
	"tags.%":               "0",
	"status":               string(Available),
	"delete_auto_snapshot": "false",
	"delete_with_instance": "false",
	"enable_auto_snapshot": "false",
}
