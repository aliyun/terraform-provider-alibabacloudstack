---
subcategory: "裸机管理 BMS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bms_machinetypes"
sidebar_current: "docs-Alibabacloudstack-datasource-bms-machinetypes"
description: |-
  查询AlibabacloudStack BMS 机型
---

# alibabacloudstack_bms_machinetypes

该数据源提供了用户可用的 BMS 机型信息。

-> **注意:** 版本 v3.16.0+ 可用。

## 示例代码

```hcl
data "alibabacloudstack_bms_machinetypes" "default" {
  name_regex = "PG"
}

output "first_machine_type" {
  value = data.alibabacloudstack_bms_machinetypes.default.machinetypes.0.name
}
```

## 参数说明

支持以下参数：

* `name_regex` - （可选）用于按机型名称过滤结果的正则表达式字符串。
* `arch_regex` - （可选）用于按 CPU 架构过滤结果的正则表达式字符串。

## 属性说明

导出以下属性：

* `ids` - 机型 ID 列表。
* `machinetypes` - 机型对象列表。每个元素包含以下属性：
  * `id` - 机型 ID。
  * `name` - 机型名称。
  * `description` - 机型描述。
  * `deploy_type` - 部署类型。有效值：`bmcp`、`bmcp_managed`、`bmcp_no_clone`、`base`、`ehpc`、`ehpc_managed`、`aspeed`。
  * `manufacturer` - 机器制造商。
  * `cpu_arch` - CPU 架构（例如 x86、ARM）。
  * `cpu_model` - CPU 型号。
  * `cpu_manufacturer` - CPU 制造商。
  * `cpu_number` - CPU 核心数。
  * `memory` - 内存大小（单位：GB）。
  * `disk` - 磁盘大小（单位：GB）。
  * `disk_type` - 磁盘类型。
  * `gpu` - GPU 型号。
  * `gpu_num` - GPU 数量。
  * `gpu_manufacturer` - GPU 制造商。
  * `video_memory` - 显存大小（单位：GB）。
  * `tflops_fp32` - TFLOPS FP32 性能指标。
  * `network_card_type` - 网卡类型。
  * `network_card_num` - 网卡数量。
  * `rated_power` - 额定功率。
  * `specification` - 规格详情。
  * `unit_num` - 单元数量。
  * `create_time` - 创建时间。
  * `update_time` - 最后更新时间。
