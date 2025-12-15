package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccAlibabacloudStackHsmInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_hsm_instance.default"

	ra := resourceAttrInit(resourceId, map[string]string{})
	commonProvider := Provider()
	yundunProvider := Provider()
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HsmService{yundunProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeHsmInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf_hsm_instance%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceHsmInstanceDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreYunCheck(t)
		},
		IDRefreshName: resourceId,
		Providers: func() map[string]*schema.Provider {
			yundunProvider.Schema["access_key"] = &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_YUNDUN_ACCESS_KEY", ""),
				Description: descriptions["access_key"],
			}
			yundunProvider.Schema["secret_key"] = &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_YUNDUN_SECRET_KEY", ""),
				Description: descriptions["secret_key"],
			}
			yundunProvider.Schema["role_arn"] = &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["assume_role_role_arn"],
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_YUNDUN_ASSUME_ROLE_ARN", ""),
			}
			return map[string]*schema.Provider{
				"alibabacloudstack":        yundunProvider,
				"alibabacloudstack-common": commonProvider,
			}
		}(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"product_code": "${data.alibabacloudstack_hsm_vendors.default.vendors.0.products.0.code}",
					"vendor_code":  "${data.alibabacloudstack_hsm_vendors.default.vendors.0.code}",
					"vsm_type":     "gvsm",
					"zone_no":      "${data.alibabacloudstack_zones.default.zones.0.id}",
					"remark":       "test-tf-hsm-instance",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_code": CHECKSET,
						"vendor_code":  CHECKSET,
						"vsm_type":     "gvsm",
						"zone_no":      CHECKSET,
						"remark":       "test-tf-hsm-instance",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "updated_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "updated_remark",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":     "${alibabacloudstack_vpc.vpc.id}",
					"vswitch_id": "${alibabacloudstack_vswitch.vsw.id}",
					"ip":         "192.168.0.100",
					"white_list": "192.168.0.123/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":     CHECKSET,
						"vswitch_id": CHECKSET,
						"ip":         "192.168.0.100",
						"white_list": "192.168.0.123/24",
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

func resourceHsmInstanceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

data "alibabacloudstack_zones" "default" {
	provider = alibabacloudstack-common
	available_resource_creation = "VSwitch"
}

data "alibabacloudstack_hsm_vendors" "default" {
}

resource "alibabacloudstack_vpc" "vpc" {	
	provider = alibabacloudstack-common
	vpc_name = var.name
	cidr_block = "192.168.0.0/16" # VPC CIDR block
}

resource "alibabacloudstack_vswitch" "vsw" {
	provider = alibabacloudstack-common
	vpc_id = alibabacloudstack_vpc.vpc.id
	cidr_block = "192.168.0.0/24" # VSwitch CIDR block
	availability_zone = data.alibabacloudstack_zones.default.zones.0.id # Availability zone
}

`, name)
}
