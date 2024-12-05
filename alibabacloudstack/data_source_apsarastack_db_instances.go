package alibabacloudstack

import (
	"encoding/json"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackDBInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDBInstancesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"engine": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string(MySQL), string(SQLServer), string(PPAS), string(PostgreSQL)}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Primary", "Readonly", "Guard", "Temp"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connection_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Standard", "Safe"}, false),
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"guard_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"temp_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"readonly_instance_ids": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackDBInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := rds.CreateDescribeDBInstancesRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.Engine = d.Get("engine").(string)
	request.DBInstanceStatus = d.Get("status").(string)
	request.DBInstanceType = d.Get("db_type").(string)
	request.VpcId = d.Get("vpc_id").(string)
	request.VSwitchId = d.Get("vswitch_id").(string)
	request.ConnectionMode = d.Get("connection_mode").(string)
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := v.(map[string]interface{})
		bs, err := json.Marshal(tagsMap)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.Tags = string(bs)
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	var dbi []rds.DBInstance

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		nameRegex = r
	}

	// ids
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[vv.(string)] = vv.(string)
		}
	}

	for {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeDBInstances(request)
		})
		response, ok := raw.(*rds.DescribeDBInstancesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_db_instances", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if !ok {
			return errmsgs.WrapErrorf(nil, errmsgs.DataDefaultErrorMsg, "alibabacloudstack_db_instances", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, "Invalid response type")
		}
		if len(response.Items.DBInstance) < 1 {
			break
		}

		for _, item := range response.Items.DBInstance {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.DBInstanceDescription) {
					continue
				}
			}

			if len(idsMap) > 0 {
				if _, ok := idsMap[item.DBInstanceId]; !ok {
					continue
				}
			}

			dbi = append(dbi, item)
		}

		if len(response.Items.DBInstance) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.PageNumber = page
	}
	return rdsInstancesDescription(d, meta, dbi)
}

func rdsInstancesDescription(d *schema.ResourceData, meta interface{}, dbi []rds.DBInstance) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}

	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, item := range dbi {
		readOnlyInstanceIDs := []string{}
		for _, id := range item.ReadOnlyDBInstanceIds.ReadOnlyDBInstanceId {
			readOnlyInstanceIDs = append(readOnlyInstanceIDs, id.DBInstanceId)
		}
		instance, err := rdsService.DescribeDBInstance(item.DBInstanceId)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		mapping := map[string]interface{}{
			"id":                    item.DBInstanceId,
			"name":                  item.DBInstanceDescription,
			"charge_type":           item.PayType,
			"db_type":               item.DBInstanceType,
			"region_id":             item.RegionId,
			"create_time":           item.CreateTime,
			"expire_time":           item.ExpireTime,
			"status":                item.DBInstanceStatus,
			"engine":                item.Engine,
			"engine_version":        item.EngineVersion,
			"net_type":              item.DBInstanceNetType,
			"connection_mode":       item.ConnectionMode,
			"instance_type":         item.DBInstanceClass,
			"availability_zone":     item.ZoneId,
			"master_instance_id":    item.MasterInstanceId,
			"guard_instance_id":     item.GuardDBInstanceId,
			"temp_instance_id":      item.TempDBInstanceId,
			"readonly_instance_ids": readOnlyInstanceIDs,
			"vpc_id":                item.VpcId,
			"vswitch_id":            item.VSwitchId,
			"connection_string":     instance.ConnectionString,
			"port":                  instance.Port,
			"instance_storage":      instance.DBInstanceStorage,
		}

		ids = append(ids, item.DBInstanceId)
		names = append(names, item.DBInstanceDescription)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
