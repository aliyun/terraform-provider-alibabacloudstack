---
subcategory: "Elasticsearch"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_elasticsearch_k8s_instance"
sidebar_current: "docs-Alibabacloudstack-elasticsearch-k8s-instance"
description: |-
  Provides a Elasticsearch on K8s instance resource.
---

# alibabacloudstack_elasticsearch_k8s_instance

Provides a Elasticsearch on K8s instance resource.

## Example Usage

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

## Argument Reference
The following arguments are supported:

* `description` - (Optional) - The description of the Elasticsearch instance. It must be 0 to 30 characters in length and can contain numbers, letters, underscores, underscores (_), and hyphens (-). It must start with a letter, a number, or a Chinese character.
* `vswitch_id` - (Required, ForceNew) - The ID of the VSwitch.
* `password` - (Sensitive, Optional) - The password for the Elasticsearch instance. One of password or kms_encrypted_password must be set.
* `kms_encrypted_password` - (Optional) - An KMS-encrypted password for the Elasticsearch instance. One of password or kms_encrypted_password must be set.
* `kms_encryption_context` - (Optional) - A map of encryption context used to decrypt the kms_encrypted_password.
* `version` - (Required, ForceNew) - The version of the Elasticsearch instance.
* `tags` - (Optional) - A mapping of tags to assign to the resource.
* `instance_charge_type` - (Optional) - The billing method of the instance. Valid values: PrePaid (subscription) and PostPaid (pay-as-you-go). Default is PostPaid.
* `period` - (Optional) - The subscription period of the instance. Valid values: 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36. Required if instance_charge_type is PrePaid.
* `data_node_amount` - (Required) - The number of data nodes. Valid range: 2 to 50.
* `data_node_spec` - (Required) - The specification of the data nodes.
* `data_node_disk_size` - (Required) - The disk size of the data nodes in GB.
* `data_node_disk_type` - (Required) - The disk type of the data nodes. Valid values: cloud_efficiency, cloud_ssd, cloud_essd.
* `data_node_disk_encrypted` - (Optional, ForceNew) - Whether the data node disks are encrypted. Default is false.
* `private_whitelist` - (Optional) - A set of private IP whitelists for the Elasticsearch instance.
* `enable_public` - (Optional) - Whether to enable public access. Default is false.
* `public_whitelist` - (Optional) - A set of public IP whitelists for the Elasticsearch instance.
* `master_node_spec` - (Optional) - The specification of the master nodes.
* `client_node_amount` - (Optional) - The number of client nodes. Valid range: 2 to 25.
* `client_node_spec` - (Optional) - The specification of the client nodes.
* `protocol` - (Optional) - The protocol used by the Elasticsearch instance. Valid values: HTTP and HTTPS. Default is HTTP.
* `zone_count` - (Optional, ForceNew) - The number of availability zones. Valid range: 1 to 3. Default is 1.
* `resource_group_id` - (Optional, ForceNew) - The ID of the resource group.
* `setting_config` - (Optional) - A map of Elasticsearch settings.


## Attributes Reference
The following attributes are exported in addition to the arguments listed above:

* `domain` - The domain name of the Elasticsearch instance.
* `port` - The port number of the Elasticsearch instance.
* `status` - The status of the Elasticsearch instance.
* `kibana_domain` - The domain name of the Kibana instance.
* `kibana_port` - The port number of the Kibana instance.
* `enable_kibana_public_network` - Whether the Kibana public network is enabled.
* `kibana_whitelist` - A set of public IP whitelists for the Kibana instance.
* `enable_kibana_private_network` - Whether the Kibana private network is enabled.
* `kibana_private_whitelist` - A set of private IP whitelists for the Kibana instance.
