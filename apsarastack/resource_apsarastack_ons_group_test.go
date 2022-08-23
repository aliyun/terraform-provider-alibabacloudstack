package apsarastack

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"
)

func (rc *resourceCheck) checkResourceOnsGroupDestroy() resource.TestCheckFunc {
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

func TestAccApsaraStackOnsGroup_basic(t *testing.T) {
	var v *OnsGroup
	resourceId := "apsarastack_ons_group.default"
	ra := resourceAttrInit(resourceId, onsGroupBasicMap)
	serviceFunc := func() interface{} {
		return &OnsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("GID-tf-testacconsgroupbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccOnsGroupConfigBasic)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceOnsGroupDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${apsarastack_ons_instance.default.id}",
					"group_id":    name,
					"remark":      "Ons_Group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"read_enable"},
			},
		},
	})

}

func testAccOnsGroupConfigBasic(name string) string {
	return fmt.Sprintf(`

variable "group_id" {
 default = "%s"
}

resource "apsarastack_ons_instance" "default" {
  tps_receive_max = 500
  tps_send_max = 500
  topic_capacity = 50
  cluster = "cluster1"
  independent_naming = "true"
  name = "${var.group_id}"
  remark = "Ons_instance"
}

`, name)
}

var onsGroupBasicMap = map[string]string{
	"instance_id": CHECKSET,
	"group_id":    CHECKSET,
	"remark":      CHECKSET,
}
