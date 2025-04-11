---
subcategory: "ESS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ess_scalinggroup_vserver_groups"
sidebar_current: "docs-Alibabacloudstack-ess-scalinggroup-vserver-groups"
description: |-
  绑定虚拟服务组到指定的缩放组
---

# alibabacloudstack_ess_scalinggroup_vserver_groups

使用Provider配置的凭证在指定的资源集下编排绑定虚拟服务组到指定的缩放组。

-> **注意**：vserver 组所属的负载均衡器必须处于 `active` 状态。

-> **注意**：如果缩放组的网络类型为 `VPC`，则 vserver 组必须在同一个 `VPC` 中。
 
-> **注意**：默认情况下，一个缩放组最多可以绑定 5 个 vserver 组。

-> **注意**：vserver 组和负载均衡器的默认组共享相同的后端服务器配额。

-> **注意**：当将 vserver 组绑定到缩放组时，现有的 ECS 实例将被添加到 vserver 组；相反，ECS 实例将从 vserver 组中移除。

-> **注意**：解除操作将在绑定操作之前执行。

-> **注意**：vserver 组由 `loadbalancer_id`、`vserver_group_id` 和 `port` 唯一定义。

-> **注意**：修改 `weight` 属性意味着先解除 vserver 组的绑定，然后再以新的权重参数重新绑定。

## 示例用法

```hcl
resource "alibabacloudstack_ess_scalinggroup_vserver_groups" "default" {
  scaling_group_id = "your_scaling_group_id"
  vserver_groups {
    loadbalancer_id = "your_loadbalancer_id"
    vserver_attributes {
      vserver_group_id = "your_vserver_group_id"
      port             = 80
      weight           = 100
    }
  }
  vserver_groups {
    loadbalancer_id = "another_loadbalancer_id"
    vserver_attributes {
      vserver_group_id = "another_vserver_group_id"
      port             = 8080
      weight           = 200
    }
  }
  force = true
}
```

### 参数参考
以下是支持的参数：

* `scaling_group_id` - (必填，变更时重建) - 缩放组的 ID。
* `vserver_groups` - (必填) - 要绑定到缩放组的一组 VServer 组。
  * `loadbalancer_id` - (必填) - 负载均衡器的 ID。
  * `vserver_attributes` - (必填) - 一组 VServer 属性。
    * `vserver_group_id` - (必填) - VServer 组的 ID。
    * `port` - (必填) - VServer 组的端口号。
    * `weight` - (必填) - VServer 组的权重。修改此属性将导致先解除绑定，然后以新的权重重新绑定。
* `force` - (可选) - 是否强制绑定或解除绑定 VServer 组。默认值为 `true`。如果设置为 `true`，即使存在依赖关系，也会强制执行绑定或解除绑定操作。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `id` - (必填，变更时重建) - ESS vserver 组绑定资源的唯一标识符（ID）。该 ID 唯一标识了绑定关系，并在创建时生成。