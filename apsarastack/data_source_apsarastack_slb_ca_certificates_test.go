package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"strings"
	"testing"
)

func TestAccApsaraStackSlbCACertificatesDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbCaCertificatesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_slb_ca_certificate.default.name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbCaCertificatesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_slb_ca_certificate.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbCaCertificatesDataSourceConfig(rand, map[string]string{
			"ids": `["${apsarastack_slb_ca_certificate.default.id}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbCaCertificatesDataSourceConfig(rand, map[string]string{
			"ids": `["${apsarastack_slb_ca_certificate.default.id}_fake"]`,
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbCaCertificatesDataSourceConfig(rand, map[string]string{
			"ids": `["${apsarastack_slb_ca_certificate.default.id}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbCaCertificatesDataSourceConfig(rand, map[string]string{
			"ids": `["${apsarastack_slb_ca_certificate.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbCaCertificatesDataSourceConfig(rand, map[string]string{
			"ids":        `["${apsarastack_slb_ca_certificate.default.id}"]`,
			"name_regex": `"${apsarastack_slb_ca_certificate.default.name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbCaCertificatesDataSourceConfig(rand, map[string]string{
			"ids":        `["${apsarastack_slb_ca_certificate.default.id}_fake"]`,
			"name_regex": `"${apsarastack_slb_ca_certificate.default.name}"`,
		}),
	}

	var existDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"certificates.#":                   "1",
			"ids.#":                            "1",
			"names.#":                          "1",
			"certificates.0.id":                CHECKSET,
			"certificates.0.name":              fmt.Sprintf("tf-testAccSlbCACertificatesDataSourceBasic-%d", rand),
			"certificates.0.fingerprint":       CHECKSET,
			"certificates.0.created_timestamp": CHECKSET,
			"certificates.0.region_id":         defaultRegionToTest,
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"certificates.#": "0",
			"ids.#":          "0",
			"names.#":        "0",
		}
	}

	var slbCaCertificatesCheckInfo = dataSourceAttr{
		resourceId:   "data.apsarastack_slb_ca_certificates.default",
		existMapFunc: existDnsRecordsMapFunc,
		fakeMapFunc:  fakeDnsRecordsMapFunc,
	}

	slbCaCertificatesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, resourceGroupIdConf, allConf)

}

func testAccCheckApsaraStackSlbCaCertificatesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSlbCACertificatesDataSourceBasic-%d"
}


resource "apsarastack_slb_ca_certificate" "default" {
  name = "${var.name}"
  ca_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
}

data "apsarastack_slb_ca_certificates" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
