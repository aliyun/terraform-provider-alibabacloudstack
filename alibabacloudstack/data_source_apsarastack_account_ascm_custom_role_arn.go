package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceAlibabacloudStackAccountAscmCustomRoleArn exposes the RAM role that the ASCM control
// plane materializes for a role held by the account the provider is authenticated with.
//
// An `ROLETYPE_ASCM` role cannot be assumed through STS directly. Whenever such a role is granted
// to a user, ASCM registers a companion RAM role inside that user's own cloud account, named
// after the binding it stands for:
//
//	acs:ram::<account uid of the organization>:role/ascm-role-<organization id>-<ascm role id>-<ascm user id>
//
// That companion role is what `alibabacloudstack_ascm_user_role` style consumers need in order to
// hand a workload credentials, and its `assume_role_policy_document` states which principal is
// allowed to assume it.
//
// The organization and the ASCM user are not accepted as arguments: GetRamRoleAndRamPolicy
// resolves both from the caller's identity, and refuses to disclose the role of another user.
// The data source therefore always answers the question "which RAM role does the account I am
// logging in with hold for this ASCM role".
func dataSourceAlibabacloudStackAccountAscmCustomRoleArn() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAccountAscmCustomRoleArnRead,

		Schema: map[string]*schema.Schema{
			"ascm_role_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ram_role_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"assume_role_policy_document": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlibabacloudStackAccountAscmCustomRoleArnRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}

	ascmRoleId := d.Get("ascm_role_id").(int)

	role, err := ascmService.GetRamRoleAndRamPolicy(ascmRoleId)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	if role.Arn == "" {
		return errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("no RAM role ARN is materialized for ascm_role_id %d of the authenticated account", ascmRoleId))
	}

	d.SetId(role.Arn)
	if err := d.Set("arn", role.Arn); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("role_name", role.RoleName); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ram_role_id", role.RoleID); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("assume_role_policy_document", role.AssumeRolePolicyDocument); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
