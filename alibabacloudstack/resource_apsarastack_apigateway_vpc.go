package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackApigatewayVpc() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackApigatewayVpcAccessCreate, resourceAlibabacloudStackApigatewayVpcAccessRead, nil, resourceAlibabacloudStackApigatewayVpcAccessDelete)
	return resource
}

func resourceAlibabacloudStackApigatewayVpcAccessCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := make(map[string]interface{})
	request["Name"] = d.Get("name")
	request["VpcId"] = d.Get("vpc_id")
	request["Port"] = d.Get("port")
	request["InstanceId"] = d.Get("instance_id")
	request["ClientToken"] = buildClientToken("SetVpcAccess")

	_, err = client.DoTeaRequest("POST", "CloudAPI", "2016-07-14", "SetVpcAccess", "", nil, nil, request)
	d.SetId(fmt.Sprintf("%s%s%s%s%s%s%s", request["Name"], COLON_SEPARATED, request["VpcId"], COLON_SEPARATED, request["InstanceId"], COLON_SEPARATED, request["Port"]))
	return nil
}

func resourceAlibabacloudStackApigatewayVpcAccessRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudApiService := CloudApiService{client}

	vpc, err := cloudApiService.DescribeApiGatewayVpcAccess(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("name", vpc.Name)
	d.Set("vpc_id", vpc.VpcId)
	d.Set("instance_id", vpc.InstanceId)
	d.Set("port", vpc.Port)

	return nil
}

func resourceAlibabacloudStackApigatewayVpcAccessDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := cloudapi.CreateRemoveVpcAccessRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.VpcId = d.Get("vpc_id").(string)
	request.InstanceId = d.Get("instance_id").(string)
	request.Port = requests.NewInteger(d.Get("port").(int))

	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.RemoveVpcAccess(request)
	})
	response, ok := raw.(*cloudapi.RemoveVpcAccessResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	return nil
}