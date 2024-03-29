package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDBAccountPrivilege() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDBAccountPrivilegeCreate,
		Read:   resourceAlibabacloudStackDBAccountPrivilegeRead,
		Update: resourceAlibabacloudStackDBAccountPrivilegeUpdate,
		Delete: resourceAlibabacloudStackDBAccountPrivilegeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"account_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"privilege": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ReadOnly", "ReadWrite"}, false),
				ForceNew:     true,
			},

			"db_names": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
		},
	}
}

func resourceAlibabacloudStackDBAccountPrivilegeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}
	instanceId := d.Get("instance_id").(string)
	account := d.Get("account_name").(string)
	privilege := d.Get("privilege").(string)
	dbList := d.Get("db_names").(*schema.Set).List()
	// wait instance running before granting
	if err := rdsService.WaitForDBInstance(instanceId, Running, DefaultLongTimeout); err != nil {
		return WrapError(err)
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", instanceId, COLON_SEPARATED, account, COLON_SEPARATED, privilege))

	if len(dbList) > 0 {
		for _, db := range dbList {
			if err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				if err := rdsService.GrantAccountPrivilege(d.Id(), db.(string)); err != nil {
					if IsExpectedErrors(err, OperationDeniedDBStatus) {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			}); err != nil {
				return WrapError(err)
			}
		}
	}

	return resourceAlibabacloudStackDBAccountPrivilegeRead(d, meta)
}

func resourceAlibabacloudStackDBAccountPrivilegeRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	rsdService := RdsService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	object, err := rsdService.DescribeDBAccountPrivilege(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.DBInstanceId)
	d.Set("account_name", object.AccountName)
	d.Set("privilege", parts[2])
	var names []string
	for _, pri := range object.DatabasePrivileges.DatabasePrivilege {
		if pri.AccountPrivilege == parts[2] {
			names = append(names, pri.DBName)
		}
	}

	if len(names) < 1 && strings.HasPrefix(object.DBInstanceId, "pgm-") {

		request := rds.CreateDescribeDatabasesRequest()
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.DBInstanceId = object.DBInstanceId
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.DescribeDatabases(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"InternalError", "OperationDenied.DBInstanceStatus"}) {
					return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, object.DBInstanceId, request.GetActionName(), AlibabacloudStackSdkGoERROR))
				}
				return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, object.DBInstanceId, request.GetActionName(), AlibabacloudStackSdkGoERROR))
			}

			addDebug(request.GetActionName(), raw, request.RpcRequest, request)

			response, _ := raw.(*rds.DescribeDatabasesResponse)
			for _, db := range response.Databases.Database {
				for _, account := range db.Accounts.AccountPrivilegeInfo {
					if account.Account == object.AccountName && (account.AccountPrivilege == parts[2] || account.AccountPrivilege == "ALL") {
						names = append(names, db.DBName)
					}
				}
			}
			return nil
		})
		if err != nil {
			return WrapError(err)
		}
	}

	d.Set("db_names", names)

	return nil
}

func resourceAlibabacloudStackDBAccountPrivilegeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}
	d.Partial(true)

	if d.HasChange("db_names") {
		parts := strings.Split(d.Id(), COLON_SEPARATED)

		o, n := d.GetChange("db_names")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		if len(remove) > 0 {
			if strings.HasPrefix(d.Id(), "pgm-") {
				return WrapError(fmt.Errorf("At present, the PostgreSql database does not support revoking the current privilege."))
			}
			// wait instance running before revoking
			if err := rdsService.WaitForDBInstance(parts[0], Running, DefaultTimeoutMedium); err != nil {
				return WrapError(err)
			}
			for _, db := range remove {
				if err := rdsService.RevokeAccountPrivilege(d.Id(), db.(string)); err != nil {
					return WrapError(err)
				}
			}
		}

		if len(add) > 0 {
			// wait instance running before granting
			if err := rdsService.WaitForDBInstance(parts[0], Running, DefaultTimeoutMedium); err != nil {
				return WrapError(err)
			}
			for _, db := range add {
				if err := rdsService.GrantAccountPrivilege(d.Id(), db.(string)); err != nil {
					return WrapError(err)
				}
			}
		}
		//d.SetPartial("db_names")
	}

	d.Partial(false)
	return resourceAlibabacloudStackDBAccountPrivilegeRead(d, meta)
}

func resourceAlibabacloudStackDBAccountPrivilegeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	object, err := rdsService.DescribeDBAccountPrivilege(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	if strings.HasPrefix(d.Id(), "pgm-") {
		return nil
	}
	var dbName string

	if len(object.DatabasePrivileges.DatabasePrivilege) > 0 {
		for _, pri := range object.DatabasePrivileges.DatabasePrivilege {
			if pri.AccountPrivilege == parts[2] {
				dbName = pri.DBName
				if err := rdsService.RevokeAccountPrivilege(d.Id(), pri.DBName); err != nil {
					return WrapError(err)
				}
			}
		}
	}

	return rdsService.WaitForAccountPrivilege(d.Id(), dbName, Deleted, DefaultTimeoutMedium)
}
