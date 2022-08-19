---
subcategory: "Server Load Balancer (SLB)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_slb_backend_servers"
sidebar_current: "docs-apsarastack-datasource-slb-backend_servers"
description: |-
    Provides a list of server load balancer backend servers to the user.
---

# apsarastack\_slb_backend_servers

This data source provides the server load balancer backend servers related to a server load balancer..

## Example Usage

```
data "apsarastack_slb_beckend_servers" "sample_ds" {
  load_balancer_id = "${apsarastack_slb.sample_slb.id}"
}

output "first_slb_backend_server_id" {
  value = "${data.apsarastack_slb_beckend_servers.sample_ds.backend_servers.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required) ID of the SLB with attachments.
* `ids` - (Optional) List of attached ECS instance IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `backend_servers` - 
  * `id` - backend server ID.
  * `weight` - Weight associated to the ECS instance.

