package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDataWorksUser() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_code": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"role_project_owner", "role_project_admin", "role_project_dev", "role_project_pe", "role_project_deploy", "role_project_guest", "role_project_security"}, false),
				},
				MinItems: 1,
			},			
			"project_member_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackDataWorksUserCreate, resourceAlibabacloudStackDataWorksUserRead, resourceAlibabacloudStackDataWorksUserUpdate, resourceAlibabacloudStackDataWorksUserDelete)
	return resource
}

func convertAscmUid2MemberUid(uid string) string {
	if uid[0] != '2' {
		return uid
	}
	return "5" + uid
}

func resourceAlibabacloudStackDataWorksUserCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "CreateProjectMember"
	request := map[string]interface{}{
		"ProjectId": d.Get("project_id"),
		"UserId":    convertAscmUid2MemberUid(d.Get("user_id").(string)),
	}

	_, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_data_works_folder", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", d.Get("project_id"), d.Get("user_id")))

	return nil
}

func resourceAlibabacloudStackDataWorksUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dataworksPublicService := DataworksService{client}
	object, err := dataworksPublicService.DescribeDataWorksUser(d.Id())
	log.Printf(fmt.Sprint(object))
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_data_works_folder dataworksPublicService.DescribeDataWorksUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("user_id", parts[1])
	d.Set("project_id", parts[0])
	d.Set("project_member_id", object["ProjectMemberId"])
	roleCode := []string{}
	for _, i := range object["ProjectRoleList"].([]interface{}) {
		roleCode = append(roleCode, i.(map[string]interface{})["ProjectRoleCode"].(string))
	}
	d.Set("role_code", roleCode)

	return nil
}

func resourceAlibabacloudStackDataWorksUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if d.HasChange("role_code") {
		o, n := d.GetChange("role_code")
		oldMap := map[string]struct{}{}
		for _, old := range o.(*schema.Set).List() {
			oldMap[old.(string)] = struct{}{}
		}
		newMap := map[string]struct{}{}
		for _, new := range n.(*schema.Set).List() {
			newMap[new.(string)] = struct{}{}
		}
		for _, v := range o.(*schema.Set).List() {
			old := v.(string)
			if _, existed := newMap[old]; !existed {
				action := "RemoveProjectMemberFromRole"
				request := map[string]interface{}{
					"ProjectId": parts[0],
					"UserId":    convertAscmUid2MemberUid(parts[1]),
					"RoleCode":  old,
				}

				_, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request)
				if err != nil {
					if e, ok := err.(*errmsgs.ComplexError); ok {
						err = e.Cause
					}
					if e, ok := err.(*tea.SDKError); ok {
						if strings.Contains(*e.Message, "user aleady in role.") {
							continue
						}
					}
					return err
				}
			}
		}
		for _, v := range n.(*schema.Set).List() {
			new := v.(string)
			if _, existed := oldMap[new]; !existed {
				action := "AddProjectMemberToRole"
				request := map[string]interface{}{
					"ProjectId": parts[0],
					"UserId":    convertAscmUid2MemberUid(parts[1]),
					"RoleCode":  new,
				}

				_, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func resourceAlibabacloudStackDataWorksUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	action := "DeleteProjectMember"
	request := map[string]interface{}{
		"ProjectId": parts[0],
		"UserId":    convertAscmUid2MemberUid(parts[1]),
	}

	_, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
