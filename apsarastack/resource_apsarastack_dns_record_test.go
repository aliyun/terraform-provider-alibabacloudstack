package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func (rc *resourceCheck) checkResourceDnsRecordDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		strs := strings.Split(rc.resourceId, ":")
		var resourceType string
		for _, str := range strs {
			if strings.Contains(str, "apsarastack_") {
				resourceType = strings.Trim(str, " ")
				break
			}
		}

		if resourceType == "" {
			return WrapError(Error("The resourceId %s is not correct and it should prefix with apsarastack_", rc.resourceId))
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}
			outValue, err := rc.callDescribeMethod(rs)
			errorValue := outValue[1]
			if !errorValue.IsNil() {
				err = errorValue.Interface().(error)
				if err != nil {
					if NotFoundError(err) {
						continue
					}
					return WrapError(err)
				}
			} else {
				return WrapError(Error("the resource %s %s was not destroyed ! ", rc.resourceId, rs.Primary.ID))
			}
		}
		return nil
	}
}

func TestAccApsaraStackDnsRecord_basic(t *testing.T) {
	var v *DnsRecord
	resourceId := "apsarastack_dns_record.default"
	ra := resourceAttrInit(resourceId, dnsRecordBasicMap)
	serviceFunc := func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprint("tf-testdnsrecordbasic.")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccDnsRecordConfigBasicConfigBasic)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDnsRecordDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_id":   "${apsarastack_dns_domain.default.domain_id}",
					"host_record": "test",
					"type":        "A",
					"ttl":         "0",
					"rr_set":      []string{"10.0.0.1", "10.0.0.3", "10.0.0.2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccDnsRecordConfigBasicConfigBasic(name string) string {
	return fmt.Sprintf(`

resource "apsarastack_dns_domain" "default" {
 domain_name = "%s"
 remark = "test_dummy_1"
}
`, name)
}

var dnsRecordBasicMap = map[string]string{
	"domain_id":   CHECKSET,
	"host_record": CHECKSET,
	"type":        CHECKSET,
}
