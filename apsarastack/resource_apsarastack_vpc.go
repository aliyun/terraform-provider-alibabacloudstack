package apsarastack

import (
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackVpcCreate,
		Read:   resourceApsaraStackVpcRead,
		Update: resourceApsaraStackVpcUpdate,
		Delete: resourceApsaraStackVpcDelete,
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

func resourceApsaraStackVpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	var response *vpc.CreateVpcResponse
	request := buildApsaraStackVpcArgs(d, meta)
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = string(client.Region)
	//request.Headers = map[string]string{"RegionId": client.RegionId}
	//request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateVpc(&args)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "UnknownError", Throttling}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*vpc.CreateVpcResponse)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_vpc", request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	d.SetId(response.VpcId)

	stateConf := BuildStateConf([]string{"Pending"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, vpcService.VpcStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceApsaraStackVpcUpdate(d, meta)
}

func resourceApsaraStackVpcRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeVpc(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource apsarastack_vpc_vpc vpcService.DescribeVpc Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cidr_block", object.CidrBlock)
	d.Set("name", object.VpcName)
	d.Set("description", object.Description)
	d.Set("router_id", object.VRouterId)
	d.Set("ipv6_cidr_block", object.Ipv6CidrBlock)
	d.Set("secondary_cidr_blocks", object.SecondaryCidrBlocks.SecondaryCidrBlock)
	d.Set("status", object.Status)
	if tag := object.Tags.Tag; tag != nil {
		d.Set("tags", vpcService.tagToMap(tag))
	}
	d.Set("user_cidrs", object.UserCidrs.UserCidr)
	d.Set("vpc_name", object.VpcName)
	d.Set("name", object.VpcName)

	request := vpc.CreateDescribeRouteTablesRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.VRouterId = object.VRouterId
	request.ResourceGroupId = object.ResourceGroupId
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
				if IsExpectedErrors(err, []string{Throttling}) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(err)
				} else {
					return resource.NonRetryableError(err)
				}
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			response, _ := raw.(*vpc.DescribeRouteTablesResponse)
			routeTabls = append(routeTabls, response.RouteTables.RouteTable...)
			total = len(response.RouteTables.RouteTable)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}

		if total < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
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

func resourceApsaraStackVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	if err := vpcService.setInstanceSecondaryCidrBlocks(d); err != nil {
		return WrapError(err)
	}
	if d.HasChange("tags") {
		if err := vpcService.SetResourceTags(d, "vpc"); err != nil {
			return WrapError(err)
		}
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceApsaraStackVpcRead(d, meta)
	}

	groupRequestUpdate := false
	groupRequest := vpc.CreateMoveResourceGroupRequest()
	groupRequest.RegionId = client.RegionId
	groupRequest.Headers = map[string]string{"RegionId": client.RegionId}
	groupRequest.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	if strings.ToLower(client.Config.Protocol) == "https" {
		groupRequest.Scheme = "https"
	} else {
		groupRequest.Scheme = "http"
	}
	groupRequest.ResourceId = d.Id()

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		groupRequestUpdate = true
	}
	groupRequest.NewResourceGroupId = d.Get("resource_group_id").(string)
	groupRequest.ResourceType = "vpc"
	if groupRequestUpdate {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.MoveResourceGroup(groupRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), groupRequest.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(groupRequest.GetActionName(), raw, groupRequest.RpcRequest, groupRequest)
	}

	attributeUpdate := false
	request := vpc.CreateModifyVpcAttributeRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.VpcId = d.Id()

	if d.HasChange("name") {
		request.VpcName = d.Get("name").(string)
		attributeUpdate = true
	}

	if d.HasChange("vpc_name") {
		request.VpcName = d.Get("vpc_name").(string)
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

	if attributeUpdate {
		if _, ok := d.GetOkExists("enable_ipv6"); ok {
			request.EnableIPv6 = requests.NewBoolean(d.Get("enable_ipv6").(bool))
		}
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyVpcAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceApsaraStackVpcRead(d, meta)
}

func resourceApsaraStackVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}
	request := vpc.CreateDeleteVpcRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.VpcId = d.Id()
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVpc(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidVpcID.NotFound", "Forbidden.VpcNotFound"}) {
				return nil
			}
			return resource.RetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{"Pending"}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Second, vpcService.VpcStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func buildApsaraStackVpcArgs(d *schema.ResourceData, meta interface{}) *vpc.CreateVpcRequest {
	client := meta.(*connectivity.ApsaraStackClient)
	request := vpc.CreateCreateVpcRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.CidrBlock = d.Get("cidr_block").(string)

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request.DryRun = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOkExists("enable_ipv6"); ok {
		request.EnableIpv6 = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}

	if v, ok := d.GetOk("user_cidrs"); ok && v != nil {
		request.UserCidr = convertListToCommaSeparate(v.([]interface{}))
	}

	if v, ok := d.GetOk("vpc_name"); ok {
		request.VpcName = v.(string)
	} else if v, ok := d.GetOk("name"); ok {
		request.VpcName = v.(string)
	}

	//request.ClientToken = buildClientToken(request.GetActionName())

	return request
}
