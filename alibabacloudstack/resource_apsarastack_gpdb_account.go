package alibabacloudstack

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackGpdbAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackGpdbAccountCreate,
		Read:   resourceAlibabacloudStackGpdbAccountRead,
		Update: resourceAlibabacloudStackGpdbAccountUpdate,
		Delete: resourceAlibabacloudStackGpdbAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[A-Za-z][\w\\_]{2,255}$`), "The description of the account. The description must be 2 to 256 characters in length and can contain letters, digits, underscores (_)."),
			},
			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9_]{1,14}[a-z0-9]$`), "The name of the account. The name must be 2 to 16 characters in length and can contain lower letters, digits, underscores (_)."),
			},
			"account_password": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(8, 32),
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackGpdbAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "CreateAccount"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("account_description"); ok {
		request["AccountDescription"] = v
	}
	request["AccountName"] = d.Get("account_name")
	request["DBInstanceId"] = d.Get("db_instance_id")
	request["AccountPassword"] = d.Get("account_password")

	_, err := client.DoTeaRequest("POST", "gpdb", "2016-05-03", action, "", nil, request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_gpdb_account", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DBInstanceId"], ":", request["AccountName"]))
	gpdbService := GpdbService{client}
	stateConf := BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, gpdbService.GpdbAccountStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return resourceAlibabacloudStackGpdbAccountRead(d, meta)
}

func resourceAlibabacloudStackGpdbAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	gpdbService := GpdbService{client}
	object, err := gpdbService.DescribeGpdbAccount(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_gpdb_account gpdbService.DescribeGpdbAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("account_name", parts[1])
	d.Set("db_instance_id", parts[0])
	d.Set("account_description", object["AccountDescription"])
	d.Set("status", convertGpdbAccountStatusResponse(object["AccountStatus"]))
	return nil
}

func resourceAlibabacloudStackGpdbAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	request := map[string]interface{}{
		"AccountName": parts[1],
		"DBInstanceId": parts[0],
	}

	update := false
	if d.HasChange("account_password") {
		update = true
		if v, ok := d.GetOk("account_password"); ok {
			request["AccountPassword"] = v
		}
	}

	if update {
		action := "ResetAccountPassword"
		_, err = client.DoTeaRequest("POST", "gpdb", "2016-05-03", action, "", nil, request)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	return resourceAlibabacloudStackGpdbAccountRead(d, meta)
}

func resourceAlibabacloudStackGpdbAccountDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourcealibabacloudstackGpdbAccount. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

func convertGpdbAccountStatusResponse(source interface{}) interface{} {
	switch source {
	case "Creating":
		return "0"
	case "Active":
		return "1"
	case "Deleting":
		return "3"
	}
	return source
}
