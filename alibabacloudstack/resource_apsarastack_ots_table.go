package alibabacloudstack

import (
	"fmt"
	"log"

	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"strconv"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackOtsTable() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"table_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"primary_key": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{string(IntegerType), string(BinaryType), string(StringType)}, false),
						},
					},
				},
				MaxItems: 4,
				ForceNew: true,
			},
			"time_to_live": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(-1, INT_MAX),
			},
			"max_version": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, INT_MAX),
			},
			"deviation_cell_version_in_sec": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringConvertInt64(),
				Default:      "86400",
			},
		},
	}
	setResourceFunc(resource, resourceAliyunOtsTableCreate, resourceAliyunOtsTableRead, resourceAliyunOtsTableUpdate, resourceAliyunOtsTableDelete)
	return resource
}

func resourceAliyunOtsTableCreate(d *schema.ResourceData, meta interface{}) error {
	tableMeta := new(tablestore.TableMeta)
	instanceName := d.Get("instance_name").(string)
	tableName := d.Get("table_name").(string)
	tableMeta.TableName = tableName
	client := meta.(*connectivity.AlibabacloudStackClient)
	otsService := OtsService{client}
	if err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, e := otsService.DescribeOtsInstance(instanceName)
		if e != nil {
			if errmsgs.NotFoundError(e) {
				return resource.RetryableError(e)
			}
			return resource.NonRetryableError(e)
		}
		return nil
	}); err != nil {
		return errmsgs.WrapError(err)
	}
	for _, primaryKey := range d.Get("primary_key").([]interface{}) {
		pk := primaryKey.(map[string]interface{})
		pkValue := otsService.getPrimaryKeyType(pk["type"].(string))
		tableMeta.AddPrimaryKeyColumn(pk["name"].(string), pkValue)
	}
	tableOption := new(tablestore.TableOption)
	tableOption.TimeToAlive = d.Get("time_to_live").(int)
	tableOption.MaxVersion = d.Get("max_version").(int)
	if deviation, ok := d.GetOk("deviation_cell_version_in_sec"); ok {
		tableOption.DeviationCellVersionInSec, _ = strconv.ParseInt(deviation.(string), 10, 64)
	}
	reservedThroughput := new(tablestore.ReservedThroughput)

	request := new(tablestore.CreateTableRequest)
	request.TableMeta = tableMeta
	request.TableOption = tableOption
	request.ReservedThroughput = reservedThroughput
	if err := resource.Retry(6*time.Minute, func() *resource.RetryError {
		raw, err := client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			return tableStoreClient.CreateTable(request)
		})
		log.Printf("====================  CreateTable Response  ===================  \n%v\n", raw)
		log.Printf("====================  CreateTable Response err  ===================  \n%v\n", err)
		if err != nil {
			errmsg := ""
			if errmsgs.IsExpectedErrors(err, errmsgs.OtsTableIsTemporarilyUnavailable) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ots_table", "CreateTable", errmsgs.AliyunTablestoreGoSdk, errmsg))
		}
		return nil
	}); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s%s%s", instanceName, COLON_SEPARATED, tableName))
	return nil
}

func resourceAliyunOtsTableRead(d *schema.ResourceData, meta interface{}) error {
	instanceName, _, err := parseId(d, meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	client := meta.(*connectivity.AlibabacloudStackClient)
	otsService := OtsService{client}
	object, err := otsService.DescribeOtsTable(d.Id())

	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("instance_name", instanceName)
	d.Set("table_name", object.TableMeta.TableName)

	var pks []map[string]interface{}
	keys := object.TableMeta.SchemaEntry
	for _, v := range keys {
		item := make(map[string]interface{})
		item["name"] = *v.Name
		item["type"] = otsService.convertPrimaryKeyType(*v.Type)
		pks = append(pks, item)
	}
	d.Set("primary_key", pks)
	d.Set("time_to_live", object.TableOption.TimeToAlive)
	d.Set("max_version", object.TableOption.MaxVersion)
	d.Set("deviation_cell_version_in_sec", strconv.FormatInt(object.TableOption.DeviationCellVersionInSec, 10))

	return nil
}

func resourceAliyunOtsTableUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChanges("time_to_live", "max_version", "deviation_cell_version_in_sec") {
		instanceName, tableName, err := parseId(d, meta)
		if err != nil {
			return err
		}
		client := meta.(*connectivity.AlibabacloudStackClient)

		request := new(tablestore.UpdateTableRequest)
		request.TableName = tableName
		tableOption := new(tablestore.TableOption)

		tableOption.TimeToAlive = d.Get("time_to_live").(int)
		tableOption.MaxVersion = d.Get("max_version").(int)
		if deviation, ok := d.GetOk("deviation_cell_version_in_sec"); ok {
			tableOption.DeviationCellVersionInSec, _ = strconv.ParseInt(deviation.(string), 10, 64)
		}

		request.TableOption = tableOption
		if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			raw, err := client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
				return tableStoreClient.UpdateTable(request)
			})
			if err != nil {
				errmsg := ""
				if errmsgs.IsExpectedErrors(err, errmsgs.OtsTableIsTemporarilyUnavailable) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "UpdateTable", errmsgs.AliyunTablestoreGoSdk, errmsg))
			}
			addDebug("UpdateTable", raw, request)
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

func resourceAliyunOtsTableDelete(d *schema.ResourceData, meta interface{}) error {
	instanceName, tableName, err := parseId(d, meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	client := meta.(*connectivity.AlibabacloudStackClient)
	otsService := OtsService{client}
	req := new(tablestore.DeleteTableRequest)
	req.TableName = tableName
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			return tableStoreClient.DeleteTable(req)
		})
		if err != nil {
			errmsg := ""
			if errmsgs.IsExpectedErrors(err, errmsgs.OtsTableIsTemporarilyUnavailable) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteTable", errmsgs.AliyunTablestoreGoSdk, errmsg))
		}
		addDebug("DeleteTable", raw, req)
		return nil
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "OTSObjectNotExist") {
			return nil
		}
		return err
	}
	return errmsgs.WrapError(otsService.WaitForOtsTable(instanceName, tableName, Deleted, DefaultTimeout))
}

func parseId(d *schema.ResourceData, meta interface{}) (instanceName, tableName string, err error) {
	split := strings.Split(d.Id(), COLON_SEPARATED)
	if len(split) == 1 {
		if meta.(*connectivity.AlibabacloudStackClient).OtsInstanceName != "" {
			tableName = split[0]
			instanceName = meta.(*connectivity.AlibabacloudStackClient).OtsInstanceName
			d.SetId(fmt.Sprintf("%s%s%s", instanceName, COLON_SEPARATED, tableName))
		} else {
			err = errmsgs.WrapError(errmsgs.Error("From Provider version 1.10.0, the provider field 'ots_instance_name' has been deprecated and " +
				"you should use resource alibabacloudstack_ots_table's new field 'instance_name' and 'table_name' to re-import this resource."))
			return
		}
	} else {
		instanceName = split[0]
		tableName = split[1]
	}

	return
}
