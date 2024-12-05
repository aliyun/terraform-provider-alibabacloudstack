package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackMaxcomputeCu() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackMaxcomputeCuCreate,
		Read:   resourceAlibabacloudStackMaxcomputeCuRead,
		Delete: resourceAlibabacloudStackMaxcomputeCuDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"cu_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 27),
			},
			"cu_num": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntAtLeast(1),
				Required:     true,
				ForceNew:     true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlibabacloudStackMaxcomputeCuCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "CreateUpdateOdpsCu"
	response := make(map[string]interface{})
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", action, "")
	mergeMaps(request.QueryParams, map[string]string{
		"CuName":          d.Get("cu_name").(string),
		"CuNum":           fmt.Sprintf("%v", d.Get("cu_num").(int)),
		"ClusterName":     d.Get("cluster_name").(string),
		"ResourceGroupId": fmt.Sprintf("%v", client.ResourceGroup),
		"RegionName":      client.RegionId,
		"Share":           "0",
	})
	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-Type"] = "application/json"

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	addDebug(action, raw, request)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_cu", action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	if bresponse.GetHttpStatus() != 200 {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_cu", action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_cu", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if fmt.Sprintf(`%v`, response["code"]) != "200" {
		return errmsgs.WrapError(errmsgs.Error("CreateUpdateOdpsCu failed for " + response["asapiErrorMessage"].(string)))
	}

	d.Set("cu_name", d.Get("cu_name").(string))

	return resourceAlibabacloudStackMaxcomputeCuRead(d, meta)
}

func resourceAlibabacloudStackMaxcomputeCuRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	maxcomputeService := MaxcomputeService{client}
	object, err := maxcomputeService.DescribeMaxcomputeCu(d.Get("cu_name").(string))
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_maxcompute_project maxcomputeService.DescribeMaxcomputeCu Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	var data map[string]interface{}
	datas := object["data"].([]interface{})
	if datas == nil || len(datas) < 1 {
		d.SetId(d.Get("id").(string))
		d.Set("cluster_name", d.Get("cluster_name").(string))
	}
	s := d.Get("cu_name").(string)
	for _, element := range datas {
		data = element.(map[string]interface{})
		if data["quota_name"].(string) != s {
			continue
		}
		d.SetId(data["id"].(string))
		max_cu, err := data["max_cu"].(json.Number).Float64()
		if err != nil {
			return errmsgs.WrapError(errmsgs.Error("illegal max_cu value"))
		}
		d.Set("cu_num", int64(max_cu))
		d.Set("cluster_name", data["cluster"].(string))
		break
	}
	return nil
}

func resourceAlibabacloudStackMaxcomputeCuDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteOdpsCu"
	request := make(map[string]interface{})
	request["CuId"] = d.Id()
	request["CuName"] = d.Get("cu_name")
	request["ClusterName"] = d.Get("cluster_name")

	response, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", action, "", nil, request)
	
	if err != nil {
		return err
	}
	if fmt.Sprintf("%v", response["code"]) == "102" || fmt.Sprintf("%v", response["code"]) == "403" {
		return nil
	}

	return nil
}
