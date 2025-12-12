---
subcategory: "Prometheus"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_prometheus_v2_contact"
sidebar_current: "docs-Alibabacloudstack-datasource-prometheus_v2_contact"
description: |-
  Queries Prometheus v2 contact information of Alibaba Cloud
---

# alibabacloudstack_prometheus_v2_contact

Queries Prometheus v2 contact information of Alibaba Cloud, used to obtain the list of created contacts.

## Example Usage

```hcl
variable "name" {
  default = "tfacc_prometheus31218"
}
resource "alibabacloudstack_prometheus_v2_contact" "default" {
  username = "tfacc-${var.name}"
  mobile   = "13812345678"
  mail     = "test@example.com"
}

data "alibabacloudstack_prometheus_v2_contacts" "default" {
  name_regex = alibabacloudstack_prometheus_v2_contact.default.username
}
```

## Argument Reference

The following arguments are supported for filtering query results:

* `ids` (Optional): A list of contact IDs, used to precisely match contacts with specified IDs.
* `name_regex` (Optional): A regular expression for username, used to fuzzy match contact names.

## Attributes Reference

The following attributes are exported:

* `id` (String): The unique identifier of the contact.
* `groups` (List): A list of group IDs to which the contact belongs.
* `mail` (String): The email address of the contact.
* `mobile` (String): The mobile phone number of the contact.
* `username` (String): The username of the contact.