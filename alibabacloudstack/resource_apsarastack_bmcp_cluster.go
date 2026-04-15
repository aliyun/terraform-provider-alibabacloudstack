package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"
)

func resourceAlibabacloudStackBmcpCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackBmcpClusterCreate,
		Read:   resourceAlibabacloudStackBmcpClusterRead,
		Delete: resourceAlibabacloudStackBmcpClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"standard_vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"evpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"machine_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_count": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"activation_code": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"switch_method": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "CHSW",
			},
			"cluster_arch_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "standard",
			},
			"is_install_yundun_aegis": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"is_create_cpfs_cluster": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			// Computed values.
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mem_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"gpu": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gpu_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gpu_model": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"flops_count": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"video_memory": {
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
	}
}

func resourceAlibabacloudStackBmcpClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	ecsService := &EcsService{client: client}
	department, err := strconv.Atoi(client.Department)
	if err != nil {
		return fmt.Errorf("failed to convert department to int: %s", err)
	}
	activationId, activationCode, err := ecsService.CreateActivation(1000, department, 4, client.RegionId)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "CreateActivation", "POST", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Query machine type info from API
	machineType := d.Get("machine_type").(string)
	nodeCount := d.Get("node_count").(int)
	vswitchId := d.Get("vswitch_id").(string)

	machineTypeInfo, err := queryMachineTypeInfo(client, machineType)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "QueryMachineType", "POST", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Build MachineTypeList with queried info
	machineTypeList := []map[string]interface{}{
		{
			"MachineType":    machineType,
			"NodeCount":      nodeCount,
			"VswitchId":      vswitchId,
			"Arch":           machineTypeInfo["cpu_arch"],
			"Gpu":            machineTypeInfo["gpu"],
			"CpuCount":       machineTypeInfo["cpu_number"],
			"MemCount":       machineTypeInfo["memory"],
			"GpuNum":         machineTypeInfo["gpu_num"],
			"VideoMemCount":  machineTypeInfo["video_memory"],
			"FlopsCount":     machineTypeInfo["tflops_fp32"],
			"Image":          "",
			"ImageType":      "",
			"UseOriginImage": true,
			"Specification":  machineTypeInfo["specification"],
		},
	}
	machineTypeListjson, err := json.Marshal(machineTypeList)
	if err != nil {
		return fmt.Errorf("failed to marshal MachineTypeList to JSON: %s", err)
	}

	// Prepare the query parameters (simple parameters)
	reqQuery := map[string]interface{}{
		"ZoneId":               d.Get("zone_id").(string),
		"ClusterName":          d.Get("cluster_name").(string),
		"StandardVswitchId":    d.Get("standard_vswitch_id").(string),
		"Password":             d.Get("password").(string),
		"VpcId":                d.Get("vpc_id").(string),
		"EvpcId":               d.Get("evpc_id").(string),
		"SwitchMethod":         d.Get("switch_method").(string),
		"ClusterArchType":      d.Get("cluster_arch_type").(string),
		"IsInstallYundunAegis": fmt.Sprintf("%v", d.Get("is_install_yundun_aegis").(bool)),
		"IsCreateCpfsCluster":  fmt.Sprintf("%v", d.Get("is_create_cpfs_cluster").(bool)),
		"EnableIpv6":           fmt.Sprintf("%v", d.Get("enable_ipv6").(bool)),
		"ActivationId":         activationId,
		"ActivationCode":       activationCode,
		"MachineTypeList":      string(machineTypeListjson),
	}

	if activationCode, ok := d.GetOk("activation_code"); ok {
		reqQuery["ActivationCode"] = activationCode.(string)
	}

	// Call the API to create the cluster
	resp, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "CreateBmcpCluster", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "CreateBmcpCluster", "POST", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Extract cluster ID from response
	clusterId, err := jmespath.Search("data.ClusterId", resp)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "CreateBmcpCluster", "data.ClusterId", resp)
	}

	if clusterId == nil {
		return fmt.Errorf("failed to retrieve ClusterId from CreateBmcpCluster response")
	}

	clusterIdStr := clusterId.(string)
	d.SetId(clusterIdStr)

	// Wait for the cluster to be active
	if err := waitForBmcpClusterActive(client, clusterIdStr, d.Timeout(schema.TimeoutCreate)); err != nil {
		return fmt.Errorf("error waiting for BMCP cluster (%s) to become active: %s", clusterIdStr, err)
	}

	return resourceAlibabacloudStackBmcpClusterRead(d, meta)
}

func resourceAlibabacloudStackBmcpClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	clusterId := d.Id()

	reqQuery := map[string]interface{}{
		"PageNumber": 1,
		"PageSize":   1,
		"ClusterId":  clusterId,
	}

	resp, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ListBmcpCluster", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "ListBmcpCluster", "GET", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Parse response
	clusterList, err := jmespath.Search("ClusterList", resp)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "ListBmcpCluster", "ClusterList", resp)
	}

	clusterListSlice, ok := clusterList.([]interface{})
	if !ok || len(clusterListSlice) == 0 {
		log.Printf("[DEBUG] BMCP cluster %s not found", clusterId)
		d.SetId("")
		return nil
	}

	cluster := clusterListSlice[0].(map[string]interface{})

	d.Set("cluster_id", cluster["ClusterId"])
	d.Set("cluster_name", cluster["ClusterName"])
	d.Set("status", cluster["Status"])
	d.Set("region_id", cluster["RegionId"])
	d.Set("vpc_id", cluster["VpcId"])
	d.Set("evpc_id", cluster["EvpcId"])
	d.Set("zone_id", cluster["Zone"])
	d.Set("node_count", cluster["NodeCount"])
	d.Set("cpu_count", formatAnyToInt(cluster["CpuCount"]))
	d.Set("mem_count", formatAnyToInt(cluster["MemCount"]))
	d.Set("flops_count", formatAnyToFloat(cluster["FlopsCount"]))
	d.Set("video_memory", formatAnyToInt(cluster["VideoMemory"]))
	d.Set("create_time", cluster["CreateTime"])
	d.Set("update_time", cluster["UpdateTime"])
	d.Set("enable_ipv6", cluster["EnableIpv6"])
	d.Set("switch_method", cluster["SwitchMethod"])

	// Set GPU list
	if gpuList, ok := cluster["Gpu"].([]interface{}); ok {
		gpuListData := make([]map[string]interface{}, 0, len(gpuList))
		for _, gpu := range gpuList {
			gpuMap := gpu.(map[string]interface{})
			gpuListData = append(gpuListData, map[string]interface{}{
				"gpu_num":   gpuMap["GpuNum"],
				"gpu_model": gpuMap["GpuModel"],
			})
		}
		d.Set("gpu", gpuListData)
	}

	return nil
}

func resourceAlibabacloudStackBmcpClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	clusterId := d.Id()
	clusterName := d.Get("cluster_name").(string)

	reqBody := map[string]interface{}{
		"ClusterId":       clusterId,
		"ClusterName":     clusterName,
		"RetentionPeriod": "168",
	}

	_, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "DeleteBmcpCluster", "", nil, reqBody, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "DeleteBmcpCluster", "POST", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Wait for the cluster to be deleted
	if err := waitForBmcpClusterDeleted(client, clusterId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return fmt.Errorf("error waiting for BMCP cluster (%s) to be deleted: %s", clusterId, err)
	}

	return nil
}

func waitForBmcpClusterActive(client *connectivity.AlibabacloudStackClient, clusterId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"creating", "initializing", "configuring"},
		Target:  []string{"active"},
		Refresh: func() (interface{}, string, error) {
			reqQuery := map[string]interface{}{
				"PageNumber": 1,
				"PageSize":   1,
				"ClusterId":  clusterId,
			}

			resp, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ListBmcpCluster", "", nil, reqQuery, nil)
			if err != nil {
				return nil, "", err
			}

			clusterList, err := jmespath.Search("ClusterList", resp)
			if err != nil {
				return nil, "", err
			}

			clusterListSlice, ok := clusterList.([]interface{})
			if !ok || len(clusterListSlice) == 0 {
				return nil, "creating", nil
			}

			cluster := clusterListSlice[0].(map[string]interface{})
			status := cluster["Status"].(string)

			return cluster, status, nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 30 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func waitForBmcpClusterDeleted(client *connectivity.AlibabacloudStackClient, clusterId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"active", "deleting"},
		Target:  []string{"deleted"},
		Refresh: func() (interface{}, string, error) {
			reqQuery := map[string]interface{}{
				"PageNumber": 1,
				"PageSize":   1,
				"ClusterId":  clusterId,
			}

			resp, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ListBmcpCluster", "", nil, reqQuery, nil)
			if err != nil {
				return nil, "deleting", nil
			}

			clusterList, err := jmespath.Search("ClusterList", resp)
			if err != nil {
				return nil, "deleting", nil
			}

			clusterListSlice, ok := clusterList.([]interface{})
			if !ok || len(clusterListSlice) == 0 {
				return "deleted", "deleted", nil
			}

			cluster := clusterListSlice[0].(map[string]interface{})
			status := cluster["Status"].(string)

			return cluster, status, nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 30 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func formatAnyToFloat(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case json.Number:
		f, _ := val.Float64()
		return f
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	default:
		return 0
	}
}

// queryMachineTypeInfo queries machine type details from ListBmcpMachineType API
func queryMachineTypeInfo(client *connectivity.AlibabacloudStackClient, machineType string) (map[string]interface{}, error) {
	reqQuery := map[string]interface{}{
		"PageNumber":  1,
		"PageSize":    PageSizeLarge,
		"DeployTypes": []string{"bmcp"},
	}

	resp, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ListBmcpMachineType", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}

	data, err := jmespath.Search("data.MachineTypeList", resp)
	if err != nil {
		return nil, err
	}

	items, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format for MachineTypeList")
	}

	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		name, _ := itemMap["name"].(string)
		if name == machineType {
			return map[string]interface{}{
				"cpu_arch":      formatAnyToString(itemMap["CPUArch"]),
				"gpu":           formatAnyToString(itemMap["gpu"]),
				"cpu_number":    formatAnyToInt(itemMap["CPUNumber"]),
				"memory":        formatAnyToInt(itemMap["memory"]),
				"gpu_num":       formatAnyToInt(itemMap["GpuNum"]),
				"video_memory":  formatAnyToInt(itemMap["VideoMemory"]),
				"tflops_fp32":   formatAnyToInt(itemMap["TflopsFP32"]),
				"specification": formatAnyToString(itemMap["Specification"]),
			}, nil
		}
	}

	return nil, fmt.Errorf("machine type %s not found", machineType)
}
