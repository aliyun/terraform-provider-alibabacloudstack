package alibabacloudstack

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAqsWebLock() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instanceid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "off",
			},
			"lock_configs": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"dir": {
							Type:     schema.TypeString,
							Required: true,
						},
						"local_backup_dir": {
							Type:     schema.TypeString,
							Required: true,
						},
						"inclusive_file_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"exclusive_dir": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"exclusive_file_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"exclusive_file": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"defence_mode": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"block", "audit"}, false),
						},
						"mode": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"blacklist", "whitelist"}, false),
						},
					},
				},
				Set: func(v interface{}) int {
					config := v.(map[string]interface{})
					return GetConfigsHash(config)
				},
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
	}
	setResourceFunc(resource, resourceAlibabacloudStackAqsWebLockCreate, resourceAlibabacloudStackAqsWebLockRead, resourceAlibabacloudStackAqsWebLockUpdate, resourceAlibabacloudStackAqsWebLockDelete)
	return resource
}

func resourceAlibabacloudStackAqsWebLockCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	aqsService := AqsService{client}
	err := aqsService.RefreshAssets("ecs")
	if err != nil {
		return errmsgs.WrapError(err)
	}
	instances, err := aqsService.DescribeCloudCenterInstances()
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
			"alibabacloudstack_aqs_weblock", "DescribeCloudCenterInstances", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	var uuid string
	for _, v := range instances {
		instance := v.(map[string]interface{})
		if instance["InstanceId"].(string) == d.Get("instanceid").(string) {
			uuid = instance["Uuid"].(string)
		}
	}
	if uuid == "" {
		return errmsgs.WrapError(fmt.Errorf("The instance %s is not exist!", d.Get("instanceid").(string)))
	}
	request := map[string]interface{}{
		"From": "sas",
		"Uuid": uuid,
	}
	lock_configs := d.Get("lock_configs").(*schema.Set).List()
	config := lock_configs[0].(map[string]interface{})
	request["Dir"] = config["dir"]
	request["LocalBackupDir"] = config["local_backup_dir"]
	request["DefenceMode"] = config["defence_mode"]
	mode := config["mode"].(string)
	if mode == "blacklist" {
		request["Mode"] = mode
		request["ExclusiveFileType"] = config["exclusive_file_type"]
		request["ExclusiveFile"] = config["exclusive_file"]
		request["ExclusiveDir"] = config["exclusive_dir"]
	} else {
		request["Mode"] = mode
		request["InclusiveFileType"] = config["inclusive_file_type"]
	}

	_, err = client.DoTeaRequest("POST", "aegis", "2016-11-11", "ModifyWebLockStart", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
			"alibabacloudstack_aqs_weblock", "ModifyWebLockStart", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(lock_configs) > 1 {
		for _, v := range lock_configs[1:] {
			config := v.(map[string]interface{})
			configReq := map[string]interface{}{
				"From":           "sas",
				"Uuid":           uuid,
				"Dir":            config["dir"],
				"LocalBackupDir": config["local_backup_dir"],
				"DefenceMode":    config["defence_mode"],
			}
			mode := config["mode"].(string)
			if mode == "blacklist" {
				configReq["Mode"] = mode
				configReq["ExclusiveFileType"] = config["exclusive_file_type"]
				configReq["ExclusiveFile"] = config["exclusive_file"]
				configReq["ExclusiveDir"] = config["exclusive_dir"]
			} else {
				configReq["Mode"] = mode
				configReq["InclusiveFileType"] = config["inclusive_file_type"]
			}
			_, err = client.DoTeaRequest("POST", "aegis", "2016-11-11", "ModifyWebLockCreateConfig", "", nil, configReq, nil)
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
					"alibabacloudstack_aqs_weblock", "ModifyWebLockCreateConfig", errmsgs.AlibabacloudStackSdkGoERROR)
			}
		}
	}
	d.SetId(uuid)
	return nil
}

func resourceAlibabacloudStackAqsWebLockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	aqsService := AqsService{client}

	object, err := aqsService.DescribeWebLockInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	instances, err := aqsService.DescribeCloudCenterInstances()
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
			"alibabacloudstack_aqs_weblock", "DescribeCloudCenterInstances", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	var instanceid string
	for _, v := range instances {
		instance := v.(map[string]interface{})
		if instance["Uuid"].(string) == d.Id() {
			instanceid = instance["InstanceId"].(string)
		}
	}
	if instanceid != "" {
		d.Set("instanceid", instanceid)
	}
	d.Set("status", object["Status"])
	d.Set("client_status", object["ClientStatus"])
	d.Set("defence_type", object["DefenceType"])
	d.Set("os", object["Os"])
	d.Set("os_name", object["OsName"])
	d.Set("dir_count", object["DirCount"])
	d.Set("service_detail", object["ServiceDetail"])
	d.Set("intranet_ip", object["IntranetIp"])
	d.Set("instance_name", object["InstanceName"])
	d.Set("audit_count", object["AuditCount"])
	d.Set("service_code", object["ServiceCode"])
	d.Set("internet_ip", object["InternetIp"])
	d.Set("service_status", object["ServiceStatus"])
	d.Set("block_count", object["BlockCount"])
	configs, err := aqsService.ListWebLockConfigs(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
			"alibabacloudstack_aqs_weblock", "ListWebLockConfigs", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	lockConfigs := make([]map[string]interface{}, 0)
	for _, v := range configs {
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
		lockConfigs = append(lockConfigs, lockConfig)

	}
	d.Set("lock_configs", lockConfigs)
	return nil
}

func resourceAlibabacloudStackAqsWebLockUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	aqsService := AqsService{client}
	if d.HasChange("status") {
		requestInfo := map[string]interface{}{
			"From":   "sas",
			"Uuid":   d.Id(),
			"Status": d.Get("status").(string),
		}

		_, err := client.DoTeaRequest("POST", "aegis", "2016-11-11", "ModifyWebLockStatus", "", nil, requestInfo, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_aqs_weblock", "ModifyWebLockStatus", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}
	if !d.IsNewResource() && d.HasChange("lock_configs") {
		o, n := d.GetChange("lock_configs")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		removed := os.Difference(ns)
		added := ns.Difference(os)
		have_interim_lockconfig := false
		if removed.Len() == d.Get("lock_configs").(*schema.Set).Len() {
			// "lock_configs" cannot be deleted to empty, add a interim lock_configs
			err := AddInterimLockConfig(d.Id(), meta)
			if err != nil {
				return errmsgs.WrapError(err)
			}
			have_interim_lockconfig = true
		}
		for _, v := range removed.List() {
			config := v.(map[string]interface{})
			id, ok := config["id"]
			if ok {
				err := aqsService.DeleteWebLockConfig(d.Id(), id)
				if err != nil {
					return errmsgs.WrapError(err)
				}
			}
		}
		for _, v := range added.List() {
			config := v.(map[string]interface{})
			configReq := map[string]interface{}{
				"From":           "sas",
				"Uuid":           d.Id(),
				"Dir":            config["dir"],
				"LocalBackupDir": config["local_backup_dir"],
				"DefenceMode":    config["defence_mode"],
			}
			mode := config["mode"].(string)
			if mode == "blacklist" {
				configReq["Mode"] = mode
				configReq["ExclusiveFileType"] = config["exclusive_file_type"]
				configReq["ExclusiveFile"] = config["exclusive_file"]
				configReq["ExclusiveDir"] = config["exclusive_dir"]
			} else {
				configReq["Mode"] = mode
				configReq["InclusiveFileType"] = config["inclusive_file_type"]
			}

			_, err := client.DoTeaRequest("POST", "aegis", "2016-11-11", "ModifyWebLockCreateConfig", "", nil, configReq, nil)
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
					"alibabacloudstack_aqs_weblock", "ModifyWebLockCreateConfig", errmsgs.AlibabacloudStackSdkGoERROR)
			}
		}
		if have_interim_lockconfig {
			err := DeleteInterimLockConfig(d.Id(), meta)
			if err != nil {
				return errmsgs.WrapError(err)
			}
		}
	}

	return nil
}

func resourceAlibabacloudStackAqsWebLockDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	aqsService := AqsService{client}
	object, err := aqsService.DescribeWebLockInstance(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if object["Status"].(string) != "off" {
		requestInfo := map[string]interface{}{
			"From":   "sas",
			"Uuid":   d.Id(),
			"Status": "off",
		}

		_, err := client.DoTeaRequest("POST", "aegis", "2016-11-11", "ModifyWebLockStatus", "", nil, requestInfo, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_aqs_weblock", "ModifyWebLockStatus", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	request := map[string]interface{}{
		"From": "sas",
		"Uuid": d.Id(),
	}

	_, err = client.DoTeaRequest("POST", "aegis", "2016-11-11", "ModifyWebLockUnbind", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
			"alibabacloudstack_aqs_weblock", "ModifyWebLockUnbind", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}

func GetConfigsHash(config map[string]interface{}) int {
	var inclusive_file_type, exclusive_dir, exclusive_file_type, exclusive_file string
	if v, ok := config["inclusive_file_type"]; ok {
		inclusive_file_type = v.(string)
	}
	if v, ok := config["exclusive_dir"]; ok {
		exclusive_dir = v.(string)
	}
	if v, ok := config["exclusive_file_type"]; ok {
		exclusive_file_type = v.(string)
	}
	if v, ok := config["exclusive_file"]; ok {
		exclusive_file = v.(string)
	}
	dir := config["dir"].(string)
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}
	value := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s", dir, config["local_backup_dir"].(string),
		inclusive_file_type, exclusive_dir, exclusive_file_type,
		exclusive_file, config["defence_mode"].(string), config["mode"].(string),
	)
	return hashcode.String(value)
}

func AddInterimLockConfig(uuid string, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := map[string]interface{}{
		"From":              "sas",
		"Uuid":              uuid,
		"Dir":               "/terraforminterimdir/dome/",
		"InclusiveFileType": "php",
		"LocalBackupDir":    "/usr/local/aegis/bak",
		"DefenceMode":       "block",
		"Mode":              "whitelist",
	}
	_, err := client.DoTeaRequest("POST", "aegis", "2016-11-11", "ModifyWebLockCreateConfig", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
			"alibabacloudstack_aqs_weblock", "ModifyWebLockCreateConfig", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}

func DeleteInterimLockConfig(uuid string, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	aqsService := AqsService{client}
	configs, err := aqsService.ListWebLockConfigs(uuid)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
			"alibabacloudstack_aqs_weblock", "ListWebLockConfigs", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	for _, v := range configs {
		config := v.(map[string]interface{})
		if config["Dir"].(string) == "/terraforminterimdir/dome/" {
			err := aqsService.DeleteWebLockConfig(uuid, config["Id"])
			if err != nil {
				return errmsgs.WrapError(err)
			}
			break
		}
	}
	return nil
}
