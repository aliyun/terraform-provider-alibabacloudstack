---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_servercertificates"
sidebar_current: "docs-Alibabacloudstack-datasource-slb-servercertificates"
description: |- 
  查询负载均衡(SLB)服务器证书
---

# alibabacloudstack_slb_servercertificates
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_slb_server_certificates`

根据指定过滤条件列出当前凭证权限可以访问的负载均衡(SLB)服务器证书列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccSlbServerCertificatesDataSourceBasic-18598"
}

resource "alibabacloudstack_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDdjCCAl4CCQCcm*******XgthAiFFjl1S9ZgdA6Zc=\n-----END CERTIFICATE-----"
  private_key        = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQ******7l3xC00BL7Z+SAJyI4QKA\n-----END RSA PRIVATE KEY-----"
}

data "alibabacloudstack_slb_server_certificates" "default" {
  name_regex = "${alibabacloudstack_slb_server_certificate.default.name}"

  output_file = "server_certificates_output.txt"
}

output "first_certificate_id" {
  value = data.alibabacloudstack_slb_server_certificates.default.certificates.0.id
}
```

## 参数参考

以下参数是支持的：

* `ids` - (可选) 用于过滤结果的 SLB 服务器证书 ID 列表。可以通过指定一个或多个证书 ID 来筛选特定的证书。
* `name_regex` - (可选，变更时重建) 用于通过 SLB 服务器证书名称过滤结果的正则表达式字符串。例如，可以使用 `example_cert.*` 来匹配所有以 `example_cert` 开头的证书名称。

## 属性参考

除了上述所有参数外，还导出以下属性：

* `names` - 匹配的 SLB 服务器证书名称列表。
* `certificates` - 匹配的 SLB 服务器证书列表。每个元素包含以下属性：
  * `id` - SLB 服务器证书的唯一标识符(ID)。
  * `name` - SLB 服务器证书的名称。
  * `fingerprint` - 服务器证书的指纹，用于验证证书的真实性。
  * `created_time` - 服务器证书的创建时间，格式为人类可读的时间(如 `2023-01-01 12:00:00`)。
  * `created_timestamp` - 服务器证书的创建时间戳，格式为 Unix 时间戳(如 `1672569600`)。

