package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackSlbListenersDataSource_http(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_slb_listeners.default"
	name := fmt.Sprintf("tf-testacc-slblisteners%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSlbListenersHttpConfigDependence)

	// Test with load balancer ID filter
	lbIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "fake-lb-id",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id":  "${alibabacloudstack_slb_listener.default.load_balancer_id}",
			"description_regex": "^${alibabacloudstack_slb_listener.default.description}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id":  "${alibabacloudstack_slb_listener.default.load_balancer_id}",
			"description_regex": "fake_description",
		}),
	}

	// Test with protocol filter
	protocolConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
			"protocol":         "http",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
			"protocol":         "fake_protocol",
		}),
	}

	// Test with protocol filter
	portConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
			"frontend_port":    80,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
			"frontend_port":    81,
		}),
	}

	// Test with both filters
	bothFiltersConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
			"protocol":         "http",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
			"protocol":         "fake_protocol",
		}),
	}

	var existSlbListenersHttpMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"load_balancer_id": CHECKSET,
			"slb_listeners.#":  "1",
		}
	}

	var fakeSlbListenersHttpMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"load_balancer_id": CHECKSET,
			"slb_listeners.#":  "0",
		}
	}

	var slbListenersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existSlbListenersHttpMapFunc,
		fakeMapFunc:  fakeSlbListenersHttpMapFunc,
	}
	slbListenersCheckInfo.dataSourceTestCheck(t, rand, lbIdConf, protocolConf, nameRegexConf, portConf, bothFiltersConf)
}

func TestAccAlibabacloudStackSlbListenersDataSource_https(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_slb_listeners.default"
	name := fmt.Sprintf("tf-testacc-slblisteners%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSlbListenersHttpsConfigDependence)

	// Test with load balancer ID filter
	lbIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}_fake",
		}),
	}

	// Test with protocol filter
	protocolConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
			"protocol":         "https",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
			"protocol":         "fake_protocol",
		}),
	}

	// Test with both filters
	bothFiltersConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
			"protocol":         "https",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
			"protocol":         "fake_protocol",
		}),
	}

	var existSlbListenersHttpsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"load_balancer_id": CHECKSET,
			"slb_listeners.#":  "1",
		}
	}

	var fakeSlbListenersHttpsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"load_balancer_id": CHECKSET,
			"slb_listeners.#":  "0",
		}
	}

	var slbListenersHttpsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existSlbListenersHttpsMapFunc,
		fakeMapFunc:  fakeSlbListenersHttpsMapFunc,
	}
	slbListenersHttpsCheckInfo.dataSourceTestCheck(t, rand, lbIdConf, protocolConf, bothFiltersConf)
}

func TestAccAlibabacloudStackSlbListenersDataSource_tcp(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_slb_listeners.default"
	name := fmt.Sprintf("tf-testacc-slblisteners%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSlbListenersTcpConfigDependence)

	// Test with load balancer ID filter
	lbIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}_fake",
		}),
	}

	// Test with protocol filter
	protocolConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
			"protocol":         "tcp",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
			"protocol":         "fake_protocol",
		}),
	}

	// Test with both filters
	bothFiltersConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
			"protocol":         "tcp",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
			"protocol":         "fake_protocol",
		}),
	}

	var existSlbListenersTcpMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"load_balancer_id": CHECKSET,
			"slb_listeners.#":  "1",
		}
	}

	var fakeSlbListenersTcpMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"load_balancer_id": CHECKSET,
			"slb_listeners.#":  "0",
		}
	}

	var slbListenersTcpCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existSlbListenersTcpMapFunc,
		fakeMapFunc:  fakeSlbListenersTcpMapFunc,
	}
	slbListenersTcpCheckInfo.dataSourceTestCheck(t, rand, lbIdConf, protocolConf, bothFiltersConf)
}

