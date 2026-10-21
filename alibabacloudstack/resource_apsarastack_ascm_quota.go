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
		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
			productName := d.Get("product_name").(string)
			targetType := d.Get("target_type").(string)
			if (productName == "RDS" || productName == "R-KVSTORE") && targetType == "" {
				return fmt.Errorf("target_type is required when product_name is %s", productName)
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
	switch productName {
	case "ecs":
		totalCpu := d.Get("total_cpu").(int)
		totalMem := d.Get("total_mem").(int)
		totalGpu := d.Get("total_gpu").(int)
		totalDiskCloudSsd := d.Get("total_disk_cloud_ssd").(int)
		totalDiskCloudEfficiency := d.Get("total_disk_cloud_efficiency").(int)
		quotaBody = fmt.Sprintf("{\"totalCpu\":%d,\"totalMem\":%d,\"totalGpu\":%d,\"totalDisk_cloud_ssd\":%d,\"totalDisk_cloud_efficiency\":%d}", totalCpu, totalMem, totalGpu, totalDiskCloudSsd, totalDiskCloudEfficiency)
	case "oss":
		totalAmount := d.Get("total_amount").(int)
		quotaBody = fmt.Sprintf("{\"totalAmount\":%d}", totalAmount)
	case "eip":
		totalEIP := d.Get("total_eip").(int)
		quotaBody = fmt.Sprintf("{\"totalEIP\":%d}", totalEIP)
	case "slb":
		totalVipPublic := d.Get("total_vip_public").(int)
		totalVipInternal := d.Get("total_vip_internal").(int)
		quotaBody = fmt.Sprintf("{\"totalVipPublic\":%d,\"totalVipInternal\":%d}", totalVipPublic, totalVipInternal)
	case "odps":
		totalCu := d.Get("total_cu").(int)
		totalDisk := d.Get("total_disk").(int)
		quotaBody = fmt.Sprintf("{\"totalCu\":%d,\"totalDisk\":%d}", totalCu, totalDisk)
	case "dds":
		totalCpu := d.Get("total_cpu").(int)
		totalMem := d.Get("total_mem").(int)
		totalDisk := d.Get("total_disk").(int)
		quotaBody = fmt.Sprintf("{\"totalCpu\":%d,\"totalMem\":%d,\"totalDisk\":%d}", totalCpu, totalMem, totalDisk)
	case "rds":
		targetType = d.Get("target_type").(string)
		totalCpu := d.Get("total_cpu").(int)
		totalMem := d.Get("total_mem").(int)
		totalDisk := d.Get("total_disk").(int)
		quotaBody = fmt.Sprintf("{\"totalCpu\":%d,\"totalMem\":%d,\"totalDisk\":%d}", totalCpu, totalMem, totalDisk)
	case "gpdb":
		totalCpu := d.Get("total_cpu").(int)
		totalMem := d.Get("total_mem").(int)
		totalDisk := d.Get("total_disk").(int)
		quotaBody = fmt.Sprintf("{\"totalCpu\":%d,\"totalMem\":%d,\"totalDisk\":%d}", totalCpu, totalMem, totalDisk)
	case "r-kvstore":
		targetType = d.Get("target_type").(string)
		totalMem := d.Get("total_mem").(int)
		quotaBody = fmt.Sprintf("{\"totalMem\":%d}", totalMem)
	case "vpc":
		totalVPC := d.Get("total_vpc").(int)
		quotaBody = fmt.Sprintf("{\"totalVPC\":%d}", totalVPC)
	default:
		log.Print("Please Enter a valid Product Name.")
		return nil
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

	d.SetId(productName + COLON_SEPARATED + quotaType + COLON_SEPARATED + quotaTypeId)

	return nil
}

func resourceAlibabacloudStackAscmQuotaUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.IsNewResource() {
		return nil
	}
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	attributeUpdate := false
	check, err := ascmService.DescribeAscmQuota(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsQuotaExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	var quotaBody string
	var targetType string
	switch did[0] {
	case "vpc":
		var totalVPC int
		if d.HasChange("total_vpc") {
			if v, ok := d.GetOk("total_vpc"); ok {
				totalVPC = v.(int)
			}
			check.Data.TotalVPC = float64(totalVPC)
			attributeUpdate = true
		} else {
			if v, ok := d.GetOk("total_vpc"); ok {
				totalVPC = v.(int)
			}
			check.Data.TotalVPC = float64(totalVPC)
		}
		quotaBody = fmt.Sprintf("{\"totalVPC\":%d}", totalVPC)
	case "rds":
		targetType = d.Get("target_type").(string)
		var totalDISK, totalCPU, totalMEM int
		if d.HasChange("total_mem") {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = float64(totalMEM)
			attributeUpdate = true
		}
		if d.HasChange("total_disk") {
			if v, ok := d.GetOk("total_disk"); ok {
				totalDISK = v.(int)
			}
			check.Data.TotalDisk = float64(totalDISK)
			attributeUpdate = true
		}
		if d.HasChange("total_cpu") {
			if v, ok := d.GetOk("total_cpu"); ok {
				totalCPU = v.(int)
			}
			check.Data.TotalCPU = float64(totalCPU)
			attributeUpdate = true
		}
		quotaBody = fmt.Sprintf("{\"totalCpu\":%d,\"totalMem\":%d,\"totalDisk\":%d}", totalCPU, totalMEM, totalDISK)
	case "eip":
		var totalEIP int
		if d.HasChange("total_eip") {
			if v, ok := d.GetOk("total_eip"); ok {
				totalEIP = v.(int)
			}
			check.Data.TotalMem = float64(totalEIP)
			attributeUpdate = true
		}
		quotaBody = fmt.Sprintf("{\"totalEIP\":%d}", totalEIP)
	case "ecs":
		var totalGPU, totalCPU, totalMEM, totalDCE, totalDCS int
		if d.HasChange("total_cpu") {
			if v, ok := d.GetOk("total_cpu"); ok {
				totalCPU = v.(int)
			}
			check.Data.TotalCPU = float64(totalCPU)
			attributeUpdate = true
		}
		if d.HasChange("total_gpu") {
			if v, ok := d.GetOk("total_gpu"); ok {
				totalGPU = v.(int)
			}
			check.Data.TotalGpu = float64(totalGPU)
			attributeUpdate = true
		}
		if d.HasChange("total_mem") {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = float64(totalMEM)
			attributeUpdate = true
		}
		if d.HasChange("total_disk_cloud_ssd") {
			if v, ok := d.GetOk("total_disk_cloud_ssd"); ok {
				totalDCS = v.(int)
			}
			check.Data.TotalDiskCloudSsd = float64(totalDCS)
			attributeUpdate = true
		}
		if d.HasChange("total_disk_cloud_efficiency") {
			if v, ok := d.GetOk("total_disk_cloud_efficiency"); ok {
				totalDCE = v.(int)
			}
			check.Data.TotalDiskCloudEfficiency = float64(totalDCE)
			attributeUpdate = true
		}
		quotaBody = fmt.Sprintf("{\"totalCpu\":%d,\"totalMem\":%d,\"totalGpu\":%d,\"totalDisk_cloud_ssd\":%d,\"totalDisk_cloud_efficiency\":%d}", totalCPU, totalMEM, totalGPU, totalDCS, totalDCE)
	case "slb":
		var totalVP, totalVI int
		if d.HasChange("total_vip_internal") {
			if v, ok := d.GetOk("total_vip_internal"); ok {
				totalVI = v.(int)
			}
			check.Data.TotalVipInternal = float64(totalVI)
			attributeUpdate = true
		}
		if d.HasChange("total_vip_public") {
			if v, ok := d.GetOk("total_vip_public"); ok {
				totalVP = v.(int)
			}
			check.Data.TotalVipPublic = float64(totalVP)
			attributeUpdate = true
		}
		quotaBody = fmt.Sprintf("{\"totalVipPublic\":%d,\"totalVipInternal\":%d}", totalVP, totalVI)
	case "oss":
		var totalAmount int
		if d.HasChange("total_amount") {
			if v, ok := d.GetOk("total_amount"); ok {
				totalAmount = v.(int)
			}
			check.Data.TotalAmount = float64(totalAmount)
			attributeUpdate = true
		}
		quotaBody = fmt.Sprintf("{\"totalAmount\":%d}", totalAmount)
	case "dds":
		var totalDISK, totalCPU, totalMEM int
		if d.HasChange("total_mem") {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = float64(totalMEM)
			attributeUpdate = true
		}
		if d.HasChange("total_disk") {
			if v, ok := d.GetOk("total_disk"); ok {
				totalDISK = v.(int)
			}
			check.Data.TotalDisk = float64(totalDISK)
			attributeUpdate = true
		}
		if d.HasChange("total_cpu") {
			if v, ok := d.GetOk("total_cpu"); ok {
				totalCPU = v.(int)
			}
			check.Data.TotalCPU = float64(totalCPU)
			attributeUpdate = true
		}
		quotaBody = fmt.Sprintf("{\"totalCpu\":%d,\"totalMem\":%d,\"totalDisk\":%d}", totalCPU, totalMEM, totalDISK)
	case "odps":
		var totalDISK, totalCU int
		if d.HasChange("total_disk") {
			if v, ok := d.GetOk("total_disk"); ok {
				totalDISK = v.(int)
			}
			check.Data.TotalDisk = float64(totalDISK)
			attributeUpdate = true
		}
		if d.HasChange("total_cu") {
			if v, ok := d.GetOk("total_cu"); ok {
				totalCU = v.(int)
			}
			check.Data.TotalCU = float64(totalCU)
			attributeUpdate = true
		}
		quotaBody = fmt.Sprintf("{\"totalCu\":%d,\"totalDisk\":%d}", totalCU, totalDISK)
	case "gpdb":
		var totalDISK, totalCPU, totalMEM int
		if d.HasChange("total_mem") {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = float64(totalMEM)
			attributeUpdate = true
		}
		if d.HasChange("total_disk") {
			if v, ok := d.GetOk("total_disk"); ok {
				totalDISK = v.(int)
			}
			check.Data.TotalDisk = float64(totalDISK)
			attributeUpdate = true
		}
		if d.HasChange("total_cpu") {
			if v, ok := d.GetOk("total_cpu"); ok {
				totalCPU = v.(int)
			}
			check.Data.TotalCPU = float64(totalCPU)
			attributeUpdate = true
		}
		quotaBody = fmt.Sprintf("{\"totalCpu\":%d,\"totalMem\":%d,\"totalDisk\":%d}", totalCPU, totalMEM, totalDISK)
	case "r-kvstore":
		targetType = d.Get("target_type").(string)
		var totalMEM int
		if d.HasChange("total_mem") {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = float64(totalMEM)
			attributeUpdate = true
		}
		quotaBody = fmt.Sprintf("{\"totalMem\":%d}", totalMEM)
	}

	if attributeUpdate {
		quotaTypeIdInt, _ := strconv.Atoi(did[2])
		var quotaBodyMap interface{}
		json.Unmarshal([]byte(quotaBody), &quotaBodyMap)
		requestBody := map[string]interface{}{
			"productName": did[0],
			"quotaType":   did[1],
			"quotaTypeId": quotaTypeIdInt,
			"regionName":  client.RegionId,
			"quotaBody":   quotaBodyMap,
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

	if did[0] == "vpc" || did[0] == "oss" || did[0] == "ecs" || did[0] == "slb" || did[0] == "eip" || did[0] == "rds" || did[0] == "odps" || did[0] == "gpdb" || did[0] == "r-kvstore" || did[0] == "dds" {
		d.Set("quota_id", object.Data.ID)
		d.Set("product_name", strings.ToUpper(did[0]))
		d.Set("quota_type", did[1])
		d.Set("quota_type_id", did[2])
		d.Set("total_cpu", int(object.Data.TotalCPU))
		d.Set("total_mem", int(object.Data.TotalMem))
		d.Set("total_gpu", int(object.Data.TotalGpu))
		d.Set("total_disk_cloud_ssd", int(object.Data.TotalDiskCloudSsd))
		d.Set("total_disk_cloud_efficiency", int(object.Data.TotalDiskCloudEfficiency))
		d.Set("total_vpc", int(object.Data.TotalVPC))
		d.Set("total_eip", int(object.Data.TotalEIP))
		d.Set("total_amount", int(object.Data.TotalAmount))
		d.Set("total_disk", int(object.Data.TotalDisk))
		d.Set("total_cu", int(object.Data.TotalCU))
		d.Set("total_vip_internal", int(object.Data.TotalVipInternal))
		d.Set("total_vip_public", int(object.Data.TotalVipPublic))
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
