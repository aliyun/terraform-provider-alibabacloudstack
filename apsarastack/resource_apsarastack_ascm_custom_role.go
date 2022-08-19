package apsarastack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
	"time"
)

func resourceApsaraStackAscmRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackAscmRoleCreate,
		Read:   resourceApsaraStackAscmRoleRead,
		Update: resourceApsaraStackAscmRoleUpdate,
		Delete: resourceApsaraStackAscmRoleDelete,
		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"organization_visibility": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_range": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"privileges": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
		},
	}
}
func resourceApsaraStackAscmRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var requestInfo *ecs.Client
	ascmService := AscmService{client}
	name := d.Get("role_name").(string)
	description := d.Get("description").(string)
	roleRange := d.Get("role_range").(string)
	var priv string
	var privs []string
	if v, ok := d.GetOk("privileges"); ok {
		privs = expandStringList(v.(*schema.Set).List())
		for i, k := range privs {
			if i != 0 {
				priv = fmt.Sprintf("%s\",\"%s", priv, k)
			} else {
				priv = k
			}
		}
	}
	check, err := ascmService.DescribeAscmCustomRole(name)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_custom_role", "role alreadyExist", ApsaraStackSdkGoERROR)
	}
	organizationvisibility := d.Get("organization_visibility").(string)
	if len(check.Data) == 0 {
		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":               client.RegionId,
			"AccessKeySecret":        client.SecretKey,
			"Product":                "ascm",
			"Action":                 "CreateRole",
			"Version":                "2019-05-10",
			"ProductName":            "ascm",
			"roleName":               name,
			"description":            description,
			"roleRange":              roleRange,
			"organizationVisibility": organizationvisibility,
			//"privileges":             fmt.Sprintf("[\"%s\"]", priv),
			"params": fmt.Sprintf("{\"privileges\":%s}", priv),
		}
		request.Method = "POST"
		request.Product = "ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "CreateRole"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_custom_role", "CreateRole", raw)
		}
		addDebug("CreateRole", raw, requestInfo, request)

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_custom_role", "CreateRole", ApsaraStackSdkGoERROR)
		}
		addDebug("CreateRole", raw, requestInfo, bresponse.GetHttpContentString())
	}
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		check, err = ascmService.DescribeAscmCustomRole(name)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	d.SetId(name + COLON_SEPARATED + fmt.Sprint(check.Data[0].ID))
	return resourceApsaraStackAscmRoleUpdate(d, meta)

}

func resourceApsaraStackAscmRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceApsaraStackAscmRoleRead(d, meta)
}

func resourceApsaraStackAscmRoleRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmCustomRole(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	log.Printf("Privileges for did[0]:%v", object.Data[0].Privileges)
	d.Set("role_name", did[0])
	d.Set("organization_visibility", object.Data[0].OrganizationVisibility)
	d.Set("role_id", object.Data[0].ID)
	d.Set("description", object.Data[0].Description)
	return nil
}

func resourceApsaraStackAscmRoleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	check, err := ascmService.DescribeAscmCustomRole(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsRoleExist", ApsaraStackSdkGoERROR)
	}
	addDebug("IsCustomRoleExist", check, requestInfo, map[string]string{"roleName": did[0]})
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Product":         "ascm",
			"Action":          "RemoveRole",
			"Version":         "2019-05-10",
			"ProductName":     "ascm",
			"roleName":        did[0],
		}

		request.Method = "POST"
		request.Product = "ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "RemoveRole"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		_, err = ascmService.DescribeAscmCustomRole(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return nil
}
