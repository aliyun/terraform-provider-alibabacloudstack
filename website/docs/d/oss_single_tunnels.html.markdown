---
page_title: "alibabacloudstack_oss_single_tunnels"
subcategory: "Object Storage Service"
description: |-
  Provides a list of OSS single tunnels in Apsara Stack Cloud.
---

# alibabacloudstack_oss_single_tunnels

Provides a list of OSS single tunnels in Apsara Stack Cloud. These tunnels represent virtual IP addresses used for connecting to OSS services across VPCs.

## Example Usage

```hcl
# Declare the data source
data "alibabacloudstack_oss_single_tunnels" "tunnels" {
  ids        = ["cluster1:vpc-abc123:vip1", "cluster2:vpc-def456:vip2"]
  name_regex = "^test-.*"
}

output "first_tunnel_id" {
  value = data.alibabacloudstack_oss_single_tunnels.tunnels.tunnels.0.id
}

output "all_tunnel_ids" {
  value = [for tunnel in data.alibabacloudstack_oss_single_tunnels.tunnels.tunnels : tunnel.id]
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of tunnel IDs. Each ID is in the format `cluster:vpc_id:vip`.
* `name_regex` - (Optional) A regex string to filter results by label. Only tunnels with labels matching this regex will be returned.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `tunnels` - A list of tunnels. Each element contains the following attributes:
  * `id` - The ID of the tunnel, formatted as `cluster:vpc_id:vip`.
  * `cluster` - The cluster to which the tunnel belongs.
  * `label` - The label of the tunnel.
  * `vip` - The virtual IP address of the tunnel.
  * `vpc_id` - The ID of the VPC associated with the tunnel.
  * `shared` - Whether the tunnel is shared (value is 0 or 1).
