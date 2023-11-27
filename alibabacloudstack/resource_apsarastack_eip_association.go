package alibabacloudstack

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEipAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackEipAssociationCreate,
		Read:   resourceAlibabacloudStackEipAssociationRead,
		Delete: resourceAlibabacloudStackEipAssociationDelete,

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
				Default:  true,
				ForceNew: true,
			},

			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackEipAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateAssociateEipAddressRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
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
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_eip_association", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	if err := vpcService.WaitForEip(request.AllocationId, InUse, 60); err != nil {
		return WrapError(err)
	}
	// There is at least 30 seconds delay for ecs instance
	if request.InstanceType == EcsInstance {
		time.Sleep(30 * time.Second)
	}

	d.SetId(request.AllocationId + ":" + request.InstanceId)

	return resourceAlibabacloudStackEipAssociationRead(d, meta)
}

func resourceAlibabacloudStackEipAssociationRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeEipAssociation(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("instance_id", object.InstanceId)
	d.Set("allocation_id", object.AllocationId)
	d.Set("instance_type", object.InstanceType)
	d.Set("force", d.Get("force").(bool))
	return nil
}

func resourceAlibabacloudStackEipAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	allocationId, instanceId := parts[0], parts[1]
	if err != nil {
		return WrapError(err)
	}

	request := vpc.CreateUnassociateEipAddressRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.AllocationId = allocationId
	request.InstanceId = instanceId
	request.Force = requests.NewBoolean(d.Get("force").(bool))
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
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectInstanceStatus", "IncorrectHaVipStatus", "TaskConflict",
				"InvalidIpStatus.HasBeenUsedBySnatTable", "InvalidIpStatus.HasBeenUsedByForwardEntry", "InvalidStatus.SnatOrDnat"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	return WrapError(vpcService.WaitForEipAssociation(d.Id(), Available, DefaultTimeoutMedium))
}
