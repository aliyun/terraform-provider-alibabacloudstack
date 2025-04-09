package alibabacloudstack

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackOtsInstance() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 16),
			},

			"accessed_by": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      AnyNetwork,
				ValidateFunc: validation.StringInSlice([]string{string(AnyNetwork), string(VpcOnly), string(VpcOrConsole)}, false),
			},

			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      OtsHighPerformance,
				ValidateFunc: validation.StringInSlice([]string{string(OtsCapacity), string(OtsHighPerformance)}, false),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Id() != ""
				},
			},
			"tags": tagsSchema(),
		},
	}
	setResourceFunc(resource, resourceAliyunOtsInstanceCreate, resourceAliyunOtsInstanceRead, resourceAliyunOtsInstanceUpdate, resourceAliyunOtsInstanceDelete)
	return resource
}

func resourceAliyunOtsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	otsService := OtsService{client}

	instanceType := d.Get("instance_type").(string)
	request := ots.CreateInsertInstanceRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ClusterType = convertInstanceType(OtsInstanceType(instanceType))

	types, err := otsService.DescribeOtsInstanceTypes()
	if err != nil {
		return errmsgs.WrapError(err)
	}
	valid := false
	for _, t := range types {
		if request.ClusterType == t {
			valid = true
			break
		}
	}
	if !valid {
		return errmsgs.WrapError(errmsgs.Error("The instance type %s is not available in the region %s.", instanceType, client.RegionId))
	}

	request.InstanceName = d.Get("name").(string)
	request.Description = d.Get("description").(string)
	request.Network = convertInstanceAccessedBy(InstanceAccessedByType(d.Get("accessed_by").(string)))

	raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.InsertInstance(request)
	})
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*ots.InsertInstanceResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.SetId(request.InstanceName)
	if err := otsService.WaitForOtsInstance(request.InstanceName, Running, DefaultTimeout/3); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func resourceAliyunOtsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	otsService := OtsService{client}
	object, err := otsService.DescribeOtsInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("name", object.InstanceName)
	d.Set("accessed_by", convertInstanceAccessedByRevert(object.Network))
	d.Set("instance_type", convertInstanceTypeRevert(object.ClusterType))
	d.Set("description", object.Description)
	d.Set("tags", otsTagsToMapFun(object.TagInfos))
	return nil
}

func resourceAliyunOtsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	otsService := OtsService{client}

	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("accessed_by") {
		request := ots.CreateUpdateInstanceRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.InstanceName = d.Id()
		request.Network = convertInstanceAccessedBy(InstanceAccessedByType(d.Get("accessed_by").(string)))

		raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.UpdateInstance(request)
		})
		if err != nil {
			errmsg := ""
			if response, ok := raw.(*ots.UpdateInstanceResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := diffTags(tagsFromMap(o), tagsFromMap(n))

		if len(remove) > 0 {
			request := ots.CreateDeleteTagsRequest()
			client.InitRpcRequest(*request.RpcRequest)
			request.InstanceName = d.Id()
			var tags []ots.DeleteTagsTagInfo
			for _, t := range remove {
				tags = append(tags, ots.DeleteTagsTagInfo{
					TagKey:   t.Key,
					TagValue: t.Value,
				})
			}
			request.TagInfo = &tags

			raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
				return otsClient.DeleteTags(request)
			})
			if err != nil {
				errmsg := ""
				if response, ok := raw.(*ots.DeleteTagsResponse); ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		if len(create) > 0 {
			request := ots.CreateInsertTagsRequest()
			client.InitRpcRequest(*request.RpcRequest)
			request.InstanceName = d.Id()
			var tags []ots.InsertTagsTagInfo
			for _, t := range create {
				tags = append(tags, ots.InsertTagsTagInfo{
					TagKey:   t.Key,
					TagValue: t.Value,
				})
			}
			request.TagInfo = &tags

			raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
				return otsClient.InsertTags(request)
			})
			if err != nil {
				errmsg := ""
				if response, ok := raw.(*ots.InsertTagsResponse); ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}
	}

	if err := otsService.WaitForOtsInstance(d.Id(), Running, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}
	d.Partial(false)
	return nil
}

func resourceAliyunOtsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	otsService := OtsService{client}
	request := ots.CreateDeleteInstanceRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceName = d.Id()

	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.DeleteInstance(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"AuthFailed", "InvalidStatus", "ValidationFailed"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if response, ok := raw.(*ots.DeleteInstanceResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errmsgs.NotFoundError(err) {
			return nil
		}
		return errmsgs.WrapError(err)
	}
	return errmsgs.WrapError(otsService.WaitForOtsInstance(d.Id(), Deleted, DefaultLongTimeout))
}