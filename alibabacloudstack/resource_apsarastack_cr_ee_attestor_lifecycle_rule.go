package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackCrEEArtifactLifecycleRule() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"REPO", "NAMESPACE"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"retention_tag_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"tag_regexp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repo_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_delete_tag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"recent_pull_keep": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"recent_push_keep": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"schedule": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_delete_untagged_manifest": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackCrEEArtifactLifecycleRuleCreate, resourceAlibabacloudStackCrEEArtifactLifecycleRuleRead, resourceAlibabacloudStackCrEEArtifactLifecycleRuleUpdate, resourceAlibabacloudStackCrEEArtifactLifecycleRuleDelete)
	return resource
}

func resourceAlibabacloudStackCrEEArtifactLifecycleRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	scope := d.Get("scope").(string)
	reqQuery := map[string]interface{}{
		"InstanceId":        d.Get("instance_id"),
		"Scope":             scope,
		"Schedule":          "MANUAL",
		"NamespaceName":     d.Get("namespace_name"),
		"RetentionTagCount": d.Get("retention_tag_count").(int),
		"EnableDeleteTag":   d.Get("enable_delete_tag").(bool),
		"RecentPullKeep":    d.Get("recent_pull_keep").(int),
		"RecentPushKeep":    d.Get("recent_push_keep").(int),
		"Auto":              false,
	}
	if v, ok := d.GetOk("tag_regexp"); ok {
		reqQuery["TagRegexp"] = v.(string)
		reqQuery["sTagRegexp"] = true
	}
	if scope == "REPO" {
		reqQuery["RepoName"] = d.Get("repo_name")
	}

	response, err := client.DoTeaRequest("POST", "cr-ee", "2018-12-01", "CreateArtifactLifecycleRule", "", nil, reqQuery, nil)
	if err != nil {
		return err
	}

	// Get RuleId from response
	ruleId, ok := response["RuleId"]
	if !ok {
		return fmt.Errorf("RuleId not found in response")
	}

	// Generate resource ID
	resourceId := fmt.Sprintf("%s:%s", d.Get("instance_id").(string), ruleId.(string))
	d.SetId(resourceId)

	return nil
}

func resourceAlibabacloudStackCrEEArtifactLifecycleRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := CrService{client}

	object, err := crService.DescribeCrEEArtifactLifecycleRule(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	if v, ok := object["Scope"]; ok && v.(string) != "" {
		d.Set("scope", v.(string))
	}
	d.Set("rule_id", object["RuleId"])
	d.Set("schedule", object["Schedule"])
	d.Set("retention_tag_count", object["RetentionTagCount"])
	d.Set("tag_regexp", object["TagRegexp"])
	d.Set("namespace_name", object["NamespaceName"])
	d.Set("repo_name", object["RepoName"])
	d.Set("enable_delete_tag", object["EnableDeleteTag"])
	d.Set("auto", object["Auto"])
	d.Set("enable_delete_untagged_manifest", object["EnableDeleteUntaggedManifest"])
	d.Set("modified_time", object["ModifiedTime"])
	d.Set("create_time", object["CreateTime"])
	d.Set("recent_pull_keep", object["RecentPullKeep"])
	d.Set("recent_push_keep", object["RecentPushKeep"])

	return nil
}

func resourceAlibabacloudStackCrEEArtifactLifecycleRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.IsNewResource() {
		return nil
	}

	if d.HasChanges("scope", "tag_regexp", "namespace_name", "retention_tag_count", "recent_pull_keep", "recent_push_keep") {
		scope := d.Get("scope").(string)
		crService := CrService{client}
		strRet := crService.ParseResourceId(d.Id())
		instanceId := strRet[0]
		ruleId := strRet[1]
		reqQuery := map[string]interface{}{
			"RuleId":            ruleId,
			"InstanceId":        instanceId,
			"Scope":             scope,
			"Schedule":          "MANUAL",
			"NamespaceName":     d.Get("namespace_name"),
			"RetentionTagCount": d.Get("retention_tag_count").(int),
			"EnableDeleteTag":   d.Get("enable_delete_tag").(bool),
			"RecentPullKeep":    d.Get("recent_pull_keep").(int),
			"RecentPushKeep":    d.Get("recent_push_keep").(int),
			"Auto":              false,
		}
		if v, ok := d.GetOk("tag_regexp"); ok {
			reqQuery["TagRegexp"] = v.(string)
			reqQuery["sTagRegexp"] = true
		}
		if scope == "REPO" {
			reqQuery["RepoName"] = d.Get("repo_name")
		}
		_, err := client.DoTeaRequest("POST", "cr-ee", "2018-12-01", "UpdateArtifactLifecycleRule", "", nil, reqQuery, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceAlibabacloudStackCrEEArtifactLifecycleRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := CrService{client}
	strRet := crService.ParseResourceId(d.Id())
	instanceId := strRet[0]
	ruleId := strRet[1]
	request := map[string]interface{}{
		"InstanceId": instanceId,
		"RuleId":     ruleId,
	}

	_, err := client.DoTeaRequest("POST", "cr-ee", "2018-12-01", "DeleteArtifactLifecycleRule", "", nil, request, nil)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteArtifactLifecycleRule", errmsgs.AlibabacloudStackSdkGoERROR, err.Error())
	}

	return nil
}
