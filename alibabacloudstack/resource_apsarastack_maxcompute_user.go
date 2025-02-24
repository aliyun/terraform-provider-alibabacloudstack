package alibabacloudstack

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackMaxcomputeUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackMaxcomputeUserCreate,
		Read:   resourceAlibabacloudStackMaxcomputeUserRead,
		Update: resourceAlibabacloudStackMaxcomputeUserUpdate,
		Delete: resourceAlibabacloudStackMaxcomputeUserDelete,
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

func resourceAlibabacloudStackMaxcomputeUserCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateOdpsUser"
	request := make(map[string]interface{})
	request["UserName"] = d.Get("user_name")
	if v, ok := d.GetOk("organization_id"); ok {
		request["OrganizationId"] = v
	} else {
		request["OrganizationId"] = client.Department
	}
	request["Description"] = d.Get("description")

	response, err = client.DoTeaRequest("POST", "ascm", "2019-05-10", action, "", nil, nil, request)

	if err != nil {
		return err
	}
	if fmt.Sprintf(`%v`, response["code"]) != "200" {
		return errmsgs.WrapError(errmsgs.Error("CreateUpdateOdpsUser failed for "))
	}

	return resourceAlibabacloudStackMaxcomputeUserRead(d, meta)
}

func resourceAlibabacloudStackMaxcomputeUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	maxcomputeService := MaxcomputeService{client}
	object, err := maxcomputeService.DescribeMaxcomputeUser(d.Get("user_name").(string))
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_maxcompute_project_user maxcomputeService.DescribeMaxcomputeUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
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

func resourceAlibabacloudStackMaxcomputeUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

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
		roleId, err := client.RoleIds()
		if err != nil {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ASCM User", "defaultRoleId")), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}

		OrganizationId := ""
		if v, ok := d.GetOk("organization_id"); ok {
			OrganizationId = strconv.Itoa(v.(int))
		} else {
			OrganizationId = client.Department
		}

		commonRequest := client.NewCommonRequest("POST", "ascm", "2019-05-10", "UpdateOdpsUser", "/ascm/manage/resource_mgmt/updateOdpsUser")
		mergeMaps(commonRequest.QueryParams, map[string]string{
			"Id":               d.Get("id").(string),
			"UserId":           d.Get("user_id").(string),
			"UserName":         d.Get("user_name").(string),
			"UserType":         d.Get("user_type").(string),
			"OrganizationId":   OrganizationId,
			"OrganizationName": d.Get("organization_name").(string),
			"Description":      d.Get("description").(string),
		})
		commonRequest.Headers["x-acs-roleid"] = strconv.Itoa(roleId)

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(commonRequest)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"Error OdpsUser Not Found"}) {
				return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
			}
			errmsg := ""
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Get("user_name").(string), "UpdateOdpsUser", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug("UpdateOdpsUser", raw, requestInfo, commonRequest)
	}
	return resourceAlibabacloudStackMaxcomputeUserRead(d, meta)
}

func resourceAlibabacloudStackMaxcomputeUserDelete(d *schema.ResourceData, meta interface{}) error {
	// ASCM does not support deletion
	return nil

	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteOdpsCu"
	request := make(map[string]interface{})
	request["CuId"] = d.Id()
	request["CuName"] = d.Get("cu_name")
	request["ClusterName"] = d.Get("cluster_name")

	response, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", action, "", nil, nil, request)

	if err != nil {
		return err
	}
	if errmsgs.IsExpectedErrorCodes(fmt.Sprintf("%v", response["code"]), []string{"102", "403"}) {
		return nil
	}
	if fmt.Sprintf(`%v`, response["code"]) != "200" {
		return errmsgs.WrapError(errmsgs.Error("DeleteOdpsCu failed for " + response["Message"].(string)))
	}
	return nil
}
