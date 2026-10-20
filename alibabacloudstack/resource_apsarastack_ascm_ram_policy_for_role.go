package alibabacloudstack

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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
	ascmService := AscmService{client}
	ram_id := d.Get("ram_policy_id").(string)
	roleid := d.Get("role_id").(int)

	// D1: AddRAMPolicyToRole only supports RAM role authorization (roleRange.rawRamRole).
	// ASCM custom roles (roleRange.allOrganizations/orgAndSubOrgs/resourceSet/userGroup) are rejected
	// with ascm.manage.InvalidParameter.RoleType.RamAuthorization.
	roleIdStr := fmt.Sprint(roleid)
	roleObj, err := ascmService.DescribeAscmRamRoleForRoleid(fmt.Sprintf("_:%d", roleid))
	if err == nil && roleObj.RoleRange != "roleRange.rawRamRole" {
		return fmt.Errorf("alibabacloudstack_ascm_ram_policy_for_role: AddRAMPolicyToRole only supports "+
			"RAM role authorization type (roleRange.rawRamRole), but role %s has roleRange=%s",
			roleIdStr, roleObj.RoleRange)
	}

	// API doc: POST /ascm/auth/role/addRAMPolicyToRole, Body: {"roleId": Integer, "ramPolicyId": Integer}
	ramPolicyIdInt, _ := strconv.Atoi(ram_id)
	requestBody := map[string]interface{}{
		"roleId":      roleid,
		"ramPolicyId": ramPolicyIdInt,
	}
	_, err = client.DoTeaRequest("POST", "ascm", "2019-05-10", "AddRAMPolicyToRole", "/ascm/auth/role/addRAMPolicyToRole", nil, nil, requestBody)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ram_policy_for_role", "AddRAMPolicyToRole", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.SetId(ram_id + COLON_SEPARATED + roleIdStr)

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
	addDebug("IsBindingExist", check, map[string]interface{}{"ramPolicyId": did[0], "roleId": did[1]})

	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		// API doc: POST /ascm/auth/role/removeRAMPolicyFromRole, Body: {"roleId": Integer, "ramPolicyId": Integer}
		roleIdInt, _ := strconv.Atoi(did[1])
		ramPolicyIdInt, _ := strconv.Atoi(did[0])
		requestBody := map[string]interface{}{
			"roleId":      roleIdInt,
			"ramPolicyId": ramPolicyIdInt,
		}
		_, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "RemoveRAMPolicyFromRole", "/ascm/auth/role/removeRAMPolicyFromRole", nil, nil, requestBody)
		if err != nil {
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ram_policy_for_role", "RemoveRAMPolicyFromRole", errmsgs.AlibabacloudStackSdkGoERROR))
		}
		check, err = ascmService.DescribeAscmRamPolicyForRole(d.Id())

		if err != nil {
			if errmsgs.NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "RemoveRAMPolicyFromRole", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
