package alibabacloudstack

import (
	"encoding/json"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/sdk_patch/datahub_patch"
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
			"output_file": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.",
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
	var addDomains = &datahub_patch.EcsDescribeEcsEbsStorageSetsResult{}
	action := "DescribeStorageSets"

	request := client.NewCommonRequest("GET", "Ecs", "2014-05-26", action, "")
	request.QueryParams["PageNumber"] = "1"
	request.QueryParams["PageSize"] = "20"

	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)
	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_listener", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

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
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}

	return nil
}
