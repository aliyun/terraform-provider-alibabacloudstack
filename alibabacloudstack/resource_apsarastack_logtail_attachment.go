package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackLogtailAttachment() *schema.Resource {
	resource := &schema.Resource{
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
	setResourceFunc(resource, resourceAlibabacloudStackLogtailAttachmentCreate, resourceAlibabacloudStackLogtailAttachmentRead, nil, resourceAlibabacloudStackLogtailAttachmentDelete)
	return resource
}

func resourceAlibabacloudStackLogtailAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	logService := LogService{client}
	project := d.Get("project").(string)
	config_name := d.Get("logtail_config_name").(string)
	group_name := d.Get("machine_group_name").(string)

	slsClient, err := logService.GetSlsDataClient(project)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	err = slsClient.ApplyConfigToMachineGroup(project, config_name, group_name)
	if err != nil {
		errmsg := ""
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_logtail_attachment", "ApplyConfigToMachineGroup", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
	}

	d.SetId(fmt.Sprintf("%s%s%s%s%s", project, COLON_SEPARATED, config_name, COLON_SEPARATED, group_name))
	return nil
}

func resourceAlibabacloudStackLogtailAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	object, err := logService.DescribeLogtailAttachment(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("project", parts[0])
	d.Set("logtail_config_name", parts[1])
	d.Set("machine_group_name", object)

	return nil
}

func resourceAlibabacloudStackLogtailAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	slsClient, err := logService.GetSlsDataClient(parts[0])
	if err != nil {
		return errmsgs.WrapError(err)
	}

	err = slsClient.RemoveConfigFromMachineGroup(parts[0], parts[1], parts[2])
	if err != nil {
		errmsg := ""
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "RemoveConfigFromMachineGroup", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
	}

	return errmsgs.WrapError(logService.WaitForLogtailAttachment(d.Id(), Deleted, DefaultTimeout))
}
