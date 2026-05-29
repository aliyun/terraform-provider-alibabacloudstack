---
subcategory: "硬件安全模块"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_hsm_cluster"
sidebar_current: "docs-alibabacloudstack-resource-hsm-cluster"
description: |-
  提供阿里云专有云HSM密码机集群资源。
---

# alibabacloudstack\_hsm\_cluster

提供一个HSM 密码机集群资源。

## 示例用法

### 创建HSM集群

```hcl
variable "name" {
  default = "tf_hsm_cluster6821719"
}

data "alibabacloudstack_zones" "default" {
  provider                    = alibabacloudstack-common
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  provider   = alibabacloudstack-common
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  provider     = alibabacloudstack-common
  vswitch_name = "${var.name}_vsw"
  vpc_id       = alibabacloudstack_vpc.default.id
  cidr_block   = "172.16.1.0/24"
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
}

data "alibabacloudstack_hsm_vendors" "default" {
}

resource "random_password" "password" {
  count            = 1
  length           = 12
  special          = true
  override_special = "!@#$^&*()_"
  min_lower        = 1
  min_upper        = 1
  min_numeric      = 1
}

resource "alibabacloudstack_hsm_instance" "default0" {
  product_code = data.alibabacloudstack_hsm_vendors.default.vendors.0.products.0.code
  vendor_code  = data.alibabacloudstack_hsm_vendors.default.vendors.0.code
  vsm_type     = "gvsm"
  zone_no      = data.alibabacloudstack_zones.default.zones.0.id
  remark       = var.name
}

resource "alibabacloudstack_hsm_instance" "default1" {
  product_code = data.alibabacloudstack_hsm_vendors.default.vendors.0.products.0.code
  vendor_code  = data.alibabacloudstack_hsm_vendors.default.vendors.0.code
  vsm_type     = "gvsm"
  zone_no      = data.alibabacloudstack_zones.default.zones.0.id
  remark       = var.name
}

resource "alibabacloudstack_hsm_cluster" "example" {
  cluster_name        = var.name
  master_instance_id  = alibabacloudstack_hsm_instance.default0.id
  vpc_id              = alibabacloudstack_vpc.default.id
  vswitch_ids         = alibabacloudstack_vswitch.default.id
  zone_nos            = data.alibabacloudstack_zones.default.zones.0.id
  password            = random_password.password.0.result
  ip_white_list       = "10.0.0.0/8"
  
  sub_instance_ids = [
    alibabacloudstack_hsm_instance.default1.id,
  ]
}
```

## 参数参考

以下参数被支持：

* `cluster_name` - (必需) HSM集群名称。长度为2-128个字符，必须以英文字母或数字开头，支持英文、数字、中文、下划线(_)和短横线(-)。
* `master_instance_id` - (必需, 强制新建) 主HSM实例ID。
* `vpc_id` - (必需, 强制新建) HSM集群所在VPC的ID。
* `vswitch_ids` - (必需, 强制新建) 与VPC关联的交换机ID。
* `zone_nos` - (必需, 强制新建) HSM集群所在的可用区。
* `password` - (必需) HSM集群管理员密码。密码长度为8-30个字符，至少包含大写字母、小写字母、数字和特殊字符(!@#$%^&*()_+-=)中的三种。
* `ip_white_list` - (可选) HSM集群的IP白名单。多个IP地址用逗号(,)分隔，白名单支持通配符(*)。默认值：0.0.0.0/0。
* `sub_instance_ids` - (可选) 属于该集群的子HSM实例ID列表。

## 属性参考

以下属性会被导出：

* `id` - HSM集群的ID。

## 导入

可以使用id导入HSM集群，例如：

```shell
$ terraform import alibabacloudstack_hsm_cluster.example cl-abc123456
```