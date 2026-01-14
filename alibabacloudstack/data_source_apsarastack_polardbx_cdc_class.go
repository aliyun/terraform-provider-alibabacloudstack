package alibabacloudstack

import (
	"sort"
	"strconv"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackPolardbxCdcClasses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackPolardbxCdcClassesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
			},
			"cpu": {
				Type:     schema.TypeInt,
				Optional: true,
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
			// Computed values.
			"cdc_classes": {
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

func dataSourceAlibabacloudStackPolardbxCdcClassesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	filterIds := getIdsStringFilter(d)

	ids := []string{}
	classes := []map[string]interface{}{}

	reqQuery := map[string]interface{}{
		"DBInstanceName":       d.Get("instance_id").(string),
	}
	response, err := client.DoTeaRequest("GET", "polardbx", "2020-02-02", "DescribeCdcClassList", "", nil, reqQuery, nil)
	if err != nil {
		if sdkErr, ok := err.(*tea.SDKError); ok && *sdkErr.Code == "DBInstance.NotFound" {
			d.Set("ids", ids)
			d.Set("cdc_classes", classes)
			d.SetId(dataResourceIdHash(ids))
			return nil
		}
		return err
	}

	for _, i := range response["Data"].(map[string]interface{})["ClassCodeList"].([]interface{}) {
		data := i.(map[string]interface{})
		id := data["ClassCode"].(string)
		if _, exists := filterIds[id]; len(filterIds) > 0 && !exists {
			continue
		}
		var cpu, memory int
		if v, err := strconv.Atoi(data["CpuCore"].(string));err != nil {
			return err
		} else {
			cpu = v
			if v, ok := d.GetOk("cpu"); ok && cpu != v.(int) {
				continue
			}
		}
		if v, err := strconv.Atoi(data["Mem"].(string));err != nil {
			return err
		} else {
			memory = v
			if v, ok := d.GetOk("memory"); ok && memory != v.(int) {
				continue
			}
		}
		classes = append(classes, map[string]interface{}{
			"id":     id,
			"cpu":    cpu,
			"memory": memory,
		})
		ids = append(ids, id)
	}

	sortedBy := d.Get("sorted_by").(string)
	if sortedBy != "" {
		sort.SliceStable(classes, func(i, j int) bool {
			switch sortedBy {
			case "CPU":
				return classes[i]["cpu"].(int) < classes[j]["cpu"].(int)
			case "Memory":
				return classes[i]["memory"].(int) < classes[j]["memory"].(int)
			}
			return false
		})
	}

	d.Set("ids", ids)
	d.Set("cdc_classes", classes)
	d.SetId(dataResourceIdHash(ids))

	return nil
}
