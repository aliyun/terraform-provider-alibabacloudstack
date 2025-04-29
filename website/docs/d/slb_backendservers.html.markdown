---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_backendservers"
sidebar_current: "docs-Alibabacloudstack-datasource-slb-backendservers"
description: |- 
  Provides a list of slb backendservers owned by an alibabacloudstack account.
---

# alibabacloudstack_slb_backendservers
-> **NOTE:** Alias name has: `alibabacloudstack_slb_backend_servers`

This data source provides a list of slb backendservers in an AlibabaCloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_slb_backendservers" "sample_ds" {
  load_balancer_id = "${alibabacloudstack_slb.sample_slb.id}"
}

output "first_slb_backend_server_id" {
  value = "${data.alibabacloudstack_slb_backendservers.sample_ds.backend_servers.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required) The ID of the traditional server load balancer instance.
* `ids` - (Optional) A list of ECS instance IDs that are attached as backend servers.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `backend_servers` - Information about each backend server. Each element contains:
  * `id` - The unique identifier of the backend server (ECS instance).
  * `weight` - The weight assigned to the backend server, which determines the proportion of traffic it receives.