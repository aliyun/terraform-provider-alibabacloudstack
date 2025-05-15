package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEdasK8sApplication() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"application_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"replicas": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"image_url": {
				Type:     schema.TypeString,
				Optional: true,
				//ConflictsWith: []string{"package_url"},
			},
			"package_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"FatJar", "War", "Image"}, false),
				Default:      "Image",
			},
			"application_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"application_descriotion": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'application_descriotion' is deprecated and will be removed in a future release. Please use new field 'application_description' instead.",
			},
			"limit_mem": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"requests_mem": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"command": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"command_args": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"internet_slb_id": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"internet_external_traffic_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Local", "Cluster"}, false),
				Default:      "Local",
			},
			"internet_scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"rr", "wrr"}, false),
				Default:      "rr",
			},
			"internet_slb_protocol": {
				Optional:     true,
				Type:         schema.TypeString,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS"}, false),
			},
			"internet_slb_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"internet_target_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"internet_service_port_infos": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"internet_slb_protocol", "internet_target_port", "internet_slb_port"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS"}, false),
						},
						"target_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"load_balancer_protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS"}, false),
						},
					},
				},
			},
			"intranet_slb_id": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"intranet_external_traffic_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Local", "Cluster"}, false),
				Default:      "Local",
			},
			"intranet_scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"rr", "wrr"}, false),
				Default:      "rr",
			},
			"intranet_slb_protocol": {
				Optional:     true,
				Type:         schema.TypeString,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS"}, false),
			},
			"intranet_slb_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"intranet_target_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"intranet_service_port_infos": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"intranet_target_port", "intranet_slb_port", "intranet_slb_protocol"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS"}, false),
						},
						"target_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"load_balancer_protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS"}, false),
						},
					},
				},
			},
			"envs": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"pre_stop": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					e := EdasService{}
					return e.PreStopEqual(old, new)
				},
			},
			"post_start": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					e := EdasService{}
					return e.PostStartEqual(old, new)
				},
			},
			"liveness": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					e := EdasService{}
					return e.LivenessEqual(old, new)
				},
			},
			"readiness": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					e := EdasService{}
					return e.ReadinessEqual(old, new)
				},
			},
			"nas_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mount_descs": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"logical_region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"package_url": {
				Type:     schema.TypeString,
				Optional: true,
				//ConflictsWith: []string{"image_url"},
			},
			"package_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				//Default:  strconv.FormatInt(time.Now().Unix(), 10),
			},
			"jdk": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"web_container": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"edas_container_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"requests_m_cpu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit_m_cpu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cr_ee_repo_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"update_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"BatchUpdate", "GrayBatchUpdate"}, false),
			},
			"update_batch": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"update_release_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"auto", "manual"}, false),
			},
			"update_batch_wait_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"update_gray": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"config_mount_descs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ConfigMap", "Secret"}, false),
						},
						"mount_path": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"local_volume": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							// ValidateFunc: validation.StringInSlice([]string{"file", "filepath"}, false),
						},
						"node_path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mount_path": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"pvc_mount_descs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pvc_name": {
							Type:     schema.TypeString,
							Required: true,
							// ValidateFunc: validation.StringInSlice([]string{"file", "filepath"}, false),
						},
						"mount_paths": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mount_path": {
										Type:     schema.TypeString,
										Required: true,
									},
									"read_only": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEdasK8sApplicationCreate, resourceAlibabacloudStackEdasK8sApplicationRead, resourceAlibabacloudStackEdasK8sApplicationUpdate, resourceAlibabacloudStackEdasK8sApplicationDelete)
	return resource
}

