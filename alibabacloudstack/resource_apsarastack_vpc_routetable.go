package alibabacloudstack

import (
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackRouteTable() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(2, 128),
				Deprecated:    "Field 'name' is deprecated and will be removed in a future release. Please use new field 'route_table_name' instead.",
				ConflictsWith: []string{"route_table_name"},
			},
			"route_table_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"name"},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
		},
	}
	setResourceFunc(resource, resourceAliyunRouteTableCreate, resourceAliyunRouteTableRead, resourceAliyunRouteTableUpdate, resourceAliyunRouteTableDelete)
	return resource
}

func resourceAliyunRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateCreateRouteTableRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.VpcId = d.Get("vpc_id").(string)
	request.RouteTableName = d.Get("name").(string)
	request.Description = d.Get("description").(string)
	request.ClientToken = buildClientToken(request.GetActionName())

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.CreateRouteTable(request)
	})
	bresponse, ok := raw.(*vpc.CreateRouteTableResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_route_table", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	d.SetId(bresponse.RouteTableId)

	if err := vpcService.WaitForRouteTable(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func resourceAliyunRouteTableRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeRouteTable(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("vpc_id", object.VpcId)
	connectivity.SetResourceData(d, object.RouteTableName, "route_table_name", "name")
	d.Set("description", object.Description)
	d.Set("tags", vpcTagsToMap(object.Tags.Tag))
	return nil
}

func resourceAliyunRouteTableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	if err := vpcService.setInstanceTags(d, TagResourceRouteTable); err != nil {
		return errmsgs.WrapError(err)
	}
	if d.IsNewResource() {
		d.Partial(false)
		return nil
	}
	request := vpc.CreateModifyRouteTableAttributesRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.RouteTableId = d.Id()

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
	}

	if d.HasChanges("name", "route_table_name") {
		request.RouteTableName = connectivity.GetResourceData(d, "route_table_name", "name").(string)
	}

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.ModifyRouteTableAttributes(request)
	})
	bresponse, ok := raw.(*vpc.ModifyRouteTableAttributesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return nil
}

func resourceAliyunRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	routeTableService := VpcService{client}
	request := vpc.CreateDeleteRouteTableRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.RouteTableId = d.Id()

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeleteRouteTable(request)
	})
	bresponse, ok := raw.(*vpc.DeleteRouteTableResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return errmsgs.WrapError(routeTableService.WaitForRouteTable(d.Id(), Deleted, DefaultTimeoutMedium))
}