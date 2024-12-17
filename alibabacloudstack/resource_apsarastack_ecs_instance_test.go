package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"log"
	"testing"

	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_instance", &resource.Sweeper{
		Name: "alibabacloudstack_instance",
		F:    testSweepInstances,
	})
}

func testSweepInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alibabacloudstack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var insts []ecs.Instance
	req := ecs.CreateDescribeInstancesRequest()
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
			return ecsClient.DescribeInstances(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Instances: %s", err)
		}
		resp, _ := raw.(*ecs.DescribeInstancesResponse)
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
			if need, err := vpcService.needSweepVpc(v.VpcAttributes.VpcId, v.VpcAttributes.VSwitchId); err == nil {
				skip = !need
			}

		}
		if skip {
			log.Printf("[INFO] Skipping Instance: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Instance: %s (%s)", name, id)
		if v.DeletionProtection {
			request := ecs.CreateModifyInstanceAttributeRequest()
			if strings.ToLower(client.Config.Protocol) == "https" {
				request.Scheme = "https"
			} else {
				request.Scheme = "http"
			}
			request.InstanceId = id
			request.DeletionProtection = requests.NewBoolean(false)
			_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceAttribute(request)
			})
			if err != nil {
				log.Printf("[ERROR] %#v", errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR))
				continue
			}
		}
		if v.InstanceChargeType == string(PrePaid) {
			request := ecs.CreateModifyInstanceChargeTypeRequest()
			if strings.ToLower(client.Config.Protocol) == "https" {
				request.Scheme = "https"
			} else {
				request.Scheme = "http"
			}
			request.InstanceIds = convertListToJsonString(append(make([]interface{}, 0, 1), id))
			request.InstanceChargeType = string(PostPaid)
			request.IncludeDataDisks = requests.NewBoolean(true)
			_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceChargeType(request)
			})
			if err != nil {
				log.Printf("[ERROR] %#v", errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR))
				continue
			}
			time.Sleep(3 * time.Second)
		}

		req := ecs.CreateDeleteInstanceRequest()
		req.InstanceId = id
		if strings.ToLower(client.Config.Protocol) == "https" {
			req.Scheme = "https"
		} else {
			req.Scheme = "http"
		}
		req.Force = requests.NewBoolean(true)
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Instance (%s (%s)): %s", name, id, err)
		} else {
			sweeped = true
		}
	}
	if sweeped {
		// Waiting 20 seconds to eusure these instances have been deleted.
		time.Sleep(20 * time.Second)
	}
	return nil
}

