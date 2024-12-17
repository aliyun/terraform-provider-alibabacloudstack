package alibabacloudstack

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_nas_file_system",
		&resource.Sweeper{
			Name: "alibabacloudstack_nas_file_system",
			F:    testSweepNasFileSystem,
		})
}

func testSweepNasFileSystem(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeFileSystems"
	request := make(map[string]interface{})
	request["RegionId"] = client.Region
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var response map[string]interface{}
	conn, err := client.NewNasClient()
	if err != nil {
		return errmsgs.WrapError(err)
	}
	ids := make([]string, 0)
	for {
		runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			log.Printf("[ERROR] Error retrieving filesystem: %s", err)
		}
		resp, err := jsonpath.Get("$.FileSystems.FileSystem", response)
		if err != nil {
			log.Println("Get $.FileSystems.FileSystem failed. err:", err)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			description, _ := item["Description"].(string)
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(description), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping FileSystem: %s (%s)", description, item["FileSystemId"])
				continue
			}
			// 删除 fileSystem 时需要先删除其挂载关系
			if v, ok := item["MountTargets"].(map[string]interface{})["MountTarget"].([]interface{}); ok && len(v) > 0 {
				log.Printf("[INFO] Delete mount targets with filesystem: %v", item["FileSystemId"])
				for _, domain := range v {
					domainInfo := domain.(map[string]interface{})
					request := map[string]interface{}{
						"FileSystemId":      item["FileSystemId"],
						"MountTargetDomain": domainInfo["MountTargetDomain"],
					}
					action := "DeleteMountTarget"
					runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
					runtime.SetAutoretry(true)
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &runtime)
					if err != nil {
						log.Printf("[ERROR] Error delete mount target: %v with filesystem: %v err: %v", domainInfo["MountTargetDomain"], item["FileSystemId"], err)
					}
				}
			}
			ids = append(ids, item["FileSystemId"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	for _, filesystemId := range ids {
		request := map[string]interface{}{
			"FileSystemId": filesystemId,
		}
		action := "DeleteFileSystem"
		runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			log.Printf("[ERROR] Error delete filesystem: %s err: %v", filesystemId, err)
		}
	}
	return nil
}

func TestAccAlibabacloudStackNasFileSystem_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeNasFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlibabacloudStackNasFileSystem%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackNasFileSystemBasicDependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type": "${data.alibabacloudstack_nas_protocols.example.protocols.0}",
					"storage_type":  "Capacity",
					"zone_id":       "${data.alibabacloudstack_nas_zones.default.zones.0.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type": CHECKSET,
						"storage_type":  "Capacity",
						"zone_id":       CHECKSET,
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
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "Update",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackNasFileSystemEncrypt(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_nas_file_system.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeNasFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlibabacloudStackNasFileSystem%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackNasFileSystemBasicDependence1)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type": "NFS",
					"storage_type":  "Capacity",
					"encrypt_type":  "0",
					"zone_id":       "${data.alibabacloudstack_nas_zones.default.zones.0.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type": CHECKSET,
						"storage_type":  "Capacity",
						"encrypt_type":  "0",
						"zone_id":       CHECKSET,
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
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "Update",
					}),
				),
			},
		},
	})
}

var AlibabacloudStackNasFileSystem0 = map[string]string{}

func AlibabacloudStackNasFileSystemBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alibabacloudstack_nas_protocols" "example" {
        type = "Capacity"
}
data "alibabacloudstack_nas_zones" "default" {
}
`, name)
}

func AlibabacloudStackNasFileSystemBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_nas_zones" "default" {
}
`, name)
}
