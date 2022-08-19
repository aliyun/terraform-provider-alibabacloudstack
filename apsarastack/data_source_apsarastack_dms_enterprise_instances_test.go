package apsarastack

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccApsaraStackDmsEnterprisesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.apsarastack_dms_enterprise_instances.default"
	name := fmt.Sprintf("tf_testAccDmsEnterpriseInstancesDataSource_%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		name, dataSourceDmsEnterpriseInstancesConfigDependence)

	searchkeyConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"search_key": "${apsarastack_dms_enterprise_instance.default.host}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"search_key": "${apsarastack_dms_enterprise_instance.default.host}-fake",
		}),
	}
	instancealiasRegexConfConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"net_type":             "CLASSIC",
			"instance_type":        "${apsarastack_dms_enterprise_instance.default.instance_type}",
			"env_type":             "test",
			"instance_alias_regex": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"net_type":             "CLASSIC",
			"instance_type":        "${apsarastack_dms_enterprise_instance.default.instance_type}",
			"env_type":             "test",
			"instance_alias_regex": name + "fake",
		}),
	}
	nameRegexConfConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"net_type":      "CLASSIC",
			"instance_type": "${apsarastack_dms_enterprise_instance.default.instance_type}",
			"env_type":      "test",
			"name_regex":    name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"net_type":      "CLASSIC",
			"instance_type": "${apsarastack_dms_enterprise_instance.default.instance_type}",
			"env_type":      "test",
			"name_regex":    name + "fake",
		}),
	}
	var existDmsEnterpriseInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#":                   "1",
			"instances.0.data_link_name":    "",
			"instances.0.database_password": CHECKSET,
			"instances.0.database_user":     "admin123",
			"instances.0.dba_id":            CHECKSET,
			"instances.0.dba_nick_name":     CHECKSET,
			"instances.0.ddl_online":        "0",
			"instances.0.ecs_instance_id":   "",
			"instances.0.ecs_region":        os.Getenv("APSARASTACK_REGION"),
			"instances.0.env_type":          "test",
			"instances.0.export_timeout":    CHECKSET,
			"instances.0.host":              CHECKSET,
			"instances.0.instance_alias":    CHECKSET,
			"instances.0.instance_id":       CHECKSET,
			"instances.0.instance_source":   "RDS",
			"instances.0.instance_type":     "mysql",
			"instances.0.port":              "3306",
			"instances.0.query_timeout":     CHECKSET,
			"instances.0.safe_rule_id":      CHECKSET,
			"instances.0.sid":               "",
			"instances.0.status":            CHECKSET,
			"instances.0.use_dsql":          "0",
			"instances.0.vpc_id":            "",
		}
	}

	var fakeDmsEnterpriseInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#": "0",
		}
	}

	var DmsEnterpriseInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDmsEnterpriseInstancesMapFunc,
		fakeMapFunc:  fakeDmsEnterpriseInstancesMapFunc,
	}

	DmsEnterpriseInstancesCheckInfo.dataSourceTestCheck(t, rand, searchkeyConf, instancealiasRegexConfConf, nameRegexConfConf)
}

func dataSourceDmsEnterpriseInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
	data "apsarastack_account" "current" {
	}
	
	resource "apsarastack_db_instance" "instance" {
	engine           = "MySQL"
	engine_version   = "5.6"
	instance_type    = "rds.mysql.t1.small"
	instance_storage = "10"
	instance_name    = "%[1]s"
	security_ips     = ["0.0.0.0/0"]
	storage_type         = "local_ssd"
	}
	
	resource "apsarastack_db_account" "account" {
	instance_id = "${apsarastack_db_instance.instance.id}"
	name        = "admin123"
	password    = "inputYourCodeHere"
	type        = "Normal"
	}

	resource "apsarastack_ascm_user" "user" {
	 cellphone_number = "13900000000"
	 email = "test@gmail.com"
	 display_name = "C2C-DELTA"
	 organization_id = 7
	 mobile_nation_code = "91"
	 login_name = "User_Dms_%[1]s"
	 login_policy_id = 1
	}

	resource "apsarastack_dms_enterprise_user" "default" {
		  uid = apsarastack_ascm_user.user.user_id
		  user_name = apsarastack_ascm_user.user.login_name
		  mobile = "15910799999"
		  role_names = ["ADMIN"]
	}


	resource "apsarastack_dms_enterprise_instance" "default" {
	  dba_uid           =  apsarastack_dms_enterprise_user.default.uid
	  host              =  "${apsarastack_db_instance.instance.connection_string}"
	  port              =  "3306"
	  network_type      =	"CLASSIC"
	  safe_rule         =	"自由操作"
	  tid               =  "1"
	  instance_type     =	 "mysql"
	  instance_source   =	 "RDS"
	  env_type          =	 "test"
	  database_user     =	 apsarastack_db_account.account.name
	  database_password =	 apsarastack_db_account.account.password
	  instance_alias    =	 "%[1]s"
	  query_timeout     =	 "70"
	  export_timeout    =	 "2000"
	  ecs_region        =	 "%[2]s"
	  ddl_online        =	 "0"
	  use_dsql          =	 "0"
	  data_link_name    =	 ""
	}
`, name, os.Getenv("APSARASTACK_REGION"))
}
