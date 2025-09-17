package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackPolardbClusterAccount() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Normal",
				ValidateFunc: validation.StringInSlice([]string{"Normal", "Lock"}, false),
			},
			"account_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"account_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_lock_state": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"UnLock", "Lock"}, false),
			},
			"account_password_valid_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackPolardbClusterAccountCreate, resourceAlibabacloudStackPolardbClusterAccountRead, resourceAlibabacloudStackPolardbClusterAccountUpdate, resourceAlibabacloudStackPolardbClusterAccountDelete)
	return resource
}

func resourceAlibabacloudStackPolardbClusterAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolardbService{client}

	// Prepare the request parameters for CreateAccount API
	reqBody := map[string]interface{}{
		"DBClusterId":        d.Get("db_cluster_id").(string),
		"AccountName":        d.Get("account_name").(string),
		"AccountType":        d.Get("account_type").(string),
		"AccountDescription": d.Get("account_description").(string),
		"AccountPassword":    d.Get("account_password").(string),
	}

	// Call the CreateAccount API
	if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "CreateAccount", "", nil, nil, reqBody); err != nil {
		return err
	}

	// Generate and set the resource ID: {DBClusterId}:{AccountName}
	dbClusterId := d.Get("db_cluster_id").(string)
	accountName := d.Get("account_name").(string)
	d.SetId(fmt.Sprintf("%s:%s", dbClusterId, accountName))

	// Wait for the account status to become Available
	stateConf := BuildStateConf([]string{"UnAvailable"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, polardbService.PolardbClusterAccountStateRefreshFunc(d.Id(), []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func resourceAlibabacloudStackPolardbClusterAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polarDBService := PolardbService{client}

	// The resource ID is in the format {DBClusterId}:{AccountName}
	object, err := polarDBService.DescribePolardbClusterAccount(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_polardb_cluster_account polarDBService.DescribePolardbClusterAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	parts := strings.SplitN(d.Id(), ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid id format, expected DBClusterId:AccountName")
		}
	d.Set("db_cluster_id",  parts[0])
	d.Set("account_name",  parts[1])
		
	// Set the resource data from the returned object
	d.Set("account_type", object["AccountType"])
	d.Set("account_description", object["AccountDescription"])
	d.Set("account_status", object["AccountStatus"])
	d.Set("account_lock_state", object["AccountLockState"])
	d.Set("account_password_valid_time", object["AccountPasswordValidTime"])

	return nil
}

func resourceAlibabacloudStackPolardbClusterAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Parse resource ID to get DBClusterId and AccountName
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return err
	}
	dbClusterId := parts[0]
	accountName := parts[1]

	// Update account lock state if changed
	if d.HasChange("account_lock_state") {
		reqQuery := map[string]interface{}{
			"DBClusterId":      dbClusterId,
			"AccountName":      accountName,
			"AccountLockState": d.Get("account_lock_state").(string),
		}

		if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "ModifyAccountLockState", "", nil, reqQuery, nil); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_polardb_cluster_account", "ModifyAccountLockState", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	// If it's a new resource, do nothing in update
	if d.IsNewResource() {
		return nil
	}
	// Update account description if changed
	if d.HasChange("account_description") {
		reqQuery := map[string]interface{}{
			"DBClusterId":        dbClusterId,
			"AccountName":        accountName,
			"AccountDescription": d.Get("account_description").(string),
		}

		if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "ModifyAccountDescription", "", nil, reqQuery, nil); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_polardb_cluster_account", "ModifyAccountDescription", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	// Update account password if changed
	if d.HasChange("account_password") {
		reqQuery := map[string]interface{}{
			"DBClusterId":        dbClusterId,
			"AccountName":        accountName,
			"NewAccountPassword": d.Get("account_password").(string),
		}

		if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "ModifyAccountPassword", "", nil, reqQuery, nil); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_polardb_cluster_account", "ModifyAccountPassword", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	return resourceAlibabacloudStackPolardbClusterAccountRead(d, meta)
}

func resourceAlibabacloudStackPolardbClusterAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	dbClusterId := d.Get("db_cluster_id").(string)
	accountName := d.Get("account_name").(string)

	reqQuery := map[string]interface{}{
		"AccountName": accountName,
		"DBClusterId": dbClusterId,
	}

	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		_, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "DeleteAccount", "", nil, reqQuery, nil)
		if err != nil {
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, accountName, "DeleteAccount", errmsgs.AlibabacloudStackSdkGoERROR, "")
			return resource.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
