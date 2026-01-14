package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackMaxcomputeProjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackMaxcomputeProjectsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				MinItems: 1,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
				ForceNew: true,
			},
			"projects": {
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
						"account": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_pk": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"core_arch": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_ids": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackMaxcomputeProjectsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("POST", "dataworks-private-cloud", "2019-01-17", "ListCalcEnginesForAscm", "")
	var name, status string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	if v, ok := d.GetOk("status"); ok {
		status = v.(string)
	}
	
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		errmsg := ""
		errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_project", "ListCalcEnginesForAscm", errmsg)
	}

	response := &MaxComputeProjectDetailResponse{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	idsMap := getIdsStringFilter(d)
	maxcomputeService := MaxcomputeService{client}
	var t []map[string]interface{}
	var ids []string
	for _, object := range response.Data.CalcEngines {
		if status != "" && status != strconv.Itoa(object.EngineStatus) {
			continue
		}
		id := fmt.Sprint(object.EngineId)
		if _, existed := idsMap[id]; len(idsMap) > 0 && !existed {
			continue
		}
		if name != "" && name != object.Name {
			continue
		}
		quota, err := maxcomputeService.ListOdpsEngineQuotaForAscm(id, object.Name)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		account, err := maxcomputeService.DescribeMaxcomputeProjectEngine(id)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		Properties, err := maxcomputeService.DescribeMaxProjectPropertiesForAscm(id, object.Name)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		vpcData := make([]string, 0)
		vpcs := make([]string, 0)
		if Properties["odps.security.vpc.whitelist"].(string) != "" {
			err = json.Unmarshal([]byte(Properties["odps.security.vpc.whitelist"].(string)), &vpcData)
			if err != nil {
				return errmsgs.WrapError(err)
			}
			for _, v := range vpcData {
				vpc := strings.Split(v, "_")
				if strings.HasPrefix(vpc[1], "vpc") {
					vpcs = append(vpcs, vpc[1])
				}
			}
		}
		project := map[string]interface{}{
			"id":         id,
			"name":       object.Name,
			"core_arch":  object.EngineInfo.DefaultClusterArch,
			"quota_id":   fmt.Sprint(quota.QuotaId),
			"disk":       int(quota.Disk * 1024),
			"account":    account.EngineInfo.TaskAk.AliyunAccount,
			"account_pk": fmt.Sprint(account.EngineInfo.TaskAk.Kp),
			"vpc_ids":    vpcs,
		}
		t = append(t, project)
		ids = append(ids, id)

	}
	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("projects", t); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
