---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_bgpgroups"
sidebar_current: "docs-Alibabacloudstack-datasource-expressconnect-bgpgroups"
description: |-
  Provides a list of expressconnect bgpgroups owned by an alibabacloudstack account.
---

# alibabacloudstack\_expressconnect\_bgpgroups

This data source provides a list of expressconnect bgpgroups in an alibabacloudstack account according to the specified filters.

## Example Usage
```
resource "alibabacloudstack_expressconnect_bgp_group" "default" {
	bgp_group_name = "${var.name}"
	description =    "${var.name}"
	local_asn =      "65534"
	peer_asn =       "10"
	router_id =      "${alibabacloudstack_expressconnect_virtualborderrouter.default.id}$"
}

data "alibabacloudstack_expressconnect_bgp_groups" "default" {
	router_id = "${alibabacloudstack_expressconnect_virtualborderrouter.default.id}"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - A list of BGP group IDs.
  * `router_id` - (Required) - ID of VBR.
  * `region_id` - (Optional) - The Region ID of the BGP group. For the list of Region IDs, see [Region and Zone](~~ 40654 ~~).
  * `name_regex` - (Optional) - A regex string to filter results by BGP group name.
  * `description_regex` - (Optional) - A regex string to filter results by BGP group description.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `bgp_groups` - A list of BGP groups.
    * `id` - The ID of the BGP group.
    * `auth_key` - The key used by the BGP group.
    * `bgp_group_id` - The ID of the BGP group.
    * `bgp_group_name` - The name of the BGP group.
    * `description` - Description of the BGP group.
    * `hold` - The hold time to wait for the incoming BGP message. If no message has been passed in after the hold time, the BGP neighbor is considered disconnected.
    * `ip_version` - IP version
    * `is_fake` - Whether the AS number is false.
    * `keepalive` - Live time.
    * `local_asn` - This end AS number.
    * `peer_asn` - The AS number of the side equipment.
    * `route_limit` - Routing restrictions.
    * `router_id` - ID of VBR.
    * `status` - The status of the resource
