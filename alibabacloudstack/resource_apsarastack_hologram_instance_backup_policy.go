package alibabacloudstack

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackHologramInstanceBackupPolicy() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hour": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 23),
			},
			"data_keep_quantity": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 31),
			},
			"week": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"0", "1", "2", "3", "4", "5", "6"}, false),
				},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackHologramInstanceBackupPolicyCreate,
		resourceAlibabacloudStackHologramInstanceBackupPolicyRead, resourceAlibabacloudStackHologramInstanceBackupPolicyUpdate, resourceAlibabacloudStackHologramInstanceBackupPolicyDelete)
	return resource
}

func resourceAlibabacloudStackHologramInstanceBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("instance_id").(string))
	return nil
}

func resourceAlibabacloudStackHologramInstanceBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	weekList := d.Get("week").(*schema.Set).List()
	weekListstring := make([]string, len(weekList))
	for i, v := range weekList {
		weekListstring[i] = v.(string)
	}
	week := strings.Join(weekListstring, ",")
	body := map[string]interface{}{
		"RegionId":         client.RegionId,
		"instanceId":       d.Id(),
		"hour":             fmt.Sprintf("%d", d.Get("hour")),
		"dataKeepQuantity": d.Get("data_keep_quantity"),
		"enabled":          d.Get("enabled"),
		"week":             week,
	}
	request := map[string]interface{}{
		"x-acs-body": body,
	}
	response, err := client.DoTeaRequest("POST", "Hologram", "2022-06-01", "UpdateScheduledBackupConfig", "/api/v1/backups/scheduledConfig", nil, nil, request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance_backup_policy", "UpdateScheduledBackupConfig", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if !response["Data"].(bool) {
		return errmsgs.WrapError(errmsgs.Error("UpdateScheduledBackupConfig Failed! %#v", response))
	}
	return nil
}

func resourceAlibabacloudStackHologramInstanceBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	hologramService := HologramService{client}
	policy, err := hologramService.DescribeHologramInstanceBackupPolicy(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance_backup_policy", "DescribeHologramInstanceBackupPolicy", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	d.Set("instance_id", policy["instance_id"])
	d.Set("hour", policy["hour"])
	d.Set("week", strings.Split(policy["week"].(string), ","))
	d.Set("data_keep_quantity", policy["data_keep_quantity"])
	d.Set("type", policy["type"])
	d.Set("enabled", policy["Enabled"])
	return nil
}

func resourceAlibabacloudStackHologramInstanceBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	body := map[string]interface{}{
		"RegionId":   client.RegionId,
		"instanceId": d.Id(),
		"enabled":    false,
	}
	request := map[string]interface{}{
		"x-acs-body": body,
	}
	response, err := client.DoTeaRequest("POST", "Hologram", "2022-06-01", "UpdateScheduledBackupConfig", "/api/v1/backups/scheduledConfig", nil, nil, request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance_backup_policy", "UpdateScheduledBackupConfig", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if !response["Data"].(bool) {
		return errmsgs.WrapError(errmsgs.Error("UpdateScheduledBackupConfig Failed! %#v", response))
	}
	return nil
}
