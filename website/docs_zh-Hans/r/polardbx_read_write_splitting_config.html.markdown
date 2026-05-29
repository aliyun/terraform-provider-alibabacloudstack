---
subcategory: "云原生分布式数据库PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_read_write_splitting_config"
sidebar_current: "docs-Alibabacloudstack-polardbx-read-write-splitting-config"
description: |-
 阿里云 PolarDBX 读写分离配置
---
# 阿里云 PolarDBX 读写分离配置

阿里云 PolarDBX 读写分离配置

## Example Usage

```hcl
resource "alibabacloudstack_polardbx_read_write_splitting_config" "example" {
  db_instance_id                  = "example-polardbx-instance-id"
  attend_htap_list                = ["example-polardbx-instance-id1", "example-polardbx-instance-id2"]
  auto_attend_htap                = true
  delay_execution_strategy        = 3
  enable_consistent_replica_read  = true
  enable_htap                     = true
  master_read_weight              = 50
  storage_delay_threshold         = 1000
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (必填, 强制新建) PolarDBX 实例的 ID。
* `attend_htap_list` - (可选) 关联实例列表
* `auto_attend_htap` - (可选) 是否自动将只读实例自动加入读写分离
* `delay_execution_strategy` - (可选) 只读流量切回主实例
* `enable_consistent_replica_read` - (可选) 数据读是否强一致性
* `enable_htap` - (可选) 是否开启读写分离
* `master_read_weight` - (可选) 只读流量占比。有效值为 1 到 100。
* `storage_delay_threshold` - (可选) 只读实例延迟阈值(s)


## Attributes Reference

The following attributes are exported:

* `id` - 资源的 ID，与 db_instance_id 相同。

## Import

可以使用 DB 实例 ID 导入 PolarDBX 读写分离配置，例如：

```bash
$ terraform import alibabacloudstack_polardbx_read_write_splitting_config.example <db_instance_id>
```