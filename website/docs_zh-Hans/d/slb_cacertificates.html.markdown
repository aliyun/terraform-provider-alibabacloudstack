---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_cacertificates"
sidebar_current: "docs-Alibabacloudstack-datasource-slb-cacertificates"
description: |- 
  查询负载均衡(SLB)CA证书
---

# alibabacloudstack_slb_cacertificates
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_slb_ca_certificates`

根据指定过滤条件列出当前凭证权限可以访问的负载均衡(SLB)CA证书列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccSlbCACertificatesDataSourceBasic-15703"
}

resource "alibabacloudstack_slb_ca_certificate" "default" {
  name              = "${var.name}"
  ca_certificate    = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq******bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
}

data "alibabacloudstack_slb_ca_certificates" "default" {
  name_regex = "${alibabacloudstack_slb_ca_certificate.default.name}"
}

output "certificate_id" {
  value = "${data.alibabacloudstack_slb_ca_certificates.default.certificates.0.id}"
}

output "certificate_name" {
  value = "${data.alibabacloudstack_slb_ca_certificates.default.certificates.0.name}"
}

output "certificate_fingerprint" {
  value = "${data.alibabacloudstack_slb_ca_certificates.default.certificates.0.fingerprint}"
}

output "certificate_region_id" {
  value = "${data.alibabacloudstack_slb_ca_certificates.default.certificates.0.region_id}"
}
```

## 参数参考

以下参数是支持的：

* `ids` - (可选) CA证书ID列表，用于过滤结果。  
* `name_regex` - (可选，变更时重建) 用于通过CA证书名称过滤结果的正则表达式字符串。  

## 属性参考

除了上述参数外，还导出以下属性：

* `names` - SLB CA证书名称列表。  
* `certificates` - SLB CA证书列表。每个元素包含以下属性：  
  * `id` - CA证书的ID。  
  * `name` - CA证书的名称。  
  * `fingerprint` - CA证书的指纹。  
  * `created_timestamp` - CA证书创建的时间戳。  
  * `region_id` - CA证书所在的区域ID。  
```