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
			"load_balancer_id": "${alibabacloudstack_slb.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "fake-lb-id",
		}),
	}

	// Test with protocol filter
	protocolConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb.default.id}",
			"protocol":         "http",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb.default.id}",
			"protocol":         "fake_protocol",
		}),
	}

	// Test with both filters
	bothFiltersConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb.default.id}",
			"protocol":         "http",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "fake-lb-id",
			"protocol":         "fake_protocol",
		}),
	}

	var existSlbListenersHttpMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"load_balancer_id":                             CHECKSET,
			"listeners.#":                                  "1",
			"listeners.0.load_balancer_id":                 CHECKSET,
			"listeners.0.listener_port":                    "80",
			"listeners.0.backend_port":                     "80",
			"listeners.0.protocol":                         "http",
			"listeners.0.bandwidth":                        "10",
			"listeners.0.scheduler":                        "wrr",
			"listeners.0.sticky_session":                   "on",
			"listeners.0.sticky_session_type":              "insert",
			"listeners.0.health_check":                     "on",
			"listeners.0.health_check_uri":                 "/cons",
			"listeners.0.description":                      name,
			"listeners.0.cookie_timeout":                   "86400",
			"listeners.0.cookie":                           CHECKSET,
			"listeners.0.x_forwarded_for.0.retrive_slb_ip": "true",
			"listeners.0.x_forwarded_for.0.retrive_slb_id": "true",
		}
	}

	var fakeSlbListenersHttpMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"load_balancer_id": CHECKSET,
			"listeners.#":      "0",
		}
	}

	var slbListenersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existSlbListenersHttpMapFunc,
		fakeMapFunc:  fakeSlbListenersHttpMapFunc,
	}
	slbListenersCheckInfo.dataSourceTestCheck(t, rand, lbIdConf, protocolConf, bothFiltersConf)
}

func TestAccAlibabacloudStackSlbListenersDataSource_https(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_slb_listeners.default"
	name := fmt.Sprintf("tf-testacc-slblisteners%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSlbListenersHttpsConfigDependence)

	// Test with load balancer ID filter
	lbIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "fake-lb-id",
		}),
	}

	// Test with protocol filter
	protocolConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb.default.id}",
			"protocol":         "https",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb.default.id}",
			"protocol":         "fake_protocol",
		}),
	}

	// Test with both filters
	bothFiltersConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb.default.id}",
			"protocol":         "https",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "fake-lb-id",
			"protocol":         "fake_protocol",
		}),
	}

	var existSlbListenersHttpsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"load_balancer_id":                             CHECKSET,
			"listeners.#":                                  "1",
			"listeners.0.load_balancer_id":                 CHECKSET,
			"listeners.0.listener_port":                    "80",
			"listeners.0.backend_port":                     "80",
			"listeners.0.protocol":                         "https",
			"listeners.0.bandwidth":                        "10",
			"listeners.0.scheduler":                        "wrr",
			"listeners.0.sticky_session":                   "on",
			"listeners.0.sticky_session_type":              "insert",
			"listeners.0.health_check":                     "on",
			"listeners.0.health_check_uri":                 "/cons",
			"listeners.0.description":                      name,
			"listeners.0.cookie_timeout":                   "86400",
			"listeners.0.cookie":                           CHECKSET,
			"listeners.0.x_forwarded_for.0.retrive_slb_ip": "true",
			"listeners.0.x_forwarded_for.0.retrive_slb_id": "true",
			"listeners.0.server_certificate_id":            CHECKSET,
			"listeners.0.enable_http2":                     "on",
			"listeners.0.tls_cipher_policy":                "tls_cipher_policy_1_0",
		}
	}

	var fakeSlbListenersHttpsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"load_balancer_id": CHECKSET,
			"listeners.#":      "0",
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
			"load_balancer_id": "${alibabacloudstack_slb.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "fake-lb-id",
		}),
	}

	// Test with protocol filter
	protocolConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb.default.id}",
			"protocol":         "tcp",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb.default.id}",
			"protocol":         "fake_protocol",
		}),
	}

	// Test with both filters
	bothFiltersConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb.default.id}",
			"protocol":         "tcp",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "fake-lb-id",
			"protocol":         "fake_protocol",
		}),
	}

	var existSlbListenersTcpMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"load_balancer_id":                CHECKSET,
			"listeners.#":                     "1",
			"listeners.0.load_balancer_id":    CHECKSET,
			"listeners.0.listener_port":       "22",
			"listeners.0.backend_port":        "22",
			"listeners.0.protocol":            "tcp",
			"listeners.0.bandwidth":           "10",
			"listeners.0.scheduler":           "wrr",
			"listeners.0.sticky_session":      "off",
			"listeners.0.health_check":        "off",
			"listeners.0.health_check_type":   "tcp",
			"listeners.0.description":         name,
			"listeners.0.persistence_timeout": "0",
		}
	}

	var fakeSlbListenersTcpMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"load_balancer_id": CHECKSET,
			"listeners.#":      "0",
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
			"load_balancer_id": "${alibabacloudstack_slb.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "fake-lb-id",
		}),
	}

	// Test with protocol filter
	protocolConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb.default.id}",
			"protocol":         "udp",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb.default.id}",
			"protocol":         "fake_protocol",
		}),
	}

	// Test with both filters
	bothFiltersConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb.default.id}",
			"protocol":         "udp",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "fake-lb-id",
			"protocol":         "fake_protocol",
		}),
	}

	var existSlbListenersUdpMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"load_balancer_id":             CHECKSET,
			"listeners.#":                  "1",
			"listeners.0.load_balancer_id": CHECKSET,
			"listeners.0.listener_port":    "11",
			"listeners.0.backend_port":     "10",
			"listeners.0.protocol":         "udp",
			"listeners.0.bandwidth":        "10",
			"listeners.0.scheduler":        "wrr",
			"listeners.0.sticky_session":   "off",
			"listeners.0.health_check":     "off",
			"listeners.0.description":      name,
		}
	}

	var fakeSlbListenersUdpMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"load_balancer_id": CHECKSET,
			"listeners.#":      "0",
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
		server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
		private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
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
	`, name)
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
