package alibabacloudstack

import (
	"fmt"
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
	ra := resourceAttrInit(resourceId, map[string]string{
		"domain_name": CHECKSET,
	})
	serviceFunc := func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tfacc%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlibabacloudStackDnsDomainDepend)
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
				Config: testAccConfig(map[string]interface{}{
					"domain_name": "${var.name}.",
					"remark":      "${var.name}_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name": name + ".",
						"remark":      name + "_remark",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "${var.name}_remark1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name": name + ".",
						"remark":      name + "_remark1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func resourceAlibabacloudStackDnsDomainDepend(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}
