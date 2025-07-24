package alibabacloudstack

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackPolardbReadWriteSplittingConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackPolardbReadWriteSplittingConnectionCreate,
		Read:   resourceAlibabacloudStackPolardbReadWriteSplittingConnectionRead,
		Update: resourceAlibabacloudStackPolardbReadWriteSplittingConnectionUpdate,
		Delete: resourceAlibabacloudStackPolardbReadWriteSplittingConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
}

func resourceAlibabacloudStackPolardbReadWriteSplittingConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "AllocateReadWriteSplittingConnection", "")

	request.QueryParams["DBInstanceId"] = Trim(d.Get("instance_id").(string))
	request.QueryParams["MaxDelayTime"] = strconv.Itoa(d.Get("max_delay_time").(int))

	prefix, ok := d.GetOk("connection_prefix")
	if ok && prefix.(string) != "" {
		request.QueryParams["ConnectionStringPrefix"] = prefix.(string)
	}

	port, ok := d.GetOk("port")
	if ok {
		request.QueryParams["Port"] = strconv.Itoa(port.(int))
	}

	request.QueryParams["DistributionType"] = d.Get("distribution_type").(string)

	if weight, ok := d.GetOk("weight"); ok && weight != nil && len(weight.(map[string]interface{})) > 0 {
		if serial, err := json.Marshal(weight); err != nil {
			return errmsgs.WrapError(err)
		} else {
			request.QueryParams["Weight"] = string(serial)
		}
	}

	if err := resource.Retry(60*time.Minute, func() *resource.RetryError {

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return resource.RetryableError(errmsgs.WrapErrorf(err, "Process Common Request Failed"))
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_db_instance", "AllocateInstancePublicConnection", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		return nil
	}); err != nil {
		return err
	}

	d.SetId(request.QueryParams["DBInstanceId"])

	return resourceAlibabacloudStackPolardbReadWriteSplittingConnectionUpdate(d, meta)
}

func resourceAlibabacloudStackPolardbReadWriteSplittingConnectionRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)

	client := meta.(*connectivity.AlibabacloudStackClient)
	PolardbService := PolardbService{client}

	res, err := PolardbService.DoPolardbDescribedbinstancenetinfoRequest(d, client, d.Id())
	if res != nil {
		for _, conn := range res.DBInstanceNetInfos.DBInstanceNetInfo {
			if conn.ConnectionStringType != "ReadWriteSplitting" {
				continue
			}
			if conn.MaxDelayTime == "" {
				continue
			}

			// if _, err := strconv.Atoi(conn.MaxDelayTime); err != nil {
			// 	return ds, err
			// }
			d.Set("instance_id", d.Id())
			d.Set("connection_string", conn.ConnectionString)
			d.Set("distribution_type", conn.DistributionType)
			if port, err := strconv.Atoi(conn.Port); err == nil {
				d.Set("port", port)
			}
			if mdt, err := strconv.Atoi(conn.MaxDelayTime); err == nil {
				d.Set("max_delay_time", mdt)
			}
			if w, ok := d.GetOk("weight"); ok {
				documented := w.(map[string]interface{})
				for _, config := range conn.DBInstanceWeights.DBInstanceWeight {
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
			submatch := dbConnectionPrefixWithSuffixRegexp.FindStringSubmatch(conn.ConnectionString)
			if len(submatch) > 1 {
				d.Set("connection_prefix", submatch[1])
			}
		}
	}

	if err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func resourceAlibabacloudStackPolardbReadWriteSplittingConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	PolardbService := PolardbService{client}
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyReadWriteSplittingConnection", "")

	request.QueryParams["DBInstanceId"] = d.Id()

	update := false

	if d.HasChange("max_delay_time") {
		request.QueryParams["MaxDelayTime"] = strconv.Itoa(d.Get("max_delay_time").(int))
		update = true
	}

	if !update && d.IsNewResource() {
		return resourceAlibabacloudStackPolardbReadWriteSplittingConnectionRead(d, meta)
	}

	if d.HasChange("weight") {
		if weight, ok := d.GetOk("weight"); ok && weight != nil && len(weight.(map[string]interface{})) > 0 {
			if serial, err := json.Marshal(weight); err != nil {
				return err
			} else {
				request.QueryParams["Weight"] = string(serial)
			}
		}
		update = true
	}

	if d.HasChange("distribution_type") {
		request.QueryParams["DistributionType"] = d.Get("distribution_type").(string)
		update = true
	}

	if update {
		// wait instance running before modifying
		if err := PolardbService.WaitForDBInstance(d.Id(), Running, 60*60); err != nil {
			return errmsgs.WrapError(err)
		}

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_db_instance", "AllocateInstancePublicConnection", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		// wait instance running after modifying
		if err := PolardbService.WaitForDBInstance(d.Id(), Running, DefaultLongTimeout); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	return resourceAlibabacloudStackPolardbReadWriteSplittingConnectionRead(d, meta)
}

func resourceAlibabacloudStackPolardbReadWriteSplittingConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ReleaseReadWriteSplittingConnection", "")
	request.QueryParams["DBInstanceId"] = d.Id()

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_db_instance", "AllocateInstancePublicConnection", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	return nil
}
