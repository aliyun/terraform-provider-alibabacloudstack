---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_bgp_networks"
description: |-
  Provides a list of expressconnect bgpnetworks owned by an alibabacloudstack account.
---

# alibabacloudstack\_expressconnect\_bgpnetworks

This data source provides a list of expressconnect bgpnetworks in an alibabacloudstack account according to the specified filters.

## Example Usage
```
data "alibabacloudstack_expressconnect_bgp_networks" "example" {
  router_id = "vbr-bp1d8yixxxxxxxxxxx"
  dst_cidr_block = "192.168.0.0/16"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - The ids of the expressconnect bgpnetworks.
  * `router_id` - (Required, ForceNew) - The ID of the router.
  * `dst_cidr_block` - (Optional) - Declared BGP network.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `bgp_networks` - The list of expressconnect bgpnetworks.
    * `id` - The ID of the expressconnect bgpnetwork.
    * `dst_cidr_block` - Declared BGP network.
    * `router_id` - The ID of the router.
    * `status` - Declared BGP network status.
