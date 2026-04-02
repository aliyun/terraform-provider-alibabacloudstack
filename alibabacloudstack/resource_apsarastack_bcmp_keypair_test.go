package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackBcmpKeyPair0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_bcmp_keypair.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccBcmpKeypairCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &BcmpService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEasyAIListKeyPairRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbcmpkey_pair%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccBcmpKeypairBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"key_pair_name": name,
					"key_file":      "./test-key.txt",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"key_pair_name": name,
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_file"},
			},
		},
	})
}

func TestAccAlibabacloudStackBcmpKeyPair_PublicKey(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_bcmp_keypair.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccBcmpKeypairCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &BcmpService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEasyAIListKeyPairRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbcmpkey_pair%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccBcmpKeypairBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"key_pair_name": name,
					"public_key":    "${alibabacloudstack_bcmp_keypair.old.public_key}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"key_pair_name": name,
						"finger_print":  CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"public_key", "key_file"},
			},
		},
	})
}

var AlibabacloudTestAccBcmpKeypairCheckmap = map[string]string{

	"key_pair_name": CHECKSET,

	"finger_print": CHECKSET,
}

func AlibabacloudTestAccBcmpKeypairBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alibabacloudstack_bcmp_keypair" "old"{
	key_pair_name = "${var.name}_old"
}

`, name)
}
