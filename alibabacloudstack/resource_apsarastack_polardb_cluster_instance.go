package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackPolardbClusterInstance() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"pay_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tde_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
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
			"db_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cpu_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"intel", "hygon"}, false),
			},
			"db_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sub_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"General", "Exclusive"}, false),
			},
			"db_node_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_read_node_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_node_num": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"proxy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"proxy_exclusive", "proxy_off"}, false),
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
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
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
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 1),
			},
			"category": {
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
			"zone_ids": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintain_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": caseInsensitiveTagsSchema(),
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
			"inode_total": {
				Type:     schema.TypeInt,
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
			"security_ips_groups": {
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				Optional: true,
				Computed: true,
				Set: func(v interface{}) int {
					ipsMap := v.(map[string]interface{})
					keys := make([]string, 0, len(ipsMap))
					for k := range ipsMap {
						keys = append(keys, k)
					}
					sort.Strings(keys)

					sortedMap := make(map[string]interface{})
					for _, k := range keys {
						value := ipsMap[k].(string)
						parts := strings.Split(value, ",")
						filteredParts := make([]string, 0)
						for _, part := range parts {
							trimmed := strings.TrimSpace(part)
							if trimmed != "" {
								filteredParts = append(filteredParts, trimmed)
							}
						}
						sort.Strings(filteredParts)
						sortedValue := strings.Join(filteredParts, ",")
						sortedMap[k] = sortedValue
					}

					bytes, _ := json.Marshal(sortedMap)
					return hashcode.String(string(bytes))
				},
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					o, n := d.GetChange("security_ips_groups")
					oldMap := o.(map[string]interface{})
					newMap := n.(map[string]interface{})

					if len(oldMap) != len(newMap) {
						return false
					}

					for k, v := range oldMap {
						nv, ok := newMap[k]
						if !ok {
							return false
						}
						ovs := strings.Split(v.(string), ",")
						filteredOvs := make([]string, 0)
						for _, part := range ovs {
							trimmed := strings.TrimSpace(part)
							if trimmed != "" {
								filteredOvs = append(filteredOvs, trimmed)
							}
						}
						sort.Strings(filteredOvs)

						nvs := strings.Split(nv.(string), ",")
						filteredNvs := make([]string, 0)
						for _, part := range nvs {
							trimmed := strings.TrimSpace(part)
							if trimmed != "" {
								filteredNvs = append(filteredNvs, trimmed)
							}
						}
						sort.Strings(filteredNvs)

						if strings.Join(filteredOvs, ",") != strings.Join(filteredNvs, ",") {
							return false
						}
					}
					return true
				},
			},
			"security_groups": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Computed: true,
				MaxItems: 10,
			},
			"ssl_enabled": {
				Type:     schema.TypeBool,
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
			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Optional: true,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackPolardbClusterInstanceCreate, resourceAlibabacloudStackPolardbClusterInstanceRead, resourceAlibabacloudStackPolardbClusterInstanceUpdate, resourceAlibabacloudStackPolardbClusterInstanceDelete)
	return resource
}

func resourceAlibabacloudStackPolardbClusterInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolardbService{client}

	request := make(map[string]interface{})
	request["PayType"] = "Postpaid"
	request["VPCId"] = d.Get("vpc_id").(string)
	request["VSwitchId"] = d.Get("vswitch_id").(string)
	request["DBVersion"] = d.Get("db_version").(string)
	request["ZoneId"] = d.Get("zone_id").(string)
	request["DBType"] = d.Get("db_type").(string)
	request["DBNodeClass"] = d.Get("db_node_class").(string)
	request["DBNodeNum"] = d.Get("db_node_num").(int)
	request["DBClusterDescription"] = d.Get("db_cluster_description").(string)
	// Required and optional parameters from schema
	if v := d.Get("tde_enabled"); v.(bool) {
		request["TDEStatus"] = true
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
	if v, ok := d.GetOk("sub_category"); ok {
		request["SubCategory"] = GetSubCategory(v)
	}
	if v, ok := d.GetOk("cpu_type"); ok {
		request["CpuType"] = v.(string)
	}
	if v, ok := d.GetOk("multi_zone"); ok {
		request["MultiZone"] = v.(string)
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
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 2*time.Minute, polardbService.PolardbClusterInstanceStateRefreshFunc(dbClusterId, []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for PolarDB shared instance %s to be Running failed: %v", dbClusterId, err)
	}
	return nil
}

func resourceAlibabacloudStackPolardbClusterInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolardbService{client}

	object, err := polardbService.DescribePolardbClusterInstance(d.Id())
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
	d.Set("db_type", object["DBType"])
	d.Set("db_cluster_network_type", object["DBClusterNetworkType"])
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
	d.Set("maintain_time", object["MaintainTime"])
	d.Set("engine", object["Engine"])

	// Handle tags if needed
	tags := make(map[string]interface{}, 0)
	if v, ok := object["Tags"].([]interface{}); ok && len(v) > 0 {
		for _, tag := range v {
			tagMap := tag.(map[string]interface{})
			tags[tagMap["Key"].(string)] = tagMap["Value"]
		}
	}
	d.Set("tags", tags)
	d.Set("blktag_total", object["BlktagTotal"])
	d.Set("storage_type", strings.ToUpper(object["StorageType"].(string)))
	d.Set("architecture", object["Architecture"])
	d.Set("db_cluster_status", object["DBClusterStatus"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("db_cluster_description", object["DBClusterDescription"])
	d.Set("proxy_cpu_cores", object["ProxyCpuCores"])
	d.Set("pay_type", object["PayType"])
	d.Set("lock_mode", object["LockMode"])
	d.Set("inode_total", object["InodeTotal"])
	storageSpace := object["StorageSpace"].(json.Number)
	if spaceValue, err := storageSpace.Int64(); err == nil {
		storage_space := spaceValue / 1024 / 1024 / 1024
		d.Set("storage_space", int(storage_space))
	}
	d.Set("creation_time", object["CreationTime"])
	// d.Set("sub_category", object["SubCategory"])
	d.Set("sql_size", object["SQLSize"])
	d.Set("proxy_type", object["ProxyType"])
	d.Set("vip", object["Vip"])

	// Get security IPs from DescribeDBClusterAccessWhitelist
	reqQuery := map[string]interface{}{"DBClusterId": d.Id()}
	response, err := client.DoTeaRequest("GET", "polardb", "2017-08-01", "DescribeDBClusterAccessWhiteList", "", nil, reqQuery, nil)
	if err != nil {
		return err
	}
	security_ips_groups := make(map[string]string)
	security_ips := response["Items"].(map[string]interface{})["DBClusterIPArray"].([]interface{})
	if len(security_ips) > 0 {
		for _, item := range security_ips {
			ipArray := item.(map[string]interface{})
			if ipArray["DBClusterIPArrayName"].(string) == "default" {
				continue
			}
			security_ips_groups[ipArray["DBClusterIPArrayName"].(string)] = ipArray["SecurityIps"].(string)
		}
	}
	d.Set("security_ips_groups", security_ips_groups)
	security_groups := response["DBClusterSecurityGroups"].(map[string]interface{})["DBClusterSecurityGroup"]
	d.Set("security_groups", security_groups)
	if err = polardbService.RefreshClusterParameters(d); err != nil {
		return errmsgs.WrapError(err)
	}
	tdeStatus, encryptionKey, err := polardbService.DescribeDBClusterTDE(d.Id())
	if err != nil {
		return err
	}
	d.Set("tde_enabled", tdeStatus == "Enabled")
	d.Set("encryption_key", encryptionKey)
	ssl_enabled, err := polardbService.DescribeDBClusterSSL(d.Id())
	if err != nil {
		return err
	}
	d.Set("ssl_enabled", ssl_enabled)
	return nil
}

func resourceAlibabacloudStackPolardbClusterInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolardbService{client}
	if d.HasChange("security_ips_groups") {
		old, new := d.GetChange("security_ips_groups")
		err := polardbService.ModifySecurityIps(d.Id(), old, new)
		if err != nil {
			return err
		}
	}

	if d.HasChange("security_groups") {
		security_groups := d.Get("security_groups").(*schema.Set).List()
		err := polardbService.ModifySecurityGroups(d.Id(), security_groups)
		if err != nil {
			return err
		}
	}

	if d.HasChange("ssl_enabled") && !(d.IsNewResource() && !d.Get("ssl_enabled").(bool)) {

		endpointId := ""
		if endpointsResponse, err := polardbService.DescribeDBClusterEndpoints(d.Id()); err == nil {
			endpoint, err := jsonpath.Get("$.Items.0.DBEndpointId", endpointsResponse)
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_polardb_cluster_instance", "DescribeDBClusterEndpoints", errmsgs.AlibabacloudStackSdkGoERROR)
			}
			endpointId = endpoint.(string)
		} else {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_polardb_cluster_instance", "DescribeDBClusterAttribute", "Endpoint not found")
		}
		enabled := d.Get("ssl_enabled").(bool)
		reqQuery := map[string]interface{}{
			"DBClusterId":  d.Id(),
			"DBEndpointId": endpointId,
			"NetType":      "Private",
		}
		if enabled {
			reqQuery["SSLEnabled"] = "Enable"
		} else {
			reqQuery["SSLEnabled"] = "Disable"
		}

		if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "ModifyDBClusterSSL", "", nil, reqQuery, nil); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_polardb_cluster_instance", "ModifyDBClusterSSL", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{"SSLModifying"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polardbService.PolardbClusterInstanceStateRefreshFunc(d.Id(), []string{"Failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	if d.HasChange("tde_enabled") {
		if o, n := d.GetChange("tde_enabled"); o.(bool) && !n.(bool) {
			return errmsgs.Error("TDE cannot be disabled after enabled, please disable TDE in console.")
		}
		role_arn, err := polardbService.CheckCloudResourceAuthorized()
		if err != nil {
			return err
		}
		reqQuery := map[string]interface{}{
			"DBClusterId": d.Id(),
			"TDEStatus":   "Enable",
			"RoleArn":     role_arn,
		}
		if v, ok := d.GetOk("encryption_key"); ok {
			reqQuery["EncryptionKey"] = v.(string)
		}
		if v, ok := d.GetOk("encrypt_algorithm"); ok {
			reqQuery["EncryptAlgorithm"] = v.(string)
		}
		if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "ModifyDBClusterTDE", "", nil, reqQuery, nil); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_cluster_instance", "ModifyDBClusterTDE", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{"TDEModifying"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polardbService.PolardbClusterInstanceStateRefreshFunc(d.Id(), []string{"Failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	if d.HasChange("parameters") {
		if err := polardbService.ModifyClusterParameters(d); err != nil {
			return err
		}
	}
	if d.HasChange("deletion_lock") {
		reqQuery := map[string]interface{}{
			"DBClusterId": d.Id(),
			"type":        "detail",
			"Protection":  d.Get("deletion_lock").(int) == 1,
		}

		if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "ModifyDBClusterDeletion", "", nil, reqQuery, nil); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_cluster_instance", "ModifyDBClusterDeletion", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	if d.IsNewResource() {
		return resourceAlibabacloudStackPolardbClusterInstanceRead(d, meta)
	}
	if d.HasChange("db_node_class") {
		// Get the current DBNodes to find the writer node
		var oldClassid, newClassid string
		newClassid = d.Get("db_node_class").(string)

		nodes := d.Get("db_nodes").([]interface{})
		for _, n := range nodes {
			node := n.(map[string]interface{})
			if node["db_node_role"].(string) == "Writer" {
				oldClassid = node["db_node_class"].(string)
			}
		}
		newClass, err := polardbService.GetPolardbClusterClassData(newClassid)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		oldClass, err := polardbService.GetPolardbClusterClassData(oldClassid)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		modifyType := GetPolardbClusterInstanceModifyType(oldClass, newClass)

		dbNodeDatas := make([]map[string]interface{}, 0)
		reqBody := map[string]interface{}{
			"DBClusterId": d.Id(),
			"ZoneId":      d.Get("zone_id"),
			"SubCategory": GetSubCategory(d.Get("sub_category")),
			"ModifyType":  modifyType,
		}
		for _, n := range nodes {
			node := n.(map[string]interface{})
			dbNodeDatas = append(dbNodeDatas, map[string]interface{}{
				"TargetClass": newClassid,
				"DBNodeId":    node["db_node_id"].(string),
			})
		}
		reqBody["DBNode"] = dbNodeDatas
		if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "ModifyDBNodesClass", "", nil, nil, reqBody); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_cluster_instance", "ModifyDBNodesClass", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{"ClassChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polardbService.PolardbClusterInstanceStateRefreshFunc(d.Id(), []string{"Failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	if d.HasChange("db_cluster_description") {
		reqQuery := map[string]interface{}{
			"DBClusterId":          d.Id(),
			"DBClusterDescription": d.Get("db_cluster_description"),
		}

		if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "ModifyDBClusterDescription", "", nil, reqQuery, nil); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_cluster_instance", "ModifyDBClusterDescription", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}
	if d.HasChange("db_node_num") {
		o, n := d.GetChange("db_node_num")
		add := n.(int) - o.(int)
		action := "CreateDBNodes"
		reqQuery := map[string]interface{}{
			"DBClusterId": d.Id(),
		}
		if add < 0 {
			action = "DeleteDBNodes"
			standbyNodes := make([]string, 0)
			for _, v := range d.Get("db_nodes").([]interface{}) {
				node := v.(map[string]interface{})
				if node["db_node_role"].(string) == "Reader" {
					standbyNodes = append(standbyNodes, node["db_node_id"].(string))
				}
			}
			if len(standbyNodes) < add*-1 {
				return errmsgs.WrapError(errmsgs.Error("The number of standby nodes cannot be less than the number of nodes to be deleted."))
			}
			reqQuery["DBNodeId"] = standbyNodes[:add*-1]
		} else {
			count := add
			for count > 0 {
				dbNodeZones := make([]string, 0)
				dbNodeZones = append(dbNodeZones, d.Get("zone_id").(string))
				reqQuery["DBNode"] = dbNodeZones
			}
		}

		if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", action, "", nil, reqQuery, nil); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_cluster_instance", action, errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}
	return resourceAlibabacloudStackPolardbClusterInstanceRead(d, meta)
}

func resourceAlibabacloudStackPolardbClusterInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
	// client := meta.(*connectivity.AlibabacloudStackClient)
	// polardbService := PolardbService{client}

	// reqQuery := map[string]interface{}{
	// 	"DBClusterId": d.Id(),
	// }

	// // Call DeleteDBCluster API
	// _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "DeleteDBCluster", "", nil, reqQuery, nil)
	// if err != nil {
	// 	if errmsgs.IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound"}) {
	// 		return nil
	// 	}
	// 	return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteDBCluster", errmsgs.AlibabacloudStackSdkGoERROR, "")
	// }

	// // Wait for the cluster to be fully deleted
	// stateConf := BuildStateConf([]string{"Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, polardbService.PolardbClusterInstanceStateRefreshFunc(d.Id(), []string{}))
	// _, err = stateConf.WaitForState()
	// return errmsgs.WrapError(err)
}

func GetPolardbClusterInstanceModifyType(oldClass, newClass map[string]interface{}) string {
	var modifyType = "Upgrade"
	oldcpu, _ := strconv.Atoi(oldClass["cpu"].(string))
	newcpu, _ := strconv.Atoi(newClass["cpu"].(string))
	oldmemory, _ := strconv.Atoi(oldClass["memory"].(string))
	newmemory, _ := strconv.Atoi(newClass["memory"].(string))
	if newcpu < oldcpu || (newcpu == oldcpu && newmemory < oldmemory) {
		modifyType = "Downgrade"
	}
	return modifyType
}

func GetSubCategory(subCategory interface{}) string {
	switch subCategory.(string) {
	case "General":
		return "normal_general"
	case "Exclusive":
		return "normal_exclusive"
	default:
		return ""
	}
}
