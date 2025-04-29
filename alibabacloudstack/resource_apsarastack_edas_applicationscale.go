package alibabacloudstack

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEdasInstanceApplicationAttachment() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ecc_info": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"deploy_group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"force_status": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"ecu_info": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEdasInstanceApplicationAttachmentCreate, resourceAlibabacloudStackEdasInstanceApplicationAttachmentRead, nil, resourceAlibabacloudStackEdasInstanceApplicationAttachmentDelete)
	return resource
}

func resourceAlibabacloudStackEdasInstanceApplicationAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	appId := d.Get("app_id").(string)
	ecuInfo := d.Get("ecu_info").([]interface{})
	aString := make([]string, len(ecuInfo))
	for i, v := range ecuInfo {
		if v != nil {
			aString[i] = v.(string)
		}
	}

	request := edas.CreateScaleOutApplicationRequest()
	client.InitRoaRequest(*request.RoaRequest)
	request.AppId = appId
	request.DeployGroup = d.Get("deploy_group").(string)
	request.EcuInfo = strings.Join(aString, ",")

	var changeOrderId string

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ScaleOutApplication(request)
	})

	bresponse, ok := raw.(*edas.ScaleOutApplicationResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_instance_application_attachment", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	changeOrderId = bresponse.ChangeOrderId
	d.SetId(appId + ":" + strings.Join(aString, ","))
	if bresponse.Code != 200 {
		return errmsgs.WrapError(errmsgs.Error("scaleOut application failed for " + bresponse.Message))
	}

	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	return nil
}

func resourceAlibabacloudStackEdasInstanceApplicationAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	strs, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	appId := strs[0]
	ecuInfo := strs[1]
	aString := strings.Split(ecuInfo, ",")

	request := edas.CreateQueryApplicationStatusRequest()
	client.InitRoaRequest(*request.RoaRequest)
	request.AppId = appId

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.QueryApplicationStatus(request)
	})
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	if err != nil {
		errmsg := ""
		if raw != nil {
			bresponse, ok := raw.(*edas.QueryApplicationStatusResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_instance_application_attachment", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	var eccs []string
	bresponse := raw.(*edas.QueryApplicationStatusResponse)
	if bresponse.Code != 200 {
		return errmsgs.WrapError(errmsgs.Error("QueryApplicationStatus failed for " + bresponse.Message))
	}
	for _, ecc := range bresponse.AppInfo.EccList.Ecc {
		for _, ecu := range aString {
			if ecu == ecc.EcuId {
				eccs = append(eccs, ecc.EccId)
			}
		}
	}
	if eccs != nil {
		d.Set("ecc_info", strings.Join(eccs, ","))
	}
	return nil
}

func resourceAlibabacloudStackEdasInstanceApplicationAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	request := edas.CreateScaleInApplicationRequest()
	client.InitRoaRequest(*request.RoaRequest)
	request.AppId = d.Get("app_id").(string)
	request.EccInfo = d.Get("ecc_info").(string)

	if v, ok := d.GetOk("force_status"); ok {
		request.ForceStatus = requests.NewBoolean(v.(bool))
	}

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ScaleInApplication(request)
	})

	bresponse, ok := raw.(*edas.ScaleInApplicationResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_instance_application_attachment", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	if bresponse.Code != 200 {
		return errmsgs.WrapError(errmsgs.Error("scaleIn application failed for " + bresponse.Message))
	}
	changeOrderId := bresponse.ChangeOrderId

	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	return nil
}