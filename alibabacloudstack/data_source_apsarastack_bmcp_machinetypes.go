package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"
)

func dataSourceAlibabacloudStackBmcpMachineTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackBmcpMachineTypesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"arch_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values.
			"machinetypes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deploy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"manufacturer": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_arch": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_model": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_manufacturer": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gpu": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gpu_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gpu_manufacturer": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"video_memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tflops_fp32": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"network_card_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_card_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rated_power": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unit_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackBmcpMachineTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := make(map[string]interface{})
	reqQuery["PageNumber"] = 1
	reqQuery["PageSize"] = PageSizeLarge
	reqQuery["DeployTypes"] = `["bmcp","bmcp_managed","bmcp_no_clone","base","ehpc","ehpc_managed","aspeed"]`

	// Call ListMachineType API
	resp, err := client.DoTeaRequest("POST", "bms", "2022-05-30", "ListMachineType", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "ListMachineType", "POST", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Parse response data
	data, err := jmespath.Search("data.data", resp)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "ListMachineType", "data.data", resp)
	}

	items, ok := data.([]interface{})
	if !ok {
		items = []interface{}{}
	}

	var filteredItems []interface{}
	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		name, _ := itemMap["name"].(string)
		cpuArch, _ := itemMap["CPUArch"].(string)

		// Apply name regex filter
		if nameRegex, ok := d.GetOk("name_regex"); ok {
			if !strings.Contains(name, nameRegex.(string)) {
				continue
			}
		}

		// Apply arch regex filter
		if archRegex, ok := d.GetOk("arch_regex"); ok {
			if !strings.Contains(cpuArch, archRegex.(string)) {
				continue
			}
		}

		filteredItems = append(filteredItems, item)
	}

	return bmcpMachineTypesAttributes(d, filteredItems)
}

func bmcpMachineTypesAttributes(d *schema.ResourceData, items []interface{}) error {
	var ids []string
	var s []map[string]interface{}

	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		var idStr string
		switch v := itemMap["ID"].(type) {
		case string:
			idStr = v
		case json.Number:
			idStr = v.String()
		case float64:
			idStr = strconv.FormatInt(int64(v), 10)
		case int:
			idStr = strconv.Itoa(v)
		default:
			idStr = fmt.Sprintf("%v", v)
		}

		id := idStr
		name, _ := itemMap["name"].(string)
		description, _ := itemMap["description"].(string)
		deployType, _ := itemMap["deployType"].(string)
		manufacturer, _ := itemMap["manufacturer"].(string)
		cpuArch, _ := itemMap["CPUArch"].(string)
		cpuModel, _ := itemMap["CPUModel"].(string)
		cpuManufacturer, _ := itemMap["cpuManufacturer"].(string)
		cpuNumber, _ := itemMap["CPUNumber"].(float64)
		memory, _ := itemMap["memory"].(float64)
		disk, _ := itemMap["disk"].(float64)
		diskType, _ := itemMap["diskType"].(string)
		gpu, _ := itemMap["gpu"].(string)
		gpuNum, _ := itemMap["GpuNum"].(float64)
		gpuManufacturer, _ := itemMap["gpuManufacturer"].(string)
		videoMemory, _ := itemMap["VideoMemory"].(float64)
		tflopsFP32, _ := itemMap["TflopsFP32"].(float64)
		networkCardType, _ := itemMap["networkCardType"].(string)
		networkCardNum, _ := itemMap["networkCardNum"].(float64)
		ratedPower, _ := itemMap["ratedPower"].(string)
		specification, _ := itemMap["Specification"].(string)
		unitNum, _ := itemMap["unitNum"].(float64)
		createTime, _ := itemMap["createTime"].(string)
		updateTime, _ := itemMap["updateTime"].(string)

		mapping := map[string]interface{}{
			"id":                id,
			"name":              name,
			"description":       description,
			"deploy_type":       deployType,
			"manufacturer":      manufacturer,
			"cpu_arch":          cpuArch,
			"cpu_model":         cpuModel,
			"cpu_manufacturer":  cpuManufacturer,
			"cpu_number":        int(cpuNumber),
			"memory":            int(memory),
			"disk":              int(disk),
			"disk_type":         diskType,
			"gpu":               gpu,
			"gpu_num":           int(gpuNum),
			"gpu_manufacturer":  gpuManufacturer,
			"video_memory":      int(videoMemory),
			"tflops_fp32":       int(tflopsFP32),
			"network_card_type": networkCardType,
			"network_card_num":  int(networkCardNum),
			"rated_power":       ratedPower,
			"specification":     specification,
			"unit_num":          int(unitNum),
			"create_time":       createTime,
			"update_time":       updateTime,
		}

		ids = append(ids, id)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("machinetypes", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
