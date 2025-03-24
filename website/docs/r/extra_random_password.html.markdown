---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_random_password"
sidebar_current: "docs-alibabacloudstack-resource-random-password"
description: |-
    Provides a random password  to the create ecs.
---

# alibabacloudstack_random_password

This resource does use a cryptographic random number generator.provides a random password for creating ECS instances

## Example Usage

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

## Argument Reference

The following arguments are supported:
* `length ` - (Required) The length of the string desired. The minimum value for length is 1 and, length must also be >= (min_upper + min_lower + min_numeric + min_special).
* `keepers` - (Optional) Arbitrary map of values that, when changed, will trigger recreation of resource. See the main provider documentation for more information.
* `lower` - (Optional) Include lowercase alphabet characters in the result. Default value is true.
* `min_lower` - (Optional) Minimum number of lowercase alphabet characters in the result. Default value is 0.
* `min_numeric` - (Optional) Minimum number of numeric characters in the result. Default value is 0.
* `min_special` - (Optional) Minimum number of special characters in the result. Default value is 0.
* `min_upper` - (Optional) Minimum number of uppercase alphabet characters in the result. Default value is 0.
* `number` - (Optional) Include numeric characters in the result. Default value is true. NOTE: This is deprecated, use numeric instead.
* `numeric` - (Optional) Include numeric characters in the result. Default value is true.
* `override_special` - (Optional) Supply your own list of special characters to use for string generation. This overrides the default character list in the special argument. The special argument must still be set to true for any overwritten characters to be used in generation.
* `special` - (Optional) Include special characters in the result. These are !@#$%&*()-_=+[]{}<>:?. Default value is true.
* `upper` - (Optional) Include uppercase alphabet characters in the result. Default value is true.
 
## Attributes Reference

The following attributes are exported:

* `bcrypt_hash` - A bcrypt hash of the generated random string.
* `id` -A static value used internally by Terraform, this should not be referenced in configurations.
* `result`  -The generated random string.


## Import  

Import is supported using the following syntax:

```
$terraform import random_password.password securepassword
```