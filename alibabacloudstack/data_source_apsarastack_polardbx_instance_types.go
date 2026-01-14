package alibabacloudstack

import (
	"regexp"
	"sort"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackPolardbxInstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackPolardbxInstanceTypesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"spec_series": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"SHARE", "SINGLE"}, false),
			},
			"engine_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"5.7", "8.0"}, false),
			},
			"cpu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cpu_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"intel", "arm64", "Hygon"}, false),
			},
			"spec_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"DN", "CN"}, false),
			},
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"sorted_by": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CPU", "Memory"}, false),
			},
			"series": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "enterprise",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"enterprise", "standard"}, false),
			},
			// Computed values.
			"instance_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackPolardbxInstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	filterIds := getIdsStringFilter(d)

	existedId := map[string]string{}
	ids := []string{}
	types := []map[string]interface{}{}

	reqQuery := map[string]interface{}{
		"pageStart":    1,
		"pageSize":     500,
		"label":        "true",
		"resourceType": "PolarDB-X",
		"status":       "Available",
		"groupFiled":   "NodeClass",
	}
	if v, ok := d.GetOk("engine_version"); ok {
		reqQuery["engineVersion"] = v
	}
	if v, ok := d.GetOk("cpu_type"); ok {
		reqQuery["cpuType"] = v
	}

	if v, ok := d.GetOk("series"); ok {
		reqQuery["series"] = v
	}
	if v, ok := d.GetOk("cpu"); ok {
		reqQuery["cpu"] = v
	}
	if v, ok := d.GetOk("memory"); ok {
		reqQuery["memory"] = v
	}
	if v, ok := d.GetOk("spec_type"); ok {
		reqQuery["specType"] = v
	}

	reqHeader := map[string]string{
		"x-acs-territory": "US",
		"x-acs-lang":      "EN",
	}
	response, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "SelectCommonSpec", "/ascm/manage/saleconf/commonSpec/group", reqHeader, reqQuery, nil)
	if err != nil {
		return err
	}

	for _, d := range response["data"].([]interface{}) {
		data := d.(map[string]interface{})
		id := data["nodeClass"].(string)
		if _, exists := filterIds[id]; len(filterIds) > 0 && !exists {
			continue
		}
		if _, exists := existedId[id]; exists {
			continue
		}
		cpu_memory, err := ExtractNumbers(data["nodeClassLabel"].(string))
		if err != nil {
			return err
		}
		var cpu, memory int
		if len(cpu_memory) >= 2 {
			cpu = cpu_memory[0]
			memory = cpu_memory[1]
		}
		types = append(types, map[string]interface{}{
			"id":     id,
			"cpu":    cpu,
			"memory": memory,
		})
		existedId[id] = ""
		ids = append(ids, id)
	}

	sortedBy := d.Get("sorted_by").(string)
	if sortedBy != "" {
		sort.SliceStable(types, func(i, j int) bool {
			switch sortedBy {
			case "CPU":
				return types[i]["cpu"].(int) < types[j]["cpu"].(int)
			case "Memory":
				return types[i]["memory"].(int) < types[j]["memory"].(int)
			}
			return false
		})
	}

	d.Set("ids", ids)
	d.Set("instance_types", types)
	d.SetId(dataResourceIdHash(ids))

	return nil
}

func ExtractNumbers(input string) ([]int, error) {
	// Define a regular expression pattern to match one or more digits
	re := regexp.MustCompile(`\d+`)

	// Find all matching strings
	matches := re.FindAllString(input, -1)

	// Convert strings to integers
	var numbers []int
	for _, match := range matches {
		num, err := strconv.Atoi(match)
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, num)
	}

	return numbers, nil
}
