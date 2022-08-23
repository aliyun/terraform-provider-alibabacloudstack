package apsarastack

import (
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackImageCopy() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackImageCopyCreate,
		Read:   resourceApsaraStackImageCopyRead,
		Update: resourceApsaraStackImageCopyUpdate,
		Delete: resourceApsaraStackImageCopyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Attribute 'name' has been deprecated from version 1.69.0. Use `image_name` instead.",
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
func resourceApsaraStackImageCopyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}

	request := ecs.CreateCopyImageRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.ImageId = d.Get("source_image_id").(string)
	request.DestinationRegionId = d.Get("destination_region_id").(string)
	request.DestinationImageName = d.Get("image_name").(string)
	request.DestinationDescription = d.Get("description").(string)
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {

		return ecsClient.CopyImage(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsrastack_image_copy", request.GetActionName(), ApsaraStackGoClientFailure)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ := raw.(*ecs.CopyImageResponse)
	d.SetId(response.ImageId)
	log.Printf("[DEBUG] state %#v", d.Id())
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, ecsService.ImageStateRefreshFuncforcopy(d.Id(), d.Get("destination_region_id").(string), []string{"CreateFailed", "UnAvailable"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceApsaraStackImageCopyRead(d, meta)
}
func resourceApsaraStackImageCopyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}
	err := ecsService.updateImage(d)
	if err != nil {
		return WrapError(err)
	}
	return resourceApsaraStackImageRead(d, meta)
}
func resourceApsaraStackImageCopyRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)

	ecsService := EcsService{client}
	object, err := ecsService.DescribeImage(d.Id(), d.Get("destination_region_id").(string))
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.ImageName)
	d.Set("image_name", object.ImageName)
	d.Set("description", object.Description)

	return WrapError(err)
}

func resourceApsaraStackImageCopyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}
	err := ecsService.deleteImageforDest(d, d.Get("destination_region_id").(string))
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), ApsaraStackSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{"Available", "CreateFailed"}, []string{"Deprecated", "UnAvailable"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, ecsService.ImageStateRefreshFuncforcopy(d.Id(), d.Get("destination_region_id").(string), []string{"CreateFailed", "UnAvailable"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceApsaraStackImageCopyRead(d, meta)
}
