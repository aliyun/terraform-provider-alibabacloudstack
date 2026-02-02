package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEdasClusterMember() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_member_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEdasClusterMemberCreate, resourceAlibabacloudStackEdasClusterMemberRead, nil, resourceAlibabacloudStackEdasClusterMemberDelete)
	return resource
}

func resourceAlibabacloudStackEdasClusterMemberCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	clusterId := d.Get("cluster_id").(string)
	instanceId := d.Get("instance_id").(string)
	request := map[string]interface{}{
		"ClusterId":   clusterId,
		"InstanceIds": instanceId,
		"DoAsync":     true,
	}

	_, err := client.DoTeaRequest("POST", "Edas", "2017-08-01", "InstallAgent", "/pop/v5/ecss/install_agent", nil, request, nil)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	id := fmt.Sprintf("%s:%s", clusterId, instanceId)
	d.SetId(id)
	edasService := EdasService{client}
	wait := incrementalWait(1*time.Second, 10*time.Second)
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		_, err := edasService.DescribeClusterMember(d.Id())
		if err != nil {
			if errmsgs.NotFoundError(err) {
				wait()
				return resource.RetryableError(err)
			}
			errmsg := ""
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DescribeClusterMember", errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		return nil
	})
	return err
}

func resourceAlibabacloudStackEdasClusterMemberRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	object, err := edasService.DescribeClusterMember(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) != 2 {
		return errmsgs.GetNotFoundErrorFromString("Resource not found")
	}

	d.Set("cluster_id", object["ClusterId"])
	d.Set("instance_id", object["EcsId"])
	d.Set("cluster_member_id", object["ClusterMemberId"])

	return nil
}

func resourceAlibabacloudStackEdasClusterMemberDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) != 2 {
		return errmsgs.GetNotFoundErrorFromString("Resource not found")
	}

	reqQuery := map[string]interface{}{
		"clusterId":       parts[0],
		"clusterMemberId": d.Get("cluster_member_id"),
	}

	_, err := client.DoTeaRequest("DELETE", "Edas", "2017-08-01", "DeleteClusterMember", "/pop/v5/resource/cluster_member", nil, reqQuery, nil)

	if err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
