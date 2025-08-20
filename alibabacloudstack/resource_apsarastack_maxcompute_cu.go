package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackMaxcomputeCu() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
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
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackMaxcomputeCuCreate, resourceAlibabacloudStackMaxcomputeCuRead, resourceAlibabacloudStackMaxcomputeCuUpdate, resourceAlibabacloudStackMaxcomputeCuDelete)
	return resource
}

func resourceAlibabacloudStackMaxcomputeCuCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "CreateUpdateOdpsCuForAscm"
	response := make(map[string]interface{})
	request := client.NewCommonRequest("POST", "dataworks-private-cloud", "2019-01-17", action, "")
	mergeMaps(request.QueryParams, map[string]string{
		"Region":      client.RegionId,
		"Action":      "CreateUpdateOdpsCuForAscm",
		"AccessKeyId": client.AccessKey,
		"CuName":      d.Get("cu_name").(string),
		"CuNum":       fmt.Sprintf("%v", d.Get("cu_num").(int)),
		"Cluster":     d.Get("cluster_name").(string),
		"ClusterName": d.Get("cluster_name").(string),
		"Share":       "0",
	})

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(action, bresponse, request, request.QueryParams)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_cu", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_cu", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	data := response["Data"].(map[string]interface{})
	d.SetId(data["CuId"].(string))

	return nil
}

func resourceAlibabacloudStackMaxcomputeCuUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if !d.IsNewResource() && d.HasChange("cu_num") {
		action := "CreateUpdateOdpsCuForAscm"
		request := client.NewCommonRequest("POST", "dataworks-private-cloud", "2019-01-17", action, "")
		mergeMaps(request.QueryParams, map[string]string{
			"Region":      client.RegionId,
			"Action":      "CreateUpdateOdpsCuForAscm",
			"AccessKeyId": client.AccessKey,
			"CuName":      d.Get("cu_name").(string),
			"CuNum":       fmt.Sprintf("%v", d.Get("cu_num").(int)),
			"Cluster":     d.Get("cluster_name").(string),
			"CuId":        d.Id(),
		})
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(action, bresponse, request, request.QueryParams)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_cu", action, errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}
	return nil
}

func resourceAlibabacloudStackMaxcomputeCuRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	maxcomputeService := MaxcomputeService{client}
	object, err := maxcomputeService.DescribeMaxcomputeCu(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_maxcompute_project maxcomputeService.DescribeMaxcomputeCu Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	var cu_num int
	switch v := object["max_cu"].(type) {
	case string:
		cu_num, err = strconv.Atoi(v)
		if err != nil {
			return errmsgs.WrapError(errmsgs.Error("illegal max_cu value"))
		}
	case json.Number:
		var floatVal float64
		floatVal, err = v.Float64()
		if err != nil {
			return errmsgs.WrapError(errmsgs.Error("illegal max_cu value"))
		}
		cu_num = int(floatVal)
	case int:
		cu_num = v
	case float64:
		cu_num = int(v)
	default:
		return errmsgs.WrapError(errmsgs.Error("illegal max_cu value type"))
	}
	d.Set("cu_num", cu_num)
	d.Set("cu_name", object["quota_name"].(string))
	d.Set("cluster_name", object["cluster"].(string))
	return nil
}

func resourceAlibabacloudStackMaxcomputeCuDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteOdpsCuForAscm"
	request := map[string]interface{}{
		"Region":      client.RegionId,
		"Action":      "CreateUpdateOdpsCuForAscm",
		"AccessKeyId": client.AccessKey,
		"CuName":      d.Get("cu_name").(string),
		"CuNum":       d.Get("cu_num").(int),
		"Cluster":     d.Get("cluster_name").(string),
		"CuId":        d.Id(),
	}

	_, err := client.DoTeaRequest("POST", "dataworks-private-cloud", "2019-01-17", action, "", nil, request, nil)

	if err != nil {
		return err
	}
	return nil
}
