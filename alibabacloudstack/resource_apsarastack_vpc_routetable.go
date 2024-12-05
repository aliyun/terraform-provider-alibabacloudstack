package alibabacloudstack

import (
	"reflect"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunRouteTableCreate,
		Read:   resourceAliyunRouteTableRead,
		Update: resourceAliyunRouteTableUpdate,
		Delete: resourceAliyunRouteTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
				Deprecated:   "Field 'name' is deprecated and will be removed in a future release. Please use 'route_table_name' instead.",
			},
			"route_table_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliyunRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateCreateRouteTableRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.VpcId = d.Get("vpc_id").(string)
	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "route_table_name", "name"); err == nil {
		request.RouteTableName = v.(string)
	} else {
		return err
	}
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

	return resourceAliyunRouteTableUpdate(d, meta)
}

func resourceAliyunRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)

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
		return resourceAliyunRouteTableRead(d, meta)
	}
	request := vpc.CreateModifyRouteTableAttributesRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.RouteTableId = d.Id()

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
	}

	if d.HasChange("name") || d.HasChange("route_table_name") {
		if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "route_table_name", "name"); err == nil {
			request.RouteTableName = v.(string)
		} else {
			return err
		}
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

	return resourceAliyunRouteTableRead(d, meta)
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
