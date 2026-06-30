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

func dataSourceAlibabacloudStackBmcpNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackBmcpNodesRead,

		Schema: map[string]*schema.Schema{
			"node_name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"node_id_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_id_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"sn_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_ip_regex": {
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
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"machine_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"machine_type_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_arch": {
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
						"gpu_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gpu_model": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackBmcpNodesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := make(map[string]interface{})
	reqQuery["PageNumber"] = 1
	reqQuery["PageSize"] = PageSizeLarge

	// Call ListBmcpInstance API
	resp, err := client.DoTeaRequest("POST", "bms", "2024-03-01", "ListBmcpInstance", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "ListBmcpInstance", "POST", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Parse response data
	data, err := jmespath.Search("data.data", resp)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "ListBmcpInstance", "data.data", resp)
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

		nodeName := formatAnyToString(itemMap["HostName"])
		nodeId := formatAnyToString(itemMap["ID"])
		clusterId := formatAnyToString(itemMap["ClusterId"])
		clusterName := formatAnyToString(itemMap["ClusterName"])
		sn := formatAnyToString(itemMap["SN"])
		vpcIp := formatAnyToString(itemMap["IP"])

		// Apply node_name_regex filter
		if nodeNameRegex, ok := d.GetOk("node_name_regex"); ok {
			if nodeNameRegex.(string) == "" {
				return errmsgs.WrapError(fmt.Errorf("node_name_regex is set but empty value, this is the second query with filter"))
			}
			if !strings.Contains(nodeName, nodeNameRegex.(string)) {
				continue
			}
		}

		// Apply node_id_regex filter
		if nodeIdRegex, ok := d.GetOk("node_id_regex"); ok {
			if !strings.Contains(nodeId, nodeIdRegex.(string)) {
				continue
			}
		}

		// Apply cluster_id_regex filter
		if clusterIdRegex, ok := d.GetOk("cluster_id_regex"); ok {
			if !strings.Contains(clusterId, clusterIdRegex.(string)) {
				continue
			}
		}

		// Apply cluster_name_regex filter
		if clusterNameRegex, ok := d.GetOk("cluster_name_regex"); ok {
			if !strings.Contains(clusterName, clusterNameRegex.(string)) {
				continue
			}
		}

		// Apply sn_regex filter
		if snRegex, ok := d.GetOk("sn_regex"); ok {
			if !strings.Contains(sn, snRegex.(string)) {
				continue
			}
		}

		// Apply vpc_ip_regex filter
		if vpcIpRegex, ok := d.GetOk("vpc_ip_regex"); ok {
			if !strings.Contains(vpcIp, vpcIpRegex.(string)) {
				continue
			}
		}

		filteredItems = append(filteredItems, item)
	}

	return bmcpNodesAttributes(d, filteredItems)
}

func bmcpNodesAttributes(d *schema.ResourceData, items []interface{}) error {
	var ids []string
	var s []map[string]interface{}

	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		id := formatAnyToString(itemMap["ID"])
		nodeId := id
		nodeName := formatAnyToString(itemMap["HostName"])
		clusterId := formatAnyToString(itemMap["ClusterId"])
		clusterName := formatAnyToString(itemMap["ClusterName"])
		instanceAlias := formatAnyToString(itemMap["InstanceAlias"])
		instanceName := formatAnyToString(itemMap["InstanceName"])
		sn := formatAnyToString(itemMap["SN"])
		vpcIp := formatAnyToString(itemMap["IP"])
		status := formatAnyToString(itemMap["Status"])
		machineType := formatAnyToString(itemMap["MachineType"])
		machineTypeName := machineType
		cpuArch := formatAnyToString(itemMap["CPUArch"])
		gpuModel := formatAnyToString(itemMap["Gpu"])
		regionId := formatAnyToString(itemMap["RegionId"])
		createTime := formatAnyToString(itemMap["CreateTime"])

		var cpuNumber, memory, disk int
		if config, ok := itemMap["Config"].(map[string]interface{}); ok {
			cpuNumber = formatAnyToInt(config["CPUNumber"])
			memory = formatAnyToInt(config["Memory"])
			disk = formatAnyToInt(config["Disk"])
		}

		gpuNum := formatAnyToInt(itemMap["GpuNum"])

		mapping := map[string]interface{}{
			"id":                id,
			"node_id":           nodeId,
			"node_name":         nodeName,
			"cluster_id":        clusterId,
			"cluster_name":      clusterName,
			"instance_alias":    instanceAlias,
			"instance_name":     instanceName,
			"sn":                sn,
			"vpc_ip":            vpcIp,
			"status":            status,
			"machine_type":      machineType,
			"machine_type_name": machineTypeName,
			"cpu_arch":          cpuArch,
			"cpu_number":        cpuNumber,
			"memory":            memory,
			"disk":              disk,
			"gpu_num":           gpuNum,
			"gpu_model":         gpuModel,
			"region_id":         regionId,
			"create_time":       createTime,
		}

		ids = append(ids, id)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("nodes", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func formatAnyToString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case json.Number:
		return val.String()
	case float64:
		return fmt.Sprintf("%v", val)
	case int:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%d", val)
	case bool:
		return fmt.Sprintf("%t", val)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func formatAnyToInt(v interface{}) int {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	case json.Number:
		i, _ := val.Int64()
		return int(i)
	case string:
		i, _ := strconv.Atoi(val)
		return i
	default:
		return 0
	}
}
