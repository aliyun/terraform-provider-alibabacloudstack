package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackEdasClusterMembers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackEdasClusterMembersRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of EDAS cluster install agent IDs. Each ID is formatted as 'ClusterId:InstanceId'.",
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"names": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of instance IDs corresponding to the retrieved agents.",
			},
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the agent, formatted as 'ClusterId:InstanceId'.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the EDAS cluster.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the ECS instance (same as cluster_member_id).",
						},
						"ecu_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the ECU.",
						},
						"ecs_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the ECS instance.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The status of the cluster member.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The timestamp when the cluster member was created.",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The timestamp when the cluster member was last updated.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackEdasClusterMembersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	clusterId := d.Get("cluster_id").(string)

	request := map[string]interface{}{
		"clusterId": clusterId,
		"pageSize":  100, // Set a large page size to get all results at once
	}
	pageNumber := 1
	result := make([]interface{}, 0)
	for {
		request["currentPage"] = pageNumber
		response, err := client.DoTeaRequest("GET", "Edas", "2017-08-01", "ListClusterMembers", "/pop/v5/resource/cluster_member_list", nil, request, nil)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if instances, err := jsonpath.Get("$.ClusterMemberPage.ClusterMemberList.ClusterMember", response); err == nil {
			result = append(result, instances.([]interface{})...)
		}
		totalSize, _ := response["TotalSize"].(json.Number).Int64()
		if int(totalSize) <= pageNumber*100 {
			break
		}
		pageNumber++

	}

	idsMap := getIdsStringFilter(d)

	var ids []string
	members := make([]interface{}, 0)
	for _, data := range result {
		member := data.(map[string]interface{})
		id := fmt.Sprintf("%s:%s", member["ClusterId"], member["InstanceId"])
		if description_regex, ok := connectivity.GetResourceDataOk(d, "description_regex", "name_regex"); ok {
			r := regexp.MustCompile(description_regex.(string))
			if !r.MatchString(member["InstanceName"].(string)) {
				continue
			}
		}

		if len(idsMap) > 0 {
			if _, exist := idsMap[id]; !exist {
				continue
			}
		}
		i := map[string]interface{}{
			"id":          id,
			"cluster_id":  member["clusterId"],
			"instance_id": member["instanceId"],
			"ecu_id":      member["ecuId"],
			"ecs_id":      member["ecsId"],
			"status":      member["status"],
			"create_time": member["createTime"],
			"update_time": member["updateTime"],
		}

		members = append(members, i)
		ids = append(ids, id)
	}

	// Set the data source ID to a hash of the IDs
	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("members", members); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
