package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/yundun_bastionhost"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_bastionhost_instance", &resource.Sweeper{
		Name: "alibabacloudstack_bastionhost_instance",
		F:    testSweepBastionhostInstances,
	})
}

func testSweepBastionhostInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	// prefixes := []string{
	// 	"tf-testAcc",
	// 	"tf_testAcc",
	// }
	request := yundun_bastionhost.CreateDescribeInstanceBastionhostRequest()
	request.PageSize = requests.NewInteger(PageSizeSmall)
	request.CurrentPage = requests.NewInteger(1)
	var instances []yundun_bastionhost.Instance

	for {
		raw, err := client.WithBastionhostClient(func(bastionhostClient *yundun_bastionhost.Client) (interface{}, error) {
			return bastionhostClient.DescribeInstanceBastionhost(request)
		})
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DataDefaultErrorMsg, "alibabacloudctack_yundun_bastionhost", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}

		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*yundun_bastionhost.DescribeInstanceBastionhostResponse)
		if len(response.Instances) < 1 {
			break
		}

		instances = append(instances, response.Instances...)

		if len(response.Instances) < PageSizeSmall {
			break
		}

		currentPageNo := request.CurrentPage
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DataDefaultErrorMsg, "alibabacloudctack_yundun_bastionhost", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}

		if page, err := getNextpageNumber(currentPageNo); err != nil {
			return errmsgs.WrapError(err)
		} else {
			request.CurrentPage = page
		}
	}

	// for _, v := range instances {
	// 	name := v.Description
	// 	skip := true
	// 	if !sweepAll() {
	// 		for _, prefix := range prefixes {
	// 			if name != "" && strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
	// 				skip = false
	// 				break
	// 			}
	// 		}
	// 		if skip {
	// 			log.Printf("[INFO] Skipping Bastionhost Instance: %s", name)
	// 			continue
	// 		}
	// 	}
	// 	log.Printf("[INFO] Deleting Bastionhost Instance %s .", v.InstanceId)

	// 	releaseReq := yundun_bastionhost.CreateRefundInstanceRequest()
	// 	releaseReq.InstanceId = v.InstanceId
	// 	_, err := client.WithBastionhostClient(func(bastionhostClient *yundun_bastionhost.Client) (interface{}, error) {
	// 		return bastionhostClient.RefundInstance(releaseReq)
	// 	})
	// 	if err != nil {
	// 		log.Printf("[ERROR] Deleting Instance %s got an error: %#v.", v.InstanceId, err)
	// 	}
	// }

	return nil
}

