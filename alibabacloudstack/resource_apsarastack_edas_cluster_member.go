package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
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
	}
	edasService := EdasService{client}
	client.Config.ClientReadTimeout = 120
	client.Config.ClientConnectTimeout = 120
	_, err := client.DoTeaRequest("POST", "Edas", "2017-08-01", "InstallAgent", "/pop/v5/ecss/install_agent", nil, request, nil)
	id := fmt.Sprintf("%s:%s", clusterId, instanceId)
	if err != nil {
		_, e := edasService.DescribeClusterMember(id)
		if e != nil {
			return errmsgs.WrapError(err)
		}
	}
	d.SetId(id)
	return nil
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
	log.Printf(" =================================================================== Reading EDAS Cluster Member: %#v", object)
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
		"ClusterId":       parts[0],
		"ClusterMemberId": d.Get("cluster_member_id"),
	}

	_, err := client.DoTeaRequest("DELETE", "Edas", "2017-08-01", "DeleteClusterMember", "/pop/v5/resource/cluster_member", nil, reqQuery, nil)

	if err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
