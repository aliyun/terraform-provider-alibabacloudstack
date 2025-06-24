package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSlbAccessLog_basic0(t *testing.T) {
	var v *Slblogsdownloadattribute
	resourceId := "alibabacloudstack_slb_access_log.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackSlbAccesslogCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeAccessLogsDownloadAttribute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-test-log%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackSlbAccesslogDependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"log_project":      "${alibabacloudstack_log_project.default.name}",
					"log_store":        "${alibabacloudstack_log_store.default.name}",
					"load_balancer_id": "${alibabacloudstack_slb.default.id}",
					"depends_on":       []string{"alibabacloudstack_slb_listener.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_project": name,
						"log_store":   name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// 该参数无回读信息
				ImportStateVerifyIgnore: []string{"role_name"},
			},
		},
	})
}

var AlibabacloudStackSlbAccesslogCheckMap = map[string]string{}

func AlibabacloudStackSlbAccesslogDependence0(name string) string {
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

`, name, ECSInstanceCommonTestCase, ServerCertificateTestCase(), RsaPrivateKeyTestCase())
}
