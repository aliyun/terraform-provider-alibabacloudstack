package alibabacloudstack

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackOnsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackOnsInstanceCreate,
		Read:   resourceAlibabacloudStackOnsInstanceRead,
		Update: resourceAlibabacloudStackOnsInstanceUpdate,
		Delete: resourceAlibabacloudStackOnsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
				ValidateFunc: validation.StringLenBetween(0, 128),
			},

			// Computed Values
			"instance_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
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
}

func resourceAlibabacloudStackOnsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *ecs.Client
	maxrtps := d.Get("tps_receive_max").(int)
	maxstps := d.Get("tps_send_max").(int)
	topiccapacity := d.Get("topic_capacity").(int)
	independentname := d.Get("independent_naming").(string)
	ins_resp := OnsInstance{}

	cluster := d.Get("cluster").(string)
	name := d.Get("name").(string)
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":          client.RegionId,
		"AccessKeySecret":   client.SecretKey,
		"Department":        client.Department,
		"ResourceGroup":     client.ResourceGroup,
		"Product":           "Ons-inner",
		"Action":            "ConsoleInstanceCreate",
		"Version":           "2018-02-05",
		"ProductName":       "Ons-inner",
		"OnsRegionId":       client.RegionId,
		"InstanceName":      name,
		"MaxReceiveTps":     fmt.Sprint(maxrtps),
		"MaxSendTps":        fmt.Sprint(maxstps),
		"TopicCapacity":     fmt.Sprint(topiccapacity),
		"Cluster":           cluster,
		"IndependentNaming": independentname,
	}
	request.Method = "POST"
	request.Product = "Ons-inner"
	request.Version = "2018-02-05"
	request.ServiceCode = "Ons-inner"
	request.Domain = client.Domain
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ConsoleInstanceCreate"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ons_instance", "ConsoleInstanceCreate", raw)
	}
	addDebug("ConsoleInstanceCreate", raw, requestInfo, request)
	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &ins_resp)
	if ins_resp.Success != true {
		return WrapErrorf(errors.New(ins_resp.Message), DefaultErrorMsg, "alibabacloudstack_ons_instance", "ConsoleInstanceCreate", AlibabacloudStackSdkGoERROR)
	}

	if err != nil {
		return WrapError(err)
	}
	d.SetId(ins_resp.Data.InstanceID)

	return resourceAlibabacloudStackOnsInstanceUpdate(d, meta)
}

func resourceAlibabacloudStackOnsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	onsService := OnsService{client}

	response, err := onsService.DescribeOnsInstance(d.Id())

	if err != nil {
		// Handle exceptions
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", response.Data.InstanceName)
	d.Set("instance_type", response.Data.InstanceType)
	d.Set("instance_status", response.Data.InstanceStatus)
	d.Set("create_time", time.Unix(response.Data.CreateTime/1000, 0).Format("2006-01-02 03:04:05"))

	return nil
}

func resourceAlibabacloudStackOnsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	onsService := OnsService{client}
	independentname := d.Get("independent_naming").(string)
	cluster := d.Get("cluster").(string)
	attributeUpdate := false
	check, err := onsService.DescribeOnsInstance(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsInstanceExist", AlibabacloudStackSdkGoERROR)
	}
	var name, remark string

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		}
		check.Data.InstanceName = name
		attributeUpdate = true
	} else {
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		}
		check.Data.InstanceName = name
	}
	var maxrtps, maxstps, topic int

	if d.HasChange("tps_receive_max") {
		if v, ok := d.GetOk("tps_receive_max"); ok {
			maxrtps = v.(int)
		}
		check.Data.TpsReceiveMax = maxrtps
		attributeUpdate = true
	} else {
		if v, ok := d.GetOk("tps_receive_max"); ok {
			maxrtps = v.(int)
		}
		check.Data.TpsReceiveMax = maxrtps
	}
	if d.HasChange("tps_send_max") {
		if v, ok := d.GetOk("tps_send_max"); ok {
			maxstps = v.(int)
		}
		check.Data.TpsMax = maxstps
		attributeUpdate = true
	} else {
		if v, ok := d.GetOk("tps_send_max"); ok {
			maxstps = v.(int)
		}
		check.Data.TpsMax = maxstps
	}
	if d.HasChange("topic_capacity") {
		if v, ok := d.GetOk("topic_capacity"); ok {
			topic = v.(int)
		}
		check.Data.TopicCapacity = topic
		attributeUpdate = true
	} else {
		if v, ok := d.GetOk("topic_capacity"); ok {
			topic = v.(int)
		}
		check.Data.TopicCapacity = topic
	}

	if d.HasChange("remark") {

		if v, ok := d.GetOk("remark"); ok {
			remark = v.(string)
		}
		check.Data.Remark = remark
		attributeUpdate = true
	} else {
		if v, ok := d.GetOk("remark"); ok {
			remark = v.(string)
		}
		check.Data.Remark = remark
	}
	topiccap := strconv.Itoa(topic)
	Maxrtps := strconv.Itoa(maxrtps)
	Maxstps := strconv.Itoa(maxstps)
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":          client.RegionId,
		"AccessKeySecret":   client.SecretKey,
		"Department":        client.Department,
		"ResourceGroup":     client.ResourceGroup,
		"Product":           "Ons-inner",
		"Action":            "ConsoleInstanceUpdate",
		"Version":           "2018-02-05",
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
	}
	request.Method = "POST"
	request.Product = "Ons-inner"
	request.Version = "2018-02-05"
	request.ServiceCode = "Ons-inner"
	request.Domain = client.Domain
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.SetHTTPSInsecure(true)
	request.ApiName = "ConsoleInstanceUpdate"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	check.Data.InstanceID = d.Id()

	if attributeUpdate {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw ConsoleInstanceUpdate : %s", raw)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ons_instance", "ConsoleInstanceCreate", raw)
		}
		addDebug(request.GetActionName(), raw, request)
		log.Printf("total QueryParams and topic %v %v", request.GetQueryParams(), topic)

	}

	return resourceAlibabacloudStackOnsInstanceRead(d, meta)
}

func resourceAlibabacloudStackOnsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	onsService := OnsService{client}
	var requestInfo *ecs.Client
	check, err := onsService.DescribeOnsInstance(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsInstanceExist", AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsInstanceExist", check, requestInfo, map[string]string{"InstanceId": d.Id()})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ons-inner",
			"Action":          "ConsoleInstanceDelete",
			"Version":         "2018-02-05",
			"ProductName":     "Ons-inner",
			"PreventCache":    "",
			"OnsRegionId":     client.RegionId,
			"InstanceId":      d.Id(),
		}

		request.Method = "POST"
		request.Product = "Ons-inner"
		request.Version = "2018-02-05"
		request.ServiceCode = "Ons-inner"
		request.Domain = client.Domain
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "ConsoleInstanceDelete"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		check, err = onsService.DescribeOnsInstance(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	return nil
}
