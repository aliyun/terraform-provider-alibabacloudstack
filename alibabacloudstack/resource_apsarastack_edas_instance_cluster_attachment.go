package alibabacloudstack

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEdasInstanceClusterAttachment() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				ForceNew: true,
			},
			// "pass_word": {
			// 	Type:     schema.TypeString,
			// 	Required: true,
			// 	ForceNew: true,
			// },
			"status_map": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Computed: true,
			},
			"ecu_map": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"cluster_member_ids": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEdasInstanceClusterAttachmentCreate, resourceAlibabacloudStackEdasInstanceClusterAttachmentRead, nil, resourceAlibabacloudStackEdasInstanceClusterAttachmentDelete)
	return resource
}

func resourceAlibabacloudStackEdasInstanceClusterAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	clusterId := d.Get("cluster_id").(string)
	instanceIds := d.Get("instance_ids").([]interface{})
	aString := make([]string, len(instanceIds))
	for i, v := range instanceIds {
		aString[i] = v.(string)
	}
	request := map[string]interface{}{
		"ClusterId":   clusterId,
		"InstanceIds": strings.Join(aString, ","),
	}
	edasService := EdasService{client}
	client.Config.ClientReadTimeout = 120
	client.Config.ClientConnectTimeout = 120
	_, err := client.DoTeaRequest("POST", "Edas", "2017-08-01", "InstallAgent", "/pop/v5/ecss/install_agent", nil, request, nil)
	id := fmt.Sprintf("%s:%s", clusterId, strings.Join(aString, ","))
	if err != nil {
		_, e := edasService.DescribeClusterMember(id)
		if e != nil {
			return errmsgs.WrapError(err)
		}
	}
	d.SetId(id)
	return nil
}

func resourceAlibabacloudStackEdasInstanceClusterAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	objects, err := edasService.DescribeClusterMember(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	strs, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	instanceIdstr := strs[1]
	instances := strings.Split(instanceIdstr, ",")
	statusMap := make(map[string]interface{})
	ecuMap := make(map[string]string)
	memMap := make(map[string]string)
	for _, v := range objects {
		member := v.(map[string]interface{})
		if strings.Contains(instanceIdstr, member["EcsId"].(string)) {
			statusMap[member["EcsId"].(string)] = member["Status"]
			ecuMap[member["EcsId"].(string)] = member["EcuId"].(string)
			memMap[member["EcsId"].(string)] = member["ClusterMemberId"].(string)
		}
	}
	d.Set("cluster_id", strs[0])
	d.Set("instance_ids", instances)
	d.Set("status_map", statusMap)
	d.Set("ecu_map", ecuMap)
	d.Set("cluster_member_ids", memMap)
	return nil
}

func resourceAlibabacloudStackEdasInstanceClusterAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) != 2 {
		return errmsgs.GetNotFoundErrorFromString("Resource not found")
	}

	cluster_member_ids := d.Get("cluster_member_ids").(map[string]interface{})
	for _, v := range cluster_member_ids {
		if v == nil {
			continue
		}

		reqQuery := map[string]interface{}{
			"ClusterId":       parts[0],
			"ClusterMemberId": v.(string),
		}

		_, err := client.DoTeaRequest("DELETE", "Edas", "2017-08-01", "DeleteClusterMember", "/pop/v5/resource/cluster_member", nil, reqQuery, nil)

		if err != nil {
			return errmsgs.WrapError(err)
		}
	}

	return nil
}