func TestAccAlibabacloudStackSlbListenersDataSource_udp(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_slb_listeners.default"
	name := fmt.Sprintf("tf-testacc-slblisteners%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSlbListenersUdpConfigDependence)

	// Test with load balancer ID filter
	lbIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}_fake",
		}),
	}

	// Test with protocol filter
	protocolConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
			"protocol":         "udp",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
			"protocol":         "fake_protocol",
		}),
	}

	// Test with both filters
	bothFiltersConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}",
			"protocol":         "udp",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_listener.default.load_balancer_id}_fake",
			"protocol":         "fake_protocol",
		}),
	}

	var existSlbListenersUdpMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"load_balancer_id": CHECKSET,
			"slb_listeners.#":  "1",
		}
	}

	var fakeSlbListenersUdpMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"load_balancer_id": CHECKSET,
			"slb_listeners.#":  "0",
		}
	}

	var slbListenersUdpCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existSlbListenersUdpMapFunc,
		fakeMapFunc:  fakeSlbListenersUdpMapFunc,
	}
	slbListenersUdpCheckInfo.dataSourceTestCheck(t, rand, lbIdConf, protocolConf, bothFiltersConf)
}

func dataSourceSlbListenersHttpConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
	data "alibabacloudstack_zones" "default" {
		available_resource_creation= "VSwitch"
	}
	
	resource "alibabacloudstack_slb" "default" {
		name = "${var.name}"
	}
	
	resource "alibabacloudstack_slb_listener" "default" {
		load_balancer_id = "${alibabacloudstack_slb.default.id}"
		backend_port = 80
		frontend_port = 80
		protocol = "http"
		sticky_session = "on"
		sticky_session_type = "insert"
		cookie = "${var.name}"
		cookie_timeout = 86400
		health_check = "on"
		health_check_uri = "/cons"
		health_check_connect_port = 20
		healthy_threshold = 8
		unhealthy_threshold = 8
		health_check_timeout = 8
		health_check_interval = 5
		health_check_http_code = "http_2xx,http_3xx"
		bandwidth = 10
		x_forwarded_for  {
			retrive_slb_ip = true
			retrive_slb_id = true
		}
		description = "${var.name}"
	}
	`, name)
}

func dataSourceSlbListenersHttpsConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
	resource "alibabacloudstack_slb" "default" {
		name = "${var.name}"
	}
	
	resource "alibabacloudstack_slb_server_certificate" "default" {
		name = "${var.name}"
		server_certificate = %s
		private_key = %s
	}
	
	resource "alibabacloudstack_slb_listener" "default" {
		load_balancer_id = "${alibabacloudstack_slb.default.id}"
		backend_port = 80
		frontend_port = 80
		protocol = "https"
		sticky_session = "on"
		sticky_session_type = "insert"
		cookie = "${var.name}"
		cookie_timeout = 86400
		health_check = "on"
		health_check_uri = "/cons"
		health_check_connect_port = 20
		healthy_threshold = 8
		unhealthy_threshold = 8
		health_check_timeout = 8
		health_check_interval = 5
		health_check_http_code = "http_2xx,http_3xx"
		bandwidth = 10
		x_forwarded_for  {
			retrive_slb_ip = true
			retrive_slb_id = true
		}
		server_certificate_id = "${alibabacloudstack_slb_server_certificate.default.id}"
		description = "${var.name}"
	}
	`, name, ServerCertificateTestCase(), RsaPrivateKeyTestCase())
}

func dataSourceSlbListenersTcpConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
	data "alibabacloudstack_zones" "default" {
		available_resource_creation= "VSwitch"
	}
	
	resource "alibabacloudstack_slb" "default" {
		name = "${var.name}"
	}
	
	resource "alibabacloudstack_slb_listener" "default" {
		load_balancer_id = "${alibabacloudstack_slb.default.id}"
		backend_port = 22
		frontend_port = 22
		protocol = "tcp"
		health_check_connect_port = 20
		healthy_threshold = 8
		unhealthy_threshold = 8
		health_check_timeout = 8
		health_check_interval = 5
		health_check_type = "tcp"
		health_check = "off"
		sticky_session = "off"
		bandwidth = 10
		description = "${var.name}"
	}
	`, name)
}

func dataSourceSlbListenersUdpConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
	data "alibabacloudstack_zones" "default" {
		available_resource_creation= "VSwitch"
	}
	
	resource "alibabacloudstack_slb" "default" {
		name = "${var.name}"
	}
	
	resource "alibabacloudstack_slb_listener" "default" {
		load_balancer_id = "${alibabacloudstack_slb.default.id}"
		backend_port = 10
		frontend_port = 11
		protocol = "udp"
		health_check_connect_port = 20
		healthy_threshold = 8
		unhealthy_threshold = 8
		health_check_timeout = 8
		health_check_interval = 5
		bandwidth = 10
		health_check = "off"
		sticky_session = "off"
		description = "${var.name}"
	}
	`, name)
}
