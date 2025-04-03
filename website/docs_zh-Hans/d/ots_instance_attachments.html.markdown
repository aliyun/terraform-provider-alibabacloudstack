---
subcategory: "OTS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ots_instance_attachments"
sidebar_current: "docs-Alibabacloudstack-datasource-ots-instance_attachments"
description: |- 
  查询表格存储（OTS）实例关联
---

# alibabacloudstack_ots_instance_attachments

根据指定过滤条件列出当前凭证权限可以访问的表格存储（OTS）实例关联列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAcc49077"
}

resource "alibabacloudstack_ots_instance" "foo" {
  name        = "${var.name}"
  description = "${var.name}"
  accessed_by = "Vpc"
  instance_type = "Capacity"
}

data "alibabacloudstack_zones" "foo" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "foo" {
  cidr_block = "172.16.0.0/16"
  name       = "${var.name}"
}

resource "alibabacloudstack_vswitch" "foo" {
  vpc_id            = "${alibabacloudstack_vpc.foo.id}"
  name              = "${var.name}"
  cidr_block        = "172.16.1.0/24"
  availability_zone = "${data.alibabacloudstack_zones.foo.zones.0.id}"
}

resource "alibabacloudstack_ots_instance_attachment" "foo" {
  instance_name = "${alibabacloudstack_ots_instance.foo.name}"
  vpc_name      = "testvpc"
  vswitch_id    = "${alibabacloudstack_vswitch.foo.id}"
}

data "alibabacloudstack_ots_instance_attachments" "default" {
  instance_name = "${alibabacloudstack_ots_instance_attachment.foo.instance_name}"
  name_regex    = "testvpc"
  output_file   = "attachments.txt"

  # 输出第一个附件的ID
  output "first_ots_attachment_id" {
    value = "${data.alibabacloudstack_ots_instance_attachments.default.attachments.0.id}"
  }
}
```

## 参数参考

以下参数是支持的：

* `instance_name` - (必填) OTS实例名称，用于筛选与该实例相关的附件。
* `name_regex` - (可选) 正则表达式字符串，用于通过VPC名称进一步过滤结果。

## 属性参考

除了上述参数外，还导出以下属性：

* `names` - 与OTS实例关联的VPC名称列表。
* `vpc_ids` - 与OTS实例关联的VPC ID列表。
* `attachments` - 实例附件列表。每个元素包含以下属性：
  * `id` - 资源ID，通常与`instance_name`相同。
  * `domain` - 实例附件的域信息。
  * `endpoint` - 实例附件的访问端点。
  * `region` - 实例附件所属的区域。
  * `instance_name` - OTS实例的名称。
  * `vpc_name` - 绑定到OTS实例的VPC名称。
  * `vpc_id` - 绑定到OTS实例的VPC ID。
```