package alibabacloudstack

import (
	//	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"

	//	"github.com/aliyun/alibaba-cloud-sdk-go/services/cs"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	KubernetesClusterNetworkTypeFlannel = "flannel"
	KubernetesClusterNetworkTypeTerway  = "terway"

	KubernetesClusterLoggingTypeSLS = "SLS"
	ClusterType                     = "Kubernetes"
	OsType                          = "Linux"
	Platform                        = "CentOS"
	RuntimeName                     = "docker"
	RuntimeVersion                  = "19.03.5"
	PortRange                       = "30000-32767"
)

var (
	KubernetesClusterNodeCIDRMasksByDefault = 24
)

func resourceAlibabacloudStackCSKubernetes() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 63),
			},
			"master_disk_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          40,
				ValidateFunc:     validation.IntBetween(40, 500),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"master_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  DiskCloudSSD,
				ValidateFunc: validation.StringInSlice([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD), string(DiskCloudPPERF), string(DiskCloudSPERF)}, false),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"delete_protection": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"num_of_nodes": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"worker_disk_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          40,
				ValidateFunc:     validation.IntBetween(20, 32768),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"worker_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  DiskCloudSSD,
				ValidateFunc: validation.StringInSlice([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD), string(DiskCloudPPERF), string(DiskCloudSPERF)}, false),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			// 			"worker_data_disk_size": {
			// 				Type:             schema.TypeInt,
			// 				Optional:         true,
			// 				Default:          40,
			// 				ValidateFunc:     validation.IntBetween(20, 32768),
			// 				DiffSuppressFunc: workerDataDiskSizeSuppressFunc,
			// 			},
			// 			"worker_data_disk_category": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 				ValidateFunc: validation.StringInSlice([]string{
			// 					string(DiskCloudEfficiency), string(DiskCloudSSD), string(DiskCloudPPERF), string(DiskCloudSPERF)}, false),
			// 				DiffSuppressFunc: csForceUpdateSuppressFunc,
			// 			},
			"worker_data_disks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Default:  "flannel",
							Optional: true,
						},
						"encrypted": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"auto_snapshot_policy_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"performance_level": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"category": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"master_storage_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"master_storage_set_partition_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 2000),
			},
			"worker_storage_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"worker_storage_set_partition_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 2000),
			},
			// 			"exclude_autoscaler_nodes": {
			// 				Type:     schema.TypeBool,
			// 				Default:  false,
			// 				Optional: true,
			// 			},
			//"worker_data_disk": {
			//	Type:     schema.TypeBool,
			//	Default:  false,
			//	Optional: true,
			//},
			// global configurations
			// Terway network
			"pod_vswitch_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MaxItems:         10,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
				ConflictsWith:    []string{"pod_cidr"},
			},
			// Flannel network
			"pod_cidr": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"service_cidr": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"node_cidr_mask": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"new_nat_gateway": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"password": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				ConflictsWith:    []string{"kms_encrypted_password", "key_name"},
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"key_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"password", "kms_encrypted_password"},
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"kms_encrypted_password": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password", "key_name"},
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
			// 			"user_ca": {
			// 				Type:             schema.TypeString,
			// 				Optional:         true,
			// 				DiffSuppressFunc: csForceUpdateSuppressFunc,
			// 			},
			"enable_ssh": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"node_port_range": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  PortRange,
			},
			"image_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: imageIdSuppressFunc,
			},
			// 			"install_cloud_monitor": {
			// 				Type:             schema.TypeBool,
			// 				Optional:         true,
			// 				Default:          true,
			// 				DiffSuppressFunc: csForceUpdateSuppressFunc,
			// 			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  ClusterType,
			},
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  OsType,
			},
			"platform": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  Platform,
			},
			// cpu policy options of kubelet
			"cpu_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "none",
				ValidateFunc: validation.StringInSlice([]string{"none", "static"}, false),
			},
			"proxy_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"iptables", "ipvs"}, false),
			},
			"addons": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Default:  "flannel",
							Optional: true,
						},
						"config": {
							Type:     schema.TypeString,
							Optional: true,
						},
						//"disabled": {
						//	Type:     schema.TypeBool,
						//	Optional: true,
						//	Default:  false,
						//},
					},
				},
			},
			"slb_internet_enabled": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			// computed parameters
			"kube_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_cert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_ca_cert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// 			"connections": {
			// 				Type:     schema.TypeList,
			// 				Computed: true,
			// 				Elem: &schema.Resource{
			// 					Schema: map[string]*schema.Schema{
			// 						"api_server_internet": {
			// 							Type:     schema.TypeString,
			// 							Computed: true,
			// 						},
			// 						"api_server_intranet": {
			// 							Type:     schema.TypeString,
			// 							Computed: true,
			// 						},
			// 						"master_public_ip": {
			// 							Type:     schema.TypeString,
			// 							Computed: true,
			// 						},
			// 						"service_domain": {
			// 							Type:     schema.TypeString,
			// 							Computed: true,
			// 						},
			// 					},
			// 				},
			// 			},
			// 			"slb_id": {
			// 				Type:       schema.TypeString,
			// 				Computed:   true,
			// 				Deprecated: "Field 'slb_id' has been deprecated from provider version 1.9.2. New field 'slb_internet' replaces it.",
			// 			},
			// 			"slb_internet": {
			// 				Type:     schema.TypeString,
			// 				Computed: true,
			// 			},
			"slb_intranet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"master_system_disk_performance_level": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"worker_system_disk_performance_level": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_enterprise_security_group": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"cloud_monitor_flags": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"nat_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"runtime": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  RuntimeName,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  RuntimeVersion,
						},
					},
				},
			},
			"master_nodes": {
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
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"worker_nodes": {
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
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			// remove parameters below
			// mix vswitch_ids between master and worker is not a good guidance to create cluster
			// 			"worker_instance_type": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 			},
			"master_instance_types": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"master_vswitch_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"worker_instance_types": {
				Type:          schema.TypeList,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"instances"},
				Optional:      true,
			},
			"worker_vswitch_ids": {
				Type:          schema.TypeList,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"instances"},
				Optional:      true,
			},
			"instances": {
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"worker_instance_types", "worker_vswitch_ids", "worker_disk_category"},
				Optional:      true,
			},
			"format_disk": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"keep_instance_name": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// 			"vswitch_ids": {
			// 				Type:     schema.TypeList,
			// 				Optional: true,
			// 				Elem: &schema.Schema{
			// 					Type:         schema.TypeString,
			// 					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
			// 				},
			// 				MinItems:         3,
			// 				MaxItems:         5,
			// 				DiffSuppressFunc: csForceUpdateSuppressFunc,
			// 				//Removed:          "Field 'vswitch_ids' has been removed from provider version 1.75.0. New field 'master_vswitch_ids' and 'worker_vswitch_ids' replace it.",
			// 			},
			"master_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.IntAtLeast(3),
			},
			// single instance type would cause extra troubles
			// 			"master_instance_type": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 			},
			// force update is a high risk operation
			// 			"force_update": {
			// 				Type:     schema.TypeBool,
			// 				Optional: true,
			// 				Default:  false,
			// 				//Removed:  "Field 'force_update' has been removed from provider version 1.75.0.",
			// 			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// single az would be never supported.
			//"vswitch_id": {
			//	Type:     schema.TypeString,
			//	Required: true,
			//	//Removed:  "Field 'vswitch_id' has been removed from provider version 1.75.0. New field 'master_vswitch_ids' and 'worker_vswitch_ids' replaces it.",
			//},
			"timeout_mins": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},
			"nodes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// too hard to use this config
			// 			"log_config": {
			// 				Type:     schema.TypeList,
			// 				Optional: true,
			// 				MaxItems: 1,
			// 				Elem: &schema.Resource{
			// 					Schema: map[string]*schema.Schema{
			// 						"type": {
			// 							Type:         schema.TypeString,
			// 							ValidateFunc: validation.StringInSlice([]string{KubernetesClusterLoggingTypeSLS}, false),
			// 							Required:     true,
			// 						},
			// 						"project": {
			// 							Type:     schema.TypeString,
			// 							Optional: true,
			// 						},
			// 					},
			// 				},
			// 				DiffSuppressFunc: csForceUpdateSuppressFunc,
			// 				//Removed:          "Field 'log_config' has been removed from provider version 1.75.0. New field 'addons' replaces it.",
			// 			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// 			"node_name_mode": {
			// 				Type:         schema.TypeString,
			// 				Optional:     true,
			// 				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^customized,[a-z0-9]([-a-z0-9\.])*,([5-9]|[1][0-2]),([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`), "Each node name consists of a prefix, an IP substring, and a suffix. For example, if the node IP address is 192.168.0.55, the prefix is aliyun.com, IP substring length is 5, and the suffix is test, the node name will be aliyun.com00055test."),
			// 			},
			"worker_ram_role_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// 			"service_account_issuer": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 				ForceNew: true,
			// 			},
			// 			"api_audiences": {
			// 				Type:     schema.TypeList,
			// 				Optional: true,
			// 				Elem: &schema.Schema{
			// 					Type: schema.TypeString,
			// 				},
			// 				ForceNew: true,
			// 			},
			"nodepool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackCSKubernetesCreate,
		resourceAlibabacloudStackCSKubernetesRead, resourceAlibabacloudStackCSKubernetesUpdate,
		resourceAlibabacloudStackCSKubernetesDelete)
	return resource
}

