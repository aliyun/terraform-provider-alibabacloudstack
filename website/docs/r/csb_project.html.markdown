---
subcategory: "CSB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_csb_project"
sidebar_current: "docs-alibabacloudstack-resource-csb-project"
description: |-
  Provides a Alibabacloudstack resource to manage CSB Project .
---

# alibabacloudstack_csb_project

This resource will help you to manager CSB Project.

For information about CSB Project and how to use it, see [Create a Project](https://help.aliyun.com/apsara/enterprise/v_3_17_0_30393230/csb/apsarastack-developer-guide/obtains-information-about-a-single-service-group.html?spm=a2c4g.14484438.10001.97)



-> **NOTE:** You need to set your registry password in CSB Project console before use this resource.

## Example Usage

Basic Usage

```

resource "alibabacloudstack_csb_project" "project" {
 "data":         "{\\\"projectName\\\":\\\"test17\\\",\\\"projectOwnerName\\\":\\\"test17\\\",\\\"projectOwnerEmail\\\":\\\"\\\",\\\"projectOwnerPhoneNum\\\":\\\"\\\",\\\"description\\\":\\\"\\\"}",
 "csb_id":       "134",
 "project_name": "test17",
}
```

## Argument Reference

The following arguments are supported:

* `data` - (Optional) Infomation of CSB Project. 
* `csb_id` - (Required, ForceNew) id of  CSB instance  where repository is created. 
* `project_name` - (Required, ForceNew) Name of CSB Project. It can contain 2 to 64 characters.


## Attributes Reference

The following attributes are exported:

* `csb_id` - The id of CSB instance. 
* `project_name` - The project name of CSB Project.
* `project_owner_name` - The project owner name of CSB Project.
* `gmt_modified` - The project modified time of CSB Project.
* `gmt_create` - The project create time of CSB Project.
* `owner_id` - The owner id of CSB Project.
* `api_num` - The api num of CSB Project.
* `user_id` - The user id of CSB Project.
* `delete_flag` - The delete flag of CSB Project.
* `cs_id` - The project id of CSB Project.
* `status` - The project status of CSB Project.
* `data` - Infomation of CSB Project. 