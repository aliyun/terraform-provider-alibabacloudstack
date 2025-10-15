package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAckTemplate0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_ack_template.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccAckTemplateCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeAckTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testaccAckTemplate%d", rand)
	modify_name := fmt.Sprintf("tf-testaccAckTemplatemodify%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccAckTemplateBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"template": "${var.template_content}",
					"description":   name,
					"name":          name,
					"template_type": "kubernetes",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template": CHECKSET,
						"description":   name,
						"name":          name,
						"template_type": "kubernetes",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"description": modify_name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": modify_name,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"template": "${var.template_content_update}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"template": CHECEKSET,
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

var AlibabacloudTestAccAckTemplateCheckmap = map[string]string{

	"template_type": "kubernetes",

	"template_id": CHECKSET,

	"template_with_hist_id": CHECKSET,

	"template_hash_code_version": CHECKSET,
	"created":                    CHECKSET,

	"acl": CHECKSET,

	"version": CHECKSET,
	"ali_uid": CHECKSET,

	"updated": CHECKSET,
}

func AlibabacloudTestAccAckTemplateBasicdependence(name string) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "%s"
}

variable "template_content" {
	default = <<EOF
	apiVersion: apps/v1
	kind: Deployment
	metadata:
	labels:
		vsw: test
	name: nginx-deployment-basic
	namespace: default
	spec:
	replicas: 1
	selector:
		matchLabels:
		vsw: test
	template:
		metadata:
		labels:
			vsw: test
		spec:
		containers:
			- command:
				- sleep
				- '10000001'
			image: >-
				registry.acs.%s/acs/busybox:1.33.1
			imagePullPolicy: IfNotPresent
			name: vsw
	EOF
}

variable "template_content_update" {
	default = <<EOF
	apiVersion: apps/v1
	kind: Deployment
	metadata:
	labels:
		vsw: test
	name: nginx-deployment-basic
	namespace: default
	spec:
	replicas: 1
	selector:
		matchLabels:
		vsw: test
	template:
		metadata:
		labels:
			vsw: test
		spec:
		containers:
			- command:
				- sleep
				- '10000002'
			image: >-
				registry.acs.%s/acs/busybox:1.33.1
			imagePullPolicy: IfNotPresent
			name: vsw
	EOF
}
`, name, os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN"), os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN"))
}
