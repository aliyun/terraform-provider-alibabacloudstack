package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDataWorksFile_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_data_works_file.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackDataWorksFileMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataworksService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDataWorksFile")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf_file%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackDataWorksFileBasicDependence0)
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
					"project_id": "${alibabacloudstack_data_works_project.default.id}",
					"folder_id":  `${alibabacloudstack_data_works_folder.default0.folder_id}`,
					"file_name":  "${var.name}",
					"file_type":  "${data.alibabacloudstack_dataworks_file_types.anyone.file_types.0.node_type_id}",
					"content": TfRawString(`<<EOF
					#!/bin/bash
					#********************************************************************#
					##author:ascm-org-1770602583560
					##create time:2026-02-13 16:08:39
					#********************************************************************#
					echo "Hello World"
					EOF`),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_name": name,
						"file_type": CHECKSET,
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
					"folder_id":   `${alibabacloudstack_data_works_folder.default1.folder_id}`,
					"file_name":   "${var.name}_update",
					"description": "${var.name} desc",
					"content": TfRawString(`<<EOF
					#!/bin/bash
					#********************************************************************#
					##author:ascm-org-1770602583560
					##create time:2026-02-13 16:08:39
					#********************************************************************#
					echo "Hello World!"
					EOF`),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_name":   name + "_update",
						"file_type":   CHECKSET,
						"description": name + " desc",
					}),
				),
			},
		},
	})
}

// AlibabacloudStackDataWorksFileMap0 is the expected state map for the DataWorks folder resource
var AlibabacloudStackDataWorksFileMap0 = map[string]string{
	"folder_id":  CHECKSET,
	"project_id": CHECKSET,
}

// AlibabacloudStackDataWorksFileBasicDependence0 returns the basic dependencies for the DataWorks folder resource test
func AlibabacloudStackDataWorksFileBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}


resource "alibabacloudstack_data_works_project" "default" {
	name =           "${var.name}"
	description =    "${var.name}_desc"
	task_auth_type = "PROJECT"
	}

resource "alibabacloudstack_data_works_business" "default" {
	project_id=  "${alibabacloudstack_data_works_project.default.id}"
	name=        "${var.name}"
	description= "${var.name}_desc"
}

resource "alibabacloudstack_data_works_folder" "default0" {
	project_id    = "${alibabacloudstack_data_works_project.default.id}"
	business_name = "${alibabacloudstack_data_works_business.default.name}"
	engine_type   = "General"
	folder_path   = "folder_test"
}

resource "alibabacloudstack_data_works_folder" "default1" {
	project_id    = "${alibabacloudstack_data_works_project.default.id}"
	business_name = "${alibabacloudstack_data_works_business.default.name}"
	engine_type   = "General"
	folder_path   = "folder_test_new"
}

data "alibabacloudstack_dataworks_file_types" "anyone" {
	name = "Shell"
	project_id = "${alibabacloudstack_data_works_project.default.id}"
}

`, name)
}
