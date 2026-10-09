package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)


// sweepPolardbxAccounts deletes all non-super accounts in the given instance.
// If rights separation is enabled, it passes the security account name.
func sweepPolardbxAccounts(client *connectivity.AlibabacloudStackClient, instanceId string) error {
	request := client.NewCommonRequest("GET", "polardbx", "2020-02-02", "DescribeAccountList", "")
	request.QueryParams["DBInstanceName"] = instanceId
	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		return fmt.Errorf("DescribeAccountList failed: %s", err)
	}
	var resp DoPolardbxDescribeAccountListResponse
	if err := json.Unmarshal(bresponse.GetHttpContentBytes(), &resp); err != nil {
		return fmt.Errorf("parse DescribeAccountList response failed: %s", err)
	}

	// Check if rights separation is enabled (security account exists)
	var securityAccountName string
	for _, acc := range resp.Data {
		if acc.AccountType == "3" {
			securityAccountName = acc.AccountName
			break
		}
	}

	// Delete non-super accounts (AccountType "1" = super, "2" = super-like; skip them)
	for _, acc := range resp.Data {
		if acc.AccountType == "1" || acc.AccountType == "2" {
			continue // skip super accounts
		}
		if acc.AccountType == "3" || acc.AccountType == "4" {
			continue // skip security and audit accounts
		}
		log.Printf("[INFO] Deleting account %s (type=%s) from instance %s", acc.AccountName, acc.AccountType, instanceId)
		delReq := client.NewCommonRequest("POST", "polardbx", "2020-02-02", "DeleteAccount", "")
		delReq.QueryParams["DBInstanceName"] = instanceId
		delReq.QueryParams["AccountName"] = acc.AccountName
		if securityAccountName != "" {
			delReq.QueryParams["SecurityAccountName"] = securityAccountName
		}
		if _, err := client.ProcessCommonRequest(delReq); err != nil {
			log.Printf("[WARN] Failed to delete account %s: %s", acc.AccountName, err)
		}
	}
	return nil
}

func TestAccAlibabacloudStackPolardbxAccount_basic0(t *testing.T) {
	var v *PolardbxAccount

	resourceId := "alibabacloudstack_polardbx_account.default"
	ra := resourceAttrInit(resourceId, PolardbxAccountbasicMap)

	serviceFunc := func() interface{} {
		return &PolardbXService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DoPolardbxDescribeAccountRequest")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_pldbx_account_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbxAccountDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)

		},
		// module name
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${local.polardbx_instance.id}",
					"account_name": "${var.name}",
					"password":     "${random_password.password.0.result}",
					"description":  "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name": name,
						"description":  name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// password cannot be read back
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "${random_password.password.1.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
		},
	})
}

func resourcePolardbxAccountDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

%s

%s

 `, name, RandomPasswordTestCase(12,2), VSwitchCommonTestCase, PolardbxReadOrCreateCommonTestCase())
}

var PolardbxAccountbasicMap = map[string]string{
	"instance_id":  CHECKSET,
	"account_name": CHECKSET,
}
