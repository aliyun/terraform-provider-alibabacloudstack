---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_bgp_group"
description: |-
  Provides a expressconnect Bgpgroup resource.
---

# alibabacloudstack\_expressconnect\_bgpgroup

Provides a expressconnect Bgpgroup resource.

## Example Usage
```
resource "alibabacloudstack_expressconnect_bgp_group" "default" {
	bgp_group_name = "${var.name}"
	description =    "${var.name}"
	local_asn =      "65534"
	peer_asn =       "10"
	router_id =      "${alibabacloudstack_expressconnect_virtualborderrouter.default.id}$"
}

```

## Argument Reference

The following arguments are supported:
  * `auth_key` - (Optional) - The key used by the BGP group.
  * `bgp_group_id` - (Optional) - The ID of the BGP group.
  * `bgp_group_name` - (Optional) - The name of the BGP group.
  * `description` - (Optional) - Description of the BGP group.
  * `ip_version` - (Optional) - IP version
  * `local_asn` - (Optional) - This end AS number.
  * `peer_asn` - (Required) - The AS number of the side equipment.
  * `region_id` - (Optional) - The Region ID of the BGP group. For the list of Region IDs, see [Region and Zone](~~ 40654 ~~).
  * `router_id` - (Required) - ID of VBR.
  * `status` - (Optional) - The status of the resource

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `bgp_group_id` - The ID of the BGP group.
  * `hold` - The hold time to wait for the incoming BGP message. If no message has been passed in after the hold time, the BGP neighbor is considered disconnected.
  * `ip_version` - IP version
  * `is_fake` - Whether the AS number is false.
  * `keepalive` - Live time.
  * `region_id` - The Region ID of the BGP group. For the list of Region IDs, see [Region and Zone](~~ 40654 ~~).
  * `route_limit` - Routing restrictions.
  * `status` - The status of the resource
