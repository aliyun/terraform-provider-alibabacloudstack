package alibabacloudstack

import (
	//	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
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
			"master_disk_encrypt_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"sm4-128", "aes-256"}, false),
			},
			"master_disk_kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"master_disk_encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"delete_protection": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"num_of_nodes": {
				Type:         schema.TypeInt,
				Default:      0,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
				Description:  "The number of worker nodes. Set to 0 to create a cluster without a default nodepool.",
			},
			"worker_disk_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          40,
				ValidateFunc:     validation.IntBetween(20, 32768),
				DiffSuppressFunc: numOfNodesZeroSuppressFunc,
				Description:      "The `worker_disk_size` field will become Computed in version 3.21.0 and will no longer support input.",
			},
			"worker_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  DiskCloudSSD,
				ValidateFunc: validation.StringInSlice([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD), string(DiskCloudPPERF), string(DiskCloudSPERF)}, false),
				DiffSuppressFunc: numOfNodesZeroSuppressFunc,
				Description:      "The `worker_disk_category` field will become Computed in version 3.21.0 and will no longer support input.",
			},
			"worker_disk_encrypt_algorithm": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"sm4-128", "aes-256"}, false),
				DiffSuppressFunc: numOfNodesZeroSuppressFunc,
				Description:      "The `worker_disk_encrypt_algorithm` field will become Computed in version 3.21.0 and will no longer support input.",
			},
			"worker_disk_kms_key_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: numOfNodesZeroSuppressFunc,
				Description:      "The `worker_disk_kms_key_id` field will become Computed in version 3.21.0 and will no longer support input.",
			},
			"worker_disk_encrypted": {
				Type:             schema.TypeBool,
				Optional:         true,
				DiffSuppressFunc: numOfNodesZeroSuppressFunc,
				Description:      "The `worker_disk_encrypted` field will become Computed in version 3.21.0 and will no longer support input.",
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
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						// Not supported
						// "encrypt_algorithm": {
						// 	Type:         schema.TypeString,
						// 	Optional:     true,
						// 	ValidateFunc: validation.StringInSlice([]string{"sm4-128", "aes-256"}, false),
						// },
					},
				},
				Description: "The `worker_data_disks` field will become Computed in version 3.21.0 and will no longer support input.",
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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The `worker_storage_set_id` field will become Computed in version 3.21.0 and will no longer support input.",
			},
			"worker_storage_set_partition_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 2000),
				Description:  "The `worker_storage_set_partition_number` field will become Computed in version 3.21.0 and will no longer support input.",
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
				DiffSuppressFunc: numOfNodesZeroSuppressFunc,
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
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "none",
				ValidateFunc:     validation.StringInSlice([]string{"none", "static"}, false),
				DiffSuppressFunc: numOfNodesZeroSuppressFunc,
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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The `worker_system_disk_performance_level` field will become Computed in version 3.21.0 and will no longer support input.",
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
				Type:             schema.TypeList,
				Elem:             &schema.Schema{Type: schema.TypeString},
				ConflictsWith:    []string{"instances"},
				Optional:         true,
				DiffSuppressFunc: numOfNodesZeroSuppressFunc,
				Description:      "The `worker_instance_types` field will become Computed in version 3.21.0 and will no longer support input.",
			},
			"worker_vswitch_ids": {
				Type:             schema.TypeList,
				Elem:             &schema.Schema{Type: schema.TypeString},
				ConflictsWith:    []string{"instances"},
				Optional:         true,
				DiffSuppressFunc: numOfNodesZeroSuppressFunc,
				Description:      "The `worker_vswitch_ids` field will become Computed in version 3.21.0 and will no longer support input.",
			},
			"instances": {
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"worker_instance_types", "worker_vswitch_ids", "worker_disk_category"},
				Optional:      true,
				Description:   "The `instances` field will become Computed in version 3.21.0 and will no longer support input.",
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
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: numOfNodesZeroSuppressFunc,
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
	request.SetContent([]byte("{}")) // Must be specified, otherwise the SDK will change the type to www-form, which will eventually cause the cr to fail with a certain random probability

	request.QueryParams = map[string]string{
		"Action":    "CreateCluster",
		"SNatEntry": "false",
	}

	body := map[string]interface{}{
		"Product":               "Cs",
		"os_type":               d.Get("os_type").(string),
		"platform":              d.Get("platform").(string),
		"cluster_type":          d.Get("cluster_type").(string),
		"region_id":             client.RegionId,
		"timeout_mins":          d.Get("timeout_mins").(int),
		"disable_rollback":      true,
		"ip_stack":              "ipv4",
		"kubernetes_version":    d.Get("version").(string),
		"container_cidr":        d.Get("pod_cidr").(string),
		"service_cidr":          d.Get("service_cidr").(string),
		"name":                  d.Get("name").(string),
		"master_instance_types": d.Get("master_instance_types").([]interface{}),
		"master_vswitch_ids":    d.Get("master_vswitch_ids").([]interface{}),
		// "num_of_nodes":                         d.Get("num_of_nodes").(int),
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
		"image_id":                             d.Get("image_id").(string),
		"image_type":                           "AliyunLinux3",
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

	defnodepool := make(map[string]interface{})
	defnodepool["nodepool_info"] = map[string]string{
		"name": "default-nodepool",
	}
	defnodepool["count"] = d.Get("num_of_nodes").(int)
	auto_scaling := map[string]interface{}{
		"enable": false,
	}
	defnodepool["auto_scaling"] = auto_scaling
	tee_config := map[string]interface{}{
		"tee_enable": false,
	}
	defnodepool["tee_config"] = tee_config
	scaling_group := map[string]interface{}{
		"platform":                      d.Get("platform").(string),
		"vpc_id":                        d.Get("vpc_id").(string),
		"vswitch_ids":                   d.Get("worker_vswitch_ids").([]interface{}),
		"instance_types":                d.Get("worker_instance_types").([]interface{}),
		"image_type":                    "AliyunLinux3",
		"system_disk_size":              d.Get("worker_disk_size").(int),
		"system_disk_category":          d.Get("worker_disk_category").(string),
		"system_disk_performance_level": d.Get("worker_system_disk_performance_level").(string),
		"internet_max_bandwidth_out":    0,
		"rds_instances":                 []string{},
	}

	if v, ok := d.GetOk("worker_disk_encrypted"); ok && v.(bool) {
		scaling_group["system_disk_encrypted"] = v.(bool)
		if v, ok := d.GetOk("worker_disk_encrypt_algorithm"); ok && v.(string) != "" {
			scaling_group["system_disk_encrypt_algorithm"] = v.(string)
		}
		if v, ok := d.GetOk("worker_disk_kms_key_id"); ok && v.(string) != "" {
			scaling_group["system_disk_kms_key_id"] = v.(string)
		}
	}
	if password, ok := d.GetOk("password"); ok && password.(string) != "" {
		body["login_Password"] = password.(string)
		scaling_group["login_Password"] = password.(string)
	} else if v, ok := d.GetOk("kms_encrypted_password"); v.(string) != "" && ok {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(v.(string), d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		password = decryptResp.Plaintext
		body["login_Password"] = password
		scaling_group["login_Password"] = password
	}
	if key_name, ok := d.GetOk("key_name"); ok && key_name != "" {
		body["key_pair"] = key_name.(string)
		scaling_group["key_pair"] = key_name.(string)
	}
	if data, ok := d.GetOk("worker_data_disks"); ok {
		data_disks := make([]map[string]interface{}, 0)
		// worker_data_disks := make([]map[string]interface{}, 0)
		for _, value := range data.([]interface{}) {
			disk := value.(map[string]interface{})
			data_disks = append(data_disks, map[string]interface{}{
				"size":                    disk["size"].(int),
				"category":                disk["category"].(string),
				"encrypted":               fmt.Sprintf("%t", disk["encrypted"].(bool)),
				"auto_snapshot_policy_id": disk["auto_snapshot_policy_id"].(string),
				"performance_level":       disk["performance_level"].(string),
				"kms_key_id":              disk["kms_key_id"].(string),
			})
			// v3.18.3 Sp01 does not support this parameter to encrypt worker nodes
			// worker_data_disks = append(worker_data_disks, map[string]interface{}{
			// 	"size":                    fmt.Sprintf("%d", disk["size"].(int)),
			// 	"category":                disk["category"].(string),
			// 	"encrypted":               fmt.Sprintf("%t", disk["encrypted"].(bool)),
			// 	"auto_snapshot_policy_id": disk["auto_snapshot_policy_id"].(string),
			// 	"performance_level":       disk["performance_level"].(string),
			// 	"kms_key_id":              disk["kms_key_id"].(string),
			// })
		}
		scaling_group["data_disks"] = data_disks
		// v3.18.3 Sp01 does not support this parameter to encrypt worker nodes
		// body["worker_data_disks"] = worker_data_disks
	}
	defnodepool["scaling_group"] = scaling_group
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
	kubernetes_config := map[string]interface{}{
		"cpu_policy":    d.Get("cpu_policy").(string),
		"cms_enabled":   false,
		"unschedulable": false,
		"labels":        []string{},
	}
	if v, ok := d.GetOk("runtime"); ok && len(v.([]interface{})) > 0 {
		all, _ := v.([]interface{})
		runtime := all[0].(map[string]interface{})
		body["runtime"] = runtime
		kubernetes_config["runtime"] = runtime["name"].(string)
		kubernetes_config["runtime_version"] = runtime["version"].(string)
		defnodepool["kubernetes_config"] = kubernetes_config
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

	if v, ok := d.GetOk("master_disk_encrypted"); ok && v.(bool) {
		body["master_system_disk_encrypted"] = fmt.Sprintf("%t", v.(bool))
		if v, ok := d.GetOk("master_disk_encrypt_algorithm"); ok && v.(string) != "" {
			body["master_system_disk_encrypt_algorithm"] = v.(string)
		}
		if v, ok := d.GetOk("master_disk_kms_key_id"); ok && v.(string) != "" {
			body["master_system_disk_kms_key_id"] = v.(string)
		}
	}

	// v3.18.3 Sp01 does not support this parameter to encrypt worker nodes
	// if v, ok := d.GetOk("worker_disk_encrypted"); ok && v.(bool) {
	// 	body["worker_system_disk_encrypted"] = fmt.Sprintf("%t", v.(bool))
	// 	if v, ok := d.GetOk("worker_disk_encrypt_algorithm"); ok && v.(string) != "" {
	// 		body["worker_system_disk_encrypt_algorithm"] = v.(string)
	// 	}
	// 	if v, ok := d.GetOk("worker_disk_kms_key_id"); ok && v.(string) != "" {
	// 		body["worker_system_disk_kms_key_id"] = v.(string)
	// 	}
	// }

	if v, ok := d.GetOk("instances"); ok {
		body["format_disk"] = d.Get("format_disk").(bool)
		body["keep_instance_name"] = d.Get("keep_instance_name").(bool)
		body["instances"] = expandStringList(v.(*schema.Set).List())
	} else {
		if d.Get("num_of_nodes").(int) > 0 {
			body["nodepools"] = []interface{}{defnodepool}
		}
		// v3.18.3 Sp01 does not support this parameter to encrypt worker nodes
		// body["worker_instance_types"] = d.Get("worker_instance_types").([]interface{})
		// body["worker_vswitch_ids"] = d.Get("worker_vswitch_ids").([]interface{})
		// body["worker_system_disk_category"] = d.Get("worker_disk_category").(string)
		// body["worker_system_disk_size"] = d.Get("worker_disk_size").(int)
		// body["worker_system_disk_performance_level"] = d.Get("worker_system_disk_performance_level").(string)
		// body["master_storage_set_id"] = d.Get("master_storage_set_id").(string)
		// body["master_storage_set_partition_number"] = d.Get("master_storage_set_partition_number").(int)
		// body["worker_storage_set_id"] = d.Get("worker_storage_set_id").(string)
		// body["worker_storage_set_partition_number"] = d.Get("worker_storage_set_partition_number").(int)
	}
	log.Printf("[DEBUG] Request body: %v", body)
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
	invoker := NewInvoker()

	var nodepoolid string
	if d.IsNewResource() {
		return nil
	}
	noUpdatesAllowedCheck(d, []string{
		"name", "master_disk_size", "master_disk_category", "master_disk_encrypt_algorithm",
		"master_disk_kms_key_id", "master_disk_encrypted", "delete_protection",
		"worker_disk_size", "worker_disk_category", "worker_disk_encrypt_algorithm",
		"worker_disk_kms_key_id", "worker_disk_encrypted", "worker_data_disks",
		"pod_vswitch_ids", "pod_cidr", "service_cidr", "node_cidr_mask",
		"new_nat_gateway", "enable_ssh", "node_port_range", "image_id",
		"version", "cluster_type", "os_type", "platform", "cpu_policy",
		"proxy_mode", "addons", "slb_internet_enabled", "master_instance_types",
		"master_vswitch_ids", "worker_instance_types", "worker_vswitch_ids",
		"instances", "format_disk", "keep_instance_name", "master_count",
		"timeout_mins", "nodes", "user_data", "cloud_monitor_flags", "runtime",
		"is_enterprise_security_group", "security_group_id",
		"master_system_disk_performance_level", "worker_system_disk_performance_level",
		"master_storage_set_id", "master_storage_set_partition_number",
		"worker_storage_set_id", "worker_storage_set_partition_number",
	})
	nodepoolid = d.Get("nodepool_id").(string)
	resourceNodepoolId := fmt.Sprintf("%s:%s", d.Id(), nodepoolid)

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
			stateConf := BuildStateConf([]string{"removing"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(resourceNodepoolId, []string{"deleting", "failed"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
			}
			if newValue == 0 {
				// If the number of worker nodes is 0, delete the default node pool
				stateConf := BuildStateConf([]string{"initial", "removing"}, []string{"active", "running"}, d.Timeout(schema.TimeoutUpdate), 15*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"failed"}))
				stateConf.NotFoundChecks = 1000
				if _, err := stateConf.WaitForState(); err != nil {
					return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
				}
				req := client.NewCommonRequest("DELETE", "CS", "2015-12-15", "DeleteClusterNodepool", fmt.Sprintf("/clusters/%s/nodepools/%s", d.Id(), nodepoolid))
				req.QueryParams["ClusterId"] = d.Id()
				req.QueryParams["NodepoolId"] = nodepoolid
				req.QueryParams["force"] = "true"
				req.QueryParams["Force"] = "true"
				body := map[string]interface{}{
					"force": true,
					"Force": true,
				}
				jsonData, err := json.Marshal(body)
				if err != nil {
					return errmsgs.WrapError(fmt.Errorf("Error marshaling to JSON: %v", err))
				}
				req.SetContentType(requests.Json)
				req.SetContent(jsonData)

				response, err := client.ProcessCommonRequest(req)
				if err != nil {
					if response == nil {
						return errmsgs.WrapErrorf(err, "Process Common Request Failed")
					}
					errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
					return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, nodepoolid, "DeleteClusterNodePool", errmsg)
				}
				stateConf = BuildStateConf([]string{"deleting", "active"}, []string{}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(resourceNodepoolId, []string{"failed"}))
				if _, err = stateConf.WaitForState(); err != nil {
					return errmsgs.WrapErrorf(err, errmsgs.IdMsg, nodepoolid)
				}
			}
		}

		if newValue > oldValue {
			if oldValue == 0 {
				// Cannot create default node pool when num_of_nodes was 0
				return errmsgs.WrapErrorf(fmt.Errorf("num_of_nodes can not be greater than before"), "scaleOutFailed %d:%d", newValue, oldValue)
			}
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
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_cs_kubernetes", "ScaleClusterNodePool", resp)
			}
			addDebug("ScaleClusterNodePool", resp, request, request.QueryParams)

			stateConf := BuildStateConf([]string{"scaling"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(resourceNodepoolId, []string{"deleting", "failed"}))

			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
			}
		}
	}
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
	d.Set("name", object.Name)
	d.Set("vpc_id", object.VpcId)
	d.Set("pod_cidr", object.ContainerCIDR)
	d.Set("version", object.InitVersion)
	d.Set("cluster_type", string(object.ClusterType))
	d.Set("security_group_id", object.SecurityGroupId)
	d.Set("delete_protection", object.DeletionProtection)
	d.Set("worker_ram_role_name", object.WorkerRamRoleName)

	// node_count, err := csService.GetCsK8sNodesCount(d.Id())
	// if err != nil {
	// 	return errmsgs.WrapError(err)
	// }
	if nodepoolid != "" {
		nodepool, err := csService.DescribeCsKubernetesNodePool(fmt.Sprintf("%s:%s", d.Id(), nodepoolid))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		d.Set("num_of_nodes", nodepool.Status.TotalNodes)

		// Read worker-related configurations from nodepool
		if nodepool.ScalingGroup.InstanceTypes != nil {
			d.Set("worker_instance_types", nodepool.ScalingGroup.InstanceTypes)
		}
		if nodepool.ScalingGroup.VswitchIds != nil {
			d.Set("worker_vswitch_ids", nodepool.ScalingGroup.VswitchIds)
		}
		if nodepool.ScalingGroup.SystemDiskCategory != "" {
			d.Set("worker_disk_category", nodepool.ScalingGroup.SystemDiskCategory)
		}
		if nodepool.ScalingGroup.SystemDiskSize > 0 {
			d.Set("worker_disk_size", nodepool.ScalingGroup.SystemDiskSize)
		}
		if nodepool.ScalingGroup.KeyPair != "" {
			d.Set("key_name", nodepool.ScalingGroup.KeyPair)
		}

		// Read configurations from kubernetes_config
		if nodepool.KubernetesConfig.CPUPolicy != "" {
			d.Set("cpu_policy", nodepool.KubernetesConfig.CPUPolicy)
		}
		d.Set("cloud_monitor_flags", nodepool.KubernetesConfig.CmsEnabled)
		if nodepool.KubernetesConfig.UserData != "" {
			d.Set("user_data", nodepool.KubernetesConfig.UserData)
		}
		// Read runtime configuration
		// Prefer nodepool.KubernetesConfig for runtime details as it's more specific to the node pool
		runtimeName := nodepool.KubernetesConfig.Runtime
		runtimeVersion := nodepool.KubernetesConfig.RuntimeVersion

		// Fallback to object.MetaData if nodepool doesn't have runtime info
		if runtimeName == "" || runtimeVersion == "" {
			if object.MetaData != "" {
				var metaDataMap map[string]interface{}
				if err := json.Unmarshal([]byte(object.MetaData), &metaDataMap); err == nil {
					if rn, ok := metaDataMap["Runtime"].(string); ok && runtimeName == "" {
						runtimeName = rn
					}
					if rv, ok := metaDataMap["RuntimeVersion"].(string); ok && runtimeVersion == "" {
						runtimeVersion = rv
					}
				}
			}
		}

		if runtimeName != "" || runtimeVersion != "" {
			runtime := []map[string]interface{}{
				{
					"name":    runtimeName,
					"version": runtimeVersion,
				},
			}
			d.Set("runtime", runtime)
		}
		sworker := make([]map[string]interface{}, 0)
		clusternode, err := csService.DescribeClusterNodes(d.Id(), nodepoolid)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		for _, k := range clusternode.Nodes {
			if k.InstanceRole == "Worker" && k.InstanceStatus == "Running" {
				WorkerNodes := map[string]interface{}{
					"id":         k.InstanceID,
					"name":       k.InstanceName,
					"private_ip": fmt.Sprintf("%s", k.IPAddress),
				}
				sworker = append(sworker, WorkerNodes)
			}
		}
		d.Set("worker_nodes", sworker)
	} else {
		d.Set("num_of_nodes", 0)
		d.Set("worker_instance_types", nil)
		d.Set("worker_vswitch_ids", nil)
		d.Set("worker_disk_category", nil)
		d.Set("worker_disk_size", nil)
		d.Set("worker_nodes", nil)
	}
	smaster := make([]map[string]interface{}, 0)
	masternodes, err := csService.DescribeClusterMasterNodes(d.Id())
	for _, k := range masternodes {
		MasterNodes := map[string]interface{}{
			"id":         k.InstanceID,
			"name":       k.InstanceName,
			"private_ip": fmt.Sprintf("%s", k.IPAddress),
		}
		smaster = append(smaster, MasterNodes)
	}

	d.Set("master_nodes", smaster)
	if err := d.Set("tags", flattenTagsConfig(object.Tags)); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func resourceAlibabacloudStackCSKubernetesDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
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
		if errmsgs.IsExpectedErrors(err, "ErrorClusterNotFound") {
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
			for _, npinfo := range nodepool.Nodepools {
				if npinfo.NodepoolInfo.Name == "default-nodepool" {
					nodepoolid = npinfo.NodepoolInfo.NodepoolID
					break
				}
			}
		}
		// if nodepoolid == "" {
		// 	return resource.RetryableError(errmsgs.WrapErrorf(fmt.Errorf("can not found default node_pool"), "DescribeClusterNodePools", nodepool.Nodepools))
		// }
		return nil
	}); err != nil {
		return "", err
	}
	return nodepoolid, nil
}

func updateKubernetesClusterTag(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	csService := CsService{client}
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
