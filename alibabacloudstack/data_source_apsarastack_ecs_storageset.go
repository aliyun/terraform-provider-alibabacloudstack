package alibabacloudstack

import (
	"encoding/json"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
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

	request := client.NewCommonRequest("GET", "Ecs", "2014-05-26", action, "")
	request.QueryParams["PageNumber"] = "1"
	request.QueryParams["PageSize"] = "20"

	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)
	raw, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
		return EcsClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_drds_instances", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(action, raw, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &addDomains)
	if err != nil {
		return errmsgs.WrapError(err)
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
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
