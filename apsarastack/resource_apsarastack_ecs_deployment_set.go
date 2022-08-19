package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"log"
	"regexp"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceApsaraStackEcsDeploymentSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackEcsDeploymentSetCreate,
		Read:   resourceApsaraStackEcsDeploymentSetRead,
		Update: resourceApsaraStackEcsDeploymentSetUpdate,
		Delete: resourceApsaraStackEcsDeploymentSetDelete,
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
				ValidateFunc: validation.StringInSlice([]string{"Default","default"}, false),
			},
			"granularity": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Host","host"}, false),
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
				ValidateFunc: validation.StringInSlice([]string{"Availability"}, false),
			},
		},
	}
}

func resourceApsaraStackEcsDeploymentSetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
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
		"AccessKeySecret":  client.SecretKey,
		"AccessKeyId":      client.AccessKey,
		"Product":          "Ecs",
		"RegionId":         client.RegionId,
		"Department":       client.Department,
		"ResourceGroup":    client.ResourceGroup,
		"Action":           action,
		"Version":          "2014-05-26",
		"DeploymentSetName":DeploymentSetName,
		"Domain":           "Default",
		"Description":      Description,
		"Granularity":      "Host",
		"OnUnableToRedeployFailedInstance":         OnUnableToRedeployFailedInstance,
		"Strategy":         Strategy,
		"ClientToken":      ClientToken,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
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
		resp := &datahub.EcsDeploymentSetCreateResult{}
		bresponse := raw.(*responses.CommonResponse)
		err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
		d.SetId(fmt.Sprint(resp.DeploymentSetId))
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ecs_deployment_set", action, ApsaraStackSdkGoERROR)
	}

	return resourceApsaraStackEcsDeploymentSetRead(d, meta)
}
func resourceApsaraStackEcsDeploymentSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsDeploymentSet(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource apsarastack_ecs_deployment_set ecsService.DescribeEcsDeploymentSet Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("domain", object.DeploymentSets.DeploymentSet[0].Domain)
	d.Set("granularity", object.DeploymentSets.DeploymentSet[0].Granularity)
	d.Set("deployment_set_name", object.DeploymentSets.DeploymentSet[0].DeploymentSetName)
	d.Set("description", object.DeploymentSets.DeploymentSet[0].DeploymentSetDescription)
	d.Set("strategy", object.DeploymentSets.DeploymentSet[0].DeploymentStrategy)
	//d.Set("DeploymentSetId", d.Get("DeploymentSetId"))

	return nil
}
func resourceApsaraStackEcsDeploymentSetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := requests.NewCommonRequest()

	update := false
		DeploymentSetId := d.Id()

	var DeploymentSetName string
	if d.HasChange("deployment_set_name") {
		update = true
		if v, ok := d.GetOk("deployment_set_name"); ok {
			DeploymentSetName  = fmt.Sprint(v.(string))
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
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "Ecs",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          action,
		"Version":         "2014-05-26",
		"DeploymentSetId": DeploymentSetId,
		"DeploymentSetName": DeploymentSetName,
		"Description": Description,
	}
	if update {

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
		}
	}
	return resourceApsaraStackEcsDeploymentSetRead(d, meta)
}
func resourceApsaraStackEcsDeploymentSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	action := "DeleteDeploymentSet"
	DeploymentSetId:= d.Id()
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
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "Ecs",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          action,
		"Version":         "2014-05-26",
		"DeploymentSetId": DeploymentSetId,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
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
	}
	return source
}
