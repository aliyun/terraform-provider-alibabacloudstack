---
subcategory: "EIP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_eip_address"
sidebar_current: "docs-Alibabacloudstack-eip-address"
description: |- 
  编排弹性公网地址
---

# alibabacloudstack_eip_address
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_eip`

使用Provider配置的凭证在指定的资源集下编排弹性公网地址。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAcceEipName5478"
}

resource "alibabacloudstack_eip_address" "default" {
  name        = var.name
  description = "This is a test EIP address."
  bandwidth   = "5"
  ip_address  = "192.168.0.1"
  tags        = {
    Environment = "Test"
    CreatedBy   = "Terraform"
  }
}
```

## 参数参考

支持以下参数：

* `name` - (可选) EIP实例的名称。该名称可以包含2到128个字符，必须仅包含字母数字字符或连字符(例如“-”、“.”、“_”)，并且不能以连字符开头或结尾，也不能以`http://`或`https://`开头。默认值为null。
* `description` - (可选) EIP实例的描述。此描述可以包含2到256个字符。它不能以`http://`或`https://`开头。创建预付费的EIP实例时，不支持设置该参数。默认值为null。
* `bandwidth` - (可选) 要指定申请的EIP的带宽峰值，单位：Mbps。如果未指定此值，默认为**5** Mbps。
  - 当`payment_type`取值为`PayAsYouGo`，且`internet_charge_type`取值为`PayByBandwidth`时，`bandwidth`取值范围为**1**~**500**。
  - 当`payment_type`取值为`PayAsYouGo`，且`internet_charge_type`取值为`PayByTraffic`时，`bandwidth`取值范围为**1**~**200**。
  - 当`payment_type`取值为`Subscription`时，`bandwidth`取值范围为**1**~**1000**。
* `ip_address` - (可选，变更时重建) EIP的IP地址。最多支持50个EIP的IP地址。如果需要指定特定的IP地址，请确保其有效且未被占用。
* `tags` - (可选，映射) 要分配给资源的标签映射。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - EIP的唯一标识符。
* `status` - EIP的状态，取值：
  - **Associating**：绑定中。
  - **Unassociating**：解绑中。
  - **InUse**：已分配。
  - **Available**：可用。
  - **Releasing**：释放中。
* `ip_address` - 分配给EIP的弹性IP地址。
* `name` - EIP实例的名称。该名称可以包含2到128个字符，必须仅包含字母数字字符或连字符(例如“-”、“.”、“_”)，并且不能以连字符开头或结尾，也不能以`http://`或`https://`开头