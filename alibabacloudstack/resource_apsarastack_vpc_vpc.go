package alibabacloudstack

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackVpcCreate,
		Read:   resourceAlibabacloudStackVpcRead,
		Update: resourceAlibabacloudStackVpcUpdate,
		Delete: resourceAlibabacloudStackVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cidr_block": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "172.16.0.0/12",
				ValidateFunc:  validateCIDRNetworkAddress,
				ConflictsWith: []string{"enable_ipv6"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' has been deprecated from provider version 1.119.0. New field 'vpc_name' instead.",
				ConflictsWith: []string{"vpc_name"},
				ValidateFunc:  validateNormalName,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"router_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"enable_ipv6": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"cidr_block"},
			},
			"ipv6_cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"router_table_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Attribute router_table_id has been deprecated and replaced with route_table_id.",
			},
			"secondary_cidr_blocks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"user_cidrs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
			},
			"vpc_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validateNormalName,
			},
		},
	}
}

func resourceAlibabacloudStackVpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	var response *vpc.CreateVpcResponse
	request := buildAlibabacloudStackVpcArgs(d, meta)
	client.InitRpcRequest(*request.RpcRequest)

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateVpc(&args)
		})
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"TaskConflict", "UnknownError", errmsgs.Throttling}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			errmsg := ""
			bresponse, ok := raw.(*vpc.CreateVpcResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_vpc", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*vpc.CreateVpcResponse)
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(response.VpcId)

	stateConf := BuildStateConf([]string{"Pending"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, vpcService.VpcStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return resourceAlibabacloudStackVpcUpdate(d, meta)
}

func resourceAlibabacloudStackVpcRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeVpc(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_vpc_vpc vpcService.DescribeVpc Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("cidr_block", object.CidrBlock)
<<<<<<< HEAD
	// d.Set("name", object.VpcName)
=======
>>>>>>> origin/rebase_add_testcase_code
	d.Set("description", object.Description)
	d.Set("router_id", object.VRouterId)
	d.Set("ipv6_cidr_block", object.Ipv6CidrBlock)
	d.Set("secondary_cidr_blocks", object.SecondaryCidrBlocks.SecondaryCidrBlock)
	d.Set("status", object.Status)
	if tag := object.Tags.Tag; tag != nil {
		d.Set("tags", vpcService.tagToMap(tag))
	}
	d.Set("user_cidrs", object.UserCidrs.UserCidr)
<<<<<<< HEAD
	connectivity.SetResourceData(d, object.VpcName, "name", "vpc_name")
	// d.Set("vpc_name", object.VpcName)
=======
	connectivity.SetResourceData(d, object.VpcName ,"vpc_name", "name")
>>>>>>> origin/rebase_add_testcase_code

	request := vpc.CreateDescribeRouteTablesRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.VRouterId = object.VRouterId
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	var routeTabls []vpc.RouteTable
	for {
		total := 0
		err = resource.Retry(6*time.Minute, func() *resource.RetryError {
			raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeRouteTables(request)
			})
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{errmsgs.Throttling}) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(err)
				} else {
					errmsg := ""
					bresponse, ok := raw.(*vpc.DescribeRouteTablesResponse)
					if ok {
						errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
					}
					return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
				}
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			response, _ := raw.(*vpc.DescribeRouteTablesResponse)
			routeTabls = append(routeTabls, response.RouteTables.RouteTable...)
			total = len(response.RouteTables.RouteTable)
			return nil
		})
		if err != nil {
			return err
		}

		if total < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return errmsgs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	// Generally, the system route table is the last one
	for i := len(routeTabls) - 1; i >= 0; i-- {
		if routeTabls[i].RouteTableType == "System" {
			d.Set("route_table_id", routeTabls[i].RouteTableId)
			d.Set("router_table_id", routeTabls[i].RouteTableId)
			break
		}
	}

	return nil
}

func resourceAlibabacloudStackVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	if err := vpcService.setInstanceSecondaryCidrBlocks(d); err != nil {
		return errmsgs.WrapError(err)
	}
	if d.HasChange("tags") {
		if err := vpcService.SetResourceTags(d, "vpc"); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlibabacloudStackVpcRead(d, meta)
	}

	groupRequestUpdate := false
	groupRequest := vpc.CreateMoveResourceGroupRequest()
	client.InitRpcRequest(*groupRequest.RpcRequest)
	groupRequest.ResourceId = d.Id()
	groupRequest.NewResourceGroupId = d.Get("resource_group_id").(string)
	groupRequest.ResourceType = "vpc"
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		groupRequestUpdate = true
	}
	if groupRequestUpdate {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.MoveResourceGroup(groupRequest)
		})
		if err != nil {
			errmsg := ""
			bresponse, ok := raw.(*vpc.MoveResourceGroupResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), groupRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(groupRequest.GetActionName(), raw, groupRequest.RpcRequest, groupRequest)
	}

	attributeUpdate := false
	request := vpc.CreateModifyVpcAttributeRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.VpcId = d.Id()

	if d.HasChanges("name","vpc_name") {
		request.VpcName = connectivity.GetResourceData(d,"vpc_name", "name").(string)
		attributeUpdate = true
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		attributeUpdate = true
	}

	if !d.IsNewResource() && d.HasChange("cidr_block") {
		request.CidrBlock = d.Get("cidr_block").(string)
		attributeUpdate = true
	}
	enable_ipv6 := d.Get("enable_ipv6").(bool)
	if attributeUpdate {
		if enable_ipv6 {
			request.EnableIPv6 = requests.NewBoolean(d.Get("enable_ipv6").(bool))
		}
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyVpcAttribute(request)
		})
		if err != nil {
			errmsg := ""
			bresponse, ok := raw.(*vpc.ModifyVpcAttributeResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAlibabacloudStackVpcRead(d, meta)
}

func resourceAlibabacloudStackVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	request := vpc.CreateDeleteVpcRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.VpcId = d.Id()
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVpc(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"InvalidVpcID.NotFound", "Forbidden.VpcNotFound"}) {
				return nil
			}
			errmsg := ""
			bresponse, ok := raw.(*vpc.DeleteVpcResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return err
	}
	stateConf := BuildStateConf([]string{"Pending"}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Second, vpcService.VpcStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}

func buildAlibabacloudStackVpcArgs(d *schema.ResourceData, meta interface{}) *vpc.CreateVpcRequest {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := vpc.CreateCreateVpcRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.CidrBlock = d.Get("cidr_block").(string)

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request.DryRun = requests.NewBoolean(v.(bool))
	}

	request.EnableIpv6 = requests.NewBoolean(d.Get("enable_ipv6").(bool))

	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}

	if v, ok := d.GetOk("user_cidrs"); ok && v != nil {
		request.UserCidr = convertListToCommaSeparate(v.([]interface{}))
	}

	if v, ok := connectivity.GetResourceDataOk(d, "vpc_name", "name"); ok {
		request.VpcName = v.(string)
	}

	return request
}
