package apsarastack

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"

	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackGpdbConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackGpdbConnectionCreate,
		Read:   resourceApsaraStackGpdbConnectionRead,
		Update: resourceApsaraStackGpdbConnectionUpdate,
		Delete: resourceApsaraStackGpdbConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
}

func resourceApsaraStackGpdbConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	gpdbService := GpdbService{client}
	instanceId := d.Get("instance_id").(string)
	prefix := d.Get("connection_prefix").(string)
	if prefix == "" {
		prefix = fmt.Sprintf("%s-tf", instanceId)
	}
	request := gpdb.CreateAllocateInstancePublicConnectionRequest()
	request.Headers["x-ascm-product-name"] = "Gpdb"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "gpdb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.DBInstanceId = instanceId
	request.ConnectionStringPrefix = prefix
	request.Port = d.Get("port").(string)
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	err := resource.Retry(8*time.Minute, func() *resource.RetryError {
		raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.AllocateInstancePublicConnection(request)
		})
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_gpdb_connection", request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", instanceId, COLON_SEPARATED, request.ConnectionStringPrefix))
	// wait instance running after allocating
	stateConf := BuildStateConf([]string{"Creating", "NetAddressCreating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gpdbService.GpdbInstanceStateRefreshFunc(instanceId, []string{"Deleting"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceApsaraStackGpdbConnectionRead(d, meta)
}

func resourceApsaraStackGpdbConnectionRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	client := meta.(*connectivity.ApsaraStackClient)
	gpdbService := GpdbService{client}
	object, err := gpdbService.DescribeGpdbConnection(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", parts[0])
	d.Set("connection_prefix", parts[1])
	d.Set("port", object.Port)
	d.Set("connection_string", object.ConnectionString)
	d.Set("ip_address", object.IPAddress)

	return nil
}

func resourceApsaraStackGpdbConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	if d.HasChange("port") {
		client := meta.(*connectivity.ApsaraStackClient)
		gpdbService := GpdbService{client}

		request := gpdb.CreateModifyDBInstanceConnectionStringRequest()
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "gpdb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.DBInstanceId = parts[0]
		object, err := gpdbService.DescribeGpdbConnection(d.Id())
		if err != nil {
			return WrapError(err)
		}
		request.CurrentConnectionString = object.ConnectionString
		request.ConnectionStringPrefix = parts[1]
		request.Port = d.Get("port").(string)

		if err := resource.Retry(8*time.Minute, func() *resource.RetryError {
			raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
				return gpdbClient.ModifyDBInstanceConnectionString(request)
			})
			if err != nil {
				if IsExpectedErrors(err, OperationDeniedDBStatus) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}

		// wait instance running after modifying
		stateConf := BuildStateConf([]string{"NET_MODIFYING"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, gpdbService.GpdbInstanceStateRefreshFunc(request.DBInstanceId, []string{"Deleting"}))

		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceApsaraStackGpdbConnectionRead(d, meta)
}

func resourceApsaraStackGpdbConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := gpdb.CreateReleaseInstancePublicConnectionRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "gpdb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.DBInstanceId = parts[0]

	gpdbService := GpdbService{client}
	request.RegionId = client.RegionId
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		object, err := gpdbService.DescribeGpdbConnection(d.Id())
		if err != nil {
			return resource.NonRetryableError(WrapError(err))
		}
		request.CurrentConnectionString = object.ConnectionString

		var raw interface{}
		raw, err = client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.ReleaseInstancePublicConnection(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound", "InvalidCurrentConnectionString.NotFound", "AtLeastOneNetTypeExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{"NetAddressDeleting"}, []string{"Running"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gpdbService.GpdbInstanceStateRefreshFunc(request.DBInstanceId, []string{"Deleting"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return WrapError(gpdbService.WaitForGpdbConnection(d.Id(), Deleted, DefaultTimeoutMedium))
}
