package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
	"time"
)

func resourceApsaraStackEcsCommand() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackEcsCommandCreate,
		Read:   resourceApsaraStackEcsCommandRead,
		Delete: resourceApsaraStackEcsCommandDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
}

func resourceApsaraStackEcsCommandCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	//var response map[string]interface{}
	action := "CreateCommand"
	response := &datahub.EcsCreate{}

	//var err error
	//config := client.Config
	//
	//log.Printf(" config.StsEndpoint : %s", config.StsEndpoint)
	//log.Printf(" before get config.RamRoleArn : %s", config.RamRoleArn)
	//
	//if config.RamRoleArn == "" {
	//	config.RamRoleArn = os.Getenv("APSARASTACK_ASSUME_ROLE_ARN")
	//}
	//
	//log.Printf(" after get config.RamRoleArn : %s", config.RamRoleArn)
	//if config.RamRoleArn != "" {
	//	log.Printf("begin getAssumeRoleAK ")
	//	config.AccessKey, config.SecretKey, config.SecurityToken, err = getAssumeRoleAK(config)
	//	if err != nil {
	//		log.Printf(" rsponse of raw ListLoginPolicies : %s", err)
	//		return err
	//	}
	//}

	//request := make(map[string]interface{})
	//conn, err := client.NewEcsClient()
	//if err != nil {
	//	return WrapError(err)
	//}
	CommandContent := d.Get("command_content").(string)

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
	request := requests.NewCommonRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Method = "POST"
	request.Product = "Ecs"
	request.Domain = client.Domain
	request.Version = "2014-05-26"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = action
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "Ecs",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          action,
		"Version":         "2014-05-26",
		"Name":            name,
		"CommandContent":  CommandContent,
		"Description":     Description,
		"EnableParameter": EnableParameter,
		"Timeout":         Timeout,
		"WorkingDir":      WorkingDir,
		"Type":            Type,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
			return EcsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw, request)
		bresponse := raw.(*responses.CommonResponse)
		err = json.Unmarshal(bresponse.GetHttpContentBytes(), response)

		//var response *ecs.CreateCommandResponse
		//response, _ := raw.(*ecs.CreateCommandResponse)
		d.SetId(fmt.Sprint(response.CommandId))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ecs_command", action, ApsaraStackSdkGoERROR)
	}

	return resourceApsaraStackEcsCommandRead(d, meta)
}
func resourceApsaraStackEcsCommandRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsCommand(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource apsarastack_ecs_command ecsService.DescribeEcsCommand Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("command_content", object.Commands.Command[0].CommandContent)
	d.Set("description", object.Commands.Command[0].Description)
	d.Set("enable_parameter", object.Commands.Command[0].EnableParameter)
	d.Set("name", object.Commands.Command[0].Name)
	d.Set("timeout", object.Commands.Command[0].Timeout)
	d.Set("type", object.Commands.Command[0].Type)
	d.Set("working_dir", object.Commands.Command[0].WorkingDir)
	return nil
}
func resourceApsaraStackEcsCommandDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	action := "DeleteCommand"
	//var response map[string]interface{}
	//conn, err := client.NewEcsClient()
	//if err != nil {
	//	return WrapError(err)
	//}
	//request := map[string]interface{}{
	//	"CommandId": d.Id(),
	//}
	request := requests.NewCommonRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Method = "POST"
	request.Product = "Ecs"
	request.Domain = client.Domain
	request.Version = "2014-05-26"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = action
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "Ecs",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          action,
		"Version":         "2014-05-26",
		"CommandId":       d.Id(),
	}

	//request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		//response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		response, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
			return EcsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidCmdId.NotFound", "InvalidRegionId.NotFound", "Operation.Forbidden"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
	}
	return nil
}