func TestAccAlibabacloudStackInstanceBasic(t *testing.T) {
	var v ecs.Instance

	resourceId := "alibabacloudstack_instance.default"
	ra := resourceAttrInit(resourceId, testAccInstanceCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsInstanceConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceInstanceBasicConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id":        "${data.alibabacloudstack_images.default.images.0.id}",
					"security_groups": []string{"${alibabacloudstack_ecs_securitygroup.default.id}"},
					"instance_type":   "${local.default_instance_type_id}",
					// "storage_set_id":       "${alibabacloudstack_ecs_ebs_storage_set.default.id}",
					// "storage_set_partition_number": "1",
					"availability_zone":    "${data.alibabacloudstack_zones.default.zones[0].id}",
					"system_disk_category": "cloud_ssd",
					"instance_name":        "${var.name}",
					//"key_name":                      "${alibabacloudstack_key_pair.default.key_name}",
					"user_data":                     "I_am_user_data",
					"security_enhancement_strategy": "Active",
					"vswitch_id":                    "${alibabacloudstack_vpc_vswitch.default.id}",
					"tags": map[string]string{
						"foo": "foo",
						"Bar": "Bar",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						//"key_name":      name,
						//"vswitch_id":    REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_enhancement_strategy", "dry_run", "user_data", "enable_ipv6", "ipv6_address_count"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]interface{}{
						"Created": "Terraform",
						"For":     "Instance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "Terraform",
						"tags.For":     "Instance test",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackInstanceVpc(t *testing.T) {
	var v ecs.Instance

	resourceId := "alibabacloudstack_instance.default"
	ra := resourceAttrInit(resourceId, testAccInstanceCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAcc%sEcsInstanceConfigVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceInstanceVpcConfigDependence)

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
					"image_id":        "${data.alibabacloudstack_images.default.images.0.id}",
					"security_groups": []string{"${alibabacloudstack_ecs_securitygroup.default.id}"},
					"instance_type":   "${local.default_instance_type_id}",

					"availability_zone":             "${data.alibabacloudstack_zones.default.zones.0.id}",
					"system_disk_category":          "cloud_ssd",
					"instance_name":                 "${var.name}",
					"security_enhancement_strategy": "Active",
					"user_data":                     "I_am_user_data",

					"vswitch_id": "${alibabacloudstack_vpc_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						//"key_name":      name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_enhancement_strategy", "user_data", "enable_ipv6", "ipv6_address_count"},
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"user_data": "I_am_user_data_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_data": "I_am_user_data_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_groups": []string{"${alibabacloudstack_ecs_securitygroup.default.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_groups.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_change",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_description",
					}),
				),
			},
			//			{
			//				Config: testAccConfig(map[string]interface{}{
			//					"internet_max_bandwidth_out": "50",
			//				}),
			//				Check: resource.ComposeTestCheckFunc(
			//					testAccCheck(map[string]string{
			//						"internet_max_bandwidth_out": "50",
			//						"private_ip":                 CHECKSET,
			//					}),
			//				),
			//			},
			//			{
			//				Config: testAccConfig(map[string]interface{}{
			//					"internet_max_bandwidth_in": "50",
			//				}),
			//				Check: resource.ComposeTestCheckFunc(
			//					testAccCheck(map[string]string{
			//						"internet_max_bandwidth_in": "50",
			//					}),
			//				),
			//			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_name": "hostNameExample",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_name": "hostNameExample",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "inputYourCod@",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "inputYourCod@",
					}),
				),
			},
			/*{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_ip": "172.16.0.10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_ip": "172.16.0.10",
					}),
				),
			},*/

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"foo": "foo",
						"Bar": "Bar",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":   "2",
						"tags.foo": "foo",
						"tags.Bar": "Bar",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackInstanceDataDisks(t *testing.T) {
	var v ecs.Instance

	resourceId := "alibabacloudstack_instance.default"
	ra := resourceAttrInit(resourceId, testAccInstanceCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAcc%sEcsInstanceDataDisks%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceInstancePrePaidConfigDependence)

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
					"image_id":        "${data.alibabacloudstack_images.default.images.0.id}",
					"security_groups": []string{"${alibabacloudstack_ecs_securitygroup.default.id}"},
					"instance_type":   "${local.default_instance_type_id}",

					"availability_zone":    "${data.alibabacloudstack_zones.default.zones.0.id}",
					"system_disk_category": "cloud_ssd",
					"instance_name":        "${var.name}",
					//"key_name":                      "${alibabacloudstack_key_pair.default.key_name}",
					"security_enhancement_strategy": "Active",
					"user_data":                     "I_am_user_data",

					"vswitch_id": "${alibabacloudstack_vpc_vswitch.default.id}",
					"data_disks": []map[string]string{
						{
							"name":        "disk1",
							"size":        "40",
							"category":    "cloud_ssd",
							"description": "disk1",
						},
						{
							"name":        "disk2",
							"size":        "40",
							"category":    "cloud_ssd",
							"description": "disk2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						//"key_name":      name,
						"user_data": "I_am_user_data",

						"data_disks.#":             "2",
						"data_disks.0.name":        "disk1",
						"data_disks.0.size":        "40",
						"data_disks.0.category":    "cloud_ssd",
						"data_disks.0.description": "disk1",
						"data_disks.1.name":        "disk2",
						"data_disks.1.size":        "40",
						"data_disks.1.category":    "cloud_ssd",
						"data_disks.1.description": "disk2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_enhancement_strategy", "data_disks", "user_data"},
			},
		},
	})
}

