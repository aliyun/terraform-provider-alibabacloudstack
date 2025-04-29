---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_domainextensions"
sidebar_current: "docs-Alibabacloudstack-datasource-slb-domainextensions"
description: |- 
  Provides a list of slb domainextensions owned by an alibabacloudstack account.
---

# alibabacloudstack_slb_domainextensions
-> **NOTE:** Alias name has: `alibabacloudstack_slb_domain_extensions`

This data source provides a list of SLB domain extensions in an AlibabacloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_slb_domainextensions" "example" {
  ids               = ["de-12345678"]
  load_balancer_id  = "lb-abc12345"
  frontend_port     = 443
}

output "slb_domain_extension_ids" {
  value = data.alibabacloudstack_slb_domainextensions.example.extensions[*].id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of SLB domain extension IDs. The query result will be filtered based on this parameter.
* `load_balancer_id` - (Required) The ID of the Server Load Balancer instance. This is used to filter domain extensions associated with the specified SLB instance.
* `frontend_port` - (Required) The frontend port used by the HTTPS listener of the SLB instance. Valid values range from 1 to 65535. This is used to filter domain extensions for the specified listener port.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `extensions` - A list of SLB domain extensions. Each element contains the following attributes:
  * `id` - The unique ID of the domain extension.
  * `domain` - The domain name associated with the domain extension.
  * `server_certificate_id` - The ID of the server certificate used by the domain name. This certificate is typically used for SSL/TLS encryption when serving HTTPS traffic.
