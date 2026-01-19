package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackSlbAccessLogsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_slb_access_logs.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-slblog-data%d", rand),
		dataSourceSlbAccessLogsDependence)

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"log_store_regex": "${alibabacloudstack_slb_access_log.default.log_store}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"log_store_regex": "${alibabacloudstack_slb_access_log.default.log_store}-fakeTestAcccc",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_slb_access_log.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_slb_access_log.default.id}-fakeTestAcccc"},
		}),
	}

	loadBalancerIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_access_log.default.load_balancer_id}",
			"ids":              []string{"${alibabacloudstack_slb_access_log.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_access_log.default.load_balancer_id}-fakeTestAcccc",
			"ids":              []string{"${alibabacloudstack_slb_access_log.default.id}-fakeTestAcccc"},
		}),
	}

	var existSlbAccessLogsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"ids.0":                          CHECKSET,
			"access_logs.#":                  "1",
			"access_logs.0.log_store":        CHECKSET,
			"access_logs.0.log_project":      CHECKSET,
			"access_logs.0.load_balancer_id": CHECKSET,
			"access_logs.0.log_type":         CHECKSET,
		}
	}

	var fakeSlbAccessLogsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "0",
			"access_logs.#": "0",
		}
	}

	var SlbAccessLogsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existSlbAccessLogsMapFunc,
		fakeMapFunc:  fakeSlbAccessLogsMapFunc,
	}

	SlbAccessLogsCheckInfo.dataSourceTestCheck(t, rand, descriptionRegexConf, idsConf, loadBalancerIdConf)
}

func dataSourceSlbAccessLogsDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

resource "alibabacloudstack_slb" "default" {
  name          = "${var.name}_slb"
  vswitch_id    = "${alibabacloudstack_vpc_vswitch.default.id}"
}
resource "alibabacloudstack_slb_server_certificate" "servercertificate" {
  name               = "slbservercertificate"
  server_certificate = %s
  private_key        = %s
}

resource "alibabacloudstack_slb_listener" "default" {
    load_balancer_id            = alibabacloudstack_slb.default.id
    server_certificate_id       = alibabacloudstack_slb_server_certificate.servercertificate.id
    sticky_session              = "off"
    sticky_session_type         = "insert"
    cookie_timeout              = 1000
    frontend_port               = 120
    backend_port                = 120
    enable_http2                = "on"
    acl_status                  = "off"
    acl_type                    = "white"
    protocol                    = "https"
    bandwidth                   = -1
    gzip                        = true
    x_forwarded_for {
        retrive_slb_ip = false
        retrive_slb_id = false
        retrive_slb_proto = false
    }
    tls_cipher_policy           = "tls_cipher_policy_1_2"
    health_check                = "on"
    health_check_type           = "http"
    health_check_uri            = "/"
    health_check_connect_port   = 20
    health_check_method         = "head"
    healthy_threshold           = "3"
    unhealthy_threshold         = "3"
    health_check_timeout        = "5"
    health_check_interval       = "2"
    health_check_http_code      = "http_2xx,http_3xx"
    description                 = "testslblistener"
	lifecycle {
	    ignore_changes = [
		logs_download_attributes
	    ]
	}
}

resource "alibabacloudstack_log_project" "default" {
	name = "${var.name}"
	description = "test"
}

resource "alibabacloudstack_log_store" "default" {
	name = "${var.name}"
	project = "${alibabacloudstack_log_project.default.name}"	
	retention_period      = "30"
	shard_count           = "2"
	enable_web_tracking   = false
	auto_split            = true
	max_split_shard_count = "64"
	append_meta           = true
}

resource "alibabacloudstack_slb_access_log" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  log_project = "${var.name}"
  log_store = "${var.name}"
}
`, name, ECSInstanceCommonTestCase, ServerCertificateTestCase(), RsaPrivateKeyTestCase())
}
