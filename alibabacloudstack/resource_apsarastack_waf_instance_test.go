package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccAlibabacloudStackWafInstance_basic(t *testing.T) {
	resourceId := "alibabacloudstack_waf_instance.default"
	var v map[string]interface{}
	ra := resourceAttrInit(resourceId, WafInstanceBasicMap)
	serviceFunc := func() interface{} {
		return &WafOpenapiService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf_testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceWafInstanceDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreYunCheck(t)
		},
		IDRefreshName: resourceId,
		Providers: func() map[string]*schema.Provider {
			commonProvider := Provider()
			yundunProvider := Provider()
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
					"vswitch_id":       "${alibabacloudstack_vswitch.vsw.id}",
					"vpc_id":           "${alibabacloudstack_vpc.vpc.id}",
					"name":             "waf_instance_test",
					"detector_specs":   "exclusive",
					"detector_version": "basic",
					"detector_nodenum": 2,
					"vpc_vswitch": []map[string]interface{}{
						{

							"vswitch_name":   "${alibabacloudstack_vswitch.vsw.vswitch_name}",
							"vswitch":        "${alibabacloudstack_vswitch.vsw.id}",
							"cidr_block":     "${alibabacloudstack_vswitch.vsw.cidr_block}",
							"available_zone": "${alibabacloudstack_vswitch.vsw.zone_id}",
							"vpc":            "${alibabacloudstack_vpc.vpc.id}",
							"vpc_name":       "${alibabacloudstack_vpc.vpc.vpc_name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":       CHECKSET,
						"name":             "waf_instance_test",
						"detector_specs":   "exclusive",
						"detector_version": "basic",
						"vpc_id":           CHECKSET,
						"detector_nodenum": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detector_nodenum": 3,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detector_nodenum": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detector_nodenum": 2,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detector_nodenum": "2",
					}),
				),
			},
		},
	})
}

func resourceWafInstanceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_zones" "default" {
	provider = alibabacloudstack-common
	available_resource_creation = "VSwitch"
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

var WafInstanceBasicMap = map[string]string{
	// "description":          CHECKSET,
	// "license_code":         "bhah_ent_50_asset",
	// "period":               "1",
	// "vswitch_id":           CHECKSET,
	// "security_group_ids.#": "1",
}
