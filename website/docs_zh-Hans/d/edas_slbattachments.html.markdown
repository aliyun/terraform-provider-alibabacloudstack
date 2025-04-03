---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_slbattachments"
sidebar_current: "docs-Alibabacloudstack-datasource-edas-slbattachments"
description: |- 
    查询企业级分布式应用服务负载均衡挂载
---

# alibabacloudstack_edas_slbattachments
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_edas_applications`

根据指定过滤条件列出当前凭证权限可以访问的企业级分布式应用服务负载均衡挂载列表。

## 示例用法
```hcl
variable "name" {
  default = "tf-testacc-edas-applications2798"
}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  name       = "${var.name}"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "cn-beijing-a"
  name             = "${var.name}"
}

resource "alibabacloudstack_slb" "default" {
  name          = "${var.name}"
  vswitch_id    = "${alibabacloudstack_vswitch.default.id}"
  address_type  = "internet"
  specification = "slb.s1.small"
}

resource "alibabacloudstack_edas_cluster" "default" {
  cluster_name = "${var.name}"
  cluster_type = 2
  network_mode = 2
  vpc_id       = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_edas_application" "default" {
  application_name = "${var.name}"
  cluster_id      = "${alibabacloudstack_edas_cluster.default.id}"
  package_type    = "JAR"
}

resource "alibabacloudstack_edas_slb_attachment" "default" {
  app_id         = "${alibabacloudstack_edas_application.default.id}"
  slb_id         = "${alibabacloudstack_slb.default.id}"
  slb_ip         = "${alibabacloudstack_slb.default.address}"
  type           = "${alibabacloudstack_slb.default.address_type}"
  listener_port  = 80
}

data "alibabacloudstack_edas_slbattachments" "default" {
  ids = ["${alibabacloudstack_edas_slb_attachment.default.id}"]
  name_regex = "${alibabacloudstack_edas_slb_attachment.default.app_id}"
  output_file = "slbattachments_output.txt"
}
```

## 参数参考
以下参数是支持的：
  * `ids` - （可选，变更时重建）SLB挂载ID列表。用于通过特定的SLB挂载ID筛选结果。
  * `name_regex` - （可选，变更时重建）用于按名称筛选结果的正则表达式字符串。当您想查找符合特定命名模式的SLB挂载时，这可能非常有用。
  
## 属性参考
除了上述参数外，还导出以下属性：
  * `names` - SLB挂载名称列表。
  * `applications` - 与SLB挂载关联的应用程序列表。列表中的每个元素都是一个包含以下键的映射：
    - `app_id` - 应用程序的ID。
    - `slb_id` - SLB的ID。
    - `slb_ip` - SLB的IP地址。
    - `type` - SLB的类型（例如，“internet”或“intranet”）。
    - `listener_port` - SLB的监听端口。
    - `vserver_group_id` - 与SLB关联的VServer组的ID。
    - `slb_status` - SLB的状态。
    - `vswitch_id` - 与SLB关联的VSwitch的ID。