package alibabacloudstack

import (
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDiskCreate,
		Read:   resourceAlibabacloudStackDiskRead,
		Update: resourceAlibabacloudStackDiskUpdate,
		Delete: resourceAlibabacloudStackDiskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'availability_zone' is deprecated and will be removed in a future release. Please use new field 'zone_id' instead.",
				ConflictsWith: []string{"zone_id"},
			},
			"zone_id": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"availability_zone"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(2, 128),
				Deprecated:    "Field 'name' is deprecated and will be removed in a future release. Please use new field 'disk_name' instead.",
				ConflictsWith: []string{"disk_name"},
			},
			"disk_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"name"},
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},

			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud", "cloud_efficiency", "cloud_ssd", "cloud_pperf", "cloud_sperf"}, false),
				Default:      DiskCloudEfficiency,
			},

			"size": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"snapshot_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"encrypted"},
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"encrypted": {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"snapshot_id"},
			},
			"encrypt_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"sm4-128", "aes-256"}, false),
			},

			"delete_auto_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"delete_with_instance": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"enable_auto_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_snapshot_policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_automated_snapshot_policy": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceAlibabacloudStackDiskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	zoneId := connectivity.GetResourceData(d, "zone_id", "availability_zone").(string)
	if err := errmsgs.CheckEmpty(zoneId, schema.TypeString, "zone_id", "availability_zone"); err != nil {
		return errmsgs.WrapError(err)
	}
	availabilityZone, err := ecsService.DescribeZone(zoneId)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	request := ecs.CreateCreateDiskRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ZoneId = availabilityZone.ZoneId

	if v, ok := d.GetOk("category"); ok && v.(string) != "" {
		category := DiskCategory(v.(string))
		if err := ecsService.DiskAvailable(availabilityZone, category); err != nil {
			return errmsgs.WrapError(err)
		}
		request.DiskCategory = v.(string)
	}

	if v, ok := d.GetOk("size"); ok {
		request.Size = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("snapshot_id"); ok && v.(string) != "" {
		request.SnapshotId = v.(string)
	}

	if v, ok := connectivity.GetResourceDataOk(d, "disk_name", "name"); ok && v.(string) != "" {
		request.DiskName = v.(string)
	}
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.Description = v.(string)
	}

	if v, ok := d.GetOk("encrypted"); ok {
		request.Encrypted = requests.NewBoolean(v.(bool))
		if v.(bool) == true {
			if j, ok4 := d.GetOk("kms_key_id"); ok4 {
				request.KMSKeyId = j.(string)
			}
			if request.KMSKeyId == "" {
				return errmsgs.WrapError(errors.New("KmsKeyId can not be empty if encrypted is set to \"true\""))
			}
			request.EncryptAlgorithm = d.Get("encrypt_algorithm").(string)
		}
	}
	if v, ok := d.GetOk("tags"); ok && len(v.(map[string]interface{})) > 0 {
		tags := make([]ecs.CreateDiskTag, len(v.(map[string]interface{})))
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, ecs.CreateDiskTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &tags
	}
	request.ClientToken = buildClientToken(request.GetActionName())
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateDisk(request)
	})
	response, ok := raw.(*ecs.CreateDiskResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_disk", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	d.SetId(response.DiskId)
	if err := ecsService.WaitForDisk(d.Id(), Available, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}

	return resourceAlibabacloudStackDiskUpdate(d, meta)
}

func resourceAlibabacloudStackDiskRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeDisk(d.Id())

	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	connectivity.SetResourceData(d, object.ZoneId, "zone_id", "availability_zone")
	d.Set("category", object.Category)
	d.Set("size", object.Size)
	d.Set("status", object.Status)
	connectivity.SetResourceData(d, object.DiskName, "disk_name", "name")
	d.Set("description", object.Description)
	d.Set("snapshot_id", object.SourceSnapshotId)
	d.Set("encrypted", object.Encrypted)
	d.Set("kms_key_id", object.KMSKeyId)
	d.Set("delete_auto_snapshot", object.DeleteAutoSnapshot)
	d.Set("delete_with_instance", object.DeleteWithInstance)
	d.Set("enable_auto_snapshot", object.EnableAutoSnapshot)
	d.Set("enable_automated_snapshot_policy", object.EnableAutomatedSnapshotPolicy)
	d.Set("auto_snapshot_policy_id", object.AutoSnapshotPolicyId)
	d.Set("tags", ecsService.tagsToMap(object.Tags.Tag))

	return nil
}

func resourceAlibabacloudStackDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	d.Partial(true)

	update := false
	request := ecs.CreateModifyDiskAttributeRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DiskId = d.Id()

	if !d.IsNewResource() && d.HasChange("disk_name") {
		request.DiskName = connectivity.GetResourceData(d, "disk_name", "name").(string)
		update = true
		//d.SetPartial("disk_name")
	}

	if !d.IsNewResource() && d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
		//d.SetPartial("description")
	}

	if d.IsNewResource() || d.HasChange("delete_auto_snapshot") {
		v := d.Get("delete_auto_snapshot")
		request.DeleteAutoSnapshot = requests.NewBoolean(v.(bool))
		update = true
		//d.SetPartial("delete_auto_snapshot")
	}

	if d.IsNewResource() || d.HasChange("delete_with_instance") {
		v := d.Get("delete_with_instance")
		request.DeleteWithInstance = requests.NewBoolean(v.(bool))
		update = true
		//d.SetPartial("delete_with_instance")
	}

	if d.IsNewResource() || d.HasChange("enable_auto_snapshot") {
		v := d.Get("enable_auto_snapshot")
		request.EnableAutoSnapshot = requests.NewBoolean(v.(bool))
		update = true
		//d.SetPartial("enable_auto_snapshot")
	}

	if update {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyDiskAttribute(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*ecs.ModifyDiskAttributeResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlibabacloudStackDiskRead(d, meta)
	}

	err := setTags(client, TagResourceDisk, d)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	//d.SetPartial("tags")

	if d.HasChange("size") {
		size := d.Get("size").(int)
		request := ecs.CreateResizeDiskRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DiskId = d.Id()
		request.NewSize = requests.NewInteger(size)
		request.Type = string(DiskResizeTypeOnline)
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ResizeDisk(request)
		})
		if errmsgs.IsExpectedErrors(err, errmsgs.DiskNotSupportOnlineChangeErrors) {
			request.Type = string(DiskResizeTypeOffline)
			raw, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ResizeDisk(request)
			})
		}
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*ecs.ResizeDiskResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("size")
	}

	d.Partial(false)
	return resourceAlibabacloudStackDiskRead(d, meta)
}

func resourceAlibabacloudStackDiskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	request := ecs.CreateDeleteDiskRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DiskId = d.Id()

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteDisk(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.DiskInvalidOperation) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*ecs.DeleteDiskResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return errmsgs.WrapError(ecsService.WaitForDisk(d.Id(), Deleted, DefaultTimeout))
}
