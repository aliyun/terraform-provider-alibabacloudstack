package alibabacloudstack

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackPolardbConnection() *schema.Resource {
	resource := &schema.Resource{
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
	setResourceFunc(resource,
		resourceAlibabacloudStackPolardbConnectionCreate,
		resourceAlibabacloudStackPolardbConnectionRead,
		resourceAlibabacloudStackPolardbConnectionUpdate,
		resourceAlibabacloudStackPolardbConnectionDelete)
	return resource
}

func resourceAlibabacloudStackPolardbConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolardbService{client}

	instanceId := d.Get("instance_id").(string)
	prefix := d.Get("connection_prefix").(string)
	if prefix == "" {
		prefix = fmt.Sprintf("%stf", instanceId)
	}
	if err := polardbService.WaitForConnectionDBInstance(d, client, instanceId, Running, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "AllocateInstancePublicConnection", "")

	request.QueryParams["DBInstanceId"] = instanceId
	request.QueryParams["Port"] = strconv.Itoa(d.Get("port").(int))
	request.QueryParams["ConnectionStringPrefix"] = prefix

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_db_instance", "AllocateInstancePublicConnection", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	d.SetId(fmt.Sprintf("%s%s%s", instanceId, COLON_SEPARATED, prefix))
	stateConf := BuildStateConfByTimes([]string{"NET_CREATING"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, polardbService.PolardbDBInstanceStateRefreshFunc(d, client, instanceId, []string{"Deleting"}), 100)
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}

func resourceAlibabacloudStackPolardbConnectionRead(d *schema.ResourceData, meta interface{}) error {
	parts, _ := ParseResourceId(d.Id(), 2)

	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbdb_instanceservice :=
		PolardbService{client}
	response, err := polardbdb_instanceservice.DescribeDBConnection(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_polardb_dbinstance", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	data := response
	d.Set("instance_id", parts[0])
	d.Set("connection_prefix", parts[1])

	if port, err := toInt(data.Port); err != nil {
		return err
	} else {
		d.Set("port", port)
	}
	d.Set("connection_string", data.ConnectionString)
	d.Set("ip_address", data.IPAddress)

	return nil
}

func resourceAlibabacloudStackPolardbConnectionUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.IsNewResource() {
		return nil
	}

	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolardbService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	if d.HasChanges("connection_prefix", "port") {
		request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyDBInstanceConnectionString", "")
		request.QueryParams["DBInstanceId"] = parts[0]
		request.QueryParams["ConnectionStringPrefix"] = d.Get("connection_prefix").(string)
		request.QueryParams["Port"] = strconv.Itoa(d.Get("port").(int))

		if v, ok := d.GetOk("connection_string"); ok {
			request.QueryParams["CurrentConnectionString"] = v.(string)
		} else {
			return fmt.Errorf("CurrentConnectionString is required")
		}

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_polardb_db_instance", "ModifyDBInstanceConnectionString", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		stateConf := BuildStateConfByTimes([]string{"NET_MODIFYING"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, polardbService.PolardbDBInstanceStateRefreshFunc(d, client, parts[0], []string{"Deleting"}), 100)
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		d.SetId(fmt.Sprintf("%s%s%s", parts[0], COLON_SEPARATED, d.Get("connection_prefix").(string)))
	}
	return nil
}

func resourceAlibabacloudStackPolardbConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolardbService{client}

	split := strings.Split(d.Id(), COLON_SEPARATED)

	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ReleaseInstancePublicConnection", "")
	request.QueryParams["DBInstanceId"] = split[0]
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := polardbService.DescribeDBConnection(d.Id())
		if err != nil {
			return resource.NonRetryableError(errmsgs.WrapError(err))
		}
		request.QueryParams["CurrentConnectionString"] = response.ConnectionString
		rsp, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), rsp, request, request.QueryParams)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "OperationDenied.DBInstanceStatus") {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if errmsgs.NotFoundError(err) || errmsgs.IsExpectedErrors(err, "InvalidCurrentConnectionString.NotFound", "AtLeastOneNetTypeExists") {
				return nil
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}
	
	stateConf := BuildStateConfByTimes([]string{"NET_DELETING"}, []string{"Running"}, d.Timeout(schema.TimeoutDelete), 10*time.Second, polardbService.PolardbDBInstanceStateRefreshFunc(d, client, split[0], []string{"Deleting"}), 100)
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}
