package alibabacloudstack

import (
	"fmt"
	"testing"

	

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackKmsAlias_basic(t *testing.T) {
	var v kms.KeyMetadata

	resourceId := "alibabacloudstack_kms_alias.default"
	ra := resourceAttrInit(resourceId, kmsAliasBasicMap)

	serviceFunc := func() interface{} {
		return &KmsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}

	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("alias/tf-testKmsAlias_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceKmsAliadConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alias_name": name,
					"key_id":     "${alibabacloudstack_kms_key.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alias_name": name,
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
					"key_id": "${alibabacloudstack_kms_key.default1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_id": CHECKSET,
					}),
				),
			},
		},
	})
}

func resourceKmsAliadConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alibabacloudstack_kms_key" "default" {}

resource "alibabacloudstack_kms_key" "default1" {}
`)
}

var kmsAliasBasicMap = map[string]string{
	"key_id": CHECKSET,
}
