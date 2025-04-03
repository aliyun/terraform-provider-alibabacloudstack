---
subcategory: "Elasticsearch"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_elasticsearch_instance"
sidebar_current: "docs-Alibabacloudstack-elasticsearch-instance"
description: |-
  Provides a Elasticsearch on K8s instance resource.
---

# alibabacloudstack_elasticsearch_instance

Provides a Elasticsearch nstance resource.


The `alibabacloudstack_elasticsearch_instance` resource allows you to manage Elasticsearch instances on Alibaba Cloud.

## Example Usage

```hcl
resource "alibabacloudstack_elasticsearch_instance" "example" {
  description            = "example_elasticsearch"
  vswitch_id             = "vsw-1234567890abcdef0"
  version                = "6.7_with_X-Pack"
  instance_charge_type   = "PostPaid"
  data_node_amount       = 2
  data_node_spec         = "elasticsearch.sn2ne.large"
  data_node_disk_size    = 20
  data_node_disk_type    = "cloud_efficiency"
  enable_public          = true
  public_whitelist       = ["0.0.0.0/0"]
  enable_kibana_public_network = true
  kibana_whitelist       = ["0.0.0.0/0"]
  zone_count             = 2
}
```

## Argument Reference
The following arguments are supported -

* `description` - (Optional) - The description of the Elasticsearch instance. It must be 0 to 30 characters in length and can contain numbers, letters, underscores, and hyphens. It must start with a letter, a number, or a Chinese character.
* `vswitch_id` - (Required) - The ID of the VSwitch in which to launch the Elasticsearch instance.
* `password` - (Optional) - The password for the Elasticsearch instance.
* `kms_encrypted_password` - (Optional) - An KMS encrypted password for the Elasticsearch instance.
* `kms_encryption_context` - (Optional) - A map of key-value pairs that can be used to provide additional context for the KMS encryption.
* `version` - (Required) - The version of Elasticsearch to deploy.
* `tags` - (Optional) - A mapping of tags to assign to the resource.
* `instance_charge_type` - (Optional) - The charge type of the instance. Valid values are PrePaid and PostPaid. Default value is PostPaid.
* `period` - (Optional) - The subscription period of the instance. Valid values are 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36. Default value is 1.
* `data_node_amount` - (Required) - The number of data nodes in the Elasticsearch cluster. Valid range is from 2 to 50.
* `data_node_spec` - (Required) - The specification of the data nodes.
* `data_node_disk_size` - (Required) - The disk size of the data nodes.
* `data_node_disk_type` - (Required) - The disk type of the data nodes.
* `data_node_disk_encrypted` - (Optional) - Whether the data node disk is encrypted. Default value is false.
* `private_whitelist` - (Optional) - The private IP whitelist for the Elasticsearch instance.
* `enable_public` - (Optional) - Whether to enable public access to the Elasticsearch instance. Default value is false.
* `public_whitelist` - (Optional) - The public IP whitelist for the Elasticsearch instance.
* `master_node_spec` - (Optional) - The specification of the master nodes.
* `client_node_amount` - (Optional) - The number of client nodes in the Elasticsearch cluster. Valid range is from 2 to 25.
* `client_node_spec` - (Optional) - The specification of the client nodes.
* `protocol` - (Optional) - The protocol used by the Elasticsearch instance. Valid values are HTTP and HTTPS. Default value is HTTP.
* `zone_count` - (Optional) - The number of zones in which to deploy the Elasticsearch instance. Valid range is from 1 to 3. Default value is 1.
* `resource_group_id` - (Optional) - The ID of the resource group to which the Elasticsearch instance belongs.
* `setting_config` - (Optional) - A map of settings to configure the Elasticsearch instance.

## Attributes Reference
The following attributes are exported -

* `domain` - The domain name of the Elasticsearch instance.
* `port` - The port number of the Elasticsearch instance.
* `status` - The status of the Elasticsearch instance.
* `kibana_domain` - The domain name of the Kibana instance.
* `kibana_port` - The port number of the Kibana instance.
* `enable_kibana_public_network` - Whether public access to the Kibana instance is enabled.
* `enable_kibana_private_network` - Whether private access to the Kibana instance is enabled.
* `kibana_whitelist` - The public IP whitelist for the Kibana instance.
* `kibana_private_whitelist` - The private IP whitelist for the Kibana instance.

## Import
Elasticsearch instances can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_elasticsearch_instance.example i-1234567890abcdef0
```