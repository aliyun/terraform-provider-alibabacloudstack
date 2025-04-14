---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_reservedinstance"
sidebar_current: "docs-Alibabacloudstack-ecs-reservedinstance"
description: |- 
  编排云服务器（Ecs）预留实例
---

# alibabacloudstack_ecs_reservedinstance
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_reserved_instance`

使用Provider配置的凭证在指定的资源集下编排云服务器（Ecs）预留实例。

## 示例用法

```hcl
variable "name" {
    default = "tf-testaccecsreserved_instance37879"
}

resource "alibabacloudstack_ecs_reservedinstance" "default" {
  instance_type      = "ecs.t6-c4m1.large"
  instance_amount    = 2
  period_unit        = "Year"
  offering_type      = "All Upfront"
  reserved_instance_name = var.name
  description        = "ReservedInstance for testing"
  zone_id           = "cn-hangzhou-i"
  scope             = "Zone"
  period            = 3
  platform          = "Linux"
  resource_group_id = "rg-acfm5xxxxxx"
}
```

## 参数说明

支持以下参数：

* `instance_type` - (必填, 变更时重建) - 匹配的按量付费实例的规格。例如，`ecs.t6-c4m1.large`。
* `scope` - (选填, 变更时重建) - 预留实例的范围。可选值：
  * `Region`：区域级。
  * `Zone`：地域级。
  默认值为 `Region`。当设置为 `Zone` 时，必须指定 `zone_id`。
* `zone_id` - (选填, 变更时重建) - 可用区ID。当 `scope` 设置为 `Zone` 时，此参数是必填的。有关可用区列表的信息，请参见 [DescribeZones](https://www.alibabacloud.com/help/doc-detail/25610.html)。
* `instance_amount` - (选填, 变更时重建) - 可以匹配同规格按量付费实例的数量。默认值为 `1`。
* `platform` - (选填, 变更时重建) - 实例使用的镜像的操作系统类型。可能值：
  * `Windows`：Windows Server 类型的操作系统。
  * `Linux`：Linux 及类 Unix 类型的操作系统。
* `period_unit` - (选填, 变更时重建) - 购买预留实例券的时长单位。取值范围：
  * 国际区域：`Year`
  * 中国区域：`Year`、`Month`
  默认值：在中国区域为 `Month`，在国际区域为 `Year`。
* `period` - (选填, 变更时重建) - 购买预留实例券的时长。取值范围：
  * 当 `PeriodUnit` 是 `Year` 时，取值范围：`1`、`3`、`5`。
  * 当 `PeriodUnit` 是 `Month` 时，取值范围是 `1`。
  默认值：`1`。
* `offering_type` - (选填, 变更时重建) - 预留实例的支付类型。可选值：
  * `No Upfront`：无需预付款。
  * `Partial Upfront`：部分预付款。
  * `All Upfront`：全额预付款。
* `reserved_instance_name` - (选填) - 预留实例名称。名称必须是一个包含 2 到 128 个字符的字符串，可以包含字母、数字、冒号 (`:`)、下划线 (`_`) 和连字符。它必须以字母开头。不能以 `http://` 或 `https://` 开头。
* `description` - (选填) - 预留实例描述。长度为 2 到 256 个英文或中文字符。不能以 `http://` 或 `https://` 开头。
* `resource_group_id` - (选填, 变更时重建) - 资源组 ID。
* `reserved_instance_id` - (选填, 变更时重建) - 预留实例的 ID。
* `name` - (选填, 已弃用) - 预留实例的名称。

### 从配置中移除 alibabacloudstack_ecs_reservedinstance

alibabacloudstack_ecs_reservedinstance 资源允许您管理您的预留实例，但 Terraform 无法销毁它。将此资源从您的配置中移除会将其从您的状态文件和管理中移除，但不会销毁预留实例。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 预留实例的 ID。
* `instance_amount` - 可以匹配同规格按量付费实例的数量。
* `platform` - 实例使用的镜像的操作系统类型。可能值：
  * `Windows`：Windows Server 类型的操作系统。
  * `Linux`：Linux 及类 Unix 类型的操作系统。
* `reserved_instance_name` - 预留实例的名称。
* `resource_group_id` - 资源组 ID。
* `reserved_instance_id` - 预留实例的 ID。