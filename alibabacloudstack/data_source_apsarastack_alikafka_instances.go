package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudAlikafkaInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlikafkaInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							// instanceName
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							// instanceName
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"selected_zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"cup_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replicas": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sasl": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"plaintext": {
							Type:     schema.TypeBool,
							Computed: true,
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
							Computed: true,
						},
						"num_partitions": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"auto_create_topics_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"num_io_threads": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"queued_max_requests": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"replica_fetch_wait_max_ms": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"replica_lag_time_max_ms": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"num_network_threads": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"log_retention_bytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"replica_fetch_max_bytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"num_replica_fetchers": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"default_replication_factor": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"offsets_retention_minutes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"background_threads": {
							Type:     schema.TypeInt,
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
				},
			},
		},
	}
}

func dataSourceAlicloudAlikafkaInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	action := "GetInstanceList"

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	request := client.NewCommonRequest("POST", "alikafka", "2019-09-16", action, "")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var bresponse *responses.CommonResponse
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		bresponse, err = client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	var instanceListResp GetInstanceListResponse
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &instanceListResp)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_alikafka_instances", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	objects := make([]InstanceVO, 0)
	for _, v := range instanceListResp.InstanceList {
		if nameRegex != nil && !nameRegex.MatchString(v.Name) {
			continue
		}

		if len(idsMap) > 0 {
			if _, ok := idsMap[v.InstanceId]; !ok {
				continue
			}
		}
		objects = append(objects, v)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)

	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		mapping := map[string]interface{}{}

		mapping["name"] = object.Name
		mapping["zone_id"] = object.ZoneId
		if object.ZoneC != "" && object.ZoneB != "" {
			mapping["selected_zones"] = []string{object.ZoneC, object.ZoneB}
		}

		mapping["id"] = object.InstanceId
		mapping["cup_type"] = object.CpuType
		mapping["spec"] = object.SpecName
		mapping["replicas"] = object.Replicas
		mapping["disk_num"] = object.DiskNum
		if object.VSwitchId != "" {
			mapping["vswitch_id"] = object.VSwitchId
		} else if object.VipInfo.VswId != "" {
			mapping["vswitch_id"] = object.VipInfo.VswId
		}
		if object.VpcId != "" {
			mapping["vpc_id"] = object.VpcId
		} else if object.VipInfo.VpcId != "" {
			mapping["vpc_id"] = object.VipInfo.VpcId
		}
		if object.VSwitchId != "" {
			mapping["vip_type"] = "SingleTunnel"
		} else {
			mapping["vip_type"] = "AnyTunnel"
		}
		enabledProtocols := object.VipInfo.EnabledProtocols
		for _, proto := range enabledProtocols {
			if proto == "SASL_SSL" {
				mapping["sasl"] = true
			} else if proto == "VPC_MODE" {
				mapping["plaintext"] = true
			}
		}

		endPointMap := object.VipInfo.EndPointMap
		if v, ok := endPointMap["SASL_SSL"]; ok {
			mapping["sasl_ssl_endpoint"] = strings.Split(v, ",")
		} else {
			mapping["sasl_ssl_endpoint"] = []string{}
		}
		if v, ok := endPointMap["SASL_PLAINTEXT"]; ok {
			mapping["sasl_plaintext_endpoint"] = strings.Split(v, ",")
		} else {
			mapping["sasl_plaintext_endpoint"] = []string{}
		}
		if v, ok := endPointMap["PLAINTEXT"]; ok {
			mapping["plaintext_endpoint"] = strings.Split(v, ",")
		} else {
			mapping["plaintext_endpoint"] = []string{}
		}

		mapping["status"] = object.ServiceStatus

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, mapping["name"])
		//id := fmt.Sprint(object["InstanceId"])

		//AlikaService := AlikafkaService{client}
		if d.Get("enable_details").(bool) {
			configMap, err := alikafkaService.DescribeAlikafkaInstanceConfigMap(mapping["id"].(string))
			if err != nil {
				// Handle exceptions
				if !d.IsNewResource() && errmsgs.NotFoundError(err) {
					log.Printf("[DEBUG] Resource alikafkaService.DescribeAliKafkaInstance Failed!!! %s", err)
					return nil
				}
				return errmsgs.WrapError(err)
			}

			if v, err := strconv.Atoi(configMap.MessageMaxBytes); err == nil {
				mapping["message_max_bytes"] = v
			}
			if v, err := strconv.Atoi(configMap.NumPartitions); err == nil {
				mapping["num_partitions"] = v
			}
			mapping["auto_create_topics_enable"] = string(configMap.AutoCreateTopicsEnable) == "true"
			if v, err := strconv.Atoi(configMap.NumIoThreads); err == nil {
				mapping["num_io_threads"] = v
			}
			if v, err := strconv.Atoi(configMap.QueuedMaxRequests); err == nil {
				mapping["queued_max_requests"] = v
			}
			if v, err := strconv.Atoi(configMap.ReplicaFetchWaitMaxMs); err == nil {
				mapping["replica_fetch_wait_max_ms"] = v
			}
			if v, err := strconv.Atoi(configMap.ReplicaLagTimeMaxMs); err == nil {
				mapping["replica_lag_time_max_ms"] = v
			}
			if v, err := strconv.Atoi(configMap.NumNetworkThreads); err == nil {
				mapping["num_network_threads"] = v
			}
			if v, err := strconv.Atoi(configMap.LogRetentionBytes); err == nil {
				mapping["log_retention_bytes"] = v
			}
			if v, err := strconv.Atoi(configMap.ReplicaFetchMaxBytes); err == nil {
				mapping["replica_fetch_max_bytes"] = v
			}
			if v, err := strconv.Atoi(configMap.NumReplicaFetchers); err == nil {
				mapping["num_replica_fetchers"] = v
			}
			if v, err := strconv.Atoi(configMap.DefaultReplicationFactor); err == nil {
				mapping["default_replication_factor"] = v
			}
			if v, err := strconv.Atoi(configMap.OffsetsRetentionMinutes); err == nil {
				mapping["offsets_retention_minutes"] = v
			}
			if v, err := strconv.Atoi(configMap.BackgroundThreads); err == nil {
				mapping["background_threads"] = v
			}
		}

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("instances", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
