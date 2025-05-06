package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackLogMachineGroup() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"identify_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      sls.MachineIDTypeIP,
				ValidateFunc: validation.StringInSlice([]string{sls.MachineIDTypeIP, sls.MachineIDTypeUserDefined}, false),
			},
			"topic": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"identify_list": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				MinItems: 1,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackLogMachineGroupCreate, resourceAlibabacloudStackLogMachineGroupRead, resourceAlibabacloudStackLogMachineGroupUpdate, resourceAlibabacloudStackLogMachineGroupDelete)
	return resource
}

func resourceAlibabacloudStackLogMachineGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	params := &sls.MachineGroup{
		Name:          d.Get("name").(string),
		MachineIDType: d.Get("identify_type").(string),
		MachineIDList: expandStringList(d.Get("identify_list").(*schema.Set).List()),
		Attribute: sls.MachinGroupAttribute{
			TopicName: d.Get("topic").(string),
		},
	}
	var requestInfo *sls.Client
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlsDataClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.CreateMachineGroup(d.Get("project").(string), params)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("CreateMachineGroup", raw, requestInfo, map[string]interface{}{
				"project":      d.Get("project").(string),
				"MachineGroup": params,
			})
		}
		return nil
	}); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s%s%s", d.Get("project").(string), COLON_SEPARATED, d.Get("name").(string)))

	return nil
}

func resourceAlibabacloudStackLogMachineGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	object, err := logService.DescribeLogMachineGroup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("project", parts[0])
	d.Set("name", object.Name)
	d.Set("identify_type", object.MachineIDType)
	d.Set("identify_list", object.MachineIDList)
	d.Set("topic", object.Attribute.TopicName)

	return nil
}

func resourceAlibabacloudStackLogMachineGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChanges("identify_type", "identify_list", "topic") {
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		client := meta.(*connectivity.AlibabacloudStackClient)
		var requestInfo *sls.Client
		params := &sls.MachineGroup{
			Name:          parts[1],
			MachineIDType: d.Get("identify_type").(string),
			MachineIDList: expandStringList(d.Get("identify_list").(*schema.Set).List()),
			Attribute: sls.MachinGroupAttribute{
				TopicName: d.Get("topic").(string),
			},
		}
		if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
			raw, err := client.WithSlsDataClient(func(slsClient *sls.Client) (interface{}, error) {
				requestInfo = slsClient
				return nil, slsClient.UpdateMachineGroup(parts[0], params)
			})
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{errmsgs.LogClientTimeout}) {
					time.Sleep(5 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			if debugOn() {
				addDebug("UpdateMachineGroup", raw, requestInfo, map[string]interface{}{
					"project":      parts[0],
					"MachineGroup": params,
				})
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func resourceAlibabacloudStackLogMachineGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	var requestInfo *sls.Client
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlsDataClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.DeleteMachineGroup(parts[0], parts[1])
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("DeleteMachineGroup", raw, requestInfo, map[string]interface{}{
				"project":      parts[0],
				"machineGroup": parts[1],
			})
		}
		return nil
	})
	if err != nil {
		return err
	}
	return errmsgs.WrapError(logService.WaitForLogMachineGroup(d.Id(), Deleted, DefaultTimeout))
}