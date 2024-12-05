package alibabacloudstack

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEdasDeployGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackEdasDeployGroupCreate,
		Read:   resourceAlibabacloudStackEdasDeployGroupRead,
		Delete: resourceAlibabacloudStackEdasDeployGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_type": {
				Type:     schema.TypeInt,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlibabacloudStackEdasDeployGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	appId := d.Get("app_id").(string)
	groupName := d.Get("group_name").(string)

	request := edas.CreateInsertDeployGroupRequest()
	client.InitRoaRequest(*request.RoaRequest)
	request.AppId = appId
	request.GroupName = groupName

	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	wait := incrementalWait(1*time.Second, 2*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.InsertDeployGroup(request)
		})
		bresponse, ok := raw.(*edas.InsertDeployGroupResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_deploy_group", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		deployGroup := bresponse.DeployGroupEntity
		d.SetId(appId + ":" + groupName + ":" + deployGroup.Id)
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		return nil
	})

	if err != nil {
		return err
	}

	return resourceAlibabacloudStackEdasDeployGroupRead(d, meta)
}

func resourceAlibabacloudStackEdasDeployGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	strs, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	appId := strs[0]
	groupId := strs[2]

	deployGroup, err := edasService.GetDeployGroup(appId, groupId)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if deployGroup == nil {
		return nil
	}

	d.Set("group_type", deployGroup.GroupType)
	d.Set("app_id", deployGroup.AppId)
	d.Set("group_name", deployGroup.GroupName)

	return nil
}

func resourceAlibabacloudStackEdasDeployGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	request := edas.CreateDeleteDeployGroupRequest()
	client.InitRoaRequest(*request.RoaRequest)
	request.AppId = d.Get("app_id").(string)
	request.GroupName = d.Get("group_name").(string)

	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	wait := incrementalWait(1*time.Second, 2*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.DeleteDeployGroup(request)
		})
		bresponse, ok := raw.(*edas.DeleteDeployGroupResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
