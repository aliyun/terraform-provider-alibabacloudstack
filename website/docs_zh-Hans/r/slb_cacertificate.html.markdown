---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_cacertificate"
sidebar_current: "docs-Alibabacloudstack-slb-cacertificate"
description: |- 
  编排负载均衡(SLB)CA证书
---

# alibabacloudstack_slb_cacertificate
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_slb_ca_certificate`

使用Provider配置的凭证在指定的资源集编排负载均衡(SLB)CA证书。

## 示例用法

### 使用CA证书内容

```hcl
variable "name" {
    default = "tf-testaccslbca_certificate67317"
}

resource "alibabacloudstack_slb_cacertificate" "default" {
  name              = var.name
  ca_certificate    = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgI+OuMs******XTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
}
```

### 使用CA证书文件

```hcl
resource "alibabacloudstack_slb_cacertificate" "file_example" {
  name           = "tf-testaccslbca_certificate_file"
  ca_certificate = file("${path.module}/ca_certificate.pem")
}
```

## 参数说明

支持以下参数：

* `name` - (可选，已弃用) CA证书的名称。**此字段已弃用**，请使用 `ca_certificate_name` 代替。与 `ca_certificate_name` 参数互斥。
* `ca_certificate_name` - (可选) CA证书的名称。与 `name` 参数互斥。
* `ca_certificate` - (必填，变更时重建) PEM格式的CA证书内容。此字段是不可变的，修改此参数将强制创建新资源。

### 参数详细说明

- **name**: (可选，**已弃用**) 指定CA证书的名称。**此字段已弃用**，建议使用 `ca_certificate_name` 代替。如果未提供，则默认使用系统生成的名称。与 `ca_certificate_name` 互斥。
- **ca_certificate_name**: (可选) 用于标识CA证书的名称。与 `name` 互斥，不能同时设置。
- **ca_certificate**: (必填) PEM格式的CA证书内容。此字段必须在资源创建时提供，并且一旦创建后无法更改。请确保提供的证书内容符合PEM格式要求。

## 属性说明

除了上述所有参数外，还导出以下属性：

* `id` - CA证书的唯一标识符(ID)，用于唯一标识该资源。
* `name` - 创建时提供的CA证书名称，可以用于其他资源或配置中引用该证书。
* `ca_certificate_name` - CA证书的名称，对于在其他资源或配置中引用该证书非常有用。

### 属性详细说明

- **id**: CA证书的唯一标识符，即API返回的 `CACertificateId`，用于标识和管理该资源。此属性在创建资源后自动填充，用户无需手动设置。
- **name**: CA证书的名称。
- **ca_certificate_name**: CA证书的名称。

## Import

SLB CA证书可以使用 `CACertificateId` 进行导入，例如：

```
$ terraform import alibabacloudstack_slb_cacertificate.example CACertificateId
```

-> **注意:** `ca_certificate` 参数无法通过导入恢复，因为API不会回显该字段。导入后可能会产生非空的执行计划。