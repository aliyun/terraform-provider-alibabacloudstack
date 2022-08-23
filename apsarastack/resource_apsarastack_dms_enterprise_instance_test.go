package apsarastack

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("apsarastack_dms_enterprise_instance", &resource.Sweeper{
		Name: "apsarastack_dms_enterprise_instance",
		F:    testSweepDMSEnterpriseInstances,
	})
}

func testSweepDMSEnterpriseInstances(region string) error {

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting ApsaraStack client.")
	}
	client := rawClient.(*connectivity.ApsaraStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"testacc",
	}
	request := map[string]interface{}{
		"InstanceState": "NORMAL",
		"PageSize":      PageSizeXLarge,
		"PageNumber":    1,
	}
	var response map[string]interface{}
	action := "ListInstances"
	conn, err := client.NewDmsenterpriseClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_dms_enterprise_instances", action, ApsaraStackSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.InstanceList.Instance", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.InstanceList.Instance", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			id := item["Host"].(string) + ":" + item["Port"].(json.Number).String()

			skip := true

			for _, prefix := range prefixes {
				if item["InstanceAlias"] != nil {
					if strings.HasPrefix(strings.ToLower(item["InstanceAlias"].(string)), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
			}
			if skip || item["InstanceAlias"] == nil {
				log.Printf("[INFO] Skipping DMS Enterprise Instances: %s", id)
				continue
			}
			action := "DeleteInstance"
			request := map[string]interface{}{
				"Host": item["Host"],
				"Port": item["Port"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete DMS Enterprise Instance (%s (%s)): %s", item["InstanceAlias"].(string), id, err)
				continue
			}
			log.Printf("[INFO] Delete DMS Enterprise Instance Success: %s ", item["InstanceAlias"].(string))
		}
		if len(result) < PageSizeXLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccApsaraStackDmsEnterprise(t *testing.T) {
	resourceId := "apsarastack_dms_enterprise_instance.default"
	var v map[string]interface{}
	ra := resourceAttrInit(resourceId, testAccCheckKeyValueInMapsForDMS)

	serviceFunc := func() interface{} {
		return &Dms_enterpriseService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccDmsEnterpriseInstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDmsConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dba_uid":           "${apsarastack_dms_enterprise_user.default.uid}",
					"host":              "${apsarastack_db_instance.instance.connection_string}",
					"port":              "3306",
					"network_type":      "CLASSIC",
					"safe_rule":         "自由操作",
					"tid":               "1",
					"instance_type":     "mysql",
					"instance_source":   "RDS",
					"env_type":          "test",
					"database_user":     "${apsarastack_db_account.account.name}",
					"database_password": "${apsarastack_db_account.account.password}",
					"instance_alias":    name,
					"query_timeout":     "70",
					"export_timeout":    "2000",
					"ecs_region":        os.Getenv("APSARASTACK_REGION"),
					"ddl_online":        "0",
					"use_dsql":          "0",
					"data_link_name":    "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dba_uid":         CHECKSET,
						"host":            CHECKSET,
						"port":            "3306",
						"network_type":    "CLASSIC",
						"safe_rule":       "自由操作",
						"tid":             "1",
						"instance_type":   "mysql",
						"instance_source": "RDS",
						"env_type":        "test",
						"database_user":   CHECKSET,
						"instance_alias":  name,
						"query_timeout":   "70",
						"export_timeout":  "2000",
						"ecs_region":      os.Getenv("APSARASTACK_REGION"),
						"ddl_online":      "0",
						"use_dsql":        "0",
						"data_link_name":  "",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"database_password", "dba_uid", "network_type", "port", "safe_rule", "tid"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"env_type": "dev",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"env_type": "dev",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_alias": "other_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_alias": "other_name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_timeout": "77",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_timeout": "77",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"dba_uid":           "${apsarastack_dms_enterprise_user.default.uid}",
					"host":              "${apsarastack_db_instance.instance.connection_string}",
					"port":              "3306",
					"network_type":      "CLASSIC",
					"safe_rule":         "自由操作",
					"tid":               "1",
					"instance_type":     "mysql",
					"instance_source":   "RDS",
					"env_type":          "test",
					"database_user":     "${apsarastack_db_account.account.name}",
					"database_password": "${apsarastack_db_account.account.password}",
					"instance_alias":    name,
					"query_timeout":     "70",
					"export_timeout":    "2000",
					"ecs_region":        os.Getenv("APSARASTACK_REGION"),
					"ddl_online":        "0",
					"use_dsql":          "0",
					"data_link_name":    "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dba_uid":         CHECKSET,
						"host":            CHECKSET,
						"port":            "3306",
						"network_type":    "CLASSIC",
						"safe_rule":       "自由操作",
						"tid":             "1",
						"instance_type":   "mysql",
						"instance_source": "RDS",
						"env_type":        "test",
						"database_user":   CHECKSET,
						"instance_alias":  name,
						"query_timeout":   "70",
						"export_timeout":  "2000",
						"ecs_region":      os.Getenv("APSARASTACK_REGION"),
						"ddl_online":      "0",
						"use_dsql":        "0",
						"data_link_name":  "",
					}),
				),
			},
		},
	})
}

var testAccCheckKeyValueInMapsForDMS = map[string]string{}

func resourceDmsConfigDependence(name string) string {
	return fmt.Sprintf(`

	variable "name" {
		  default = "%v"
	}
	data "apsarastack_account" "current" {
	}

	resource "apsarastack_db_instance" "instance" {
	engine           = "MySQL"
	engine_version   = "5.6"
	instance_type    = "rds.mysql.t1.small"
	instance_storage = "10"
	instance_name    = "${var.name}"
	security_ips     = ["0.0.0.0/0"]
	storage_type         = "local_ssd"
	}
	
	resource "apsarastack_db_account" "account" {
	instance_id = "${apsarastack_db_instance.instance.id}"
	name        = "tftest123"
	password    = "inputYourCodeHere"
	type        = "Normal"
	}

	resource "apsarastack_ascm_user" "user" {
	 cellphone_number = "13900000000"
	 email = "test@gmail.com"
	 display_name = "C2C-DELTA"
	 organization_id = 33
	 mobile_nation_code = "91"
	 login_name = "User_Dms_${var.name}"
	 login_policy_id = 1
	}

	resource "apsarastack_dms_enterprise_user" "default" {
		  uid = apsarastack_ascm_user.user.user_id
		  user_name = apsarastack_ascm_user.user.login_name
		  mobile = "15910799999"
		  role_names = ["ADMIN"]
	}
	

	`, name)
}
