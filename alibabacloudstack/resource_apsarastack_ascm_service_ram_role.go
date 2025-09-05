package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackAscmServiceRamRole() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ram_roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"product_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aliyun_user_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"assume_role_policy_document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policies": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"region": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"policy_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"policy_document": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"policy_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"default_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"aliyun_user_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ram_group_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ascm_ram_policy_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"attach_date": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"resource_set_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"privilege_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ram_role_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAscmServiceRamRoleCreate,
		resourceAlibabacloudStackAscmServiceRamRoleRead, nil, resourceAlibabacloudStackAscmServiceRamRoleDelete)
	return resource
}

func resourceAlibabacloudStackAscmServiceRamRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	product_name := d.Get("product_name").(string)
	organization_id := d.Get("organization_id").(string)
	id := fmt.Sprintf("%s:%s", organization_id, product_name)
	ascmService := AscmService{client}
	object, _ := ascmService.ListRAMServiceRoles(id)
	if len(object.Data) > 0 {
		d.SetId(id)
		return nil
	}
	request := make(map[string]interface{})
	request["productName"] = product_name
	request["organizationIdList"] = []string{organization_id}
	_, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "BatchCreateRAMServiceRole", "/ascm/auth/ramServiceRole/batchCreateRAMServiceRole", nil, nil, request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_ram_service_role", "BatchCreateRAMServiceRole", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	d.SetId(id)
	return nil
}

func resourceAlibabacloudStackAscmServiceRamRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.ListRAMServiceRoles(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	var ram_roles []map[string]interface{}
	for _, v := range object.Data {
		// Set policies list
		var policies []map[string]interface{}
		if len(v.Policies) > 0 {
			for _, policy := range v.Policies {
				policies = append(policies, map[string]interface{}{
					"policy_id":          policy.ID,
					"region":             policy.Region,
					"policy_name":        policy.PolicyName,
					"description":        policy.Description,
					"policy_document":    policy.PolicyDocument,
					"policy_type":        policy.PolicyType,
					"default_version":    policy.DefaultVersion,
					"aliyun_user_id":     policy.AliyunUserId,
					"ram_group_id":       policy.RamGroupId,
					"ascm_ram_policy_id": policy.AscmRamPolicyId,
					"attach_date":        policy.AttachDate,
					"resource_set_id":    policy.ResourceSetId,
					"privilege_id":       policy.PrivilegeId,
					"ram_role_id":        policy.RamRoleId,
				})
			}
		}
		ram_roles = append(ram_roles, map[string]interface{}{
			"id":                          v.ID,
			"arn":                         v.Arn,
			"region":                      v.Region,
			"product_name":                v.Product,
			"role_id":                     v.RoleId,
			"role_name":                   v.RoleName,
			"role_type":                   v.RoleType,
			"description":                 v.Description,
			"aliyun_user_id":              v.AliyunUserId,
			"assume_role_policy_document": v.AssumeRolePolicyDocument,
			"organization_id":             v.OrganizationId,
			"organization_name":           v.OrganizationName,
			"policies":                    policies,
		})
	}
	d.Set("ram_roles", ram_roles)
	d.Set("organization_id", fmt.Sprint(object.Data[0].OrganizationId))
	d.Set("product_name", object.Data[0].Product)
	return nil
}

func resourceAlibabacloudStackAscmServiceRamRoleDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
