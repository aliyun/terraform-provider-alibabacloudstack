package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackSlbAccessLogsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_slb_access_logs.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%s-%d", defaultRegionToTest, rand),
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
			"access_logs.0.log_store":        "testtf1",
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
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key        = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
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

// resource "alibabacloudstack_log_project" "default" {
// 	name = "${var.name}_project"
// 	description = "test"
// }

// resource "alibabacloudstack_log_store" "default" {
// 	name = "${var.name}_store"
// 	project = "${alibabacloudstack_log_project.default.name}"	
// 	retention_period      = "30"
// 	shard_count           = "2"
// 	enable_web_tracking   = false
// 	auto_split            = true
// 	max_split_shard_count = "64"
// 	append_meta           = true
// }

resource "alibabacloudstack_slb_access_log" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  log_project = "testtf"
  log_store = "testtf1"
}
`, name, ECSInstanceCommonTestCase)
}
