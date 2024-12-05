package alibabacloudstack

import (
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEdasInstanceClusterAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackEdasInstanceClusterAttachmentCreate,
		Read:   resourceAlibabacloudStackEdasInstanceClusterAttachmentRead,
		Delete: resourceAlibabacloudStackEdasInstanceClusterAttachmentDelete,
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
				Type:     schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Computed: true,
				ForceNew: true,
			},
			"ecu_map": {
				Type:     schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
				ForceNew: true,
			},
			"cluster_member_ids": {
				Type:     schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlibabacloudStackEdasInstanceClusterAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	clusterId := d.Get("cluster_id").(string)
	instanceIds := d.Get("instance_ids").([]interface{})
	aString := make([]string, len(instanceIds))
	for i, v := range instanceIds {
		aString[i] = v.(string)
	}

	request := edas.CreateInsertClusterMemberRequest()
	client.InitRoaRequest(*request.RoaRequest)
	request.ClusterId = clusterId
	request.Password = d.Get("pass_word").(string)
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
		bresponse, ok := raw.(*edas.InsertClusterMemberResponse)
		if err != nil {
			if errmsgs.IsThrottling(err) {
				wait()
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_instance_cluster_attachment", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)

		if bresponse.Code != 200 {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.Error("insert instances to cluster failed for " + bresponse.Message + " " + errmsg))
		}

		d.SetId(clusterId + ":" + strings.Join(aString, ","))
		return nil
	})
	if err != nil {
		return err
	}

	var cnt int
	ImportSuccessFlag := false
	for {
		if cnt >= 5 {
			break
		}
		requestList := edas.CreateListClusterMembersRequest()
		client.InitRoaRequest(*request.RoaRequest)
		requestList.ClusterId = clusterId
		rawList, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.ListClusterMembers(requestList)
		})
		bresponseList, ok := rawList.(*edas.ListClusterMembersResponse)
		addDebug(requestList.GetActionName(), rawList, requestList.RoaRequest, requestList)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponseList.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_List_attachment", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		strs, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		instanceIdstr := strs[1]
		for _, member := range bresponseList.ClusterMemberPage.ClusterMemberList.ClusterMember {
			log.Printf("===================================  instance status: %d ecsId : %s", member.Status, member.EcsId)
			if strings.Contains(instanceIdstr, member.EcsId) {
				if member.Status == 1 {
					ImportSuccessFlag = true
					break
				}
				if member.Status == 7 {
					errmsg := ""
					if ok {
						errmsg = errmsgs.GetBaseResponseErrorMessage(bresponseList.BaseResponse)
					}
					return errmsgs.Error("Instance:`%s` Import Timeout! " + errmsg, member.EcsId)
				}
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
	return resourceAlibabacloudStackEdasInstanceClusterAttachmentRead(d, meta)
}

func resourceAlibabacloudStackEdasInstanceClusterAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	strs, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	clusterId := strs[0]
	regionId := client.RegionId
	instanceIdstr := strs[1]

	request := edas.CreateListClusterMembersRequest()
	client.InitRoaRequest(*request.RoaRequest)
	request.ClusterId = clusterId
	request.RegionId = regionId
	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListClusterMembers(request)
	})
	bresponse, ok := raw.(*edas.ListClusterMembersResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_instance_cluster_attachment", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	statusMap := make(map[string]int)
	ecuMap := make(map[string]string)
	memMap := make(map[string]string)
	for _, member := range bresponse.ClusterMemberPage.ClusterMemberList.ClusterMember {
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

func resourceAlibabacloudStackEdasInstanceClusterAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	memIds := d.Get("cluster_member_ids").(map[string]interface{})
	for instanceId, memberId := range memIds {
		request := edas.CreateDeleteClusterMemberRequest()
		client.InitRoaRequest(*request.RoaRequest)
		request.ClusterId = d.Get("cluster_id").(string)
		request.ClusterMemberId = memberId.(string)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(1*time.Minute, func() *resource.RetryError {
			raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
				return edasClient.DeleteClusterMember(request)
			})
			bresponse, ok := raw.(*edas.DeleteClusterMemberResponse)
			if err != nil {
				if errmsgs.IsThrottling(err) {
					wait()
					return resource.RetryableError(err)
				}
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_instance_cluster_attachment", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			addDebug(request.GetActionName(), raw, request.RoaRequest, request)
			if strings.Contains(bresponse.Message, "there are still applications deployed in this cluster") {
				err = errmsgs.Error("there are still applications deployed in this cluster")
				return resource.RetryableError(err)
			} else if bresponse.Code != 200 {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return resource.NonRetryableError(errmsgs.Error("delete instance:" + instanceId + " from cluster failed for " + bresponse.Message + " " + errmsg))
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}
