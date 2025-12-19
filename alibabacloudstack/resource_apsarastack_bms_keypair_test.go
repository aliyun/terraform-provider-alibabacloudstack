package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackBmsKeypair_basic(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_bms_keypair.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &BmsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeKeyPair")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("tf-bms-keypair%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, BmsKeypairDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key", "public_key"},
			},
		},
	})
}

func TestAccAlibabacloudStackBmsKeypair_PublicKey(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_bms_keypair.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &BmsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeKeyPair")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("tf-bms-keypair%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, BmsKeypairDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":       "${var.name}",
					"public_key": "${var.public_key}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key", "public_key"},
			},
		},
	})
}

func BmsKeypairDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "public_key" {
  default = "%s"
}
`, name, RsaPublicKeyTestCase())
}
