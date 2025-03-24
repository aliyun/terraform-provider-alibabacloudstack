---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_slbattachment"
sidebar_current: "docs-Alibabacloudstack-edas-slbattachment"
description: |- 
  使用Provider配置的凭证在指定的资源集下编排企业级分布式应用服务（Edas） Slbattachment resource.
---

# alibabacloudstack_edas_slbattachment
-> **NOTE:** Alias name has: `alibabacloudstack_edas_slb_attachment`

使用Provider配置的凭证在指定的资源集下编排企业级分布式应用服务（Edas） Slbattachment resource.

## Example Usage

Basic Usage

```hcl
resource "alibabacloudstack_edas_slbattachment" "default" {
  app_id           = var.app_id
  slb_id           = var.slb_id
  slb_ip           = var.slb_ip
  type             = var.type
  listener_port    = var.listener_port
  vserver_group_id = var.vserver_group_id
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The ID of the application to which you want to bind an SLB instance.
* `slb_id` - (Required, ForceNew) The ID of the SLB instance that is going to be bound.
* `slb_ip` - (Required, ForceNew) The IP address that is allocated to the bound SLB instance.
* `type` - (Required, ForceNew) The network type of the SLB instance. Valid values:
  * `internet`: Internet instance.
  * `intranet`: Intranet instance.
* `listener_port` - (Optional, ForceNew) The listening port for the bound SLB instance.
* `vserver_group_id` - (Optional, ForceNew) The ID of the virtual server (VServer) group associated with the intranet SLB instance.
* `slb_status` - (ForceNew) Running status of the SLB instance. Possible values include:
  * `Inactive`: The instance is stopped, and the listener will not monitor or forward traffic.
  * `Active`: The instance is running. After the instance is created, the default state is active.
  * `Locked`: The instance is locked, usually due to overdue payments or being locked by Alibaba Cloud.
  * `Expired`: The instance has expired.
* `vswitch_id` - (ForceNew) The ID of the VSwitch in the VPC to which the SLB instance belongs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the resource, formulated as `<app_id>:<slb_id>`.
* `slb_status` - The current running status of the SLB instance.
* `vswitch_id` - The ID of the VSwitch in the VPC associated with the SLB instance.
