package alibabacloudstack

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackKVStoreBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackKVStoreBackupPolicyCreate,
		Read:   resourceAlibabacloudStackKVStoreBackupPolicyRead,
		Update: resourceAlibabacloudStackKVStoreBackupPolicyUpdate,
		Delete: resourceAlibabacloudStackKVStoreBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backup_time": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Default:      "02:00Z-03:00Z",
				Deprecated:   "Field 'backup_time' is deprecated and will be removed in a future release. Please use new field 'preferred_backup_time' instead.",
				ConflictsWith: []string{"preferred_backup_time"},
			},
			"preferred_backup_time": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Default:      "02:00Z-03:00Z",
				ConflictsWith: []string{"backup_time"},
			},
			"backup_period": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				Deprecated: "Field 'backup_period' is deprecated and will be removed in a future release. Please use new field 'preferred_backup_period' instead.",
				ConflictsWith: []string{"preferred_backup_period"},
			},
			"preferred_backup_period": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				ConflictsWith: []string{"backup_period"},
			},
		},
	}
}

func resourceAlibabacloudStackKVStoreBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("instance_id").(string))
	return resourceAlibabacloudStackKVStoreBackupPolicyUpdate(d, meta)
}

func resourceAlibabacloudStackKVStoreBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	kvstoreService := KvstoreService{client}

	object, err := kvstoreService.DescribeKVstoreBackupPolicy(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("instance_id", d.Id())
	connectivity.SetResourceData(d, object.PreferredBackupTime, "preferred_backup_time", "backup_time");
	connectivity.SetResourceData(d, strings.Split(object.PreferredBackupPeriod, ","), "preferred_backup_period", "backup_period")

	return nil
}

func resourceAlibabacloudStackKVStoreBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("preferred_backup_time") || d.HasChange("preferred_backup_period") || d.HasChange("backup_time") || d.HasChange("backup_period") {
		client := meta.(*connectivity.AlibabacloudStackClient)
		kvstoreService := KvstoreService{client}

		request := r_kvstore.CreateModifyBackupPolicyRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.InstanceId = d.Id()
		request.PreferredBackupTime = connectivity.GetResourceData(d, "preferred_backup_time", "backup_time")
		periodList =  expandStringList(connectivity.GetResourceData(d, "preferred_backup_period", "backup_period").(*schema.Set).List())
		request.PreferredBackupPeriod = strings.Join(periodList, COMMA_SEPARATED)

		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.ModifyBackupPolicy(request)
		})
		response, ok := raw.(*r_kvstore.ModifyBackupPolicyResponse)
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		// There is a random error and need waiting some seconds to ensure the update is success
		_, err = kvstoreService.DescribeKVstoreBackupPolicy(d.Id())
		if err != nil {
			return errmsgs.WrapError(err)
		}
	}

	return resourceAlibabacloudStackKVStoreBackupPolicyRead(d, meta)
}

func resourceAlibabacloudStackKVStoreBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := r_kvstore.CreateModifyBackupPolicyRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = d.Id()

	request.PreferredBackupTime = "01:00Z-02:00Z"
	request.PreferredBackupPeriod = "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday"

	raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.ModifyBackupPolicy(request)
	})
	response, ok := raw.(*r_kvstore.ModifyBackupPolicyResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	return nil
}
