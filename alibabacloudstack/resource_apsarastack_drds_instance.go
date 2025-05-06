package alibabacloudstack

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDRDSInstance() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 129),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"specification": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
				ForceNew: true,
				Default:  PostPaid,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_series": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"drds.sn2.4c16g", "drds.sn2.8c32g", "drds.sn2.16c64g", "drds.sn1.32c64g"}, false),
				ForceNew:     true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackDRDSInstanceCreate, resourceAlibabacloudStackDRDSInstanceRead, resourceAlibabacloudStackDRDSInstanceUpdate, resourceAlibabacloudStackDRDSInstanceDelete)
	return resource
}

func resourceAlibabacloudStackDRDSInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	drdsService := DrdsService{client}

	request := drds.CreateCreateDrdsInstanceRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.Description = d.Get("description").(string)
	request.Type = "PRIVATE"
	request.ZoneId = d.Get("zone_id").(string)
	request.Specification = d.Get("specification").(string)
	request.PayType = d.Get("instance_charge_type").(string)
	request.VswitchId = d.Get("vswitch_id").(string)
	request.InstanceSeries = d.Get("instance_series").(string)
	request.Quantity = "1"

	if request.VswitchId != "" {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitch(request.VswitchId)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.VpcId = vsw.VpcId
	}
	request.ClientToken = buildClientToken(request.GetActionName())

	if request.PayType == string(PostPaid) {
		request.PayType = "drdsPost"
	}
	if request.PayType == string(PrePaid) {
		request.PayType = "drdsPre"
	}

	raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.CreateDrdsInstance(request)
	})
	bresponse, ok := raw.(*drds.CreateDrdsInstanceResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_drds_instance", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	idList := bresponse.Data.DrdsInstanceIdList.DrdsInstanceIdList
	if len(idList) != 1 {
		return errmsgs.WrapError(errmsgs.Error("failed to get DRDS instance id and response. DrdsInstanceIdList is %#v", idList))
	}
	d.SetId(idList[0])

	stateConf := BuildStateConf([]string{"DO_CREATE"}, []string{"RUN"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, drdsService.DrdsInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func resourceAlibabacloudStackDRDSInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	drdsService := DrdsService{client}

	configItem := make(map[string]string)
	if d.HasChange("description") {
		request := drds.CreateModifyDrdsInstanceDescriptionRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DrdsInstanceId = d.Id()
		request.Description = d.Get("description").(string)
		configItem["description"] = request.Description

		raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
			return drdsClient.ModifyDrdsInstanceDescription(request)
		})
		bresponse, ok := raw.(*drds.ModifyDrdsInstanceDescriptionResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if err := drdsService.WaitDrdsInstanceConfigEffect(
		d.Id(), configItem, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return errmsgs.WrapError(err)
	}
	stateConf := BuildStateConf([]string{}, []string{"RUN"}, d.Timeout(schema.TimeoutUpdate), 3*time.Second, drdsService.DrdsInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func resourceAlibabacloudStackDRDSInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	drdsService := DrdsService{client}

	object, err := drdsService.DescribeDrdsInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	data := object.Data
	d.Set("zone_id", data.ZoneId)
	d.Set("description", data.Description)

	return nil
}

func resourceAlibabacloudStackDRDSInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	drdsService := DrdsService{client}

	request := drds.CreateRemoveDrdsInstanceRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DrdsInstanceId = d.Id()

	raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.RemoveDrdsInstance(request)
	})
	bresponse, ok := raw.(*drds.RemoveDrdsInstanceResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDrdsInstanceId.NotFound"}) {
			return nil
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if !bresponse.Success {
		return errmsgs.WrapError(errmsgs.Error("failed to delete instance timeout "+"and got an error: %#v", err))
	}

	stateConf := BuildStateConf([]string{"RUN", "DO_CREATE", "EXCEPTION", "EXPIRE", "DO_RELEASE", "RELEASE", "UPGRADE", "DOWNGRADE", "VersionUpgrade", "VersionRollback", "RESTART"}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Second, drdsService.DrdsInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err = stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}