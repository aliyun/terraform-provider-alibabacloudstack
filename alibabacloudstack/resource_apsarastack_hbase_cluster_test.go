package alibabacloudstack

import (
	"testing"

	

	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("aalibabacloudstack_hbase_instance", &resource.Sweeper{
		Name: "alibabacloudstack_hbase_instance",
		F:    testSweepHBaseInstances,
	})
}

func testSweepHBaseInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var insts []hbase.Instance
	req := hbase.CreateDescribeInstancesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
			return hbaseClient.DescribeInstances(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving HBase Instances: %s", err)
		}
		resp, _ := raw.(*hbase.DescribeInstancesResponse)
		if resp == nil || len(resp.Instances.Instance) < 1 {
			break
		}
		insts = append(insts, resp.Instances.Instance...)

		if len(resp.Instances.Instance) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	sweeped := false
	vpcService := VpcService{client}
	for _, v := range insts {
		name := v.InstanceName
		id := v.InstanceId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a slb name is set by other service, it should be fetched by vswitch name and deleted.
		if skip {
			if need, err := vpcService.needSweepVpc(v.VpcId, ""); err == nil {
				skip = !need
			}

		}

		if skip {
			log.Printf("[INFO] Skipping Hbase Instance: %s (%s)", name, id)
			continue
		}

		log.Printf("[INFO] Deleting HBase Instance: %s (%s)", name, id)
		req := hbase.CreateDeleteInstanceRequest()
		req.ClusterId = id
		_, err := client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
			return hbaseClient.DeleteInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Hbase Instance (%s (%s)): %s", name, id, err)
		} else {
			sweeped = true
		}
	}
	if sweeped {
		// Waiting 30 seconds to eusure these DB instances have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func AlibabacloudStackHbaseBasicDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
	default = "%s"
}
data "alibabacloudstack_zones" "default" {}

data "alibabacloudstack_vpcs" "default" {
	name_regex = "default-NODELETING"
}
resource "alibabacloudstack_vpc" "default" {
name       = var.name
cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = data.alibabacloudstack_zones.default.ids.0
  name              = "${var.name}"
}

resource "alibabacloudstack_security_group" "default" {
	count = 2
	vpc_id = alibabacloudstack_vpc.default.id
	name = var.name
}
`, name)
}

func TestAccAlibabacloudStackHBaseInstanceVpc(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_hbase_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HBaseService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeHBaseInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackHbaseBasicDependence)
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
					"name":                  "${var.name}",
					"engine_version":        "2.0",
					"zone_id":               "${data.alibabacloudstack_zones.default.zones.0.id}",
					"master_instance_type":  "hbase.sn1.large",
					"core_instance_type":    "hbase.sn1.large",
					"core_disk_type":        "cloud_efficiency",
					"vswitch_id":            "${alibabacloudstack_vswitch.default.id}",
					"immediate_delete_flag": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                     name,
						"engine":                   "hbase",
						"engine_version":           "2.0",
						"core_instance_type":       "hbase.sn1.large",
						"core_disk_type":           "cloud_efficiency",
						"vswitch_id":               CHECKSET,
						"immediate_delete_flag":    "true",
						"core_instance_quantity":   "2",
						"cold_storage_size":        "0",
						"deletion_protection":      "true",
						"zone_id":                  CHECKSET,
						"master_instance_quantity": CHECKSET,
						"maintain_start_time":      CHECKSET,
						"maintain_end_time":        CHECKSET,
						"pay_type":                 "PostPaid",
						"ip_white":                 CHECKSET,
						"security_groups.#":        "0",
						"ui_proxy_conn_addrs.#":    "1",
						"zk_conn_addrs.#":          "3",
						"slb_conn_addrs.#":         "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"immediate_delete_flag"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_start_time": "04:00Z",
					"maintain_end_time":   "06:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_start_time": "04:00Z",
						"maintain_end_time":   "06:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"core_disk_size": "440",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"core_disk_size": "440",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_white": "192.168.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_white": "192.168.1.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account":  "admin",
					"password": "inputYourCodeHere",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account":  "admin",
						"password": "inputYourCodeHere",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_storage_size": "800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_storage_size": "800",
					}),
				),
			},
			//{
			//	Config: resourceHBaseConfigPrePaid,
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"name":                "tf_testAccHBase_vpc_change_name",
			//			"maintain_start_time": "04:00Z",
			//			"maintain_end_time":   "06:00Z",
			//			"tags.%":              "2",
			//			"tags.Created":        "TF",
			//			"tags.For":            "acceptance test",
			//			"ip_white":            "192.168.1.1",
			//			"cold_storage_size":   "800",
			//			"pay_type":            "PrePaid",
			//		}),
			//	),
			//},

			{
				Config: testAccConfig(map[string]interface{}{
					"security_groups": []string{"${alibabacloudstack_security_group.default.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_groups.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                name,
					"maintain_start_time": "14:00Z",
					"maintain_end_time":   "16:00Z",
					"security_groups":     []string{"${alibabacloudstack_security_group.default.0.id}", "${alibabacloudstack_security_group.default.1.id}"},
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "acceptance test 123",
					},
					"core_disk_size":    "480",
					"ip_white":          "192.168.1.2",
					"account":           "adminu",
					"password":          "inputYourCodeHere",
					"cold_storage_size": "900",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                name,
						"maintain_start_time": "14:00Z",
						"maintain_end_time":   "16:00Z",
						"security_groups.#":   "2",
						"tags.%":              "2",
						"tags.Created":        "TF-update",
						"tags.For":            "acceptance test 123",
						"core_disk_size":      "480",
						"ip_white":            "192.168.1.2",
						"account":             "adminu",
						"password":            "inputYourCodeHere",
						"cold_storage_size":   "900",
					}),
				),
			},
		},
	})
}
