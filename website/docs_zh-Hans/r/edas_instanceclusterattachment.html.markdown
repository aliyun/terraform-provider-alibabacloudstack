---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_instanceclusterattachment"
sidebar_current: "docs-Alibabacloudstack-edas-instanceclusterattachment"
description: |- 
  编排绑定企业级分布式应用服务（Edas）实例和集群
---

# alibabacloudstack_edas_instance_cluster_attachment
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_edas_instanceclusterattachment`

使用Provider配置的凭证在指定的资源集下编排绑定企业级分布式应用服务（Edas）实例和集群。

## 示例用法

```hcl
variable "name" {
  default = "tf-testacc-edasicattachment19002"
}

variable "password" {
  default = "Li65272237###"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
} 

resource "alibabacloudstack_vpc" "default" {
  name        = var.name
  cidr_block  = "10.1.0.0/21"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id              = alibabacloudstack_vpc.default.id
  cidr_block          = "10.1.1.0/24"
  availability_zone   = data.alibabacloudstack_zones.default.zones[0].id
  name               = var.name
}

resource "alibabacloudstack_security_group" "default" {
  name       = var.name
  description= "New security group"
  vpc_id     = alibabacloudstack_vpc.default.id
}

resource "alibabacloudstack_instance" "default" {
  vswitch_id         = alibabacloudstack_vswitch.default.id
  image_id           = "centos_7_7_x64_20G_alibase_20200426.vhd"
  availability_zone  = data.alibabacloudstack_zones.default.zones[0].id
  system_disk_category = "cloud_efficiency"
  system_disk_size  = 60
  instance_type      = "ecs.n4v2.xlarge"

  security_groups    = [alibabacloudstack_security_group.default.id]
  instance_name      = var.name
  tags = {
    Name = "TerraformTest-instance"
  }
}

resource "alibabacloudstack_edas_cluster" "default" {
  cluster_name = var.name
  cluster_type = 2
  network_mode = 2
  vpc_id       = alibabacloudstack_vpc.default.id
}

resource "alibabacloudstack_edas_instance_cluster_attachment" "default" {
  cluster_id = alibabacloudstack_edas_cluster.default.id
  instance_ids = [
                   alibabacloudstack_instance.default.id
                 ]
  pass_word = var.password
}
```

## 参数参考

支持以下参数：
  * `cluster_id` - (必填, 变更时重建) - 要附加实例的集群ID。
  * `instance_ids` - (必填, 变更时重建) - 将附加到指定集群的ECS实例ID列表。
  * `pass_word` - (必填, 变更时重建) - 集群中ECS实例的登录密码。在附加过程中需要此密码以确保安全访问。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `id` - 资源的唯一标识符，由 `<cluster_id>:<instance_id1,instance_id2,...>` 组成。
  * `status_map` - 表示集群中每个实例状态的映射。键是实例ID，值表示状态：`1`(运行中)，`0`(转换中)，`-1`(失败)，和 `-2`(离线)。
  * `ecu_map` - 将每个实例链接到其对应的ECU(弹性计算单元)的映射。键是实例ID，值是ECU ID。
  * `cluster_member_ids` - 与每个实例关联的集群成员ID的映射。键是实例ID，值是集群成员ID。