package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDBDatabase() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"data_base_instance_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ConflictsWith: []string{"instance_id"},
			},
			"instance_id": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'instance_id' is deprecated and will be removed in a future release. Please use new field 'data_base_instance_id' instead.",
				ConflictsWith: []string{"data_base_instance_id"},
			},

			"data_base_name": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' is deprecated and will be removed in a future release. Please use new field 'data_base_name' instead.",
				ConflictsWith: []string{"data_base_name"},
			},

			"character_set_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"character_set"},
			},
			"character_set": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'character_set' is deprecated and will be removed in a future release. Please use new field 'character_set_name' instead.",
				ConflictsWith: []string{"character_set_name"},
			},

			"data_base_description": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"description"},
			},
			"description": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'description' is deprecated and will be removed in a future release. Please use new field 'data_base_description' instead.",
				ConflictsWith: []string{"data_base_description"},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackDBDatabaseCreate, resourceAlibabacloudStackDBDatabaseRead, resourceAlibabacloudStackDBDatabaseUpdate, resourceAlibabacloudStackDBDatabaseDelete)
	return resource
}

func resourceAlibabacloudStackDBDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := rds.CreateCreateDatabaseRequest()
	client.InitRpcRequest(*request.RpcRequest)

	request.DBInstanceId = connectivity.GetResourceData(d, "data_base_instance_id", "instance_id").(string)
	if err := errmsgs.CheckEmpty(request.DBInstanceId, schema.TypeString, "data_base_instance_id", "instance_id"); err != nil {
		return errmsgs.WrapError(err)
	}
	request.DBName = connectivity.GetResourceData(d, "data_base_name", "name").(string)
	if err := errmsgs.CheckEmpty(request.DBName, schema.TypeString, "data_base_name", "name"); err != nil {
		return errmsgs.WrapError(err)
	}
	request.CharacterSetName = connectivity.GetResourceData(d, "character_set_name", "character_set").(string)
	if err := errmsgs.CheckEmpty(request.CharacterSetName, schema.TypeString, "character_set_name", "character_set"); err != nil {
		return errmsgs.WrapError(err)
	}
	if v, ok := connectivity.GetResourceDataOk(d, "data_base_description", "description"); ok && v.(string) != "" {
		request.DBDescription = v.(string)
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.CreateDatabase(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*rds.CreateDatabaseResponse)
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

	d.SetId(fmt.Sprintf("%s%s%s", request.DBInstanceId, COLON_SEPARATED, request.DBName))

	return nil
}

func resourceAlibabacloudStackDBDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	rsdService := RdsService{client}
	object, err := rsdService.DescribeDBDatabase(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	connectivity.SetResourceData(d, object.DBInstanceId, "data_base_instance_id", "instance_id")
	connectivity.SetResourceData(d, object.DBName, "data_base_name", "name")
	connectivity.SetResourceData(d, object.CharacterSetName, "character_set_name", "character_set")
	connectivity.SetResourceData(d, object.DBDescription, "data_base_description", "description")

	return nil
}

func resourceAlibabacloudStackDBDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if (d.HasChanges("data_base_description", "description")) && !d.IsNewResource() {
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request := rds.CreateModifyDBDescriptionRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = parts[0]
		request.DBName = parts[1]

		request.DBDescription = connectivity.GetResourceData(d, "data_base_description", "description").(string)

		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ModifyDBDescription(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*rds.ModifyDBDescriptionResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return nil
}

func resourceAlibabacloudStackDBDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request := rds.CreateDeleteDatabaseRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = parts[0]
	request.DBName = parts[1]
	// wait instance status is running before deleting database
	if err := rdsService.WaitForDBInstance(parts[0], Running, 1800); err != nil {
		return errmsgs.WrapError(err)
	}
	raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DeleteDatabase(request)
	})
	if err != nil {
		if errmsgs.NotFoundError(err) || errmsgs.IsExpectedErrors(err, []string{"InvalidDBName.NotFound"}) {
			return nil
		}
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*rds.DeleteDatabaseResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return errmsgs.WrapError(rdsService.WaitForDBDatabase(d.Id(), Deleted, DefaultTimeoutMedium))
}