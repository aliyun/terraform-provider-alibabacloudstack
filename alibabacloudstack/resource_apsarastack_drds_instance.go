package alibabacloudstack

import (
	"encoding/json"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDrdsInstance() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
				ForceNew:     true,
				Default:      PostPaid,
			},
			"specification": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 129),
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_series": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The `instance_series` property is no longer a required field. Selecting the `specification` will automatically retrieve the corresponding value.",
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackDRDSInstanceCreate, resourceAlibabacloudStackDRDSInstanceRead, resourceAlibabacloudStackDRDSInstanceUpdate, resourceAlibabacloudStackDRDSInstanceDelete)
	return resource
}

func resourceAlibabacloudStackDRDSInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	drdsService := DrdsService{client}

	action := "CreateDrdsInstance"

	var series string
	if v, err := drdsService.DescribeInstanceSpecification(d.Get("specification").(string)); err != nil {
		return err
	} else {
		series = v["seriesId"].(string)
	}

	reqQuery := map[string]interface{}{
		"Type":           "PRIVATE",
		"InstanceSeries": series,
		"Specification":  d.Get("specification").(string),
		"Description":    d.Get("description").(string),
		"Quantity":       "1",
		"ZoneId":         d.Get("zone_id").(string),
		"PayType":        d.Get("instance_charge_type").(string),
	}

	if v, ok := d.GetOk("vswitch_id"); ok && v.(string) != "" {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitch(v.(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		reqQuery["VpcId"] = vsw.VpcId
		reqQuery["VswitchId"] = v.(string)
		reqQuery["InstanceNetworkType"] = "VPC"
	} else {
		reqQuery["InstanceNetworkType"] = "CLASSIC"
	}

	if reqQuery["PayType"] == string(PostPaid) {
		reqQuery["PayType"] = "drdsPost"
	}
	if reqQuery["PayType"] == string(PrePaid) {
		reqQuery["PayType"] = "drdsPre"
	}
	
	if v, ok := d.GetOk("master_instance_id"); ok {
		// Only work for Readonly Instance
		reqQuery["MasterInstId"] = v
	}

	response, err := client.DoTeaRequest("POST", "Drds", "2019-01-23", action, "", nil, reqQuery, nil)
	if err != nil {
		return err
	}

	idList := response["Data"].(map[string]interface{})["DrdsInstanceIdList"].(map[string]interface{})["DrdsInstanceId"].([]interface{})
	if len(idList) != 1 {
		return errmsgs.WrapError(errmsgs.Error("failed to get DRDS instance id and response. DrdsInstanceIdList is %#v", idList))
	}
	d.SetId(idList[0].(string))

	stateConf := BuildStateConf([]string{"DO_CREATE"}, []string{"RUN"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, drdsService.DrdsInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func resourceAlibabacloudStackDRDSInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	drdsService := DrdsService{client}
	
	if d.IsNewResource() {
		return nil
	}
	
	if d.HasChange("description") {
		request := drds.CreateModifyDrdsInstanceDescriptionRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DrdsInstanceId = d.Id()
		request.Description = d.Get("description").(string)

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

	if d.HasChange("specification") {
		specification := d.Get("specification").(string)
		var series, chargeType string
		if v, err := drdsService.DescribeInstanceSpecification(d.Get("specification").(string)); err != nil {
			return err
		} else {
			series = v["seriesId"].(string)
		}
		if d.Get("instance_charge_type").(string) == string(PostPaid) {
			chargeType = "POSTPAY"
		} else {
			chargeType = "PREPAY"
		}
		
		reqQuery := map[string]interface{}{
			"commodityCode":"drdsPost",
			"data": map[string]interface{}{
				"drds_instance_type":"private",
				"orderType":"UPGRADE",
				"drds_region":client.RegionId,
				"drds_zone":d.Get("zone_id").(string),
				"instId":d.Id(),
				"drds_instance_series":series,
				"drds_instance_spec":specification,
				"chargeType":chargeType,
			},
		}
		orders, err := json.Marshal(reqQuery)
		if err != nil {
			return err
		}
		reqQuery = map[string]interface{}{
			"DrdsInstanceId":d.Id(),
			"Orders": string(orders),
		}
		if _, err := client.DoTeaRequest("POST", "Drds", "2019-01-23", "UpgradeDrdsInstance", "", nil, reqQuery, nil); err != nil {
			return err
		}
		stateConf := BuildStateConf([]string{"CHANGE_GRADE"}, []string{"RUN"}, d.Timeout(schema.TimeoutUpdate), 3*time.Second, drdsService.DrdsInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
			
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
	d.Set("specification", data.InstanceSpec)
	for _, vip := range data.Vips.Vip {
		if vip.VswitchId!= ""{
			d.Set("vswitch_id", vip.VswitchId)
			break
		}
	}
	
	if data.MasterInstanceId != "" {
		d.Set("master_instance_id", data.MasterInstanceId) 
	}
	
	if data.CommodityCode == "drdsPost" {
		d.Set("instance_charge_type", string(PostPaid))
	} else if data.CommodityCode == "drdsPre" {
		d.Set("instance_charge_type", string(PrePaid))
	}

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
		if errmsgs.IsExpectedErrors(err, "InvalidDrdsInstanceId.NotFound") {
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
