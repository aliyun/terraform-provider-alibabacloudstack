package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackOtsInstanceAttachmentBasic(t *testing.T) {
	var v ots.VpcInfo

	resourceId := "alibabacloudstack_ots_instance_attachment.default"
	ra := resourceAttrInit(resourceId, otsInstanceAttachmentBasicMap)

	serviceFunc := func() interface{} {
		return &OtsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOtsInstanceAttachmentConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alibabacloudstack_ots_instance.default.name}",
					"vpc_name":      "test",
					"vswitch_id":    "${alibabacloudstack_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"vpc_name":      "test",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackOtsInstanceAttachmentHighPerformance(t *testing.T) {
	var v ots.VpcInfo

	resourceId := "alibabacloudstack_ots_instance_attachment.default"
	ra := resourceAttrInit(resourceId, otsInstanceAttachmentBasicMap)

	serviceFunc := func() interface{} {
		return &OtsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOtsInstanceAttachmentConfigDependenceHighperformance)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alibabacloudstack_ots_instance.default.name}",
					"vpc_name":      "test",
					"vswitch_id":    "${alibabacloudstack_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"vpc_name":      "test",
					}),
				),
			},
		},
	})
}

func resourceOtsInstanceAttachmentConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}

	resource "alibabacloudstack_ots_instance" "default" {
	  name = "${var.name}"
	  description = "${var.name}"
	  accessed_by = "Vpc"
	  instance_type = "%s"
	}

	data "alibabacloudstack_zones" "default" {
	  available_resource_creation = "VSwitch"
	}
	resource "alibabacloudstack_vpc" "default" {
	  cidr_block = "172.16.0.0/16"
	  name = "${var.name}"
	}

	resource "alibabacloudstack_vswitch" "default" {
	  vpc_id = "${alibabacloudstack_vpc.default.id}"
	  name = "${var.name}"
	  cidr_block = "172.16.1.0/24"
	  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	}
	`, name, string(OtsCapacity))
}

func resourceOtsInstanceAttachmentConfigDependenceHighperformance(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}

	resource "alibabacloudstack_ots_instance" "default" {
	  name = "${var.name}"
	  description = "${var.name}"
	  accessed_by = "Vpc"
	  instance_type = "%s"
	}

	data "alibabacloudstack_zones" "default" {
	  available_resource_creation = "VSwitch"
	}
	resource "alibabacloudstack_vpc" "default" {
	  cidr_block = "172.16.0.0/16"
	  name = "${var.name}"
	}

	resource "alibabacloudstack_vswitch" "default" {
	  vpc_id = "${alibabacloudstack_vpc.default.id}"
	  name = "${var.name}"
	  cidr_block = "172.16.1.0/24"
	  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	}
	`, name, string(OtsHighPerformance))
}

var otsInstanceAttachmentBasicMap = map[string]string{
	"instance_name": CHECKSET,
	"vpc_name":      CHECKSET,
	"vswitch_id":    CHECKSET,
	"vpc_id":        CHECKSET,
}
