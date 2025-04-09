---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_application"
sidebar_current: "docs-Alibabacloudstack-edas-application"
description: |- 
  使用Provider配置的凭证在指定的资源集下编排企业级分布式应用服务（Edas） Application resource.
---

# alibabacloudstack_edas_application

使用Provider配置的凭证在指定的资源集下编排企业级分布式应用服务（Edas） application resource.

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
    component_id     = "7"
    application_name = "${var.name}"
    package_type     = "WAR"
    cluster_id       = "${alibabacloudstack_edas_cluster.default.id}"
    build_pack_id   = 1
    descriotion      = "Test Description"
    health_check_url = "/health"
    group_id         = "group-id-12345"
    package_version  = "v1.0.0"
    war_url          = "http://example.com/app.war"
}
```

## Argument Reference

The following arguments are supported:

* `application_name` - (Required) The name of the EDAS application. Only letters, numbers, '-', and '_' are allowed. The length cannot exceed 36 characters.
* `package_type` - (Required, ForceNew) The type of the package for the deployment of the application. Valid values are `WAR` and `JAR`.
* `cluster_id` - (Required, ForceNew) The ID of the cluster where the application will be deployed. If not specified, the default cluster will be used.
* `build_pack_id` - (Optional) The package ID of the EDAS container. This is required when creating a High-speed Service Framework (HSF) application.
* `component_id` - (Optional) The ID of the component in the container where the application is going to be deployed. When deploying a native Dubbo or Spring Cloud application using a WAR package for the first time, you must specify the version of the Apache Tomcat component based on the deployed application. You can call the `ListClusterOperation` interface to query the components.
* `descriotion` - (Optional) A description of the application.
* `health_check_url` - (Optional) The URL used for health checking of the application.
* `region_id` - (Optional) The ID of the region where the application will be created. You can call the `ListUserDefineRegion` operation to query the region ID.
* `group_id` - (Optional) The ID of the instance group where the application will be deployed. Set this parameter to `all` if you want to deploy the application to all groups.
* `package_version` - (Optional) The version of the application that you want to deploy. It must be unique for every application. The length cannot exceed 64 characters. We recommend using a timestamp.
* `war_url` - (Optional) The address to store the uploaded web application (WAR) package for application deployment. This parameter is required when the `deployType` parameter is set as `url`.
* `ecu_info` - (Optional) Information about the Elastic Compute Unit (ECU) associated with the application.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the resource. The value is formulated as `app_Id`.