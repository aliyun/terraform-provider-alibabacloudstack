package alibabacloudstack

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// productFieldMappings defines product-specific mapping from old schema fields to API keys
var productFieldMappings = map[string]map[string]string{
	"ecs": {
		"total_cpu":                   "totalCpu",
		"total_mem":                   "totalMem",
		"total_gpu":                   "totalGpu",
		"total_disk_cloud_ssd":        "totalDisk_cloud_ssd",
		"total_disk_cloud_efficiency": "totalDisk_cloud_efficiency",
	},
	"oss":       {"total_amount": "totalAmount"},
	"eip":       {"total_eip": "totalEIP"},
	"slb":       {"total_vip_public": "totalVipPublic", "total_vip_internal": "totalVipInternal"},
	"odps":      {"total_cu": "totalCu", "total_disk": "totalDisk"},
	"dds":       {"total_cpu": "totalCpu", "total_mem": "totalMem", "total_disk": "totalDisk"},
	"rds":       {"total_cpu": "totalCpu", "total_mem": "totalMem", "total_disk": "totalDisk"},
	"gpdb":      {"total_cpu": "totalCpu", "total_mem": "totalMem", "total_disk": "totalDisk"},
	"r-kvstore": {"total_mem": "totalMem"},
	"vpc":       {"total_vpc": "totalVPC"},
}

func resourceAlibabacloudStackAscmQuota() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ECS", "OSS", "VPC", "DRDS", "RDS", "SLB", "ODPS", "EIP", "GPDB", "R-KVSTORE", "NAS", "DDS"}, false),
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
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "This field is deprecated. Only the current region is supported. This field will be removed in version 3.21.0.",
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
		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
			productName := d.Get("product_name").(string)
			targetType := d.Get("target_type").(string)
			if (productName == "RDS" || productName == "R-KVSTORE") && targetType == "" {
				return fmt.Errorf("target_type is required when product_name is %s", productName)
			}

			// Get product-specific field mapping
			fieldMapping, ok := productFieldMappings[strings.ToLower(productName)]
			if !ok {
				return fmt.Errorf("unsupported product_name: %s", productName)
			}

			// Check if at least one old field is defined with value > 0
			atLeastOneSet := false
			for oldField := range fieldMapping {
				if v, ok := d.GetOk(oldField); ok && v.(int) > 0 {
					atLeastOneSet = true
					break
				}
			}

			if !atLeastOneSet {
				return fmt.Errorf("at least one quota item must be set for product %s", productName)
			}

			return nil
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAscmQuotaCreate, resourceAlibabacloudStackAscmQuotaRead, resourceAlibabacloudStackAscmQuotaUpdate, resourceAlibabacloudStackAscmQuotaDelete)
	return resource
}

func resourceAlibabacloudStackAscmQuotaCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	productName := strings.ToLower(d.Get("product_name").(string))
	quotaType := d.Get("quota_type").(string)
	quotaTypeId := d.Get("quota_type_id").(string)

	var quotaBody string
	var targetType string

	// Build quotaBody from old fixed fields
	quotaBodyMap := make(map[string]int)
	if fieldMapping, ok := productFieldMappings[productName]; ok {
		for oldField, newKey := range fieldMapping {
			if v := d.Get(oldField).(int); v > 0 {
				quotaBodyMap[newKey] = v
			}
		}
		quotaBodyBytes, _ := json.Marshal(quotaBodyMap)
		quotaBody = string(quotaBodyBytes)
	} else {
		log.Print("Please Enter a valid Product Name.")
		return nil
	}

	// For RDS/R-KVSTORE, need targetType
	if productName == "rds" || productName == "r-kvstore" {
		targetType = d.Get("target_type").(string)
	}

	// Validate that at least one quota item is set (items with value <= 0 don't count)
	if quotaBody == "" || quotaBody == "{}" {
		return fmt.Errorf("at least one quota item must be set for product %s", productName)
	}

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "CreateQuota", "/ascm/manage/quota/add")
	request.QueryParams["ProductName"] = productName
	request.QueryParams["QuotaBody"] = quotaBody
	request.QueryParams["QuotaType"] = quotaType
	request.QueryParams["QuotaTypeId"] = quotaTypeId
	request.QueryParams["RegionName"] = client.RegionId
	if targetType != "" {
		request.QueryParams["TargetType"] = targetType
	}
	log.Printf("[DEBUG] CreateQuota request QueryParams: %v", request.QueryParams)

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		log.Printf("[ERROR] CreateQuota response body: %s", bresponse.GetHttpContentString())
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_quota", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	log.Printf("[DEBUG] CreateQuota response: %s", bresponse.GetHttpContentString())

	d.SetId(productName + COLON_SEPARATED + quotaType + COLON_SEPARATED + quotaTypeId)

	return nil
}

func resourceAlibabacloudStackAscmQuotaUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.IsNewResource() {
		return nil
	}
	client := meta.(*connectivity.AlibabacloudStackClient)
	did := strings.Split(d.Id(), COLON_SEPARATED)

	var quotaBody string
	var targetType string

	// Build quotaBody from old fixed fields: send all fields for this product
	// to prevent API from resetting unchanged fields to -1
	quotaBodyMap := make(map[string]int)
	if fieldMapping, ok := productFieldMappings[did[0]]; ok {
		for oldField, newKey := range fieldMapping {
			v := d.Get(oldField).(int)
			if v > 0 {
				quotaBodyMap[newKey] = v
			} else {
				quotaBodyMap[newKey] = -1
			}
		}
	}
	if len(quotaBodyMap) == 0 {
		return nil
	}
	quotaBodyBytes, _ := json.Marshal(quotaBodyMap)
	quotaBody = string(quotaBodyBytes)

	// For RDS/R-KVSTORE, need targetType
	if did[0] == "rds" || did[0] == "r-kvstore" {
		targetType = d.Get("target_type").(string)
	}

	quotaTypeIdInt, _ := strconv.Atoi(did[2])
	var quotaBodyObj interface{}
	json.Unmarshal([]byte(quotaBody), &quotaBodyObj)
	requestBody := map[string]interface{}{
		"productName": did[0],
		"quotaType":   did[1],
		"quotaTypeId": quotaTypeIdInt,
		"regionName":  client.RegionId,
		"quotaBody":   quotaBodyObj,
	}
	if targetType != "" {
		requestBody["targetType"] = targetType
	}
	bodyJSON, _ := json.Marshal(requestBody)
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "UpdateQuota", "/ascm/manage/quota/update")
	request.Content = bodyJSON
	request.SetContentType("application/json")
	log.Printf("[DEBUG] UpdateQuota request body: %s", string(bodyJSON))

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		log.Printf("[ERROR] UpdateQuota response body: %s", bresponse.GetHttpContentString())
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_quota", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	return nil
}

func resourceAlibabacloudStackAscmQuotaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmQuota(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	if object.Code != "200" {
		d.SetId("")
		return nil
	}

	// Set common fields for all products
	d.Set("quota_id", object.Data.ID)
	d.Set("product_name", strings.ToUpper(did[0]))
	d.Set("quota_type", did[1])
	d.Set("quota_type_id", did[2])

	// Set old fixed fields from API response
	// If API returns -1, set to 0 (delete the quota item)
	if object.Data.TotalCPU > 0 {
		d.Set("total_cpu", int(object.Data.TotalCPU))
	} else {
		d.Set("total_cpu", 0)
	}
	if object.Data.TotalMem > 0 {
		d.Set("total_mem", int(object.Data.TotalMem))
	} else {
		d.Set("total_mem", 0)
	}
	if object.Data.TotalGpu > 0 {
		d.Set("total_gpu", int(object.Data.TotalGpu))
	} else {
		d.Set("total_gpu", 0)
	}
	if object.Data.TotalDiskCloudSsd > 0 {
		d.Set("total_disk_cloud_ssd", int(object.Data.TotalDiskCloudSsd))
	} else {
		d.Set("total_disk_cloud_ssd", 0)
	}
	if object.Data.TotalDiskCloudEfficiency > 0 {
		d.Set("total_disk_cloud_efficiency", int(object.Data.TotalDiskCloudEfficiency))
	} else {
		d.Set("total_disk_cloud_efficiency", 0)
	}
	if object.Data.TotalVPC > 0 {
		d.Set("total_vpc", int(object.Data.TotalVPC))
	} else {
		d.Set("total_vpc", 0)
	}
	if object.Data.TotalEIP > 0 {
		d.Set("total_eip", int(object.Data.TotalEIP))
	} else {
		d.Set("total_eip", 0)
	}
	if object.Data.TotalAmount > 0 {
		d.Set("total_amount", int(object.Data.TotalAmount))
	} else {
		d.Set("total_amount", 0)
	}
	if object.Data.TotalDisk > 0 {
		d.Set("total_disk", int(object.Data.TotalDisk))
	} else {
		d.Set("total_disk", 0)
	}
	if object.Data.TotalCU > 0 {
		d.Set("total_cu", int(object.Data.TotalCU))
	} else {
		d.Set("total_cu", 0)
	}
	if object.Data.TotalVipInternal > 0 {
		d.Set("total_vip_internal", int(object.Data.TotalVipInternal))
	} else {
		d.Set("total_vip_internal", 0)
	}
	if object.Data.TotalVipPublic > 0 {
		d.Set("total_vip_public", int(object.Data.TotalVipPublic))
	} else {
		d.Set("total_vip_public", 0)
	}

	return nil
}

func resourceAlibabacloudStackAscmQuotaDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	check, err := ascmService.DescribeAscmQuota(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[WARN] DeleteQuota: quota %s not found, treating as already deleted", d.Id())
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, did[0], "IsQuotaExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	addDebug("IsQuotaExist", check, nil, map[string]string{"productName": did[0]})
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "DeleteQuota", "/ascm/manage/quota/delete")
		request.FormParams["productName"] = did[0]
		request.FormParams["quotaType"] = did[1]
		request.FormParams["quotaTypeId"] = did[2]
		request.FormParams["regionName"] = client.RegionId

		_, err := client.ProcessCommonRequest(request)
		if err != nil {
			if strings.Contains(err.Error(), "EntityNotExist") || strings.Contains(err.Error(), "quota_not_found") {
				log.Printf("[WARN] DeleteQuota: quota %s already deleted (%v)", d.Id(), err)
				return nil
			}
			log.Printf("[ERROR] DeleteQuota response: %v", err)
			return resource.RetryableError(err)
		}
		check, err := ascmService.DescribeAscmQuota(d.Id())

		if err != nil {
			if errmsgs.NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		if check.Data.QuotaTypeID != 0 {
			return resource.RetryableError(errmsgs.Error("Trying to delete Quota %#v successfully.", d.Id()))
		}
		return nil
	})
	return nil
}
