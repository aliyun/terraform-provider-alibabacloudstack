---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_edas_namespace"
sidebar_current: "docs-Alibabacloudstack-edas-namespace"
description: |- 
  使用Provider配置的凭证在指定的资源集下编排企业级分布式应用服务（Edas） Namespace resource.
---

# alibabacloudstack_edas_namespace

使用Provider配置的凭证在指定的资源集下编排企业级分布式应用服务（Edas） Namespace resource.

For information about EDAS Namespace and how to use it, see [What is Namespace](https://www.alibabacloud.com/help/en/enterprise-distributed-application-service/latest/insertorupdateregion).

-> **NOTE:** Available since v3.16.5

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alibabacloudstack_edas_namespace&exampleId=34281039-bffb-a43d-3670-ce75c36528dc9c56a834&activeTab=example&spm=docs.r.edas_namespace.0.34281039bf&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alibabacloudstack" {
  region = var.region
}

variable "region" {
  default = "cn-hangzhou"
}
variable "name" {
  default = "tfexample"
}

resource "alibabacloudstack_edas_namespace" "default" {
  debug_enable         = false
  description          = var.name
  namespace_logical_id = "${var.region}:${var.name}"
  namespace_name       = var.name
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the namespace. It can be up to `128` characters in length.
* `namespace_logical_id` - (Required, ForceNew) The ID of the namespace.  
  - For custom namespaces, the format is `region ID:namespace identifier`, e.g., `cn-beijing:tdy218`.
  - For default namespaces, the format is just the `region ID`, e.g., `cn-beijing`.
* `namespace_name` - (Required) The name of the namespace. It can be up to `63` characters in length.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier (ID) of the namespace in Terraform.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Default to 1 minute) Used when creating the Namespace.
* `delete` - (Default to 1 minute) Used when deleting the Namespace.
* `update` - (Default to 1 minute) Used when updating the Namespace.