package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDrdsReadonlyInstance() *schema.Resource {
	resource := resourceAlibabacloudStackDrdsInstance()
	resource.Schema["master_instance_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	}

	return resource
}
