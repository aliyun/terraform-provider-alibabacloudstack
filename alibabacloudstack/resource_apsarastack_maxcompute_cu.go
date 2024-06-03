package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
	product := "ascm"
	response := make(map[string]interface{})
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = product
	request.Domain = client.Domain
	request.Version = "2019-05-10"
	request.ApiName = action
	request.RegionId = client.RegionId
	request.QueryParams = map[string]string{
		"CuName":          d.Get("cu_name").(string),
		"CuNum":           fmt.Sprintf("%v", d.Get("cu_num").(int)),
		"ClusterName":     d.Get("cluster_name").(string),
		"Department":      fmt.Sprintf("%v", client.Department),
		"OrganizationId":  fmt.Sprintf("%v", client.Department),
		"RegionId":        client.RegionId,
		"ResourceGroupId": fmt.Sprintf("%v", client.ResourceGroup),
		"RegionName":      client.RegionId,
		"Share":           "0",
		"Product":         product,
	}
	request.Headers = map[string]string{
		"RegionId":           client.RegionId,
		"x-acs-content-type": "application/json",
		"Content-Type":       "application/json",
	}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	bresponse, _ := raw.(*responses.CommonResponse)
	addDebug(action, raw, request)
	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_maxcompute_cu", "CreateUpdateOdpsCu", AlibabacloudStackSdkGoERROR)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_maxcompute_cu", action, AlibabacloudStackSdkGoERROR)
	}
	if fmt.Sprintf(`%v`, response["code"]) != "200" {
		return WrapError(Error("CreateUpdateOdpsCu failed for " + response["asapiErrorMessage"].(string)))
	}

	d.Set("cu_name", d.Get("cu_name").(string))

	return resourceAlibabacloudStackMaxcomputeCuRead(d, meta)
}

func resourceAlibabacloudStackMaxcomputeCuRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	maxcomputeService := MaxcomputeService{client}
	object, err := maxcomputeService.DescribeMaxcomputeCu(d.Get("cu_name").(string))
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_maxcompute_project maxcomputeService.DescribeMaxcomputeCu Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
			return WrapError(Error("illegal max_cu value"))
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
	var response map[string]interface{}
	conn, err := client.NewOdpsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"CuId":        d.Id(),
		"CuName":      d.Get("cu_name"),
		"ClusterName": d.Get("cluster_name"),
		"Product":     "ascm",
		"RegionId":    client.RegionId,
		"RegionName":  client.RegionId,
	}

	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequestWithOrg(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"500"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}
	if IsExpectedErrorCodes(fmt.Sprintf("%v", response["code"]), []string{"102", "403"}) {
		return nil
	}

	return nil
}
