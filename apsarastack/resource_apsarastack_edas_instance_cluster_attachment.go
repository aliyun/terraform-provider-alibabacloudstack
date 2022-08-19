package apsarastack

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackEdasInstanceClusterAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackEdasInstanceClusterAttachmentCreate,
		Read:   resourceApsaraStackEdasInstanceClusterAttachmentRead,
		Delete: resourceApsaraStackEdasInstanceClusterAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
			"pass_word": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status_map": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Computed: true,
				ForceNew: true,
			},
			"ecu_map": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
				ForceNew: true,
			},
			"cluster_member_ids": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceApsaraStackEdasInstanceClusterAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	edasService := EdasService{client}

	clusterId := d.Get("cluster_id").(string)
	instanceIds := d.Get("instance_ids").([]interface{})
	aString := make([]string, len(instanceIds))
	for i, v := range instanceIds {
		aString[i] = v.(string)
	}

	request := edas.CreateInsertClusterMemberRequest()
	request.ClusterId = clusterId
	request.RegionId = client.RegionId
	request.Password = d.Get("pass_word").(string)
	request.Headers["x-ascm-product-name"] = "Edas"
	request.Headers["x-acs-organizationid"] = client.Department
	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	request.InstanceIds = strings.Join(aString, ",")
	request.SetReadTimeout(30 * time.Second)

	if err := edasService.SyncResource("ecs"); err != nil {
		return err
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.InsertClusterMember(request)
		})
		if err != nil {
			if IsThrottling(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		response, _ := raw.(*edas.InsertClusterMemberResponse)

		if response.Code != 200 {
			return resource.NonRetryableError(Error("insert instances to cluster failed for " + response.Message))
		}

		d.SetId(clusterId + ":" + strings.Join(aString, ","))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsaraStack_edas_instance_cluster_attachment", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	//使用这种方法有问题，超过120s还没有导入成功，则用例失败，需要改成
	//time.Sleep(120 * time.Second)
	var cnt int
	ImportSuccessFlag := false
	for {
		if cnt >= 10{
			break
		}
		requestList := edas.CreateListClusterMembersRequest()
		requestList.RegionId = client.RegionId
		requestList.ClusterId = clusterId
		requestList.Headers["x-ascm-product-name"] = "Edas"
		requestList.Headers["x-acs-organizationid"] = client.Department
		requestList.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
		rawList, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.ListClusterMembers(requestList)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_edas_List_attachment", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		strs, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return WrapError(err)
		}

		instanceIdstr := strs[1]
		responseList := rawList.(*edas.ListClusterMembersResponse)
		for _, member := range responseList.ClusterMemberPage.ClusterMemberList.ClusterMember {
			if strings.Contains(instanceIdstr, member.EcsId) {
				if member.EcuId != "" {
					ImportSuccessFlag = true
					break
				}
			}
		}
		if ImportSuccessFlag == true {
			break
		}
		time.Sleep(30 * time.Second)
		cnt++
	}
	return resourceApsaraStackEdasInstanceClusterAttachmentRead(d, meta)
}

func resourceApsaraStackEdasInstanceClusterAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	edasService := EdasService{client}

	strs, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	clusterId := strs[0]
	regionId := client.RegionId
	instanceIdstr := strs[1]

	request := edas.CreateListClusterMembersRequest()
	request.RegionId = regionId
	request.ClusterId = clusterId
	request.Headers["x-ascm-product-name"] = "Edas"
	request.Headers["x-acs-organizationid"] = client.Department
	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListClusterMembers(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_edas_instance_cluster_attachment", request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	statusMap := make(map[string]int)
	ecuMap := make(map[string]string)
	memMap := make(map[string]string)
	response := raw.(*edas.ListClusterMembersResponse)
	for _, member := range response.ClusterMemberPage.ClusterMemberList.ClusterMember {
		if strings.Contains(instanceIdstr, member.EcsId) {
			statusMap[member.EcsId] = member.Status
			ecuMap[member.EcsId] = member.EcuId
			memMap[member.EcsId] = member.ClusterMemberId
		}
	}

	d.Set("status_map", statusMap)
	d.Set("ecu_map", ecuMap)
	d.Set("cluster_member_ids", memMap)

	return nil
}

//有问题 单个实例删除失败会影响整个过程
func resourceApsaraStackEdasInstanceClusterAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	edasService := EdasService{client}

	memIds := d.Get("cluster_member_ids").(map[string]interface{})
	for instanceId, memberId := range memIds {
		request := edas.CreateDeleteClusterMemberRequest()
		request.RegionId = client.RegionId
		request.ClusterId = d.Get("cluster_id").(string)
		request.ClusterMemberId = memberId.(string)
		request.Headers["x-ascm-product-name"] = "Edas"
		request.Headers["x-acs-organizationid"] = client.Department
		request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(1*time.Minute, func() *resource.RetryError {
			raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
				return edasClient.DeleteClusterMember(request)

			})
			if err != nil {
				if IsThrottling(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RoaRequest, request)
			response, _ := raw.(*edas.DeleteClusterMemberResponse)
			if strings.Contains(response.Message, "there are still applications deployed in this cluster") {
				err = Error("there are still applications deployed in this cluster")
				return resource.RetryableError(err)
			} else if response.Code != 200 {
				return resource.NonRetryableError(Error("delete instance:" + instanceId + " from cluster failed for " + response.Message))
			}

			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_edas_instance_cluster_attachment", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
	}

	return nil
}
