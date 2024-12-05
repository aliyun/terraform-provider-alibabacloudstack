package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackEcsDeploymentSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackEcsDeploymentSetCreate,
		Read:   resourceAlibabacloudStackEcsDeploymentSetRead,
		Update: resourceAlibabacloudStackEcsDeploymentSetUpdate,
		Delete: resourceAlibabacloudStackEcsDeploymentSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"deployment_set_name": {
				Type:         schema.TypeString,
				Optional:     true,
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
}

func resourceAlibabacloudStackEcsDeploymentSetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	//var response map[string]interface{}
	action := "CreateDeploymentSet"
	//request := make(map[string]interface{})
	//conn, err := client.NewEcsClient()

	var DeploymentSetName string
	if v, ok := d.GetOk("deployment_set_name"); ok {
		DeploymentSetName = fmt.Sprint(v.(string))
	}
	var Description string
	if v, ok := d.GetOk("description"); ok {
		Description = fmt.Sprint(v.(string))
	}
	//var Domain string
	//if v, ok := d.GetOk("domain"); ok {
	//	Domain = fmt.Sprint(v.(string))
	//}
	//var Granularity string
	//if v, ok := d.GetOk("granularity"); ok {
	//	Granularity = fmt.Sprint(v.(string))
	//}
	var OnUnableToRedeployFailedInstance string
	if v, ok := d.GetOk("on_unable_to_redeploy_failed_instance"); ok {
		OnUnableToRedeployFailedInstance = fmt.Sprint(v.(string))
	}
	//request["RegionId"] = client.RegionId
	var Strategy string
	if v, ok := d.GetOk("strategy"); ok {
		Strategy = fmt.Sprint(v.(string))
	}
	ClientToken := buildClientToken("CreateDeploymentSet")
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
		
		
		"Product":                          "Ecs",
		"RegionId":                         client.RegionId,
		"Department":                       client.Department,
		"ResourceGroup":                    client.ResourceGroup,
		"Action":                           action,
		"Version":                          "2014-05-26",
		"DeploymentSetName":                DeploymentSetName,
		"Domain":                           "Default",
		"Description":                      Description,
		"Granularity":                      d.Get("granularity").(string),
		"OnUnableToRedeployFailedInstance": OnUnableToRedeployFailedInstance,
		"Strategy":                         Strategy,
		"ClientToken":                      ClientToken,
	}
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)
	raw, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
		return EcsClient.ProcessCommonRequest(request)
	})
	addDebug(action, raw, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ecs_deployment_set", action, AlibabacloudStackSdkGoERROR)
	}
	resp := &datahub.EcsDeploymentSetCreateResult{}
	bresponse := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ecs_deployment_set", action, AlibabacloudStackSdkGoERROR)
	}
	d.SetId(fmt.Sprint(resp.DeploymentSetId))

	return resourceAlibabacloudStackEcsDeploymentSetRead(d, meta)
}
func resourceAlibabacloudStackEcsDeploymentSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsDeploymentSet(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_ecs_deployment_set ecsService.DescribeEcsDeploymentSet Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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

	ecsService := EcsService{client}
	if d.HasChange("tags") {
		if err := ecsService.SetResourceTagsNew(d, "deployment_set"); err != nil {
			return WrapError(err)
		}
	}

	request := requests.NewCommonRequest()

	update := false
	DeploymentSetId := d.Id()

	var DeploymentSetName string
	if d.HasChange("deployment_set_name") {
		update = true
		if v, ok := d.GetOk("deployment_set_name"); ok {
			DeploymentSetName = fmt.Sprint(v.(string))
		}
	}
	//var Description string
	//if d.HasChange("description") {
	//	update = true
	//	if v, ok := d.GetOk("description"); ok {
	//		Description  = fmt.Sprint(v.(string))
	//	}
	//}
	Description := d.Get("description").(string)
	action := "ModifyDeploymentSetAttribute"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Method = "POST"
	request.Product = "Ecs"
	request.Domain = client.Domain
	request.Version = "2014-05-26"
	request.ApiName = action
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		
		
		"Product":           "Ecs",
		"RegionId":          client.RegionId,
		"Department":        client.Department,
		"ResourceGroup":     client.ResourceGroup,
		"Action":            action,
		"Version":           "2014-05-26",
		"DeploymentSetId":   DeploymentSetId,
		"DeploymentSetName": DeploymentSetName,
		"Description":       Description,
	}
	if update {

		response, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
			return EcsClient.ProcessCommonRequest(request)
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
		}
	}
	return resourceAlibabacloudStackEcsDeploymentSetRead(d, meta)
}
func resourceAlibabacloudStackEcsDeploymentSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteDeploymentSet"
	DeploymentSetId := d.Id()
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
	request.ApiName = action
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		
		
		"Product":         "Ecs",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          action,
		"Version":         "2014-05-26",
		"DeploymentSetId": DeploymentSetId,
	}
	response, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
		return EcsClient.ProcessCommonRequest(request)
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
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
