package alibabacloudstack

// Generated By apsara-orchestration-generator
// Product POLARDB Resouce Database
import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackPolardbDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackPolardbDatabaseCreate,
		Read:   resourceAlibabacloudStackPolardbDatabaseRead,
		Update: resourceAlibabacloudStackPolardbDatabaseUpdate,
		Delete: resourceAlibabacloudStackPolardbDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{

			"accounts": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"account": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"account_privilege": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"account_privilege_detail": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"character_set_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"data_base_description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"data_base_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"data_base_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackPolardbDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "CreateDatabase", "")
	PolardbCreatedatabaseResponse := PolardbCreatedatabaseResponse{}

	if v, ok := d.GetOk("character_set_name"); ok && v != "" {
		request.QueryParams["CharacterSetName"] = v.(string)
	} else {
		return fmt.Errorf("CharacterSetName is required")
	}

	if v, ok := d.GetOk("data_base_description"); ok && v != "" {
		request.QueryParams["DBDescription"] = v.(string)
	}

	if v, ok := d.GetOk("data_base_instance_id"); ok && v != "" {
		request.QueryParams["DBInstanceId"] = v.(string)
	} else {
		return fmt.Errorf("DataBaseInstanceId is required")
	}

	if v, ok := d.GetOk("data_base_name"); ok && v != "" {
		request.QueryParams["DBName"] = v.(string)
	} else {
		return fmt.Errorf("DataBaseName is required")
	}

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_database", "CreateDatabase", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbCreatedatabaseResponse)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"alibabacloudstack_polardb_database", "CreateDatabase", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	data_base_instance_id := d.Get("data_base_instance_id").(string)

	data_base_name := d.Get("data_base_name").(string)

	d.SetId(fmt.Sprintf("%s", data_base_instance_id+":"+data_base_name))
	return resourceAlibabacloudStackPolardbDatabaseUpdate(d, meta)

}

func resourceAlibabacloudStackPolardbDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.HasChanges("data_base_description") {
		request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyDBDescription", "")
		PolardbModifydbdescriptionResponse := PolardbModifydbdescriptionResponse{}

		if v, ok := d.GetOk("data_base_description"); ok {
			request.QueryParams["DBDescription"] = v.(string)
		} else {
			return fmt.Errorf("DataBaseDescription is required")
		}

		if v, ok := d.GetOk("data_base_instance_id"); ok {
			request.QueryParams["DBInstanceId"] = v.(string)
		} else {
			return fmt.Errorf("DataBaseInstanceId is required")
		}

		if v, ok := d.GetOk("data_base_name"); ok {
			request.QueryParams["DBName"] = v.(string)
		} else {
			return fmt.Errorf("DataBaseName is required")
		}

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_polardb_database", "ModifyDBDescription", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbModifydbdescriptionResponse)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_polardb_database", "ModifyDBDescription", errmsgs.AlibabacloudStackSdkGoERROR)
		}

	}

	return resourceAlibabacloudStackPolardbDatabaseRead(d, meta)
}

func resourceAlibabacloudStackPolardbDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbdatabaseservice :=
		PolardbService{client}
	response, err := polardbdatabaseservice.DescribeDBDatabase(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_polardb_database", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	data := response.Databases.Database[0]
	d.Set("character_set_name", data.CharacterSetName)

	d.Set("data_base_description", data.DBDescription)

	d.Set("data_base_instance_id", data.DBInstanceId)

	d.Set("data_base_name", data.DBName)

	d.Set("engine", data.Engine)

	d.Set("status", data.DBStatus)

	return nil
}

func resourceAlibabacloudStackPolardbDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "DeleteDatabase", "")

	if v, ok := d.GetOk("data_base_instance_id"); ok && v != "" {
		request.QueryParams["DBInstanceId"] = v.(string)
	} else {
		return fmt.Errorf("DataBaseInstanceId is required")
	}

	if v, ok := d.GetOk("data_base_name"); ok && v != "" {
		request.QueryParams["DBName"] = v.(string)
	} else {
		return fmt.Errorf("DataBaseName is required")
	}

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_database", "DeleteDatabase", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	return nil
}

type PolardbCreatedatabaseResponse struct {
	RequestId string `json:"RequestId"`
}

type PolardbCopydatabasebetweeninstancesResponse struct {
	RequestId    string `json:"RequestId"`
	DBInstanceId string `json:"DBInstanceId"`
}
type PolardbCreateonlinedatabasetaskResponse struct {
	RequestId string `json:"RequestId"`
}
type PolardbModifydbdescriptionResponse struct {
	RequestId string `json:"RequestId"`
}

type PolardbDeletedatabaseResponse struct {
	RequestId string `json:"RequestId"`
}
