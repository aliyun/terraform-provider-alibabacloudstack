---
subcategory: "BMCP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bmcp_machinetypes"
sidebar_current: "docs-Alibabacloudstack-datasource-bmcp-machinetypes"
description: |-
  查询裸金属计算平台(BMCP)机型规格
---

# alibabacloudstack_bmcp_machinetypes

根据指定过滤条件列出当前凭证权限可以访问的BMCP机型规格列表。

## 示例用法

### 查询所有机型规格

```hcl
data "alibabacloudstack_bmcp_machinetypes" "default" {}
```

### 按名称正则表达式过滤

```hcl
data "alibabacloudstack_bmcp_machinetypes" "name_filtered" {
  name_regex = "PG"
}
```

### 按架构正则表达式过滤

```hcl
data "alibabacloudstack_bmcp_machinetypes" "arch_filtered" {
  arch_regex = "x86"
}
```

### 按最小标准实例数过滤

```hcl
data "alibabacloudstack_bmcp_machinetypes" "min_count_filtered" {
  min_standard_instance_count = 1
}
```

### 按最大标准实例数过滤

```hcl
data "alibabacloudstack_bmcp_machinetypes" "max_count_filtered" {
  max_standard_instance_count = 100
}
```

### 按标准实例数范围过滤

```hcl
data "alibabacloudstack_bmcp_machinetypes" "range_filtered" {
  min_standard_instance_count = 1
  max_standard_instance_count = 100
}
```

### 组合过滤条件

```hcl
data "alibabacloudstack_bmcp_machinetypes" "combined" {
  name_regex                  = "PG"
  arch_regex                  = "x86"
  min_standard_instance_count = 1
}
```

### 动态过滤（先查询再过滤）

```hcl
# 首先查询所有机型规格获取真实数据
data "alibabacloudstack_bmcp_machinetypes" "all" {
}

# 使用真实数据中的值进行动态过滤
data "alibabacloudstack_bmcp_machinetypes" "filtered" {
  name_regex = data.alibabacloudstack_bmcp_machinetypes.all.machinetypes.0.name
}
```

## 参数说明

以下参数是支持的：

* `name_regex` - (选填, 变更时重建) 用于按名称过滤结果的正则表达式字符串。
* `arch_regex` - (选填, 变更时重建) 用于按CPU架构过滤结果的正则表达式字符串。
* `min_standard_instance_count` - (选填, 变更时重建) 最小标准实例数过滤。仅返回标准实例数大于或等于此值的机型规格。
* `max_standard_instance_count` - (选填, 变更时重建) 最大标准实例数过滤。仅返回标准实例数小于或等于此值的机型规格。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 机型规格ID列表。
* `machinetypes` - 机型规格列表。每个元素包含以下属性：
  * `id` - 机型规格的ID。
  * `name` - 机型规格的名称。
  * `description` - 机型规格的描述信息。
  * `deploy_type` - 机型规格的部署类型。
  * `manufacturer` - 机型规格的制造商。
  * `cpu_arch` - 机型规格的CPU架构。
  * `cpu_model` - 机型规格的CPU型号。
  * `cpu_manufacturer` - 机型规格的CPU制造商。
  * `cpu_number` - CPU数量。
  * `memory` - 内存大小(GB)。
  * `disk` - 磁盘大小。
  * `disk_type` - 磁盘类型。
  * `gpu` - GPU型号。
  * `gpu_num` - GPU数量。
  * `gpu_manufacturer` - GPU制造商。
  * `video_memory` - 显存大小(GB)。
  * `tflops_fp32` - TFLOPS FP32性能。
  * `network_card_type` - 网卡类型。
  * `network_card_num` - 网卡数量。
  * `rated_power` - 额定功率。
  * `specification` - 机型规格的详细规格。
  * `unit_num` - 单元数量。
  * `standard_instance_count` - 该机型规格可用的标准实例数量。
  * `default_fmin` - 默认FMin值。
  * `create_time` - 创建时间。
  * `update_time` - 更新时间。
