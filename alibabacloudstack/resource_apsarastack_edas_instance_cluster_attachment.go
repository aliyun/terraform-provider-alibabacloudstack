package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
	id := fmt.Sprintf("%s:%s", clusterId, strings.Join(aString, ","))
	d.SetId(id)

	// Retry InstallAgent for transient errors like 503
	var installErr error
	_ = resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, installErr = client.DoTeaRequest("POST", "Edas", "2017-08-01", "InstallAgent", "/pop/v5/ecss/install_agent", nil, request, nil)
		if installErr != nil {
			log.Printf("[DEBUG] InstallAgent error for %s: %v", id, installErr)
			if errmsgs.NeedRetry(installErr) {
				return resource.RetryableError(installErr)
			}
			return resource.NonRetryableError(installErr)
		}
		return nil
	})
	if installErr != nil {
		return errmsgs.WrapError(installErr)
	}

	// Wait until all instances appear in cluster member list
	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		objects, e := edasService.DescribeClusterMember(id)
		if e != nil {
			log.Printf("[DEBUG] DescribeClusterMember error for %s: %v", id, e)
			if errmsgs.NotFoundError(e) || errmsgs.NeedRetry(e) {
				return resource.RetryableError(e)
			}
			return resource.NonRetryableError(e)
		}
		instanceIdstr := strings.Join(aString, ",")
		foundCount := 0
		for _, v := range objects {
			member := v.(map[string]interface{})
			if ecsId, ok := member["EcsId"].(string); ok && strings.Contains(instanceIdstr, ecsId) {
				foundCount++
			}
		}
		log.Printf("[DEBUG] DescribeClusterMember for %s: total members=%d, found=%d/%d", id, len(objects), foundCount, len(aString))
		if foundCount < len(aString) {
			return resource.RetryableError(fmt.Errorf("waiting for all instances to be attached to cluster, found %d/%d", foundCount, len(aString)))
		}
		return nil
	})
	if err != nil {
		return errmsgs.WrapError(err)
	}

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
