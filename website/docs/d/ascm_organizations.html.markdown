---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_organizations"
sidebar_current: "docs-apsarastack-datasource-ascm-organizations"
description: |-
    Provides a list of organizations to the user.
---

# apsarastack\_ascm_organizations

This data source provides the organizations of the current Apsara Stack Cloud user.

## Example Usage

```
resource "apsarastack_ascm_organization" "default" {
  name = "Test_org"
}
output "orgres" {
  value = apsarastack_ascm_organization.default.*
}
data "apsarastack_ascm_organizations" "default" {
    name_regex = apsarastack_ascm_organization.default.name
    parent_id = apsarastack_ascm_organization.default.parent_id
}
output "orgs" {
  value = data.apsarastack_ascm_organizations.default.*
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of organizations IDs.
* `name_regex` - (Optional) A regex string to filter results by organization name.
* `parent_id` - (Optional) Filter the results by the specified organization parent ID.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `organizations` - A list of organizations. Each element contains the following attributes:
  * `id` - ID of the organization.
  * `name` - organization name.
  * `cuser_id` - Id of a Cuser.
  * `muser_id` - Id of a Muser.
  * `alias` - alias for the Organization.
  * `parent_id` - Parent id of an Organization.
 