func resourceAlibabacloudStackEdasK8sApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}
	// request := edas.CreateInsertK8sApplicationRequest()
	request := client.NewCommonRequest("POST", "Edas", "2017-08-01", "InsertK8sApplication", "/pop/v5/k8s/acs/create_k8s_app")
	packageType := d.Get("package_type").(string)

	request.QueryParams["AppName"] = d.Get("application_name").(string)
	request.QueryParams["RegionId"] = client.RegionId
	request.QueryParams["PackageType"] = packageType
	request.QueryParams["ClusterId"] = d.Get("cluster_id").(string)
	if strings.ToLower(packageType) == "image" {
		if v, ok := d.GetOk("image_url"); !ok {
			return errmsgs.WrapError(errmsgs.Error("image_url is needed for creating image k8s application"))
		} else {
			request.QueryParams["ImageUrl"] = v.(string)
			if strings.HasPrefix(v.(string), "cr-ee.registry") {
				if crid, ok := d.GetOk("cr_ee_repo_id"); ok && crid.(string) != "" {
					request.QueryParams["crInstanceId"] = crid.(string)
				} else {
					return errmsgs.WrapError(errmsgs.Error("`cr_ee_repo_id` is needed for the image repo is enterprise-edition"))
				}
			}
		}
	} else {
		if v, ok := d.GetOk("package_url"); !ok {
			return errmsgs.WrapError(errmsgs.Error("package_url is needed for creating fatjar k8s application"))
		} else {
			request.QueryParams["PackageUrl"] = v.(string)
			request.QueryParams["ImageUrl"] = v.(string)
		}
		if v, ok := d.GetOk("package_version"); ok {
			request.QueryParams["PackageVersion"] = v.(string)
		}
		if v, ok := d.GetOk("jdk"); !ok {
			return errmsgs.WrapError(errmsgs.Error("jdk is needed for creating non-image k8s application"))
		} else {
			request.QueryParams["JDK"] = v.(string)
		}
		if strings.ToLower(packageType) == "war" {
			var webContainer string
			var edasContainer string
			if v, ok := d.GetOk("web_container"); ok {
				webContainer = v.(string)
			}
			if v, ok := d.GetOk("edas_container_version"); ok {
				edasContainer = v.(string)
			}
			if len(webContainer) == 0 && len(edasContainer) == 0 {
				return errmsgs.WrapError(errmsgs.Error("web_container or edas_container_version is needed for creating war k8s application"))
			}
			request.QueryParams["WebContainer"] = webContainer
			request.QueryParams["EdasContainerVersion"] = edasContainer
		}
	}

	request.QueryParams["Replicas"] = fmt.Sprintf("%d", d.Get("replicas").(int))

	if v, ok := connectivity.GetResourceDataOk(d, "application_description", "application_descriotion"); ok {
		request.QueryParams["ApplicationDescription"] = v.(string)
	}

	if v, ok := d.GetOk("limit_mem"); ok {
		request.QueryParams["LimitMem"] = fmt.Sprintf("%d", v.(int))
	}

	if v, ok := d.GetOk("requests_mem"); ok {
		request.QueryParams["RequestsMem"] = fmt.Sprintf("%d", v.(int))
	}

	if v, ok := d.GetOk("command"); ok {
		request.QueryParams["Command"] = v.(string)
	}

	if v, ok := d.GetOk("command_args"); ok {
		commands, err := edasService.GetK8sCommandArgs(v.([]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.QueryParams["CommandArgs"] = commands
	}

	if v, ok := d.GetOk("envs"); ok {
		envs, err := edasService.GetK8sEnvs(v.(map[string]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.QueryParams["Envs"] = envs
	}

	if v, ok := d.GetOk("pre_stop"); ok {
		request.QueryParams["PreStop"] = v.(string)
	}

	if v, ok := d.GetOk("post_start"); ok {
		request.QueryParams["PostStart"] = v.(string)
	}

	if v, ok := d.GetOk("liveness"); ok {
		request.QueryParams["Liveness"] = v.(string)
	}

	if v, ok := d.GetOk("readiness"); ok {
		request.QueryParams["Readiness"] = v.(string)
	}

	if v, ok := d.GetOk("nas_id"); ok {
		request.QueryParams["NasId"] = v.(string)
	}

	if v, ok := d.GetOk("mount_descs"); ok {
		request.QueryParams["MountDescs"] = v.(string)
	}

	if v, ok := d.GetOk("config_mount_descs"); ok {
		configmaps, err := edasService.GetK8sConfigMaps(v.([]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.QueryParams["ConfigMountDescs"] = configmaps
	}

	if v, ok := d.GetOk("local_volume"); ok {
		local_volumes, err := edasService.GetK8sLocalVolumes(v.([]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.QueryParams["LocalVolume"] = local_volumes
	}

	if v, ok := d.GetOk("pvc_mount_descs"); ok {
		pvc_mount_descs, err := edasService.GetK8sPvcMountDescs(v.([]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.QueryParams["PvcMountDescs"] = pvc_mount_descs
	}

	if v, ok := d.GetOk("namespace"); ok {
		request.QueryParams["Namespace"] = v.(string)
	}

	if v, ok := d.GetOk("logical_region_id"); ok {
		request.QueryParams["LogicalRegionId"] = v.(string)
	}

	if v, ok := d.GetOk("requests_m_cpu"); ok {
		request.QueryParams["RequestsmCpu"] = fmt.Sprintf("%d", v.(int))
	}

	if v, ok := d.GetOk("limit_m_cpu"); ok {
		request.QueryParams["LimitmCpu"] = fmt.Sprintf("%d", v.(int))
	}
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug("InsertK8sApplication", bresponse, request.QueryParams, request)
	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_k8s_application", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	var response map[string]interface{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if fmt.Sprint(response["Code"]) != "200" {
		return errmsgs.WrapError(fmt.Errorf("Create k8s application failed for %s", response["Message"].(string)))
	}
	appId := response["ApplicationInfo"].(map[string]interface{})["AppId"].(string)
	changeOrderId := response["ApplicationInfo"].(map[string]interface{})["ChangeOrderId"].(string)
	d.SetId(appId)

	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	return nil
}

func resourceAlibabacloudStackEdasK8sApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	response, err := edasService.DescribeEdasK8sApplication(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_edas_k8s_application ecsService.DescribeEdasK8sApplication Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("application_name", response.App.ApplicationName)
	d.Set("cluster_id", response.App.ClusterId)
	d.Set("replicas", response.App.Instances)
	d.Set("package_type", response.App.ApplicationType)
	if d.Get("package_type").(string) == "docker" {
		d.Set("image_url", response.ImageInfo.ImageUrl)
	}
	envs := make(map[string]string)
	for _, e := range response.App.EnvList.Env {
		envs[e.Name] = e.Value
	}
	d.Set("envs", envs)
	d.Set("command", response.App.Cmd)
	d.Set("command_args", response.App.CmdArgs.CmdArg)
	d.Set("requests_m_cpu", response.App.RequestCpuM)
	d.Set("limit_m_cpu", response.App.LimitCpuM)
	d.Set("limit_mem", response.App.LimitMem)
	d.Set("requests_mem", response.App.RequestMem)

	allDeploy := response.DeployGroups.DeployGroup
	for _, v := range allDeploy {
		if len(v.PackageVersion) > 0 {
			d.Set("package_version", v.PackageVersion)
		}

		for _, c := range v.Components.ComponentsItem {
			if strings.Contains(c.ComponentKey, "JDK") {
				d.Set("jdk", c.ComponentKey)
			}
		}
	}

	if len(response.App.EdasContainerVersion) > 0 {
		d.Set("edas_container_version", response.App.EdasContainerVersion)
	}
	intranet_slbs := make([]map[string]interface{}, 0)
	internet_slbs := make([]map[string]interface{}, 0)
	intranet_slb_id := ""
	internet_slb_id := ""
	intranet_external_traffic_policy := "Local"
	internet_external_traffic_policy := "Local"
	if response.App.SlbInfo != "" {
		var slbinfos []interface{}
		err = json.Unmarshal([]byte(response.App.SlbInfo), &slbinfos)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		for _, slbinfo := range slbinfos {
			slb := slbinfo.(map[string]interface{})
			if slb["addressType"] == "intranet" {
				intranet_slb_id = slb["slbId"].(string)
				intranet_external_traffic_policy = slb["externalTrafficPolicy"].(string)
				for _, service_port := range slb["portMappings"].([]interface{}) {
					info := service_port.(map[string]interface{})
					port_data := info["servicePort"].(map[string]interface{})
					intranet_slbs = append(intranet_slbs, map[string]interface{}{
						"protocol":               port_data["protocol"].(string),
						"port":                   port_data["port"].(int),
						"target_port":            port_data["targetPort"].(int),
						"load_balancer_protocol": info["loadBalancerProtocol"].(string),
					})

				}
			} else if slb["addressType"] == "internet" {
				internet_slb_id = slb["slbId"].(string)
				internet_external_traffic_policy = slb["externalTrafficPolicy"].(string)
				for _, service_port := range slb["portMappings"].([]interface{}) {
					info := service_port.(map[string]interface{})
					port_data := info["servicePort"].(map[string]interface{})
					internet_slbs = append(internet_slbs, map[string]interface{}{
						"protocol":               port_data["protocol"].(string),
						"port":                   port_data["port"].(int),
						"target_port":            port_data["targetPort"].(int),
						"load_balancer_protocol": info["loadBalancerProtocol"].(string),
					})
				}
			}
		}
	}
	d.Set("intranet_slb_id", intranet_slb_id)
	d.Set("intranet_external_traffic_policy", intranet_external_traffic_policy)
	d.Set("intranet_service_port_infos", intranet_slbs)
	if len(intranet_slbs) == 1 {
		d.Set("intranet_slb_protocol", intranet_slbs[0]["protocol"].(string))
		d.Set("intranet_target_port", intranet_slbs[0]["target_port"].(int))
		d.Set("intranet_slb_port", intranet_slbs[0]["port"].(int))
	}
	d.Set("internet_slb_id", internet_slb_id)
	d.Set("internet_external_traffic_policy", internet_external_traffic_policy)
	d.Set("internet_service_port_infos", internet_slbs)
	if len(internet_slbs) == 1 {
		d.Set("internet_slb_protocol", internet_slbs[0]["protocol"].(string))
		d.Set("internet_target_port", internet_slbs[0]["target_port"].(int))
		d.Set("internet_slb_port", internet_slbs[0]["port"].(int))
	}
	if len(response.Conf.PreStop) > 0 {
		d.Set("pre_stop", response.Conf.PreStop)
	}
	if len(response.Conf.PostStart) > 0 {
		d.Set("post_start", response.Conf.PostStart)
	}
	if len(response.Conf.Liveness) > 0 {
		d.Set("liveness", response.Conf.Liveness)
	}
	if len(response.Conf.Readiness) > 0 {
		d.Set("readiness", response.Conf.Readiness)
	}
	d.Set("namespace", response.App.K8sNamespace)
	if len(response.Conf.K8sVolumeInfo) > 0 {
		k8sVolumeInfo := make(map[string]interface{})
		err = json.Unmarshal([]byte(response.Conf.K8sVolumeInfo), &k8sVolumeInfo)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		configMountDescs, ok := k8sVolumeInfo["configMountDescs"]
		if ok {
			configmaps := make([]ConfigMaps, 0)
			err = json.Unmarshal([]byte(configMountDescs.(string)), &configmaps)
			if err != nil {
				return errmsgs.WrapError(err)
			}
			config_mount_descs := make([]map[string]interface{}, 0)
			for _, v := range configmaps {
				config_mount_descs = append(config_mount_descs, map[string]interface{}{
					"name":       v.Name,
					"type":       v.Type,
					"mount_path": v.MountPath,
				})
			}
			d.Set("config_mount_descs", config_mount_descs)
		}

		pvcMountDescs, ok := k8sVolumeInfo["pvcMountDescs"]
		if ok {
			pvcmounts := make([]PvcMountDescs, 0)
			err = json.Unmarshal([]byte(pvcMountDescs.(string)), &pvcmounts)
			if err != nil {
				return errmsgs.WrapError(err)
			}
			pvc_mount_descs := make([]map[string]interface{}, 0)
			for _, v := range pvcmounts {
				mountPaths := make([]map[string]interface{}, 0)
				for _, vm := range v.MountPaths {
					mountPaths = append(mountPaths, map[string]interface{}{
						"mount_path": vm.MountPath,
						"read_only":  vm.ReadOnly,
					})
				}
				pvc_mount_descs = append(pvc_mount_descs, map[string]interface{}{
					"pvc_name":    v.PvcName,
					"mount_paths": mountPaths,
				})
			}
			d.Set("pvc_mount_descs", pvc_mount_descs)
		}
	}
	if response.Conf.K8sLocalvolumeInfo != "" {
		K8sLocalvolumeInfo := make(map[string]interface{})
		err = json.Unmarshal([]byte(response.Conf.K8sLocalvolumeInfo), &K8sLocalvolumeInfo)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		localVolumeDOs, ok := K8sLocalvolumeInfo["localVolumeDOs"]
		local_volumes := make([]map[string]string, 0)
		if ok {
			localVolumes := localVolumeDOs.([]interface{})
			for _, lv := range localVolumes {
				v := lv.(map[string]interface{})
				local_volume := map[string]string{
					"mount_path": v["mountPath"].(string),
					"node_path":  v["nodePath"].(string),
				}
				lv_type, ok := v["type"]
				if ok {
					local_volume["type"] = lv_type.(string)

				} else {
					local_volume["type"] = ""
				}
				local_volumes = append(local_volumes, local_volume)
			}
		}
		d.Set("local_volume", local_volumes)
	}
	return nil
}

func resourceAlibabacloudStackEdasK8sApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}
	var partialKeys []string
	request := client.NewCommonRequest("POST", "Edas", "2017-08-01", "DeployK8sApplication", "/pop/v5/k8s/acs/k8s_apps")
	request.QueryParams["RegionId"] = client.RegionId
	request.QueryParams["AppId"] = d.Id()

	// 检查该app是否已经绑定了slb
	appobj, err := edasService.DescribeEdasK8sApplication(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	intranet_slb_unset := true
	internet_slb_unset := true
	if appobj.App.SlbInfo != "" {
		var slbinfos []interface{}
		err = json.Unmarshal([]byte(appobj.App.SlbInfo), &slbinfos)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		for _, slbinfo := range slbinfos {
			slb := slbinfo.(map[string]interface{})
			if slb["addressType"] == "intranet" {
				intranet_slb_unset = false
			} else if slb["addressType"] == "internet" {
				internet_slb_unset = false
			}
		}
	}
	d.Partial(true)
	if intranet_slb_unset {
		bind_slb_err := K8sBindSlb("intranet", intranet_slb_unset, d, meta)
		if bind_slb_err != nil {
			return errmsgs.WrapError(bind_slb_err)
		}
	} else {
		if d.Get("intranet_slb_id") == "" {
			err = DeleteK8sSlb("intranet", d, meta)
			if err != nil {
				return errmsgs.WrapError(err)
			}
		}
	}

	if internet_slb_unset {
		bind_slb_err := K8sBindSlb("internet", internet_slb_unset, d, meta)
		if bind_slb_err != nil {
			return errmsgs.WrapError(bind_slb_err)
		}
	} else {
		if d.Get("internet_slb_id") == "" {
			err = DeleteK8sSlb("internet", d, meta)
			if err != nil {
				return errmsgs.WrapError(err)
			}
		}
	}

	packageType, err := edasService.QueryK8sAppPackageType(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if strings.ToLower(packageType) == "image" {
		if d.HasChange("image_url") {
			partialKeys = append(partialKeys, "image_url")
		}
		request.QueryParams["Image"] = d.Get("image_url").(string)
		if len(request.QueryParams["Image"]) == 0 {
			return errmsgs.WrapError(errmsgs.Error("image_url is needed for creating image k8s application"))
		}
	} else {
		if d.HasChange("package_url") {
			partialKeys = append(partialKeys, "package_url")
		}
		request.QueryParams["PackageUrl"] = d.Get("package_url").(string)
		if len(request.QueryParams["PackageUrl"]) == 0 {
			return errmsgs.WrapError(errmsgs.Error("package_url is needed for creating fatjar k8s application"))
		}
		if d.HasChange("package_version") {
			partialKeys = append(partialKeys, "package_version")
		}
		request.QueryParams["PackageVersion"] = d.Get("package_version").(string)

		if d.HasChange("jdk") {
			partialKeys = append(partialKeys, "jdk")
		}
		request.QueryParams["JDK"] = d.Get("jdk").(string)
		if len(request.QueryParams["JDK"]) == 0 {
			return errmsgs.WrapError(errmsgs.Error("jdk is needed for creating non-image k8s application"))
		}
		if strings.ToLower(packageType) == "war" {
			var webContainer string
			var edasContainer string
			if d.HasChange("web_container") {
				partialKeys = append(partialKeys, "web_container")
			}
			webContainer = d.Get("web_container").(string)

			if d.HasChange("edas_container_version") {
				partialKeys = append(partialKeys, "edas_container_version")
			}
			edasContainer = d.Get("edas_container_version").(string)
			if len(webContainer) == 0 && len(edasContainer) == 0 {
				return errmsgs.WrapError(errmsgs.Error("web_container or edas_container_version is needed for updating war k8s application"))
			}
			request.QueryParams["WebContainer"] = webContainer
			request.QueryParams["EdasContainerVersion"] = edasContainer
		}
	}

	if d.HasChange("replicas") {
		partialKeys = append(partialKeys, "replicas")
	}
	replicas := d.Get("replicas").(int)
	request.QueryParams["Replicas"] = fmt.Sprintf("%d", replicas)

	if d.HasChange("limit_mem") {
		partialKeys = append(partialKeys, "limit_mem")
		request.QueryParams["MemoryLimit"] = fmt.Sprintf("%d", d.Get("limit_mem").(int))
	}

	if d.HasChange("requests_mem") {
		partialKeys = append(partialKeys, "requests_mem")
		request.QueryParams["MemoryRequest"] = fmt.Sprintf("%d", d.Get("requests_mem").(int))
	}

	if d.HasChange("command") {
		partialKeys = append(partialKeys, "command")
		request.QueryParams["Command"] = d.Get("command").(string)
	}

	if d.HasChange("command_args") {
		partialKeys = append(partialKeys, "command_args")
		commands, err := edasService.GetK8sCommandArgsForDeploy(d.Get("command_args").([]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.QueryParams["Args"] = commands
	}

	if d.HasChange("envs") {
		partialKeys = append(partialKeys, "envs")
		envs, err := edasService.GetK8sEnvs(d.Get("envs").(map[string]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.QueryParams["Envs"] = envs
	}

	if d.HasChange("pre_stop") {
		if !edasService.PreStopEqual(d.GetChange("pre_stop")) {
			partialKeys = append(partialKeys, "pre_stop")
			request.QueryParams["PreStop"] = d.Get("pre_stop").(string)
		}
	}

	if d.HasChange("post_start") {
		if !edasService.PostStartEqual(d.GetChange("post_start")) {
			partialKeys = append(partialKeys, "post_start")
			request.QueryParams["PostStart"] = d.Get("post_start").(string)
		}
	}

	if d.HasChange("liveness") {
		if !edasService.LivenessEqual(d.GetChange("liveness")) {
			partialKeys = append(partialKeys, "liveness")
			request.QueryParams["Liveness"] = d.Get("liveness").(string)
		}
	}

	if d.HasChange("readiness") {
		if !edasService.ReadinessEqual(d.GetChange("readiness")) {
			partialKeys = append(partialKeys, "readiness")
			request.QueryParams["Readiness"] = d.Get("readiness").(string)
		}
	}

	if d.HasChange("nas_id") {
		partialKeys = append(partialKeys, "nas_id")
		request.QueryParams["NasId"] = d.Get("nas_id").(string)
	}

	if d.HasChange("mount_descs") {
		partialKeys = append(partialKeys, "mount_descs")
		request.QueryParams["MountDescs"] = d.Get("mount_descs").(string)
	}

	if d.HasChange("config_mount_descs") {
		configmaps, err := edasService.GetK8sConfigMaps(d.Get("config_mount_descs").([]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.QueryParams["ConfigMountDescs"] = configmaps
	}

	if d.HasChange("local_volume") {
		local_volumes, err := edasService.GetK8sLocalVolumes(d.Get("local_volume").([]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.QueryParams["LocalVolume"] = local_volumes
	}

	if d.HasChange("pvc_mount_descs") {
		pvc_mount_descs, err := edasService.GetK8sPvcMountDescs(d.Get("pvc_mount_descs").([]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.QueryParams["PvcMountDescs"] = pvc_mount_descs
	}

	if d.HasChange("requests_m_cpu") {
		partialKeys = append(partialKeys, "requests_m_cpu")
		request.QueryParams["McpuRequest"] = fmt.Sprintf("%d", d.Get("requests_m_cpu").(int))
	}

	if d.HasChange("limit_m_cpu") {
		partialKeys = append(partialKeys, "limit_m_cpu")
		request.QueryParams["McpuLimit"] = fmt.Sprintf("%d", d.Get("limit_m_cpu").(int))
	}

	if len(partialKeys) > 0 && !d.IsNewResource() {
		if v, ok := d.GetOk("update_type"); ok && v.(string) != "" && replicas > 1 {
			partialKeys = append(partialKeys, "update_type")
			update_type := d.Get("update_type").(string)
			update_batch := d.Get("update_batch").(int)
			if update_batch < 2 {
				return errmsgs.WrapError(errmsgs.Error("`update_batch` must be greater than 2"))
			}
			update_release_type := d.Get("update_release_type").(string)
			gray_update_strategy := ""
			if v, ok := d.GetOk("update_gray"); ok {
				update_gray := v.(int)
				gray_update_strategy = fmt.Sprintf(",\"grayUpdate\":{\"gray\":%d}", update_gray)
			}
			if update_release_type == "auto" {
				update_batch_wait_time := d.Get("update_batch_wait_time").(int)
				request.QueryParams["UpdateStrategy"] = fmt.Sprintf("{\"type\":\"%s\",\"batchUpdate\":{\"batch\":%d,\"releaseType\":\"%s\",\"batchWaitTime\":%d}%s}", update_type, update_batch, update_release_type, update_batch_wait_time, gray_update_strategy)
			} else {
				request.QueryParams["UpdateStrategy"] = fmt.Sprintf("{\"type\":\"%s\",\"batchUpdate\":{\"batch\":%d,\"releaseType\":\"%s\"}%s}", update_type, update_batch, update_release_type, gray_update_strategy)
			}
		}
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request)

		if err != nil {
			errmsg := ""
			if bresponse != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		response := make(map[string]interface{})
		_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		changeOrderId := response["ChangeOrderId"].(string)
		if fmt.Sprint(response["Code"]) != "200" {
			return errmsgs.WrapError(errmsgs.Error("deploy k8s application failed for " + response["Message"].(string)))
		}

		if changeOrderId != "" {
			stateConf := BuildStateConf([]string{"0", "1", "9"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
			}
		}
		//for _, key := range partialKeys {
		//	d.SetPartial(key)
		//}
	}
	d.Partial(false)
	return nil
}

func resourceAlibabacloudStackEdasK8sApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	// request := edas.CreateDeleteK8sApplicationRequest()
	request := client.NewCommonRequest("DELETE", "Edas", "2017-08-01", "DeleteK8sApplication", "/pop/v5/k8s/acs/k8s_apps")
	request.QueryParams["RegionId"] = client.RegionId
	request.QueryParams["AppId"] = d.Id()
	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-Type"] = "application/json"
	wait := incrementalWait(1*time.Second, 2*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			errmsg := ""
			if bresponse != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)

			return resource.NonRetryableError(err)
		}
		response := make(map[string]interface{})
		_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		changeOrderId := response["ChangeOrderId"].(string)
		if fmt.Sprint(response["Code"]) != "200" {
			return resource.NonRetryableError(errmsgs.Error("Delete k8s application failed for " + response["Message"].(string)))
		}

		if changeOrderId != "" {
			stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id()))
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func K8sBindSlb(net_type string, isnew bool, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}
	action := "BindK8sSlb"
	method := "POST"
	if !isnew {
		action = "UpdateK8sSlb"
		method = "PUT"
	}
	bind_slb := false
	request := client.NewCommonRequest(method, "Edas", "2017-08-01", action, "/pop/v5/k8s/acs/k8s_slb_binding")
	request.QueryParams["RegionId"] = client.RegionId
	request.QueryParams["AppId"] = d.Id()
	request.QueryParams["ClusterId"] = d.Get("cluster_id").(string)
	request.QueryParams["Type"] = net_type
	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-Type"] = "application/json"
	if v, ok := d.GetOk(fmt.Sprintf("%s_slb_id", net_type)); ok {
		request.QueryParams["SlbId"] = v.(string)
		bind_slb = true
	} else {
		request.QueryParams["Scheduler"] = d.Get(fmt.Sprintf("%s_scheduler", net_type)).(string)
		request.QueryParams["ExternalTrafficPolicy"] = d.Get(fmt.Sprintf("%s_external_traffic_policy", net_type)).(string)
		service_port_infos := make([]map[string]interface{}, 0)
		if v, ok := d.GetOk(fmt.Sprintf("%s_service_port_infos", net_type)); ok && len(v.([]interface{})) > 0 {
			for _, info := range v.([]interface{}) {
				service_port_info := info.(map[string]interface{})
				service_port_infos = append(service_port_infos, map[string]interface{}{
					"protocol":             service_port_info["protocol"].(string),
					"port":                 service_port_info["port"].(int),
					"targetPort":           service_port_info["target_port"].(int),
					"loadBalancerProtocol": service_port_info["load_balancer_protocol"].(string),
				})
			}
		} else {
			v1, ok1 := d.GetOk(fmt.Sprintf("%s_target_port", net_type))
			v2, ok2 := d.GetOk(fmt.Sprintf("%s_slb_port", net_type))
			v3, ok3 := d.GetOk(fmt.Sprintf("%s_slb_protocol", net_type))
			if ok1 && ok2 && ok3 {
				service_port_infos = append(service_port_infos, map[string]interface{}{
					"protocol":             v3.(string),
					"port":                 v2.(int),
					"targetPort":           v1.(int),
					"loadBalancerProtocol": v3.(string),
				})
			}
		}
		if len(service_port_infos) > 0 {
			data, err := json.Marshal(service_port_infos)
			if err != nil {
				return fmt.Errorf("slb service_port_infos data to marshal JSON failed: %w", err)
			}
			request.QueryParams["ServicePortInfos"] = string(data)
			bind_slb = true
		}

	}
	if bind_slb {
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(fmt.Sprintf("BindK8sSlb: %s", net_type), bresponse, request)
		if err != nil {
			errmsg := ""
			if bresponse != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_k8s_application", "BindK8sSlb", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		response := make(map[string]interface{})
		_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if fmt.Sprint(response["Code"]) != "200" {
			return errmsgs.WrapError(fmt.Errorf("BindK8sSlb Failed , response: %#v", response))
		}
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(response["ChangeOrderId"].(string), []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapError(fmt.Errorf("BindK8sSlb Failed , response: %#v", response))
		}
		time.Sleep(time.Duration(10) * time.Second)
	}
	return nil
}

func DeleteK8sSlb(net_type string, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "BindK8sSlb"
	request := client.NewCommonRequest("DELETE", "Edas", "2017-08-01", action, "/pop/v5/k8s/acs/k8s_slb_binding")
	request.QueryParams["RegionId"] = client.RegionId
	request.QueryParams["AppId"] = d.Id()
	request.QueryParams["ClusterId"] = d.Get("cluster_id").(string)
	request.QueryParams["Type"] = net_type
	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-Type"] = "application/json"
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug("DeleteK8sSlb", bresponse, request)
	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_k8s_application", "DeleteK8sSlb", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	response := make(map[string]interface{})
	_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if fmt.Sprint(response["Code"]) != "200" {
		return errmsgs.WrapError(fmt.Errorf("BindK8sSlb Failed , response: %#v", response))
	}
	return nil
}
