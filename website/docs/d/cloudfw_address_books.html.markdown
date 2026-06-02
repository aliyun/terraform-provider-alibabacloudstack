---
subcategory: "Cloud Firewall"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudfw_address_books"
description: |-
  Provides information about Cloud Firewall address books
---

# alibabacloudstack_cloudfw_address_books

This data source provides information about Cloud Firewall address books that have been created.

## Example Usage

```hcl

variable "name" {
  default = "tf-testacc-addrbook-11865"
}

resource "alibabacloudstack_cloudfw_address_book" "default" {
  group_type   = "ip"
  group_name   = var.name
  address_list = ["100.100.100.100/30"]
  description  = "test address book"
}

data "alibabacloudstack_cloudfw_address_books" "default" {
  group_type = "ip"
  ids        = ["${alibabacloudstack_cloudfw_address_book.default.id}"]
}

```

## Argument Reference

The following arguments support filtering query results:

* `contain_port` (Optional) - Filter address books containing specific ports.
* `group_type` (Optional) - The type of address book, valid values are `ip` or `port`.
* `ids` (Optional) - A list of address book UUIDs for filtering results.
* `name_regex` (Optional) - Filter address book names using regular expressions.
* `query` (Optional) - Query address books by name keywords.

## Attributes Reference

The following attributes are exported:

* `id` (String) - The address book ID, formatted as `{group_type}:{group_uuid}`.
* `address_list` (List) - A list of IP addresses or ports in the address book.
* `address_list_count` (Integer) - The number of entries in the address list.
* `auto_add_tag_ecs` (Integer) - Whether to automatically add tagged ECS instances (0: No, 1: Yes).
* `description` (String) - The description of the address book.
* `global` (Integer) - Whether it is a global address book (0: No, 1: Yes).
* `group_name` (String) - The name of the address book.
* `group_type` (String) - The type of address book (ip or port).
* `group_uuid` (String) - The unique identifier UUID of the address book.
* `reference_count` (Integer) - The number of times the address book is referenced.
* `tag_relation` (String) - Tag relationship configuration information.