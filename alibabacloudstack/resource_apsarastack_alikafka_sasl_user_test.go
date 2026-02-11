package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAlikafkaSaslUser_basic(t *testing.T) {

	//var v *alikafka.SaslUserVO
	var v map[string]interface{}
	resourceId := "alibabacloudstack_alikafka_sasl_user.default"
	ra := resourceAttrInit(resourceId, alikafkaSaslUserBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(10000, 20000)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkasasluserbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaSaslUserConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${local.alikafka_instnace_id}",
					"username":    "${var.name}",
					"password":    "${random_password.password.0.result}",
					"type":        "scram",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"username": name,
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
					"username": name + "_new",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"username": name + "_new",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"password": "${random_password.password.1.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})

}

func resourceAlikafkaSaslUserConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}
%s

%s

%s
`, name, DataZoneCommonTestCase, KafkaCommonTestCase(), RandomPasswordTestCase(12, 2))
}

var alikafkaSaslUserBasicMap = map[string]string{
	"username": "${var.name}",
}
