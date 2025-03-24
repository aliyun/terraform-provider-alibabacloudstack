---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_dedicatedhost"
sidebar_current: "docs-Alibabacloudstack-ecs-dedicatedhost"
description: |- 
  编排云服务器（Ecs）专有宿主机
---

# alibabacloudstack_ecs_dedicatedhost
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ecs_dedicated_host`

使用Provider配置的凭证在指定的资源集下编排云服务器（Ecs）专有宿主机。

## 示例用法

### 基本使用

```hcl
resource "alibabacloudstack_ecs_dedicated_host" "default" {
  dedicated_host_type = "ddh.g5"
  payment_type        = "PostPaid"
  tags = {
    Create = "Terraform"
    For    = "DDH"
  }
  description         = "From_Terraform"
  dedicated_host_name = "dedicated_host_name"
}
```

### 预付费专用主机

```hcl
resource "alibabacloudstack_ecs_dedicated_host" "prepaid" {
  dedicated_host_type = "ddh.g5"
  payment_type        = "PrePaid"
  auto_renew          = true
  auto_renew_period   = 1
  sale_cycle          = "Month"
  expired_time        = 12
  tags = {
    Create = "Terraform"
    For    = "DDH"
  }
  description         = "Prepaid_DDH"
  dedicated_host_name = "prepaid_dedicated_host"
}
```

### 自定义CPU超分比率

```hcl
resource "alibabacloudstack_ecs_dedicated_host" "custom_cpu" {
  dedicated_host_type = "ddh.c6s"
  payment_type        = "PostPaid"
  cpu_over_commit_ratio = 4
  tags = {
    Create = "Terraform"
    For    = "DDH"
  }
  description         = "Custom_CPU_Ratio"
  dedicated_host_name = "custom_cpu_ratio_host"
}
```

### 指定区域和自动释放时间

```hcl
resource "alibabacloudstack_ecs_dedicated_host" "with_zone" {
  dedicated_host_type = "ddh.c5"
  payment_type        = "PostPaid"
  zone_id            = "cn-hangzhou-a"
  auto_release_time  = "2025-12-31T23:59:59Z"
  dedicated_host_name = "test-name"
}
```

## 参数参考

支持以下参数：

* `action_on_maintenance` - (选填) - 当专有宿主机发生故障或者在线修复时，为其所宿实例设置迁移方案。取值范围：
  * `Migrate`: 迁移实例到其他物理机并重新启动实例。当专有宿主机上挂载云盘存储时，默认值：Migrate。
  * `Stop`: 在当前专有宿主机上停止实例，确认无法修复专有宿主机后，迁移实例到其他物理机并重新启动实例。当专有宿主机上挂载本地盘存储时，默认值：Stop。

* `auto_placement` - (选填) - 专有宿主机是否加入自动部署资源池。当您在专有宿主机上创建实例，却不指定**DedicatedHostId**时，阿里云将自动从加入资源池的专有宿主机中，为您选取适合的宿主机部署实例，更多信息，请参见[自动部署功能介绍](https://help.aliyun.com/document_detail/118938.html)。取值范围：
  * `on`: 加入自动部署资源池。
  * `off`: 不加入自动部署资源池。
  默认值：`on`。

* `auto_release_time` - (选填) - 专有宿主机的自动释放时间。按照ISO 8601标准表示，并使用UTC+0时间，格式为`yyyy-MM-ddTHH:mm:ssZ`。
  * 必须晚于当前时间起算的半小时及以后。
  * 必须早于当前时间起算的三年及以前。
  * 如果参数值中的秒(ss)不是00，则自动取为00。
  * 如果不输入`AutoReleaseTime`参数，表示取消自动释放，专有宿主机在预约时间点不再自动释放。

* `auto_renew` - (选填) - 是否自动续费包年包月专有宿主机。取值范围：
  * `true`: 自动续费包年包月专有宿主机。
  * `false`: 不自动续费包年包月专有宿主机。
  默认值：`false`。

* `auto_renew_period` - (选填) - 单次自动续费的周期。当参数**AutoRenew**为`true`时，**AutoRenewPeriod**参数方可生效，并为必选参数。取值范围：
  * PeriodUnit=Week时：1、2、3。
  * PeriodUnit=Month时：1、2、3、6、12。
  
* `cpu_over_commit_ratio` - (选填) - CPU超卖比。仅自定义规格g6s、c6s、r6s支持设置CPU超卖比。取值范围：1~5。CPU超卖比影响DDH的可用vCPU数，一台DDH的可用vCPU数=物理CPU核数*2*CPU超卖比。例如，g6s的物理CPU核数为52，如果设置CPU超卖比为4，则DDH创建完成后vCPU总数显示为416。

* `dedicated_host_cluster_id` - (选填) - 专有宿主机所属的专有宿主机集群ID。

* `dedicated_host_name` - (选填) - 专有宿主机的名称。长度为2~128个字符，支持Unicode中letter分类下的字符(其中包括英文、中文和数字等)。可以包含半角冒号(:)、下划线(_)、半角句号(.)或者短划线(-)。

* `dedicated_host_type` - (必填, 变更时重建) - 专有宿主机的规格。您可以调用[DescribeDedicatedHostTypes](https://help.aliyun.com/document_detail/134240.html)接口获得最新的专有宿主机规格列表。

* `description` - (选填) - 专有宿主机的描述。长度为2~256个英文或中文字符，不能以`http://`或`https://`开头。

