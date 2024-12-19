package alibabacloudstack

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAlikafkaTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackAlikafkaTopicCreate,
		Update: resourceAlibabacloudStackAlikafkaTopicUpdate,
		Read:   resourceAlibabacloudStackAlikafkaTopicRead,
		Delete: resourceAlibabacloudStackAlikafkaTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topic": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"local_topic": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"compact_topic": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"partition_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      12,
				ValidateFunc: validation.IntBetween(0, 360),
			},
			"remark": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlibabacloudStackAlikafkaTopicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	instanceId := d.Get("instance_id").(string)
	topic := d.Get("topic").(string)

	request := alikafka.CreateCreateTopicRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = instanceId
	request.Topic = topic
	request.LocalTopic = requests.NewBoolean(d.Get("local_topic").(bool))
	request.CompactTopic = requests.NewBoolean(d.Get("compact_topic").(bool))
	request.PartitionNum = strconv.Itoa(d.Get("partition_num").(int))
	log.Printf("------------------ LocalTopic:%t CompactTopic:%t PartitionNum:%d", d.Get("local_topic").(bool), d.Get("local_topic").(bool), d.Get("partition_num").(int))
	// if v, ok := d.GetOk("local_topic"); ok {
	// 	request.LocalTopic = requests.NewBoolean(v.(bool))
	// }
	// if v, ok := d.GetOk("compact_topic"); ok {
	// 	request.CompactTopic = requests.NewBoolean(v.(bool))
	// }
	// if v, ok := d.GetOk("partition_num"); ok {
	// 	request.PartitionNum = strconv.Itoa(v.(int))
	// }
	if v, ok := d.GetOk("remark"); ok {
		request.Remark = v.(string)
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.CreateTopic(request)
		})
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_alikafka_topic", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.SetId(instanceId + ":" + topic)

	return resourceAlibabacloudStackAlikafkaTopicUpdate(d, meta)
}

func resourceAlibabacloudStackAlikafkaTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}
	d.Partial(true)
	// if err := alikafkaService.SetResourceTags(d, "topic"); err != nil {
	// 	return errmsgs.WrapError(err)
	// }
	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlibabacloudStackAlikafkaTopicRead(d, meta)
	}

	instanceId := d.Get("instance_id").(string)
	if d.HasChange("remark") {
		remark := d.Get("remark").(string)
		topic := d.Get("topic").(string)

		modifyRemarkRequest := alikafka.CreateModifyTopicRemarkRequest()
		client.InitRpcRequest(*modifyRemarkRequest.RpcRequest)
		modifyRemarkRequest.InstanceId = instanceId
		modifyRemarkRequest.Topic = topic
		modifyRemarkRequest.Remark = remark

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
				return alikafkaClient.ModifyTopicRemark(modifyRemarkRequest)
			})
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(modifyRemarkRequest.GetActionName(), raw, modifyRemarkRequest.RpcRequest, modifyRemarkRequest)
			return nil
		})
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), modifyRemarkRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	if d.HasChange("partition_num") {
		o, n := d.GetChange("partition_num")
		oldPartitionNum := o.(int)
		newPartitionNum := n.(int)

		if newPartitionNum < oldPartitionNum {
			return errmsgs.WrapError(errors.New("partition_num only support adjust to a greater value."))
		} else {
			topic := d.Get("topic").(string)

			modifyPartitionReq := alikafka.CreateModifyPartitionNumRequest()
			client.InitRpcRequest(*modifyPartitionReq.RpcRequest)
			modifyPartitionReq.InstanceId = instanceId
			modifyPartitionReq.Topic = topic
			modifyPartitionReq.AddPartitionNum = requests.NewInteger(newPartitionNum - oldPartitionNum)

			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
					return alikafkaClient.ModifyPartitionNum(modifyPartitionReq)
				})
				if err != nil {
					if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
						time.Sleep(10 * time.Second)
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(modifyPartitionReq.GetActionName(), raw, modifyPartitionReq.RpcRequest, modifyPartitionReq)
				return nil
			})
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), modifyPartitionReq.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
			}
		}
	}

	d.Partial(false)
	return resourceAlibabacloudStackAlikafkaTopicRead(d, meta)
}

func resourceAlibabacloudStackAlikafkaTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	object, err := alikafkaService.DescribeAlikafkaTopic(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("instance_id", object.InstanceId)
	d.Set("topic", object.Topic)
	d.Set("local_topic", object.LocalTopic)
	// d.Set("compact_topic", object.CompactTopic)
	// d.Set("partition_num", object.PartitionNum)
	// d.Set("remark", object.Remark)

	// tags, err := alikafkaService.ListTagResources(d.Id(), "topic")
	// if err == nil {
	// 	d.Set("tags", tagsToMap(tags))
	// }

	return nil
}

func resourceAlibabacloudStackAlikafkaTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	instanceId := parts[0]
	topic := parts[1]

	request := alikafka.CreateDeleteTopicRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = instanceId
	request.Topic = topic

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.DeleteTopic(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return errmsgs.WrapError(alikafkaService.WaitForAlikafkaTopic(d.Id(), Deleted, DefaultTimeoutMedium))
}
