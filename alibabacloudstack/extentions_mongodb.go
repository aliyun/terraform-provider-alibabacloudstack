package alibabacloudstack

import (
	"fmt"
	"regexp"
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

var (
MongoDBChangingStatus = []string{"Creating","NodeCreating", "NodeDeleting", "DBInstanceClassChanging", "DBInstanceNetTypeChanging",
	"NET_DELETING", "NET_CREATING", "NET_MODIFYING", "TDEModifying", "CONFIG_SWITCHING"}
)

func validateShardeNodeDescription() schema.SchemaValidateFunc{
	regex := regexp.MustCompile(`^[\p{L}_-][\w\p{L}-]{1,255}$`)
	errmesg :="The value must be 2 to 256 characters in length, and can contain letters, digits, underscores (_), and hyphen (-). It must start with a letter, and cannot start with http:// or https://."
	return validation.StringMatch(regex, errmesg)
}

func validateShardeNodeConnectionString() schema.SchemaValidateFunc{
	regex := regexp.MustCompile(`^[a-z]([a-z0-9-]{6,62})[a-z0-9]$`)
	errmesg :="The connection string must be 8 to 64 characters in length and can contain lowercase letters, digits, and hyphens (-). It must start with a lowercase letter and end with a lowercase letter or digit."
	return validation.StringMatch(regex, errmesg)
}

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
