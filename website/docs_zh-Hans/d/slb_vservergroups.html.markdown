---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_vservergroups"
sidebar_current: "docs-Alibabacloudstack-datasource-slb-vservergroups"
description: |- 
  查询负载均衡(SLB)虚拟服务组
---

# alibabacloudstack_slb_vservergroups
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_slb_server_groups`

根据指定过滤条件列出当前凭证权限可以访问的负载均衡(SLB)虚拟服务组列表。

## 示例用法

```hcl
data "alibabacloudstack_slb_vservergroups" "sample_ds" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  ids             = ["vsg-12345678", "vsg-abcdefg"]
  name_regex      = "^group-.*"

  output_file = "slb_vservergroups_output.txt"
}

output "first_slb_vserver_group_id" {
  value = data.alibabacloudstack_slb_vservergroups.sample_ds.slb_server_groups[0].id
}
```

## 参数说明

以下参数是支持的：

* `load_balancer_id` - (必填) 负载均衡实例的ID。
* `ids` - (可选) 用于过滤结果的SLB VServer组ID列表。此参数允许用户通过指定的VServer组ID来筛选结果。
* `name_regex` - (可选，变更时重建) 用于通过VServer组名称过滤结果的正则表达式字符串。此参数允许用户通过名称模式匹配来筛选VServer组。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - SLB VServer组ID列表。此属性包含所有符合条件的VServer组的ID。
* `names` - SLB VServer组名称列表。此属性包含所有符合条件的VServer组的名称。
* `slb_server_groups` - SLB VServer组列表。每个元素包含以下属性：
  * `id` - VServer组ID。此属性标识唯一的VServer组。
  * `name` - VServer组名称。此属性表示VServer组的名称。
  * `servers` - 与该组关联的ECS实例列表。每个元素包含以下属性：
    * `instance_id` - ECS实例的ID。此属性标识附加到VServer组的ECS实例。
    * `port` - 后端服务器使用的端口号。此属性表示ECS实例在VServer组中使用的端口。
    * `weight` - ECS实例的权重。此属性表示ECS实例在负载均衡中的权重值。