package alibabacloudstack

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/sdk_patch/datahub_patch"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackEcsCommand() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"command_content": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enable_parameter": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  60,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"RunBatScript", "RunPowerShellScript", "RunShellScript"}, false),
			},
			"working_dir": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEcsCommandCreate, resourceAlibabacloudStackEcsCommandRead, nil, resourceAlibabacloudStackEcsCommandDelete)
	return resource
}

func resourceAlibabacloudStackEcsCommandCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	response := &datahub_patch.EcsCreate{}

	CommandContent := d.Get("command_content").(string)
	contentBytes := []byte(CommandContent)
	encoded := base64.StdEncoding.EncodeToString(contentBytes)
	var Description string
	if v, ok := d.GetOk("description"); ok {
		Description = fmt.Sprint(v.(string))
	}
	var EnableParameter string
	if v, ok := d.GetOkExists("enable_parameter"); ok {
		EnableParameter = fmt.Sprint(v.(bool))
	}
	name := d.Get("name").(string)
	var Timeout string
	if v, ok := d.GetOk("timeout"); ok {
		Timeout = fmt.Sprint(v.(int))
	}
	Type := d.Get("type").(string)
	var WorkingDir string
	if v, ok := d.GetOk("working_dir"); ok {
		WorkingDir = fmt.Sprint(v.(string))
	}

	request := client.NewCommonRequest("POST", "Ecs", "2014-05-26", "CreateCommand", "")
	mergeMaps(request.QueryParams, map[string]string{
		"Name":            name,
		"CommandContent":  encoded,
		"Description":     Description,
		"EnableParameter": EnableParameter,
		"Timeout":         Timeout,
		"WorkingDir":      WorkingDir,
		"Type":            Type,
	})

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if bresponse == nil {
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, "Process Common Request Failed"))
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ecs_command", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		err = json.Unmarshal(bresponse.GetHttpContentBytes(), response)
		d.SetId(fmt.Sprint(response.CommandId))
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceAlibabacloudStackEcsCommandRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsCommand(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_ecs_command ecsService.DescribeEcsCommand Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	command_content := object.Commands.Command[0].CommandContent
	decodedBytes, err := base64.StdEncoding.DecodeString(command_content)
	if err == nil {
		command_content = string(decodedBytes)
	}
	d.Set("command_content", command_content)
	d.Set("description", object.Commands.Command[0].Description)
	d.Set("enable_parameter", object.Commands.Command[0].EnableParameter)
	d.Set("name", object.Commands.Command[0].Name)
	d.Set("timeout", object.Commands.Command[0].Timeout)
	d.Set("type", object.Commands.Command[0].Type)
	d.Set("working_dir", object.Commands.Command[0].WorkingDir)
	return nil
}

func resourceAlibabacloudStackEcsCommandDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("POST", "Ecs", "2014-05-26", "DeleteCommand", "")
	request.QueryParams["CommandId"] = d.Id()

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if bresponse == nil {
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, "Process Common Request Failed"))
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ecs_command", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidCmdId.NotFound", "InvalidRegionId.NotFound", "Operation.Forbidden"}) {
			return nil
		}
		return err
	}
	return nil
}
