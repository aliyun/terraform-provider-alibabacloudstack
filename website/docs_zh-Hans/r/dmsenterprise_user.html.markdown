---
subcategory: "DMSEnterprise"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dmsenterprise_user"
sidebar_current: "docs-Alibabacloudstack-dmsenterprise-user"
description: |- 
  编排企业版数据库管理用户
---

# alibabacloudstack_dmsenterprise_user
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_dms_enterprise_user`

使用Provider配置的凭证在指定的资源集下编排企业版数据库管理用户。

## 示例用法

```terraform
variable "name" {
    default = "tf-testaccdms_enterpriseuser93463"
}

resource "alibabacloudstack_dmsenterprise_user" "default" {
  mobile          = "11111111111"
  uid             = "265530631068325049"
  user_name       = "rdktest"
  role_names      = ["DBA"]
  status          = "NORMAL"
  max_execute_count = 100
  max_result_count  = 500
}
```

## 参数参考

支持以下参数：

* `uid` - (必填，变更时重建) - 用户的阿里云UID。此字段在创建后无法修改。
* `user_name` - (选填)- 用户的昵称。
* `mobile` - (选填)- 用户的钉钉号或手机号码。
* `role_names` - (选填)- 用户扮演的角色列表。例如：`["DBA"]`。
* `status` - (选填)- DMS Enterprise用户的状态。有效值为：`NORMAL`(正常)、`DISABLE`(禁用)。
* `max_execute_count` - (选填)- 用户当天允许执行的最大查询次数。
* `max_result_count` - (选填)- 用户当天可以查询的最大行数。
* `tid` - (选填)- 租户ID。取自系统右上角头像处悬停展示的租户ID信息，详情请参见[查看租户信息](https://www.alibabacloud.com/help/doc-detail/181330.htm)。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 用户的阿里云唯一ID (UID)。其值与`uid`相同。
* `nick_name` - 用户的昵称(已弃用，建议使用`user_name`代替)。
* `role_names` - 用户当前扮演的角色列表。
* `status` - DMS Enterprise用户的当前状态。
* `mobile` - 用户的钉钉号或手机号码。