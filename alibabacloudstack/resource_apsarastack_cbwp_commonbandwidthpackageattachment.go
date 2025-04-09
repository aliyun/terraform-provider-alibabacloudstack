package alibabacloudstack

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackCommonBandwidthPackageAttachment() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bandwidth_package_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackCommonBandwidthPackageAttachmentCreate, 
		resourceAlibabacloudStackCommonBandwidthPackageAttachmentRead, nil, 
		resourceAlibabacloudStackCommonBandwidthPackageAttachmentDelete)
	return resource
}

func resourceAlibabacloudStackCommonBandwidthPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	request := vpc.CreateAddCommonBandwidthPackageIpRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.BandwidthPackageId = Trim(d.Get("bandwidth_package_id").(string))
	request.IpInstanceId = Trim(d.Get("instance_id").(string))

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.AddCommonBandwidthPackageIp(request)
	})

	var response *vpc.AddCommonBandwidthPackageIpResponse
	var ok bool
	if raw != nil {
		response, ok = raw.(*vpc.AddCommonBandwidthPackageIpResponse)
	}
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_common_bandwidth_package_attachment", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	//check the common bandwidth package attachment
	d.SetId(request.BandwidthPackageId + COLON_SEPARATED + request.IpInstanceId)
	if err := vpcService.WaitForCommonBandwidthPackageAttachment(d.Id(), Available, 5*DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func resourceAlibabacloudStackCommonBandwidthPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	bandwidthPackageId, ipInstanceId := parts[0], parts[1]
	_, err = vpcService.DescribeCommonBandwidthPackageAttachment(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("bandwidth_package_id", bandwidthPackageId)
	d.Set("instance_id", ipInstanceId)
	return nil
}

func resourceAlibabacloudStackCommonBandwidthPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	bandwidthPackageId, ipInstanceId := parts[0], parts[1]

	request := vpc.CreateRemoveCommonBandwidthPackageIpRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.BandwidthPackageId = bandwidthPackageId
	request.IpInstanceId = ipInstanceId

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.RemoveCommonBandwidthPackageIp(request)
		})

		var response *vpc.RemoveCommonBandwidthPackageIpResponse
		var ok bool
		if raw != nil {
			response, ok = raw.(*vpc.RemoveCommonBandwidthPackageIpResponse)
		}
		//Waiting for unassociate the common bandwidth package
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"TaskConflict"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return err
	}
	return errmsgs.WrapError(vpcService.WaitForCommonBandwidthPackageAttachment(d.Id(), Deleted, DefaultTimeoutMedium))
}
