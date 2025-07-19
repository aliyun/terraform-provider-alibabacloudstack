package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type MongoDBShardingNodeType string

const (
	MongoDBShardingNodeMongos = MongoDBShardingNodeType("mongos")
	MongoDBShardingNodeShard  = MongoDBShardingNodeType("shard")
)

func auditFilterSchema(availableRoleType []string) *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role_type": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(availableRoleType, false),
				},
				"filters": {
					Type: schema.TypeSet,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validation.StringInSlice([]string{"admin", "slow", "query", "insert", "update", "delete", "command"}, false),
					},
					Required: true,
				},
			},
		},
		Optional: true,
		Computed: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			if d.Get("audit_status").(string) != "Enable" {
				return true
			}
			return old == new
		},
	}
}
