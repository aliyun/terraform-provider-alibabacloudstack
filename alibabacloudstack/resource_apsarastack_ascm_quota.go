package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
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
	}
	setResourceFunc(resource, resourceAlibabacloudStackAscmQuotaCreate, resourceAlibabacloudStackAscmQuotaRead, resourceAlibabacloudStackAscmQuotaUpdate, resourceAlibabacloudStackAscmQuotaDelete)
	return resource
}

func resourceAlibabacloudStackAscmQuotaCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	productName := strings.ToUpper(d.Get("product_name").(string))
	quotaType := d.Get("quota_type").(string)
	quotaTypeId := d.Get("quota_type_id").(string)

	var quotaBody string
	var targetType string
	switch productName {
	case "ECS":
		totalCpu := d.Get("total_cpu").(int)
		totalMem := d.Get("total_mem").(int)
		totalGpu := d.Get("total_gpu").(int)
		totalDiskCloudSsd := d.Get("total_disk_cloud_ssd").(int)
		totalDiskCloudEfficiency := d.Get("total_disk_cloud_efficiency").(int)
		quotaBody = fmt.Sprintf("{\"totalCpu\":%d,\"totalMem\":%d,\"totalGpu\":%d,\"totalDisk_cloud_ssd\":%d,\"totalDisk_cloud_efficiency\":%d}", totalCpu, totalMem, totalGpu, totalDiskCloudSsd, totalDiskCloudEfficiency)
	case "OSS":
		totalAmount := d.Get("total_amount").(int)
		quotaBody = fmt.Sprintf("{\"totalAmount\":%d}", totalAmount)
	case "EIP":
		totalEIP := d.Get("total_eip").(int)
		quotaBody = fmt.Sprintf("{\"totalEIP\":%d}", totalEIP)
	case "SLB":
		totalVipPublic := d.Get("total_vip_public").(int)
		totalVipInternal := d.Get("total_vip_internal").(int)
		quotaBody = fmt.Sprintf("{\"totalVipPublic\":%d,\"totalVipInternal\":%d}", totalVipPublic, totalVipInternal)
	case "ODPS":
		totalCu := d.Get("total_cu").(int)
		totalDisk := d.Get("total_disk").(int)
		quotaBody = fmt.Sprintf("{\"totalCu\":%d,\"totalDisk\":%d}", totalCu, totalDisk)
	case "DDS":
		totalCpu := d.Get("total_cpu").(int)
		totalMem := d.Get("total_mem").(int)
		totalDisk := d.Get("total_disk").(int)
		quotaBody = fmt.Sprintf("{\"totalCpu\":%d,\"totalMem\":%d,\"totalDisk\":%d}", totalCpu, totalMem, totalDisk)
	case "RDS":
		targetType = d.Get("target_type").(string)
		totalCpu := d.Get("total_cpu").(int)
		totalMem := d.Get("total_mem").(int)
		totalDisk := d.Get("total_disk").(int)
		quotaBody = fmt.Sprintf("{\"totalCpu\":%d,\"totalMem\":%d,\"totalDisk\":%d}", totalCpu, totalMem, totalDisk)
	case "GPDB":
		totalCpu := d.Get("total_cpu").(int)
		totalMem := d.Get("total_mem").(int)
		totalDisk := d.Get("total_disk").(int)
		quotaBody = fmt.Sprintf("{\"totalCpu\":%d,\"totalMem\":%d,\"totalDisk\":%d}", totalCpu, totalMem, totalDisk)
	case "R-KVSTORE":
		targetType = d.Get("target_type").(string)
		totalMem := d.Get("total_mem").(int)
		quotaBody = fmt.Sprintf("{\"totalMem\":%d}", totalMem)
	case "VPC":
		totalVPC := d.Get("total_vpc").(int)
		quotaBody = fmt.Sprintf("{\"totalVPC\":%d}", totalVPC)
	default:
		log.Print("Please Enter a valid Product Name.")
		return nil
	}

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "CreateQuota", "/ascm/manage/quota/add")
	request.QueryParams["regionName"] = client.RegionId
	request.QueryParams["quotaType"] = quotaType
	request.QueryParams["quotaTypeId"] = quotaTypeId
	request.QueryParams["productName"] = productName
	request.QueryParams["targetType"] = ""
	request.QueryParams["quotaBody"] = quotaBody
	request.QueryParams["targetType"] = targetType

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf(" response of raw CreateQuota : %s", raw)

	if err != nil {
		errmsg := ""
		if bresponse, ok := raw.(*responses.CommonResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_quota", "CreateQuota", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("CreateQuota", raw, nil, request)

	bresponse, ok := raw.(*responses.CommonResponse)
	if !ok || bresponse.GetHttpStatus() != 200 {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_quota", "CreateQuota", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("CreateQuota", raw, nil, bresponse.GetHttpContentString())

	d.SetId(productName + COLON_SEPARATED + quotaType + COLON_SEPARATED + quotaTypeId)

	return nil
}

func resourceAlibabacloudStackAscmQuotaUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	attributeUpdate := false
	check, err := ascmService.DescribeAscmQuota(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsQuotaExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	var quotaBody string
	switch did[0] {
	case "VPC":
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
		quotaBody = fmt.Sprintf("{\"totalVPC\":%d}", totalVPC)
	case "RDS":
		var totalDISK, totalCPU, totalMEM int
		if d.HasChange("total_mem") {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = totalMEM
			attributeUpdate = true
		}
		if d.HasChange("total_disk") {
			if v, ok := d.GetOk("total_disk"); ok {
				totalDISK = v.(int)
			}
			check.Data.TotalDisk = totalDISK
			attributeUpdate = true
		}
		if d.HasChange("total_cpu") {
			if v, ok := d.GetOk("total_cpu"); ok {
				totalCPU = v.(int)
			}
			check.Data.TotalCPU = totalCPU
			attributeUpdate = true
		}
		quotaBody = fmt.Sprintf("{\"totalCpu\":%d,\"totalMem\":%d,\"totalDisk\":%d}", totalCPU, totalMEM, totalDISK)
	case "EIP":
		var totalEIP int
		if d.HasChange("total_eip") {
			if v, ok := d.GetOk("total_eip"); ok {
				totalEIP = v.(int)
			}
			check.Data.TotalMem = totalEIP
			attributeUpdate = true
		}
		quotaBody = fmt.Sprintf("{\"totalEIP\":%d}", totalEIP)
	case "ECS":
		var totalGPU, totalCPU, totalMEM, totalDCE, totalDCS int
		if d.HasChange("total_cpu") {
			if v, ok := d.GetOk("total_cpu"); ok {
				totalCPU = v.(int)
			}
			check.Data.TotalCPU = totalCPU
			attributeUpdate = true
		}
		if d.HasChange("total_gpu") {
			if v, ok := d.GetOk("total_gpu"); ok {
				totalGPU = v.(int)
			}
			check.Data.TotalGpu = totalGPU
			attributeUpdate = true
		}
		if d.HasChange("total_mem") {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = totalMEM
			attributeUpdate = true
		}
		if d.HasChange("total_disk_cloud_ssd") {
			if v, ok := d.GetOk("total_disk_cloud_ssd"); ok {
				totalDCS = v.(int)
			}
			check.Data.TotalDiskCloudSsd = totalDCS
			attributeUpdate = true
		}
		if d.HasChange("total_disk_cloud_efficiency") {
			if v, ok := d.GetOk("total_disk_cloud_efficiency"); ok {
				totalDCE = v.(int)
			}
			check.Data.TotalDiskCloudEfficiency = totalDCE
			attributeUpdate = true
		}
		quotaBody = fmt.Sprintf("{\"totalCpu\":%d,\"totalMem\":%d,\"totalGpu\":%d,\"totalDisk_cloud_ssd\":%d,\"totalDisk_cloud_efficiency\":%d}", totalCPU, totalMEM, totalGPU, totalDCS, totalDCE)
	case "SLB":
		var totalVP, totalVI int
		if d.HasChange("total_vip_internal") {
			if v, ok := d.GetOk("total_vip_internal"); ok {
				totalVI = v.(int)
			}
			check.Data.TotalVipInternal = totalVI
			attributeUpdate = true
		}
		if d.HasChange("total_vip_public") {
			if v, ok := d.GetOk("total_vip_public"); ok {
				totalVP = v.(int)
			}
			check.Data.TotalVipPublic = totalVP
			attributeUpdate = true
		}
		quotaBody = fmt.Sprintf("{\"totalVipPublic\":%d,\"totalVipInternal\":%d}", totalVP, totalVI)
	case "OSS":
		var totalAmount int
		if d.HasChange("total_amount") {
			if v, ok := d.GetOk("total_amount"); ok {
				totalAmount = v.(int)
			}
			check.Data.TotalAmount = totalAmount
			attributeUpdate = true
		}
		quotaBody = fmt.Sprintf("{\"totalAmount\":%d}", totalAmount)
	case "DDS":
		var totalDISK, totalCPU, totalMEM int
		if d.HasChange("total_mem") {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = totalMEM
			attributeUpdate = true
		}
		if d.HasChange("total_disk") {
			if v, ok := d.GetOk("total_disk"); ok {
				totalDISK = v.(int)
			}
			check.Data.TotalDisk = totalDISK
			attributeUpdate = true
		}
		if d.HasChange("total_cpu") {
			if v, ok := d.GetOk("total_cpu"); ok {
				totalCPU = v.(int)
			}
			check.Data.TotalCPU = totalCPU
			attributeUpdate = true
		}
		quotaBody = fmt.Sprintf("{\"totalCpu\":%d,\"totalMem\":%d,\"totalDisk\":%d}", totalCPU, totalMEM, totalDISK)
	case "ODPS":
		var totalDISK, totalCU int
		if d.HasChange("total_disk") {
			if v, ok := d.GetOk("total_disk"); ok {
				totalDISK = v.(int)
			}
			check.Data.TotalDisk = totalDISK
			attributeUpdate = true
		}
		if d.HasChange("total_cu") {
			if v, ok := d.GetOk("total_cu"); ok {
				totalCU = v.(int)
			}
			check.Data.TotalCU = totalCU
			attributeUpdate = true
		}
		quotaBody = fmt.Sprintf("{\"totalCu\":%d,\"totalDisk\":%d}", totalCU, totalDISK)
	case "GPDB":
		var totalDISK, totalCPU, totalMEM int
		if d.HasChange("total_mem") {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = totalMEM
			attributeUpdate = true
		}
		if d.HasChange("total_disk") {
			if v, ok := d.GetOk("total_disk"); ok {
				totalDISK = v.(int)
			}
			check.Data.TotalDisk = totalDISK
			attributeUpdate = true
		}
		if d.HasChange("total_cpu") {
			if v, ok := d.GetOk("total_cpu"); ok {
				totalCPU = v.(int)
			}
			check.Data.TotalCPU = totalCPU
			attributeUpdate = true
		}
		quotaBody = fmt.Sprintf("{\"totalCpu\":%d,\"totalMem\":%d,\"totalDisk\":%d}", totalCPU, totalMEM, totalDISK)
	case "R-KVSTORE":
		var totalMEM int
		if d.HasChange("total_mem") {
			if v, ok := d.GetOk("total_mem"); ok {
				totalMEM = v.(int)
			}
			check.Data.TotalMem = totalMEM
			attributeUpdate = true
		}
		quotaBody = fmt.Sprintf("{\"totalMem\":%d}", totalMEM)
	}

	if attributeUpdate {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "UpdateQuota", "/ascm/manage/quota/update")
		request.QueryParams["regionName"] = client.RegionId
		request.QueryParams["quotaType"] = did[1]
		request.QueryParams["quotaTypeId"] = did[2]
		request.QueryParams["productName"] = did[0]
		request.QueryParams["targetType"] = ""
		request.QueryParams["quotaBody"] = quotaBody

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*responses.CommonResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_quota", "UpdateQuota", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, nil, request)
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

	if did[0] == "VPC" || did[0] == "OSS" || did[0] == "ECS" || did[0] == "SLB" || did[0] == "EIP" || did[0] == "RDS" || did[0] == "ODPS" || did[0] == "GPDB" || did[0] == "R-KVSTORE" || did[0] == "DDS" {
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

func resourceAlibabacloudStackAscmQuotaDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	check, err := ascmService.DescribeAscmQuota(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, did[0], "IsQuotaExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	addDebug("IsQuotaExist", check, requestInfo, map[string]string{"productName": did[0]})
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "DeleteQuota", "/ascm/manage/quota/delete")
		request.QueryParams["productName"] = did[0]
		request.QueryParams["quotaType"] = did[1]
		request.QueryParams["quotaTypeId"] = did[2]

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
			return resource.RetryableError(errmsgs.Error("Trying to delete Quota %#v successfully.", d.Id()))
		}
		return nil
	})
	return nil
}