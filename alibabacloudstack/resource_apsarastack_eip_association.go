package alibabacloudstack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEipAssociation() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allocation_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEipAssociationCreate, resourceAlibabacloudStackEipAssociationRead, resourceAlibabacloudStackEipAssociationUpdate, resourceAlibabacloudStackEipAssociationDelete)
	return resource
}

func resourceAlibabacloudStackEipAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateAssociateEipAddressRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.AllocationId = Trim(d.Get("allocation_id").(string))
	request.InstanceId = Trim(d.Get("instance_id").(string))
	request.InstanceType = EcsInstance
	request.ClientToken = buildClientToken(request.GetActionName())

	if strings.HasPrefix(request.InstanceId, "lb-") {
		request.InstanceType = SlbInstance
	}
	if strings.HasPrefix(request.InstanceId, "ngw-") {
		request.InstanceType = Nat
	}
	if instanceType, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = instanceType.(string)
	}
	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.AssociateEipAddress(request)
		})
		var response *vpc.AssociateEipAddressResponse
		var ok bool
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"TaskConflict"}) {
				return resource.RetryableError(err)
			}
			if raw != nil {
				response, ok = raw.(*vpc.AssociateEipAddressResponse)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_eip_association", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return err
	}

	if err := vpcService.WaitForEip(request.AllocationId, InUse, 60); err != nil {
		return errmsgs.WrapError(err)
	}
	// There is at least 30 seconds delay for ecs instance
	if request.InstanceType == EcsInstance {
		time.Sleep(30 * time.Second)
	}

	d.SetId(request.AllocationId + ":" + request.InstanceId)

	return nil
}

func resourceAlibabacloudStackEipAssociationUpdate(d *schema.ResourceData, meta interface{}) error {
	// force是个删除时的动作参数，允许修改
	return nil
}

func resourceAlibabacloudStackEipAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeEipAssociation(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("instance_id", object.InstanceId)
	d.Set("allocation_id", object.AllocationId)
	d.Set("instance_type", object.InstanceType)
	return nil
}

func resourceAlibabacloudStackEipAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	allocationId, instanceId := parts[0], parts[1]

	request := vpc.CreateUnassociateEipAddressRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.AllocationId = allocationId
	request.InstanceId = instanceId
	if force, ok := d.GetOk("force"); ok {
		request.Force = requests.NewBoolean(force.(bool))
	}
	request.InstanceType = EcsInstance
	request.ClientToken = buildClientToken(request.GetActionName())

	if strings.HasPrefix(instanceId, "lb-") {
		request.InstanceType = SlbInstance
	}
	if strings.HasPrefix(instanceId, "ngw-") {
		request.InstanceType = Nat
	}
	if instanceType, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = instanceType.(string)
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.UnassociateEipAddress(request)
		})
		var response *vpc.UnassociateEipAddressResponse
		var ok bool
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"IncorrectInstanceStatus", "IncorrectHaVipStatus", "TaskConflict",
				"InvalidIpStatus.HasBeenUsedBySnatTable", "InvalidIpStatus.HasBeenUsedByForwardEntry", "InvalidStatus.SnatOrDnat"}) {
				return resource.RetryableError(err)
			}
			if raw != nil {
				response, ok = raw.(*vpc.UnassociateEipAddressResponse)
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
	return errmsgs.WrapError(vpcService.WaitForEipAssociation(d.Id(), Available, DefaultTimeoutMedium))
}
