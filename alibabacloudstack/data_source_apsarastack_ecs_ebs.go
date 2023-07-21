package alibabacloudstack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"strings"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackEcsEbsStorageSets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackEcsEbsStorageSetsRead,
		Schema: map[string]*schema.Schema{
			"storage_set_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"maxpartition_number": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"storage_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values
			"storages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_set_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_set_partition_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackEcsEbsStorageSetsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var addDomains = &datahub.EcsDescribeEcsEbsStorageSetsResult{}
	action := "DescribeStorageSets"
	request := requests.NewCommonRequest()
	request.Method = "GET"
	request.Product = "Ecs"
	request.Domain = client.Domain
	request.Version = "2014-05-26"
	request.RegionId = client.RegionId
	request.ApiName = action
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers["x-ascm-product-name"] = "Ecs"
	request.Headers["x-acs-organizationId"] = client.Department
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"RegionId":        client.RegionId,
		"Product":         "Ecs",
		"Version":         "2014-05-26",
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          action,
		"PageNumber":      "1",
	}
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(20)
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	raw, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
		return EcsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_drds_instances", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(action, raw, request)

	response, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(response.GetHttpContentBytes(), &addDomains)
	//v, err := jsonpath.Get("$.Commands.Command", bresponse)
	if err != nil {
		return WrapError(err)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range addDomains.StorageSets.StorageSet {
		mapping := map[string]interface{}{
			"storage_set_id":               object.StorageSetId,
			"storage_set_partition_number": object.StorageSetPartitionNumber,
			"storage_set_name":             object.StorageSetName,
		}
		ids = append(ids, object.StorageSetId)
		names = append(names, object.StorageSetName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("storages", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
