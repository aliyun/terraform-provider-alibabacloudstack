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
)

func resourceAlibabacloudStackApiGatewayV2RouteGroup() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"gw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"base_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"group_id": {
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
			"editable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackApiGatewayV2RouteGroupCreate, resourceAlibabacloudStackApiGatewayV2RouteGroupRead, resourceAlibabacloudStackApiGatewayV2RouteGroupUpdate, resourceAlibabacloudStackApiGatewayV2RouteGroupDelete)
	return resource
}

func resourceAlibabacloudStackApiGatewayV2RouteGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	apiGatewayV2Service := ApiGatewayV2Service{client}

	request := make(map[string]interface{})
	request["name"] = d.Get("name")
	request["basePath"] = d.Get("base_path")
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("domain_ids"); ok && v.(*schema.Set).Len() > 0 {
		request["domainIds"] = v.(*schema.Set).List()
	}
	gwInstanceId := d.Get("gw_instance_id").(string)
	request["gwInstanceId"] = gwInstanceId

	resp, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "CreateGroup", "", nil, request, nil)
	if err != nil {
		return err
	}

	groupId, ok := resp["data"].(string)
	if !ok {
		return fmt.Errorf("failed to get groupId from response")
	}

	// Construct resource ID using gwInstanceId and groupId
	resourceId := fmt.Sprintf("%s:%s", gwInstanceId, groupId)
	d.SetId(resourceId)

	// Wait for the resource to be available
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 3*time.Second,
		apiGatewayV2Service.ApiGatewayV2RouteGroupStateRefreshFunc(resourceId, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for ApiGateway V2 Route Group (%s) to be created failed: %v", resourceId, err)
	}

	// Set the computed group_id field
	if err := d.Set("group_id", groupId); err != nil {
		return err
	}

	return nil
}

func resourceAlibabacloudStackApiGatewayV2RouteGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	apiGatewayV2Service := ApiGatewayV2Service{client}

	object, err := apiGatewayV2Service.DescribeRouteGroup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_api_gateway_v2_route_group apiGatewayV2Service.DescribeRouteGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("gw_instance_id", object["gwInstanceId"])
	d.Set("name", object["name"])
	d.Set("base_path", object["basePath"])
	d.Set("description", object["description"])
	d.Set("group_id", object["groupId"])
	d.Set("create_time", object["createTime"])
	d.Set("update_time", object["updateTime"])
	d.Set("editable", object["editable"])

	// Set domains
	if domains, ok := object["domains"].([]interface{}); ok {
		domainList := make([]map[string]interface{}, 0)
		for _, domain := range domains {
			if domainMap, ok := domain.(map[string]interface{}); ok {
				domainList = append(domainList, map[string]interface{}{
					"protocol":     domainMap["protocol"],
					"create_time":  domainMap["createTime"],
					"domain":       domainMap["domain"],
					"domain_id":    domainMap["domainId"],
				})
			}
		}
		d.Set("domains", domainList)
	}

	// Set domain_ids
	if domains, ok := object["domains"].([]interface{}); ok {
		domainIds := make([]string, 0)
		for _, domain := range domains {
			if domainMap, ok := domain.(map[string]interface{}); ok {
				domainIds = append(domainIds, domainMap["domainId"].(string))
			}
		}
		d.Set("domain_ids", domainIds)
	}

	return nil
}

