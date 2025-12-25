package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSlbListener0(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_slb_listener.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbListenerCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeSlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc-slblistener%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbListenerBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"load_balancer_id": "${alibabacloudstack_slb_loadbalancer.default.id}",
					"bandwidth":        "10",
					"frontend_port":    "80",
					"backend_port":     "80",
					"sticky_session":   "off",
					// "sticky_session_type": "",
					"health_check": "off",
					"protocol":     "http",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"load_balancer_id": CHECKSET,
						"bandwidth":        "10",
						"frontend_port":    "80",
						"backend_port":     "80",
						"sticky_session":   "off",
						"health_check":     "off",
						"protocol":         "http",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// delete_protection_validation is a local attribute and cannot be loaded from remote
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
		},
	})
}

var AlibabacloudTestAccSlbListenerCheckmap = map[string]string{

	// "description": CHECKSET,

	// "scheduler": CHECKSET,

	// "acl_id": CHECKSET,

	// "server_group_id": CHECKSET,

	// "load_balancer_id": CHECKSET,

	// "acl_status": CHECKSET,

	// "bandwidth": CHECKSET,

	// "acl_type": CHECKSET,
}

func AlibabacloudTestAccSlbListenerBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alibabacloudstack_slb_server_certificate" "default" {
	name = "${var.name}"
	server_certificate = %s
	private_key = %s
  }

resource "alibabacloudstack_slb_acl" "default" {
	name = "${var.name}"
	ip_version = "ipv4"
  }

resource "alibabacloudstack_slb_loadbalancer" "default" {
	name = "${var.name}"
	// vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	//address_type       = "internet"
  	specification        = "slb.s2.small"
  }

`, name, ServerCertificateTestCase(), RsaPrivateKeyTestCase())
}
func TestAccAlibabacloudStackSlbListener1(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_slb_listener.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbListenerCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeSlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc-slblistener%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbListenerBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"load_balancer_id": "${alibabacloudstack_slb_loadbalancer.default.id}",
					"bandwidth":        "10",
					"frontend_port":    "80",
					"backend_port":     "80",
					"sticky_session":   "off",
					// "sticky_session_type": "",
					"health_check":          "off",
					"protocol":              "https",
					"server_certificate_id": "${alibabacloudstack_slb_server_certificate.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"load_balancer_id": CHECKSET,
						"bandwidth":        "10",
						"frontend_port":    "80",
						"backend_port":     "80",
						"sticky_session":   "off",
						// "sticky_session_type": "",
						"health_check": "off",
						"protocol":     "https",
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackSlbListener2(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_slb_listener.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbListenerCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeSlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc-slblistener%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbListenerBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alibabacloudstack_slb_loadbalancer.default.id}",
					"bandwidth":        "10",
					"frontend_port":    "80",
					"backend_port":     "80",
					"sticky_session":   "off",
					// "sticky_session_type": "",
					"health_check": "on",
					"protocol":     "tcp",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id": CHECKSET,
						"bandwidth":        "10",
						"frontend_port":    "80",
						"backend_port":     "80",
						// "sticky_session_type": "",
						"health_check": "on",
						"protocol":     "tcp",
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackSlbListener3(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_slb_listener.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbListenerCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeSlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc-slblistener%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbListenerLogStoredependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"load_balancer_id": "${alibabacloudstack_slb_loadbalancer.default.id}",
					"protocol":         "http",
					"bandwidth":        "10",
					"frontend_port":    "80",
					"backend_port":     "80",
					"acl_status":       "on",
					"sticky_session":   "off",
					"health_check":     "off",
					"acl_type":         "white",
					"description":      "testcreate",
					"scheduler":        "wrr",
					"acl_id":           "${alibabacloudstack_slb_acl.default.id}",
					"logs_download_attributes": []map[string]string{
						{
							"log_store":   "${alibabacloudstack_log_store.default.name}",
							"log_project": "${alibabacloudstack_log_project.default.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"load_balancer_id":                     CHECKSET,
						"protocol":                             "http",
						"bandwidth":                            "10",
						"frontend_port":                        "80",
						"backend_port":                         "80",
						"acl_status":                           "on",
						"acl_type":                             "white",
						"description":                          "testcreate",
						"scheduler":                            "wrr",
						"acl_id":                               CHECKSET,
						"logs_download_attributes.#":           "1",
						"logs_download_attributes.0.log_store": name,
					}),
				),
			},
		},
	})
}

func AlibabacloudTestAccSlbListenerLogStoredependence(name string) string {
	return AlibabacloudTestAccSlbListenerBasicdependence(name) + `
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
	}`
}
