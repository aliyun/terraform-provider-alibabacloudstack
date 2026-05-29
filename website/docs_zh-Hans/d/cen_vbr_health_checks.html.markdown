---
subcategory: "云企业网"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_vbr_health_checks"
sidebar_current: "docs-alibabacloudstack-datasource-cen-vbr-health-checks"
description: |-
  提供CEN VBR健康检查列表
---

# alibabacloudstack\_cen\_vbr\_health\_checks

该数据源根据指定的过滤条件提供阿里云账户中的CEN VBR健康检查列表。

## 示例用法

```terraform
variable "name" {
  default = "tf-testAccCenVbrHealthChecks"
}

resource "alibabacloudstack_cen_instance" "default" {
  cen_instance_name = var.name
  description       = var.name
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = "pc-123456"
  virtual_border_router_name = var.name
  vlan_id                    = 1991
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alibabacloudstack_cen_transit_router_vbr_attachment" "default" {
  vbr_id            = alibabacloudstack_express_connect_virtual_border_router.default.id
  cen_id            = alibabacloudstack_cen_instance.default.id
  transit_router_id = alibabacloudstack_cen_instance.default.transit_router_id
}

resource "alibabacloudstack_cen_vbr_health_check" "default" {
  cen_id                 = alibabacloudstack_cen_transit_router_vbr_attachment.default.cen_id
  vbr_instance_id        = alibabacloudstack_express_connect_virtual_border_router.default.id
  health_check_source_ip = "172.16.0.1"
  health_check_target_ip = "10.0.0.1"
  healthy_threshold      = 3
  health_check_interval  = 2
  health_check_only      = true
  depends_on             = [alibabacloudstack_cen_transit_router_vbr_attachment.default]
}

data "alibabacloudstack_cen_vbr_health_checks" "default" {
  cen_id          = alibabacloudstack_cen_vbr_health_check.default.cen_id
  vbr_instance_id = alibabacloudstack_cen_vbr_health_check.default.vbr_instance_id
  ids             = ["${alibabacloudstack_cen_vbr_health_check.default.cen_id}:${alibabacloudstack_cen_vbr_health_check.default.vbr_instance_id}"]
}
```

## 参数参考

以下参数被支持：

* `cen_id` - (可选, 强制新资源) CEN实例的ID。
* `vbr_instance_id` - (可选, 强制新资源) VBR实例的ID。
* `ids` - (可选, 强制新资源) VBR健康检查ID列表，格式为`<cen_id>:<vbr_instance_id>`。

## 属性参考

以下属性被导出：

* `ids` - VBR健康检查ID列表。
* `vbr_health_checks` - VBR健康检查列表。每个元素包含以下属性：
  * `cen_id` - CEN实例的ID。
  * `vbr_instance_id` - VBR实例的ID。
  * `vbr_instance_region_id` - VBR实例的区域ID。
  * `health_check_source_ip` - 健康检查的源IP地址。
  * `health_check_target_ip` - 健康检查的目标IP地址。
  * `health_check_interval` - 健康检查的间隔。
  * `healthy_threshold` - 链路被认定为健康之前的连续健康检查成功次数。
  * `health_check_only` - 是否仅执行健康检查。
```