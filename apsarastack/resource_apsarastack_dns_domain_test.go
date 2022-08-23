package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func (rc *resourceCheck) checkResourceDnsDomainDestroy() resource.TestCheckFunc {
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

func TestAccApsaraStackDnsDomain_basic(t *testing.T) {
	var v *DnsDomains
	resourceId := "apsarastack_dns_domain.default"
	ra := resourceAttrInit(resourceId, dnsDomainBasicMap)
	serviceFunc := func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: resourceApsaraStackDns_Domain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

const resourceApsaraStackDns_Domain = `
resource "apsarastack_dns_domain" "default" {
	domain_name = "testdummy."
	remark = "test_dummy_1"
}
`

var dnsDomainBasicMap = map[string]string{
	"domain_name": CHECKSET,
	"remark":      CHECKSET,
}
