---
subcategory: "Log Service (SLS)"
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
```
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

* `name` - (Required, ForceNew) The name of the log project. It is the only in one Alibabacloudstack account.
* `description` - (Optional) Description of the log project.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log project. It same as its name.
* `name` - Log project name.
* `description` - Log project description.