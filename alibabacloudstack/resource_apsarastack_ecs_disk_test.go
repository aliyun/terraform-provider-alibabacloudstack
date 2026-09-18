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

func testAccCheckDiskDestroyOrg(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_disk" {
			continue
		}
		client := testAccOrgProviders["alibabacloudstack"].Meta().(*connectivity.AlibabacloudStackClient)
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
	name := "tf-testAccDiskConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEcsDiskConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id": "${data.alibabacloudstack_zones.default.zones.0.id}",
					"size":    "50",
				}),
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
				Config: testAccConfig(map[string]interface{}{
					"size":      "70",
					"disk_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size":      "70",
						"disk_name": "tf-testAccDiskConfig",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_name":   "tf-testAccDiskConfig",
						"description": "tf-testAccDiskConfig_description",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypted":            "false",
					"delete_auto_snapshot": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_auto_snapshot": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_with_instance": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_with_instance": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_auto_snapshot": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_auto_snapshot": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_name":            "tf-testAccDiskConfig_all",
					"description":          "nothing",
					"delete_auto_snapshot": "false",
					"delete_with_instance": "false",
					"enable_auto_snapshot": "false",
					"tags": map[string]string{
						"name1": "name1",
						"name2": "name2",
						"name3": "name3",
					},
				}),
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
	name := "tf-testAccDiskConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEcsDiskMultiConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":       "3",
					"zone_id":     "${data.alibabacloudstack_zones.default.zones.0.id}",
					"size":        "50",
					"disk_name":   "${var.name}_multi",
					"description": "nothing",
					"category":    "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
					"encrypted":   "false",
				}),
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
		return &EcsService{testAccOrgProviders["alibabacloudstack"].Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serverFunc)
	ra := resourceAttrInit(resourceId, map[string]string{})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "testAccDiskConfig_encrypted"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEcsDiskEncryptedConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckOrg(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccOrgProviders,
		CheckDestroy:  testAccCheckDiskDestroyOrg,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              "testAccDiskConfig_encrypted",
					"availability_zone": "${data.alibabacloudstack_zones.default.zones.0.id}",
					"size":              "50",
					"category":          "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
					"encrypted":         true,
					"kms_key_id":        "${alibabacloudstack_kms_key.key.id}",
				}),
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

func resourceEcsDiskConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s
`, name, DataAlibabacloudstackVswitchZones)
}

func resourceEcsDiskMultiConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s
`, name, DataAlibabacloudstackVswitchZones)
}

func resourceEcsDiskEncryptedConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

%s
`, name, DataZoneCommonTestCase, KeyCommonTestCase)
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
