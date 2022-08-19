---
subcategory: "Cloud Monitor"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_cms_project_meta"
sidebar_current: "docs-apsarastack-datasource-cms-project-meta"
description: |-
    Provides a list of project meta owned by an Apsarastack Cloud account.
---

# apsarastack\_cms\_project\_meta

Provides a list of project meta owned by an Apsarastack Cloud account.

## Example Usage

Basic Usage

```
data "apsarastack_cms_project_meta" "default" {
  name_regex = "OSS"
}

output "project_meta" {
  value = data.apsarastack_cms_project_meta.default.*
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter results by project meta description.

## Attributes Reference

The following attributes are exported:

* `resources` - A list of cms project meta. Each element contains the following attributes:
    * `description` - Description for a project meta.
    * `labels` - Labels for a cms project meta. A tag of a metric is used as a special mark of alerts triggered by the metric. The format is `[{"name":"Tag name","value":"Tag value"}, {"name":"Tag name","value":"Tag value"}]`.
    * `namespace` - The namespace of the service, which is used to distinguish between services. Generally, the value is in the format acs_Abbreviation of the service name .
  