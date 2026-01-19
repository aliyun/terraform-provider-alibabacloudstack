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
					"health_check":                 "off",
					"protocol":                     "http",
					"master_slave_server_group_id": "${alibabacloudstack_slb_master_slave_server_group.default.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"load_balancer_id": CHECKSET,
						"bandwidth":        "10",
						"frontend_port":    "80",
						"backend_port":     "80",
						// "forward_port":     "8080",
						"sticky_session": "off",
						"health_check":   "off",
						"protocol":       "http",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"master_slave_server_group_id": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"master_slave_server_group_id": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"x_forwarded_for": []map[string]interface{}{{
						"retrive_slb_ip":    false,
						"retrive_slb_id":    false,
						"retrive_slb_proto": false,
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"x_forwarded_for.0.retrive_slb_ip":    "false",
						"x_forwarded_for.0.retrive_slb_id":    "false",
						"x_forwarded_for.0.retrive_slb_proto": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"x_forwarded_for": []map[string]interface{}{{
						"retrive_slb_ip":    true,
						"retrive_slb_id":    true,
						"retrive_slb_proto": true,
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"x_forwarded_for.0.retrive_slb_ip":    "true",
						"x_forwarded_for.0.retrive_slb_id":    "true",
						"x_forwarded_for.0.retrive_slb_proto": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"x_forwarded_for": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"x_forwarded_for.0.retrive_slb_ip":    REMOVEKEY,
						"x_forwarded_for.0.retrive_slb_id":    REMOVEKEY,
						"x_forwarded_for.0.retrive_slb_proto": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"health_check":                 "on",
					"health_check_domain":          `example.com`,
					"health_check_method":          "get",
					"master_slave_server_group_id": "${alibabacloudstack_slb_master_slave_server_group.default.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check":        "on",
						"health_check_domain": `example.com`,
						"health_check_method": "get",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// delete_protection_validation is a local attribute and cannot be loaded from remote
				ImportStateVerifyIgnore: []string{"delete_protection_validation", "master_slave_server_group_id"},
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

func AlibabacloudTestAccSlbListenerBasicdependence2(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alibabacloudstack_slb_server_certificate" "default" {
	name = "${var.name}"
	server_certificate = %s
	private_key = %s
  }

resource "alibabacloudstack_slb_ca_certificate" "default" {
  name = "${var.name}"
  ca_certificate = %s
}

resource "alibabacloudstack_slb_acl" "default" {
	name = "${var.name}"
	ip_version = "ipv4"
  }

resource "alibabacloudstack_slb_acl" "new" {
	name = "${var.name}-new"
	ip_version = "ipv4"
  }

resource "alibabacloudstack_slb_loadbalancer" "default" {
	name = "${var.name}"
	// vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	//address_type       = "internet"
  	specification        = "slb.s2.small"
  }

`, name, ServerCertificateTestCase(), RsaPrivateKeyTestCase(), ServerCertificateTestCase())
}

func AlibabacloudTestAccSlbListenerBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s

resource "alibabacloudstack_slb_server_certificate" "default" {
	name = "${var.name}"
	server_certificate = %s
	private_key = %s
  }

resource "alibabacloudstack_slb_ca_certificate" "default" {
  name = "${var.name}"
  ca_certificate = %s
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

resource "alibabacloudstack_ecs_instance" "new" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 40
  system_disk_name     = "test_sys_diskv2"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id    = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false
  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

resource "alibabacloudstack_slb_master_slave_server_group" "default" {
  count = 2
  load_balancer_id = "${alibabacloudstack_slb_loadbalancer.default.id}"
  name = "${var.name}"
  servers {
      server_id = "${alibabacloudstack_ecs_instance.default.id}"
      port = 80
      weight = 100
      server_type = "Master"
  }
  servers {
      server_id = "${alibabacloudstack_ecs_instance.new.id}"
      port = 80
      weight = 100
      server_type = "Slave"
  }
}

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

`, name, ECSInstanceCommonTestCase, ServerCertificateTestCase(), RsaPrivateKeyTestCase(), ServerCertificateTestCase())
}

func AlibabacloudTestAccSlbListenerBasicdependence3(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s

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

`, name, ECSInstanceCommonTestCase, ServerCertificateTestCase(), RsaPrivateKeyTestCase())
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
					"ca_certificate_id":     "${alibabacloudstack_slb_ca_certificate.default.id}",
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

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbListenerBasicdependence2)
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

					"load_balancer_id":      "${alibabacloudstack_slb_loadbalancer.default.id}",
					"protocol":              "https",
					"server_certificate_id": "${alibabacloudstack_slb_server_certificate.default.id}",
					"bandwidth":             "10",
					"frontend_port":         "80",
					"backend_port":          "80",
					"acl_status":            "on",
					"sticky_session":        "off",
					"health_check":          "off",
					"acl_type":              "white",
					"description":           "testcreate",
					"scheduler":             "wrr",
					"acl_id":                "${alibabacloudstack_slb_acl.default.id}",
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
						"protocol":                             "https",
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
			{
				Config: testAccConfig(map[string]interface{}{
					"logs_download_attributes": []map[string]string{
						{
							"log_store":   "${alibabacloudstack_log_store.new.name}",
							"log_project": "${alibabacloudstack_log_project.new.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logs_download_attributes.#":           "1",
						"logs_download_attributes.0.log_store": name + "-new",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"logs_download_attributes": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"logs_download_attributes.#":           REMOVEKEY,
						"logs_download_attributes.0.log_store": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackSlbListener4(t *testing.T) {

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

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbListenerBasicdependence2)
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
			{
				Config: testAccConfig(map[string]interface{}{
					"persistence_timeout": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"persistence_timeout": "10",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// delete_protection_validation is a local attribute and cannot be loaded from remote
				ImportStateVerifyIgnore: []string{"server_group_id", "delete_protection_validation", "sticky_session"},
			},
		},
	})
}

func TestAccAlibabacloudStackSlbListener5(t *testing.T) {

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

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbListenerBasicdependence2)
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

					"load_balancer_id":             "${alibabacloudstack_slb_loadbalancer.default.id}",
					"protocol":                     "http",
					"delete_protection_validation": "true",
					"bandwidth":                    "10",
					"frontend_port":                "80",
					"backend_port":                 "80",
					// "listener_forward":             "on",
					// "forward_port":                 "443",
					"acl_status":     "on",
					"sticky_session": "off",
					"health_check":   "off",
					"acl_type":       "white",
					"description":    "testcreate",
					"scheduler":      "wrr",
					"acl_id":         "${alibabacloudstack_slb_acl.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"load_balancer_id": CHECKSET,
						"protocol":         "http",
						"bandwidth":        "10",
						"frontend_port":    "80",
						"backend_port":     "80",
						// "listener_forward": "on",
						// "forward_port":     "443",
						"acl_status":  "on",
						"acl_type":    "white",
						"description": "testcreate",
						"scheduler":   "wrr",
						"acl_id":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"acl_id":   "${alibabacloudstack_slb_acl.new.id}",
					"acl_type": "black",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"acl_id":   CHECKSET,
						"acl_type": "black",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"acl_status": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"acl_status": "off",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackSlbListener6(t *testing.T) {

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

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbListenerBasicdependence3)
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

					"load_balancer_id":             "${alibabacloudstack_slb_listener.new.load_balancer_id}",
					"bandwidth":                    "10",
					"frontend_port":                "80",
					"backend_port":                 "80",
					"delete_protection_validation": "true",
					"protocol":                     "http",
					"listener_forward":             "on",
					"forward_port":                 "443",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"load_balancer_id": CHECKSET,
						"bandwidth":        "10",
						"frontend_port":    "80",
						"backend_port":     "80",
						"listener_forward": "on",
						"forward_port":     "443",
						"protocol":         "http",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackSlbListener7(t *testing.T) {

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

					"load_balancer_id":             "${alibabacloudstack_slb_loadbalancer.default.id}",
					"bandwidth":                    "10",
					"frontend_port":                "80",
					"backend_port":                 "80",
					"sticky_session":               "off",
					"delete_protection_validation": "true",
					"cookie_timeout":               "86400",
					"cookie":                       "ALB",
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
						// "forward_port":     "8080",
						"sticky_session": "off",
						"health_check":   "off",
						"protocol":       "http",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_id": "${alibabacloudstack_slb_server_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_id": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_id": REMOVEKEY,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"sticky_session":      "on",
					"sticky_session_type": "insert",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sticky_session":      "on",
						"sticky_session_type": "insert",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"sticky_session_type": "server",
					"cookie":              "cookie-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sticky_session_type": "server",
						"cookie":              "cookie-test",
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

func AlibabacloudTestAccSlbListenerLogStoredependence(name string) string {
	return AlibabacloudTestAccSlbListenerBasicdependence2(name) + `
	resource "alibabacloudstack_log_project" "default" {
		name = "${var.name}"
		description = "test"
	}
	resource "alibabacloudstack_log_project" "new" {
		name = "${var.name}-new"
		description = "test-new"
	}
	resource "alibabacloudstack_log_store" "new" {
		name = "${var.name}-new"
		project = "${alibabacloudstack_log_project.new.name}"
		retention_period      = "30"
		shard_count           = "2"
		enable_web_tracking   = false
		auto_split            = true
		max_split_shard_count = "64"
		append_meta           = true
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
`
}
