package alibabacloudstack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"regexp"
)

func dataSourceAlibabacloudStackAscmMeteringQueryEcs() *schema.Resource {
	return &schema.Resource{
		Read:    dataSourceAlibabacloudStackAscmMeteringQueryEcsRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"product_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ins_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type_family": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"os_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"org_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_g_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sys_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gpu_amount": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAscmMeteringQueryEcsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	starttime := d.Get("start_time").(string)
	endtime := d.Get("end_time").(string)

	request := client.NewCommonRequest("GET", "ascm", "2019-05-10", "MeteringWebQuery", "")
	request.QueryParams["StartTime"] = starttime
	request.QueryParams["EndTime"] = endtime
	request.QueryParams["productName"] = "ECS"

	response := MeteringQueryDataEcs{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" rsponse of raw MeteringWebQuery : %s", raw)
		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*responses.CommonResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_metering_query", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		bresponse, _ := raw.(*responses.CommonResponse)

		_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)

		if response.Success == true || response.Data[0].StartTime == starttime {
			break
		}
	}
	var r *regexp.Regexp
	if rt, ok := d.GetOk("name_regex"); ok && rt.(string) != "" {
		r = regexp.MustCompile(rt.(string))
	}
	var ids []string
	var s []map[string]interface{}

	for _, data := range response.Data {
		if r != nil && !r.MatchString(data.StartTime) {
			continue
		}
		mapping := map[string]interface{}{
			"private_ip_address":    data.PrivateIPAddress,
			"instance_type_family":  data.InstanceTypeFamily,
			"memory":                data.Memory,
			"cpu":                   data.CPU,
			"os_name":               data.OSName,
			"org_name":              data.OrgName,
			"instance_network_type": data.InstanceNetworkType,
			"eip_address":           data.EipAddress,
			"resource_g_name":       data.ResourceGName,
			"instance_type":         data.InstanceType,
			"status":                data.Status,
			"sys_disk_size":         data.SysDiskSize,
			"gpu_amount":            data.GPUAmount,
			"instance_name":         data.InstanceName,
			"vpc_id":                data.VpcID,
			"data_disk_size":        data.DataDiskSize,
			"start_time":            data.StartTime,
			"end_time":              data.EndTime,
			"create_time":           data.CreateTime,
		}

		ids = append(ids, data.InsID)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("data", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}
	return nil
}
