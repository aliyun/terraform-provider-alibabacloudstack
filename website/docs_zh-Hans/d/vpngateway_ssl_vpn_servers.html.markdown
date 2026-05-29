---
subcategory: "VPN网关"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_sslvpnclientcerts"
sidebar_current: "docs-Alibabacloudstack-datasource-vpngateway-sslvpnclientcerts"
description: |-
  提供当前 AlibabacloudStack 账户下所有 SSL VPN 客户端证书的列表。
---

# alibabacloudstack\_vpngateway\_sslvpnclientcerts

该数据源用于根据指定过滤条件获取 VPNGateway 下的 SSL VPN 客户端证书列表。

## 示例用法

```hcl
variable "name" {
		default = "tf_testAccVpngatewaySslvpnclientcertDataSource_7150678"
	}
	
	
data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}


resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}


resource "alibabacloudstack_vpn_gateway" "default" {
 name                 = "${var.name}"
 vpc_id               = "${alibabacloudstack_vpc_vpc.default.id}"
 bandwidth            = 10
 instance_charge_type = "PostPaid"
 enable_ssl           = true
 enable_ipsec		  = true
 vswitch_id			  = "${alibabacloudstack_vpc_vswitch.default.id}"
}

	
	resource "alibabacloudstack_vpngateway_ssl_vpnserver" "default" {
		client_ip_pool =  "10.8.0.0/24"
	
		local_subnet =  "192.168.1.0/24"
	
		ssl_vpn_server_name =  "${var.name}"
		vpn_gateway_id =     "${alibabacloudstack_vpn_gateway.default.id}"
		}

resource "alibabacloudstack_vpngateway_sslvpnclientcert" "default" {
	ssl_vpn_client_cert_name= "${var.name}"

	ssl_vpn_server_id= "${alibabacloudstack_vpngateway_ssl_vpnserver.default.id}"

	vpn_gateway_id= "${alibabacloudstack_vpn_gateway.default.id}"
}



data "alibabacloudstack_vpngateway_sslvpnclientcerts" "default" {
  ids = [
          "${alibabacloudstack_vpngateway_sslvpnclientcert.default.id}"
        ]
}
```
## 参数说明
以下参数支持配置：

* `ids` - (可选) 客户端证书 ID 列表。若设置此参数，数据源将返回匹配的客户端证书信息。
* `name_regex` - (可选, ForceNew) 用于通过名称模式匹配客户端证书的正则表达式。可以使用该参数筛选出符合特定命名规则的证书。
## 属性说明
除上述参数外，还将导出以下属性：

* `ssl_vpn_client_certs` - SSL VPN 客户端证书列表
* `id` - SSL VPN 客户端证书的 ID
* `ca_cert` - CA 证书内容。
* `client_cert` - 客户端证书内容。
* `client_config` - 客户端配置文件。
* `client_key` - 客户端私钥。
* `create_time` - SSL 客户端证书创建时间。
* `end_time` - SSL 客户端证书过期时间。
* `ssl_vpn_client_cert_id` - SSL 客户端证书的唯一标识 ID。
* `ssl_vpn_client_cert_name` - 客户端证书的名称。
* `ssl_vpn_server_id` - 所属 SSL 服务器的 ID。
* `status` - 客户端证书的当前状态。