package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAlikafkaSaslAcl() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"acl_resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Group", "Topic"}, false),
			},
			"acl_resource_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"acl_resource_pattern_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"LITERAL", "PREFIXED"}, false),
			},
			"acl_operation_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Read", "Write"}, false),
			},
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAlikafkaSaslAclCreate, resourceAlibabacloudStackAlikafkaSaslAclRead, nil, resourceAlibabacloudStackAlikafkaSaslAclDelete)
	return resource
}

func resourceAlibabacloudStackAlikafkaSaslAclCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	instanceId := d.Get("instance_id").(string)
	regionId := client.RegionId
	username := d.Get("username").(string)
	aclResourceType := d.Get("acl_resource_type").(string)
	aclResourceName := d.Get("acl_resource_name").(string)
	aclResourcePatternType := d.Get("acl_resource_pattern_type").(string)
	aclOperationType := d.Get("acl_operation_type").(string)

	request := alikafka.CreateCreateAclRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = instanceId
	request.RegionId = regionId
	request.Username = username
	request.AclResourceType = aclResourceType
	request.AclResourceName = aclResourceName
	request.AclResourcePatternType = aclResourcePatternType
	request.AclOperationType = aclOperationType
	request.Domain = client.Config.Endpoints[connectivity.ALIKAFKACode]
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.CreateAcl(request)
		})
		bresponse, ok := raw.(*alikafka.CreateAclResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_alikafka_sasl_acl", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return err
	}

	// Server may have cache, sleep a while.
	time.Sleep(60 * time.Second)
	d.SetId(fmt.Sprintf("%s:%s:%s:%s:%s:%s", instanceId, username, aclResourceType, aclResourceName, aclResourcePatternType, aclOperationType))
	return nil
}

func resourceAlibabacloudStackAlikafkaSaslAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	parts, err := ParseResourceId(d.Id(), 6)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	object, err := alikafkaService.DescribeAlikafkaSaslAcl(d.Id())
	if err != nil {
		// Handle exceptions
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("instance_id", parts[0])
	d.Set("username", object.Username)
	d.Set("acl_resource_type", object.AclResourceType)
	d.Set("acl_resource_name", object.AclResourceName)
	d.Set("acl_resource_pattern_type", object.AclResourcePatternType)
	d.Set("acl_operation_type", object.AclOperationType)
	d.Set("host", object.Host)

	return nil
}

func resourceAlibabacloudStackAlikafkaSaslAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	parts, err := ParseResourceId(d.Id(), 6)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	instanceId := parts[0]
	username := parts[1]
	aclResourceType := parts[2]
	aclResourceName := parts[3]
	aclResourcePatternType := parts[4]
	aclOperationType := parts[5]

	request := alikafka.CreateDeleteAclRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = instanceId
	request.RegionId = client.RegionId
	request.Username = username
	request.AclResourceType = aclResourceType
	request.AclResourceName = aclResourceName
	request.AclResourcePatternType = aclResourcePatternType
	request.AclOperationType = aclOperationType
	request.Domain = client.Config.Endpoints[connectivity.ALIKAFKACode]
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.DeleteAcl(request)
		})
		bresponse, ok := raw.(*alikafka.DeleteAclResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return err
	}

	// Server may have cache, sleep a while.
	time.Sleep(60 * time.Second)
	return errmsgs.WrapError(alikafkaService.WaitForAlikafkaSaslAcl(d.Id(), Deleted, DefaultTimeoutMedium))
}