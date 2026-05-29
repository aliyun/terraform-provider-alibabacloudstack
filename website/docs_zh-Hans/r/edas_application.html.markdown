---
subcategory: "企业级分布式应用服务 EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_application"
sidebar_current: "docs-Alibabacloudstack-resource-edas-application"
description: |- 
  编排企业级分布式应用服务（Edas）应用
---

# alibabacloudstack_edas_application

使用 Provider 配置的凭证在指定的资源集下编排企业级分布式应用服务（Edas）应用。

## 示例用法

以下是一个完整的示例，展示如何创建一个 EDAS 应用程序资源：

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

## 参数说明

支持以下参数：

* `application_name` - (必填) EDAS 应用程序的名称。仅允许字母、数字、'-' 和 '_'。长度不得超过 36 个字符。
* `package_type` - (必填, 变更时重建) 应用程序部署的包类型。有效值为 `WAR`、`JAR` 和 `Image`。
* `cluster_id` - (必填, 变更时重建) 应用程序将要部署的集群 ID。
* `build_pack_id` - (选填) EDAS 容器的构建包 ID。在创建高速服务框架（HSF）应用程序时需要此参数。
* `component_id` - (选填) 应用程序将要部署的容器组件 ID。当首次以 WAR 包部署原生 Dubbo 或 Spring Cloud 应用程序时，必须根据部署的应用程序指定 Apache Tomcat 组件的版本。可以调用 `ListClusterOperation` 接口查询组件。
* `descriotion` - (选填) 应用程序的描述。
* `health_check_url` - (选填) 用于应用程序健康检查的 URL。
* `ecu_info` - (选填) 与应用程序关联的弹性计算单元（ECU）的信息列表。
* `group_id` - (选填) 应用程序将要部署的实例组 ID。如果希望将应用程序部署到所有组，请将此参数设置为 `all`。
* `package_version` - (选填) 您要部署的应用程序版本。它对于每个应用程序必须是唯一的。长度不得超过 64 个字符。建议使用时间戳。
* `war_url` - (选填) 用于应用程序部署的上传 Web 应用程序（WAR）包的存储地址。当 `deployType` 参数设置为 `url` 时，需要提供此参数。
* `logical_region_id` - (选填, 已弃用) 该字段在专有云中不支持，将在版本 3.21.0 中移除。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 资源的 ID。其值为应用的 AppId。

## Import

EDAS Application 可以通过 AppId 导入，例如：

```
$ terraform import alibabacloudstack_edas_application.example app-id-12345
```