func TestAccAlibabacloudStackInstanceTypeUpdate(t *testing.T) {
	var v ecs.Instance

	resourceId := "alibabacloudstack_instance.default"
	ra := resourceAttrInit(resourceId, testAccInstanceCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsInstanceConfigInstanceType%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceInstanceTypeConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id":                      "${data.alibabacloudstack_images.default.images.0.id}",
					"system_disk_category":          "cloud_ssd",
					"system_disk_size":              "40",
					"instance_type":                 "${data.alibabacloudstack_instance_types.new1.instance_types.0.id}",
					"instance_name":                 "${var.name}",
					"security_groups":               []string{"${alibabacloudstack_ecs_securitygroup.default.id}"},
					"security_enhancement_strategy": "Active",
					"user_data":                     "I_am_user_data",
					"vswitch_id":                    "${alibabacloudstack_vpc_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "${data.alibabacloudstack_instance_types.new2.instance_types.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackInstanceMulti(t *testing.T) {
	var v ecs.Instance

	resourceId := "alibabacloudstack_instance.default.2"
	ra := resourceAttrInit(resourceId, testAccInstanceCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAcc%sEcsInstanceConfigMulti%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceInstanceVpcConfigDependence)

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
					"count":                "3",
					"image_id":             "${data.alibabacloudstack_images.default.images.0.id}",
					"security_groups":      []string{"${alibabacloudstack_ecs_securitygroup.default.id}"},
					"instance_type":        "${local.default_instance_type_id}",
					"availability_zone":    "${data.alibabacloudstack_zones.default.zones.0.id}",
					"system_disk_category": "cloud_ssd",
					"instance_name":        "${var.name}",
					//"key_name":                      "${alibabacloudstack_key_pair.default.key_name}",
					"security_enhancement_strategy": "Active",
					"user_data":                     "I_am_user_data",

					"vswitch_id": "${alibabacloudstack_vpc_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						//"key_name":      name,
						//"role_name":     name,
					}),
				),
			},
			/*{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},*/
		},
	})
}

func resourceInstanceVpcConfigDependence(name string) string {
	return fmt.Sprintf(`
%s

%s

%s

variable "name" {
	default = "%s"
}

//resource "alibabacloudstack_key_pair" "default" {
//	key_name = "${var.name}"
//}

`, SecurityGroupCommonTestCase , DataAlibabacloudstackImages , DataAlibabacloudstackInstanceTypes, name)
}

func resourceInstancePrePaidConfigDependence(name string) string {
	return fmt.Sprintf(`
%s

%s

%s

variable "name" {
	default = "%s"
}
//resource "alibabacloudstack_key_pair" "default" {
//	key_name = "${var.name}"
//}

`, SecurityGroupCommonTestCase , DataAlibabacloudstackImages , DataAlibabacloudstackInstanceTypes, name)
}

func resourceInstanceBasicConfigDependence(name string) string {
	return fmt.Sprintf(`
%s

%s

%s

variable "name" {
	default = "%s"
}

//resource "alibabacloudstack_key_pair" "default" {
//	key_name = "${var.name}"
//}

`, SecurityGroupCommonTestCase, DataAlibabacloudstackImages, DataAlibabacloudstackInstanceTypes, name)
}

func resourceInstanceTypeConfigDependence(name string) string {
	return fmt.Sprintf(`
%s

%s

%s

variable "name" {
	default = "%s"
}

data "alibabacloudstack_instance_types" "new1" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count    = 2
  sorted_by         = "Memory"
}
data "alibabacloudstack_instance_types" "new2" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count    = 4
  sorted_by         = "Memory"
}



`, SecurityGroupCommonTestCase , DataAlibabacloudstackImages , DataAlibabacloudstackInstanceTypes, name)
}

var testAccInstanceCheckMap = map[string]string{
	"image_id":                      CHECKSET,
	"instance_type":                 CHECKSET,
	"security_groups.#":             "1",
	"enable_ipv6":                   "false",
	"ipv6_address_count":            "0",
	"availability_zone":             CHECKSET,
	"system_disk_category":          "cloud_ssd",
	"security_enhancement_strategy": "Active",
	"vswitch_id":                    CHECKSET,
	"user_data":                     "I_am_user_data",

	"description":      "",
	"host_name":        CHECKSET,
	"password":         "",
	"system_disk_size": "40",

	"data_disks.#": NOSET,
	//"tags.%":       NOSET,
	//"tags.%": CHECKSET,

	"private_ip":                 CHECKSET,
	"status":                     "Running",
	"internet_max_bandwidth_out": "0",
	"force_delete":               NOSET,
}
