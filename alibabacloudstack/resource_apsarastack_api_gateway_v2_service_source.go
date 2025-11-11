package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAPIGatewayV2ServiceSource() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"1", "2", "3", "5"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_type_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"edas_end_point_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"check_type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"nacos_access_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nacos_secret_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"max_connection": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_idle_connection": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"connection_idle_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"database_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nacos_registry": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"edas_name_space_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"edas_access_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"edas_secret_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"edas_end_point": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"eureka_registry": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jdbc_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAPIGatewayV2ServiceSourceCreate, resourceAlibabacloudStackAPIGatewayV2ServiceSourceRead, resourceAlibabacloudStackAPIGatewayV2ServiceSourceUpdate, resourceAlibabacloudStackAPIGatewayV2ServiceSourceDelete)
	return resource
}

func resourceAlibabacloudStackAPIGatewayV2ServiceSourceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	requestBody := make(map[string]interface{})
	source_struct := make(map[string]interface{})
	source_struct["mode"] = "endPoint"
	if v, ok := d.GetOk("edas_end_point_port"); ok {
		source_struct["edasEndPointPort"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		source_struct["type"] = v
	}
	if v, ok := d.GetOk("check_type"); ok {
		source_struct["checkType"] = v
	}
	if v, ok := d.GetOk("nacos_access_key"); ok {
		source_struct["nacosAccessKey"] = v
	}
	if v, ok := d.GetOk("nacos_secret_key"); ok {
		source_struct["nacosSecretKey"] = v
	}
	if v, ok := d.GetOk("max_connection"); ok {
		source_struct["maxConnection"] = v
	}
	if v, ok := d.GetOk("max_idle_connection"); ok {
		source_struct["maxIdleConnection"] = v
	}
	if v, ok := d.GetOk("connection_idle_time"); ok {
		source_struct["connectionIdleTime"] = v
	}
	if v, ok := d.GetOk("database_type"); ok {
		source_struct["databaseType"] = v
	}
	if v, ok := d.GetOk("nacos_registry"); ok {
		source_struct["nacosRegistry"] = v
	}
	if v, ok := d.GetOk("edas_name_space_id"); ok {
		source_struct["edasNameSpaceId"] = v
	}
	if v, ok := d.GetOk("edas_access_key"); ok {
		source_struct["edasAccessKey"] = v
	}
	if v, ok := d.GetOk("edas_secret_key"); ok {
		source_struct["edasSecretKey"] = v
	}
	if v, ok := d.GetOk("edas_end_point"); ok {
		source_struct["edasEndPoint"] = v
	}
	if v, ok := d.GetOk("eureka_registry"); ok {
		source_struct["eurekaRegistry"] = v
	}
	if v, ok := d.GetOk("jdbc_url"); ok {
		source_struct["jdbcUrl"] = v
	}
	if v, ok := d.GetOk("username"); ok {
		source_struct["username"] = v
	}
	if v, ok := d.GetOk("password"); ok {
		source_struct["password"] = v
	}
	gwInstanceId := d.Get("instance_id").(string)
	requestBody["gwInstanceId"] = gwInstanceId
	requestBody["sourceStruct"] = source_struct
	requestBody["sourceName"] = d.Get("source_name")
	requestBody["sourceType"] = d.Get("source_type")
	requestBody["description"] = d.Get("description")

	// Call the create API
	resp, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "CreateSource", "/source/createSource", nil, nil, requestBody)
	if err != nil {
		return err
	}

	// Extract sourceId from response data
	sourceId, ok := resp["data"].(string)
	if !ok {
		return fmt.Errorf("failed to extract sourceId from response")
	}

	// Construct resource ID
	resourceId := fmt.Sprintf("%s:%s", gwInstanceId, sourceId)
	d.SetId(resourceId)
	return resourceAlibabacloudStackAPIGatewayV2ServiceSourceRead(d, meta)
}

