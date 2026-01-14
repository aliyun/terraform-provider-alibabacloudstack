package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackLindormInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackLindormInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsValidRegExp,
				Deprecated:    "Field 'name_regex' is deprecated and will be removed in a future release. Please use new field 'description_regex' instead.",
				ConflictsWith: []string{"description_regex"},
			},
			"description_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsValidRegExp,
				ConflictsWith: []string{"name_regex"},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_brand": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ascm_create_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ali_uid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackLindormInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"PageSize":   10000,
		"PageNumber": 1,
		"Department": client.Department,
		"RegionId":   client.RegionId,
	}
	response := make(map[string]interface{})
	// Call request_params_handler

	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v.(string)
	}
	// request["PageNumber"] = "10000"
	// request["PageSize"] = "1"
	response, err := client.DoTeaRequest("POST", "hitsdb", "2020-06-15", "GetLindormInstanceList", "", nil, request, nil)
	if err != nil {
		if response == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		return errmsgs.WrapError(err)
	}
	idsMap := getIdsStringFilter(d)

	var ids []string
	datas := make([]interface{}, 0)
	for _, v := range response["InstanceList"].([]interface{}) {
		data := v.(map[string]interface{})
		if description_regex, ok := connectivity.GetResourceDataOk(d, "description_regex", "name_regex"); ok {
			r := regexp.MustCompile(description_regex.(string))
			if !r.MatchString(data["InstanceAlias"].(string)) {
				continue
			}
		}

		if len(idsMap) > 0 {
			if _, exist := idsMap[data["InstanceId"].(string)]; !exist {
				continue
			}
		}

		i := map[string]interface{}{
			"id":               data["InstanceId"],
			"cpu_brand":        data["CpuBrand"],
			"instance_storage": data["InstanceStorage"],
			"zone_id":          data["ZoneId"],
			"instance_id":      data["InstanceId"],
			"create_time":      data["CreateTime"],
			"ascm_create_user": data["AscmCreateUser"],
			"instance_alias":   data["InstanceAlias"],
			"network_type":     data["NetworkType"],
			"service_type":     data["ServiceType"],
			"engine_type":      data["EngineType"],
			"ali_uid":          data["AliUid"],
		}

		datas = append(datas, i)

		ids = append(ids, data["InstanceId"].(string))
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", datas); err != nil {
		return err
	}
	if err := d.Set("ids", ids); err != nil {
		return err
	}
	return nil
}
