package alibabacloudstack

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackOnsInstance() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(3, 64),
			},

			"tps_receive_max": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"tps_send_max": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"topic_capacity": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"independent_naming": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},

			// Computed Values
			"instance_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"instance_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackOnsInstanceCreate, resourceAlibabacloudStackOnsInstanceRead, resourceAlibabacloudStackOnsInstanceUpdate, resourceAlibabacloudStackOnsInstanceDelete)
	return resource
}

func resourceAlibabacloudStackOnsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	maxrtps := d.Get("tps_receive_max").(int)
	maxstps := d.Get("tps_send_max").(int)
	topiccapacity := d.Get("topic_capacity").(int)
	independentname := d.Get("independent_naming").(string)
	ins_resp := OnsInstance{}

	cluster := d.Get("cluster").(string)
	name := d.Get("name").(string)
	request := client.NewCommonRequest("POST", "Ons-inner", "2018-02-05", "ConsoleInstanceCreate", "")
	mergeMaps(request.QueryParams, map[string]string{
		"ProductName":       "Ons-inner",
		"OnsRegionId":       client.RegionId,
		"InstanceName":      name,
		"MaxReceiveTps":     fmt.Sprint(maxrtps),
		"MaxSendTps":        fmt.Sprint(maxstps),
		"TopicCapacity":     fmt.Sprint(topiccapacity),
		"Cluster":           cluster,
		"IndependentNaming": independentname,
	})

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	log.Printf(" response of raw ConsoleInstanceCreate : %s", bresponse)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ons_instance", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &ins_resp)
	if ins_resp.Success != true {
		return errmsgs.WrapErrorf(errors.New(ins_resp.Message), errmsgs.DefaultErrorMsg, "alibabacloudstack_ons_instance", "ConsoleInstanceCreate", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.SetId(ins_resp.Data.InstanceID)

	return nil
}

func resourceAlibabacloudStackOnsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	onsService := OnsService{client}

	response, err := onsService.DescribeOnsInstance(d.Id())

	if err != nil {
		// Handle exceptions
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("name", response.Data.InstanceName)
	d.Set("instance_type", response.Data.InstanceType)
	d.Set("instance_status", response.Data.InstanceStatus)
	d.Set("remark", response.Data.Remark)
	d.Set("tps_receive_max", response.Data.TpsReceiveMax)
	d.Set("tps_send_max", response.Data.TpsSendMax)
	d.Set("independent_naming", fmt.Sprintf("%t", response.Data.IndependentNaming))
	d.Set("cluster", response.Data.Cluster)
	d.Set("topic_capacity", response.Data.TopicCapacity)
	d.Set("create_time", time.Unix(response.Data.CreateTime/1000, 0).Format("2006-01-02 03:04:05"))

	return nil
}

func resourceAlibabacloudStackOnsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	onsService := OnsService{client}
	independentname := d.Get("independent_naming").(string)
	cluster := d.Get("cluster").(string)
	_, err := onsService.DescribeOnsInstance(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsInstanceExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if d.HasChanges("name", "tps_receive_max", "tps_send_max", "topic_capacity", "remark") {
		name := d.Get("name").(string)
		remark := d.Get("remark").(string)
		topiccap := strconv.Itoa(d.Get("topic_capacity").(int))
		Maxrtps := strconv.Itoa(d.Get("tps_receive_max").(int))
		Maxstps := strconv.Itoa(d.Get("tps_send_max").(int))
		request := client.NewCommonRequest("POST", "Ons-inner", "2018-02-05", "ConsoleInstanceUpdate", "")
		mergeMaps(request.QueryParams, map[string]string{
			"Remark":            remark,
			"InstanceName":      name,
			"OnsRegionId":       client.RegionId,
			"PreventCache":      "",
			"MaxReceiveTps":     Maxrtps,
			"MaxSendTps":        Maxstps,
			"Cluster":           cluster,
			"IndependentNaming": independentname,
			"InstanceId":        d.Id(),
			"TopicCapacity":     topiccap,
		})

		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		log.Printf(" response of raw DescribeProjectMeta : %s", bresponse)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ons_instance", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		log.Printf("total QueryParams and topic %v %v", request.GetQueryParams(), name)

	}

	return nil
}

func resourceAlibabacloudStackOnsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	onsService := OnsService{client}
	var requestInfo *ons.Client
	check, err := onsService.DescribeOnsInstance(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsInstanceExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsInstanceExist", check, requestInfo, map[string]string{"InstanceId": d.Id()})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		request := client.NewCommonRequest("POST", "Ons-inner", "2018-02-05", "ConsoleInstanceDelete", "")
		request.QueryParams["OnsRegionId"] = client.RegionId
		request.QueryParams["InstanceId"] = d.Id()
		request.QueryParams["PreventCache"] = ""

		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		log.Printf(" response of raw ConsoleInstanceDelete : %s", bresponse)
		if err != nil {
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ons_instance", "ConsoleInstanceDelete", errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		check, err = onsService.DescribeOnsInstance(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	return nil
}
