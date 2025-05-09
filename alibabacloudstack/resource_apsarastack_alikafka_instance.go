package alibabacloudstack

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var AlikafkaInstanceSpecMap = map[string]string{
	"Broker4C16G":   "KAFKA.2XUNIT",
	"Broker8C32G":   "KAFKA.6XUNIT",
	"Broker16C64G":  "KAFKA.12XUNIT",
	"Broker32C128G": "KAFKA.48XUNIT",
	"Broker64C256G": "KAFKA.100XUNIT",
}

func resourceAlibabacloudStackAlikafkaInstance() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				// instanceName
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"selected_zones": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 2,
				MinItems: 2,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cup_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"spec": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Broker4C16G", "Broker8C32G", "Broker16C64G",
					"Broker32C128G", "Broker64C256G",
				}, false),
			},
			"replicas": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.IntAtLeast(3),
				ForceNew:     true,
			},
			"disk_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntAtLeast(1),
				ForceNew:     true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"sasl": {
				Type:         schema.TypeBool,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"sasl", "plaintext"},
			},
			"plaintext": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vip_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sasl_ssl_endpoint": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"sasl_plaintext_endpoint": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"message_max_bytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"num_partitions": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"auto_create_topics_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"num_io_threads": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"queued_max_requests": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"replica_fetch_wait_max_ms": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"replica_lag_time_max_ms": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"num_network_threads": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"log_retention_bytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"replica_fetch_max_bytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"num_replica_fetchers": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"default_replication_factor": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"offsets_retention_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"background_threads": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"plaintext_endpoint": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAlikafkaInstanceCreate, resourceAlibabacloudStackAlikafkaInstanceRead,
		resourceAlibabacloudStackAlikafkaInstanceUpdate, resourceAlibabacloudStackAlikafkaInstanceDelete)
	return resource
}

func resourceAlibabacloudStackAlikafkaInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}
	var err error

	createInstanceAction := "CreateInstance"
	createInstanceResponse := make(map[string]interface{})
	createInstanceRequest := make(map[string]interface{})

	createInstanceRequest["InstanceName"] = d.Get("name").(string)
	createInstanceRequest["ZoneId"] = d.Get("zone_id").(string)

	if v, ok := d.GetOk("selected_zones"); ok {
		zones := v.([]string)
		if zones[0] == zones[1] {
			return errmsgs.WrapError(fmt.Errorf("Please select two differencent zones"))
		}
		createInstanceRequest["ZoneLength"] = 3
		createInstanceRequest["InstanceType"] = 3
		createInstanceRequest["zoneB"] = zones[0]
		createInstanceRequest["zoneC"] = zones[1]
	} else {
		createInstanceRequest["ZoneLength"] = 1
		createInstanceRequest["InstanceType"] = 1
	}

	if v, ok := d.GetOk("cup_type"); ok {
		createInstanceRequest["CpuType"] = v.(string)
	}

	if v, ok := d.GetOk("spec"); ok {
		createInstanceRequest["Spec"] = AlikafkaInstanceSpecMap[v.(string)]
	}

	if v, ok := d.GetOk("replicas"); ok {
		createInstanceRequest["Replicas"] = v.(int)
	}

	if v, ok := d.GetOk("disk_num"); ok {
		createInstanceRequest["DiskNum"] = v.(int)
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		createInstanceRequest["VipType"] = "SingleTunnel"
		if v, ok := d.GetOk("vpc_id"); ok {
			createInstanceRequest["VpcId"] = v.(string)
		} else {
			vpcService := VpcService{client}
			vswitch, err := vpcService.DescribeVSwitch(d.Get("vswitch_id").(string))
			if err != nil {
				if errmsgs.NotFoundError(err) {
					d.SetId("")
					return nil
				}
				return errmsgs.WrapError(err)
			}
			createInstanceRequest["VpcId"] = vswitch.VpcId
		}
		createInstanceRequest["VSwitchId"] = v.(string)
	} else {
		createInstanceRequest["VipType"] = "AnyTunnel"
	}

	endpointTypes := make([]string, 0)
	if v, ok := d.GetOk("sasl"); ok && v.(bool) {
		endpointTypes = append(endpointTypes, "SASL")
	}
	if v, ok := d.GetOk("plaintext"); ok && v.(bool) {
		endpointTypes = append(endpointTypes, "PLAINTEXT")
	}
	createInstanceRequest["EndpointTypes"] = strings.Join(endpointTypes, ",")

	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		createInstanceResponse, err = client.DoTeaRequest("POST", "alikafka", "2019-09-16", "CreateInstance", "", nil, createInstanceRequest, nil)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || errmsgs.NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alicloud_alikafka_instance", createInstanceAction, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if fmt.Sprint(createInstanceResponse["Success"]) == "false" {
		return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", createInstanceAction, createInstanceResponse))
	}

	d.SetId(fmt.Sprint(createInstanceResponse["instanceId"]))

	// 3. wait until running
	stateConf := BuildStateConf([]string{}, []string{"5"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	if _, ok := d.GetOk("vswitch_id"); ok {
		stateConf := BuildStateConf([]string{}, []string{"create"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, alikafkaService.AliKafkaInstanceVipStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	return nil

}

func resourceAlibabacloudStackAlikafkaInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	object, err := alikafkaService.DescribeAlikafkaInstance(d.Id())
	if err != nil {
		// Handle exceptions
		if !d.IsNewResource() && errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alikafkaService.DescribeAliKafkaInstance Failed!!! %s", err)
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("name", object.Name)
	d.Set("zone_id", object.ZoneId)
	if object.ZoneC != "" && object.ZoneB != "" {
		d.Set("seleced_zones", []string{object.ZoneC, object.ZoneB})
	}
	d.Set("cup_type", object.CpuType)
	d.Set("spec", object.SpecName)
	d.Set("replicas", object.Replicas)
	d.Set("disk_num", object.DiskNum)

	if object.VSwitchId != "" {
		d.Set("vswitch_id", object.VSwitchId)
	} else if object.VipInfo.VswId != "" {
		d.Set("vswitch_id", object.VipInfo.VswId)
	}
	if object.VpcId != "" {
		d.Set("vpc_id", object.VpcId)
	} else if object.VipInfo.VpcId != "" {
		d.Set("vpc_id", object.VipInfo.VpcId)
	}

	if object.VSwitchId != "" {
		d.Set("vip_type", "SingleTunnel")
	} else {
		d.Set("vip_type", "AnyTunnel")
	}
	enabledProtocols := object.VipInfo.EnabledProtocols
	for _, proto := range enabledProtocols {
		if proto == "SASL_SSL" {
			d.Set("sasl", true)
		} else if proto == "VPC_MODE" {
			d.Set("plaintext", true)
		}
	}

	endPointMap := object.VipInfo.EndPointMap
	if v, ok := endPointMap["SASL_SSL"]; ok {
		d.Set("sasl_ssl_endpoint", strings.Split(v, ","))
	} else {
		d.Set("sasl_ssl_endpoint", []string{})
	}
	if v, ok := endPointMap["SASL_PLAINTEXT"]; ok {
		d.Set("sasl_plaintext_endpoint", strings.Split(v, ","))
	} else {
		d.Set("sasl_plaintext_endpoint", []string{})
	}
	if v, ok := endPointMap["PLAINTEXT"]; ok {
		d.Set("plaintext_endpoint", strings.Split(v, ","))
	} else {
		d.Set("plaintext_endpoint", []string{})
	}

	d.Set("status", object.ServiceStatus)

	configMap, err := alikafkaService.DescribeAlikafkaInstanceConfigMap(d.Id())
	if err != nil {
		// Handle exceptions
		if !d.IsNewResource() && errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alikafkaService.DescribeAliKafkaInstance Failed!!! %s", err)
			return nil
		}
		return errmsgs.WrapError(err)
	}

	if v, err := strconv.Atoi(configMap.MessageMaxBytes); err == nil {
		d.Set("message_max_bytes", v)
	}
	if v, err := strconv.Atoi(configMap.NumPartitions); err == nil {
		d.Set("num_partitions", v)
	}
	d.Set("auto_create_topics_enable", string(configMap.AutoCreateTopicsEnable) == "true")
	if v, err := strconv.Atoi(configMap.NumIoThreads); err == nil {
		d.Set("num_io_threads", v)
	}
	if v, err := strconv.Atoi(configMap.QueuedMaxRequests); err == nil {
		d.Set("queued_max_requests", v)
	}
	if v, err := strconv.Atoi(configMap.ReplicaFetchWaitMaxMs); err == nil {
		d.Set("replica_fetch_wait_max_ms", v)
	}
	if v, err := strconv.Atoi(configMap.ReplicaLagTimeMaxMs); err == nil {
		d.Set("replica_lag_time_max_ms", v)
	}
	if v, err := strconv.Atoi(configMap.NumNetworkThreads); err == nil {
		d.Set("num_network_threads", v)
	}
	if v, err := strconv.Atoi(configMap.LogRetentionBytes); err == nil {
		d.Set("log_retention_bytes", v)
	}
	if v, err := strconv.Atoi(configMap.ReplicaFetchMaxBytes); err == nil {
		d.Set("replica_fetch_max_bytes", v)
	}
	if v, err := strconv.Atoi(configMap.NumReplicaFetchers); err == nil {
		d.Set("num_replica_fetchers", v)
	}
	if v, err := strconv.Atoi(configMap.DefaultReplicationFactor); err == nil {
		d.Set("default_replication_factor", v)
	}
	if v, err := strconv.Atoi(configMap.OffsetsRetentionMinutes); err == nil {
		d.Set("offsets_retention_minutes", v)
	}
	if v, err := strconv.Atoi(configMap.BackgroundThreads); err == nil {
		d.Set("background_threads", v)
	}

	return nil
}

func resourceAlibabacloudStackAlikafkaInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}
	// 	var err error

	if d.IsNewResource() {
		return nil
	}

	var configKeys = []string{"message.max.bytes", "num.partitions",
		"auto.create.topics.enable", "num.io.threads", "queued.max.requests",
		"replica.fetch.wait.max.ms", "replica.lag.time.max.ms", "num.network.threads",
		"log.retention.bytes", "replica.fetch.max.bytes", "num.replica.fetchers",
		"default.replication.factor", "offsets.retention.minutes", "background.threads",
	}
	for _, configKey := range configKeys {
		schemaName := strings.Replace(configKey, ".", "_", -1)
		if v, ok := d.GetOk(schemaName); ok && d.HasChange(schemaName) {
			var value string
			if schemaName == "auto.create.topics.enable" {
				if v.(bool) {
					value = "true"
				} else {
					value = "false"
				}
			} else {
				value = strconv.Itoa(v.(int))
			}
			action := "UpdateInstanceConfig"
			request := alikafkaService.client.NewCommonRequest("POST", "alikafka", "2019-09-16", action, "")
			request.QueryParams["InstanceId"] = d.Id()
			request.QueryParams["Config"] = configKey
			request.QueryParams["Value"] = value
			//不判断是否失败，如果失败了等待二次apply
			client.ProcessCommonRequest(request)
		}
	}
	return nil
}

func resourceAlibabacloudStackAlikafkaInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	action := "DeleteInstance"
	var err error
	var response map[string]interface{}
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}

	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.DoTeaRequest("POST", "alikafka", "2019-09-16", action, "", nil, nil, request)
		if err != nil {
			if errmsgs.NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}
