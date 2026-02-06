package alibabacloudstack

import (
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEdasApplicationPackageAttachment() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"package_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"war_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"last_package_version": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEdasApplicationPackageAttachmentCreate,
		resourceAlibabacloudStackEdasApplicationPackageAttachmentRead, nil, resourceAlibabacloudStackEdasApplicationPackageAttachmentDelete)
	return resource
}

func resourceAlibabacloudStackEdasApplicationPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	appId := d.Get("app_id").(string)
	groupId := d.Get("group_id").(string)

	request := edas.CreateDeployApplicationRequest()
	client.InitRoaRequest(*request.RoaRequest)
	request.AppId = appId
	request.GroupId = groupId
	request.WarUrl = d.Get("war_url").(string)

	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	var packageVersion string
	if v, ok := d.GetOk("package_version"); ok {
		packageVersion = v.(string)
	} else {
		packageVersion = strconv.FormatInt(time.Now().Unix(), 10)
	}
	request.DeployType = "url"
	request.PackageVersion = packageVersion
	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.DeployApplication(request)
	})
	if err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*edas.DeployApplicationResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.DeployApplicationResponse)
	changeOrderId := response.ChangeOrderId
	if response.Code != 200 {
		return errmsgs.WrapError(errmsgs.Error("deploy application failed for " + response.Message))
	}

	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	d.SetId(appId + ":" + groupId + ":" + packageVersion)

	return nil
}

func resourceAlibabacloudStackEdasApplicationPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}
	pastr := strings.Split(d.Id(), ":")
	appId := pastr[0]
	groupId := pastr[1]
	package_version := pastr[2]
	d.Set("package_version", package_version)
	d.Set("app_id", appId)
	d.Set("group_id", groupId)
	version, err := edasService.GetLastPackgeVersion(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("last_package_version", version)

	return nil
}

func resourceAlibabacloudStackEdasApplicationPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}
	pastr := strings.Split(d.Id(), ":")
	appId := pastr[0]
	groupId := pastr[1]
	packageVersion := d.Get("last_package_version").(string)

	if len(packageVersion) == 0 {
		return nil
	}

	request := edas.CreateRollbackApplicationRequest()
	client.InitRoaRequest(*request.RoaRequest)
	request.AppId = appId
	request.HistoryVersion = packageVersion
	request.GroupId = groupId

	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.RollbackApplication(request)
	})
	if err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*edas.RollbackApplicationResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.RollbackApplicationResponse)
	changeOrderId := response.ChangeOrderId
	if response.Code != 200 && !strings.Contains(response.Message, "ex.app.deploy.group.empty") {
		return errmsgs.WrapError(errmsgs.Error("deploy application failed for " + response.Message))
	}

	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	return nil
}