type Response struct {
	RequestId string `json:"request_id"`
}
type ClusterCommonResponse struct {
	Response
	ClusterID  string `json:"cluster_id"`
	Token      string `json:"token,omitempty"`
	TaskId     string `json:"task_id,omitempty"`
	InstanceId string `json:"instanceId"`
}

func resourceAlibabacloudStackCSKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("check meta %v", meta)
	client := meta.(*connectivity.AlibabacloudStackClient)
	csService := CsService{client}
	invoker := NewInvoker()

	request := client.NewCommonRequest("POST", "CS", "2015-12-15", "CreateCluster", "/clusters")
	request.SetContentType("application/json")
	request.SetContent([]byte("{}")) // 必须指定，否则SDK会将类型修改为www-form，最终导致cr有一定的随机概率失败

	request.QueryParams = map[string]string{
		"Action":    "CreateCluster",
		"SNatEntry": "false",
	}

	body := map[string]interface{}{
		"Product":                              "Cs",
		"os_type":                              d.Get("os_type").(string),
		"platform":                             d.Get("platform").(string),
		"cluster_type":                         d.Get("cluster_type").(string),
		"region_id":                            client.RegionId,
		"timeout_mins":                         d.Get("timeout_mins").(int),
		"disable_rollback":                     true,
		"kubernetes_version":                   d.Get("version").(string),
		"container_cidr":                       d.Get("pod_cidr").(string),
		"service_cidr":                         d.Get("service_cidr").(string),
		"name":                                 d.Get("name").(string),
		"master_instance_types":                d.Get("master_instance_types").([]interface{}),
		"master_vswitch_ids":                   d.Get("master_vswitch_ids").([]interface{}),
		"num_of_nodes":                         d.Get("num_of_nodes").(int),
		"master_count":                         d.Get("master_count").(int),
		"snat_entry":                           d.Get("new_nat_gateway").(bool),
		"endpoint_public_access":               d.Get("slb_internet_enabled").(bool),
		"ssh_flags":                            d.Get("enable_ssh").(bool),
		"master_system_disk_category":          d.Get("master_disk_category").(string),
		"master_system_disk_size":              d.Get("master_disk_size").(int),
		"deletion_protection":                  d.Get("delete_protection").(bool),
		"node_cidr_mask":                       d.Get("node_cidr_mask").(string),
		"vpcid":                                d.Get("vpc_id").(string),
		"proxy_mode":                           d.Get("proxy_mode").(string),
		"user_data":                            d.Get("user_data").(string),
		"node_port_range":                      d.Get("node_port_range").(string),
		"cpu_policy":                           d.Get("cpu_policy").(string),
		"cloud_monitor_flags":                  d.Get("cloud_monitor_flags").(bool),
		"master_system_disk_performance_level": d.Get("master_system_disk_performance_level").(string),
		"worker_system_disk_performance_level": d.Get("worker_system_disk_performance_level").(string),
		"image_id":                             d.Get("image_id").(string),
	}

	pod := 0
	if v, ok := d.GetOk("addons"); ok {
		body["addons"] = v
		if all, ok := v.([]interface{}); ok {
			for _, a := range all {
				log.Printf("check addon %v", a)
				if addon, ok := a.(map[string]interface{}); ok && addon["name"] == "terway-eniip" {
					pod = 1
					log.Printf("pod request true")
				}
			}
		}
	}
	if data, ok := d.GetOk("worker_data_disks"); ok {
		body["worker_data_disks"] = data
	}
	if v, ok := d.GetOk("tags"); ok {
		var tags = []map[string]interface{}{}
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"key":   key,
				"value": value.(string),
			})
		}
		body["tags"] = tags
	}
	if v, ok := d.GetOk("runtime"); ok {
		all, _ := v.([]interface{})
		for _, runtime := range all {
			body["runtime"] = runtime
		}
	}
	if v, ok := d.GetOk("is_enterprise_security_group"); ok && v.(bool) {
		if v, ok := d.GetOk("security_group_id"); ok && v.(string) != "" {
			return fmt.Errorf("security_group_id must be `` or nil when is_enterprise_security_group is `true`")
		}
		body["is_enterprise_security_group"] = v.(bool)
	} else {
		if v, ok := d.GetOk("security_group_id"); ok && v.(string) != "" {
			body["security_group_id"] = d.Get("security_group_id").(string)
		} else {
			return fmt.Errorf("security_group_id must be set when is_enterprise_security_group is `false` or not set")
		}
	}
	if v, ok := d.GetOk("pod_vswitch_ids"); ok && pod == 1 {
		body["pod_vswitch_ids"] = expandStringList(v.(*schema.Set).List())
	}

	if password := d.Get("password").(string); password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return errmsgs.WrapError(err)
			}
			password = decryptResp.Plaintext
		}
		body["login_Password"] = password
	} else {
		body["login_Password"] = password
	}
	if key_name, ok := d.GetOk("key_name"); ok && key_name != "" {
		body["key_pair"] = key_name.(string)
	}

	if v, ok := d.GetOk("instances"); ok {
		body["format_disk"] = d.Get("format_disk").(bool)
		body["keep_instance_name"] = d.Get("keep_instance_name").(bool)
		body["instances"] = expandStringList(v.(*schema.Set).List())
	} else {
		body["worker_instance_types"] = d.Get("worker_instance_types").([]interface{})
		body["worker_vswitch_ids"] = d.Get("worker_vswitch_ids").([]interface{})
		body["worker_system_disk_category"] = d.Get("worker_disk_category").(string)
		body["worker_system_disk_size"] = d.Get("worker_disk_size").(int)
		body["master_storage_set_id"] = d.Get("master_storage_set_id").(string)
		body["master_storage_set_partition_number"] = d.Get("master_storage_set_partition_number").(int)
		body["worker_storage_set_id"] = d.Get("worker_storage_set_id").(string)
		body["worker_storage_set_partition_number"] = d.Get("worker_storage_set_partition_number").(int)
	}
	log.Printf("[DEBUG] Request body: %s", body)
	if data, err := json.Marshal(body); err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	} else {
		request.SetContent(data)
	}
	var err error
	err = nil
	var cluster *responses.CommonResponse
	if err = invoker.Run(func() error {
		cluster, err = client.ProcessCommonRequest(request)
		addDebug("CreateKubernetesCluster", cluster, request, request.QueryParams)
		return err
	}); err != nil {
		//return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_cs_kubernetes", "CreateKubernetesCluster", raw)
		return err
	}

	clusterresponse := ClusterCommonResponse{}
	if cluster.IsSuccess() == false {
		//return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm", "API Action", cluster.GetHttpContentString())
		return err
	}
	ok := json.Unmarshal(cluster.GetHttpContentBytes(), &clusterresponse)
	if ok != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_cs_kubernetes", "ParseKubernetesClusterResponse", cluster)
	}
	d.SetId(clusterresponse.ClusterID)

	stateConf := BuildStateConf([]string{"initial", " "}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 15*time.Minute, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	stateConf.NotFoundChecks = 1000
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}

func resourceAlibabacloudStackCSKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	csService := CsService{client}
	d.Partial(true)
	invoker := NewInvoker()

	var nodepoolid string
	if d.IsNewResource() {
		var err error
		nodepoolid, err = getDefaultNodePoolId(csService, d.Id())
		if err != nil {
			return err
		}
	} else {
		nodepoolid = d.Get("nodepool_id").(string)
	}
	if d.HasChange("num_of_nodes") && !d.IsNewResource() {
		password := d.Get("password").(string)
		if password == "" {
			if v := d.Get("kms_encrypted_password").(string); v != "" {
				kmsService := KmsService{client}
				decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
				if err != nil {
					return errmsgs.WrapError(err)
				}
				password = decryptResp.Plaintext
			}
		}

		oldV, newV := d.GetChange("num_of_nodes")
		oldValue, ok := oldV.(int)
		if ok != true {
			return errmsgs.WrapErrorf(fmt.Errorf("num_of_nodes old value can not be parsed"), "parseError %d", oldValue)
		}
		newValue, ok := newV.(int)
		if ok != true {
			return errmsgs.WrapErrorf(fmt.Errorf("num_of_nodes new value can not be parsed"), "parseError %d", newValue)
		}

		if newValue < oldValue {
			//return errmsgs.WrapErrorf(fmt.Errorf("num_of_nodes can not be less than before"), "scaleOutFailed %d:%d", newValue, oldValue)
			object, err := csService.DescribeClusterNodes(d.Id(), nodepoolid)
			if err != nil {
				if errmsgs.NotFoundError(err) {
					d.SetId("")
					return nil
				}
				return errmsgs.WrapError(err)
			}
			var allNodeName []string
			for _, value := range object.Nodes {
				allNodeName = append(allNodeName, value.NodeName)
			}
			count := oldValue - newValue
			removeNodesName := allNodeName[:count]
			if len(removeNodesName) > 0 {
			}
			req := client.NewCommonRequest("POST", "CS", "2015-12-15", "RemoveClusterNodes", fmt.Sprintf("/api/v2/clusters/%s/nodes/remove", d.Id()))
			body := fmt.Sprintf("{\"%s\":%t,\"%s\":%t,\"%s\":%q}",
				"release_node", true,
				"drain_node", true,
				"nodes", removeNodesName,
			)
			log.Printf("[DEBUG]RemoveClusterNodes Request body: %s", body)
			req.SetContent([]byte(body))
			req.Headers["x-acs-content-type"] = "application/json"
			var resp *responses.CommonResponse
			if err := invoker.Run(func() error {
				resp, err = client.ProcessCommonRequest(req)
				return err
			}); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, nodepoolid, "DeleteKubernetesClusterNodes", errmsgs.AlibabacloudStackSdkGoERROR)
			}
			if resp.IsSuccess() == false {
				//return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm", "API Action", cluster.GetHttpContentString())
				return err
			}
			stateConf := BuildStateConf([]string{"removing"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(nodepoolid, d.Id(), []string{"deleting", "failed"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
			}
		}

		if newValue > oldValue {
			request := client.NewCommonRequest("POST", "CS", "2015-12-15", "ScaleClusterNodePool", fmt.Sprintf("/clusters/%s/nodepools/%s", d.Id(), nodepoolid))
			body := fmt.Sprintf("{\"%s\":%d}",
				"count", int64(newValue)-int64(oldValue),
			)
			request.QueryParams["NodepoolId"] = nodepoolid
			request.QueryParams["ClusterId"] = d.Id()
			request.SetContent([]byte(body))
			request.Headers["x-acs-content-type"] = "application/json"
			//var err error
			var resp *responses.CommonResponse
			if err := invoker.Run(func() error {
				var err error
				resp, err = client.ProcessCommonRequest(request)
				return err
			}); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_cs_kubernetes", "CreateKubernetesCluster", resp)
			}

			if debugOn() {
				resizeRequestMap := make(map[string]interface{})
				resizeRequestMap["ClusterId"] = d.Id()
				resizeRequestMap["Args"] = request.GetQueryParams()
				addDebug("ScaleClusterNodePool", resp, resizeRequestMap)
			}

			stateConf := BuildStateConf([]string{"scaling"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(nodepoolid, d.Id(), []string{"deleting", "failed"}))

			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
			}
			//d.SetPartial("num_of_nodes")
		}
	}

	d.Partial(false)
	return nil
}

func resourceAlibabacloudStackCSKubernetesRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	csService := CsService{client}
	object, err := csService.DescribeCsKubernetes(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	nodepoolid, err := getDefaultNodePoolId(csService, d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("nodepool_id", nodepoolid)
	clusternode, err := csService.DescribeClusterNodes(d.Id(), nodepoolid)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("name", object.Name)
	// node_count, err := csService.GetCsK8sNodesCount(d.Id())
	// if err != nil {
	// 	return errmsgs.WrapError(err)
	// }
	nodepool, err := csService.DescribeCsKubernetesNodePool(nodepoolid, d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("num_of_nodes", nodepool.Status.HealthyNodes)
	//d.Set("id", object.ClusterId)
	//d.Set("state", object.State)
	d.Set("vpc_id", object.VpcId)
	//d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("pod_cidr", object.ContainerCIDR)
	d.Set("version", object.CurrentVersion)
	d.Set("delete_protection", object.DeletionProtection)
	d.Set("version", object.InitVersion)
	var smaster, sworker []map[string]interface{}
	//var MasterNodes, WorkerNodes map[string]interface{}
	for _, k := range clusternode.Nodes {
		if k.InstanceRole == "Master" {
			MasterNodes := map[string]interface{}{
				"id":         k.InstanceID,
				"name":       k.InstanceName,
				"private_ip": fmt.Sprintf("%s", k.IPAddress),
			}
			smaster = append(smaster, MasterNodes)
		} else {
			WorkerNodes := map[string]interface{}{
				"id":         k.InstanceID,
				"name":       k.InstanceName,
				"private_ip": fmt.Sprintf("%s", k.IPAddress),
			}
			sworker = append(sworker, WorkerNodes)
		}
	}

	d.Set("master_nodes", smaster)
	d.Set("worker_nodes", sworker)
	if err := d.Set("tags", flattenTagsConfig(object.Tags)); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func resourceAlibabacloudStackCSKubernetesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	csService := CsService{client}
	invoker := NewInvoker()
	body := fmt.Sprintf("{\"%s\":\"%t\",\"%s\":\"%s\"}", "keep_slb", false, "ClusterId", d.Id())
	request := client.NewCommonRequest("DELETE", "CS", "2015-12-15", "DeleteCluster", fmt.Sprintf("/clusters/%s", d.Id()))
	request.QueryParams["ClusterId"] = d.Id()
	request.SetContent([]byte(body))
	request.Headers["x-acs-content-type"] = "application/json"
	var response *responses.CommonResponse
	err := resource.Retry(30*time.Minute, func() *resource.RetryError {
		if err := invoker.Run(func() error {
			var err error
			response, err = client.ProcessCommonRequest(request)
			addDebug("DeleteCluster", response, request, request.QueryParams)
			return err
		}); err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteCluster", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}

	stateConf := BuildStateConf([]string{"running", "deleting", "initial"}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Minute, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"delete_failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}

func getDefaultNodePoolId(csService CsService, clusterId string) (string, error) {
	var nodepoolid string
	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		nodepool, err := csService.DescribeClusterNodePools(clusterId)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		for _, k := range nodepool.Nodepools {
			//Considering multiple nodepools
			if k.NodepoolInfo.IsDefault {
				nodepoolid = k.NodepoolInfo.NodepoolID
			}
		}
		if nodepoolid == "" {
			if len(nodepool.Nodepools) == 1 && nodepool.Nodepools[0].NodepoolInfo.Name == "default-nodepool" {
				nodepoolid = nodepool.Nodepools[0].NodepoolInfo.NodepoolID
			} else {
				return resource.RetryableError(errmsgs.WrapErrorf(fmt.Errorf("can not found default node_pool"), "DescribeClusterNodePools", nodepool.Nodepools))
			}
		}
		return nil
	}); err != nil {
		return "", err
	}
	return nodepoolid, nil
}

func updateKubernetesClusterTag(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	csService := CsService{client}
	d.Partial(true)
	invoker := NewInvoker()
	request := client.NewCommonRequest("POST", "CS", "2015-12-15", "ModifyClusterTags", fmt.Sprintf("/clusters/%s/tags", d.Id()))
	tagss := make([]interface{}, 0)
	if v, ok := d.GetOk("tags"); ok && len(v.(map[string]interface{})) > 0 {
		for key, value := range v.(map[string]interface{}) {
			tagss = append(tagss, Tag{
				Key:   key,
				Value: value.(string),
			})
		}
	}
	tagsBytes, _ := json.Marshal(map[string]interface{}{"tags": tagss})
	request.SetContent([]byte(tagsBytes))
	request.QueryParams["ClusterId"] = d.Id()
	request.QueryParams["ProductName"] = "cs"
	request.QueryParams["SignatureVersion"] = "1.0"
	var err error
	var raw *responses.CommonResponse
	if err = invoker.Run(func() error {
		raw, err = client.ProcessCommonRequest(request)
		addDebug("ModifyClusterTags", raw, request, request.QueryParams)
		return err
	}); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_cs_kubernetes", "ModifyClusterTags", raw)
	}
	stateConf := BuildStateConf([]string{"scaling"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}
