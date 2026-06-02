---
subcategory: "Enterprise Distributed Application Service (EDAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_application"
sidebar_current: "docs-Alibabacloudstack-resource-edas-application"
description: |- 
  Provides a Edas Application resource.
---

# alibabacloudstack_edas_application

Provides a Edas application resource.

## Example Usage

Basic Usage

```hcl
variable "name" {
    default = "tf-testacc-edasapplicationbasic4966"
}

resource "alibabacloudstack_vpc" "default" {
    cidr_block = "172.16.0.0/12"
    name       = "${var.name}"
}

resource "alibabacloudstack_edas_cluster" "default" {
    cluster_name = "${var.name}"
    cluster_type = 2
    network_mode = 2
    vpc_id       = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_edas_application" "default" {
    component_id     = 7
    application_name = "${var.name}"
    package_type     = "WAR"
    cluster_id       = "${alibabacloudstack_edas_cluster.default.id}"
    build_pack_id    = 1
    descriotion      = "Test Description"
    group_id         = "all"
    package_version  = "v1.0.0"
    war_url          = "http://example.com/app.war"
    health_check_url = "/health"
}
```

## Argument Reference

The following arguments are supported:

* `application_name` - (Required) The name of the EDAS application. Only letters, numbers, '-', and '_' are allowed. The length cannot exceed 36 characters.
* `package_type` - (Required, ForceNew) The type of the package for the deployment of the application. Valid values are `WAR`, `JAR`, and `Image`.
* `cluster_id` - (Required, ForceNew) The ID of the cluster where the application will be deployed.
* `build_pack_id` - (Optional) The build pack ID of the EDAS container. This is required when creating a High-speed Service Framework (HSF) application.
* `component_id` - (Optional) The ID of the component in the container where the application is going to be deployed. When deploying a native Dubbo or Spring Cloud application using a WAR package for the first time, you must specify the Apache Tomcat component version based on the deployed application. You can call the `ListClusterOperation` API to query the components.
* `descriotion` - (Optional) A description of the application.
* `health_check_url` - (Optional) The URL used for health checking of the application.
* `ecu_info` - (Optional) A list of Elastic Compute Unit (ECU) information associated with the application.
* `group_id` - (Optional) The ID of the instance group where the application will be deployed. Set this parameter to `all` if you want to deploy the application to all groups.
* `package_version` - (Optional) The version of the application that you want to deploy. It must be unique for every application. The length cannot exceed 64 characters. A timestamp is recommended.
* `war_url` - (Optional) The storage address of the uploaded web application (WAR) package for application deployment. This parameter is needed when the `deployType` parameter is set to `url`.
* `logical_region_id` - (Optional, Deprecated) This field is not supported on ApsaraStack and will be removed in version 3.21.0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the resource. The value is the application's AppId.

## Import

EDAS Application can be imported using the AppId, e.g.

```
$ terraform import alibabacloudstack_edas_application.example app-id-12345
```
