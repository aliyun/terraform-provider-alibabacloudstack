package alibabacloudstack

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func (rc *resourceCheck) checkResourceDnsDomainDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		strs := strings.Split(rc.resourceId, ":")
		var resourceType string
		for _, str := range strs {
			if strings.Contains(str, "alibabacloudstack_") {
				resourceType = strings.Trim(str, " ")
				break
			}
		}

		if resourceType == "" {
			return errmsgs.WrapError(errmsgs.Error("The resourceId %s is not correct and it should prefix with alibabacloudstack_", rc.resourceId))
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
					if errmsgs.NotFoundError(err) {
						continue
					}
					return errmsgs.WrapError(err)
				}
			} else {
				return errmsgs.WrapError(errmsgs.Error("the resource %s %s was not destroyed ! ", rc.resourceId, rs.Primary.ID))
			}
		}
		return nil
	}
}

func TestAccAlibabacloudStackDnsDomain_basic(t *testing.T) {
	var v *DnsDomains
	resourceId := "alibabacloudstack_dns_domain.default"
	ra := resourceAttrInit(resourceId, dnsDomainBasicMap)
	serviceFunc := func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: resourceAlibabacloudStackDns_Domain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: resourceAlibabacloudStackDns_Domain2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

const resourceAlibabacloudStackDns_Domain = `
resource "alibabacloudstack_dns_domain" "default" {
	domain_name = "testdummy."

}
`
const resourceAlibabacloudStackDns_Domain2 = `
resource "alibabacloudstack_dns_domain" "default" {
	domain_name = "testdummy."
     remark = "test_dummy_1"

}
`

var dnsDomainBasicMap = map[string]string{
	"domain_name": CHECKSET,
}
