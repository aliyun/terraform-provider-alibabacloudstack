package apsarastack

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"

	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceApsaraStackAdbConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackAdbConnectionCreate,
		Read:   resourceApsaraStackAdbConnectionRead,
		Delete: resourceApsaraStackAdbConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"connection_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9\\-]{4,28}[a-z0-9]$`), "The prefix must be 6 to 30 characters in length, and can contain lowercase letters, digits, and hyphens (-), must start with a letter and end with a digit or letter."),
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceApsaraStackAdbConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	adbService := AdbService{client}
	dbClusterId := d.Get("db_cluster_id").(string)
	prefix := d.Get("connection_prefix").(string)
	if prefix == "" {
		prefix = fmt.Sprintf("%stf", dbClusterId)
	}

	request := adb.CreateAllocateClusterPublicConnectionRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = dbClusterId
	request.ConnectionStringPrefix = prefix
	request.Headers["x-ascm-product-name"] = "adb"
	request.Headers["x-acs-organizationid"] = client.Department
	var raw interface{}
	var err error
	err = resource.Retry(8*time.Minute, func() *resource.RetryError {
		raw, err = client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
			return adbClient.AllocateClusterPublicConnection(request)
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
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_adb_connection", request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", dbClusterId, COLON_SEPARATED, request.ConnectionStringPrefix))

	if err := adbService.WaitForAdbConnection(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	// wait instance running after allocating
	if err := adbService.WaitForCluster(dbClusterId, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return resourceApsaraStackAdbConnectionRead(d, meta)
}

func resourceApsaraStackAdbConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	adbService := AdbService{client}
	object, err := adbService.DescribeAdbConnection(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("db_cluster_id", parts[0])
	d.Set("connection_prefix", parts[1])
	d.Set("port", object.Port)
	d.Set("connection_string", object.ConnectionString)
	d.Set("ip_address", object.IPAddress)

	return nil
}

func resourceApsaraStackAdbConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	adbService := AdbService{client}

	split := strings.Split(d.Id(), COLON_SEPARATED)
	request := adb.CreateReleaseClusterPublicConnectionRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = split[0]
	request.Headers["x-ascm-product-name"] = "adb"
	request.Headers["x-acs-organizationid"] = client.Department

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var raw interface{}
		raw, err := client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
			return adbClient.ReleaseClusterPublicConnection(request)
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
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	return adbService.WaitForAdbConnection(d.Id(), Deleted, DefaultTimeoutMedium)
}
