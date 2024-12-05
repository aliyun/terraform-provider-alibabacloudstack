package alibabacloudstack

import (
	"fmt"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackApigatewayVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackApigatewayVpcAccessCreate,
		Read:   resourceAlibabacloudStackApigatewayVpcAccessRead,
		Delete: resourceAlibabacloudStackApigatewayVpcAccessDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
}

func resourceAlibabacloudStackApigatewayVpcAccessCreate(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*connectivity.AlibabacloudStackClient)
	//
	//request := cloudapi.CreateSetVpcAccessRequest()
	//
	//request.RegionId = client.RegionId
	//request.ReadTimeout = 600000
	//request.Name = d.Get("name").(string)
	//request.VpcId = d.Get("vpc_id").(string)
	//request.InstanceId = d.Get("instance_id").(string)
	//request.Port = requests.NewInteger(d.Get("port").(int))
	//request.Headers = map[string]string{
	//	"RegionId": client.RegionId,
	//}
	//request.QueryParams = map[string]string{
	//	
	//	
	//	"Product":         "CloudAPI",
	//	"RegionId":        client.RegionId,
	//	"Department":      client.Department,
	//	"ResourceGroup":   client.ResourceGroup,
	//	"Action":          "SetVpcAccess",
	//	"Version":         "2016-07-14",
	//	"SignatureVersion": "2.1",
	//}
	//
	//raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
	//	return cloudApiClient.SetVpcAccess(request)
	//})
	//if err != nil {
	//	return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_api_gateway_vpc", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	//}
	//addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "SetVpcAccess"
	request := make(map[string]interface{})
	request["Product"] = "CloudAPI"
	request["product"] = "CloudAPI"
	request["OrganizationId"] = client.Department
	request["RegionId"] = client.RegionId
	request["Name"] = d.Get("name")
	request["VpcId"] = d.Get("vpc_id")
	request["Port"] = d.Get("port")
	request["InstanceId"] = d.Get("instance_id")
	conn, err := client.NewCloudApiClient()
	if err != nil {
		return WrapError(err)
	}
	request["ClientToken"] = buildClientToken("SetVpcAccess")
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequesttowpoint1(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-07-14"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	d.SetId(fmt.Sprintf("%s%s%s%s%s%s%s", request["Name"], COLON_SEPARATED, request["VpcId"], COLON_SEPARATED, request["InstanceId"], COLON_SEPARATED, request["Port"]))
	return resourceAlibabacloudStackApigatewayVpcAccessRead(d, meta)
}

func resourceAlibabacloudStackApigatewayVpcAccessRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudApiService := CloudApiService{client}

	vpc, err := cloudApiService.DescribeApiGatewayVpcAccess(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
	request.RegionId = client.RegionId
	request.VpcId = d.Get("vpc_id").(string)
	request.InstanceId = d.Get("instance_id").(string)
	request.Port = requests.NewInteger(d.Get("port").(int))
	request.Headers = map[string]string{
		"RegionId": client.RegionId,
	}
	request.QueryParams = map[string]string{
		
		
		"Product":         "CloudAPI",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "RemoveVpcAccess",
		"Version":         "2016-07-14",
	}
	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.RemoveVpcAccess(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return nil

}
