---
subcategory: "ROS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ros_template"
sidebar_current: "docs-Alibabacloudstack-ros-template"
description: |- 
  编排资源编排（ROS）模板
---

# alibabacloudstack_ros_template

使用Provider配置的凭证在指定的资源集编排资源编排（ROS）模板。

## 示例用法

### 基础用法

以下示例展示了如何通过 `template_body` 参数直接定义 ROS 模板内容：

```terraform
variable "name" {
    default = "tf-testaccrostemplate20209"
}

resource "alibabacloudstack_ros_template" "basic" {
  template_name = "MyTemplateTest12"
  description   = "这是一个测试的 ROS 模板。"
  template_body = <<EOF
    {
      "ROSTemplateFormatVersion": "2015-09-01",
      "Parameters": {
        "InstanceType": {
          "Type": "String",
          "Description": "ECS 实例的实例类型。"
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

### 高级用法 - 使用 `template_url`

以下示例展示了如何通过 `template_url` 参数从外部 URL 加载模板内容：

```terraform
resource "alibabacloudstack_ros_template" "from_url" {
  template_name = "MyTemplateFromURL"
  description   = "此模板从指定的 URL 加载。"
  template_url  = "http://example.com/path/to/template.json"
}
```

## 参数参考

支持以下参数：

* `description` - (可选) 模板的描述。描述长度最多为256个字符。
* `template_body` - (可选) 包含模板主体的结构。模板主体必须为1到524,288字节之间。如果模板主体长度过长，建议将参数添加到HTTP POST请求体中，以避免因URL长度过长导致请求失败。必须指定 `template_body` 或 `template_url` 中的一个，但不能同时指定两者。
* `template_name` - (必填) 模板的名称。名称最多可以包含255个字符，并且可以包含数字、字母、连字符(`-`)和下划线(`_`)。它必须以数字或字母开头。
* `template_url` - (可选) 包含模板主体的文件的URL。该URL必须指向位于Web服务器(HTTP或HTTPS)或阿里云OSS上的存储空间(例如，`oss://ros/stack-policy/demo`，`oss://ros/stack-policy/demo?RegionId=cn-hangzhou`)。模板的最大长度为524,288字节。如果未指定OSS区域，则默认使用接口的 `RegionId`。只能指定 `template_body` 或 `template_url` 中的一个。
* `tags` - (可选) 要分配给资源的标签映射。

## 属性参考

除了上述所有参数外，还导出以下属性：

* `id` - ROS模板资源的ID。此值等同于 `template_id`。
* `template_id` - ROS模板的唯一标识符。