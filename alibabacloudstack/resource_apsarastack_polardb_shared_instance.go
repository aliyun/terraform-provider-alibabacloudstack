package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackPolardbSharedInstance() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"pay_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tde_status": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"department": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"multi_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"creation_option": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"db_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cpu_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"db_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sub_category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"db_node_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_node_num": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"proxy_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"storage_space": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"db_instance_net_type": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_cluster_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sub_domain": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_id_slave2": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_length": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ipv6_cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"deletion_lock": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_pay_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"proxy_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_cluster_network_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_latest_version": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"storage_max": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"db_nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_node_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"added_cpu_cores": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_node_role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"imci_switch": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_iops": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"db_node_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scc_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"failover_priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"hot_replica_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_weight": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"host_arch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone_ids": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintain_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"blktag_total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"architecture": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_cluster_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"inode_total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"db_version_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sql_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"proxy_cpu_cores": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_ips": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_cluster_ip_array_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl_enabled": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encrypt_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackPolardbSharedInstanceCreate, resourceAlibabacloudStackPolardbSharedInstanceRead, resourceAlibabacloudStackPolardbSharedInstanceUpdate, resourceAlibabacloudStackPolardbSharedInstanceDelete)
	return resource
}

func resourceAlibabacloudStackPolardbSharedInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolarDBService{client}

	request := make(map[string]interface{})

	// Required and optional parameters from schema
	if v, ok := d.GetOk("pay_type"); ok {
		request["PayType"] = v.(string)
	}
	if v, ok := d.GetOkExists("tde_status"); ok {
		request["TDEStatus"] = v.(bool)
	}
	if v, ok := d.GetOk("region_id"); ok {
		request["RegionId"] = v.(string)
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v.(string)
	}
	if v, ok := d.GetOk("creation_option"); ok {
		request["CreationOption"] = v.(string)
	}
	if v, ok := d.GetOk("db_type"); ok {
		request["DBType"] = v.(string)
	}
	if v, ok := d.GetOk("db_version"); ok {
		request["DBVersion"] = v.(string)
	}
	if v, ok := d.GetOk("db_node_class"); ok {
		request["DBNodeClass"] = v.(string)
	}
	if v, ok := d.GetOk("db_node_num"); ok {
		request["DBNodeNum"] = v.(string)
	}
	if v, ok := d.GetOk("proxy_type"); ok {
		request["ProxyType"] = v.(string)
	}
	if v, ok := d.GetOk("storage_type"); ok {
		request["StorageType"] = v.(string)
	}
	if v, ok := d.GetOk("storage_space"); ok {
		request["StorageSpace"] = v.(int)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VPCId"] = v.(string)
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v.(string)
	}
	if v, ok := d.GetOk("db_cluster_description"); ok {
		request["DBClusterDescription"] = v.(string)
	}

	// Additional parameters from the example that may not be in schema but required by API
	if v, ok := d.GetOk("sub_category"); ok {
		request["SubCategory"] = v.(string)
	}
	if v, ok := d.GetOk("cpu_type"); ok {
		request["CpuType"] = v.(string)
	}
	if v, ok := d.GetOk("multi_zone"); ok {
		request["MultiZone"] = v.(string)
	}
	if v, ok := d.GetOk("department"); ok {
		request["Department"] = v.(int)
	}
	if v, ok := d.GetOk("resource_group"); ok {
		request["ResourceGroup"] = v.(int)
	}
	if v, ok := d.GetOk("db_instance_net_type"); ok {
		request["DBInstanceNetType"] = v.(int)
	}
	if v, ok := d.GetOk("sub_domain"); ok {
		request["SubDomain"] = v.(string)
	}
	if v, ok := d.GetOk("zone_id_slave2"); ok {
		request["ZoneIdSlave2"] = v.(string)
	}
	if v, ok := d.GetOk("vswitch_zone_id"); ok {
		request["VswitchZoneId"] = v.(string)
	}
	if v, ok := d.GetOk("zone_length"); ok {
		request["ZoneLength"] = v.(int)
	}
	if v, ok := d.GetOk("cidr_block"); ok {
		request["CidrBlock"] = v.(string)
	}
	if v, ok := d.GetOk("ipv6_cidr_block"); ok {
		request["Ipv6CidrBlock"] = v.(string)
	}

	// Call CreateDBCluster API
	resp, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "CreateDBCluster", "", nil, request, nil)
	if err != nil {
		return err
	}

	// Extract DBClusterId from response
	dbClusterId, ok := resp["DBClusterId"].(string)
	if !ok {
		return fmt.Errorf("failed to get DBClusterId from CreateDBCluster response")
	}

	// Set temporary ID
	d.SetId(dbClusterId)

	// Wait for the cluster to be in Running state
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, polardbService.PolardbSharedInstanceStateRefreshFunc(dbClusterId, []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for PolarDB shared instance %s to be Running failed: %v", dbClusterId, err)
	}

	// Set the final ID
	d.SetId(dbClusterId)

	return nil
}

func resourceAlibabacloudStackPolardbSharedInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolarDBService{client}

	object, err := polardbService.DescribePolarDBSharedInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("db_cluster_id", object["DBClusterId"])
	d.Set("deletion_lock", object["DeletionLock"])
	d.Set("category", object["Category"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("storage_pay_type", object["StoragePayType"])
	d.Set("db_type", object["DBType"])
	d.Set("db_cluster_network_type", object["DBClusterNetworkType"])
	d.Set("is_latest_version", object["IsLatestVersion"])
	d.Set("storage_max", object["StorageMax"])
	d.Set("db_version", object["DBVersion"])

	if v, ok := object["DBNodes"].([]interface{}); ok {
		dbNodes := make([]map[string]interface{}, 0)
		for _, item := range v {
			node := item.(map[string]interface{})
			dbNode := map[string]interface{}{
				"db_node_status":    node["DBNodeStatus"],
				"zone_id":           node["ZoneId"],
				"max_connections":   node["MaxConnections"],
				"added_cpu_cores":   node["AddedCpuCores"],
				"db_node_role":      node["DBNodeRole"],
				"imci_switch":       node["ImciSwitch"],
				"db_node_id":        node["DBNodeId"],
				"max_iops":          node["MaxIOPS"],
				"db_node_class":     node["DBNodeClass"],
				"creation_time":     node["CreationTime"],
				"scc_mode":          node["SccMode"],
				"failover_priority": node["FailoverPriority"],
				"hot_replica_mode":  node["HotReplicaMode"],
				"server_weight":     node["ServerWeight"],
			}
			dbNodes = append(dbNodes, dbNode)
		}
		d.Set("db_nodes", dbNodes)
	}

	d.Set("zone_ids", object["ZoneIds"])
	d.Set("host_arch", object["HostArch"])
	d.Set("maintain_time", object["MaintainTime"])
	d.Set("engine", object["Engine"])

	// Handle tags if needed
	if v, ok := object["Tags"].([]interface{}); ok && len(v) > 0 {
		tags := make([]string, 0)
		for _, tag := range v {
			tagMap := tag.(map[string]interface{})
			tags = append(tags, fmt.Sprintf("%s:%s", tagMap["Key"], tagMap["Value"]))
		}
		d.Set("tags", tags)
	}

	d.Set("blktag_total", object["BlktagTotal"])
	d.Set("storage_type", object["StorageType"])
	d.Set("architecture", object["Architecture"])
	d.Set("vpc_id", object["VPCId"])
	d.Set("db_cluster_status", object["DBClusterStatus"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("db_cluster_description", object["DBClusterDescription"])
	d.Set("proxy_cpu_cores", object["ProxyCpuCores"])
	d.Set("pay_type", object["PayType"])
	d.Set("lock_mode", object["LockMode"])
	d.Set("storage_used", object["StorageUsed"])
	d.Set("inode_total", object["InodeTotal"])
	d.Set("storage_space", object["StorageSpace"])
	d.Set("db_version_status", object["DBVersionStatus"])
	d.Set("creation_time", object["CreationTime"])
	d.Set("sub_category", object["SubCategory"])
	d.Set("sql_size", object["SQLSize"])
	d.Set("region_id", object["RegionId"])
	d.Set("proxy_type", object["ProxyType"])
	d.Set("expire_time", object["ExpireTime"])
	d.Set("vip", object["Vip"])

	// Get security IPs from DescribeDBClusterAccessWhitelist
	reqQuery := map[string]interface{}{"DBClusterId": d.Id()}
	response, err := client.DoTeaRequest("GET", "polardb", "2017-08-01", "DescribeDBClusterAccessWhitelist", "", nil, reqQuery, nil)
	if err != nil {
		return err
	}

	if items, ok := response["Items"].([]interface{}); ok && len(items) > 0 {
		for _, item := range items {
			ipArray := item.(map[string]interface{})
			if ipArray["DBClusterIPArrayName"] == d.Get("db_cluster_ip_array_name") ||
				(d.Get("db_cluster_ip_array_name").(string) == "" && ipArray["DBClusterIPArrayName"] == "default") {
				d.Set("security_ips", ipArray["SecurityIps"])
				d.Set("db_cluster_ip_array_name", ipArray["DBClusterIPArrayName"])
				break
			}
		}
	}

	return nil
}

func resourceAlibabacloudStackPolardbSharedInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolarDBService{client}

	if d.IsNewResource() {
		return resourceAlibabacloudStackPolardbSharedInstanceRead(d, meta)
	}

	d.Partial(true)

	if d.HasChange("db_node_class") {
		// Get the current DBNodes to find the writer node
		object, err := polardbService.DescribePolarDBSharedInstance(d.Id())
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_polardb_shared_instance", "DescribeDBClusterAttribute", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		dbNodes := object["DBNodes"].([]interface{})
		var writerNode map[string]interface{}
		for _, node := range dbNodes {
			nodeMap := node.(map[string]interface{})
			if nodeMap["DBNodeRole"].(string) == "Writer" {
				writerNode = nodeMap
				break
			}
		}

		if writerNode == nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_polardb_shared_instance", "DescribeDBClusterAttribute", "Writer node not found")
		}

		oldClass, newClass := d.GetChange("db_node_class")
		modifyType := "Upgrade"
		if oldClass.(string) > newClass.(string) {
			modifyType = "Downgrade"
		}

		reqBody := map[string]interface{}{
			"DBClusterId": d.Id(),
			"ZoneId":      d.Get("zone_id"),
			"SubCategory": d.Get("sub_category"),
			"DBNode": map[string]interface{}{
				"TargetClass": newClass,
				"DBNodeId":    writerNode["DBNodeId"],
			},
			"ModifyType": modifyType,
		}

		if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "ModifyDBNodesClass", "", nil, nil, reqBody); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_shared_instance", "ModifyDBNodesClass", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{"ClassChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polardbService.PolardbSharedInstanceStateRefreshFunc(d.Id(), []string{"Failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		d.SetPartial("db_node_class")
	}

	if d.HasChange("deletion_lock") {
		reqQuery := map[string]interface{}{
			"DBClusterId": d.Id(),
			"type":        "detail",
			"Protection":  d.Get("deletion_lock").(int) == 1,
		}

		if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "ModifyDBClusterDeletion", "", nil, reqQuery, nil); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_shared_instance", "ModifyDBClusterDeletion", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		d.SetPartial("deletion_lock")
	}

	if d.HasChange("db_cluster_description") {
		reqQuery := map[string]interface{}{
			"DBClusterId":          d.Id(),
			"DBClusterDescription": d.Get("db_cluster_description"),
		}

		if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "ModifyDBClusterDescription", "", nil, reqQuery, nil); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_shared_instance", "ModifyDBClusterDescription", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		d.SetPartial("db_cluster_description")
	}

	if d.HasChange("security_ips") || d.HasChange("db_cluster_ip_array_name") {
		reqQuery := map[string]interface{}{
			"DBClusterId":          d.Id(),
			"SecurityIps":          d.Get("security_ips"),
			"DBClusterIPArrayName": d.Get("db_cluster_ip_array_name"),
		}

		if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "ModifyDBClusterAccessWhiteList", "", nil, reqQuery, nil); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_shared_instance", "ModifyDBClusterAccessWhiteList", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		d.SetPartial("security_ips")
		d.SetPartial("db_cluster_ip_array_name")
	}

	if d.HasChange("ssl_enabled") {
		// Get endpoint ID
		object, err := polardbService.DescribePolarDBSharedInstance(d.Id())
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_polardb_shared_instance", "DescribeDBClusterAttribute", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		endpointId := ""
		if endpoints, ok := object["Endpoints"].([]interface{}); ok && len(endpoints) > 0 {
			endpointId = endpoints[0].(map[string]interface{})["DBEndpointId"].(string)
		} else {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_polardb_shared_instance", "DescribeDBClusterAttribute", "Endpoint not found")
		}

		reqQuery := map[string]interface{}{
			"RegionId":     d.Get("region_id"),
			"DBClusterId":  d.Id(),
			"SSLEnabled":   d.Get("ssl_enabled"),
			"DBEndpointId": endpointId,
			"NetType":      "Private",
		}

		if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "ModifyDBClusterSSL", "", nil, reqQuery, nil); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_shared_instance", "ModifyDBClusterSSL", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{"SSLModifying"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polardbService.PolardbSharedInstanceStateRefreshFunc(d.Id(), []string{"Failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		d.SetPartial("ssl_enabled")
	}

	if d.HasChange("tde_status") {
		reqQuery := map[string]interface{}{
			"RegionId":         d.Get("region_id"),
			"DBClusterId":      d.Id(),
			"TDEStatus":        "Enable",
			"RoleArn":          d.Get("role_arn"),
			"EncryptionKey":    d.Get("encryption_key"),
			"EncryptAlgorithm": d.Get("encrypt_algorithm"),
		}

		if !d.Get("tde_status").(bool) {
			reqQuery["TDEStatus"] = "Disable"
		}

		if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "ModifyDBClusterTDE", "", nil, reqQuery, nil); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_shared_instance", "ModifyDBClusterTDE", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{"TDEModifying"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polardbService.PolardbSharedInstanceStateRefreshFunc(d.Id(), []string{"Failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		d.SetPartial("tde_status")
	}

	d.Partial(false)
	return resourceAlibabacloudStackPolardbSharedInstanceRead(d, meta)
}

func resourceAlibabacloudStackPolardbSharedInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolarDBService{client}

	reqQuery := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	// Call DeleteDBCluster API
	_, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "DeleteDBCluster", "", nil, reqQuery, nil)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound"}) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteDBCluster", errmsgs.AlibabacloudStackSdkGoERROR, "")
	}

	// Wait for the cluster to be fully deleted
	stateConf := BuildStateConf([]string{"Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, polardbService.PolardbSharedInstanceStateRefreshFunc(d.Id(), []string{}))
	_, err = stateConf.WaitForState()
	return errmsgs.WrapError(err)
}
