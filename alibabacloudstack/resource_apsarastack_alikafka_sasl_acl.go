package alibabacloudstack

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAlikafkaSaslAcl() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"acl_resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Group", "Topic"}, false),
			},
			"acl_resource_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"acl_resource_pattern_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"LITERAL", "PREFIXED"}, false),
			},
			"acl_operation_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Read", "Write"}, true),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool{
					return strings.EqualFold(oldValue, newValue)
				},
				DiffSuppressOnRefresh: true,
			},
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "*",
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAlikafkaSaslAclCreate, resourceAlibabacloudStackAlikafkaSaslAclRead, nil, resourceAlibabacloudStackAlikafkaSaslAclDelete)
	return resource
}

func resourceAlibabacloudStackAlikafkaSaslAclCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"AclOperationType":       d.Get("acl_operation_type"),
		"AclResourceName":        d.Get("acl_resource_name"),
		"AclResourcePatternType": d.Get("acl_resource_pattern_type"),
		"AclResourceType":        d.Get("acl_resource_type"),
		"InstanceId":             d.Get("instance_id"),
		"Username":               d.Get("username"),
		"Host":                   d.Get("host"),
	}
	_, err := client.DoTeaRequest("POST", "alikafka", "2019-09-16", "CreateAcl", "", nil, reqQuery, nil)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", d.Get("instance_id"), d.Get("username"), d.Get("acl_resource_type"), d.Get("acl_resource_name"), d.Get("acl_resource_pattern_type"), d.Get("acl_operation_type"), d.Get("host")))
	return nil
}

func resourceAlibabacloudStackAlikafkaSaslAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	object, err := alikafkaService.DescribeAlikafkaSaslAcl(d.Id())
	if err != nil {
		// Handle exceptions
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	
	parts, err := ParseResourceId(d.Id(), 7)

	d.Set("instance_id", parts[0])
	d.Set("username", object["username"])
	d.Set("acl_resource_type", object["aclResourceType"])
	d.Set("acl_resource_name", object["aclResourceName"])
	d.Set("acl_resource_pattern_type", object["aclResourcePatternType"])
	d.Set("acl_operation_type", object["aclOperationType"])
	d.Set("host", object["host"])

	return nil
}

func resourceAlibabacloudStackAlikafkaSaslAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts, err := ParseResourceId(d.Id(), 7)

	reqQuery := map[string]interface{}{
		"AclResourceType":        parts[2],
		"AclResourcePatternType": parts[4],
		"AclOperationType":       parts[5],
		"AclResourceName":        parts[3],
		"Username":               parts[1],
		"Host":                   parts[6],
		"InstanceId":             parts[0],
	}

	_, err = client.DoTeaRequest("GET", "alikafka", "2019-09-16", "DeleteAcl", "", nil, reqQuery, nil)

	return err

}
