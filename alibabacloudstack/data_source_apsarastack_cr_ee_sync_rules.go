package alibabacloudstack

import (
	"regexp"
	"sort"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackCrEeSyncRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCrEeSyncRulesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"namespace_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repo_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repo_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_namespace_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_repo_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_filter": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sync_scope": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sync_direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sync_trigger": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackCrEeSyncRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := &CrService{client}
	instanceId := d.Get("instance_id").(string)

	var (
		namespaceName    string
		repoName         string
		targetInstanceId string
	)

	if v, ok := d.GetOk("namespace_name"); ok {
		namespaceName = v.(string)
	}
	if v, ok := d.GetOk("repo_name"); ok {
		repoName = v.(string)
	}
	if v, ok := d.GetOk("target_instance_id"); ok {
		targetInstanceId = v.(string)
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var idsMap map[string]string
	if v, ok := d.GetOk("ids"); ok {
		idsMap = make(map[string]string)
		for _, vv := range v.([]interface{}) {
			idsMap[vv.(string)] = vv.(string)
		}
	}

	pageNo, pageSize := 1, 50
	var syncRules []cr_ee.SyncRulesItem
	for {
		request := cr_ee.CreateListRepoSyncRuleRequest()
		client.InitRpcRequest(*request.RpcRequest)

		request.RegionId = crService.client.RegionId
		request.InstanceId = instanceId
		if namespaceName != "" {
			request.NamespaceName = namespaceName
		}
		if repoName != "" {
			request.RepoName = repoName
		}
		if targetInstanceId != "" {
			request.TargetInstanceId = targetInstanceId
		}

		request.PageNo = requests.NewInteger(pageNo)
		request.PageSize = requests.NewInteger(pageSize)
		raw, err := crService.client.WithCrEeClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.ListRepoSyncRule(request)
		})
		bresponse, ok := raw.(*cr_ee.ListRepoSyncRuleResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cr_ee_sync_rules", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		// if !response.ListRepoSyncRuleIsSuccess {
		// 	return crService.wrapCrServiceError("alibabacloudstack_cr_ee_sync_rules", request.GetActionName(), response.Code)
		// }

		for _, rule := range bresponse.SyncRules {
			if nameRegex != nil && !nameRegex.MatchString(rule.SyncRuleName) {
				continue
			}
			if idsMap != nil && idsMap[rule.SyncRuleId] == "" {
				continue
			}
			syncRules = append(syncRules, rule)
		}

		if len(bresponse.SyncRules) < pageSize {
			break
		}

		pageNo++
	}

	sort.SliceStable(syncRules, func(i, j int) bool {
		return syncRules[i].CreateTime < syncRules[j].CreateTime
	})

	ids := make([]string, len(syncRules))
	names := make([]string, len(syncRules))
	rulesMaps := make([]map[string]interface{}, len(syncRules))
	for i, r := range syncRules {
		ids[i] = r.SyncRuleId
		names[i] = r.SyncRuleName
		m := make(map[string]interface{})
		m["region_id"] = r.LocalRegionId
		m["instance_id"] = r.LocalInstanceId
		m["namespace_name"] = r.LocalNamespaceName
		m["repo_name"] = r.LocalRepoName
		m["target_region_id"] = r.TargetRegionId
		m["target_instance_id"] = r.TargetInstanceId
		m["target_namespace_name"] = r.TargetNamespaceName
		m["target_repo_name"] = r.TargetRepoName
		m["tag_filter"] = r.TagFilter
		m["sync_scope"] = r.SyncScope
		m["sync_direction"] = r.SyncDirection
		m["sync_trigger"] = r.SyncTrigger
		m["id"] = r.SyncRuleId
		m["name"] = r.SyncRuleName
		rulesMaps[i] = m
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("rules", rulesMaps); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), rulesMaps); err != nil {
			return err
		}
	}

	return nil
}
