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

func (rc *resourceCheck) checkResourceOnsTopicDestroy() resource.TestCheckFunc {
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

func TestAccAlibabacloudStackOnsTopic_basic(t *testing.T) {
	var v *Topic
	resourceId := "alibabacloudstack_ons_topic.default"
	ra := resourceAttrInit(resourceId, onsTopicBasicMap)
	serviceFunc := func() interface{} {
		return &OnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000,20000)
	name := fmt.Sprintf("tf-testacconstopicbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccOnsTopicConfigBasic)

	ResourceTest(t, resource.TestCase{
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
					"instance_id":  "${alibabacloudstack_ons_instance.default.id}",
					"topic":        name,
					"remark":       "Ons_topic",
					"message_type": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"perm"},
			},
		},
	})

}

func testAccOnsTopicConfigBasic(name string) string {
	return fmt.Sprintf(`

variable "topic" {
 default = "%s"
}

resource "alibabacloudstack_ons_instance" "default" {
  tps_receive_max = 500
  tps_send_max = 500
  topic_capacity = 50
  cluster = "cluster1"
  independent_naming = "true"
  name = "${var.topic}"
  remark = "Ons_instance"
}

`, name)
}

var onsTopicBasicMap = map[string]string{
	"instance_id":  CHECKSET,
	"topic":        CHECKSET,
	"message_type": CHECKSET,
	"remark":       CHECKSET,
}
