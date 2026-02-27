package alibabacloudstack

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackAscmUserRoleBinding() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"login_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_ids": {
				Type:          schema.TypeSet,
				Computed:      true,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeInt},
				MinItems:      1,
				MaxItems:      10,
				ConflictsWith: []string{"role_id"},
				Deprecated:    "Starting from version 3.21.0, the role_ids parameter is no longer supported for configuration.",
			},
			"role_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"role_ids"},
			},
		},
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			if d.Id() == "" && d.Get("role_ids").(*schema.Set).Len() > 0 {
				log.Printf("[WARN] Resource created with 'role_ids' will not support deletion. " +
					"Consider using 'role_id' or 'ascm_user' instead.")
			}
			return nil
		},
		DeprecationMessage: "ascm_user already includes corresponding functions. This resource may be removed in future versions.",
	}

	setResourceFunc(resource, resourceAlibabacloudStackAscmUserRoleBindingCreate, resourceAlibabacloudStackAscmUserRoleBindingRead, resourceAlibabacloudStackAscmUserRoleBindingUpdate, resourceAlibabacloudStackAscmUserRoleBindingDelete)
	return resource
}

func resourceAlibabacloudStackAscmUserRoleBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	lname := d.Get("login_name").(string)
	var roleids []int
	if v, ok := d.GetOk("role_ids"); ok {
		for _, id := range v.(*schema.Set).List() {
			roleids = append(roleids, id.(int))
		}
	} else if v, ok := d.GetOk("role_id"); ok {
		roleids = append(roleids, v.(int))
	}
	log.Printf("roleids is %v", roleids)
	for _, roleId := range roleids {
		if err := ascmService.AscmUserRoleBinding(lname, roleId); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("role_id"); ok {
		d.SetId(fmt.Sprintf("%s:%d", lname, v.(int)))
	} else {
		d.SetId(lname)
	}
	return nil
}

func resourceAlibabacloudStackAscmUserRoleBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmUserRoleBinding(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	if len(object.Data) == 0 {
		d.SetId("")
		return nil
	}

	loginName := object.Data[0].LoginName
	d.Set("login_name", loginName)

	roleIds := make([]int, 0)
	for _, role := range object.Data[0].Roles {
		roleIds = append(roleIds, role.ID)
	}
	parts := strings.Split(d.Id(), ":")
	if len(parts) == 2 {
		roleId, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}
		d.Set("role_id", roleId)
	}
	d.Set("role_ids", roleIds)
	return nil
}

func resourceAlibabacloudStackAscmUserRoleBindingUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}

	lname := d.Get("login_name").(string)

	if d.HasChange("role_ids") {
		o, n := d.GetChange("role_ids")
		oldValue := make(map[int]struct{})
		newValue := make(map[int]struct{})
		for _, v := range o.(*schema.Set).List() {
			oldValue[v.(int)] = struct{}{}
		}
		if len(oldValue) == 0 {
			if object, err := ascmService.DescribeAscmUserRoleBinding(d.Id()); err == nil && len(object.Data) > 0 {
				for _, role := range object.Data[0].Roles {
					oldValue[role.ID] = struct{}{}
				}
			}
		}
		for _, v := range n.(*schema.Set).List() {
			newValue[v.(int)] = struct{}{}
		}
		toRemove := []int{}
		toAdd := []int{}
		for key := range newValue {
			if _, exist := oldValue[key]; !exist {
				toAdd = append(toAdd, key)
			}
		}
		for key := range oldValue {
			if _, exist := newValue[key]; !exist {
				toRemove = append(toRemove, key)
			}
		}
		if len(oldValue) == 10 && len(toAdd) > 0 && len(toRemove) > 0 {
			if err := ascmService.AscmUserRoleUnbinding(lname, toRemove[0]); err != nil {
				return err
			}
			toRemove = toRemove[1:]
		}
		minLen := min(len(toAdd), len(toRemove))
		for i := 0; i < minLen; i++ {
			if err := ascmService.AscmUserRoleBinding(lname, toAdd[i]); err != nil {
				return err
			}
			if err := ascmService.AscmUserRoleUnbinding(lname, toRemove[i]); err != nil {
				return err
			}
		}
		for i := minLen; i < len(toAdd); i++ {
			if err := ascmService.AscmUserRoleBinding(lname, toAdd[i]); err != nil {
				return err
			}
		}
		for i := minLen; i < len(toRemove); i++ {
			if err := ascmService.AscmUserRoleUnbinding(lname, toRemove[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func resourceAlibabacloudStackAscmUserRoleBindingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	parts := strings.Split(d.Id(), ":")

	if len(parts) == 2 {
		roleId, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("invalid role ID in resource ID: %s", d.Id())
		}
		return ascmService.AscmUserRoleUnbinding(parts[0], roleId)
	}
	return nil
}

func (s *AscmService) AscmUserRoleBinding(userName string, roleId int) error {
	requestBody := map[string]interface{}{
		"loginName": userName,
		"roleId":    roleId,
	}
	if _, err := s.client.DoTeaRequest("POST", "ascm", "2019-05-10", "AddRoleToUser", "/ascm/auth/role/addRoleToUser", nil, nil, requestBody); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_role_binding", "AddRoleToUser", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}

func (s *AscmService) AscmUserRoleUnbinding(userName string, roleId int) error {
	requestBody := map[string]interface{}{
		"loginName": userName,
		"roleId":    roleId,
	}
	if _, err := s.client.DoTeaRequest("POST", "ascm", "2019-05-10", "RemoveRoleFromUser", "/ascm/auth/role/removeRoleFromUser", nil, nil, requestBody); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_role_binding", "RemoveRoleFromUser", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
