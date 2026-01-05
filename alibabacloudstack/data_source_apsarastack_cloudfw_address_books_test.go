package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func resourceCloudfwAddressBookConfigDependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-addrbook-%d"
}

resource "alibabacloudstack_cloudfw_address_book" "default" {
    group_type = "ip"
    group_name = var.name
    address_list = ["100.100.100.100/30"]
    description = "test address book"
}

data "alibabacloudstack_cloudfw_address_books" "default" {
    %s
}
`, rand, strings.Join(pairs, "\n    "))
}

func TestAccAlibabacloudStackCloudfwAddressBooksDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_cloudfw_address_books.default"
	rand := getAccTestRandInt(10000, 20000)
	testDataSourceAttr := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"address_books.#":                    "1",
				"address_books.0.group_name":         fmt.Sprintf("tf-testacc-addrbook-%d", rand),
				"address_books.0.group_type":         "ip",
				"address_books.0.address_list.0":     "100.100.100.100/30",
				"address_books.0.description":        "test address book",
				"address_books.0.address_list_count": "1",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"address_books.#": "0",
			}
		},
		Providers: testYunDunProviders(),
	}

	testDataSourceAttr.dataSourceTestCheck(t, rand,
		// Test query by ids
		dataSourceTestAccConfig{
			existConfig: resourceCloudfwAddressBookConfigDependenceNew(rand, map[string]string{
				"group_type": `"ip"`,
				"ids":        `["${alibabacloudstack_cloudfw_address_book.default.id}"]`,
			}),
			fakeConfig: resourceCloudfwAddressBookConfigDependenceNew(rand, map[string]string{
				"ids": `["${alibabacloudstack_cloudfw_address_book.default.id}_fake"]`,
			}),
		},
		// Test query by name_regex
		dataSourceTestAccConfig{
			existConfig: resourceCloudfwAddressBookConfigDependenceNew(rand, map[string]string{
				"group_type": `"ip"`,
				"name_regex": `"${alibabacloudstack_cloudfw_address_book.default.group_name}"`,
			}),
			fakeConfig: resourceCloudfwAddressBookConfigDependenceNew(rand, map[string]string{
				"group_type": `"ip"`,
				"name_regex": `"${alibabacloudstack_cloudfw_address_book.default.group_name}_fake"`,
			}),
		},
		// Test query by group_type
		dataSourceTestAccConfig{
			existConfig: resourceCloudfwAddressBookConfigDependenceNew(rand, map[string]string{
				"group_type": `"ip"`,
			}),
			fakeConfig: resourceCloudfwAddressBookConfigDependenceNew(rand, map[string]string{
				"group_type": `"port"`,
			}),
		},
		// Test query by query keyword
		dataSourceTestAccConfig{
			existConfig: resourceCloudfwAddressBookConfigDependenceNew(rand, map[string]string{
				"group_type": `"ip"`,
				"query":      `"${alibabacloudstack_cloudfw_address_book.default.group_name}"`,
			}),
			fakeConfig: resourceCloudfwAddressBookConfigDependenceNew(rand, map[string]string{
				"group_type": `"ip"`,
				"query":      `"fake-${alibabacloudstack_cloudfw_address_book.default.group_name}"`,
			}),
		},
	)
}
