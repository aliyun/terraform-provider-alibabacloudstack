---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_applications"
sidebar_current: "docs-Alibabacloudstack-datasource-edas-applications"
description: |- 
  查询企业级分布式应用服务应用
---

# alibabacloudstack_edas_applications
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_edas_slbattachments`

根据指定过滤条件列出当前凭证权限可以访问的企业级分布式应用服务应用列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testacc-edas-applications2798"
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
  application_name = "${var.name}"
  cluster_id      = "${alibabacloudstack_edas_cluster.default.id}"
  package_type    = "WAR"
  build_pack_id   = "15"
}

data "alibabacloudstack_edas_applications" "default" {
  ids        = ["${alibabacloudstack_edas_application.default.id}"]
  name_regex = "${alibabacloudstack_edas_application.default.application_name}"
  output_file = "edas_applications_output.txt"
}

output "application_names" {
  value = data.alibabacloudstack_edas_applications.default.names
}
```

## 参数参考

以下参数是支持的：

* `ids` - (可选) 应用程序ID列表，用于过滤结果。如果未提供，则会考虑所有应用程序。
* `name_regex` - (可选) 用于按应用程序名称过滤结果的正则表达式字符串。

## 属性说明

除了上述参数外，还导出以下属性：

* `names` - 所有匹配的EDAS应用程序的名称列表。
* `ids` - 所有匹配的EDAS应用程序的ID列表。
* `applications` - EDAS应用程序列表。每个元素包含以下属性：
  * `app_name` - EDAS应用程序的名称。仅允许字母、数字、'-'和'_'。长度不能超过36个字符。
  * `app_id` - 应用程序的ID。
  * `application_type` - 部署应用程序的包类型。有效值为：  
    - `WAR`：表示WAR包部署。  
    - `JAR`：表示JAR包部署。
  * `build_package_id` - 容器版本ID（即EDAS容器的包ID）。
  * `cluster_id` - 应用程序所属的集群ID。
  * `cluster_type` - 应用程序所属的集群类型。有效值为：  
    - `1`：Swarm集群。  
    - `2`：ECS集群。  
    - `3`：Kubernetes集群。
  * `region_id` - 应用程序所在的区域ID。