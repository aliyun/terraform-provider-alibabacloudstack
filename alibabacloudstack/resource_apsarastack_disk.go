package alibabacloudstack

import (
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
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

			"tags": tagsSchema(),
			"storage_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"storage_set_partition_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 2000),
			},
		},
	}
}

func resourceAlibabacloudStackDiskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	availabilityZone, err := ecsService.DescribeZone(d.Get("availability_zone").(string))
	if err != nil {
		return WrapError(err)
	}

	request := ecs.CreateCreateDiskRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.ZoneId = availabilityZone.ZoneId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	if v, ok := d.GetOk("category"); ok && v.(string) != "" {
		category := DiskCategory(v.(string))
		if err := ecsService.DiskAvailable(availabilityZone, category); err != nil {
			return WrapError(err)
		}
		request.DiskCategory = v.(string)
	}

	request.Size = requests.NewInteger(d.Get("size").(int))

	if v, ok := d.GetOk("snapshot_id"); ok && v.(string) != "" {
		request.SnapshotId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
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
				return WrapError(errors.New("KmsKeyId can not be empty if encrypted is set to \"true\""))
			}
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
	i := d.Get("storage_set_partition_number").(int)
	if v := d.Get("storage_set_id").(string); v != "" {
		request.StorageSetId = v
		if i >= 1 {
			request.StorageSetPartitionNumber = requests.NewInteger(d.Get("storage_set_partition_number").(int))
		} else {
			 return WrapError(errors.New("cant empty storage_set_partition_number when you set storage_set_id and >=2"))
		}
	}
	request.ClientToken = buildClientToken(request.GetActionName())
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateDisk(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_disk", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.CreateDiskResponse)
	d.SetId(response.DiskId)
	if err := ecsService.WaitForDisk(d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAlibabacloudStackDiskUpdate(d, meta)
}

func resourceAlibabacloudStackDiskRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeDisk(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("availability_zone", object.ZoneId)
	d.Set("category", object.Category)
	d.Set("size", object.Size)
	d.Set("status", object.Status)
	d.Set("name", object.DiskName)
	d.Set("description", object.Description)
	d.Set("snapshot_id", object.SourceSnapshotId)
	d.Set("encrypted", object.Encrypted)
	d.Set("kms_key_id", object.KMSKeyId)
	d.Set("delete_auto_snapshot", object.DeleteAutoSnapshot)
	d.Set("delete_with_instance", object.DeleteWithInstance)
	d.Set("enable_auto_snapshot", object.EnableAutoSnapshot)
	d.Set("tags", ecsService.tagsToMap(object.Tags.Tag))
	d.Set("storage_set_id", object.StorageSetId)

	return nil
}

func resourceAlibabacloudStackDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	d.Partial(true)

	update := false
	request := ecs.CreateModifyDiskAttributeRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.DiskId = d.Id()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	if !d.IsNewResource() && d.HasChange("name") {
		request.DiskName = d.Get("name").(string)
		update = true
		//d.SetPartial("name")
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlibabacloudStackDiskRead(d, meta)
	}

	err := setTags(client, TagResourceDisk, d)
	if err != nil {
		return WrapError(err)
	}
	//d.SetPartial("tags")

	if d.HasChange("size") {
		size := d.Get("size").(int)
		request := ecs.CreateResizeDiskRequest()
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.DiskId = d.Id()
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.NewSize = requests.NewInteger(size)
		request.Type = string(DiskResizeTypeOnline)
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ResizeDisk(request)
		})
		if IsExpectedErrors(err, DiskNotSupportOnlineChangeErrors) {
			request.Type = string(DiskResizeTypeOffline)
			raw, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ResizeDisk(request)
			})
		}
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
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
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.RegionId = client.RegionId
	request.DiskId = d.Id()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteDisk(request)
		})
		if err != nil {
			if IsExpectedErrors(err, DiskInvalidOperation) {
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
	return WrapError(ecsService.WaitForDisk(d.Id(), Deleted, DefaultTimeout))
}
