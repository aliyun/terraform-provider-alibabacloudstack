---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_applications"
sidebar_current: "docs-Alibabacloudstack-datasource-edas-applications"
description: |- 
  Provides a list of edas applications owned by an alibabacloudstack account.
---

# alibabacloudstack_edas_applications
-> **NOTE:** Alias name has: `alibabacloudstack_edas_slbattachments`

This data source provides a list of edas applications in an alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_edas_applications" "applications" {
  ids        = ["app12345"]
  name_regex = "example-application-.*"
  output_file = "application_list.txt"
}

output "first_application_name" {
  value = data.alibabacloudstack_edas_applications.applications.applications[0].app_name
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of application IDs to filter results. If not provided, all applications will be considered.
* `name_regex` - (Optional) A regex string to filter results by the application name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `names` - A list of names of all matched EDAS applications.
* `ids` - A list of IDs of all matched EDAS applications.
* `applications` - A list of EDAS applications. Each element contains the following attributes:
  * `app_name` - The name of the EDAS application. Only letters, numbers, '-', and '_' are allowed. The length cannot exceed 36 characters.
  * `app_id` - The ID of the application.
  * `application_type` - The type of the package for the deployment of the application. Valid values are `WAR` and `JAR`.
  * `build_package_id` - The package ID of Enterprise Distributed Application Service (EDAS) Container.
  * `cluster_id` - The ID of the cluster that the application belongs to.
  * `cluster_type` - The type of the cluster that the application belongs to. Valid values: 
    * `1`: Swarm cluster.
    * `2`: ECS cluster.
    * `3`: Kubernetes cluster.
  * `region_id` - The ID of the region where the application is located.