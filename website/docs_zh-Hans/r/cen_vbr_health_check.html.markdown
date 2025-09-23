---
subcategory: "CEN"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_vbr_health_check"
sidebar_current: "docs-alibabacloudstack-resource-cen-vbr-health-check"
description: |-
  提供阿里云CEN VBR健康检查资源
---

# alibabacloudstack\_cen\_vbr\_health\_check

提供CEN VBR健康检查资源，用于为CEN VBR配置健康检查。

## 示例用法

```terraform
variable "name" {
  default = "tf-testAcc9527"
}


resource "alibabacloudstack_cen_instance" "default" {
  cen_instance_name = var.name
  description       = var.name
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = "pc123456"
  virtual_border_router_name = var.name
  vlan_id                    = 1991
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alibabacloudstack_cen_transit_router_vbr_attachment" "default" {
    vbr_id = "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
	cen_id = "${alibabacloudstack_cen_instance.default.id}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}

resource "alibabacloudstack_cen_vbr_health_check" "default" {
    cen_id                 = "${alibabacloudstack_cen_transit_router_vbr_attachment.default.cen_id}"
    vbr_instance_id        = "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
    health_check_source_ip = "172.16.0.1"
    health_check_target_ip = "10.0.0.1"
    healthy_threshold      = 3
    health_check_interval  = 2
    health_check_only      = true
    depends_on             = [alibabacloudstack_cen_transit_router_vbr_attachment.default]
}

```

## 参数参考

以下参数被支持：

* `cen_id` - (必填, 强制新资源) CEN实例的ID。
* `health_check_source_ip` - (必填, 强制新资源) 健康检查的源IP地址。
* `health_check_target_ip` - (必填, 强制新资源) 健康检查的目标IP地址。
* `health_check_interval` - (必填, 强制新资源) 健康检查的间隔。有效值：1到3。
* `healthy_threshold` - (必填, 强制新资源) 链路被认定为健康之前的连续健康检查成功次数。有效值：1到8。
* `health_check_only` - (必填, 强制新资源) 是否仅执行健康检查。有效值：`true`和`false`。
* `vbr_instance_id` - (必填, 强制新资源) VBR实例的ID。

## 属性参考

以下属性被导出：

* `id` - 资源的ID，格式为`<cen_id>:<vbr_instance_id>`。
* `vbr_instance_owner_id` - VBR实例的所有者ID。
* `link_status` - 链路状态。
* `delay` - 健康检查的延迟。
* `packet_loss` - 健康检查的丢包率。

## 导入

可以使用id导入CEN VBR健康检查，例如：

```bash
$ terraform import alibabacloudstack_cen_vbr_health_check.default cen-abc123456:vbr-def789012
```