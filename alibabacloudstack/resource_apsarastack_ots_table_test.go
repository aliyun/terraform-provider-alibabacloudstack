package alibabacloudstack

import (
	"fmt"

	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackOtsTable_basic(t *testing.T) {
	var v *tablestore.DescribeTableResponse

	resourceId := "alibabacloudstack_ots_table.default"
	ra := resourceAttrInit(resourceId, otsTableBasicMap)
	serviceFunc := func() interface{} {
		return &OtsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOtsTableConfigDependence)

	ResourceTest(t, resource.TestCase{
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
					"table_name":    "${var.name}",
					"primary_key": []map[string]interface{}{
						{
							"name": "pk1",
							"type": "Integer",
						},
					},
					"time_to_live": "-1",
					"max_version":  "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-" + name,
						"table_name":    name,
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
					"deviation_cell_version_in_sec": "86401",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deviation_cell_version_in_sec": "86401",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_to_live": "86401",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_to_live": "86401",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_version": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_version": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_to_live":                  "-1",
					"max_version":                   "1",
					"deviation_cell_version_in_sec": "86400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_to_live":                  "-1",
						"max_version":                   "1",
						"deviation_cell_version_in_sec": "86400",
					}),
				),
			},
		},
	})
}

//func TestAccAlibabacloudStackOtsTable_highPerformance(t *testing.T) {
//	var v *tablestore.DescribeTableResponse
//
//	resourceId := "alibabacloudstack_ots_table.default"
//	ra := resourceAttrInit(resourceId, otsTableBasicMap)
//	serviceFunc := func() interface{} {
//		return &OtsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
//	}
//	rc := resourceCheckInit(resourceId, &v, serviceFunc)
//
//	rac := resourceAttrCheckInit(rc, ra)
//
//	testAccCheck := rac.resourceAttrMapUpdateSet()
//	rand := getAccTestRandInt(10000, 99999)
//	name := fmt.Sprintf("testAcc%d", rand)
//	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOtsTableConfigDependenceHighperformance)
//
//	ResourceTest(t, resource.TestCase{
//		PreCheck: func() {
//			testAccPreCheck(t)

//		},
//		// module name
//		IDRefreshName: resourceId,
//		Providers:     testAccProviders,
//		CheckDestroy:  rac.checkResourceDestroy(),
//		Steps: []resource.TestStep{
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"instance_name": "${alibabacloudstack_ots_instance.default.name}",
//					"table_name":    "${var.name}",
//					"primary_key": []map[string]interface{}{
//						{
//							"name": "pk1",
//							"type": "Integer",
//						},
//					},
//					"time_to_live": "-1",
//					"max_version":  "1",
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"instance_name": "tf-" + name,
//						"table_name":    name,
//					}),
//				),
//			},
//			{
//				ResourceName:      resourceId,
//				ImportState:       true,
//				ImportStateVerify: true,
//			},
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"deviation_cell_version_in_sec": "86401",
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"deviation_cell_version_in_sec": "86401",
//					}),
//				),
//			},
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"time_to_live": "86401",
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"time_to_live": "86401",
//					}),
//				),
//			},
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"max_version": "2",
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"max_version": "2",
//					}),
//				),
//			},
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"time_to_live":                  "-1",
//					"max_version":                   "1",
//					"deviation_cell_version_in_sec": "86400",
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"time_to_live":                  "-1",
//						"max_version":                   "1",
//						"deviation_cell_version_in_sec": "86400",
//					}),
//				),
//			},
//		},
//	})
//}

//func TestAccAlibabacloudStackOtsTable_multi(t *testing.T) {
//	var v *tablestore.DescribeTableResponse
//
//	resourceId := "alibabacloudstack_ots_table.default.4"
//	ra := resourceAttrInit(resourceId, otsTableBasicMap)
//	serviceFunc := func() interface{} {
//		return &OtsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
//	}
//	rc := resourceCheckInit(resourceId, &v, serviceFunc)
//
//	rac := resourceAttrCheckInit(rc, ra)
//
//	testAccCheck := rac.resourceAttrMapUpdateSet()
//	rand := getAccTestRandInt(10000, 99999)
//	name := fmt.Sprintf("testAcc%d", rand)
//	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOtsTableConfigDependence)
//
//	ResourceTest(t, resource.TestCase{
//		PreCheck: func() {
//			testAccPreCheck(t)

//			},
//			// module name
//			IDRefreshName: resourceId,
//			Providers:     testAccProviders,
//			CheckDestroy:  rac.checkResourceDestroy(),
//			Steps: []resource.TestStep{
//				{
//					Config: testAccConfig(map[string]interface{}{
//						"instance_name": "${alibabacloudstack_ots_instance.default.name}",
//						"table_name":    "${var.name}${count.index}",
//						"primary_key": []map[string]interface{}{
//							{
//								"name": "pk1",
//								"type": "Integer",
//							},
//						},
//						"time_to_live": "-1",
//						"max_version":  "1",
//						"count":        "5",
//					}),
//					Check: resource.ComposeTestCheckFunc(
//						testAccCheck(nil),
//					),
//				},
//			},
//		})
//	}
func resourceOtsTableConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}
	resource "alibabacloudstack_ots_instance" "default" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "%s"
	  tags = {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}
	`, name, string(OtsCapacity))
}

func resourceOtsTableConfigDependenceHighperformance(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}
	resource "alibabacloudstack_ots_instance" "default" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "%s"
	  tags = {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}
	`, name, string(OtsCapacity))
}

var otsTableBasicMap = map[string]string{
	"primary_key.#":                 "1",
	"primary_key.0.name":            "pk1",
	"primary_key.0.type":            "Integer",
	"time_to_live":                  "-1",
	"max_version":                   "1",
	"deviation_cell_version_in_sec": "86400",
}

func testAccCheckOtsTableExist(n string, table *tablestore.DescribeTableResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found OTS table: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no OTS table ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		otsService := OtsService{client}

		response, err := otsService.DescribeOtsTable(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error finding OTS table %s: %#v", rs.Primary.ID, err)
		}

		table = response
		return nil
	}
}

func testAccCheckOtsTableDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_ots_table" {
			continue
		}

		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		otsService := OtsService{client}

		if _, err := otsService.DescribeOtsTable(rs.Primary.ID); err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("error! Ots table still exists")
	}

	return nil
}
