---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_random_password"
sidebar_current: "docs-alibabacloudstack-resource-random-password"
description: |-
    提供一个随机密码以创建ECS实例。
---

# extra docs:alibabacloudstack ecs instance with random password

此资源使用加密随机数生成器，为创建ECS实例提供随机密码。

## 示例用法

```
resource "random_password" "password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
output "password" {  
  value = random_password.password.result  
}

resource "alibabacloudstack_instance" "default" {
  availability_zone = "cn-wulan-env212-amtest212001-a"
  instance_type = "ecs.s7-k-c1m1.large"
  system_disk_category = "cloud_efficiency"
  system_disk_size = 20
  system_disk_name = "test"
  data_disks{
    name = "disk1"
    category = "cloud_efficiency"
    size = 20
    delete_with_instance = false
    encrypted = false 
  }
  security_groups = [alibabacloudstack_security_group.default.id]
  password = random_password.password.result
  description = "description"
  tags = {
      testYx = "123"
    }
  image_id = "anolisos_7_9_x64_20G_anck_alibase_20220727.vhd"
  vswitch_id = alibabacloudstack_vswitch.default.id
  instance_name = "test7"
  is_outdated = false
}
```

## 参数参考

以下参数被支持：
* `length` - (必填) 期望的字符串长度。length 的最小值为 1，并且 length 必须也 >= (min_upper + min_lower + min_numeric + min_special)。
* `keepers` - (可选) 任意值映射，当其更改时，将触发资源的重新创建。有关更多信息，请参阅主提供商文档。
* `lower` - (可选) 在结果中包含小写字母字符。默认值为 true。
* `min_lower` - (可选) 结果中的最小小写字母字符数。默认值为 0。
* `min_numeric` - (可选) 结果中的最小小写数字字符数。默认值为 0。
* `min_special` - (可选) 结果中的最小特殊字符数。默认值为 0。
* `min_upper` - (可选) 结果中的最小大写字母字符数。默认值为 0。
* `number` - (可选) 在结果中包含数字字符。默认值为 true。注意：此参数已弃用，改用 numeric。
* `numeric` - (可选) 在结果中包含数字字符。默认值为 true。
* `override_special` - (可选) 提供您自己的特殊字符列表用于字符串生成。这会覆盖特殊参数中的默认字符列表。special 参数仍必须设置为 true，以便任何重写的字符在生成中使用。
* `special` - (可选) 在结果中包含特殊字符。这些是 !@#$%&*()-_=+[]{}<>:?。默认值为 true。
* `upper` - (可选) 在结果中包含大写字母字符。默认值为 true。

## 属性参考

以下属性被导出：

* `bcrypt_hash` - 生成的随机字符串的 bcrypt 哈希。
* `id` - Terraform 内部使用的静态值，不应在配置中引用。
* `result` - 生成的随机字符串。

## 导入

支持使用以下语法导入：

```bash
$ terraform import random_password.password securepassword
```