func resourceAlibabacloudStackAPIGatewayV2ServiceSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	apiGatewayV2Service := &ApiGateWayV2Service{client}

	object, err := apiGatewayV2Service.DescribeApiGatewayV2ServiceSource(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_api_gateway_v2_service_source apiGatewayV2Service.DescribeApiGatewayV2ServiceSource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("instance_id", strings.Split(d.Id(), ":")[0])
	d.Set("source_id", object["sourceId"])
	d.Set("source_name", object["sourceName"])
	d.Set("source_type", object["sourceType"])
	d.Set("description", object["description"])
	d.Set("create_time", object["createTime"])
	d.Set("update_time", object["updateTime"])
	d.Set("source_type_name", object["sourceTypeName"])
	sourceStruct := object["sourceStruct"].(map[string]interface{})
	if v, exists := sourceStruct["edasEndPointPort"]; exists {
		d.Set("edas_end_point_port", v)
	}
	if v, exists := sourceStruct["type"]; exists {
		d.Set("type", v)
	}
	if v, exists := sourceStruct["checkType"]; exists {
		d.Set("check_type", v)
	}
	if v, exists := sourceStruct["nacosAccessKey"]; exists {
		d.Set("nacos_access_key", v)
	}
	if v, exists := sourceStruct["nacosSecretKey"]; exists {
		d.Set("nacos_secret_key", v)
	}
	if v, exists := sourceStruct["maxConnection"]; exists {
		d.Set("max_connection", v)
	}
	if v, exists := sourceStruct["maxIdleConnection"]; exists {
		d.Set("max_idle_connection", v)
	}
	if v, exists := sourceStruct["connectionIdleTime"]; exists {
		d.Set("connection_idle_time", v)
	}
	if v, exists := sourceStruct["databaseType"]; exists {
		d.Set("database_type", v)
	}
	if v, exists := sourceStruct["nacosRegistry"]; exists {
		d.Set("nacos_registry", v)
	}
	if v, exists := sourceStruct["edasNameSpaceId"]; exists {
		d.Set("edas_name_space_id", v)
	}
	if v, exists := sourceStruct["edasAccessKey"]; exists {
		d.Set("edas_access_key", v)
	}
	if v, exists := sourceStruct["edasSecretKey"]; exists {
		d.Set("edas_secret_key", v)
	}
	if v, exists := sourceStruct["edasEndPoint"]; exists {
		d.Set("edas_end_point", v)
	}
	if v, exists := sourceStruct["eurekaRegistry"]; exists {
		d.Set("eureka_registry", v)
	}
	if v, exists := sourceStruct["jdbcUrl"]; exists {
		urls := strings.Split(v.(string), "?")
		d.Set("jdbc_url", urls[0])
	}
	if v, exists := sourceStruct["username"]; exists {
		d.Set("username", v)
	}
	if v, exists := sourceStruct["password"]; exists {
		d.Set("password", v)
	}
	return nil
}

func resourceAlibabacloudStackAPIGatewayV2ServiceSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	gwInstanceId := parts[0]
	sourceId := parts[1]

	if d.IsNewResource() {
		return nil
	}
	if d.HasChanges("source_name", "source_type", "description", "edas_end_point_port", "type",
		"check_type", "nacos_access_key", "nacos_secret_key", "max_connection", "max_idle_connection",
		"connection_idle_time", "database_type", "nacos_registry", "edas_name_space_id", "edas_access_key",
		"edas_secret_key", "edas_end_point", "eureka_registry", "jdbc_url", "username", "password") {

		requestBody := make(map[string]interface{})
		source_struct := make(map[string]interface{})
		source_struct["mode"] = "endPoint"
		if v, ok := d.GetOk("edas_end_point_port"); ok {
			source_struct["edasEndPointPort"] = v
		}
		if v, ok := d.GetOk("type"); ok {
			source_struct["type"] = v
		}
		if v, ok := d.GetOk("check_type"); ok {
			source_struct["checkType"] = v
		}
		if v, ok := d.GetOk("nacos_access_key"); ok {
			source_struct["nacosAccessKey"] = v
		}
		if v, ok := d.GetOk("nacos_secret_key"); ok {
			source_struct["nacosSecretKey"] = v
		}
		if v, ok := d.GetOk("max_connection"); ok {
			source_struct["maxConnection"] = v
		}
		if v, ok := d.GetOk("max_idle_connection"); ok {
			source_struct["maxIdleConnection"] = v
		}
		if v, ok := d.GetOk("connection_idle_time"); ok {
			source_struct["connectionIdleTime"] = v
		}
		if v, ok := d.GetOk("database_type"); ok {
			source_struct["databaseType"] = v
		}
		if v, ok := d.GetOk("nacos_registry"); ok {
			source_struct["nacosRegistry"] = v
		}
		if v, ok := d.GetOk("edas_name_space_id"); ok {
			source_struct["edasNameSpaceId"] = v
		}
		if v, ok := d.GetOk("edas_access_key"); ok {
			source_struct["edasAccessKey"] = v
		}
		if v, ok := d.GetOk("edas_secret_key"); ok {
			source_struct["edasSecretKey"] = v
		}
		if v, ok := d.GetOk("edas_end_point"); ok {
			source_struct["edasEndPoint"] = v
		}
		if v, ok := d.GetOk("eureka_registry"); ok {
			source_struct["eurekaRegistry"] = v
		}
		if v, ok := d.GetOk("jdbc_url"); ok {
			source_struct["jdbcUrl"] = v
		}
		if v, ok := d.GetOk("username"); ok {
			source_struct["username"] = v
		}
		if v, ok := d.GetOk("password"); ok {
			source_struct["password"] = v
		}
		requestBody["gwInstanceId"] = gwInstanceId
		requestBody["sourceId"] = sourceId
		requestBody["sourceStruct"] = source_struct
		requestBody["sourceName"] = d.Get("source_name")
		requestBody["sourceType"] = d.Get("source_type")
		requestBody["description"] = d.Get("description")

		_, err = client.DoTeaRequest("POST", "csb2", "2023-02-06", "ModifySource", "/source/modifySource", nil, nil, requestBody)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_api_gateway_v2_service_source", "ModifySource", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}
	return nil
}

func resourceAlibabacloudStackAPIGatewayV2ServiceSourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	gwInstanceId := parts[0]
	sourceId := parts[1]

	reqQuery := map[string]interface{}{
		"sourceId":     sourceId,
		"gwInstanceId": gwInstanceId,
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "DeleteSource", "/source/deleteSource", nil, nil, reqQuery)
		if err != nil {
			err = errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteSource", errmsgs.AlibabacloudStackSdkGoERROR)
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errmsgs.WrapError(err)
	}
	return errmsgs.WrapError(err)
}
