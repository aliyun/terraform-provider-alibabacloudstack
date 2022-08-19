---
subcategory: "Container Registry (CR)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_cr_repo"
sidebar_current: "docs-apsarastack-resource-container-registry"
description: |-
  Provides a Apsarastack resource to manage Container Registry repositories.
---

# apsarastack\_cr\_repo

This resource will help you to manager Container Registry repositories.


## Example Usage

Basic Usage

```
resource "apsarastack_cr_namespace" "my-namespace" {
  name               = "my-namespace"
  auto_create        = false
  default_visibility = "PUBLIC"
}

resource "apsarastack_cr_repo" "my-repo" {
  namespace = apsarastack_cr_namespace.my-namespace.name
  name      = "my-repo"
  summary   = "this is summary of my new repo"
  repo_type = "PUBLIC"
  detail    = "this is a public repo"
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, ForceNew) Name of container registry namespace where repository is located.
* `name` - (Required, ForceNew) Name of container registry repository.
* `summary` - (Required) The repository general information. It can contain 1 to 80 characters.
* `repo_type` - (Required) `PUBLIC` or `PRIVATE`, repo's visibility.
* `detail` - (Optional) The repository specific information. MarkDown format is supported, and the length limit is 2000.

## Attributes Reference

The following attributes are exported:

* `id` - The id of Container Registry repository. The value is in format `namespace/repository`.
* `domain_list` - The repository domain list.
  * `public` - Domain of public endpoint.
  * `internal` - Domain of internal endpoint, only in some regions.
  * `vpc` - Domain of vpc endpoint.

