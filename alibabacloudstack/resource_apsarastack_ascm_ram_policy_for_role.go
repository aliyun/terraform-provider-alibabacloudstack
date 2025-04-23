package alibabacloudstack

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackAscmRamPolicyForRole() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ram_policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAscmRamPolicyForRoleCreate, resourceAlibabacloudStackAscmRamPolicyForRoleRead, nil, resourceAlibabacloudStackAscmRamPolicyForRoleDelete)
	return resource
}

func resourceAlibabacloudStackAscmRamPolicyForRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ram_id := d.Get("ram_policy_id").(string)
	roleid := d.Get("role_id").(int)

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "AddRAMPolicyToRole", "/ascm/auth/role/addRAMPolicyToRole")
	request.QueryParams["RamPolicyId"] = ram_id
	request.QueryParams["RoleId"] = fmt.Sprint(roleid)

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ram_policy_for_role", "AddRAMPolicyToRole", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug("AddRAMPolicyToRole", raw, request)
	if bresponse.GetHttpStatus() != 200 {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ram_policy_for_role", "AddRAMPolicyToRole", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("AddRAMPolicyToRole", raw, request, bresponse.GetHttpContentString())

	d.SetId(ram_id + COLON_SEPARATED + fmt.Sprint(roleid))

	return nil
}

func resourceAlibabacloudStackAscmRamPolicyForRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	_, err := ascmService.DescribeAscmRamPolicyForRole(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("ram_policy_id", did[0])
	role_id, _ := strconv.Atoi(did[1])
	d.Set("role_id", role_id)

	return nil
}

func resourceAlibabacloudStackAscmRamPolicyForRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	did := strings.Split(d.Id(), COLON_SEPARATED)

	check, err := ascmService.DescribeAscmRamPolicyForRole(d.Id())

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsBindingExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsBindingExist", check, map[string]string{"ramPolicyId": did[0]})

	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveRAMPolicyFromRole", "/ascm/auth/role/removeRAMPolicyFromRole")
		request.QueryParams["ramPolicyId"] = did[0]
		request.QueryParams["roleId"] = did[1]

		raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "RemoveRAMPolicyFromRole", errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		check, err = ascmService.DescribeAscmRamPolicyForRole(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "RemoveRAMPolicyFromRole", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}