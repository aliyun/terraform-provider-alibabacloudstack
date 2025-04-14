---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_dedicatedhosts"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-dedicatedhosts"
description: |- 
  查询云服务专有宿主机
---

# alibabacloudstack_ecs_dedicatedhosts
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ecs_dedicated_hosts`

根据指定过滤条件列出当前凭证权限可以访问的ecs 专有宿主机列表。

## 示例用法

```hcl
# 创建一个专有宿主机
resource "alibabacloudstack_ecs_dedicated_host" "default" {
  dedicated_host_type = "ddh.g5"
  description        = "From_Terraform"
  dedicated_host_name = "tf_testAccEcsDedicatedHostsDataSource_4158616"
  action_on_maintenance = "Migrate"
  tags = {
    Create = "TF"
    For    = "ddh-test"
  }
}

# 使用数据源查询专有宿主机
data "alibabacloudstack_ecs_dedicated_hosts" "default" {
  ids = ["${alibabacloudstack_ecs_dedicated_host.default.id}"]
  name_regex = "tf_testAcc.*"
  status     = "Available"
  zone_id    = "cn-hangzhou-a"
  resource_group_id = "rg-acfmv7lxp9fjw"

  tags = {
    Create = "TF"
    For    = "ddh-test"
  }
}

# 输出第一个专有宿主机的ID
output "first_dedicated_host_id" {
  value = "${data.alibabacloudstack_ecs_dedicated_hosts.default.hosts.0.id}"
}
```

## 参数说明

以下参数是支持的：

* `name_regex` - (选填, 变更时重建) - 按专有宿主机名称过滤结果的正则表达式字符串。
* `ids` - (选填, 变更时重建) - 专有宿主机ID列表。最多可以输入100个预付费专有宿主机ID。多个专有宿主机ID用一个格式类似`["dh-xxxxxxxxx", "dh-yyyyyyyyy", … "dh-zzzzzzzzz"]`的JSON数组表示，ID之间用半角逗号(,)隔开。
* `dedicated_host_id` - (选填, 变更时重建) - 专有宿主机的编号。
* `dedicated_host_name` - (选填, 变更时重建) - 专有宿主机的名称。长度为2~128个字符，支持Unicode中letter分类下的字符(其中包括英文、中文和数字等)。可以包含半角冒号(:)、下划线(_)、半角句号(.)或者短划线(-)。
* `dedicated_host_type` - (选填, 变更时重建) - 专有宿主机的规格。您可以调用[DescribeDedicatedHostTypes](https://help.aliyun.com/document_detail/134240.html)接口获得最新的专有宿主机规格列表。
* `operation_locks` - (选填, 变更时重建) - 资源锁定信息。
  * `lock_reason` - (选填, 变更时重建) - 专有宿主机资源被锁定的原因。
* `resource_group_id` - (选填, 变更时重建) - 专有宿主机所在资源组ID。使用该参数过滤资源时，资源数量不能超过1000个。不支持默认资源组过滤。
* `status` - (选填, 变更时重建) - 专有宿主机的使用状态。取值范围：
  * `Available`: 运行中。专有宿主机的正常运行状态。
  * `UnderAssessment`: 物理机风险，即故障潜伏期，其物理机处于可用状态，但可能导致专有宿主机中的ECS实例出现问题。
  * `PermanentFailure`: 永久性故障，专有宿主机不可用。
  * `TempUnavailable`: 宿主机临时不可用。
  * `Redeploying`: 宿主机恢复中。默认值：`Available`。
* `zone_id` - (选填, 变更时重建) - 可用区ID。您可以调用[DescribeZones](https://help.aliyun.com/document_detail/25610.html)查看最新的阿里云可用区列表。
* `tags` - (选填) - 要分配给资源的标签映射。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 专有宿主机ID列表。
* `names` - 专有宿主机名称列表。
* `hosts` - 专有宿主机列表。每个元素包含以下属性：
  * `id` - 专有宿主机的ID。
  * `action_on_maintenance` - 当专有宿主机发生故障或需要在线修复时，用于迁移专有宿主机上的实例的策略。有效值：
    * `Migrate`: 实例迁移到另一台物理服务器并重新启动。
    * `Stop`: 实例停止。如果专有宿主机无法修复，实例将迁移到另一台物理机并重新启动。
  * `auto_placement` - 指定是否将专有宿主机添加到自动部署的资源池中。有效值：
    * `on`: 将专有宿主机添加到自动部署的资源池中。
    * `off`: 不将专有宿主机添加到自动部署的资源池中。
  * `auto_release_time` - 专有宿主机的自动释放时间。指定ISO 8601标准的时间格式为`yyyy-MM-ddTHH:mm:ssZ`。时间必须是UTC+0。
  * `capacity` - 容量。
    * `available_local_storage` - 剩余本地磁盘容量。单位：GiB。
    * `available_memory` - 剩余内存容量，单位：GiB。
    * `available_vcpus` - 剩余vCPU核心数。
    * `available_vgpus` - 可用虚拟GPU的数量。
    * `local_storage_category` - 本地磁盘类型。
    * `total_local_storage` - 本地磁盘的总容量，单位：GiB。
    * `total_memory` - 总内存容量，单位：GiB。
    * `total_vcpus` - vCPU核心总数。
    * `total_vgpus` - 虚拟GPU的总数。
  * `cores` - 核心数。
  * `cpu_over_commit_ratio` - CPU过载比。仅自定义规格g6s、c6s和r6s支持设置CPU超卖比。有效值：1到5。
  * `dedicated_host_id` - 专有宿主机的ID。
  * `dedicated_host_name` - 专有宿主机的名称。
  * `dedicated_host_type` - 专有宿主机的类型。
  * `description` - 专有宿主机的描述。
  * `expired_time` - 订阅专有宿主机的过期时间。
  * `gpu_spec` - GPU型号。
  * `machine_id` - 专有宿主机的机器代码。
  * `network_attributes` - 网络属性。
    * `slb_udp_timeout` - Server Load Balancer (SLB)与专有宿主机之间的UDP会话超时时间。单位：秒。
    * `udp_timeout` - 用户与Apsara Stack Cloud服务在专有宿主机上的UDP会话超时时间。单位：秒。
  * `operation_locks` - 操作锁。
    * `lock_reason` - 专有宿主机资源被锁定的原因。
  * `payment_type` - 专有宿主机的计费方式。
  * `physical_gpus` - 物理GPU的数量。
  * `resource_group_id` - 专有宿主机所属的资源组ID。
  * `sale_cycle` - 订阅计费方法的单位。
  * `sockets` - 物理CPU的数量。
  * `status` - 专有宿主机的服务状态。
  * `supported_custom_instance_type_families` - 专有宿主机支持的自定义实例类型族。
  * `supported_instance_type_families` - 支持的实例类型族。
  * `supported_instance_types_list` - 支持的实例类型列表。
  * `tags` - 标签。
  * `zone_id` - 专有宿主机所在的可用区ID。