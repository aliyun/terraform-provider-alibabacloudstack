package alibabacloudstack

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackRouterInterfaceConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackRouterInterfaceConnectionCreate,
		Read:   resourceAlibabacloudStackRouterInterfaceConnectionRead,
		Delete: resourceAlibabacloudStackRouterInterfaceConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"opposite_interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"opposite_router_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(VRouter), string(VBR)}, false),
				Default:  VRouter,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.HasChange("opposite_interface_owner_id")
				},
			},
			"opposite_router_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.HasChange("opposite_interface_owner_id")
				},
			},
			"opposite_interface_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackRouterInterfaceConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	oppositeId := d.Get("opposite_interface_id").(string)
	interfaceId := d.Get("interface_id").(string)
	object, err := vpcService.DescribeRouterInterface(interfaceId, client.RegionId)
	if err != nil {
		return WrapError(err)
	}

	// At present, the interface with "active/inactive" status can not be modify opposite connection information
	// and it is RouterInterface product limitation.
	if object.OppositeInterfaceId == oppositeId {
		if object.Status == string(Active) {
			return WrapError(Error("The specified router interface connection has existed, and please import it using id %s.", interfaceId))
		}
		if object.Status == string(Inactive) {
			if err = vpcService.ActivateRouterInterface(interfaceId); err != nil {
				return WrapError(err)
			}
			d.SetId(object.RouterInterfaceId)
			if err = vpcService.WaitForRouterInterfaceConnection(d.Id(), client.RegionId, Active, DefaultTimeout); err != nil {
				return WrapError(err)
			}
			return resourceAlibabacloudStackRouterInterfaceConnectionRead(d, meta)
		}
	}

	request := vpc.CreateModifyRouterInterfaceAttributeRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}

	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.RouterInterfaceId = interfaceId
	request.OppositeInterfaceId = oppositeId

	if owner_id, ok := d.GetOk("opposite_interface_owner_id"); ok && owner_id.(string) != "" {
		request.OppositeInterfaceOwnerId = requests.Integer(owner_id.(string))
		if v, o := d.GetOk("opposite_router_type"); !o || v.(string) == "" {
			return WrapError(Error("'opposite_router_type' is required when 'opposite_interface_owner_id' is set."))
		} else {
			request.OppositeRouterType = v.(string)
		}

		if v, o := d.GetOk("opposite_router_id"); !o || v.(string) == "" {
			return WrapError(Error("'opposite_router_id' is required when 'opposite_interface_owner_id' is set."))
		} else {
			request.OppositeRouterId = v.(string)
		}
	} else {
		owner := object.OppositeInterfaceOwnerId
		if owner == "" {
			owner, err = client.AccountId()
			if err != nil {
				return WrapError(err)
			}
		}
		if owner == "" {
			return WrapError(Error("Opposite router interface owner id is empty. Please use field 'opposite_interface_owner_id' or globle field 'account_id' to set."))
		}
		oppositeRi, err := vpcService.DescribeRouterInterface(oppositeId, object.OppositeRegionId)
		if err != nil {
			return WrapError(err)
		}
		request.OppositeRouterId = oppositeRi.RouterId
		request.OppositeRouterType = oppositeRi.RouterType
		request.OppositeInterfaceOwnerId = requests.Integer(owner)
	}

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.ModifyRouterInterfaceAttribute(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_router_interface_connection", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	d.SetId(interfaceId)

	if err = vpcService.WaitForRouterInterfaceConnection(d.Id(), client.RegionId, Idle, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	if object.Role == string(InitiatingSide) {
		connectRequest := vpc.CreateConnectRouterInterfaceRequest()
		connectRequest.RegionId = client.RegionId
		if strings.ToLower(client.Config.Protocol) == "https" {
			connectRequest.Scheme = "https"
		} else {
			connectRequest.Scheme = "http"
		}
		connectRequest.Headers = map[string]string{"RegionId": client.RegionId}
		connectRequest.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

		connectRequest.RouterInterfaceId = interfaceId
		if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
			raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.ConnectRouterInterface(connectRequest)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectOppositeInterfaceInfo.NotSet"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(connectRequest.GetActionName(), raw, connectRequest.RpcRequest, connectRequest)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), connectRequest.GetActionName(), AlibabacloudStackSdkGoERROR)
		}

		if err := vpcService.WaitForRouterInterfaceConnection(d.Id(), client.RegionId, Active, DefaultTimeout); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), connectRequest.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
	}

	return resourceAlibabacloudStackRouterInterfaceConnectionRead(d, meta)
}

func resourceAlibabacloudStackRouterInterfaceConnectionRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeRouterInterfaceConnection(d.Id(), client.RegionId)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if object.Status == string(Inactive) {
		if err := vpcService.ActivateRouterInterface(d.Id()); err != nil {
			return WrapError(err)
		}
		if err := vpcService.WaitForRouterInterfaceConnection(d.Id(), client.RegionId, Active, DefaultTimeout); err != nil {
			return WrapError(err)
		}
	}

	d.Set("interface_id", object.RouterInterfaceId)
	d.Set("opposite_interface_id", object.OppositeInterfaceId)
	d.Set("opposite_router_type", object.OppositeRouterType)
	d.Set("opposite_router_id", object.OppositeRouterId)
	d.Set("opposite_interface_owner_id", object.OppositeInterfaceOwnerId)

	return nil

}

func resourceAlibabacloudStackRouterInterfaceConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeRouterInterfaceConnection(d.Id(), client.RegionId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
	}

	if object.Status == string(Idle) {
		d.SetId("")
		return nil
	}

	// At present, the interface with "active/inactive" status can not be modify opposite connection information
	// and it is RouterInterface product limitation. So, the connection delete action is only modifying it to inactive.
	if object.Status == string(Active) {
		if err := vpcService.DeactivateRouterInterface(d.Id()); err != nil {
			return WrapError(err)
		}
	}

	return WrapError(vpcService.WaitForRouterInterfaceConnection(d.Id(), client.RegionId, Inactive, DefaultTimeoutMedium))
}
