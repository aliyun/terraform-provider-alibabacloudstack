---
subcategory: "Hologram"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_hologram_instance_backup_policy"
sidebar_current: "docs-alibabacloudstack-resource-hologram-instance-backup-policy"
description: |-
  提供Hologram实例备份策略资源
---

# alibabacloudstack_hologram_instance_backup_policy

提供Hologram实例备份策略资源。

## 示例用法

```hcl
variable "name" {
  default = "tf-testacc"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

data "alibabacloudstack_hologram_clusters" "default" {
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
  lifecycle {
      ignore_changes = [
        secondary_cidr_blocks,
        tags
      ]
  }
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  vswitch_name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  lifecycle {
      ignore_changes = [
        tags
      ]
  }
}

resource "alibabacloudstack_hologram_instance" "example" {
  zone_id =  "${data.alibabacloudstack_zones.default.zones.0.id}"
  instance_name = "${var.name}"
  compute_type = "Standard"
  cpu = "intel"
  node = 2
  cluster = "${data.alibabacloudstack_hologram_clusters.default.clusters.0.id}"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id = "${data.alibabacloudstack_vpc_vswitch.default.id}"
}

resource "alibabacloudstack_hologram_instance_backup_policy" "example" {
  instance_id        = "${alibabacloudstack_hologram_instance.example.id}"
  hour               = 2
  data_keep_quantity = 7
  week               = ["0", "2", "4", "6"]
}
```

## 参数说明

以下参数是支持的：

* `instance_id` - （必选，不可变）Hologram实例的ID。
* `hour` - （必选）执行备份的小时时间。有效值：0-23。
* `data_keep_quantity` - （必选）备份数据保留天数。有效值：1-31。
* `week` - （必选）执行备份的星期几。有效值："0"（星期日），"1"（星期一），"2"（星期二），"3"（星期三），"4"（星期四），"5"（星期五），"6"（星期六）。
* `enabled` - （可选）是否启用备份策略。默认值为`true`。

## 属性导出

以下属性会被导出：

* `id` - 资源ID，与`instance_id`相同。

## 导入

可以使用id（instance_id）导入Hologram实例备份策略，例如：

```shell
$ terraform import alibabacloudstack_hologram_instance_backup_policy.example example-instance-id
```