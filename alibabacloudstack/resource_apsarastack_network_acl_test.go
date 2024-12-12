package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_network_acl", &resource.Sweeper{
		Name: "alibabacloudstack_network_acl",
		F:    testSweepNetworkAcl,
	})
}

func testSweepNetworkAcl(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeNetworkAcls"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return errmsgs.WrapError(err)
	}
	networkAclIds := make([]string, 0)
	for {
		runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			log.Printf("Error retrieving network acl: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.NetworkAcls.NetworkAcl", response)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, action, "$.NetworkAcls.NetworkAcl", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := fmt.Sprint(item["NetworkAclName"])
			id := fmt.Sprint(item["NetworkAclId"])
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Network Acl: %s (%s)", name, id)
				continue
			}
			networkAclIds = append(networkAclIds, id)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	vpcService := VpcService{client}
	for _, id := range networkAclIds {
		//	Delete attach resources
		object, err := vpcService.DescribeNetworkAcl(id)
		if err != nil {
			log.Println("DescribeNetworkAcl failed", err)
		}
		deleteResources, _ := object["Resources"].(map[string]interface{})["Resource"].([]interface{})
		if len(deleteResources) > 0 {
			request := map[string]interface{}{
				"NetworkAclId": id,
			}
			resourcesMaps := make([]map[string]interface{}, 0)
			for _, resources := range deleteResources {
				resourcesArg := resources.(map[string]interface{})
				resourcesMap := map[string]interface{}{
					"ResourceId":   resourcesArg["ResourceId"],
					"ResourceType": resourcesArg["ResourceType"],
				}
				resourcesMaps = append(resourcesMaps, resourcesMap)
			}
			request["Resource"] = resourcesMaps
			request["RegionId"] = client.RegionId
			action := "UnassociateNetworkAcl"
			runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
				if err != nil {
					if errmsgs.NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				log.Println("UnassociateNetworkAcl failed", err)
			}
			stateConf := BuildStateConf([]string{}, []string{"Available"}, 5*time.Minute, 5*time.Second, vpcService.NetworkAclStateRefreshFunc(id, []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				log.Println("UnassociateNetworkAcl failed", err)
			}
		}

		log.Printf("[INFO] Deleting Network Acl: (%s)", id)
		request := map[string]interface{}{
			"NetworkAclId": id,
		}
		action := "DeleteNetworkAcl"
		request["RegionId"] = client.RegionId
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
			if err != nil {
				if errmsgs.NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Network Acl (%s): %s", id, err)
		}
	}
	return nil
}

func TestAccAlibabacloudStackVpcNetworkAcl_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_network_acl.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackNetworkAclMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeNetworkAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snetworkacl%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackNetworkAclBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":           "${alibabacloudstack_vpc.default.id}",
					"network_acl_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":           CHECKSET,
						"network_acl_name": name,
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
					"description": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"egress_acl_entries": []map[string]interface{}{
						{
							"description":            "engress test",
							"destination_cidr_ip":    "10.0.0.0/24",
							"network_acl_entry_name": "tf-testacc78924",
							"policy":                 "accept",
							"port":                   "20/80",
							"protocol":               "tcp",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"egress_acl_entries.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ingress_acl_entries": []map[string]interface{}{
						{
							"description":            "ingress test",
							"network_acl_entry_name": "tf-testacc78999",
							"policy":                 "accept",
							"port":                   "20/80",
							"protocol":               "tcp",
							"source_cidr_ip":         "10.0.0.0/24",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ingress_acl_entries.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_acl_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_acl_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resources": []map[string]interface{}{
						{
							"resource_id":   "${alibabacloudstack_vswitch.default0.id}",
							"resource_type": "VSwitch",
						},
						{
							"resource_id":   "${alibabacloudstack_vswitch.default1.id}",
							"resource_type": "VSwitch",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resources.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resources": []map[string]interface{}{
						{
							"resource_id":   "${alibabacloudstack_vswitch.default0.id}",
							"resource_type": "VSwitch",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resources.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      name,
					"network_acl_name": name,
					"ingress_acl_entries": []map[string]interface{}{
						{
							"description":            "ingress test change",
							"network_acl_entry_name": "tf-testacc78999",
							"policy":                 "accept",
							"port":                   "20/80",
							"protocol":               "tcp",
							"source_cidr_ip":         "10.0.0.0/24",
						},
					},
					"egress_acl_entries": []map[string]interface{}{
						{
							"description":            "engress test change",
							"destination_cidr_ip":    "10.0.0.0/24",
							"network_acl_entry_name": "tf-testacc78924",
							"policy":                 "accept",
							"port":                   "20/80",
							"protocol":               "tcp",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":           name,
						"network_acl_name":      name,
						"ingress_acl_entries.#": "1",
						"egress_acl_entries.#":  "1",
					}),
				),
			},
		},
	})
}

var AlibabacloudStackNetworkAclMap0 = map[string]string{}

func AlibabacloudStackNetworkAclBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%[1]s"
		}
variable "name_change" {
			default = "%[1]s_change"
		}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}
resource "alibabacloudstack_vpc" "default" {
  cidr_block = "192.168.0.0/16"
  name = var.name
}
resource "alibabacloudstack_vswitch" "default0" {
  vpc_id            = alibabacloudstack_vpc.default.id
  name      = var.name
  cidr_block        = cidrsubnets(alibabacloudstack_vpc.default.cidr_block, 4, 4)[0]
  availability_zone           = data.alibabacloudstack_zones.default.ids.0
}
resource "alibabacloudstack_vswitch" "default1" {
  vpc_id            = alibabacloudstack_vpc.default.id
  name      = var.name_change
  cidr_block        = cidrsubnets(alibabacloudstack_vpc.default.cidr_block, 4, 4)[1]
  availability_zone           = data.alibabacloudstack_zones.default.ids.0
}

`, name)
}
