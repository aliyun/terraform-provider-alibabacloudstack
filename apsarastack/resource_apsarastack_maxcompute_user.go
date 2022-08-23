package apsarastack

import (
	"fmt"
	"log"
	"strconv"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
)

func resourceApsaraStackMaxcomputeUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackMaxcomputeUserCreate,
		Read:   resourceApsaraStackMaxcomputeUserRead,
		Update: resourceApsaraStackMaxcomputeUserUpdate,
		Delete: resourceApsaraStackMaxcomputeUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"user_pk": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_type": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"organization_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"organization_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 255),
			},
		},
	}
}

func resourceApsaraStackMaxcomputeUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var response map[string]interface{}
	action := "CreateOdpsUser"
	product := "ascm"
	request := make(map[string]interface{})
	conn, err := client.NewAscmClient()
	if err != nil {
		return WrapError(err)
	}
	request["UserName"] = d.Get("user_name")
	if v, ok := d.GetOk("organization_id"); ok {
		request["OrganizationId"] = v
	} else {
		request["OrganizationId"] = client.Department
	}
	request["Description"] = d.Get("description")

	request["RegionId"] = client.RegionId
	request["RegionName"] = client.RegionId
	request["Product"] = product
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequestWithOrg(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_maxcompute_project", action, ApsaraStackSdkGoERROR)
	}
	if fmt.Sprintf(`%v`, response["code"]) != "200" {
		return WrapError(Error("CreateUpdateOdpsCu failed for " + response["Message"].(string)))
	}

	return resourceApsaraStackMaxcomputeUserRead(d, meta)
}

func resourceApsaraStackMaxcomputeUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	maxcomputeService := MaxcomputeService{client}
	object, err := maxcomputeService.DescribeMaxcomputeUser(d.Get("user_name").(string))
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource apsarastack_maxcompute_project_user maxcomputeService.DescribeMaxcomputeUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.SetId(strconv.Itoa(object.Data[0].ID))
	d.Set("user_id", object.Data[0].UserID)
	d.Set("user_pk", object.Data[0].UserPK)
	d.Set("user_name", object.Data[0].UserName)
	d.Set("user_type", object.Data[0].UserType)
	d.Set("organization_id", object.Data[0].OrganizationId)
	d.Set("organization_name", object.Data[0].OrganizationName)
	d.Set("description", object.Data[0].Description)
	return nil
}
func resourceApsaraStackMaxcomputeUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	update := false
	if d.HasChange("user_name") {
		update = true
	}
	if d.HasChange("organization_id") {
		update = true
	}
	if d.HasChange("organization_name") {
		update = true
	}
	if d.HasChange("description") {
		update = true
	}
	if update {
		var requestInfo *ecs.Client
		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		roleId, err := client.RoleIds()
		if err != nil {
			err = WrapErrorf(Error(GetNotFoundMessage("ASCM User", "defaultRoleId")), NotFoundMsg, ProviderERROR)
			return WrapErrorf(err, NotFoundMsg, ApsaraStackSdkGoERROR)
		}

		request.Method = "POST"
		request.Product = "ascm"
		request.Version = "2019-05-10"
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		OrganizationId := ""
		if v, ok := d.GetOk("organization_id"); ok {
			OrganizationId = strconv.Itoa(v.(int))
		} else {
			OrganizationId = client.Department
		}
		request.ApiName = "UpdateOdpsUser"
		request.Headers = map[string]string{
			"RegionId":              client.RegionId,
			"x-acs-roleid":          strconv.Itoa(roleId),
			"x-acs-resourcegroupid": client.ResourceGroup,
			"x-acs-regionid":        client.RegionId,
			"x-acs-organizationid":  client.Department,
		}
		request.QueryParams = map[string]string{
			"RegionId":         client.RegionId,
			"Id":               d.Get("id").(string),
			"UserId":           d.Get("user_id").(string),
			"UserName":         d.Get("user_name").(string),
			"UserType":         d.Get("user_type").(string),
			"OrganizationId":   OrganizationId,
			"OrganizationName": d.Get("organization_name").(string),
			"Description":      d.Get("description").(string),
			"Product":          "ascm",
			"ResourceGroupId":  client.ResourceGroup,
			"action":           "UpdateOdpsUser",
			"Version":          "2019-05-10",
		}
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Error OdpsUser Not Found"}) {
				return WrapErrorf(err, NotFoundMsg, ApsaraStackSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Get("user_name").(string), "UpdateOdpsUser", ApsaraStackSdkGoERROR)

		}
		addDebug("UpdateOdpsUser", raw, requestInfo, request)
	}
	return resourceApsaraStackMaxcomputeUserRead(d, meta)
}

func resourceApsaraStackMaxcomputeUserDelete(d *schema.ResourceData, meta interface{}) error {
	//ASCM不支持删除
	return nil
	client := meta.(*connectivity.ApsaraStackClient)
	action := "DeleteOdpsCu"
	var response map[string]interface{}
	conn, err := client.NewOdpsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"CuId":        d.Id(),
		"CuName":      d.Get("cu_name"),
		"ClusterName": d.Get("cluster_name"),
		"Product":     "ascm",
		"RegionId":    client.RegionId,
		"RegionName":  client.RegionId,
	}

	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"500"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
	}
	if IsExpectedErrorCodes(fmt.Sprintf("%v", response["code"]), []string{"102", "403"}) {
		return nil
	}
	if fmt.Sprintf(`%v`, response["code"]) != "200" {
		return WrapError(Error("DeleteOdpsCu failed for " + response["Message"].(string)))
	}
	return nil
}
