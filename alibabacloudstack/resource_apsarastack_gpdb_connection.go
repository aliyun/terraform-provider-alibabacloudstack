package alibabacloudstack

import (
	"fmt"
	"strconv"
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
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1000, 65534),
				Default:      3306,
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
	request.Port = strconv.Itoa(d.Get("port").(int))

	err := resource.Retry(8*time.Minute, func() *resource.RetryError {
		raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.AllocateInstancePublicConnection(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			response, ok := raw.(*gpdb.AllocateInstancePublicConnectionResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_gpdb_connection", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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
	if port, err := toInt(object.Port); err != nil {
		return err
	} else {
		d.Set("port", port)
	}
	d.Set("connection_string", object.ConnectionString)
	d.Set("ip_address", object.IPAddress)

	stateConf := BuildStateConf([]string{"NetAddressCreating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, gpdbService.GpdbInstanceStateRefreshFunc(parts[0], []string{"Deleting"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func resourceAlibabacloudStackGpdbConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	if d.IsNewResource() {
		return nil
	}

	if d.HasChanges("connection_prefix", "port") {

		client := meta.(*connectivity.AlibabacloudStackClient)
		gpdbService := GpdbService{client}

		object, err := gpdbService.DescribeGpdbConnection(d.Id())
		if err != nil {
			return errmsgs.WrapError(err)
		}
		reqQuery := map[string]interface{}{
			"DBInstanceId":            parts[0],
			"CurrentConnectionString": object.ConnectionString,
			"ConnectionStringPrefix":  d.Get("connection_prefix"),
			"Port":                    strconv.Itoa(d.Get("port").(int)),
			"DBInstanceNetType":       "public",
		}

		if err := resource.Retry(8*time.Minute, func() *resource.RetryError {
			_, err := client.DoTeaRequest("POST", "gpdb", "2016-05-03", "ModifyDBInstanceConnectionString", "", nil, reqQuery, nil)

			if err != nil {
				if errmsgs.IsExpectedErrors(err, errmsgs.OperationDeniedDBStatus) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		}); err != nil {
			return err
		}

		// wait instance running after modifying
		stateConf := BuildStateConf([]string{"NET_MODIFYING", "NetAddressCreating"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, gpdbService.GpdbInstanceStateRefreshFunc(parts[0], []string{"Deleting"}))

		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}

		d.SetId(fmt.Sprintf("%s%s%s", reqQuery["DBInstanceId"], COLON_SEPARATED, reqQuery["ConnectionStringPrefix"]))
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
			errmsg := ""
			response, ok := raw.(*gpdb.ReleaseInstancePublicConnectionResponse)
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
	stateConf := BuildStateConf([]string{"NetAddressDeleting"}, []string{"Running"}, d.Timeout(schema.TimeoutDelete), 10*time.Second, gpdbService.GpdbInstanceStateRefreshFunc(request.DBInstanceId, []string{"Deleting"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return errmsgs.WrapError(gpdbService.WaitForGpdbConnection(d.Id(), Deleted, DefaultTimeoutMedium))
}
