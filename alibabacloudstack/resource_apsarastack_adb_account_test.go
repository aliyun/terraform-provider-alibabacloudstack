package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAdbAccount0(t *testing.T) {
	var v *adb.DBAccount

	resourceId := "alibabacloudstack_adb_account.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccAdbAccountCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoAdbDescribeaccountsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tfaccount%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccAdbAccountBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":    "${local.adb_instance_id}",
					"account_name":     "${var.name}",
					"account_password": "${random_password.password.0.result}",
					"account_description": "${var.name}_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name":  name,
						"account_description": name+"_desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password":       REMOVEKEY,
					"kms_encrypted_password": "${random_password.password.0.result}",
					"kms_encryption_context":  "terraform-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "kms_encrypted_password", "kms_encryption_context"},
			},
		},
	})
}

var AlibabacloudTestAccAdbAccountCheckmap = map[string]string{
	"account_description": CHECKSET,
	"status":              CHECKSET,
	"db_cluster_id":       CHECKSET,
	"account_type":        CHECKSET,
	"account_name":        CHECKSET,
}

func AlibabacloudTestAccAdbAccountBasicdependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
			default = "%s"
		}

		%s
		
		%s
		
	`, name, AdbCommonTestCase(false), RandomPasswordTestCase(12, 1))
}
