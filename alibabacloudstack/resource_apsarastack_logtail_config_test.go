package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_logtail_config", &resource.Sweeper{
		Name: "alibabacloudstack_logtail_config",
		F:    testSweepLogConfigs,
	})
}

func testSweepLogConfigs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	logService := LogService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	raw, err := logService.GetSlsDataClient("")
	if err != nil {
		log.Printf("[ERROR] Error getting SLS client: %s", errmsgs.WrapError(err))
		return err
	}

	// Get all projects
	projects, err := raw.ListProject()
	if err != nil {
		log.Printf("[ERROR] Error retrieving Log Projects: %s", errmsgs.WrapError(err))
		return err
	}

	for _, v := range projects {
		name := v
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				cfNameList, _, cfErr := raw.ListConfig(name, 0, 100)
				if cfErr != nil {
					log.Printf("[ERROR] Error retrieving Log config: %s", errmsgs.WrapError(cfErr))
					continue
				}
				for _, cfName := range cfNameList {
					log.Printf("[INFO] Deleting Log config: %s", cfName)
					if delErr := raw.DeleteConfig(name, cfName); delErr != nil {
						log.Printf("[ERROR] Failed to delete Log Config (%s): %s", cfName, delErr)
					}
				}
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Log Project: %s", name)
			continue
		}
		log.Printf("[INFO] Deleting Log Project: %s", name)
		if delErr := raw.DeleteProject(name); delErr != nil {
			log.Printf("[ERROR] Failed to delete Log Project (%s): %s", name, delErr)
		}
	}
	return nil
}

