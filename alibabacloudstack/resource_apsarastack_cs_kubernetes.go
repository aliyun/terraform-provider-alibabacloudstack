package alibabacloudstack

import (
	//	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"

	//	"github.com/aliyun/alibaba-cloud-sdk-go/services/cs"
	"regexp"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
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
	return &schema.Resource{
		Create: resourceAlibabacloudStackCSKubernetesCreate,
		Read:   resourceAlibabacloudStackCSKubernetesRead,
		Update: resourceAlibabacloudStackCSKubernetesUpdate,
		Delete: resourceAlibabacloudStackCSKubernetesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
				Type:     schema.TypeInt,
				Required: true,
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
			"worker_data_disk_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          40,
				ValidateFunc:     validation.IntBetween(20, 32768),
				DiffSuppressFunc: workerDataDiskSizeSuppressFunc,
			},
			"worker_data_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD), string(DiskCloudPPERF), string(DiskCloudSPERF)}, false),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
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
			"exclude_autoscaler_nodes": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
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
				ConflictsWith:    []string{"key_name", "kms_encrypted_password"},
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
			"user_ca": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
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
			"install_cloud_monitor": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
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
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_server_internet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_server_intranet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"slb_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field 'slb_id' has been deprecated from provider version 1.9.2. New field 'slb_internet' replaces it.",
			},
			"slb_internet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slb_intranet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"is_enterprise_security_group"},
			},
			"master_system_disk_performance_level": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"is_enterprise_security_group"},
			},
			"worker_system_disk_performance_level": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_enterprise_security_group": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"security_group_id"},
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
			"worker_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
				ConflictsWith: []string{"worker_instance_types", "worker_vswitch_ids", "worker_data_disk_category", "worker_disk_category"},
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
			"vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems:         3,
				MaxItems:         5,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
				//Removed:          "Field 'vswitch_ids' has been removed from provider version 1.75.0. New field 'master_vswitch_ids' and 'worker_vswitch_ids' replace it.",
			},
			"master_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			// single instance type would cause extra troubles
			"master_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// force update is a high risk operation
			"force_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				//Removed:  "Field 'force_update' has been removed from provider version 1.75.0.",
			},
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
				//Removed:  "Field 'nodes' has been removed from provider version 1.9.4. New field 'master_nodes' replaces it.",
			},
			// too hard to use this config
			"log_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{KubernetesClusterLoggingTypeSLS}, false),
							Required:     true,
						},
						"project": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				DiffSuppressFunc: csForceUpdateSuppressFunc,
				//Removed:          "Field 'log_config' has been removed from provider version 1.75.0. New field 'addons' replaces it.",
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_name_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^customized,[a-z0-9]([-a-z0-9\.])*,([5-9]|[1][0-2]),([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`), "Each node name consists of a prefix, an IP substring, and a suffix. For example, if the node IP address is 192.168.0.55, the prefix is aliyun.com, IP substring length is 5, and the suffix is test, the node name will be aliyun.com00055test."),
			},
			"worker_ram_role_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_account_issuer": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"api_audiences": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
			},
			"nodepool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
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
	var requestInfo *cs.Client
	var raw interface{}

	timeout := d.Get("timeout_mins").(int)
	Name := d.Get("name").(string)
	OsType := d.Get("os_type").(string)
	Platform := d.Get("platform").(string)
	mastercount := d.Get("master_count").(int)
	msysdiskcat := d.Get("master_disk_category").(string)
	msysdisksize := d.Get("master_disk_size").(int)
	wsysdisksize := d.Get("worker_disk_size").(int)
	masterstoragesetid := d.Get("master_storage_set_id").(string)
	masterstoragesetnumber := d.Get("master_storage_set_partition_number").(int)
	workerstoragesetid := d.Get("worker_storage_set_id").(string)
	workerstoragesetnumber := d.Get("worker_storage_set_partition_number").(int)
	wsysdiskcat := d.Get("worker_disk_category").(string)
	delete_pro := d.Get("delete_protection").(bool)
	KubernetesVersion := d.Get("version").(string)
	IsEnterpriseSecurityGroup := d.Get("is_enterprise_security_group").(bool)
	addons := make([]cs.Addon, 0)
	type WorkerData struct {
		Size                 int
		Encrypted            bool
		AutoSnapshotPolicyId string
		PerformanceLevel     string
		Category             string
	}
	var req, workerdisks string
	var pod, attachinst int
	if v, ok := d.GetOk("addons"); ok {
		all, ok := v.([]interface{})
		if ok {
			for i, a := range all {
				addon, _ := a.(map[string]interface{})
				log.Printf("check addon %v", addon)
				if addon["name"] == "terway-eniip" {
					pod = 1
				}
				log.Printf("check req id %v", i)
				if i == 0 {
					req = fmt.Sprintf("{\"name\" : \"%s\",\"config\": \"%s\"}", addon["name"].(string), addon["config"].(string))
				} else {
					req = fmt.Sprintf("%s,{\"name\" : \"%s\",\"config\": %q}", req, addon["name"].(string), addon["config"].(string))
				}
			}
		}
	}
	if v, ok := d.GetOk("worker_data_disks"); ok {
		all, ok := v.([]interface{})
		if ok {
			for i, a := range all {
				disk, _ := a.(map[string]interface{})
				if i == 0 {
					workerdisks = fmt.Sprintf("{\"size\" : \"%d\",\"encrypted\": \"%t\",\"performance_level\": \"%s\",\"auto_snapshot_policy_id\": \"%s\",\"category\": \"%s\"}", disk["size"].(int), disk["encrypted"].(bool), disk["performance_level"].(string), disk["auto_snapshot_policy_id"].(string), disk["category"].(string))
				} else {
					workerdisks = fmt.Sprintf("%s,{\"size\" : \"%d\",\"encrypted\": \"%t\",\"performance_level\": \"%s\",\"auto_snapshot_policy_id\": \"%s\",\"category\": \"%s\"}", req, disk["size"].(int), disk["encrypted"].(bool), disk["performance_level"].(string), disk["auto_snapshot_policy_id"].(string), disk["category"].(string))
				}
				log.Printf("checking workerdatadisks %v", workerdisks)

			}
		}
	}
	udata := d.Get("user_data").(string)
	log.Printf("checking addons %v", addons)
	log.Printf("check req final %s", req)
	var runtime string
	if v, ok := d.GetOk("runtime"); ok {
		all, _ := v.([]interface{})
		for _, a := range all {
			run, _ := a.(map[string]interface{})
			runtime = fmt.Sprintf("\"name\": \"%s\", \"version\": \"%s\"", run["name"].(string), run["version"].(string))
		}
	}
	log.Printf("checking runtime %v", runtime)
	var tags string
	tagss := make([]interface{}, 0)
	if v, ok := d.GetOk("tags"); ok && len(v.(map[string]interface{})) > 0 {
		for key, value := range v.(map[string]interface{}) {
			tagss = append(tagss, cs.Tag{
				Key:   key,
				Value: value.(string),
			})
		}
	}
	tagsBytes, _ := json.Marshal(tagss)
	tags = string(tagsBytes)
	log.Printf("checking tags %v", tags)
	proxy_mode := d.Get("proxy_mode").(string)
	VpcId := d.Get("vpc_id").(string)
	ImageId := d.Get("image_id").(string)
	var LoginPassword string
	if password := d.Get("password").(string); password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			password = decryptResp.Plaintext
		}
		LoginPassword = password
	} else {
		LoginPassword = password
	}
	nodecidr := d.Get("node_cidr_mask").(string)
	enabSsh := d.Get("enable_ssh").(bool)
	end := d.Get("slb_internet_enabled").(bool)
	SnatEntry := d.Get("new_nat_gateway").(bool)
	scdir := d.Get("service_cidr").(string)
	pcidr := d.Get("pod_cidr").(string)
	NumOfNodes := int64(d.Get("num_of_nodes").(int))
	MasterSystemDiskPerformanceLevel := d.Get("master_system_disk_performance_level").(string)
	WorkerSystemDiskPerformanceLevel := d.Get("worker_system_disk_performance_level").(string)
	CloudMonitorFlags := d.Get("cloud_monitor_flags").(bool)
	var secgroup string
	var SecurityGroup bool
	if _, ok := d.GetOk("is_enterprise_security_group"); ok {
		secgroup = "is_enterprise_security_group"
		SecurityGroup = d.Get("is_enterprise_security_group").(bool)

	} else if _, ok := d.GetOk("security_group_id"); ok {
		secgroup = "security_group_id"
		SecurityGroup = d.Get("security_group_id").(bool)
	}

	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	var wvid, mvid, winst, minst, podid, inst string
	var formatDisk, retainIname bool

	wvids := d.Get("worker_vswitch_ids").([]interface{})
	for i, k := range wvids {
		if i == 0 {
			wvid = fmt.Sprintf("%s", k)
		} else {
			wvid = fmt.Sprintf("%s\",\"%s", wvid, k)
		}
	}
	log.Printf("new worker vids %v ", wvid)
	mvids := d.Get("master_vswitch_ids").([]interface{})
	for i, k := range mvids {
		if i == 0 {
			mvid = fmt.Sprintf("%s", k)
		} else {
			mvid = fmt.Sprintf("%s\",\"%s", mvid, k)
		}
	}
	log.Printf("master vswids %v", mvid)
	winsts := d.Get("worker_instance_types").([]interface{})
	for i, k := range winsts {
		if i == 0 {
			winst = fmt.Sprintf("%s", k)
		} else {
			winst = fmt.Sprintf("%s\",\"%s", winst, k)
		}

	}
	log.Printf("new worker inst %v ", winst)
	insrsas := d.Get("master_instance_types").([]interface{})
	for i, k := range insrsas {
		if i == 0 {
			minst = fmt.Sprintf("%s", k)
		} else {
			minst = fmt.Sprintf("%s\",\"%s", minst, k)
		}
		log.Printf("instances %d %v", i, k)
		//minst = fmt.Sprintf("%s\",\"%s", minst, k)
	}
	//log.Printf("new master inst %v ",insrsas)
	log.Printf("new master inst %v ", minst)

	var insts, podids []string
	if v, ok := d.GetOk("instances"); ok {
		attachinst = 1
		formatDisk = d.Get("format_disk").(bool)
		retainIname = d.Get("keep_instance_name").(bool)
		insts = expandStringList(v.(*schema.Set).List())
		fmt.Print("checking instances attached: ", insts)
		for i, k := range insts {
			if i != 0 {
				inst = fmt.Sprintf("%s\",\"%s", inst, k)

			} else {
				inst = k
			}
		}
	}
	log.Printf("pod request %d", pod)
	clustertype := d.Get("cluster_type").(string)
	log.Printf("pod is %d", pod)
	if pod == 1 {
		if v, ok := d.GetOk("pod_vswitch_ids"); ok {
			log.Printf("123123podid is %v\n", podid)
			podids = expandStringList(v.(*schema.Set).List())
			log.Print("checking pod vsw ids: ", podids)
			for i, k := range podids {
				if i != 0 {
					podid = fmt.Sprintf("%s\",\"%s", podid, k)
					//minst=strings.Join(minsts,",")

				} else {
					podid = k
				}
			}
		}
	}
	nodeportrange := d.Get("node_port_range").(string)
	cpuPolicy := d.Get("cpu_policy").(string)
	log.Printf("wswitchids %v mswitchids %v", wvid, mvid)
	log.Printf("winsts %v minsts %v", winst, minst)
	if attachinst == 1 {
		if pod == 0 {
			request.QueryParams = map[string]string{
				"RegionId":         client.RegionId,
				"AccessKeySecret":  client.SecretKey,
				"Product":          "Cs",
				"Department":       client.Department,
				"ResourceGroup":    client.ResourceGroup,
				"Action":           "CreateCluster",
				"AccountInfo":      "123456",
				"Version":          "2015-12-15",
				"SignatureVersion": "1.0",
				"ProductName":      "cs",
				"X-acs-body": fmt.Sprintf("{\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":%d,\"%s\":%t,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":[\"%s\"],\"%s\":[\"%s\"],\"%s\":\"%s\",\"%s\":%d,\"%s\":%d,\"%s\":%t,\"%s\":%t,\"%s\":%t,\"%s\":\"%s\",\"%s\":%d,\"%s\":%t,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":[%s],\"%s\":\"%s\",\"%s\":[\"%s\"],\"%s\":%t,\"%s\":%t,\"%s\":\"%s\",\"%s\":{%s},\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":[%s],\"%s\":%t,\"%s\":%t,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":%t,\"%s\":\"%s\",\"%s\":%s}",
					"Product", "Cs",
					"os_type", OsType,
					"platform", Platform,
					"cluster_type", clustertype,
					"region_id", client.RegionId,
					"timeout_mins", timeout,
					"disable_rollback", true,
					"kubernetes_version", KubernetesVersion,
					"container_cidr", pcidr,
					"service_cidr", scdir,
					"name", Name,
					"master_instance_types", minst,
					"master_vswitch_ids", mvid,
					"login_Password", LoginPassword,
					"num_of_nodes", NumOfNodes,
					"master_count", mastercount,
					"snat_entry", SnatEntry,
					"endpoint_public_access", end,
					"ssh_flags", enabSsh,
					"master_system_disk_category", msysdiskcat,
					"master_system_disk_size", msysdisksize,
					"deletion_protection", delete_pro,
					"node_cidr_mask", nodecidr,
					"vpcid", VpcId,
					"addons", req,
					"proxy_mode", proxy_mode,
					"instances", inst,
					"format_disk", formatDisk,
					"keep_instance_name", retainIname,
					"user_data", udata,
					"runtime", runtime,
					"node_port_range", nodeportrange,
					"cpu_policy", cpuPolicy,
					"worker_data_disks", workerdisks,

					secgroup, SecurityGroup,
					"cloud_monitor_flags", CloudMonitorFlags,
					"master_system_disk_performance_level", MasterSystemDiskPerformanceLevel,
					"worker_system_disk_performance_level", WorkerSystemDiskPerformanceLevel,
					"is_enterprise_security_group", IsEnterpriseSecurityGroup,
					"image_id", ImageId,
					"tags", tags,
				),
			}
		} else {
			request.QueryParams = map[string]string{
				"RegionId":         client.RegionId,
				"AccessKeySecret":  client.SecretKey,
				"Product":          "Cs",
				"Department":       client.Department,
				"ResourceGroup":    client.ResourceGroup,
				"Action":           "CreateCluster",
				"AccountInfo":      "123456",
				"Version":          "2015-12-15",
				"SignatureVersion": "1.0",
				"ProductName":      "cs",
				"X-acs-body": fmt.Sprintf("{\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":%d,\"%s\":%t,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":[\"%s\"],\"%s\":[\"%s\"],\"%s\":\"%s\",\"%s\":%d,\"%s\":%d,\"%s\":%t,\"%s\":%t,\"%s\":%t,\"%s\":\"%s\",\"%s\":%d,\"%s\":%t,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":[%s],\"%s\":[\"%s\"],\"%s\":\"%s\",\"%s\":[\"%s\"],\"%s\":%t,\"%s\":%t,\"%s\":\"%s\",\"%s\":{%s},\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":[%s],\"%s\":%t,\"%s\":%t,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":%t,\"%s\":\"%s\",\"%s\":%s}",
					"Product", "Cs",
					"os_type", OsType,
					"platform", Platform,
					"cluster_type", clustertype,
					"region_id", client.RegionId,
					"timeout_mins", timeout,
					"disable_rollback", true,
					"kubernetes_version", KubernetesVersion,
					"container_cidr", pcidr,
					"service_cidr", scdir,
					"name", Name,
					"master_instance_types", minst,
					"master_vswitch_ids", mvid,
					"login_Password", LoginPassword,
					"num_of_nodes", NumOfNodes,
					"master_count", mastercount,
					"snat_entry", SnatEntry,
					"endpoint_public_access", end,
					"ssh_flags", enabSsh,
					"master_system_disk_category", msysdiskcat,
					"master_system_disk_size", msysdisksize,
					"deletion_protection", delete_pro,
					"node_cidr_mask", nodecidr,
					"vpcid", VpcId,
					"addons", req,
					"pod_vswitch_ids", podid,
					"proxy_mode", proxy_mode,
					"instances", inst,
					"format_disk", formatDisk,
					"keep_instance_name", retainIname,
					"user_data", udata,
					"runtime", runtime,
					"node_port_range", nodeportrange,
					"cpu_policy", cpuPolicy,
					"worker_data_disks", workerdisks,
					secgroup, SecurityGroup,
					"cloud_monitor_flags", CloudMonitorFlags,
					"master_system_disk_performance_level", MasterSystemDiskPerformanceLevel,
					"worker_system_disk_performance_level", WorkerSystemDiskPerformanceLevel,
					"is_enterprise_security_group", IsEnterpriseSecurityGroup,
					"image_id", ImageId,
					"tags", tags,
				),
			}
		}
	} else {
		if pod == 0 {
			request.QueryParams = map[string]string{
				"RegionId":         client.RegionId,
				"AccessKeySecret":  client.SecretKey,
				"Product":          "Cs",
				"Department":       client.Department,
				"ResourceGroup":    client.ResourceGroup,
				"Action":           "CreateCluster",
				"AccountInfo":      "123456",
				"Version":          "2015-12-15",
				"SignatureVersion": "1.0",
				"ProductName":      "cs",
				"X-acs-body": fmt.Sprintf("{\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":%d,\"%s\":%t,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":[\"%s\"],\"%s\":[\"%s\"],\"%s\":[\"%s\"],\"%s\":[\"%s\"],\"%s\":\"%s\",\"%s\":%d,\"%s\":%d,\"%s\":%t,\"%s\":%t,\"%s\":%t,\"%s\":\"%s\",\"%s\":%d,\"%s\":\"%s\",\"%s\":%d,\"%s\":%t,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":[%s],\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":{%s},\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":%t,\"%s\":%t,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":[%s],\"%s\":%t,\"%s\":\"%s\",\"%s\":%s,\"%s\":\"%s\",\"%s\":%d,\"%s\":\"%s\",\"%s\":%d}",
					"Product", "Cs",
					"os_type", OsType,
					"platform", Platform,
					"cluster_type", clustertype,
					"region_id", client.RegionId,
					"timeout_mins", timeout,
					"disable_rollback", true,
					"kubernetes_version", KubernetesVersion,
					"container_cidr", pcidr,
					"service_cidr", scdir,
					"name", Name,
					"master_instance_types", minst,
					"worker_instance_types", winst,
					"master_vswitch_ids", mvid,
					"worker_vswitch_ids", wvid,
					"login_Password", LoginPassword,
					"num_of_nodes", NumOfNodes,
					"master_count", mastercount,
					"snat_entry", SnatEntry,
					"endpoint_public_access", end,
					"ssh_flags", enabSsh,
					"master_system_disk_category", msysdiskcat,
					"master_system_disk_size", msysdisksize,
					"worker_system_disk_category", wsysdiskcat,
					"worker_system_disk_size", wsysdisksize,
					"deletion_protection", delete_pro,
					"node_cidr_mask", nodecidr,
					"vpcid", VpcId,
					"addons", req,
					"proxy_mode", proxy_mode,
					"user_data", udata,
					"runtime", runtime,
					"node_port_range", nodeportrange,
					"cpu_policy", cpuPolicy,
					secgroup, SecurityGroup,
					"cloud_monitor_flags", CloudMonitorFlags,
					"master_system_disk_performance_level", MasterSystemDiskPerformanceLevel,
					"worker_system_disk_performance_level", WorkerSystemDiskPerformanceLevel,
					"worker_data_disks", workerdisks,
					"is_enterprise_security_group", IsEnterpriseSecurityGroup,
					"image_id", ImageId,
					"tags", tags,
					"master_storage_set_id", masterstoragesetid,
					"master_storage_set_partition_number", masterstoragesetnumber,
					"worker_storage_set_id", workerstoragesetid,
					"worker_storage_set_partition_number", workerstoragesetnumber,
				),
			}
		} else {
			request.QueryParams = map[string]string{
				"RegionId":         client.RegionId,
				"AccessKeySecret":  client.SecretKey,
				"Product":          "Cs",
				"Department":       client.Department,
				"ResourceGroup":    client.ResourceGroup,
				"Action":           "CreateCluster",
				"AccountInfo":      "123456",
				"Version":          "2015-12-15",
				"SignatureVersion": "1.0",
				"ProductName":      "cs",
				"X-acs-body": fmt.Sprintf("{\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":%d,\"%s\":%t,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":[\"%s\"],\"%s\":[\"%s\"],\"%s\":[\"%s\"],\"%s\":[\"%s\"],\"%s\":\"%s\",\"%s\":%d,\"%s\":%d,\"%s\":%t,\"%s\":%t,\"%s\":%t,\"%s\":\"%s\",\"%s\":%d,\"%s\":\"%s\",\"%s\":%d,\"%s\":%t,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":[%s],\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":{%s},\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":%t,\"%s\":%t,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":[%s],\"%s\":[\"%s\"],\"%s\":%t,\"%s\":\"%s\",\"%s\":%s,\"%s\":\"%s\",\"%s\":%d,\"%s\":\"%s\",\"%s\":%d}",
					"Product", "Cs",
					"os_type", OsType,
					"platform", Platform,
					"cluster_type", clustertype,
					"region_id", client.RegionId,
					"timeout_mins", timeout,
					"disable_rollback", true,
					"kubernetes_version", KubernetesVersion,
					"container_cidr", pcidr,
					"service_cidr", scdir,
					"name", Name,
					"master_instance_types", minst,
					"worker_instance_types", winst,
					"master_vswitch_ids", mvid,
					"worker_vswitch_ids", wvid,
					"login_Password", LoginPassword,
					"num_of_nodes", NumOfNodes,
					"master_count", mastercount,
					"snat_entry", SnatEntry,
					"endpoint_public_access", end,
					"ssh_flags", enabSsh,
					"master_system_disk_category", msysdiskcat,
					"master_system_disk_size", msysdisksize,
					"worker_system_disk_category", wsysdiskcat,
					"worker_system_disk_size", wsysdisksize,
					"deletion_protection", delete_pro,
					"node_cidr_mask", nodecidr,
					"vpcid", VpcId,
					"addons", req,
					"proxy_mode", proxy_mode,
					"user_data", udata,
					"runtime", runtime,
					"node_port_range", nodeportrange,
					"cpu_policy", cpuPolicy,
					secgroup, SecurityGroup,
					"cloud_monitor_flags", CloudMonitorFlags,
					"master_system_disk_performance_level", MasterSystemDiskPerformanceLevel,
					"worker_system_disk_performance_level", WorkerSystemDiskPerformanceLevel,
					"worker_data_disks", workerdisks,
					"pod_vswitch_ids", podid,
					"is_enterprise_security_group", IsEnterpriseSecurityGroup,
					"image_id", ImageId,
					"tags", tags,
					"master_storage_set_id", masterstoragesetid,
					"master_storage_set_partition_number", masterstoragesetnumber,
					"worker_storage_set_id", workerstoragesetid,
					"worker_storage_set_partition_number", workerstoragesetnumber,
				),
			}
		}
	}
	request.Method = "POST"        // Set request method
	request.Product = "CS"         // Specify product
	request.Version = "2015-12-15" // Specify product version
	request.ServiceCode = "cs"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	} // Set request scheme. Default: http
	request.ApiName = "CreateCluster"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.Headers = map[string]string{"x-acs-asapi-gateway-version": "3.0"}

	var err error
	err = nil
	if err = invoker.Run(func() error {
		raw, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		return err
	}); err != nil {
		//return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_cs_kubernetes", "CreateKubernetesCluster", raw)
		return err
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["RegionId"] = common.Region(client.RegionId)
		requestMap["Params"] = request.GetQueryParams()
		addDebug("CreateKubernetesCluster", raw, requestInfo, requestMap)
	}

	if err != nil {
		//return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_cs_kubernetes", request.GetActionName(), AlibabacloudStackSdkGoERROR)
		return err
	}

	if debugOn() {
		addDebug("CreateKubernetesCluster", raw, request)
	}
	clusterresponse := ClusterCommonResponse{}
	cluster, _ := raw.(*responses.CommonResponse)
	if cluster.IsSuccess() == false {
		//return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm", "API Action", cluster.GetHttpContentString())
		return err
	}
	ok := json.Unmarshal(cluster.GetHttpContentBytes(), &clusterresponse)
	if ok != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_cs_kubernetes", "ParseKubernetesClusterResponse", raw)
	}
	d.SetId(clusterresponse.ClusterID)

	stateConf := BuildStateConf([]string{"initial", " "}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 15*time.Minute, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlibabacloudStackCSKubernetesUpdate(d, meta)
}

func resourceAlibabacloudStackCSKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	csService := CsService{client}
	d.Partial(true)
	var raw interface{}
	invoker := NewInvoker()

	nodepool, err := csService.DescribeClusterNodePools(d.Id())
	if err != nil {
		return WrapError(err)
	}
	var nodepoolid string
	for _, k := range nodepool.Nodepools {
		//Considering multiple nodepools
		if k.NodepoolInfo.IsDefault {
			nodepoolid = k.NodepoolInfo.NodepoolID
		}
	}
	if nodepoolid == "" {
		return WrapErrorf(fmt.Errorf("can not found default node_pool"), "DescribeClusterNodePools", nodepool.Nodepools)
	}
	d.Set("nodepool_id", nodepoolid)
	if d.HasChange("num_of_nodes") && !d.IsNewResource() {
		password := d.Get("password").(string)
		if password == "" {
			if v := d.Get("kms_encrypted_password").(string); v != "" {
				kmsService := KmsService{client}
				decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
				if err != nil {
					return WrapError(err)
				}
				password = decryptResp.Plaintext
			}
		}

		oldV, newV := d.GetChange("num_of_nodes")
		oldValue, ok := oldV.(int)
		if ok != true {
			return WrapErrorf(fmt.Errorf("num_of_nodes old value can not be parsed"), "parseError %d", oldValue)
		}
		newValue, ok := newV.(int)
		if ok != true {
			return WrapErrorf(fmt.Errorf("num_of_nodes new value can not be parsed"), "parseError %d", newValue)
		}

		if newValue < oldValue {
			//return WrapErrorf(fmt.Errorf("num_of_nodes can not be less than before"), "scaleOutFailed %d:%d", newValue, oldValue)
			object, err := csService.DescribeClusterNodes(d.Id(), nodepoolid)
			if err != nil {
				if NotFoundError(err) {
					d.SetId("")
					return nil
				}
				return WrapError(err)
			}
			var allNodeName []string
			for _, value := range object.Nodes {
				allNodeName = append(allNodeName, value.NodeName)
			}
			count := oldValue - newValue
			removeNodesName := allNodeName[:count]
			if len(removeNodesName) > 0 {
			}
			req := requests.NewCommonRequest()
			if csService.client.Config.Insecure {
				req.SetHTTPSInsecure(csService.client.Config.Insecure)
			}
			req.QueryParams = map[string]string{
				"RegionId":         csService.client.RegionId,
				"AccessKeySecret":  csService.client.SecretKey,
				"Product":          "CS",
				"Department":       csService.client.Department,
				"ResourceGroup":    csService.client.ResourceGroup,
				"Action":           "RemoveClusterNodes",
				"AccountInfo":      "123456",
				"Version":          "2015-12-15",
				"SignatureVersion": "1.0",
				"ProductName":      "cs",

				"X-acs-body": fmt.Sprintf("{\"%s\":%t,\"%s\":%t,\"%s\":%q,\"%s\":\"%s\"}",

					"release_node", true,
					"drain_node", true,
					"nodes", removeNodesName,
					"ClusterId", d.Id(),
				),
			}
			req.Method = "POST"        // Set request method
			req.Product = "CS"         // Specify product
			req.Version = "2015-12-15" // Specify product version
			req.ServiceCode = "cs"
			if strings.ToLower(csService.client.Config.Protocol) == "https" {
				req.Scheme = "https"
			} else {
				req.Scheme = "http"
			} // Set request scheme. Default: http
			req.ApiName = "RemoveClusterNodes"
			req.Headers = map[string]string{"RegionId": csService.client.RegionId}
			req.Headers = map[string]string{"x-acs-asapi-gateway-version": "3.0"}
			if err := invoker.Run(func() error {
				var err error
				raw, err = csService.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
					return ecsClient.ProcessCommonRequest(req)
				})

				return err
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, nodepoolid, "DeleteKubernetesClusterNodes", DenverdinoAliyungo)
			}
			resp, _ := raw.(*responses.CommonResponse)
			if resp.IsSuccess() == false {
				//return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm", "API Action", cluster.GetHttpContentString())
				return err
			}
			stateConf := BuildStateConf([]string{"removing"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(nodepoolid, d.Id(), []string{"deleting", "failed"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		if newValue > oldValue {
			request := requests.NewCommonRequest()
			if client.Config.Insecure {
				request.SetHTTPSInsecure(client.Config.Insecure)
			}
			request.QueryParams = map[string]string{
				"RegionId":         client.RegionId,
				"AccessKeySecret":  client.SecretKey,
				"Product":          "CS",
				"Department":       client.Department,
				"ResourceGroup":    client.ResourceGroup,
				"Action":           "ScaleClusterNodePool",
				"AccountInfo":      "123456",
				"Version":          "2015-12-15",
				"SignatureVersion": "1.0",
				"ProductName":      "cs",
				"NodepoolId":       nodepoolid,
				"ClusterId":        d.Id(),
				"X-acs-body": fmt.Sprintf("{\"%s\":%d}",

					"count", int64(newValue)-int64(oldValue),
				),
			}
			request.Method = "POST"        // Set request method
			request.Product = "CS"         // Specify product
			request.Version = "2015-12-15" // Specify product version
			request.ServiceCode = "cs"
			if strings.ToLower(client.Config.Protocol) == "https" {
				request.Scheme = "https"
			} else {
				request.Scheme = "http"
			} // Set request scheme. Default: http
			request.ApiName = "ScaleClusterNodePool"
			request.Headers = map[string]string{"RegionId": client.RegionId}
			request.Headers = map[string]string{"x-acs-asapi-gateway-version": "3.0"}
			//var err error
			err = nil
			if err = invoker.Run(func() error {
				raw, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
					return ecsClient.ProcessCommonRequest(request)
				})
				return err
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_cs_kubernetes", "CreateKubernetesCluster", raw)
			}

			if debugOn() {
				resizeRequestMap := make(map[string]interface{})
				resizeRequestMap["ClusterId"] = d.Id()
				resizeRequestMap["Args"] = request.GetQueryParams()
				addDebug("ResizeKubernetesCluster", raw, resizeRequestMap)
			}

			stateConf := BuildStateConf([]string{"scaling"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))

			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			//d.SetPartial("num_of_nodes")
		}
	}

	d.Partial(false)
	return resourceAlibabacloudStackCSKubernetesRead(d, meta)

}

func resourceAlibabacloudStackCSKubernetesRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)

	client := meta.(*connectivity.AlibabacloudStackClient)
	csService := CsService{client}
	object, err := csService.DescribeCsKubernetes(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	nodepoolid := d.Get("nodepool_id").(string)
	clusternode, err := csService.DescribeClusterNodes(d.Id(), nodepoolid)
	if err != nil {
		return WrapError(err)
	}
	d.Set("name", object.Name)
	//d.Set("id", object.ClusterId)
	//d.Set("state", object.State)
	d.Set("vpc_id", object.VpcId)
	//d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("pod_cidr", object.ContainerCIDR)
	d.Set("version", object.CurrentVersion)
	d.Set("delete_protection", object.DeletionProtection)
	d.Set("version", object.InitVersion)
	var sc, smaster, sworker []map[string]interface{}
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
	MasterApi := strings.Replace(object.MasterUrl, "\\", "", -1)
	urlMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(MasterApi), &urlMap)
	intranetApi, ok := urlMap["intranet_api_server_endpoint"]
	if !ok {
		intranetApi = ""
	}
	Connections := map[string]interface{}{
		"api_server_intranet": intranetApi,
	}
	sc = append(sc, Connections)
	d.Set("connections", sc)
	d.Set("master_nodes", smaster)
	d.Set("worker_nodes", sworker)
	if err := d.Set("tags", flattenTagsConfig(object.Tags)); err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAlibabacloudStackCSKubernetesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	csService := CsService{client}
	invoker := NewInvoker()
	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":         client.RegionId,
		"AccessKeySecret":  client.SecretKey,
		"Product":          "CS",
		"Department":       client.Department,
		"ResourceGroup":    client.ResourceGroup,
		"Action":           "DeleteCluster",
		"Version":          "2015-12-15",
		"SignatureVersion": "1.0",
		"ProductName":      "cs",
		"ClusterId":        d.Id(),
		"X-acs-body":       fmt.Sprintf("{\"%s\":\"%t\",\"%s\":\"%s\"}", "keep_slb", false, "ClusterId", d.Id()),
	}
	request.Method = "POST"        // Set request method
	request.Product = "Cs"         // Specify product
	request.Version = "2015-12-15" // Specify product version
	request.ServiceCode = "cs"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	} // Set request scheme. Default: http
	request.ApiName = "DeleteCluster"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.Headers = map[string]string{"x-acs-asapi-gateway-version": "3.0"}
	var response interface{}
	err := resource.Retry(30*time.Minute, func() *resource.RetryError {
		if err := invoker.Run(func() error {
			raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
				return csClient.ProcessCommonRequest(request)
			})
			response = raw
			return err
		}); err != nil {
			return resource.RetryableError(err)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["ClusterId"] = d.Id()
			addDebug("DeleteCluster", response, d.Id(), requestMap)
		}
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteCluster", AlibabacloudStackLogGoSdkERROR)
	}

	stateConf := BuildStateConf([]string{"running", "deleting", "initial"}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Minute, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"delete_failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
func updateKubernetesClusterTag(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	csService := CsService{client}
	d.Partial(true)
	var raw interface{}
	invoker := NewInvoker()
	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	var tags string
	tagss := make([]interface{}, 0)
	if v, ok := d.GetOk("tags"); ok && len(v.(map[string]interface{})) > 0 {
		for key, value := range v.(map[string]interface{}) {
			tagss = append(tagss, cs.Tag{
				Key:   key,
				Value: value.(string),
			})
		}
	}
	tagsBytes, _ := json.Marshal(tagss)
	tags = string(tagsBytes)
	log.Printf("checking tags %v", tags)
	request.QueryParams = map[string]string{
		"RegionId":         client.RegionId,
		"AccessKeySecret":  client.SecretKey,
		"Product":          "CS",
		"Department":       client.Department,
		"ResourceGroup":    client.ResourceGroup,
		"Action":           "ModifyClusterTags",
		"Version":          "2015-12-15",
		"SignatureVersion": "1.0",
		"ProductName":      "cs",
		"ClusterId":        d.Id(),
		"X-acs-body": fmt.Sprintf("{\"%s\":%s}",
			"tags", tags,
		),
	}
	request.Method = "POST"        // Set request method
	request.Product = "CS"         // Specify product
	request.Version = "2015-12-15" // Specify product version
	request.ServiceCode = "cs"
	request.PathPattern = fmt.Sprintf("/clusters/%s/tags", d.Id())
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	} // Set request scheme. Default: http
	request.ApiName = "ModifyClusterTags"
	request.Headers = map[string]string{
		"RegionId":                    client.RegionId,
		"Authorization":               "AuthorizationString",
		"Content-Type":                "application/json",
		"x-acs-asapi-gateway-version": "3.0",
	}
	var err error
	if err = invoker.Run(func() error {
		raw, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			//ecsClient.Domain = "cs.inter.env17e.shuguang.com"
			return ecsClient.ProcessCommonRequest(request)
		})
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_cs_kubernetes", "ModifyClusterTags", raw)
	}
	if debugOn() {
		resizeRequestMap := make(map[string]interface{})
		resizeRequestMap["ClusterId"] = d.Id()
		resizeRequestMap["Args"] = request.GetQueryParams()
		addDebug("ModifyClusterTags", raw, resizeRequestMap)
	}
	stateConf := BuildStateConf([]string{"scaling"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
