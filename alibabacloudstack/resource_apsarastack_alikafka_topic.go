package alibabacloudstack

import (
	"errors"
	"strconv"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
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
	regionId := client.RegionId
	topic := d.Get("topic").(string)

	request := alikafka.CreateCreateTopicRequest()
	request.InstanceId = instanceId
	request.RegionId = regionId
	request.Topic = topic

	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "alikafka",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "CreateTopic",
		"Version":         "2019-09-16",
	}
	if v, ok := d.GetOk("local_topic"); ok {
		request.LocalTopic = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("compact_topic"); ok {
		request.CompactTopic = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("partition_num"); ok {
		request.PartitionNum = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("remark"); ok {
		request.Remark = v.(string)
	}
	request.Domain = client.Config.AlikafkaOpenAPIEndpoint
	request.QueryParams["Product"] = "alikafka"
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.CreateTopic(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_alikafka_topic", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	d.SetId(instanceId + ":" + topic)

	return resourceAlibabacloudStackAlikafkaTopicUpdate(d, meta)
}

func resourceAlibabacloudStackAlikafkaTopicUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}
	d.Partial(true)
	if err := alikafkaService.setInstanceTags(d, TagResourceTopic); err != nil {
		return WrapError(err)
	}
	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlibabacloudStackAlikafkaTopicRead(d, meta)
	}

	instanceId := d.Get("instance_id").(string)
	if d.HasChange("remark") {
		remark := d.Get("remark").(string)
		topic := d.Get("topic").(string)
		modifyRemarkRequest := alikafka.CreateModifyTopicRemarkRequest()
		modifyRemarkRequest.InstanceId = instanceId
		modifyRemarkRequest.RegionId = client.RegionId
		modifyRemarkRequest.Topic = topic
		modifyRemarkRequest.Remark = remark

		

		modifyRemarkRequest.QueryParams = map[string]string{
			"AccessKeySecret": client.SecretKey,
			"AccessKeyId":     client.AccessKey,
			"Product":         "alikafka",
			"RegionId":        client.RegionId,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Action":          "ModifyTopicRemark",
			"Version":         "2019-09-16",
		}

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
				return alikafkaClient.ModifyTopicRemark(modifyRemarkRequest)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser}) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(modifyRemarkRequest.GetActionName(), raw, modifyRemarkRequest.RpcRequest, modifyRemarkRequest)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyRemarkRequest.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		//d.SetPartial("remark")
	}

	if d.HasChange("partition_num") {
		o, n := d.GetChange("partition_num")
		oldPartitionNum := o.(int)
		newPartitionNum := n.(int)

		if newPartitionNum < oldPartitionNum {
			return WrapError(errors.New("partition_num only support adjust to a greater value."))
		} else {
			topic := d.Get("topic").(string)
			modifyPartitionReq := alikafka.CreateModifyPartitionNumRequest()
			modifyPartitionReq.InstanceId = instanceId
			modifyPartitionReq.RegionId = client.RegionId
			modifyPartitionReq.Topic = topic
			modifyPartitionReq.AddPartitionNum = requests.NewInteger(newPartitionNum - oldPartitionNum)

			


			modifyPartitionReq.QueryParams = map[string]string{
				"AccessKeySecret": client.SecretKey,
				"AccessKeyId":     client.AccessKey,
				"Product":         "alikafka",
				"RegionId":        client.RegionId,
				"Department":      client.Department,
				"ResourceGroup":   client.ResourceGroup,
				"Action":          "ModifyPartitionNum",
				"Version":         "2019-09-16",
			}

			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
					return alikafkaClient.ModifyPartitionNum(modifyPartitionReq)
				})
				if err != nil {
					if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
						time.Sleep(10 * time.Second)
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(modifyPartitionReq.GetActionName(), raw, modifyPartitionReq.RpcRequest, modifyPartitionReq)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyPartitionReq.GetActionName(), AlibabacloudStackSdkGoERROR)
			}
			//d.SetPartial("partition_num")
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
		// Handle exceptions
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.InstanceId)
	d.Set("topic", object.Topic)
	d.Set("local_topic", object.LocalTopic)
	//d.Set("compact_topic", object.CompactTopic)
	//d.Set("partition_num", object.PartitionNum)
	//d.Set("remark", object.Remark)

	//tags, err := alikafkaService.DescribeTags(d.Id(), nil, TagResourceTopic)
	//if err != nil {
	//	return WrapError(err)
	//}
	//d.Set("tags", alikafkaService.tagsToMap(tags))

	return nil
}

func resourceAlibabacloudStackAlikafkaTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	instanceId := parts[0]
	topic := parts[1]

	request := alikafka.CreateDeleteTopicRequest()
	request.Topic = topic
	request.InstanceId = instanceId
	request.RegionId = client.RegionId
	request.Domain = client.Config.AlikafkaOpenAPIEndpoint
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "alikafka",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "DeleteTopic",
		"Version":         "2019-09-16",
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.DeleteTopic(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	return WrapError(alikafkaService.WaitForAlikafkaTopic(d.Id(), Deleted, DefaultTimeoutMedium))
}
