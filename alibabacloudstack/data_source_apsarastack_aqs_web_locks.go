package alibabacloudstack

import (
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackAqsWebLocks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAqsWebLocksRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of Web Lock Config IDs.",
			},
			"instanceid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"weblocks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The UUID of the associated server.",
						},
						"lock_configs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The ID of the web lock configuration.",
									},
									"dir": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The directory path being protected.",
									},
									"inclusive_file_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "List of file types included in protection, separated by semicolons.",
									},
									"exclusive_file": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "List of excluded files, separated by semicolons.",
									},
									"exclusive_dir": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "List of excluded directories, separated by semicolons.",
									},
									"defence_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The defense mode: 'block' or others.",
									},
									"mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The operation mode: 'whitelist' or others.",
									},
									"local_backup_dir": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The local backup directory for locked files.",
									},
									"exclusive_file_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "List of file types excluded from protection, separated by semicolons.",
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"defence_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dir_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"service_detail": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"intranet_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"audit_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"service_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"block_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAqsWebLocksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	idsMap := make(map[string]string)
	var filterUuid string
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	aqsService := AqsService{client}
	if v, ok := d.GetOk("instanceid"); ok && v.(string) != "" {
		instances, err := aqsService.DescribeCloudCenterInstances()
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_aqs_weblock", "DescribeCloudCenterInstances", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		for _, v := range instances {
			instance := v.(map[string]interface{})
			if instance["InstanceId"].(string) == d.Get("instanceid").(string) {
				filterUuid = instance["Uuid"].(string)
			}
		}
		if filterUuid == "" {
			d.SetId("")
			return nil
		}
	}
	request := map[string]interface{}{
		"From":        "sas",
		"CurrentPage": 1,
		"PageSize":    1000,
	}

	response, err := client.DoTeaRequest("GET", "aegis", "2016-11-11", "DescribeWebLockBindList", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	bindlist, ok := response["BindList"].([]interface{})
	if !ok && len(bindlist) == 0 {
		d.SetId("")
		return nil
	}
	ids := make([]string, 0)
	weblocks := make([]map[string]interface{}, 0)
	for _, v := range bindlist {
		bind := v.(map[string]interface{})
		uuid := bind["Uuid"].(string)
		if len(idsMap) > 0 {
			if _, ok := idsMap[uuid]; !ok {
				continue
			}
		}
		if filterUuid != "" && uuid != filterUuid {
			continue
		}
		configs := make([]map[string]interface{}, 0)
		req := map[string]interface{}{
			"From": "sas",
			"Uuid": uuid,
		}

		result, err := client.DoTeaRequest("GET", "aegis", "2016-11-11", "DescribeWebLockConfigList", "", nil, req, nil)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		configList, ok := result["ConfigList"].([]interface{})
		if ok {
			for _, v := range configList {
				config := v.(map[string]interface{})
				lockConfig := map[string]interface{}{
					"id":               config["Id"],
					"dir":              config["Dir"],
					"defence_mode":     config["DefenceMode"],
					"mode":             config["Mode"],
					"local_backup_dir": config["LocalBackupDir"],
				}
				if v, ok := config["InclusiveFileType"]; ok && v.(string) != "" {
					lockConfig["inclusive_file_type"] = v
				}
				if v, ok := config["ExclusiveFile"]; ok && v.(string) != "" {
					lockConfig["exclusive_file"] = v
				}
				if v, ok := config["ExclusiveDir"]; ok && v.(string) != "" {
					lockConfig["exclusive_dir"] = v
				}
				if v, ok := config["ExclusiveFileType"]; ok && v.(string) != "" {
					lockConfig["exclusive_file_type"] = v
				}
				configs = append(configs, lockConfig)

			}
		}
		weblock := map[string]interface{}{
			"id":             uuid,
			"uuid":           uuid,
			"lock_configs":   configs,
			"status":         bind["Status"],
			"client_status":  bind["ClientStatus"],
			"defence_type":   bind["DefenceType"],
			"os":             bind["Os"],
			"os_name":        bind["OsName"],
			"dir_count":      bind["DirCount"],
			"service_detail": bind["ServiceDetail"],
			"intranet_ip":    bind["IntranetIp"],
			"instance_name":  bind["InstanceName"],
			"audit_count":    bind["AuditCount"],
			"service_code":   bind["ServiceCode"],
			"internet_ip":    bind["InternetIp"],
			"service_status": bind["ServiceStatus"],
			"block_count":    bind["BlockCount"],
		}
		weblocks = append(weblocks, weblock)
		ids = append(ids, uuid)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("weblocks", weblocks); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
