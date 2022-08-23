package apsarastack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
	"time"
)

func resourceApsaraStackAscmQuota() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackAscmQuotaCreate,
		Read:   resourceApsaraStackAscmQuotaRead,
		Update: resourceApsaraStackAscmQuotaUpdate,
		Delete: resourceApsaraStackAscmQuotaDelete,
		Schema: map[string]*schema.Schema{
			"product_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{"ECS", "OSS", "VPC", "DRDS", "RDS", "SLB",
					"ODPS", "EIP", "GPDB", "R-KVSTORE", "NAS", "DDS"}, false),
			},
			"quota_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"organization", "resourceGroup"}, false),
			},
			"quota_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"total_cpu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"total_mem": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"total_gpu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"total_disk_cloud_ssd": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"total_disk_cloud_efficiency": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"total_amount": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"total_vpc": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"total_disk": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"total_vip_internal": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"total_vip_public": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"total_cu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"total_eip": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"quota_type_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceApsaraStackAscmQuotaCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var requestInfo *ecs.Client
	productName := strings.ToUpper(d.Get("product_name").(string))
	quotaType := d.Get("quota_type").(string)
	quotaTypeId := d.Get("quota_type_id").(string)
	if productName == "ECS" {
		totalCpu := d.Get("total_cpu").(int)
		totalMem := d.Get("total_mem").(int)
		totalGpu := d.Get("total_gpu").(int)
		totalDiskCloudSsd := d.Get("total_disk_cloud_ssd").(int)
		totalDiskCloudEfficiency := d.Get("total_disk_cloud_efficiency").(int)
		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "CreateQuota",
			"Version":         "2019-05-10",
			"regionName":      client.RegionId,
			"quotaType":       quotaType,
			"quotaTypeId":     quotaTypeId,
			"productName":     productName,
			"targetType":      "",
			"quotaBody": fmt.Sprintf("{\"%s\":\"%d\",\"%s\":\"%d\",\"%s\":\"%d\",\"%s\":\"%d\",\"%s\":\"%d\"}",
				"totalCpu", totalCpu,
				"totalMem", totalMem,
				"totalGpu", totalGpu,
				"totalDisk_cloud_ssd", totalDiskCloudSsd,
				"totalDisk_cloud_efficiency", totalDiskCloudEfficiency,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "CreateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw CreateQuota : %s", raw)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", raw)
		}
		addDebug("CreateQuota", raw, requestInfo, request)

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", ApsaraStackSdkGoERROR)
		}
		addDebug("CreateQuota", raw, requestInfo, bresponse.GetHttpContentString())

	} else if productName == "OSS" {
		totalAmount := d.Get("total_amount").(int)
		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"regionName":      client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "CreateQuota",
			"Version":         "2019-05-10",
			"quotaType":       quotaType,
			"quotaTypeId":     quotaTypeId,
			"productName":     productName,
			"targetType":      "",
			"quotaBody": fmt.Sprintf("{\"%s\":\"%d\"}",
				"totalAmount", totalAmount,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "CreateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", raw)
		}
		addDebug("CreateQuota", raw, requestInfo, request)

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", ApsaraStackSdkGoERROR)
		}
		addDebug("CreateQuota", raw, requestInfo, bresponse.GetHttpContentString())

	} else if productName == "EIP" {
		totalEIP := d.Get("total_eip").(int)
		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"regionName":      client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "CreateQuota",
			"Version":         "2019-05-10",
			"quotaType":       quotaType,
			"quotaTypeId":     quotaTypeId,
			"productName":     productName,
			"targetType":      "",
			"quotaBody": fmt.Sprintf("{\"%s\":\"%d\"}",
				"totalEIP", totalEIP,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "CreateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", raw)
		}
		addDebug("CreateQuota", raw, requestInfo, request)

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", ApsaraStackSdkGoERROR)
		}
		addDebug("CreateQuota", raw, requestInfo, bresponse.GetHttpContentString())

	} else if productName == "SLB" {
		totalVipPublic := d.Get("total_vip_public").(int)
		totalVipInternal := d.Get("total_vip_internal").(int)
		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"regionName":      client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "CreateQuota",
			"Version":         "2019-05-10",
			"quotaType":       quotaType,
			"quotaTypeId":     quotaTypeId,
			"productName":     productName,
			"targetType":      "",
			"quotaBody": fmt.Sprintf("{\"%s\":\"%d\",\"%s\":\"%d\"}",
				"totalVipPublic", totalVipPublic,
				"totalVipInternal", totalVipInternal,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "CreateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", raw)
		}
		addDebug("CreateQuota", raw, requestInfo, request)

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", ApsaraStackSdkGoERROR)
		}
		addDebug("CreateQuota", raw, requestInfo, bresponse.GetHttpContentString())

	} else if productName == "ODPS" {
		totalCu := d.Get("total_cu").(int)
		totalDisk := d.Get("total_disk").(int)
		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"regionName":      client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "CreateQuota",
			"Version":         "2019-05-10",
			"quotaType":       quotaType,
			"quotaTypeId":     quotaTypeId,
			"productName":     productName,
			"targetType":      "",
			"quotaBody": fmt.Sprintf("{\"%s\":\"%d\",\"%s\":\"%d\"}",
				"totalCu", totalCu,
				"totalDisk", totalDisk,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "CreateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", raw)
		}
		addDebug("CreateQuota", raw, requestInfo, request)

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", ApsaraStackSdkGoERROR)
		}
		addDebug("CreateQuota", raw, requestInfo, bresponse.GetHttpContentString())

	} else if productName == "DDS" {
		totalCpu := d.Get("total_cpu").(int)
		totalMem := d.Get("total_mem").(int)
		totalDisk := d.Get("total_disk").(int)
		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"regionName":      client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "CreateQuota",
			"Version":         "2019-05-10",
			"quotaType":       quotaType,
			"quotaTypeId":     quotaTypeId,
			"productName":     productName,
			"targetType":      "",
			"quotaBody": fmt.Sprintf("{\"%s\":\"%d\",\"%s\":\"%d\",\"%s\":\"%d\"}",
				"totalCpu", totalCpu,
				"totalMem", totalMem,
				"totalDisk", totalDisk,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "CreateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", raw)
		}
		addDebug("CreateQuota", raw, requestInfo, request)

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", ApsaraStackSdkGoERROR)
		}
		addDebug("CreateQuota", raw, requestInfo, bresponse.GetHttpContentString())

	} else if productName == "RDS" {
		targetType := d.Get("target_type").(string)
		totalCpu := d.Get("total_cpu").(int)
		totalMem := d.Get("total_mem").(int)
		totalDisk := d.Get("total_disk").(int)
		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"regionName":      client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "CreateQuota",
			"Version":         "2019-05-10",
			"quotaType":       quotaType,
			"quotaTypeId":     quotaTypeId,
			"productName":     productName,
			"targetType":      targetType,
			"quotaBody": fmt.Sprintf("{\"%s\":\"%d\",\"%s\":\"%d\",\"%s\":\"%d\"}",
				"totalCpu", totalCpu,
				"totalMem", totalMem,
				"totalDisk", totalDisk,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "CreateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", raw)
		}
		addDebug("CreateQuota", raw, requestInfo, request)

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", ApsaraStackSdkGoERROR)
		}
		addDebug("CreateQuota", raw, requestInfo, bresponse.GetHttpContentString())

	} else if productName == "GPDB" {
		totalCpu := d.Get("total_cpu").(int)
		totalMem := d.Get("total_mem").(int)
		totalDisk := d.Get("total_disk").(int)
		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"regionName":      client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "CreateQuota",
			"Version":         "2019-05-10",
			"quotaType":       quotaType,
			"quotaTypeId":     quotaTypeId,
			"productName":     productName,
			"targetType":      "",
			"quotaBody": fmt.Sprintf("{\"%s\":\"%d\",\"%s\":\"%d\",\"%s\":\"%d\"}",
				"totalCpu", totalCpu,
				"totalMem", totalMem,
				"totalDisk", totalDisk,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "CreateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", raw)
		}
		addDebug("CreateQuota", raw, requestInfo, request)

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", ApsaraStackSdkGoERROR)
		}
		addDebug("CreateQuota", raw, requestInfo, bresponse.GetHttpContentString())

	} else if productName == "R-KVSTORE" {
		targetType := d.Get("target_type").(string)
		totalMem := d.Get("total_mem").(int)
		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"regionName":      client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "CreateQuota",
			"Version":         "2019-05-10",
			"quotaType":       quotaType,
			"quotaTypeId":     quotaTypeId,
			"productName":     productName,
			"targetType":      targetType,
			"quotaBody": fmt.Sprintf("{\"%s\":\"%d\"}",
				"totalMem", totalMem,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "CreateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", raw)
		}
		addDebug("CreateQuota", raw, requestInfo, request)

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", ApsaraStackSdkGoERROR)
		}
		addDebug("CreateQuota", raw, requestInfo, bresponse.GetHttpContentString())

	} else if productName == "VPC" {
		totalVPC := d.Get("total_vpc").(int)
		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"regionName":      client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "CreateQuota",
			"Version":         "2019-05-10",
			"quotaType":       quotaType,
			"quotaTypeId":     quotaTypeId,
			"productName":     productName,
			"targetType":      "",
			"quotaBody": fmt.Sprintf("{\"%s\":%d}",
				"totalVPC", totalVPC,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "CreateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", raw)
		}
		addDebug("CreateQuota", raw, requestInfo, request)

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "CreateQuota", ApsaraStackSdkGoERROR)
		}
		addDebug("CreateQuota", raw, requestInfo, bresponse.GetHttpContentString())
	} else {
		log.Print("Please Enter a valid Product Name.")
	}

	d.SetId(productName + COLON_SEPARATED + quotaType + COLON_SEPARATED + quotaTypeId)

	return resourceApsaraStackAscmQuotaUpdate(d, meta)

}

func resourceApsaraStackAscmQuotaUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	attributeUpdate := false
	check, err := ascmService.DescribeAscmQuota(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsQuotaExist", ApsaraStackSdkGoERROR)
	}
	if did[0] == "VPC" {
		var totalVPC int

		if d.HasChange("total_vpc") {
			if v, ok := d.GetOk("total_vpc"); ok {
				totalVPC = v.(int)
			}
			check.Data.TotalVPC = totalVPC
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_vpc"); ok {
				totalVPC = v.(int)
			}
			check.Data.TotalVPC = totalVPC
		}
		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"regionName":      client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "UpdateQuota",
			"Version":         "2019-05-10",
			"quotaType":       did[1],
			"quotaTypeId":     did[2],
			"productName":     did[0],
			"targetType":      "",
			"quotaBody": fmt.Sprintf("{\"%s\":\"%d\"}",
				"totalVPC", totalVPC,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "UpdateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		if attributeUpdate {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ProcessCommonRequest(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "UpdateQuota", raw)
			}
			addDebug(request.GetActionName(), raw, request)
		}
	}

	if did[0] == "RDS" {
		var totalDISK, totalCPU, totalMEM int

		if d.HasChange("total_mem") {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = totalMEM
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = totalMEM
		}
		if d.HasChange("total_disk") {
			if v, ok := d.GetOk("total_disk"); ok {
				totalDISK = v.(int)
			}
			check.Data.TotalDisk = totalDISK
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_disk"); ok {
				totalDISK = v.(int)
			}
			check.Data.TotalDisk = totalDISK
		}
		if d.HasChange("total_cpu") {
			if v, ok := d.GetOk("total_cpu"); ok {
				totalCPU = v.(int)
			}
			check.Data.TotalCPU = totalCPU
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_cpu"); ok {
				totalCPU = v.(int)
			}
			check.Data.TotalCPU = totalCPU
		}
		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"regionName":      client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "UpdateQuota",
			"Version":         "2019-05-10",
			"quotaType":       did[1],
			"quotaTypeId":     did[2],
			"productName":     did[0],
			"targetType":      "MySql",
			"quotaBody": fmt.Sprintf("{\"%s\":\"%d\",\"%s\":\"%d\",\"%s\":\"%d\"}",
				"totalCpu", totalCPU,
				"totalMem", totalMEM,
				"totalDisk", totalDISK,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "UpdateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		if attributeUpdate {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ProcessCommonRequest(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "UpdateQuota", raw)
			}
			addDebug(request.GetActionName(), raw, request)
		}
	}
	if did[0] == "EIP" {
		var totalEIP int

		if d.HasChange("total_eip") {
			if v, ok := d.GetOk("total_eip"); ok {
				totalEIP = v.(int)
			}
			check.Data.TotalMem = totalEIP
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_eip"); ok {
				totalEIP = v.(int)
			}
			check.Data.TotalMem = totalEIP
		}
		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"regionName":      client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "UpdateQuota",
			"Version":         "2019-05-10",
			"quotaType":       did[1],
			"quotaTypeId":     did[2],
			"productName":     did[0],
			"targetType":      "",
			"quotaBody": fmt.Sprintf("{\"%s\":\"%d\"}",
				"totalEIP", totalEIP,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "UpdateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		if attributeUpdate {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ProcessCommonRequest(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "UpdateQuota", raw)
			}
			addDebug(request.GetActionName(), raw, request)
		}
	}
	if did[0] == "ECS" {
		var totalGPU, totalCPU, totalMEM, totalDCE, totalDCS int

		if d.HasChange("total_cpu") {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = totalMEM
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = totalMEM
		}
		if d.HasChange("total_gpu") {
			if v, ok := d.GetOk("total_gpu"); ok {
				totalGPU = v.(int)
			}
			check.Data.TotalGpu = totalGPU
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_gpu"); ok {
				totalGPU = v.(int)
			}
			check.Data.TotalGpu = totalGPU
		}
		if d.HasChange("total_cpu") {
			if v, ok := d.GetOk("total_cpu"); ok {
				totalCPU = v.(int)
			}
			check.Data.TotalCPU = totalCPU
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_cpu"); ok {
				totalCPU = v.(int)
			}
			check.Data.TotalCPU = totalCPU
		}
		if d.HasChange("total_disk_cloud_ssd") {
			if v, ok := d.GetOk("total_disk_cloud_ssd"); ok {
				totalDCS = v.(int)
			}
			check.Data.TotalDiskCloudSsd = totalDCS
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_disk_cloud_ssd"); ok {
				totalDCS = v.(int)
			}
			check.Data.TotalDiskCloudSsd = totalDCS
		}
		if d.HasChange("total_disk_cloud_efficiency") {
			if v, ok := d.GetOk("total_disk_cloud_efficiency"); ok {
				totalDCE = v.(int)
			}
			check.Data.TotalDiskCloudEfficiency = totalDCE
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_disk_cloud_efficiency"); ok {
				totalDCE = v.(int)
			}
			check.Data.TotalDiskCloudEfficiency = totalDCE
		}
		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"regionName":      client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "UpdateQuota",
			"Version":         "2019-05-10",
			"quotaType":       did[1],
			"quotaTypeId":     did[2],
			"productName":     did[0],
			"targetType":      "",
			"quotaBody": fmt.Sprintf("{\"%s\":\"%d\",\"%s\":\"%d\",\"%s\":\"%d\",\"%s\":\"%d\",\"%s\":\"%d\"}",
				"totalCpu", totalCPU,
				"totalMem", totalMEM,
				"totalGpu", totalGPU,
				"totalDisk_cloud_ssd", totalDCS,
				"totalDisk_cloud_efficiency", totalDCE,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "UpdateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		if attributeUpdate {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ProcessCommonRequest(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "UpdateQuota", raw)
			}
			addDebug(request.GetActionName(), raw, request)
		}
	}
	if did[0] == "SLB" {
		var totalVP, totalVI int

		if d.HasChange("total_vip_internal") {
			if v, ok := d.GetOk("total_vip_internal"); ok {
				totalVI = v.(int)
			}
			check.Data.TotalVipInternal = totalVI
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_vip_internal"); ok {
				totalVI = v.(int)
			}
			check.Data.TotalVipInternal = totalVI
		}
		if d.HasChange("total_vip_public") {
			if v, ok := d.GetOk("total_vip_public"); ok {
				totalVP = v.(int)
			}
			check.Data.TotalVipPublic = totalVP
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_vip_public"); ok {
				totalVP = v.(int)
			}
			check.Data.TotalVipPublic = totalVP
		}

		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"regionName":      client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "UpdateQuota",
			"Version":         "2019-05-10",
			"quotaType":       did[1],
			"quotaTypeId":     did[2],
			"productName":     did[0],
			"targetType":      "",
			"quotaBody": fmt.Sprintf("{\"%s\":\"%d\",\"%s\":\"%d\"}",
				"totalVipPublic", totalVP,
				"totalVipInternal", totalVI,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "UpdateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		if attributeUpdate {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ProcessCommonRequest(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "UpdateQuota", raw)
			}
			addDebug(request.GetActionName(), raw, request)
		}
	}
	if did[0] == "OSS" {
		var totalAmount int

		if d.HasChange("total_amount") {
			if v, ok := d.GetOk("total_amount"); ok {
				totalAmount = v.(int)
			}
			check.Data.TotalAmount = totalAmount
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_amount"); ok {
				totalAmount = v.(int)
			}
			check.Data.TotalAmount = totalAmount
		}

		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"regionName":      client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "UpdateQuota",
			"Version":         "2019-05-10",
			"quotaType":       did[1],
			"quotaTypeId":     did[2],
			"productName":     did[0],
			"targetType":      "",
			"quotaBody": fmt.Sprintf("{\"%s\":\"%d\"}",
				"totalAmount", totalAmount,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "UpdateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		if attributeUpdate {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ProcessCommonRequest(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "UpdateQuota", raw)
			}
			addDebug(request.GetActionName(), raw, request)
		}
	}
	if did[0] == "DDS" {
		var totalDISK, totalCPU, totalMEM int

		if d.HasChange("total_mem") {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = totalMEM
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = totalMEM
		}
		if d.HasChange("total_disk") {
			if v, ok := d.GetOk("total_disk"); ok {
				totalDISK = v.(int)
			}
			check.Data.TotalDisk = totalDISK
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_disk"); ok {
				totalDISK = v.(int)
			}
			check.Data.TotalDisk = totalDISK
		}
		if d.HasChange("total_cpu") {
			if v, ok := d.GetOk("total_cpu"); ok {
				totalCPU = v.(int)
			}
			check.Data.TotalCPU = totalCPU
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_cpu"); ok {
				totalCPU = v.(int)
			}
			check.Data.TotalCPU = totalCPU
		}
		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"regionName":      client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "UpdateQuota",
			"Version":         "2019-05-10",
			"quotaType":       did[1],
			"quotaTypeId":     did[2],
			"productName":     did[0],
			"targetType":      "",
			"quotaBody": fmt.Sprintf("{\"%s\":\"%d\",\"%s\":\"%d\",\"%s\":\"%d\"}",
				"totalCpu", totalCPU,
				"totalMem", totalMEM,
				"totalDisk", totalDISK,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "UpdateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		if attributeUpdate {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ProcessCommonRequest(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "UpdateQuota", raw)
			}
			addDebug(request.GetActionName(), raw, request)
		}
	}
	if did[0] == "ODPS" {
		var totalDISK, totalCU int

		if d.HasChange("total_disk") {
			if v, ok := d.GetOk("total_disk"); ok {
				totalDISK = v.(int)
			}
			check.Data.TotalDisk = totalDISK
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_disk"); ok {
				totalDISK = v.(int)
			}
			check.Data.TotalDisk = totalDISK
		}
		if d.HasChange("total_cu") {
			if v, ok := d.GetOk("total_cu"); ok {
				totalCU = v.(int)
			}
			check.Data.TotalCU = totalCU
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_cu"); ok {
				totalCU = v.(int)
			}
			check.Data.TotalCU = totalCU
		}
		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"regionName":      client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "UpdateQuota",
			"Version":         "2019-05-10",
			"quotaType":       did[1],
			"quotaTypeId":     did[2],
			"productName":     did[0],
			"targetType":      "",
			"quotaBody": fmt.Sprintf("{\"%s\":\"%d\",\"%s\":\"%d\"}",
				"totalCu", totalCU,
				"totalDisk", totalDISK,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "UpdateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		if attributeUpdate {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ProcessCommonRequest(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "UpdateQuota", raw)
			}
			addDebug(request.GetActionName(), raw, request)
		}
	}
	if did[0] == "GPDB" {
		var totalDISK, totalCPU, totalMEM int

		if d.HasChange("total_mem") {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = totalMEM
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = totalMEM
		}
		if d.HasChange("total_disk") {
			if v, ok := d.GetOk("total_disk"); ok {
				totalDISK = v.(int)
			}
			check.Data.TotalDisk = totalDISK
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_disk"); ok {
				totalDISK = v.(int)
			}
			check.Data.TotalDisk = totalDISK
		}
		if d.HasChange("total_cpu") {
			if v, ok := d.GetOk("total_cpu"); ok {
				totalCPU = v.(int)
			}
			check.Data.TotalCPU = totalCPU
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_cpu"); ok {
				totalCPU = v.(int)
			}
			check.Data.TotalCPU = totalCPU
		}
		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"regionName":      client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "UpdateQuota",
			"Version":         "2019-05-10",
			"quotaType":       did[1],
			"quotaTypeId":     did[2],
			"productName":     did[0],
			"targetType":      "",
			"quotaBody": fmt.Sprintf("{\"%s\":\"%d\",\"%s\":\"%d\",\"%s\":\"%d\"}",
				"totalCpu", totalCPU,
				"totalMem", totalMEM,
				"totalDisk", totalDISK,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "UpdateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		if attributeUpdate {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ProcessCommonRequest(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "UpdateQuota", raw)
			}
			addDebug(request.GetActionName(), raw, request)
		}
	}
	if did[0] == "R-KVSTORE" {
		var totalMEM int

		if d.HasChange("total_mem") {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = totalMEM
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = totalMEM
		}

		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"regionName":      client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ascm",
			"Action":          "UpdateQuota",
			"Version":         "2019-05-10",
			"quotaType":       did[1],
			"quotaTypeId":     did[2],
			"productName":     did[0],
			"targetType":      "redis",
			"quotaBody": fmt.Sprintf("{\"%s\":\"%d\"}",
				"totalMem", totalMEM,
			),
		}
		request.Method = "POST"
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "UpdateQuota"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		if attributeUpdate {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ProcessCommonRequest(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_quota", "UpdateQuota", raw)
			}
			addDebug(request.GetActionName(), raw, request)
		}
	}

	return resourceApsaraStackAscmQuotaRead(d, meta)

}

func resourceApsaraStackAscmQuotaRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmQuota(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if object.Code != "200" {
		d.SetId("")
		return nil
	}

	if did[0] == "VPC" {
		d.Set("quota_id", object.Data.ID)
		d.Set("quota_type", did[1])
		d.Set("quota_type_id", did[2])
		return nil
	} else if did[0] == "OSS" {
		d.Set("quota_id", object.Data.ID)
		d.Set("quota_type", did[1])
		d.Set("quota_type_id", did[2])
		return nil
	} else if did[0] == "ECS" {
		d.Set("quota_id", object.Data.ID)
		d.Set("quota_type", did[1])
		d.Set("quota_type_id", did[2])
		return nil
	} else if did[0] == "SLB" {
		d.Set("quota_id", object.Data.ID)
		d.Set("quota_type", did[1])
		d.Set("quota_type_id", did[2])

		return nil
	} else if did[0] == "EIP" {
		d.Set("quota_id", object.Data.ID)
		d.Set("quota_type", did[1])
		d.Set("quota_type_id", did[2])

		return nil
	} else if did[0] == "RDS" {
		d.Set("quota_id", object.Data.ID)
		d.Set("quota_type", did[1])
		d.Set("quota_type_id", did[2])

		return nil
	} else if did[0] == "ODPS" {
		d.Set("quota_id", object.Data.ID)
		d.Set("quota_type", did[1])
		d.Set("quota_type_id", did[2])

		return nil
	} else if did[0] == "GPDB" {
		d.Set("quota_id", object.Data.ID)
		d.Set("quota_type", did[1])
		d.Set("quota_type_id", did[2])

		return nil
	} else if did[0] == "R-KVSTORE" {
		d.Set("quota_id", object.Data.ID)
		d.Set("quota_type", did[1])
		d.Set("quota_type_id", did[2])

		return nil
	} else if did[0] == "DDS" {
		d.Set("quota_id", object.Data.ID)
		d.Set("quota_type", did[1])
		d.Set("quota_type_id", did[2])

		return nil
	} else {

		d.Set("region_name", object.Data.RegionName)
		d.Set("product_name", object.Data.ProductName)
		d.Set("cluster_name", object.Data.Cluster)
		d.Set("total_cpu", object.Data.TotalCPU)
		d.Set("total_mem", object.Data.TotalMem)
		d.Set("total_gpu", object.Data.TotalGpu)
		d.Set("total_disk_cloud_ssd", object.Data.TotalDiskCloudSsd)
		d.Set("total_disk_cloud_efficiency", object.Data.TotalDiskCloudEfficiency)
		d.Set("total_amount", object.Data.TotalAmount)
	}

	return nil
}
func resourceApsaraStackAscmQuotaDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	check, err := ascmService.DescribeAscmQuota(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, did[0], "IsQuotaExist", ApsaraStackSdkGoERROR)
	}

	addDebug("IsQuotaExist", check, requestInfo, map[string]string{"productName": did[0]})
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"RegionName ":     client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Product":         "ascm",
			"Action":          "DeleteQuota",
			"Version":         "2019-05-10",
			"ProductName":     "ascm",
			"productName":     did[0],
			"QuotaType":       did[1],
			"QuotaTypeId":     did[2],
		}

		request.Method = "POST"
		request.Product = "ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		request.Scheme = "http"
		request.ApiName = "DeleteQuota"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		check, err := ascmService.DescribeAscmQuota(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		if check.Data.QuotaTypeID != 0 {
			return resource.RetryableError(Error("Trying to delete Quota %#v successfully.", d.Id()))
		}
		return nil
	})
	return nil
}
