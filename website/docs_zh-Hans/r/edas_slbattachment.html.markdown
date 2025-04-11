---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_slbattachment"
sidebar_current: "docs-Alibabacloudstack-edas-slbattachment"
description: |- 
  编排绑定企业级分布式应用服务（Edas）应用和负载均衡
---

# alibabacloudstack_edas_slbattachment
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_edas_slb_attachment`

使用Provider配置的凭证在指定的资源集下编排绑定企业级分布式应用服务（Edas）应用和负载均衡。

## 示例用法

### 基础用法

```hcl
resource "alibabacloudstack_edas_slbattachment" "default" {
  app_id           = var.app_id
  slb_id           = var.slb_id
  slb_ip           = var.slb_ip
  type             = var.type
  listener_port    = var.listener_port
  vserver_group_id = var.vserver_group_id
  slb_status       = var.slb_status
  vswitch_id       = var.vswitch_id
}
```

## 参数说明

支持以下参数：

* `app_id` - (必填，变更时重建) 应用的ID。
* `slb_id` - (必填，变更时重建) SLB实例的ID。
* `slb_ip` - (必填，变更时重建) 分配给绑定SLB实例的IP地址。
* `type` - (必填，变更时重建) SLB实例的网络类型。有效值：
  * `internet`: 外网实例。
  * `intranet`: 内网实例。
* `listener_port` - (选填，变更时重建) 绑定SLB实例的监听端口。
* `vserver_group_id` - (选填，变更时重建) 与内网SLB实例关联的虚拟服务器(VServer)组的ID。
* `slb_status` - (变更时重建) SLB实例的运行状态。可能的值包括：
  * `Inactive`: 实例已停止，监听器将不监控或转发流量。
  * `Active`: 实例正在运行。实例创建后，默认状态为活动状态。
  * `Locked`: 实例被锁定，通常是由于欠费或被阿里云锁定。
  * `Expired`: 实例已过期。
* `vswitch_id` - (变更时重建) SLB实例所属的VPC中的交换机ID。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 资源的唯一标识符，格式为 `<app_id>:<slb_id>`。
* `slb_status` - 当前SLB实例的运行状态。
* `vswitch_id` - 与SLB实例关联的VPC中的交换机ID。