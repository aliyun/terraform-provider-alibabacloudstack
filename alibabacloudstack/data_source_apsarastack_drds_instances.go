package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackDRDSInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDRDSInstancesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.StringIsValidRegExp,
				Deprecated:  "Field 'name_regex' is deprecated and will be removed in a future release. Please use new field 'description_regex' instead.",
				ConflictsWith: []string{"description_regex"},
			},
			"description_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.StringIsValidRegExp,
				ConflictsWith: []string{"description_regex"},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			// Computed values
			"descriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackDRDSInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := drds.CreateDescribeDrdsInstancesRequest()
	client.InitRpcRequest(*request.RpcRequest)

	var dbi []drds.Instance
	var regexString *regexp.Regexp
	nameRegex := connectivity.GetResourceData(d, "description_regex", "name_regex").(string)
	if r, err := regexp.Compile(nameRegex); err == nil {
		regexString = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsInstances(request)
	})
	response, ok := raw.(*drds.DescribeDrdsInstancesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_drds_instances", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	for _, item := range response.Instances.Instance {
		if regexString != nil {
			if !regexString.MatchString(item.Description) {
				continue
			}
		}

		if len(idsMap) > 0 {
			if _, ok := idsMap[item.DrdsInstanceId]; !ok {
				continue
			}
		}

		dbi = append(dbi, item)
	}
	return drdsInstancesDescription(d, dbi)
}

func drdsInstancesDescription(d *schema.ResourceData, dbi []drds.Instance) error {
	var ids []string
	var descriptions []string
	var s []map[string]interface{}
	for _, item := range dbi {
		mapping := map[string]interface{}{
			"id":            item.DrdsInstanceId,
			"description":   item.Description,
			"type":          item.Type,
			"create_time":   item.CreateTime,
			"status":        item.Status,
			"network_type":  item.NetworkType,
			"zone_id":       item.ZoneId,
			"version":       item.Version,
		}
		ids = append(ids, item.DrdsInstanceId)
		descriptions = append(descriptions, item.Description)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("descriptions", descriptions); err != nil {
		return errmsgs.WrapError(err)
	}
	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
