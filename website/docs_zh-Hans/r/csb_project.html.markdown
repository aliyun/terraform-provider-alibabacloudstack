---
subcategory: "CSB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_csb_project"
sidebar_current: "docs-alibabacloudstack-resource-csb-project"
description: |-
  编排CSB项目
---

# alibabacloudstack_csb_project

使用Provider配置的凭证在指定的资源集下编排CSB项目。

有关 CSB 项目的更多信息以及如何使用它，请参阅 [创建项目](https://help.aliyun.com/apsara/enterprise/v_3_17_0_30393230/csb/apsarastack-developer-guide/obtains-information-about-a-single-service-group.html?spm=a2c4g.14484438.10001.97)



-> **注意：** 在使用此资源之前，您需要在 CSB 项目控制台中设置您的注册密码。

## 示例用法

### 基础用法

```

resource "alibabacloudstack_csb_project" "project" {
 "data":         "{\\\"projectName\\\":\\\"test17\\\",\\\"projectOwnerName\\\":\\\"test17\\\",\\\"projectOwnerEmail\\\":\\\"\\\",\\\"projectOwnerPhoneNum\\\":\\\"\\\",\\\"description\\\":\\\"\\\"}",
 "csb_id":       "134",
 "project_name": "test17",
}
```

## 参数参考

支持以下参数：

* `data` - (可选) CSB 项目的详细信息。
* `csb_id` - (必填，变更时重建) 创建存储库的 CSB 实例的 id。
* `project_name` - (必填，变更时重建) CSB 项目的名称。它可以包含 2 到 64 个字符。


## 属性参考

导出以下属性：

* `csb_id` - CSB 实例的 id。
* `project_name` - CSB 项目的项目名称。
* `project_owner_name` - CSB 项目的项目所有者名称。
* `gmt_modified` - CSB 项目的项目修改时间。
* `gmt_create` - CSB 项目的项目创建时间。
* `owner_id` - CSB 项目的拥有者 id。
* `api_num` - CSB 项目的 API 数量。
* `user_id` - CSB 项目的用户 id。
* `delete_flag` - CSB 项目的删除标志。
* `cs_id` - CSB 项目的项目 id。
* `status` - CSB 项目的项目状态。
* `data` - CSB 项目的详细信息。