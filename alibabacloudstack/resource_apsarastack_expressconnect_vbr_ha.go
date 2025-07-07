package alibabacloudstack

import (
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackExpressconnectVbrHa() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"vbr_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			
			"peer_vbr_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackExpressconnectVbrHaCreate,
		resourceAlibabacloudStackExpressconnectVbrHaRead,
		nil,
		resourceAlibabacloudStackExpressconnectVbrHaDelete)
	return resource
}

func resourceAlibabacloudStackExpressconnectVbrHaCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"Name": d.Get("name"),
		"VbrId": d.Get("vbr_id"),
		"PeerVbrId": d.Get("peer_vbr_id"),
	}
	
	if v, ok:= d.GetOk("description"); ok {
		reqQuery["Description"] = v
	}
	
	response, err := client.DoTeaRequest("POST", "Vpc", "2016-04-28", "CreateVbrHa", "", nil, reqQuery, nil)
	if err != nil {
		return err
	}
	
	vbrHaId := response["VbrHaId"].(string)
	
	d.SetId(vbrHaId)
	
	expressconnectservice := ExpressconnectService{client}

	stateConf := BuildStateConf([]string{"Creating"}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, expressconnectservice.ExpressconnectVbrHaStateRefreshFunc(vbrHaId, []string{"Failed"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, vbrHaId)
	}
	return nil

}

func resourceAlibabacloudStackExpressconnectVbrHaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	express_connectbgp_peerservice := ExpressconnectService{client}
	vbrHa, err := express_connectbgp_peerservice.DoVpcDescribeVbrHaRequest(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_expressconnect_vbr_ha", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.Set("name", vbrHa["Name"])
	d.Set("vbr_id", vbrHa["VbrId"])
	d.Set("peer_vbr_id", vbrHa["PeerVbrId"])
	d.Set("description", vbrHa["Description"])
	return nil
}

func resourceAlibabacloudStackExpressconnectVbrHaDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("POST", "Vpc", "2016-04-28", "DeleteVbrHa", "")

	//调用request_params_handler

	request.QueryParams["InstanceId"] = d.Id()

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_expressconnect_vbr_ha", "DeleteVbrHa", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	expressconnectservice := ExpressconnectService{client}
	stateConf := BuildStateConf([]string{"Deleting"}, []string{}, d.Timeout(schema.TimeoutCreate), 10*time.Second, expressconnectservice.ExpressconnectVbrHaStateRefreshFunc(d.Id(), []string{"Failed"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}
