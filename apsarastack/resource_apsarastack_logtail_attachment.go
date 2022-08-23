package apsarastack

import (
	"fmt"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackLogtailAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackLogtailAttachmentCreate,
		Read:   resourceApsaraStackLogtailAttachmentRead,
		Delete: resourceApsaraStackLogtailAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logtail_config_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"machine_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceApsaraStackLogtailAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	project := d.Get("project").(string)
	config_name := d.Get("logtail_config_name").(string)
	group_name := d.Get("machine_group_name").(string)
	var requestInfo *sls.Client
	raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		requestInfo = slsClient
		return nil, slsClient.ApplyConfigToMachineGroup(project, config_name, group_name)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_logtail_attachment", "ApplyConfigToMachineGroup", ApsaraStackLogGoSdkERROR)
	}
	if debugOn() {
		addDebug("ApplyConfigToMachineGroup", raw, requestInfo, map[string]string{
			"project":   project,
			"confName":  config_name,
			"groupName": group_name,
		})
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", project, COLON_SEPARATED, config_name, COLON_SEPARATED, group_name))
	return resourceApsaraStackLogtailAttachmentRead(d, meta)
}

func resourceApsaraStackLogtailAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	object, err := logService.DescribeLogtailAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("project", parts[0])
	d.Set("logtail_config_name", parts[1])
	d.Set("machine_group_name", object)

	return nil
}

func resourceApsaraStackLogtailAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	var requestInfo *sls.Client
	raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		requestInfo = slsClient
		return nil, slsClient.RemoveConfigFromMachineGroup(parts[0], parts[1], parts[2])
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "RemoveConfigFromMachineGroup", ApsaraStackLogGoSdkERROR)
	}
	if debugOn() {
		addDebug("RemoveConfigFromMachineGroup", raw, requestInfo, map[string]string{
			"project":   parts[0],
			"confName":  parts[1],
			"groupName": parts[2],
		})
	}
	return WrapError(logService.WaitForLogtailAttachment(d.Id(), Deleted, DefaultTimeout))

}
