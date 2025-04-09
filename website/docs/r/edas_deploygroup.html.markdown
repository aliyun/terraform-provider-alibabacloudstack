---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_deploygroup"
sidebar_current: "docs-Alibabacloudstack-edas-deploygroup"
description: |- 
  使用Provider配置的凭证在指定的资源集下编排企业级分布式应用服务（Edas） Deploygroup resource.
---

# alibabacloudstack_edas_deploygroup
-> **NOTE:** Alias name has: `alibabacloudstack_edas_deploy_group`

使用Provider配置的凭证在指定的资源集下编排企业级分布式应用服务（Edas） Deploygroup resource.

## Example Usage

```hcl
variable "name" {
  default = "tf-testacc-edasdeploygroupbasic4916"
}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  name       = var.name
}

resource "alibabacloudstack_edas_cluster" "default" {
  cluster_name = var.name
  cluster_type = 2
  network_mode = 2
  vpc_id       = alibabacloudstack_vpc.default.id
}

resource "alibabacloudstack_edas_application" "default" {
  application_name = var.name
  cluster_id      = alibabacloudstack_edas_cluster.default.id
  package_type    = "JAR"
  build_pack_id   = "15"
}

resource "alibabacloudstack_edas_deploy_group" "default" {
  app_id       = alibabacloudstack_edas_application.default.id
  group_name   = var.name
  group_type   = 2 # Traffic Management Enable Grayscale
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The ID of the application that you want to deploy. This is the unique identifier for the application in EDAS. 
* `group_name` - (Required, ForceNew) The name of the instance group that you want to create. It must be unique within the application and cannot be modified after creation. 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the resource. It is formulated as `<app_id>:<group_name>:<group_id>`.
* `group_type` - The type of the instance group. This attribute reflects the value set during creation and indicates the grouping behavior: 
  - `0`: Default Grouping.
  - `1`: Grayscale is not enabled for traffic management.
  - `2`: Traffic Management Enable Grayscale.

This attribute helps identify the configuration of the deploy group and its capabilities within the EDAS environment.