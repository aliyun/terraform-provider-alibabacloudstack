package alibabacloudstack

import (
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackSwitchCreate,
		Read:   resourceAlibabacloudStackSwitchRead,
		Update: resourceAlibabacloudStackSwitchUpdate,
		Delete: resourceAlibabacloudStackSwitchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cidr_block": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateSwitchCIDRNetworkAddress,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlibabacloudStackSwitchCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateCreateVSwitchRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.VpcId = Trim(d.Get("vpc_id").(string))
	request.ZoneId = d.Get("availability_zone").(string)
	request.CidrBlock = Trim(d.Get("cidr_block").(string))

	if v, ok := d.GetOk("name"); ok && v != "" {
		request.VSwitchName = v.(string)
	}

	if v, ok := d.GetOk("description"); ok && v != "" {
		request.Description = v.(string)
	}
	request.ClientToken = buildClientToken(request.GetActionName())

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateVSwitch(&args)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "UnknownError", "InvalidStatus.RouteEntry",
				"InvalidCidrBlock.Overlapped", Throttling, "OperationFailed.IdempotentTokenProcessing"}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.CreateVSwitchResponse)
		log.Printf("Vswitch Request %s", response)
		d.SetId(response.VSwitchId)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_vswitch", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, vpcService.VSwitchStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlibabacloudStackSwitchUpdate(d, meta)
}

func resourceAlibabacloudStackSwitchRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	vswitch, err := vpcService.DescribeVSwitch(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("availability_zone", vswitch.ZoneId)
	d.Set("vpc_id", vswitch.VpcId)
	d.Set("cidr_block", vswitch.CidrBlock)
	d.Set("name", vswitch.VSwitchName)
	listTagResourcesObject, err := vpcService.ListTagResources(d.Id(), "VSWITCH")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", tagsToMap(listTagResourcesObject))
	d.Set("description", vswitch.Description)
	return nil
}

func resourceAlibabacloudStackSwitchUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	if d.HasChange("tags") {
		if err := vpcService.SetResourceTags(d, "VSWITCH"); err != nil {
			return WrapError(err)
		}
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlibabacloudStackSwitchRead(d, meta)
	}
	update := false
	request := vpc.CreateModifyVSwitchAttributeRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}

	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.VSwitchId = d.Id()

	if d.HasChange("name") {
		request.VSwitchName = d.Get("name").(string)
		update = true
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}
	if update {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyVSwitchAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlibabacloudStackSwitchRead(d, meta)
}

func resourceAlibabacloudStackSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	request := vpc.CreateDeleteVSwitchRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}

	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.VSwitchId = d.Id()
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVSwitch(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidRegionId.NotFound"}) {
				return resource.NonRetryableError(err)
			}
			if IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
				return nil
			}

			return resource.RetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Second, vpcService.VSwitchStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