func TestAccAlibabacloudStackBastionhostInstance_basic(t *testing.T) {
	var v yundun_bastionhost.Instance
	resourceId := "alibabacloudstack_bastionhost_instance.default"
	ra := resourceAttrInit(resourceId, bastionhostInstanceBasicMap)
	serviceFunc := func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf_testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBastionhostInstanceDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreYunCheck(t)
		},
		IDRefreshName: resourceId,
		Providers: func() map[string]*schema.Provider {
			commonProvider := Provider()
			yundunProvider := Provider()
			yundunProvider.Schema["access_key"] = &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_YUNDUN_ACCESS_KEY", ""),
				Description: descriptions["access_key"],
			}
			yundunProvider.Schema["secret_key"] = &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_YUNDUN_SECRET_KEY", ""),
				Description: descriptions["secret_key"],
			}
			yundunProvider.Schema["role_arn"] = &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["assume_role_role_arn"],
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_YUNDUN_ASSUME_ROLE_ARN", ""),
			}
			return map[string]*schema.Provider{
				"alibabacloudstack":        yundunProvider,
				"alibabacloudstack-common": commonProvider,
			}
		}(),
		// resource "alibabacloudstack_bastionhost_instance" "default" {
		// 	vswitch_id = alibabacloudstack_vswitch.vsw.id
		// 	license_code = "bastionhostah_small_lic"
		// 	vpc_id = alibabacloudstack_vpc.vpc.id
		// 	asset = "50"
		// 	highavailability = "false"
		// 	disasterrecovery = "false"
		// 	provider = alibabacloudstack
		//   }
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id":       "${alibabacloudstack_vswitch.vsw.id}",
					"vpc_id":           "${alibabacloudstack_vpc.vpc.id}",
					"license_code":     "bastionhostah_small_lic",
					"highavailability": "false",
					"disasterrecovery": "false",
					"provider":         "alibabacloudstack",
					"asset":            "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":       CHECKSET,
						"asset":            "50",
						"highavailability": "false",
						"disasterrecovery": "false",
						"license_code":     "bastionhostah_small_lic",
						"vpc_id":           CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      "bastionhost_instance_test",
					"asset":            "60",
					"highavailability": "true",
					"disasterrecovery": "true",
					"license_code":     "bastionhostah_large_lic",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      "bastionhost_instance_test",
						"asset":            "60",
						"highavailability": "true",
						"disasterrecovery": "true",
						"license_code":     "bastionhostah_large_lic",
					}),
				),
			},
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"description": "${var.name}_update",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"description": name + "_update",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"description": "${var.name}",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"description": name,
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"license_code": "bhah_ent_100_asset",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"license_code": "bhah_ent_100_asset",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"security_group_ids": []string{"${alibabacloudctack_security_group.default.1.id}"},
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"security_group_ids.#": "1",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"tags": map[string]string{
			// 			"Created": "TF",
			// 			"For":     "acceptance-test",
			// 		},
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"tags.%":       "2",
			// 			"tags.Created": "TF",
			// 			"tags.For":     "acceptance-test",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"tags": map[string]string{
			// 			"Created": "TF",
			// 			"For":     "acceptance-test",
			// 			"Updated": "TF",
			// 		},
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"tags.%":       "3",
			// 			"tags.Created": "TF",
			// 			"tags.For":     "acceptance-test",
			// 			"tags.Updated": "TF",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"renew_period":        "2",
			// 		"renewal_period_unit": "M",
			// 		"renewal_status":      "AutoRenewal",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"renew_period":        "2",
			// 			"renewal_period_unit": "M",
			// 			"renewal_status":      "AutoRenewal",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"renewal_status": "NotRenewal",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"renewal_status": "NotRenewal",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"ad_auth_server": []map[string]interface{}{
			// 			{
			// 				"server":         "192.168.1.1",
			// 				"standby_server": "192.168.1.3",
			// 				"port":           "80",
			// 				"domain":         "domain",
			// 				"account":        "cn=Manager,dc=test,dc=com",
			// 				"password":       "YouPassword123",
			// 				"filter":         "objectClass=person",
			// 				"name_mapping":   "nameAttr",
			// 				"email_mapping":  "emailAttr",
			// 				"mobile_mapping": "mobileAttr",
			// 				"is_ssl":         "true",
			// 				"base_dn":        "dc=test,dc=com",
			// 			},
			// 		},
			// 		"ldap_auth_server": []map[string]interface{}{
			// 			{
			// 				"server":             "192.168.1.1",
			// 				"standby_server":     "192.168.1.3",
			// 				"port":               "80",
			// 				"login_name_mapping": "uid",
			// 				"account":            "cn=Manager,dc=test,dc=com",
			// 				"password":           "YouPassword123",
			// 				"filter":             "objectClass=person",
			// 				"name_mapping":       "nameAttr",
			// 				"email_mapping":      "emailAttr",
			// 				"mobile_mapping":     "mobileAttr",
			// 				"is_ssl":             "true",
			// 				"base_dn":            "dc=test,dc=com",
			// 			},
			// 		},
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"ad_auth_server.#":   "1",
			// 			"ldap_auth_server.#": "1",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"ad_auth_server": []map[string]interface{}{
			// 			{
			// 				"server":   "192.168.1.1",
			// 				"port":     "80",
			// 				"is_ssl":   "false",
			// 				"domain":   "domain",
			// 				"account":  "cn=Manager,dc=test,dc=com",
			// 				"password": "YouPassword123",
			// 				"base_dn":  "dc=test,dc=com",
			// 			},
			// 		},
			// 		"ldap_auth_server": []map[string]interface{}{
			// 			{
			// 				"server":   "192.168.1.1",
			// 				"port":     "80",
			// 				"password": "YouPassword123",
			// 				"account":  "cn=Manager,dc=test,dc=com",
			// 				"base_dn":  "dc=test,dc=com",
			// 			},
			// 		},
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"ad_auth_server.#":   "1",
			// 			"ldap_auth_server.#": "1",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"resource_group_id":  "${data.alibabacloudctack_resource_manager_resource_groups.default.ids.0}",
			// 		"description":        "${var.name}",
			// 		"license_code":       "bhah_ent_200_asset",
			// 		"security_group_ids": []string{"${alibabacloudctack_security_group.default.0.id}", "${alibabacloudctack_security_group.default.1.id}"},
			// 		"tags":               REMOVEKEY,
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"resource_group_id":    CHECKSET,
			// 			"description":          name,
			// 			"license_code":         "bhah_ent_200_asset",
			// 			"security_group_ids.#": "2",
			// 			"tags.%":               REMOVEKEY,
			// 			"tags.Created":         REMOVEKEY,
			// 			"tags.For":             REMOVEKEY,
			// 			"tags.Updated":         REMOVEKEY,
			// 		}),
			// 	),
			// },
			// {
			// 	ResourceName:      resourceId,
			// 	ImportState:       true,
			// 	ImportStateVerify: false,
			// },
		},
	})
}

