---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_disks"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-disks"
description: |- 
  查询云服务器磁盘
---

# alibabacloudstack_ecs_disks
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_disks`

根据指定过滤条件列出当前凭证权限可以访问的云服务器磁盘列表。

## 示例用法

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

variable "name" {
  default = "tf-testAccCheckAlibabacloudStackDisksDataSource_ids-75584"
}

resource "alibabacloudstack_ecs_disk" "default" {
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  category         = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  name             = "${var.name}"
  description      = "${var.name}_description"
  size             = "20"
  tags = {
    Name  = "TerraformTest"
    Name1 = "TerraformTest"
  }
}

data "alibabacloudstack_ecs_disks" "default" {
  ids = ["${alibabacloudstack_ecs_disk.default.id}"]
  name_regex = "tf-testAccCheckAlibabacloudStackDisksDataSource_ids-75584"
  type       = "offline"
  category   = "cloud_ssd"
  instance_id = "i-bp1234567890abcdefg"

  tags = {
    Name  = "TerraformTest"
    Name1 = "TerraformTest"
  }

  output_file = "disks_output.txt"
}

output "first_disk_id" {
  value = "${data.alibabacloudstack_ecs_disks.default.disks.0.id}"
}
```

## 参数参考

以下参数是支持的：

* `ids` - (可选) 磁盘ID列表。用于精确筛选磁盘。
* `name_regex` - (可选，强制更新)用于通过磁盘名称筛选结果的正则表达式字符串。
* `type` - (可选，强制更新)扩容磁盘的方式。有效值：
  * `offline`(默认)：离线扩容。扩容后，您必须在控制台[重启实例](https://help.aliyun.com/document_detail/25440.html)或者调用API [RebootInstance](https://help.aliyun.com/document_detail/25502.html)使操作生效。
  * `online`：在线扩容，无需重启实例即可完成扩容。磁盘类型支持高效云盘、SSD云盘、ESSD云盘和弹性临时盘。
* `category` - (可选，强制更新)磁盘种类。有效值：
  * `all`：所有云盘以及本地盘。
  * `cloud`：普通云盘。
  * `cloud_efficiency`：高效云盘。
  * `cloud_ssd`：SSD盘。
  * `cloud_essd`：ESSD云盘。
  * `cloud_auto`：ESSD AutoPL云盘。
  * `local_ssd_pro`：I/O密集型本地盘。
  * `local_hdd_pro`：吞吐密集型本地盘。
  * `cloud_essd_entry`：ESSD Entry云盘。
  * `elastic_ephemeral_disk_standard`：弹性临时盘-标准版。
  * `elastic_ephemeral_disk_premium`：弹性临时盘-高级版。
  * `ephemeral`：(已停售)本地盘。
  * `ephemeral_ssd`：(已停售)本地SSD盘。
  默认值：`all`。
* `instance_id` - (可选，强制更新)创建一块包年包月磁盘，并自动挂载到指定的包年包月实例(InstanceId)上。
  * 设置实例ID后，会忽略您设置的ResourceGroupId、Tag.N.Key、Tag.N.Value、ClientToken和KMSKeyId参数。
  * 您不能同时指定ZoneId和InstanceId。
  默认值：空，代表创建的是按量付费云盘，云盘所属地由RegionId和ZoneId确定。
* `tags` - (可选) 分配给磁盘的标签映射。

## 属性参考

除了上述参数外，还导出以下属性：

* `disks` - 磁盘列表。每个元素包含以下属性：
  * `id` - 磁盘的ID。
  * `name` - 磁盘的名称。
  * `description` - 磁盘的描述。长度为2~256个英文或中文字符，不能以`http://`和`https://`开头。
  * `region_id` - 块存储所属的地域ID。您可以调用[DescribeRegions](https://help.aliyun.com/document_detail/25609.html)查看最新的阿里云地域列表。
  * `availability_zone` - 磁盘的可用区。
  * `status` - 磁盘状态。更多信息，请参见[云盘状态](https://help.aliyun.com/document_detail/25689.html)。取值范围：
    * `In_use`：使用中。
    * `Available`：待挂载。
    * `Attaching`：挂载中。
    * `Detaching`：卸载中。
    * `Creating`：创建中。
    * `ReIniting`：初始化中。
    * `All`：所有状态。默认值：`All`。
  * `type` - 扩容磁盘的方式。取值范围：
    * `offline`(默认)：离线扩容。扩容后，您必须在控制台[重启实例](https://help.aliyun.com/document_detail/25440.html)或者调用API [RebootInstance](https://help.aliyun.com/document_detail/25502.html)使操作生效。
    * `online`：在线扩容，无需重启实例即可完成扩容。磁盘类型支持高效云盘、SSD云盘、ESSD云盘和弹性临时盘。
  * `category` - 磁盘种类。取值范围同`category`参数的有效值。
  * `size` - 磁盘的大小(单位为GiB)。具体取值范围如下：
    * 系统盘：
      * 普通云盘：20~500。
      * ESSD云盘：
        * PL0：1~2048。
        * PL1：20~2048。
        * PL2：461~2048。
        * PL3：1261~2048。
      * ESSD AutoPL 云盘：1~2048。
      * 其他云盘类型：20~2048。
    * 数据盘：
      * 高效云盘(cloud_efficiency)：20~32768。
      * SSD云盘(cloud_ssd)：20~32768。
      * ESSD云盘(cloud_essd)：具体取值范围与`PerformanceLevel`的取值有关。可以调用[DescribeDisks](https://help.aliyun.com/document_detail/25514.html)查询云盘信息，再根据查询结果中的`PerformanceLevel`参数查看取值。
        * PL0：1~65536。
        * PL1：20~65536。
        * PL2：461~65536。
        * PL3：1261~65536。
      * 普通云盘(cloud)：5~2000。
      * ESSD AutoPL云盘(cloud_auto)：1~65536。
      * ESSD Entry云盘(cloud_essd_entry)：10~32768。
      * 弹性临时盘-标准版(elastic_ephemeral_disk_standard)：64～8,192。
      * 弹性临时盘-高级版(elastic_ephemeral_disk_premium)：64～8,192。
  * `image_id` - 创建磁盘所使用的镜像ID。除非磁盘是从镜像创建的，否则它为null。
  * `snapshot_id` - 创建磁盘所使用的快照。如果未使用快照创建磁盘，则它为null。
  * `instance_id` - 相关实例的ID。除非`status`为`In_use`，否则它为`null`。
  * `kms_key_id` - 对应于数据盘的KMS密钥ID。
  * `creation_time` - 磁盘的创建时间。
  * `attached_time` - 磁盘的附加时间。
  * `detached_time` - 磁盘的分离时间。
  * `storage_set_id` - 磁盘所属的存储集ID。
  * `expiration_time` - 磁盘的过期时间。
  * `tags` - 分配给磁盘的标签映射。