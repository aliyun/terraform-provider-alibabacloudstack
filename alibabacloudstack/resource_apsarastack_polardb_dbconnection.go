package alibabacloudstack

import (
	"fmt"
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
	setResourceFunc(resource, 
		resourceAlibabacloudStackPolardbConnectionCreate,
		resourceAlibabacloudStackPolardbConnectionRead,
		resourceAlibabacloudStackPolardbConnectionUpdate,
		resourceAlibabacloudStackPolardbConnectionDelete)
	return resource
}

func resourceAlibabacloudStackPolardbConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	PolardbService := PolardbService{client}

	instanceId := d.Get("instance_id").(string)
	prefix := d.Get("connection_prefix").(string)
	if prefix == "" {
		prefix = fmt.Sprintf("%stf", instanceId)
	}
	if err := PolardbService.WaitForConnectionDBInstance(d, client, instanceId, Running, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "AllocateInstancePublicConnection", "")

	request.QueryParams["DBInstanceId"] = instanceId
	request.QueryParams["Port"] = d.Get("port").(string)
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
	if err := PolardbService.WaitForDBConnection(d, client, d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := PolardbService.WaitForConnectionDBInstance(d, client, instanceId, Running, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func resourceAlibabacloudStackPolardbConnectionRead(d *schema.ResourceData, meta interface{}) error {
	submatch := dbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}
	parts, _ := ParseResourceId(d.Id(), 2)

	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbdb_instanceservice :=
		PolardbService{client}
	response, err := polardbdb_instanceservice.DescribeDBConnection(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_polardb_dbinstance", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	data := response
	d.Set("instance_id", parts[0])
	d.Set("connection_prefix", parts[1])

	d.Set("port", data.DBInstanceNetInfos.DBInstanceNetInfo[0].Port)
	d.Set("connection_string", data.DBInstanceNetInfos.DBInstanceNetInfo[0].ConnectionString)
	d.Set("ip_address", data.DBInstanceNetInfos.DBInstanceNetInfo[0].IPAddress)

	return nil
}

func resourceAlibabacloudStackPolardbConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolardbService{client}

	submatch := dbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	if d.HasChanges("connection_string", "port") {
		request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyDBInstanceConnectionString", "")
		request.QueryParams["DBInstanceId"] = parts[0]
		request.QueryParams["ConnectionStringPrefix"] = parts[1]

		if v, ok := d.GetOk("connection_string"); ok {
			request.QueryParams["CurrentConnectionString"] = v.(string)
		} else {
			return fmt.Errorf("CurrentConnectionString is required")
		}

		if v, ok := d.GetOk("port"); ok {
			request.QueryParams["Port"] = v.(string)
		} else {
			return fmt.Errorf("Port is required")
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

		// wait instance running after modifying
		if err := polardbService.WaitForConnectionDBInstance(d, client, parts[0], Running, DefaultTimeoutMedium); err != nil {
			return errmsgs.WrapError(err)
		}
	}
	return nil
}

func resourceAlibabacloudStackPolardbConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	PolardbService := PolardbService{client}

	submatch := dbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	split := strings.Split(d.Id(), COLON_SEPARATED)

	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ReleaseInstancePublicConnection", "")
	request.QueryParams["DBInstanceId"] = split[0]
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := PolardbService.DescribeDBConnection(d.Id())
		if err != nil {
			return resource.NonRetryableError(errmsgs.WrapError(err))
		}
		request.QueryParams["CurrentConnectionString"] = response.DBInstanceNetInfos.DBInstanceNetInfo[0].ConnectionString
		_, err = client.ProcessCommonRequest(request)

		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if errmsgs.NotFoundError(err) || errmsgs.IsExpectedErrors(err, []string{"InvalidCurrentConnectionString.NotFound", "AtLeastOneNetTypeExists"}) {
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

	return PolardbService.WaitForDBConnection(d, client, d.Id(), Deleted, DefaultTimeoutMedium)
}
