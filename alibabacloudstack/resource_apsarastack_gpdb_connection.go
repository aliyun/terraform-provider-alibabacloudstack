package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackGpdbConnection() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"connection_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 31),
			},
			"port": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDBConnectionPort,
				Default:      "3306",
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackGpdbConnectionCreate, resourceAlibabacloudStackGpdbConnectionRead, resourceAlibabacloudStackGpdbConnectionUpdate, resourceAlibabacloudStackGpdbConnectionDelete)
	return resource
}

func resourceAlibabacloudStackGpdbConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	gpdbService := GpdbService{client}
	instanceId := d.Get("instance_id").(string)
	prefix := d.Get("connection_prefix").(string)
	if prefix == "" {
		prefix = fmt.Sprintf("%s-tf", instanceId)
	}
	request := gpdb.CreateAllocateInstancePublicConnectionRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = instanceId
	request.ConnectionStringPrefix = prefix
	request.Port = d.Get("port").(string)

	err := resource.Retry(8*time.Minute, func() *resource.RetryError {
		raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.AllocateInstancePublicConnection(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response, ok := raw.(*gpdb.AllocateInstancePublicConnectionResponse)
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_gpdb_connection", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s%s%s", instanceId, COLON_SEPARATED, request.ConnectionStringPrefix))
	// wait instance running after allocating
	stateConf := BuildStateConf([]string{"Creating", "NetAddressCreating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gpdbService.GpdbInstanceStateRefreshFunc(instanceId, []string{"Deleting"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func resourceAlibabacloudStackGpdbConnectionRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	client := meta.(*connectivity.AlibabacloudStackClient)
	gpdbService := GpdbService{client}
	object, err := gpdbService.DescribeGpdbConnection(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("instance_id", parts[0])
	d.Set("connection_prefix", parts[1])
	d.Set("port", object.Port)
	d.Set("connection_string", object.ConnectionString)
	d.Set("ip_address", object.IPAddress)

	return nil
}

func resourceAlibabacloudStackGpdbConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	if d.HasChange("port") {
		client := meta.(*connectivity.AlibabacloudStackClient)
		gpdbService := GpdbService{client}

		request := gpdb.CreateModifyDBInstanceConnectionStringRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = parts[0]
		object, err := gpdbService.DescribeGpdbConnection(d.Id())
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.CurrentConnectionString = object.ConnectionString
		request.ConnectionStringPrefix = parts[1]
		request.Port = d.Get("port").(string)

		if err := resource.Retry(8*time.Minute, func() *resource.RetryError {
			raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
				return gpdbClient.ModifyDBInstanceConnectionString(request)
			})
			if err != nil {
				if errmsgs.IsExpectedErrors(err, errmsgs.OperationDeniedDBStatus) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response, ok := raw.(*gpdb.ModifyDBInstanceConnectionStringResponse)
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			return nil
		}); err != nil {
			return err
		}

		// wait instance running after modifying
		stateConf := BuildStateConf([]string{"NET_MODIFYING"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, gpdbService.GpdbInstanceStateRefreshFunc(request.DBInstanceId, []string{"Deleting"}))

		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	return nil
}

func resourceAlibabacloudStackGpdbConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	request := gpdb.CreateReleaseInstancePublicConnectionRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = parts[0]

	gpdbService := GpdbService{client}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		object, err := gpdbService.DescribeGpdbConnection(d.Id())
		if err != nil {
			return resource.NonRetryableError(errmsgs.WrapError(err))
		}
		request.CurrentConnectionString = object.ConnectionString

		var raw interface{}
		raw, err = client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.ReleaseInstancePublicConnection(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response, ok := raw.(*gpdb.ReleaseInstancePublicConnectionResponse)
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound", "InvalidCurrentConnectionString.NotFound", "AtLeastOneNetTypeExists"}) {
			return nil
		}
		return err
	}
	stateConf := BuildStateConf([]string{"NetAddressDeleting"}, []string{"Running"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gpdbService.GpdbInstanceStateRefreshFunc(request.DBInstanceId, []string{"Deleting"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return errmsgs.WrapError(gpdbService.WaitForGpdbConnection(d.Id(), Deleted, DefaultTimeoutMedium))
}
