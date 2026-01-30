package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackOtsInstanceAttachmentsDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_ots_instance_attachments.default"
	name := fmt.Sprintf("tfvpcatt%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceOtsInstanceAttachmentsConfigDependence)

	instanceNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alibabacloudstack_ots_instance_attachment.default.instance_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "fakeinstname",
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alibabacloudstack_ots_instance_attachment.default.instance_name}",
			"name_regex":    "${alibabacloudstack_ots_instance_attachment.default.vpc_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alibabacloudstack_ots_instance_attachment.default.instance_name}",
			"name_regex":    "${alibabacloudstack_ots_instance_attachment.default.vpc_name}-fake",
		}),
	}

	var existOtsInstanceAttachmentsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                     "1",
			"names.0":                     name,
			"vpc_ids.#":                   "1",
			"vpc_ids.0":                   CHECKSET,
			"attachments.#":               "1",
			"attachments.0.id":            CHECKSET,
			"attachments.0.domain":        CHECKSET,
			"attachments.0.endpoint":      CHECKSET,
			"attachments.0.region":        CHECKSET,
			"attachments.0.instance_name": name,
			"attachments.0.vpc_name":      name,
			"attachments.0.vpc_id":        CHECKSET,
		}
	}

	var fakeOtsInstanceAttachmentsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":       "0",
			"vpc_ids.#":     "0",
			"attachments.#": "0",
		}
	}

	var otsInstanceAttachmentsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOtsInstanceAttachmentsMapFunc,
		fakeMapFunc:  fakeOtsInstanceAttachmentsMapFunc,
	}
	otsInstanceAttachmentsCheckInfo.dataSourceTestCheck(t, rand, instanceNameConf, nameRegexConf)
}

func dataSourceOtsInstanceAttachmentsConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}

	%s

	data "alibabacloudstack_ots_clusters" "default" {}

	resource "alibabacloudstack_ots_instance" "default" {
	  name = var.name
	  description   = var.name
	  specification  = "HYBRID"
	}
	
	resource "alibabacloudstack_ots_instance_attachment" "default" {
		instance_name = "${alibabacloudstack_ots_instance.default.name}"
		vpc_name      = "${var.name}"
		vswitch_id    = "${alibabacloudstack_vpc_vswitch.default.id}"
	}
	
	`, name, VSwitchCommonTestCase)
}
