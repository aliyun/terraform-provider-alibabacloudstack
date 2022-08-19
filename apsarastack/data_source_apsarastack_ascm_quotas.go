package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strings"
)

func dataSourceApsaraStackQuotas() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackQuotasRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"product_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"quota_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"quota_type_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"quota_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_type_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"used_vip_public": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"allocate_vip_internal": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"allocate_vip_public": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_vip_public": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_vip_internal": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_vpc": {
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
							Computed: true,
						},
						"total_disk_cloud_ssd": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_disk_cloud_efficiency": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_amount": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_disk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_cu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_eip": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"used_disk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"allocate_disk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"allocate_cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"used_mem": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"target_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceApsaraStackQuotasRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.Method = "GET"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.ApiName = "GetQuota"
	productName := d.Get("product_name").(string)
	quotaType := d.Get("quota_type").(string)
	quotaTypeId := d.Get("quota_type_id").(string)
	targetType := d.Get("target_type").(string)
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeyId": client.AccessKey, "AccessKeySecret": client.SecretKey,
		"Product":       "ascm",
		"RegionId":      client.RegionId,
		"Department":    client.Department,
		"ResourceGroup": client.ResourceGroup,
		"Action":        "GetQuota",
		"Version":       "2019-05-10",
		"productName":   productName,
		"quotaType":     quotaType,
		"quotaTypeId":   quotaTypeId,
		"regionName":    client.RegionId,
		"targetType":    targetType,
	}
	response := AscmQuota{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw GetQuota : %s", raw)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_ascm_quotas", request.GetActionName(), ApsaraStackSdkGoERROR)
		}

		bresponse, _ := raw.(*responses.CommonResponse)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return WrapError(err)
		}
		if response.Code == "200" {
			break
		}

	}

	var ids []string
	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"id":                          response.Data.ID,
		"quota_type":                  response.Data.QuotaType,
		"quota_type_id":               fmt.Sprint(response.Data.QuotaTypeID),
		"target_type":                 response.Data.TargetType,
		"used_vip_public":             response.Data.UsedVipPublic,
		"allocate_vip_internal":       response.Data.AllocateVipInternal,
		"allocate_vip_public":         response.Data.AllocateVipPublic,
		"total_vip_public":            response.Data.TotalVipPublic,
		"total_vip_internal":          response.Data.TotalVipInternal,
		"region":                      response.Data.Region,
		"total_vpc":                   response.Data.TotalVPC,
		"total_cpu":                   response.Data.TotalCPU,
		"total_cu":                    response.Data.TotalCU,
		"total_disk":                  response.Data.TotalDisk,
		"total_mem":                   response.Data.TotalMem,
		"used_mem":                    response.Data.UsedMem,
		"total_gpu":                   response.Data.TotalGpu,
		"total_amount":                response.Data.TotalAmount,
		"total_disk_cloud_ssd":        response.Data.TotalDiskCloudSsd,
		"used_disk":                   response.Data.UsedDisk,
		"allocate_disk":               response.Data.AllocateDisk,
		"allocate_cpu":                response.Data.AllocateCPU,
		"total_eip":                   response.Data.TotalEIP,
		"total_disk_cloud_efficiency": response.Data.TotalDiskCloudEfficiency,
	}

	ids = append(ids, fmt.Sprint(response.Data.ID))
	s = append(s, mapping)

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("quotas", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
