---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_protocols"
sidebar_current: "docs-alibabacloudstack-datasource-nas-protocols"
description: |-
    Provides a list of FileType owned by an Alibaba Cloud account.
---

# alibabacloudstack\_nas_protocols

Provide  a data source to retrieve the type of protocol used to create NAS file system.


## Example Usage

```terraform
data "alibabacloudstack_nas_protocols" "default" {
  type        = "Performance"
  zone_id     = "cn-beijing-e"
  output_file = "protocols.txt"
}

output "nas_protocols_protocol" {
  value = "${data.alibabacloudstack_nas_protocols.default.protocols.0}"
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required) The file system type. Valid Values: `Performance` and `Capacity`.  
* `zone_id` - (Optional) String to filter results by zone id. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `protocols` - A list of supported protocol type..
