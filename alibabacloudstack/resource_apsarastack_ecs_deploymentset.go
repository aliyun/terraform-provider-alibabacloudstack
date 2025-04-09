package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/sdk_patch/datahub_patch"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackEcsDeploymentSet() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"deployment_set_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^([\w\\:\-]){2,128}$`), "\t\nThe name of the deployment set.\n\nThe name must be 2 to 128 characters in length and can contain letters, digits, colons (:), underscores (_), and hyphens (-)."),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"domain": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Default", "default"}, false),
			},
			"granularity": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Host", "Rack", "Switch"}, false),
			},
			"on_unable_to_redeploy_failed_instance": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"CancelMembershipAndStart", "KeepStopped"}, false),
			},
			"strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Availability", "LooseDispersion"}, false),
			},
			"tags": tagsSchema(),
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEcsDeploymentSetCreate, resourceAlibabacloudStackEcsDeploymentSetRead, resourceAlibabacloudStackEcsDeploymentSetUpdate, resourceAlibabacloudStackEcsDeploymentSetDelete)
	return resource
}

func resourceAlibabacloudStackEcsDeploymentSetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "CreateDeploymentSet"

	var DeploymentSetName string
	if v, ok := d.GetOk("deployment_set_name"); ok {
		DeploymentSetName = fmt.Sprint(v.(string))
	}
	var Description string
	if v, ok := d.GetOk("description"); ok {
		Description = fmt.Sprint(v.(string))
	}
	var OnUnableToRedeployFailedInstance string
	if v, ok := d.GetOk("on_unable_to_redeploy_failed_instance"); ok {
		OnUnableToRedeployFailedInstance = fmt.Sprint(v.(string))
	}
	var Strategy string
	if v, ok := d.GetOk("strategy"); ok {
		Strategy = fmt.Sprint(v.(string))
	}
	ClientToken := buildClientToken("CreateDeploymentSet")

	request := client.NewCommonRequest("POST", "Ecs", "2014-05-26", action, "")
	mergeMaps(request.QueryParams, map[string]string{
		"DeploymentSetName":                DeploymentSetName,
		"Domain":                           "Default",
		"Description":                      Description,
		"Granularity":                      d.Get("granularity").(string),
		"OnUnableToRedeployFailedInstance": OnUnableToRedeployFailedInstance,
		"Strategy":                         Strategy,
		"ClientToken":                      ClientToken,
	})
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)
	raw, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
		return EcsClient.ProcessCommonRequest(request)
	})
	addDebug(action, raw, request)
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ecs_deployment_set", action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	resp := &datahub_patch.EcsDeploymentSetCreateResult{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ecs_deployment_set", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	d.SetId(fmt.Sprint(resp.DeploymentSetId))

	return nil
}

func resourceAlibabacloudStackEcsDeploymentSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsDeploymentSet(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_ecs_deployment_set ecsService.DescribeEcsDeploymentSet Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("domain", convertEcsDeploymentSetDomainResponse(object["Domain"]))
	d.Set("granularity", convertEcsDeploymentSetGranularityResponse(object["Granularity"]))
	d.Set("deployment_set_name", object["DeploymentSetName"])
	d.Set("description", object["DeploymentSetDescription"])
	d.Set("strategy", object["DeploymentStrategy"])

	if object["Tags"] != nil {
		tags := object["Tags"].(map[string]interface{})["Tag"]
		d.Set("tags", tagsToMap(tags))
	}

	return nil
}

func resourceAlibabacloudStackEcsDeploymentSetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	err := setTags(client, "deployment_set", d)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	update := false
	DeploymentSetId := d.Id()
	Description := ""
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			Description = fmt.Sprint(v.(string))
		}
	}
	action := "ModifyDeploymentSetAttribute"

	request := client.NewCommonRequest("POST", "Ecs", "2014-05-26", action, "")
	request.QueryParams["DeploymentSetId"] = DeploymentSetId
	request.QueryParams["Description"] = Description
	if update {
		response, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
			return EcsClient.ProcessCommonRequest(request)
		})
		addDebug(action, response, request, request.QueryParams)
		bresponse, ok := response.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}
	return nil
}

func resourceAlibabacloudStackEcsDeploymentSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteDeploymentSet"
	DeploymentSetId := d.Id()

	request := client.NewCommonRequest("POST", "Ecs", "2014-05-26", action, "")
	request.QueryParams["DeploymentSetId"] = DeploymentSetId
	response, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
		return EcsClient.ProcessCommonRequest(request)
	})
	addDebug(action, response, request)
	bresponse, ok := response.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	return nil
}

func convertEcsDeploymentSetDomainResponse(source interface{}) interface{} {
	switch source {
	case "default":
		return "Default"
	}
	return source
}

func convertEcsDeploymentSetGranularityResponse(source interface{}) interface{} {
	switch source {
	case "host":
		return "Host"
	case "rack":
		return "Rack"
	case "switch":
		return "Switch"
	}
	return source
}