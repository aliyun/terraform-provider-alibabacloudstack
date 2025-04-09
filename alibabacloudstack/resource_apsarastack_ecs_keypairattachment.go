package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackKeyPairAttachment() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
				ForceNew:     true,
			},
			"instance_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				ForceNew: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackKeyPairAttachmentCreate, resourceAlibabacloudStackKeyPairAttachmentRead, nil, resourceAlibabacloudStackKeyPairAttachmentDelete)
	return resource
}

func resourceAlibabacloudStackKeyPairAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	keyName := d.Get("key_name").(string)
	instanceIds := d.Get("instance_ids").(*schema.Set).List()
	force := d.Get("force").(bool)
	idsMap := make(map[string]string)
	var newIds []string
	if force {
		ids, _, err := ecsService.QueryInstancesWithKeyPair("", keyName)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		for _, id := range ids {
			idsMap[id] = id
		}
		for _, id := range instanceIds {
			if _, ok := idsMap[id.(string)]; !ok {
				newIds = append(newIds, id.(string))
			}
		}
	}

	if err := ecsService.AttachKeyPair(keyName, instanceIds); err != nil {
		return errmsgs.WrapError(err)
	}

	if force {
		request := ecs.CreateRebootInstanceRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.ForceStop = requests.NewBoolean(true)
		for _, id := range newIds {
			request.InstanceId = id
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.RebootInstance(request)
			})
			if err != nil {
				errmsg := ""
				if bresponse, ok := raw.(*ecs.RebootInstanceResponse); ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}
		for _, id := range newIds {
			if err := ecsService.WaitForEcsInstance(id, Running, DefaultLongTimeout); err != nil {
				return errmsgs.WrapError(err)
			}
		}
	}

	d.SetId(keyName + ":" + convertListToJsonString(instanceIds))

	return nil
}

func resourceAlibabacloudStackKeyPairAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	keyName := strings.Split(d.Id(), ":")[0]
	object, err := ecsService.DescribeKeyPairAttachment(d.Id())

	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("key_name", object.KeyPairName)
	if ids, ok := d.GetOk("instance_ids"); ok {
		d.Set("instance_ids", ids)
	} else {
		ids, _, err := ecsService.QueryInstancesWithKeyPair("", keyName)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		d.Set("instance_ids", ids)
	}
	return nil
}

func resourceAlibabacloudStackKeyPairAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	keyName := strings.Split(d.Id(), ":")[0]
	instanceIds := strings.Split(d.Id(), ":")[1]

	request := ecs.CreateDetachKeyPairRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.KeyPairName = keyName

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		request.InstanceIds = instanceIds
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DetachKeyPair(request)
		})
		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*ecs.DetachKeyPairResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		instance_ids, _, err := ecsService.QueryInstancesWithKeyPair(instanceIds, keyName)
		if err != nil {
			return resource.NonRetryableError(errmsgs.WrapError(err))
		}
		if len(instance_ids) > 0 {
			var ids []interface{}
			for _, id := range instance_ids {
				ids = append(ids, id)
			}
			instanceIds = convertListToJsonString(ids)
			return resource.RetryableError(errmsgs.WrapError(fmt.Errorf("detach Key Pair timeout and the instances including %s has not yet been detached. ", instanceIds)))
		}

		return nil
	})
}
