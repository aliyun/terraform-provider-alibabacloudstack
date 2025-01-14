package alibabacloudstack

import (
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackImageCopy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackImageCopyCreate,
		Read:   resourceAlibabacloudStackImageCopyRead,
		Update: resourceAlibabacloudStackImageCopyUpdate,
		Delete: resourceAlibabacloudStackImageCopyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"source_image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_region_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' is deprecated and will be removed in a future release. Please use new field 'image_name' instead.",
				ConflictsWith: []string{"image_name"},
			},
			"image_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlibabacloudStackImageCopyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	request := ecs.CreateCopyImageRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ImageId = d.Get("source_image_id").(string)
	request.DestinationRegionId = d.Get("destination_region_id").(string)
	request.DestinationImageName = connectivity.GetResourceData(d, "image_name", "name").(string)
	request.DestinationDescription = d.Get("description").(string)
	request.ResourceGroupId = client.Config.ResourceGroupId
	if v, ok := d.GetOk("kms_key_id"); ok && v != "" {
		request.KMSKeyId = v.(string)
	}
	if v, ok := d.GetOk("encrypted"); ok {
		request.Encrypted = requests.NewBoolean(v.(bool))
	}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CopyImage(request)
	})
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*ecs.CopyImageResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "apsrastack_image_copy", request.GetActionName(), errmsgs.AlibabacloudStackGoClientFailure, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ := raw.(*ecs.CopyImageResponse)
	d.SetId(response.ImageId)
	log.Printf("[DEBUG] state %#v", d.Id())
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 10*time.Minute, ecsService.ImageStateRefreshFuncforcopy(d.Id(), d.Get("destination_region_id").(string), []string{"CreateFailed", "UnAvailable"}))
	stateConf.NotFoundChecks = 1000
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return resourceAlibabacloudStackImageCopyRead(d, meta)
}

func resourceAlibabacloudStackImageCopyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	err := ecsService.updateImage(d)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	return resourceAlibabacloudStackImageRead(d, meta)
}

func resourceAlibabacloudStackImageCopyRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)

	ecsService := EcsService{client}
	object, err := ecsService.DescribeImage(d.Id(), d.Get("destination_region_id").(string))
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	connectivity.SetResourceData(d, object.ImageName, "image_name", "name")
	d.Set("description", object.Description)

	return errmsgs.WrapError(err)
}

func resourceAlibabacloudStackImageCopyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		request := ecs.CreateDeleteImageRequest()
		request.ImageId = d.Id()
		return ecsClient.DeleteImage(request)
	})
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*ecs.DeleteImageResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	stateConf := BuildStateConf([]string{"Available", "CreateFailed"}, []string{"Deprecated", "UnAvailable"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, ecsService.ImageStateRefreshFuncforcopy(d.Id(), d.Get("destination_region_id").(string), []string{"CreateFailed", "UnAvailable"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return resourceAlibabacloudStackImageCopyRead(d, meta)
}
