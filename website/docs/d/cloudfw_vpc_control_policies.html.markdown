---
subcategory: "Cloud Firewall"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudfw_vpc_control_policies"
sidebar_current: "docs-Alibabacloudstack-datasource-cloudfw-vpc-control-policies"
description: |-
  Query VPC control policies of Cloud Firewall.
---

# alibabacloudstack_cloudfw_vpc_control_policies

This data source queries VPC control policies of Cloud Firewall to obtain information about configured VPC firewall access control rules.

## Example Usage

```hcl
variable "name" {
  default = "tf-testacc-vpccontrolpolicies-12994"
}

resource "alibabacloudstack_cloudfw_vpc_control_policy" "default" {
  destination      = "0.0.0.0/16"
  description      = var.name
  application_name = "ANY"
  source_type      = "net"
  dest_port        = "33/33"
  acl_action       = "log"
  destination_type = "net"

  source         = "0.0.0.0/16"
  dest_port_type = "port"
  proto          = "UDP"
  application_id = "0"
  release        = true
}

data "alibabacloudstack_cloudfw_vpc_control_policies" "default" {
  ids = ["${alibabacloudstack_cloudfw_vpc_control_policy.default.id}"]
}
```

## Argument Reference

The following arguments support filtering query results:

* `destination` (String, Optional): The destination IP address or address group, used to filter policies matching the destination address.
* `ids` (List, Optional): A list of policy IDs to filter specific policies. If not specified, all policies are returned.
* `name_regex` (String, Optional): A regular expression for the policy description, used to filter policies with matching descriptions.
* `proto` (String, Optional): The protocol type, such as "TCP" or "UDP", used to filter policies with the specified protocol.
* `source` (String, Optional): The source IP address or address group, used to filter policies matching the source address.
* `acl_action` (String, Optional): The action of the access control rule, such as "accept", "drop", or "log", used to filter policies with the specified action.

## Attributes Reference

The following attributes are exported:

* `id` (String): The unique identifier of the data source, generated based on a hash of the matched policy IDs.
* `ids` (List): A list of matched policy IDs.
* `names` (List): A list of matched policy descriptions.
* `policies` (List): A list of matched policies. Each policy contains the following attributes:
  * `acl_uuid` (String): The unique ID of the policy.
  * `application_id` (String): The ID of the associated application.
  * `application_name` (String): The name of the associated application.
  * `description` (String): The description of the policy.
  * `dest_port` (String): The destination port or port range.
  * `dest_port_group` (String): The name of the destination port group.
  * `dest_port_group_ports` (List): A list of ports included in the destination port group.
  * `dest_port_type` (String): The type of destination port, such as "port".
  * `destination` (String): The destination IP address or address group.
  * `destination_group_cidrs` (List): A list of destination CIDRs when the destination type is a group.
  * `destination_type` (String): The type of destination, such as "net".
  * `direction` (String): The traffic direction, such as "inout".
  * `hit_times` (Integer): The number of times the policy has been matched.
  * `order` (Integer): The order of the policy in the rule list.
  * `proto` (String): The protocol used, such as "TCP" or "UDP".
  * `release` (String): Whether the policy is published (enabled).
  * `source` (String): The source IP address or address group.
  * `source_group_cidrs` (List): A list of source CIDRs when the source type is a group.
  * `source_type` (String): The type of source, such as "net".