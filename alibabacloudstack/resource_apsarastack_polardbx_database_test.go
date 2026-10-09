package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)


// sweepPolardbxDatabases deletes all databases in the given instance.
func sweepPolardbxDatabases(client *connectivity.AlibabacloudStackClient, instanceId string) error {
	request := client.NewCommonRequest("GET", "polardbx", "2020-02-02", "DescribeDbList", "")
	request.QueryParams["DBInstanceName"] = instanceId
	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		return fmt.Errorf("DescribeDbList failed: %s", err)
	}
	var resp PolardbxDescribeDbListResponse
	if err := json.Unmarshal(bresponse.GetHttpContentBytes(), &resp); err != nil {
		return fmt.Errorf("parse DescribeDbList response failed: %s", err)
	}
	for _, db := range resp.Data {
		log.Printf("[INFO] Deleting database %s from instance %s", db.DBName, instanceId)
		delReq := client.NewCommonRequest("POST", "polardbx", "2020-02-02", "DeleteDB", "")
		delReq.QueryParams["DBInstanceName"] = instanceId
		delReq.QueryParams["DbName"] = db.DBName
		if _, err := client.ProcessCommonRequest(delReq); err != nil {
			log.Printf("[WARN] Failed to delete database %s: %s", db.DBName, err)
		}
	}
	return nil
}

func TestAccAlibabacloudStackPolardbxDatabase_basic0(t *testing.T) {
	var v *PolardbxDatabase

	resourceId := "alibabacloudstack_polardbx_database.default"
	ra := resourceAttrInit(resourceId, map[string]string{})

	serviceFunc := func() interface{} {
		return &PolardbXService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DoPolardbxDescribeDbListRequest")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_polardbx_db_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbxDatabaseDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)

		},
		// module name
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       "${local.polardbx_instance.id}",
					"database_name":     "${var.name}",
					"encode":            "utf8mb4",
					"description":       "${var.name}",
					"mode":              "auto",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"database_name":     name,
						"encode":            "utf8mb4",
						"description":       name,
						"mode":              "auto",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("%s_update", name),
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// mode cannot be read back
				ImportStateVerifyIgnore: []string{"mode"},
			},
		},
	})
}

func resourcePolardbxDatabaseDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

%s

 `, name, VSwitchCommonTestCase, PolardbxReadOrCreateCommonTestCase())
}