func resourceAlibabacloudStackApiGatewayV2RouteGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	apiGatewayV2Service := ApiGatewayV2Service{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return err
	}
	gwInstanceId := parts[0]
	groupId := parts[1]

	if d.IsNewResource() {
		return nil
	}

	if d.HasChanges("name", "base_path", "description") {
		// Prepare the base request for ModifyGroup API
		modifyReq := map[string]interface{}{
			"groupId":       groupId,
			"gwInstanceId":  gwInstanceId,
			"name":          d.Get("name"),
			"basePath":      d.Get("base_path"),
			"description":   d.Get("description"),
		}

		// Fetch current domain information to include in the modify request
		describeReq := map[string]interface{}{
			"groupId":       groupId,
			"gwInstanceId":  gwInstanceId,
		}
		raw, err := client.DoTeaRequest("GET", "csb2", "2023-02-06", "GetGroup", "", nil, describeReq, nil)
		if err != nil {
			return fmt.Errorf("failed to get group details: %v", err)
		}

		if v, ok := raw["data"]; ok && v != nil {
			data := v.(map[string]interface{})
			if domains, ok := data["domains"].([]interface{}); ok {
				modifyReq["domains"] = domains
			}
			if createTime, ok := data["createTime"].(string); ok {
				modifyReq["createTime"] = createTime
			}
			if editable, ok := data["editable"].(bool); ok {
				modifyReq["editable"] = editable
			}
		}

		// Include domain IDs from current state
		if v, ok := d.GetOk("domain_ids"); ok {
			modifyReq["domainIds"] = v.(*schema.Set).List()
		}

		_, err = client.DoTeaRequest("POST", "csb2", "2023-02-06", "ModifyGroup", "", nil, nil, modifyReq)
		if err != nil {
			return fmt.Errorf("failed to modify route group: %v", err)
		}

		// Wait for the update to complete
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 3*time.Second,
			apiGatewayV2Service.ApiGatewayV2RouteGroupStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("waiting for route group to become available failed: %v", err)
		}
	}

	if d.HasChange("domain_ids") {
		oldDomainIds, newDomainIds := d.GetChange("domain_ids")
		oldSet := oldDomainIds.(*schema.Set)
		newSet := newDomainIds.(*schema.Set)

		// Domains to add
		addDomains := newSet.Difference(oldSet).List()
		// Domains to remove
		removeDomains := oldSet.Difference(newSet).List()

		if len(addDomains) > 0 {
			addReq := map[string]interface{}{
				"groupId":       groupId,
				"gwInstanceId":  gwInstanceId,
				"domainIds":     addDomains,
			}
			_, err = client.DoTeaRequest("POST", "csb2", "2023-02-06", "AddGroupDomain", "", nil, nil, addReq)
			if err != nil {
				return fmt.Errorf("failed to add domains to route group: %v", err)
			}
		}

		if len(removeDomains) > 0 {
			for _, domainId := range removeDomains {
				deleteReq := map[string]interface{}{
					"groupId":       groupId,
					"gwInstanceId":  gwInstanceId,
					"domainId":      domainId,
				}
				_, err = client.DoTeaRequest("POST", "csb2", "2023-02-06", "DeleteGroupDomain", "", nil, nil, deleteReq)
				if err != nil {
					return fmt.Errorf("failed to delete domain from route group: %v", err)
				}
			}
		}

		// Wait for the domain changes to complete
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 3*time.Second,
			apiGatewayV2Service.ApiGatewayV2RouteGroupStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("waiting for route group domain changes to complete failed: %v", err)
		}
	}

	return resourceAlibabacloudStackApiGatewayV2RouteGroupRead(d, meta)
}

func resourceAlibabacloudStackApiGatewayV2RouteGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	apiGatewayService := ApiGatewayV2Service{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	gwInstanceId := parts[0]
	groupId := parts[1]

	reqQuery := map[string]interface{}{
		"gwInstanceId": gwInstanceId,
		"groupId":      groupId,
	}

	_, err = client.DoTeaRequest("POST", "csb2", "2023-02-06", "DeleteGroup", "", nil, reqQuery, nil)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteGroup", AlibabacloudStackApiGatewayV2OperationFailed)
	}

	stateConf := BuildStateConf([]string{"Available"}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Second, apiGatewayService.ApiGatewayV2RouteGroupStateRefreshFunc(d.Id(), []string{}))
	_, err = stateConf.WaitForState()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteGroup", AlibabacloudStackApiGatewayV2OperationFailed)
	}

	return nil
}

