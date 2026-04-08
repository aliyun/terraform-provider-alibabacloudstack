package alibabacloudstack

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const dbConnectionSuffixRegex = "\\.mysql\\.([a-zA-Z0-9\\-]+\\.){0,1}rds\\.aliyuncs\\.com"
const dbConnectionIdWithSuffixRegex = "^([a-zA-Z0-9\\-_]+:[a-zA-Z0-9\\-_]+)" + dbConnectionSuffixRegex + "$"

var dbConnectionIdWithSuffixRegexp = regexp.MustCompile(dbConnectionIdWithSuffixRegex)

func resourceAlibabacloudStackDBConnection() *schema.Resource {
	resource := &schema.Resource{
		SchemaVersion: 1, // Schema version for state migration support
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceAlibabacloudStackDBConnectionResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: migrateDBConnectionStateV0ToV1,
				Version: 0,
			},
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "public",
				ValidateFunc: validation.StringInSlice([]string{"public", "private"}, false),
				Description:  "Network type of the DB connection. Valid values: public (public network), private (private network).",
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
	setResourceFunc(resource, resourceAlibabacloudStackDBConnectionCreate, resourceAlibabacloudStackDBConnectionRead, resourceAlibabacloudStackDBConnectionUpdate, resourceAlibabacloudStackDBConnectionDelete)
	return resource
}

func resourceAlibabacloudStackDBConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}
	instanceId := d.Get("instance_id").(string)
	networkType := d.Get("network_type").(string)

	// Handle different network types: public (allocate new connection) vs private (query existing)
	if networkType == "public" {
		prefix := d.Get("connection_prefix").(string)
		if prefix == "" {
			prefix = instanceId
		}

		request := rds.CreateAllocateInstancePublicConnectionRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = instanceId
		request.ConnectionStringPrefix = prefix
		request.Port = fmt.Sprint(d.Get("port").(int))

		var raw interface{}
		var err error
		err = resource.Retry(8*time.Minute, func() *resource.RetryError {
			raw, err = client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.AllocateInstancePublicConnection(request)
			})
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			if err != nil {
				if errmsgs.IsExpectedErrors(err, errmsgs.OperationDeniedDBStatus...) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*rds.AllocateInstancePublicConnectionResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_db_connection", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

	}
	d.SetId(fmt.Sprintf("%s%s%s", instanceId, COLON_SEPARATED, networkType))
	if err := rdsService.WaitForDBConnection(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}
	// Wait for instance to be running after allocating connection
	if err := rdsService.WaitForDBInstance(instanceId, Running, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

// resourceAlibabacloudStackDBConnectionRead reads the DB connection information from the API
func resourceAlibabacloudStackDBConnectionRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}
	object, err := rdsService.DescribeDBConnection(d.Id())

	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("instance_id", parts[0])
	d.Set("connection_prefix", GetRdsConnectionPrefix(object.ConnectionString))
	port, err := strconv.Atoi(object.Port)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("port", port)
	d.Set("network_type", strings.ToLower(object.IPType))
	d.Set("connection_string", object.ConnectionString)
	d.Set("ip_address", object.IPAddress)

	return nil
}

// resourceAlibabacloudStackDBConnectionUpdate updates the DB connection configuration
func resourceAlibabacloudStackDBConnectionUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.IsNewResource() {
		return nil
	}

	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	if d.IsNewResource() && d.Get("network_type").(string) == "public" {
		return nil
	}
	if d.IsNewResource() && d.Get("network_type").(string) == "private" && d.Get("connection_prefix").(string) == "" && d.Get("connection_prefix").(string) != d.Get("instance_id").(string) {
		return nil
	}

	// Only update if port or connection_prefix has changed
	if d.HasChanges("port", "connection_prefix") {
		connection_string := d.Get("connection_string").(string)
		if connection_string == "" {
			object, err := rdsService.DescribeDBConnection(d.Id())
			if err != nil {
				return errmsgs.WrapError(err)
			}
			connection_string = object.ConnectionString
		}
		prefix := d.Get("connection_prefix").(string)
		if prefix == "" {
			prefix = d.Get("instance_id").(string)
		}
		request := rds.CreateModifyDBInstanceConnectionStringRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = parts[0]
		request.CurrentConnectionString = connection_string
		request.ConnectionStringPrefix = d.Get("connection_prefix").(string)
		request.Port = fmt.Sprint(d.Get("port").(int))

		if err := resource.Retry(8*time.Minute, func() *resource.RetryError {
			raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.ModifyDBInstanceConnectionString(request)
			})
			if err != nil {
				if errmsgs.IsExpectedErrors(err, errmsgs.OperationDeniedDBStatus...) {
					return resource.RetryableError(err)
				}
				errmsg := ""
				if raw != nil {
					response, ok := raw.(*rds.ModifyDBInstanceConnectionStringResponse)
					if ok {
						errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
					}
				}
				err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return err
		}

		// Wait for instance to be running after modifying connection
		if err := rdsService.WaitForDBInstance(request.DBInstanceId, Running, DefaultTimeoutMedium); err != nil {
			return errmsgs.WrapError(err)
		}
	}
	return nil
}

// resourceAlibabacloudStackDBConnectionDelete deletes the DB connection
func resourceAlibabacloudStackDBConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Get("network_type").(string) == "private" {
		// private connection strings can`t delete
		return nil
	}

	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}
	split := strings.Split(d.Id(), COLON_SEPARATED)
	request := rds.CreateReleaseInstancePublicConnectionRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = split[0]

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		object, err := rdsService.DescribeDBConnection(d.Id())
		if err != nil {
			return resource.NonRetryableError(errmsgs.WrapError(err))
		}
		request.CurrentConnectionString = object.ConnectionString
		var raw interface{}
		raw, err = client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ReleaseInstancePublicConnection(request)
		})

		if err != nil {
			if errmsgs.IsExpectedErrors(err, "OperationDenied.DBInstanceStatus") {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*rds.ReleaseInstancePublicConnectionResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			if errmsgs.NotFoundError(err) || errmsgs.IsExpectedErrors(err, "InvalidCurrentConnectionString.NotFound", "AtLeastOneNetTypeExists") {
				return nil
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return err
	}
	// Wait for connection to be deleted
	return rdsService.WaitForDBConnection(d.Id(), Deleted, DefaultTimeoutMedium)
}

func GetRdsConnectionPrefix(connecntion string) string {
	stringList := strings.Split(connecntion, ".")
	if len(stringList) > 0 {
		return stringList[0]
	}
	return ""
}

// resourceAlibabacloudStackDBConnectionResourceV0 returns the schema for version 0 of the DB connection resource.
// Version 0 schema only has instance_id in the ID, without network_type field.
func resourceAlibabacloudStackDBConnectionResourceV0() *schema.Resource {
	return &schema.Resource{
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
}

// migrateDBConnectionStateV0 migrates state from version 0 to version 1.
// Version 0: ID is just instance_id, no network_type field
// Version 1: ID is instance_id:network_type, with network_type defaulting to "public"
func migrateDBConnectionStateV0ToV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	// Get the old ID (which is just the instance_id in version 0)
	oldID := rawState["id"].(string)
	parts := strings.Split(oldID, COLON_SEPARATED)

	// Set network_type to "public" as default for migrated resources
	rawState["network_type"] = "public"

	// Update the ID to the new format: instance_id:network_type
	newID := fmt.Sprintf("%s%s%s", parts[0], COLON_SEPARATED, "public")
	rawState["id"] = newID

	return rawState, nil
}
