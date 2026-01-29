package alibabacloudstack

import (
	"encoding/json"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackOtsInstance() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 16),
			},
			"cluster_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"cluster_type"},
			},
			"cluster_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"cluster_name"},
			},
			"alias_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"specification": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vcu_quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sp_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackOtsInstanceCreate, resourceAlibabacloudStackOtsInstanceRead, resourceAlibabacloudStackOtsInstanceUpdate, resourceAlibabacloudStackOtsInstanceDelete)
	return resource
}

func resourceAlibabacloudStackOtsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	otsService := OtsService{client}

	instanceName := d.Get("name").(string)
	request := map[string]interface{}{
		"InstanceName": instanceName,
	}

	if v, ok := d.GetOk("cluster_name"); ok {
		request["ClusterName"] = v.(string)
	}
	if v, ok := d.GetOk("cluster_type"); ok {
		request["ClusterType"] = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		request["InstanceDescription"] = v.(string)
	}

	// Call the create API
	_, err := client.DoTeaRequest("POST", "tablestore", "2020-12-09", "CreateInstance", "/v2/openapi/createinstance", nil, nil, request)
	if err != nil {
		return err
	}

	// Set the resource ID based on the instance name

	d.SetId(instanceName)

	// Wait for the instance to be ready
	stateConf := BuildStateConf([]string{"creating"}, []string{"normal"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, otsService.OtsInstanceStateRefreshFunc(instanceName, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return err
	}

	return resourceAlibabacloudStackOtsInstanceRead(d, meta)
}

func resourceAlibabacloudStackOtsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	otsService := &OtsService{client}

	object, err := otsService.DescribeOtsInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_ots_instance otsService.DescribeOtsInstanceNew Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	// Set the basic attributes from the response
	d.Set("name", object["InstanceName"])
	d.Set("alias_name", object["AliasName"])
	d.Set("network", object["Network"])
	d.Set("description", object["InstanceDescription"])
	d.Set("cluster_name", object["ClusterName"])
	d.Set("storage_type", object["StorageType"])
	d.Set("create_time", object["CreateTime"])
	d.Set("table_quota", object["TableQuota"])
	d.Set("vcu_quota", object["VCUQuota"])
	d.Set("user_id", object["UserId"])
	d.Set("sp_instance_id", object["SPInstanceId"])
	d.Set("specification", object["InstanceSpecification"])

	// Handle tags
	if tags, ok := object["Tags"].([]interface{}); ok {
		tagsList := make([]map[string]interface{}, 0, len(tags))
		for _, tag := range tags {
			if tagMap, isMap := tag.(map[string]interface{}); isMap {
				tagItem := make(map[string]interface{})
				if tagKey, exists := tagMap["TagKey"]; exists {
					tagItem["tag_key"] = tagKey
				}
				if tagValue, exists := tagMap["TagValue"]; exists {
					tagItem["tag_value"] = tagValue
				}
				tagsList = append(tagsList, tagItem)
			}
		}
		d.Set("tags", tagsList)
	}

	return nil
}

func resourceAlibabacloudStackOtsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	otsService := OtsService{client}

	if d.IsNewResource() && d.Get("alias_name").(string) == "" {
		return nil
	}

	requestInfo := make(map[string]interface{})

	// Add instance name to request
	requestInfo["InstanceName"] = d.Id()

	// Check if alias_name has changed
	if d.HasChanges("alias_name", "description", "network") {
		if v, ok := d.GetOk("alias_name"); ok {
			requestInfo["AliasName"] = v.(string)
		}
		if v, ok := d.GetOk("description"); ok {
			requestInfo["InstanceDescription"] = v.(string)
		}
		if v, ok := d.GetOk("network"); ok {
			requestInfo["Network"] = v.(string)
		}

		// Call the update API
		_, err := client.DoTeaRequest("POST", "Tablestore", "2020-12-09", "UpdateInstance", "/v2/openapi/updateinstance", nil, nil, requestInfo)
		if err != nil {
			return err
		}

		// Wait for the instance to be updated by checking its status
		stateConf := BuildStateConf([]string{"updating"}, []string{"normal"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, otsService.OtsInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return err
		}
	}

	return nil
}

func resourceAlibabacloudStackOtsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	otsService := OtsService{client}

	if err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		request := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
		request.QueryParams["OpenApiAction"] = "DeleteInstance"
		request.QueryParams["ProductName"] = "ots"

		params := map[string]string{
			"Department":    client.Department,
			"ResourceGroup": client.ResourceGroup,
			"RegionId":      client.RegionId,
			"InstanceName":  d.Id(),
		}

		if content, err := json.Marshal(params); err != nil {
			return resource.NonRetryableError(err)
		} else {
			request.QueryParams["Params"] = string(content)
		}

		bresponse, err := client.ProcessCommonRequest(request)

		if err != nil {
			if bresponse == nil {
				return resource.RetryableError(errmsgs.WrapErrorf(err, "Process Common Request Failed"))
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteInstance", errmsgs.AlibabacloudStackOssGoSdk, errmsg))
		}
		return nil
	}); err != nil {
		return err
	}

	stateConf := BuildStateConf([]string{"deleting"}, []string{""}, d.Timeout(schema.TimeoutDelete), 10*time.Second, otsService.OtsInstanceStateRefreshFunc(d.Id(), []string{}))
	_, err := stateConf.WaitForState()
	return errmsgs.WrapError(err)
}
