package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAlikafkaSaslAcl_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_alikafka_sasl_acl.default"
	ra := resourceAttrInit(resourceId, alikafkaSaslAclBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(10000, 20000)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-kafkasaslacl%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaSaslAclConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":               "${local.alikafka_instnace_id}",
					"username":                  "${alibabacloudstack_alikafka_sasl_user.default.username}",
					"acl_resource_type":         "Topic",
					"acl_resource_name":         "${alibabacloudstack_alikafka_topic.default.topic}",
					"acl_resource_pattern_type": "LITERAL",
					"acl_operation_type":        "WRITE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"username":          name,
						"acl_resource_name": name,
					}),
				),
			},

			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"acl_resource_pattern_type": "PREFIXED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_resource_pattern_type": "PREFIXED",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"acl_operation_type": "READ",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_operation_type": "READ",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_resource_type": "Group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_resource_type": "Group",
					}),
				),
			},
		},
	})

}

func resourceAlikafkaSaslAclConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

%s

%s

%s

resource "alibabacloudstack_alikafka_topic" "default" {
  instance_id = local.alikafka_instnace_id
  topic = "${var.name}"
  remark = "topic-remark"
}

resource "alibabacloudstack_alikafka_sasl_user" "default" {
  instance_id = local.alikafka_instnace_id
  username = "${var.name}"
  password = random_password.password.0.result
  type     = "scram"
}
`, name, DataZoneCommonTestCase, KafkaCommonTestCase(), RandomPasswordTestCase(12, 1))
}

var alikafkaSaslAclBasicMap = map[string]string{
	"username":                  "${var.name}",
	"acl_resource_type":         "Topic",
	"acl_resource_name":         "${var.name}",
	"acl_resource_pattern_type": "LITERAL",
	"acl_operation_type":        "WRITE",
}
