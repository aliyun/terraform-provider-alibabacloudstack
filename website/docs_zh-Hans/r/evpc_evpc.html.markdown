---
subcategory: "EasyAI"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_evpc_evpc"
sidebar_current: "docs-Alibabacloudstack-evpc-evpc"
description: |- 
  编排EVPC资源
---

# alibabacloudstack_evpc_evpc

编排EVPC（Elastic Virtual Private Cloud，弹性虚拟私有云）资源。

## 示例用法

```hcl
resource "alibabacloudstack_evpc_evpc" "default" {
  evpc_name   = "my-evpc"
  description = "我的EVPC实例"
}
```

## 参数说明

支持以下参数：

* `evpc_name` - (必填) EVPC的名称。长度必须为2到128个字符。
* `description` - (选填) EVPC的描述。长度必须为0到256个字符。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - EVPC的ID。
* `evpc_id` - EVPC的ID。
* `evpc_name` - EVPC的名称。
* `status` - EVPC的状态。
* `description` - EVPC的描述。
* `cidr` - EVPC的CIDR块。
* `tenant_id` - EVPC的租户ID。
* `department` - EVPC的部门ID。
* `department_name` - EVPC的部门名称。
* `region_id` - EVPC的地域ID。
* `resource_group` - EVPC的资源组ID。
* `resource_group_name` - EVPC的资源组名称。
* `cluster_id` - EVPC的集群ID。
* `ascm_create_user` - 创建EVPC的用户。
* `create_time` - EVPC的创建时间。
* `update_time` - EVPC的最后更新时间。
