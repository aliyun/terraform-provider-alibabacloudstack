---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_disk"
sidebar_current: "docs-Alibabacloudstack-ecs-disk"
description: |- 
  编排云服务器（Ecs）磁盘
---

# alibabacloudstack_ecs_disk
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_disk`

使用Provider配置的凭证在指定的资源集下编排云服务器（Ecs）磁盘。

## 示例用法
```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_disk" "default" {
  zone_id           = "${data.alibabacloudstack_zones.default.zones.0.id}"
  size              = 50
  disk_name         = "example-disk"
  description       = "This is an example disk"
  category          = "cloud_ssd"
  encrypted         = true
  encrypt_algorithm = "aes-256"
  delete_auto_snapshot = true
  delete_with_instance = false
  enable_auto_snapshot = true
}
```

## 参数参考

支持以下参数：

* `zone_id` - (选填, 变更时重建) - 创建按量付费磁盘的可用区 ID。如果您不设置 `InstanceId`，则 `ZoneId` 为必填参数。您不能同时指定 `ZoneId` 和 `InstanceId`。
* `name` - (选填) - 磁盘名称。长度为2~128个字符，必须包含字母或数字，可以包含短划线(-)、点(.)、下划线(_)。不能以短划线、点或下划线开头或结尾，且不能以 `http://` 或 `https://` 开头。默认值为空。
* `disk_name` - (选填) - 磁盘名称。长度为2~128个字符，支持Unicode中letter分类下的字符(其中包括英文、中文和数字等)。可以包含冒号(:)、下划线(_)、句号(.)或者短划线(-)。
* `description` - (选填) - 磁盘描述。长度为2~256个字符，不能以 `http://` 或 `https://` 开头。默认值为空。
* `category` - (选填, 变更时重建) - 磁盘种类。取值范围：
  - all：所有云盘以及本地盘。
  - cloud：普通云盘。
  - cloud_efficiency：高效云盘。
  - cloud_ssd：SSD云盘。
  - cloud_essd：ESSD云盘。
  - cloud_auto：ESSD AutoPL云盘。
  - local_ssd_pro：I/O密集型本地盘。
  - local_hdd_pro：吞吐密集型本地盘。
  - cloud_essd_entry：ESSD Entry云盘。
  - elastic_ephemeral_disk_standard：弹性临时盘-标准版。
  - elastic_ephemeral_disk_premium：弹性临时盘-高级版。
  - ephemeral：(已停售)本地盘。
  - ephemeral_ssd：(已停售)本地SSD盘。
  默认值：all。
* `size` - (选填) - 希望扩容到的磁盘容量大小。单位为GiB。取值范围如下：
  - 系统盘：
    - 普通云盘：20~500。
    - ESSD云盘：
      - PL0：1~2048。
      - PL1：20~2048。
      - PL2：461~2048。
      - PL3：1261~2048。
    - ESSD AutoPL 云盘：1~2048。
    - 其他云盘类型：20~2048。
  - 数据盘：
    - 高效云盘(cloud_efficiency)：20~32768。
    - SSD云盘(cloud_ssd)：20~32768。
    - ESSD云盘(cloud_essd)：具体取值范围与 `PerformanceLevel` 的取值有关。可以调用 [DescribeDisks](https://help.aliyun.com/document_detail/25514.html) 查询云盘信息，再根据查询结果中的 `PerformanceLevel` 参数查看取值。
      - PL0：1~32768。
      - PL1：20~32768。
      - PL2：461~32768。
      - PL3：1261~32768。
    - 普通云盘(cloud)：5~2000。
    - ESSD AutoPL云盘(cloud_auto)：1~32768。
    - ESSD Entry云盘(cloud_essd_entry)：10~32768。
    - 弹性临时盘-标准版(elastic_ephemeral_disk_standard)：64～8,192。
    - 弹性临时盘-高级版(elastic_ephemeral_disk_premium)：64～8,192。
  > 指定的新磁盘容量必须比原磁盘容量大。
* `snapshot_id` - (选填) - 创建云盘使用的快照。2013年7月15日及以前的快照不能用来创建云盘。`SnapshotId` 参数和 `Size` 参数存在以下限制：
  - 如果 `SnapshotId` 参数对应的快照容量大于设置的 `Size` 参数值，实际创建的云盘大小为指定快照的大小。
  - 如果 `SnapshotId` 参数对应的快照容量小于设置的 `Size` 参数值，实际创建的云盘大小为指定的 `Size` 参数值。
  - 不支持使用快照创建弹性临时盘。
* `kms_key_id` - (选填) - 云盘使用的KMS密钥ID。
* `encrypted` - (选填, 变更时重建) - 是否加密云盘。取值范围：
  - true：是。
  - false：否。
  默认值：false。
* `encrypt_algorithm` - (选填, 变更时重建) - 磁盘加密方式。有效值：
  - sm4-128
  - aes-256
  默认值：aes-256。
* `delete_auto_snapshot` - (选填) - 释放云盘时，是否会同时释放自动快照。取值范围：
  - true：是。
  - false：否。
  默认值：false。
* `delete_with_instance` - (选填) - 磁盘是否随实例释放。默认值：无，表示不改变当前的值。开启多重挂载特性的云盘，不支持设置该参数。在下列两种情况下，将参数 `DeleteWithInstance` 设置成 `false` 时会报错：
  - 磁盘的种类(category)为本地盘(ephemeral)时。
  - 磁盘的种类(category)为普通云盘(cloud)，且不可以卸载(Portable=false)时。
* `enable_auto_snapshot` - (选填) - 云盘是否启用自动快照策略功能。取值范围：
  - true：启用。
  - false：关闭。
  默认值：无，表示不改变当前的值。> 创建后的云盘默认启用自动快照策略功能。您只需要为云盘绑定自动快照策略即可正常使用。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `availability_zone` - 磁盘所隶属的可用区。
* `zone_id` - 在指定可用区内创建一块按量付费磁盘。如果您不设置 `InstanceId`，则 `ZoneId` 为必填参数。您不能同时指定 `ZoneId` 和 `InstanceId`。
* `name` - 磁盘名称。
* `disk_name` - 磁盘名称。长度为2~128个字符，支持Unicode中letter分类下的字符(其中包括英文、中文和数字等)。可以包含冒号(:)、下划线(_)、句号(.)或者短划线(-)。默认值：空。
* `status` - 磁盘状态。更多信息，请参见[云盘状态](https://help.aliyun.com/document_detail/25689.html)。取值范围：
  - In_use：使用中。
  - Available：待挂载。
  - Attaching：挂载中。
  - Detaching：卸载中。
  - Creating：创建中。
  - ReIniting：初始化中。
  - All：所有状态。
  默认值：All。
* `auto_snapshot_policy_id` - 根据自动快照策略ID查询云盘。
* `enable_automated_snapshot_policy` - 云盘是否设置了自动快照策略。取值范围：
  - true：已设置。
  - false：未设置。
  默认值：false。