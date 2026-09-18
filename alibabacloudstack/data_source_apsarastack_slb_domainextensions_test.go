package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackSlbDomainExtensionsDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_slb_domain_extensions.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccCheckAlibabacloudStackSlbDataSourceBasic-%d", rand),
		dataSourceSlbDomainExtensionsConfigDependence)

	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_domain_extension.default.load_balancer_id}",
			"frontend_port":    "${alibabacloudstack_slb_domain_extension.default.frontend_port}",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_domain_extension.default.load_balancer_id}",
			"frontend_port":    "${alibabacloudstack_slb_domain_extension.default.listener_port}",
			"ids":              []string{"${alibabacloudstack_slb_domain_extension.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_domain_extension.default.load_balancer_id}",
			"frontend_port":    "${alibabacloudstack_slb_domain_extension.default.listener_port}",
			"ids":              []string{"${alibabacloudstack_slb_domain_extension.default.id}_fake"},
		}),
	}

	var existDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"extensions.#":                       "1",
			"ids.#":                              "1",
			"extensions.0.id":                    CHECKSET,
			"extensions.0.domain":                "www.test.com",
			"extensions.0.server_certificate_id": CHECKSET,
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"slb_rules.#": "0",
			"ids.#":       "0",
			"names.#":     "0",
		}
	}

	var slbDomainExtensionsCheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_slb_domain_extensions.default",
		existMapFunc: existDnsRecordsMapFunc,
		fakeMapFunc:  fakeDnsRecordsMapFunc,
	}

	slbDomainExtensionsCheckInfo.dataSourceTestCheck(t, rand, basicConf, allConf)
}

func dataSourceSlbDomainExtensionsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alibabacloudstack_slb_server_certificate" "default" {
	name = "${var.name}"
	server_certificate = %s
	private_key = %s
}

resource "alibabacloudstack_slb_loadbalancer" "default" {
	name = "${var.name}"
	specification        = "slb.s2.small"
}

%s

resource "alibabacloudstack_slb_server_group" "default" {
  vserver_group_name = "${var.name}"
  load_balancer_id = "${alibabacloudstack_slb_loadbalancer.default.id}"
  servers {
      server_ids = ["${alibabacloudstack_ecs_instance.default.id}"]
      port = 80
      weight = 100
    }
}

resource "alibabacloudstack_slb_listener" "new" {
  load_balancer_id        = "${alibabacloudstack_slb_loadbalancer.default.id}"
  bandwidth               = "10"
  frontend_port           = "443"
  backend_port            = "80"
  sticky_session          = "off"
  health_check            = "off"
  protocol                = "https"
  server_certificate_id   = "${alibabacloudstack_slb_server_certificate.default.id}"
}

resource "alibabacloudstack_slb_domain_extension" "default" {
  		load_balancer_id      = "${alibabacloudstack_slb_listener.new.load_balancer_id}"
  		frontend_port         = "${alibabacloudstack_slb_listener.new.frontend_port}"
  		domain                = "www.test.com"
  		server_certificate_id = "${alibabacloudstack_slb_listener.new.server_certificate_id}"
}
`, name, ServerCertificateTestCase(), RsaPrivateKeyTestCase(), ECSInstanceCommonTestCase)
}
