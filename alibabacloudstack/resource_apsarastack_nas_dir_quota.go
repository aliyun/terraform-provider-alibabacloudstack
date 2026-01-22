package alibabacloudstack

import (
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackNasDirQuota() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quotas": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"quota_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Enforcement", "Accounting"}, false),
						},
						"user_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Uid", "Gid", "AllUsers"}, false),
						},
						"user_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"size_limit": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"file_count_limit": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackNasDirQuotaCreate, resourceAlibabacloudStackNasDirQuotaRead, resourceAlibabacloudStackNasDirQuotaUpdate, resourceAlibabacloudStackNasDirQuotaDelete)
	return resource
}

func resourceAlibabacloudStackNasDirQuotaCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	nasService := NasService{client}
	fileSystemId := d.Get("file_system_id").(string)
	path := d.Get("path").(string)
	id := fmt.Sprintf("%s:%s", fileSystemId, path)
	request := map[string]interface{}{
		"FileSystemId": fileSystemId,
		"Path":         path,
		"UserType":     "AllUsers",
		"QuotaType":    "Accounting",
	}
	_, err := client.DoTeaRequest("POST", "nas", "2017-06-26", "SetDirQuota", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
			"alibabacloudstack_nas_dir_quota", "SetDirQuota", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{"Initializing"}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, nasService.NasDirQuotaStateRefreshFunc(id, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	quotas := d.Get("quotas").(*schema.Set).List()
	hasAllUser := false
	for _, v := range quotas {
		quota := v.(map[string]interface{})
		if quota["user_type"] == "AllUsers" {
			hasAllUser = true
		}
		request := map[string]interface{}{
			"FileSystemId": fileSystemId,
			"Path":         path,
			"UserType":     quota["user_type"].(string),
			"QuotaType":    quota["quota_type"].(string),
			"UserId":       "",
		}
		if v, ok := quota["user_id"]; ok {
			request["UserId"] = v.(string)
		}
		if v, ok := quota["size_limit"]; ok && quota["quota_type"].(string) == "Enforcement" {
			request["SizeLimit"] = v.(int)
		}
		if v, ok := quota["file_count_limit"]; ok && quota["quota_type"].(string) == "Enforcement" {
			request["FileCountLimit"] = v.(int)
		}
		_, err := client.DoTeaRequest("POST", "nas", "2017-06-26", "SetDirQuota", "", nil, request, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_nas_dir_quota", "SetDirQuota", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{"Initializing"}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, nasService.NasDirQuotaStateRefreshFunc(id, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	if !hasAllUser {
		request := map[string]interface{}{
			"FileSystemId": fileSystemId,
			"Path":         path,
			"UserType":     "AllUsers",
		}
		_, err := client.DoTeaRequest("POST", "nas", "2017-06-26", "CancelDirQuota", "", nil, request, nil)
		if err != nil {
		}
	}
	d.SetId(id)
	return nil
}

func resourceAlibabacloudStackNasDirQuotaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	nasService := NasService{client}

	object, err := nasService.DescribeNasDirQuota(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_nas_dir_quota nasService.DescribeNasDirQuota Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	fileSystemId := parts[0]
	d.Set("file_system_id", fileSystemId)
	d.Set("path", object["Path"])
	d.Set("status", object["Status"])
	userQuotaInfos, _ := object["UserQuotaInfos"].([]interface{})
	quotas := make([]map[string]interface{}, 0)
	for _, v := range userQuotaInfos {
		userQuotaInfo := v.(map[string]interface{})
		quota := map[string]interface{}{
			"user_type":  userQuotaInfo["UserType"],
			"quota_type": userQuotaInfo["QuotaType"],
		}
		if v, ok := userQuotaInfo["UserId"]; ok {
			quota["user_id"] = v
		}
		if v, ok := userQuotaInfo["SizeLimit"]; ok && fmt.Sprint(v) != "-1" {
			quota["size_limit"] = v
		}
		if v, ok := userQuotaInfo["FileCountLimit"]; ok && fmt.Sprint(v) != "-1" {
			quota["file_count_limit"] = v
		}

		quotas = append(quotas, quota)
	}
	d.Set("quotas", quotas)
	return nil
}

func resourceAlibabacloudStackNasDirQuotaUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.HasChanges("quotas") {
		nasService := NasService{client}
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		fileSystemId := parts[0]
		path := parts[1]
		log.Printf("[DEBUG] Resource alibabacloudstack_nas_dir_quota update quotas")
		o, n := d.GetChange("quotas")
		old, new := o.(*schema.Set), n.(*schema.Set)
		add := new.Difference(old)
		remove := old.Difference(new)
		updated_userdata := make([]string, 0)
		for _, v := range add.List() {
			quota := v.(map[string]interface{})
			userdata := quota["user_type"].(string)
			request := map[string]interface{}{
				"FileSystemId": fileSystemId,
				"Path":         path,
				"UserType":     quota["user_type"].(string),
				"QuotaType":    quota["quota_type"].(string),
			}
			if v, ok := quota["user_id"]; ok {
				request["UserId"] = v.(string)
				userdata = userdata + ":" + v.(string)
			}
			if v, ok := quota["size_limit"]; ok && quota["quota_type"].(string) == "Enforcement" {
				request["SizeLimit"] = v.(int)
			}
			if v, ok := quota["file_count_limit"]; ok && quota["quota_type"].(string) == "Enforcement" {
				request["FileCountLimit"] = v.(int)
			}
			_, err := client.DoTeaRequest("POST", "nas", "2017-06-26", "SetDirQuota", "", nil, request, nil)
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
					"alibabacloudstack_nas_dir_quota", "SetDirQuota", errmsgs.AlibabacloudStackSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{"Initializing"}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, nasService.NasDirQuotaStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
			}
			updated_userdata = append(updated_userdata, userdata)
		}
		for _, v := range remove.List() {
			quota := v.(map[string]interface{})
			userdata := quota["user_type"].(string)
			request := map[string]interface{}{
				"FileSystemId": fileSystemId,
				"Path":         path,
				"UserType":     quota["user_type"].(string),
				"QuotaType":    quota["quota_type"].(string),
			}
			if v, ok := quota["user_id"]; ok {
				request["UserId"] = v.(string)
				userdata = userdata + ":" + v.(string)
			}
			if slices.Contains(updated_userdata, userdata) {
				// already updated
				continue
			}
			_, err := client.DoTeaRequest("POST", "nas", "2017-06-26", "CancelDirQuota", "", nil, request, nil)
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
					"alibabacloudstack_nas_dir_quota", "CancelDirQuota", errmsgs.AlibabacloudStackSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{"Initializing"}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, nasService.NasDirQuotaStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
			}
		}
	}
	return nil
}

func resourceAlibabacloudStackNasDirQuotaDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	nasService := NasService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	fileSystemId := parts[0]
	path := parts[1]
	object, err := nasService.DescribeNasDirQuota(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			return nil
		}
		return errmsgs.WrapError(err)
	}

	userQuotaInfos, ok := object["UserQuotaInfos"].([]interface{})
	if !ok {
		return nil
	}
	for _, v := range userQuotaInfos {
		userQuotaInfo := v.(map[string]interface{})
		request := map[string]interface{}{
			"FileSystemId": fileSystemId,
			"Path":         path,
			"UserType":     userQuotaInfo["UserType"].(string),
			"QuotaType":    userQuotaInfo["QuotaType"].(string),
		}
		if v, ok := userQuotaInfo["UserId"]; ok {
			request["UserId"] = v.(string)
		}
		_, err := client.DoTeaRequest("POST", "nas", "2017-06-26", "CancelDirQuota", "", nil, request, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_nas_dir_quota", "CancelDirQuota", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		time.Sleep(3 * time.Second)
	}
	return nil
}
