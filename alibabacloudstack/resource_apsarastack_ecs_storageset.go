package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/sdk_patch/datahub_patch"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEcsEbsStorageSets() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"storage_set_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"maxpartition_number": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"storage_set_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEcsEbsStorageSetsCreate, resourceAlibabacloudStackEcsEbsStorageSetsRead, nil, resourceAlibabacloudStackEcsEbsStorageSetsDelete)
	return resource
}

func resourceAlibabacloudStackEcsEbsStorageSetsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "CreateStorageSet"
	response := &datahub_patch.EcsStorageSetsCreate{}

	StorageSetName := d.Get("storage_set_name").(string)
	MaxPartitionNumber := d.Get("maxpartition_number").(string)
	ZoneId := d.Get("zone_id").(string)

	request := client.NewCommonRequest("POST", "Ecs", "2014-05-26", action, "")
	request.QueryParams["StorageSetName"] = StorageSetName
	request.QueryParams["MaxPartitionNumber"] = MaxPartitionNumber
	request.QueryParams["ZoneId"] = ZoneId

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
			return EcsClient.ProcessCommonRequest(request)
		})
		addDebug(action, raw, request)
		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			if errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ecs_command", action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		err = json.Unmarshal(bresponse.GetHttpContentBytes(), response)
		d.SetId(fmt.Sprint(response.StorageSetId))
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceAlibabacloudStackEcsEbsStorageSetsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsEbsStorageSet(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_ecs_command ecsService.DescribeEcsCommand Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("storage_set_name", object.StorageSets.StorageSet[0].StorageSetName)
	d.Set("zone_id", object.StorageSets.StorageSet[0].ZoneId)
	d.Set("maxpartition_number", strconv.Itoa(object.StorageSets.StorageSet[0].StorageSetPartitionNumber))
	d.Set("storage_set_id", object.StorageSets.StorageSet[0].StorageSetId)
	return nil
}

func resourceAlibabacloudStackEcsEbsStorageSetsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteStorageSet"

	request := client.NewCommonRequest("POST", "Ecs", "2014-05-26", action, "")
	request.QueryParams["StorageSetId"] = d.Id()

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(EcsClient *ecs.Client) (interface{}, error) {
			return EcsClient.ProcessCommonRequest(request)
		})
		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			if errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(action, raw, request)
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidCmdId.NotFound", "InvalidRegionId.NotFound", "Operation.Forbidden"}) {
			return nil
		}
		return err
	}
	return nil
}