* `dry_run` - (选填) - 是否只预检此次请求。取值范围：
  * `true`: 发送检查请求，不会查询资源状况。检查项包括AccessKey是否有效、RAM用户的授权情况和是否填写了必填参数。如果检查不通过，则返回对应错误。如果检查通过，会返回错误码`DryRunOperation`。
  * `false`: 发送正常请求，通过检查后返回2XX的HTTP状态码并直接查询资源状况。
  默认值为`false`。

* `expired_time` - (选填) - 续费周期。取值范围：
  * PeriodUnit=Week时：1、2、3、4。
  * PeriodUnit=Month时：1、2、3、4、5、6、7、8、9、12、24、36、48、60。
  * PeriodUnit=Year时：1、2、3、4、5。

* `min_quantity` - (选填) - 指定专有宿主机的最小购买数量。取值范围：1~100。

* `network_attributes` - (选填) - 宿主机网络参数
  * `udp_timeout` - (选填) - 用户与阿里云服务在专用主机上的UDP会话超时时间。单位：秒。有效值：`15` 到 `310`。
  * `slb_udp_timeout` - (选填) - SLB与专用主机之间的UDP会话超时时间。单位：秒。有效值：`15` 到 `310`。

* `payment_type` - (选填) - 专有宿主机的计费方式。取值范围：
  * `PrePaid`: 包年包月。
  * `PostPaid`: 按量付费。
  默认值：`PostPaid`。

* `resource_group_id` - (选填) - 专有宿主机所在资源组ID。使用该参数过滤资源时，资源数量不能超过1000个。

* `sale_cycle` - (选填) - 续费时长单位。取值范围：
  * `Week`
  * `Month`
  * `Year`
  默认值：`Month`。

* `zone_id` - (选填, 变更时重建) - 可用区ID。您可以调用[DescribeZones](https://help.aliyun.com/document_detail/25610.html)查看最新的阿里云可用区列表。

* `tags` - (选填) - 分配给资源的标签映射。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 专有宿主机ID。
* `status` - 专有宿主机的使用状态。取值范围：
  * `Available`: 运行中。专有宿主机的正常运行状态。
  * `UnderAssessment`: 物理机风险，即故障潜伏期，其物理机处于可用状态，但可能导致专有宿主机中的ECS实例出现问题。
  * `PermanentFailure`: 永久性故障，专有宿主机不可用。
  * `TempUnavailable`: 宿主机临时不可用。
  * `Redeploying`: 宿主机恢复中。
  默认值：`Available`。

### 超时

`timeouts`块允许您为某些动作指定[超时](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts)：

* `create` - (默认为11分钟)用于创建专有宿主机时。
* `delete` - (默认为1分钟)用于删除专有宿主机时。
* `update` - (默认为11分钟)用于更新专有宿主机时。

## 导入

ECS专有宿主机可以通过id导入，例如：

```bash
$ terraform import alibabacloudstack_ecs_dedicated_host.default dh-2zedmxxxx
```