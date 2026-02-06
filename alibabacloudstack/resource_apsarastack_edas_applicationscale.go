package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

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
	groupId := d.Get("deploy_group").(string)
	aString := make([]string, len(ecuInfo))
	for i, v := range ecuInfo {
		if v != nil {
			aString[i] = v.(string)
		}
	}

	request := edas.CreateScaleOutApplicationRequest()
	client.InitRoaRequest(*request.RoaRequest)
	request.AppId = appId
	request.DeployGroup = groupId
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
	d.SetId(appId + ":" + groupId + ":" + strings.Join(aString, ","))
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

	strs, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	appId := strs[0]
	groupId := strs[1]
	ecuInfo := strs[2]
	aString := strings.Split(ecuInfo, ",")
	d.Set("app_id", appId)
	d.Set("ecu_info", aString)
	d.Set("deploy_group", groupId)
	appinfo, err := edasService.DescribeApplicationStatus(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	eccs := make([]string, 0)
	for _, ecc := range appinfo.EccList.Ecc {
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
	strs, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	appId := strs[0]
	req := map[string]interface{}{
		"AppId":   appId,
		"EccInfo": d.Get("ecc_info").(string),
	}
	if v, ok := d.GetOk("force_status"); ok {
		req["ForceStatus"] = v
	}
	response, err := client.DoTeaRequest("POST", "Edas", "2017-08-01", "ScaleInApplication", "/pop/v5/changeorder/co_scale_in", nil, req, nil)

	if err != nil {
		errmsg := errmsgs.GetAsapiErrorMessage(response)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_instance_application_attachment", "ScaleInApplication", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	if fmt.Sprint(response["Code"]) != "200" {
		return errmsgs.WrapError(errmsgs.Error("scaleIn application failed for " + response["Message"].(string)))
	}

	changeOrderId, ok := response["ChangeOrderId"].(string)

	if ok && changeOrderId != "" {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	return nil
}
