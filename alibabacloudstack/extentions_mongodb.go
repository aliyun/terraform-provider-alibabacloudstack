package alibabacloudstack

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type MongoDBShardingNodeType string

const (
	MongoDBShardingNodeMongos = MongoDBShardingNodeType("mongos")
	MongoDBShardingNodeShard  = MongoDBShardingNodeType("shard")
	MongoDBShardingNodeCs     = MongoDBShardingNodeType("configserver")
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
		Set: func(i interface{}) int {
			m := i.(map[string]interface{})
			filters := []string{}
			for _, f := range m["filters"].(*schema.Set).List() {
				filters = append(filters, f.(string))
			}
			sort.Strings(filters)
			hashString := fmt.Sprintf("%s@%s", m["role_type"].(string), strings.Join(filters, ","))
			return hashcode.String(hashString)
		},
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			if d.Get("audit_status").(string) != "Enable" {
				return true
			}
			return old == new
		},
	}
}