func TestAccAlibabacloudStackBastionhostInstance_PublicAccess(t *testing.T) {
	var v yundun_bastionhost.Instance
	resourceId := "alibabacloudctack_bastionhost_instance.default"
	ra := resourceAttrInit(resourceId, bastionhostInstanceBasicMap)
	serviceFunc := func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf_testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBastionhostInstanceDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"license_code":         "bhah_ent_50_asset",
					"period":               "1",
					"plan_code":            "cloudbastion",
					"storage":              "5",
					"bandwidth":            "10",
					"description":          "${var.name}",
					"vswitch_id":           "${local.vswitch_id}",
					"security_group_ids":   []string{"${alibabacloudctack_security_group.default.0.id}"},
					"enable_public_access": "false",
					"public_white_list":    []string{"192.168.0.0/16"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":          name,
						"period":               "1",
						"plan_code":            "cloudbastion",
						"security_group_ids.#": "1",
						"enable_public_access": "false",
						"public_white_list.#":  "1",
						"public_white_list.0":  "192.168.0.0/16",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_public_access": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_public_access": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"public_white_list": []string{"192.168.0.0/18"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"public_white_list.#": "1",
						"public_white_list.0": "192.168.0.0/18",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"period", "storage", "bandwidth"},
			},
		},
	})
}

func resourceBastionhostInstanceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_zones" "default" {
	provider = alibabacloudstack-common
	available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "vpc" {	
	provider = alibabacloudstack-common
	vpc_name = var.name
	cidr_block = "192.168.0.0/16" #vpc口段
}
resource "alibabacloudstack_vswitch" "vsw" {
	provider = alibabacloudstack-common
	vpc_id = alibabacloudstack_vpc.vpc.id
	cidr_block = "192.168.0.0/16" #⽹段
	availability_zone = data.alibabacloudstack_zones.default.zones.0.id #可⽤区
}
`, name)
}

var bastionhostInstanceBasicMap = map[string]string{
	// "description":          CHECKSET,
	// "license_code":         "bhah_ent_50_asset",
	// "period":               "1",
	// "vswitch_id":           CHECKSET,
	// "security_group_ids.#": "1",
}
