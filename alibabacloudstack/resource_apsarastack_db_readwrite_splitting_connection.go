package alibabacloudstack

import (
	"encoding/json"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const dbConnectionPrefixWithSuffixRegex = "^([a-zA-Z0-9\\-_]+)" + dbConnectionSuffixRegex + "$"

var dbConnectionPrefixWithSuffixRegexp = regexp.MustCompile(dbConnectionPrefixWithSuffixRegex)

func resourceAlibabacloudStackDBReadWriteSplittingConnection() *schema.Resource {
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
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 31),
			},
			"distribution_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Standard", "Custom"}, false),
			},
			"weight": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"max_delay_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackDBReadWriteSplittingConnectionCreate, 
		resourceAlibabacloudStackDBReadWriteSplittingConnectionRead, 
		resourceAlibabacloudStackDBReadWriteSplittingConnectionUpdate, 
		resourceAlibabacloudStackDBReadWriteSplittingConnectionDelete)
	return resource
}

func resourceAlibabacloudStackDBReadWriteSplittingConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}

	request := rds.CreateAllocateReadWriteSplittingConnectionRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = Trim(d.Get("instance_id").(string))
	request.MaxDelayTime = strconv.Itoa(d.Get("max_delay_time").(int))

	prefix, ok := d.GetOk("connection_prefix")
	if ok && prefix.(string) != "" {
		request.ConnectionStringPrefix = prefix.(string)
	}

	port, ok := d.GetOk("port")
	if ok {
		request.Port = strconv.Itoa(port.(int))
	}

	request.DistributionType = d.Get("distribution_type").(string)

	if weight, ok := d.GetOk("weight"); ok && weight != nil && len(weight.(map[string]interface{})) > 0 {
		if serial, err := json.Marshal(weight); err != nil {
			return errmsgs.WrapError(err)
		} else {
			request.Weight = string(serial)
		}
	}

	if err := resource.Retry(60*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.AllocateReadWriteSplittingConnection(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.DBReadInstanceNotReadyStatus) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if response, ok := raw.(*rds.AllocateReadWriteSplittingConnectionResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return err
	}

	d.SetId(request.DBInstanceId)

	// wait read write splitting connection ready after creation
	// for it may take up to 10 hours to create a readonly instance
	if err := rdsService.WaitForDBReadWriteSplitting(request.DBInstanceId, "", 60*60*10); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func resourceAlibabacloudStackDBReadWriteSplittingConnectionRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}

	err := rdsService.WaitForDBReadWriteSplitting(d.Id(), "", DefaultLongTimeout)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	object, err := rdsService.DescribeDBReadWriteSplittingConnection(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}

	d.Set("instance_id", d.Id())
	d.Set("connection_string", object.ConnectionString)
	d.Set("distribution_type", object.DistributionType)
	if port, err := strconv.Atoi(object.Port); err == nil {
		d.Set("port", port)
	}
	if mdt, err := strconv.Atoi(object.MaxDelayTime); err == nil {
		d.Set("max_delay_time", mdt)
	}
	if w, ok := d.GetOk("weight"); ok {
		documented := w.(map[string]interface{})
		for _, config := range object.DBInstanceWeights.DBInstanceWeight {
			if config.Availability != "Available" {
				delete(documented, config.DBInstanceId)
				continue
			}
			if config.Weight != "0" {
				if _, ok := documented[config.DBInstanceId]; ok {
					documented[config.DBInstanceId] = config.Weight
				}
			}
		}
		d.Set("weight", documented)
	}
	submatch := dbConnectionPrefixWithSuffixRegexp.FindStringSubmatch(object.ConnectionString)
	if len(submatch) > 1 {
		d.Set("connection_prefix", submatch[1])
	}

	return nil
}

func resourceAlibabacloudStackDBReadWriteSplittingConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}

	request := rds.CreateModifyReadWriteSplittingConnectionRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = d.Id()

	update := false

	if d.HasChange("max_delay_time") {
		request.MaxDelayTime = strconv.Itoa(d.Get("max_delay_time").(int))
		update = true
	}

	if !update && d.IsNewResource() {
		return nil
	}

	if d.HasChange("weight") {
		if weight, ok := d.GetOk("weight"); ok && weight != nil && len(weight.(map[string]interface{})) > 0 {
			if serial, err := json.Marshal(weight); err != nil {
				return err
			} else {
				request.Weight = string(serial)
			}
		}
		update = true
	}

	if d.HasChange("distribution_type") {
		request.DistributionType = d.Get("distribution_type").(string)
		update = true
	}

	if update {
		// wait instance running before modifying
		if err := rdsService.WaitForDBInstance(request.DBInstanceId, Running, 60*60); err != nil {
			return errmsgs.WrapError(err)
		}

		if err := resource.Retry(30*time.Minute, func() *resource.RetryError {
			raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.ModifyReadWriteSplittingConnection(request)
			})
			if err != nil {
				if errmsgs.IsExpectedErrors(err, errmsgs.OperationDeniedDBStatus) || errmsgs.IsExpectedErrors(err, errmsgs.DBReadInstanceNotReadyStatus) {
					return resource.RetryableError(err)
				}
				errmsg := ""
				if response, ok := raw.(*rds.ModifyReadWriteSplittingConnectionResponse); ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return err
		}

		// wait instance running after modifying
		if err := rdsService.WaitForDBInstance(request.DBInstanceId, Running, DefaultTimeoutMedium); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	return nil
}

func resourceAlibabacloudStackDBReadWriteSplittingConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}

	request := rds.CreateReleaseReadWriteSplittingConnectionRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = d.Id()

	if err := resource.Retry(30*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ReleaseReadWriteSplittingConnection(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			if errmsgs.NotFoundError(err) || errmsgs.IsExpectedErrors(err, []string{"InvalidRwSplitNetType.NotFound"}) {
				return nil
			}
			errmsg := ""
			if response, ok := raw.(*rds.ReleaseReadWriteSplittingConnectionResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return err
	}

	return errmsgs.WrapError(rdsService.WaitForDBReadWriteSplitting(d.Id(), Deleted, DefaultLongTimeout))
}
