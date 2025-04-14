---
subcategory: "Elasticsearch"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_elasticsearch_k8s_instance"
sidebar_current: "docs-Alibabacloudstack-elasticsearch-k8s-instance"
description: |-
  编排Kubernetes的Elasticsearch实例
---

# alibabacloudstack_elasticsearch_k8s_instance

使用Provider配置的凭证在指定的资源集下编排Kubernetes的Elasticsearch实例。

## 示例用法

```hcl
resource "alibabacloudstack_vpc" "default" {
  name       = "tf-testacc-vpc"
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "cn-hangzhou-e"
  name              = "tf-testacc-vswitch"
}

resource "alibabacloudstack_elasticsearch_k8s_instance" "default" {
  description            = "tf-testacc-elasticsearch-on-k8s"
  vswitch_id             = "${alibabacloudstack_vswitch.default.id}"
  version                = "6.7.0_with_X-Pack"
  instance_charge_type   = "PostPaid"
  period                 = 1
  data_node_amount       = 3
  data_node_spec         = "elasticsearch.sn2.medium"
  data_node_disk_size    = 20
  data_node_disk_type    = "cloud_efficiency"
  master_node_spec       = "elasticsearch.sn2.medium"
  client_node_amount     = 2
  client_node_spec       = "elasticsearch.sn2.medium"
  enable_public          = true
  public_whitelist       = ["0.0.0.0/0"]
  enable_kibana_public_network = true
  kibana_whitelist       = ["0.0.0.0/0"]
  zone_count             = 1
  resource_group_id      = "rg-acfmyu465pju6decvztf"
  setting_config         = {
    "cluster.routing.allocation.disk.watermark.low" = "85%"
    "cluster.routing.allocation.disk.watermark.high" = "90%"
  }
  tags = {
    Name = "tf-testacc-elasticsearch-on-k8s"
  }
}
```

## 参数说明

以下是支持的参数：

* `description` - (可选) - Elasticsearch 实例的描述。长度必须为 0 到 30 个字符，可以包含数字、字母、下划线、短划线 (-) 和中文字符。必须以字母、数字或中文字符开头。
* `vswitch_id` - (必填, 变更时重建) - VSwitch 的 ID。
* `password` - (敏感, 可选) - Elasticsearch 实例的密码。password 和 kms_encrypted_password 必须设置其中一个。
* `kms_encrypted_password` - (可选) - 使用 KMS 加密的 Elasticsearch 实例密码。password 和 kms_encrypted_password 必须设置其中一个。
* `kms_encryption_context` - (可选) - 用于解密 kms_encrypted_password 的加密上下文映射。
* `version` - (必填, 变更时重建) - Elasticsearch 实例的版本。
* `tags` - (可选) - 要分配给资源的标签映射。
* `instance_charge_type` - (可选) - 实例的计费方式。有效值：PrePaid(包年包月)和 PostPaid(按量付费)。默认是 PostPaid。
* `period` - (可选) - 实例的订阅周期。有效值：1、2、3、4、5、6、7、8、9、12、24、36。如果 instance_charge_type 是 PrePaid，则此参数是必填的。
* `data_node_amount` - (必填) - 数据节点的数量。有效范围：2 到 50。
* `data_node_spec` - (必填) - 数据节点的规格。
* `data_node_disk_size` - (必填) - 数据节点的磁盘大小，单位为 GB。
* `data_node_disk_type` - (必填) - 数据节点的磁盘类型。有效值：cloud_efficiency、cloud_ssd、cloud_essd。
* `data_node_disk_encrypted` - (可选, 变更时重建) - 是否对数据节点磁盘进行加密。默认为 false。
* `private_whitelist` - (可选) - Elasticsearch 实例的私有 IP 白名单集合。
* `enable_public` - (可选) - 是否启用公网访问。默认为 false。
* `public_whitelist` - (可选) - Elasticsearch 实例的公网 IP 白名单集合。
* `master_node_spec` - (可选) - 主节点的规格。
* `client_node_amount` - (可选) - 客户端节点的数量。有效范围：2 到 25。
* `client_node_spec` - (可选) - 客户端节点的规格。
* `protocol` - (可选) - Elasticsearch 实例使用的协议。有效值：HTTP 和 HTTPS。默认为 HTTP。
* `zone_count` - (可选, 变更时重建) - 可用区数量。有效范围：1 到 3。默认为 1。
* `resource_group_id` - (可选, 变更时重建) - 资源组的 ID。
* `setting_config` - (可选) - Elasticsearch 设置的映射。

## 属性说明

除了上述参数外，还导出以下属性：

* `domain` - Elasticsearch 实例的域名。
* `port` - Elasticsearch 实例的端口号。
* `status` - Elasticsearch 实例的状态。
* `kibana_domain` - Kibana 实例的域名。
* `kibana_port` - Kibana 实例的端口号。
* `enable_kibana_public_network` - 是否启用了 Kibana 公网。
* `kibana_whitelist` - Kibana 实例的公网 IP 白名单集合。
* `enable_kibana_private_network` - 是否启用了 Kibana 私网。
* `kibana_private_whitelist` - Kibana 实例的私有 IP 白名单集合。
