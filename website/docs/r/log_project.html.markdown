---
subcategory: "Simple Log Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_log_project"
sidebar_current: "docs-alibabacloudstack-resource-log-project"
description: |-
  Provides a Alibabacloudstack log project resource.
---

# alibabacloudstack_log_project

The project is the resource management unit in Log Service and is used to isolate and control resources.
You can manage all the logs, and the related log sources of an application by using projects.

## Example Usage

Basic Usage
To invoke this resource, you need to set the provider parameter: sls_openapi_endpoint
```hcl
provider "alibabacloudstack" {
  sls_openapi_endpoint = "var.sls_openapi_endpoint"
  ...
}

resource "alibabacloudstack_log_project" "example" {
  name        = "tf-log"
  description = "created by terraform"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the log project. It must be unique within your Alibaba Cloud account. Changing this forces a new resource to be created.
* `description` - (Optional) The description of the log project.
* `cluster_name` - (Optional, Computed) The name of the cluster where the log project is located. This attribute is returned by the API and cannot be manually set.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log project. It is the same as the `name`.
* `name` - The name of the log project.
* `description` - The description of the log project.
* `cluster_name` - The name of the cluster where the log project is located.

## Import

Log Project can be imported using the name, e.g.

```
$ terraform import alibabacloudstack_log_project.example my-log-project
```