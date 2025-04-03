---
subcategory: "ROS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ros_template"
sidebar_current: "docs-Alibabacloudstack-ros-template"
description: |- 
  Provides a ros Template resource.
---

# alibabacloudstack_ros_template

Provides a ROS Template resource.

For information about ROS Template and how to use it, see [What is Template](https://www.alibabacloud.com/help/en/doc-detail/141851.htm).

## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_ros_template" "example" {
  template_name = "MyTemplateTest12"
  description   = "This is a test ROS template."
  template_body = <<EOF
    {
      "ROSTemplateFormatVersion": "2015-09-01",
      "Parameters": {
        "InstanceType": {
          "Type": "String",
          "Description": "The instance type of the ECS instance."
        }
      },
      "Resources": {
        "EcsInstance": {
          "Type": "ALIYUN::ECS::Instance",
          "Properties": {
            "InstanceType": { "Ref": "InstanceType" }
          }
        }
      }
    }
    EOF
}
```

Advanced Usage with `template_url`

```terraform
resource "alibabacloudstack_ros_template" "example_with_url" {
  template_name = "MyTemplateFromURL"
  description   = "This template is loaded from a URL."
  template_url  = "http://example.com/path/to/template.json"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the template. The description can be up to 256 characters in length.
* `template_body` - (Optional) The structure that contains the template body. The template body must be 1 to 524,288 bytes in length. If the length of the template body is longer than required, we recommend that you add parameters to the HTTP POST request body to avoid request failures due to excessive length of URLs. You must specify one of the `template_body` or `template_url` parameters, but you cannot specify both of them.
* `template_name` - (Required) The name of the template. The name can be up to 255 characters in length and can contain digits, letters, hyphens (`-`), and underscores (`_`). It must start with a digit or letter.
* `template_url` - (Optional) The URL of the file that contains the template body. The URL must point to the storage space located on the Web server (HTTP or HTTPS) or Alibaba Cloud OSS (for example, `oss://ros/stack-policy/demo`, `oss://ros/stack-policy/demo?RegionId=cn-hangzhou`). The maximum length of the template is 524,288 bytes. If the OSS region is not specified, the `RegionId` of the interface is used by default. You must specify only one of the `template_body` or `template_url` parameters.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the ROS Template resource. This value is equivalent to `template_id`.
* `template_id` - The unique identifier for the ROS Template.