func TestAccAlibabacloudStackLogTail_basic(t *testing.T) {
	var v *sls.LogConfig
	resourceId := "alibabacloudstack_logtail_config.default"
	ra := resourceAttrInit(resourceId, logTailMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("testtf-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogTailDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project":      "${alibabacloudstack_log_project.default.name}",
					"logstore":     "${alibabacloudstack_log_store.default.name}",
					"input_type":   "file",
					"name":         name,
					"output_type":  "LogService",
					"input_detail": `{\"discardUnmatch\":false,\"enableRawLog\":true,\"fileEncoding\":\"gbk\",\"filePattern\":\"access.log\",\"logPath\":\"/logPath\",\"logType\":\"json_log\",\"maxDepth\":10,\"topicFormat\":\"default\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"input_detail": "{\"discardUnmatch\":false,\"enableRawLog\":true,\"fileEncoding\":\"gbk\",\"filePattern\":\"access.log\",\"logPath\":\"/logPath\",\"logType\":\"json_log\",\"maxDepth\":10,\"topicFormat\":\"default\"}",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"input_detail"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"input_type":   "file",
					"input_detail": `{\"autoExtend\":true,\"discardUnmatch\":true,\"enableRawLog\":true,\"fileEncoding\":\"utf8\",\"filePattern\":\"*\",\"key\":[\"test\",\"test2\"],\"logPath\":\"/logPath\",\"logType\":\"delimiter_log\",\"maxDepth\":999,\"quote\":\"\\\"\",\"separator\":\",\",\"timeKey\":\"test\",\"topicFormat\":\"default\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"input_type":   "file",
						"input_detail": "{\"autoExtend\":true,\"discardUnmatch\":true,\"enableRawLog\":true,\"fileEncoding\":\"utf8\",\"filePattern\":\"*\",\"key\":[\"test\",\"test2\"],\"logPath\":\"/logPath\",\"logType\":\"delimiter_log\",\"maxDepth\":999,\"quote\":\"\\\"\",\"separator\":\",\",\"timeKey\":\"test\",\"topicFormat\":\"default\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"input_type":   "file",
					"input_detail": `{\"discardUnmatch\":false,\"enableRawLog\":true,\"fileEncoding\":\"gbk\",\"filePattern\":\"access.log\",\"logPath\":\"/logPath\",\"logType\":\"json_log\",\"maxDepth\":10,\"topicFormat\":\"default\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"input_detail": "{\"discardUnmatch\":false,\"enableRawLog\":true,\"fileEncoding\":\"gbk\",\"filePattern\":\"access.log\",\"logPath\":\"/logPath\",\"logType\":\"json_log\",\"maxDepth\":10,\"topicFormat\":\"default\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"input_type":   "file",
					"input_detail": `{\"adjustTimezone\":false,\"advanced\":{\"blacklist\":{},\"k8s\":{\"ExternalEnvTag\":{}},\"tail_size_kb\":1024},\"discardUnmatch\":false,\"dockerExcludeEnv\":{},\"dockerExcludeLabel\":{},\"dockerFile\":false,\"dockerIncludeEnv\":{},\"dockerIncludeLabel\":{},\"enableRawLog\":false,\"fileEncoding\":\"utf8\",\"filePattern\":\"apsarach.log\",\"formatNickname\":\"customized\",\"key\":[\"remote_addr\",\"remote_ident\",\"remote_user\",\"time_local\",\"request_method\",\"request_uri\",\"request_protocol\",\"status\",\"response_size_bytes\"],\"logPath\":\"/var/log\",\"logTimezone\":\"\",\"logType\":\"common_reg_log\",\"maxDepth\":10,\"preserve\":true,\"preserveDepth\":1,\"regex\":\"([0-9.-]+)\\\\s([\\\\w.-]+)\\\\s([\\\\w.-]+)\\\\s(\\\\[[^\\\\[\\\\]]+\\\\]|-)\\\\s\\\"((?:[^\\\"]|\\\\\\\")+)\\\\s((?:[^\\\"]|\\\\\\\")+)\\\\s((?:[^\\\"]|\\\\\\\")+)\\\"\\\\s(\\\\d{3}|-)\\\\s(\\\\d+|-).*\",\"topicFormat\":\"none\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"input_detail": "{\"adjustTimezone\":false,\"advanced\":{\"blacklist\":{},\"k8s\":{\"ExternalEnvTag\":{}},\"tail_size_kb\":1024},\"discardUnmatch\":false,\"dockerExcludeEnv\":{},\"dockerExcludeLabel\":{},\"dockerFile\":false,\"dockerIncludeEnv\":{},\"dockerIncludeLabel\":{},\"enableRawLog\":false,\"fileEncoding\":\"utf8\",\"filePattern\":\"apsarach.log\",\"formatNickname\":\"customized\",\"key\":[\"remote_addr\",\"remote_ident\",\"remote_user\",\"time_local\",\"request_method\",\"request_uri\",\"request_protocol\",\"status\",\"response_size_bytes\"],\"logPath\":\"/var/log\",\"logTimezone\":\"\",\"logType\":\"common_reg_log\",\"maxDepth\":10,\"preserve\":true,\"preserveDepth\":1,\"regex\":\"([0-9.-]+)\\\\s([\\\\w.-]+)\\\\s([\\\\w.-]+)\\\\s(\\\\[[^\\\\[\\\\]]+\\\\]|-)\\\\s\\\"((?:[^\\\"]|\\\\\\\")+)\\\\s((?:[^\\\"]|\\\\\\\")+)\\\\s((?:[^\\\"]|\\\\\\\")+)\\\"\\\\s(\\\\d{3}|-)\\\\s(\\\\d+|-).*\",\"topicFormat\":\"none\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"input_type":   "file",
					"input_detail": `{\"adjustTimezone\":false,\"advanced\":{\"blacklist\":{},\"k8s\":{\"ExternalEnvTag\":{}},\"tail_size_kb\":1024},\"discardUnmatch\":false,\"dockerExcludeEnv\":{},\"dockerExcludeLabel\":{},\"dockerFile\":false,\"dockerIncludeEnv\":{},\"dockerIncludeLabel\":{},\"enableRawLog\":false,\"fileEncoding\":\"utf8\",\"filePattern\":\"*.log\",\"key\":[\"content\"],\"logBeginRegex\":\".*\",\"logPath\":\"/var/log\",\"logTimezone\":\"\",\"logType\":\"common_reg_log\",\"maxDepth\":10,\"preserve\":true,\"preserveDepth\":1,\"regex\":\"(.*)\",\"timeFormat\":\"\",\"topicFormat\":\"none\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"input_detail": "{\"adjustTimezone\":false,\"advanced\":{\"blacklist\":{},\"k8s\":{\"ExternalEnvTag\":{}},\"tail_size_kb\":1024},\"discardUnmatch\":false,\"dockerExcludeEnv\":{},\"dockerExcludeLabel\":{},\"dockerFile\":false,\"dockerIncludeEnv\":{},\"dockerIncludeLabel\":{},\"enableRawLog\":false,\"fileEncoding\":\"utf8\",\"filePattern\":\"*.log\",\"key\":[\"content\"],\"logBeginRegex\":\".*\",\"logPath\":\"/var/log\",\"logTimezone\":\"\",\"logType\":\"common_reg_log\",\"maxDepth\":10,\"preserve\":true,\"preserveDepth\":1,\"regex\":\"(.*)\",\"timeFormat\":\"\",\"topicFormat\":\"none\"}",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackLogTail_plugin(t *testing.T) {
	var v *sls.LogConfig
	resourceId := "alibabacloudstack_logtail_config.default"
	ra := resourceAttrInit(resourceId, logTailMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("testtf-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogTailDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project":      "${alibabacloudstack_log_project.default.name}",
					"logstore":     "${alibabacloudstack_log_store.default.name}",
					"input_type":   "plugin",
					"name":         name,
					"output_type":  "LogService",
					"input_detail": `{\"LocalStorage\":true,\"filterRegex\":[],\"filterKey\":[],\"topicFormat\":\"none\",\"discardUnmatch\":false,\"plugin\":{\"inputs\":[{\"type\":\"service_mysql\",\"detail\":{\"Address\":\"************.mysql.rds.aliyuncs.com\",\"User\":\"****\",\"Password\":\"*******\",\"DataBase\":\"****\",\"Limit\":true,\"PageSize\":100,\"StateMent\":\"select * from db.VersionOs where time > ?\",\"CheckPoint\":true,\"CheckPointColumn\":\"time\",\"CheckPointStart\":\"2018-01-01 00:00:00\",\"CheckPointSavePerPage\":true,\"CheckPointColumnType\":\"time\",\"IntervalMs\":60000}}]}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"input_type":   "plugin",
						"input_detail": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"input_detail": `{\"LocalStorage\":true,\"filterRegex\":[],\"filterKey\":[],\"topicFormat\":\"none\",\"discardUnmatch\":false,\"plugin\":{\"inputs\":[{\"type\":\"service_mysql\",\"detail\":{\"Address\":\"************.mysql.rds.aliyuncs.com\",\"User\":\"****\",\"Password\":\"*******\",\"DataBase\":\"****\",\"Limit\":true,\"PageSize\":10,\"StateMent\":\"select * from db.VersionOs where time > ?\",\"CheckPoint\":true,\"CheckPointColumn\":\"time\",\"CheckPointStart\":\"2018-01-01 00:00:00\",\"CheckPointSavePerPage\":true,\"CheckPointColumnType\":\"time\",\"IntervalMs\":60000}}]}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"input_detail": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"input_detail"},
			},
		},
	})
}

func resourceLogTailDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alibabacloudstack_log_project" "default"{
	name = "${var.name}"
	description = "create by terraform"
}
resource "alibabacloudstack_log_store" "default"{
  	project = "${alibabacloudstack_log_project.default.name}"
  	name = "${var.name}"
  	retention_period = 3650
  	shard_count = 3
  	auto_split = true
  	max_split_shard_count = 60
  	append_meta = true
}
`, name)
}

var logTailMap = map[string]string{
	"name":         CHECKSET,
	"project":      CHECKSET,
	"logstore":     CHECKSET,
	"input_type":   "file",
	"output_type":  "LogService",
	"input_detail": CHECKSET,
}
