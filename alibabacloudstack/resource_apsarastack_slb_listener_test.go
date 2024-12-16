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
	name := fmt.Sprintf("tf-testacc%sslblistener%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbListenerBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"load_balancer_id": "${alibabacloudstack_slb.default.id}",
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
	server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
	private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
  }

resource "alibabacloudstack_slb_acl" "default" {
	name = "${var.name}"
	ip_version = "ipv4"
  }

resource "alibabacloudstack_slb" "default" {
	name = "${var.name}"
	// vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	address_type       = "internet"
  	specification        = "slb.s2.small"
  }

`, name)
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
	name := fmt.Sprintf("tf-testacc%sslblistener%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbListenerBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"load_balancer_id": "${alibabacloudstack_slb.default.id}",
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
	name := fmt.Sprintf("tf-testacc%sslblistener%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbListenerBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alibabacloudstack_slb.default.id}",
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
	name := fmt.Sprintf("tf-testacc%sslblistener%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbListenerBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"load_balancer_id": "${alibabacloudstack_slb.default.id}",

					"protocol":      "http",
					"bandwidth":     "10",
					"frontend_port": "80",
					"backend_port":  "80",

					"acl_status":     "on",
					"sticky_session": "off",
					"health_check":   "off",

					"acl_type": "white",

					"description": "testcreate",

					"scheduler": "wrr",

					"acl_id": "${alibabacloudstack_slb_acl.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"load_balancer_id": CHECKSET,

						"listener_protocol": "http",
						"bandwidth":         "10",
						"frontend_port":     "80",
						"backend_port":      "80",

						"acl_status": "on",

						"acl_type": "white",

						"description": "testcreate",

						"scheduler": "wrr",

						"acl_id": CHECKSET,
					}),
				),
			},
		},
	})
}
