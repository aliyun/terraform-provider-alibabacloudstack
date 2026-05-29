---
subcategory: "高速通道"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_vbr_ha"
sidebar_current: "docs-Alibabacloudstack-expressconnect-vbr_ha"
description: |-
  高速通道 VBR 快速倒换组
---

# alibabacloudstack_expressconnect_vbr_ha

使用 Provider 编排高速通道 VBR 快速倒换组（VBR Failover Group）资源。

-> **注意:** 该资源也可以使用以下别名引用：
-> - `alibabacloudstack_vbr_ha`

## 示例用法

```hcl
variable "name" {
  default = "tf-testaccexpressconnect-vbrha1926"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default1" {
	physical_connection_id =     ""
	vlan_id =                    1926
	local_gateway_ip =           "10.0.0.1"
	peer_gateway_ip =            "10.0.0.2"
	peering_subnet_mask =        "255.255.255.252"
	virtual_border_router_name = "${var.name}_1"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default2" {
	physical_connection_id =     ""
	vlan_id =                    1929
	local_gateway_ip =           "10.1.0.1"
	peer_gateway_ip =            "10.1.0.2"
	peering_subnet_mask =        "255.255.255.252"
	virtual_border_router_name = "${var.name}_2"
}

resource "alibabacloudstack_expressconnect_vbr_ha" "default" {
  name        = var.name
  vbr_id      = alibabacloudstack_express_connect_virtual_border_router.default1.id
  peer_vbr_id = alibabacloudstack_express_connect_virtual_border_router.default2.id
  description = var.name
}
```

## 参数参考

支持以下参数：

* `name` -（必填，ForceNew）VBR 快速倒换组名称。
* `vbr_id` -（必填，ForceNew）VBR 实例 ID。
* `peer_vbr_id` -（必填，ForceNew）VBR 快速倒换组中另一个 VBR 的实例 ID。
* `description` -（选填，ForceNew）VBR 快速倒换组的描述信息。长度为 2～256 个字符，必须以字母或中文开头，但不能以 `http://` 或 `https://` 开头。

## 属性参考

除上述参数外，该资源还会导出以下属性：

* `id` - VBR 快速倒换组 ID。

## Import

高速通道 VBR 快速倒换组可以使用 VBR HA ID（例如 `vbrha-xxxxxxxxx`）进行导入，例如：

```
$ terraform import alibabacloudstack_expressconnect_vbr_ha.example vbrha-xxxxxxxxx
```
