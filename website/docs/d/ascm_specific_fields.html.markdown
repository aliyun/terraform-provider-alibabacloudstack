---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_specific_fields"
sidebar_current: "docs-apsarastack-datasource-ascm-specific-fields"
description: |-
    Provides a list of specific fields to the user.
---

# apsarastack\_ascm_specific_fields

This data source provides the specific fields of the current Apsara Stack Cloud user.

## Example Usage

```
data "apsarastack_ascm_specific_fields" "specifields" {
  group_filed ="storageType"
  resource_type ="OSS"
  output_file = "fields"
}
output "specifields" {
  value = data.apsarastack_ascm_specific_fields.specifields.*
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of specific fields IDs.
* `group_filed` - (Required) The field for which to query valid values.
* `resource_type` - (Required) Filter the results by the specified resource type. Valid values: OSS, ADB, DRDS, SLB, NAT, MAXCOMPUTE, POSTGRESQL, ECS, RDS, IPSIX, REDIS, MONGODB, and HITSDB.
* `label` - (Optional) Specifies whether to internationalize the field. Valid values: true and false.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `specific_fields` - A list of specific fields.
