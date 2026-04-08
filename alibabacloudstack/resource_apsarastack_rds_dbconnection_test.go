package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccAlibabacloudStackDBConnectionPublic tests public network connection
func TestAccAlibabacloudStackDBConnectionPublic(t *testing.T) {
	var v *rds.DBInstanceNetInfo
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tftestAccDBconnectionPublic%d", rand)
	var basicMap = map[string]string{
		"instance_id":       CHECKSET,
		"connection_string": CHECKSET,
		"port":              CHECKSET,
		"ip_address":        CHECKSET,
	}
	resourceId := "alibabacloudstack_db_connection.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBConnectionConfigDependence)
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
					"instance_id":       "${alibabacloudstack_db_instance.default.id}",
					"network_type":      "public",
					"connection_prefix": fmt.Sprintf("tftest%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": fmt.Sprintf("tftest%d", rand),
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
					"port":              "3333",
					"connection_prefix": fmt.Sprintf("tftest%d2", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port":              "3333",
						"connection_prefix": fmt.Sprintf("tftest%d2", rand),
					}),
				),
			},
		},
	})
}

// TestAccAlibabacloudStackDBConnectionPrivate tests private network connection
func TestAccAlibabacloudStackDBConnectionPrivate(t *testing.T) {
	var v *rds.DBInstanceNetInfo
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-testAccDBconnectionPrivate%d", rand)

	var basicMap = map[string]string{
		"instance_id":       CHECKSET,
		"connection_string": CHECKSET,
		"port":              CHECKSET,
		"ip_address":        CHECKSET,
	}
	resourceId := "alibabacloudstack_db_connection.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBConnectionConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		// private connection strings can`t delete
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alibabacloudstack_db_instance.default.id}",
					"network_type": "private",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port":              "3333",
					"connection_prefix": fmt.Sprintf("tftest%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port":              "3333",
						"connection_prefix": fmt.Sprintf("tftest%d", rand),
					}),
				),
			},
		},
	})
}

func resourceDBConnectionConfigDependence(name string) string {
	return fmt.Sprintf(`
	

	variable "creation" {
		default = "Rds"
	}

	variable "name" {
		default = "%s"
	}
	%s
	%s 
	`, name, VSwitchCommonTestCase, RdsMysqlCommonTestCase())
}
