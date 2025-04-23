---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_servercertificate"
sidebar_current: "docs-Alibabacloudstack-slb-servercertificate"
description: |- 
  编排负载均衡(SLB)服务器证书
---

# alibabacloudstack_slb_servercertificate
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_slb_server_certificate`

使用Provider配置的凭证在指定的资源集编排负载均衡(SLB)服务器证书。

## 示例用法

### 示例 1：使用server_certificate/private_key作为字符串内容

```hcl
# 创建一个使用字符串内容的服务器证书
resource "alibabacloudstack_slb_servercertificate" "example" {
  name                = "slbservercertificate-example"
  server_certificate  = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgI+OuMs******XTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key         = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0knDrlNdiys******ErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
  server_certificate_name = "example-cert-name"
}
```

### 示例 2：使用server_certificate/private_key从文件中读取

```hcl
# 创建一个使用文件内容的服务器证书
resource "alibabacloudstack_slb_servercertificate" "example" {
  name                = "slbservercertificate-example"
  server_certificate  = file("${path.module}/server_certificate.pem")
  private_key         = file("${path.module}/private_key.pem")
  server_certificate_name = "example-cert-name"
}
```

### 示例 3：完整示例

```hcl
variable "name" {
    default = "tf-testaccslbserver_certificate70376"
}

resource "alibabacloudstack_slb_servercertificate" "default" {
  server_certificate  = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgI+OuMs******XTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key         = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0knDrlNdiys******ErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
  server_certificate_name = "test-cert-name"
}
```

## 参数参考

支持以下参数：

* `name` - (选填) - 服务器证书的名称。如果未提供，Terraform将自动生成一个唯一名称。
* `server_certificate_name` - (选填) - 服务器证书的名称。这可以用来在SLB服务中标识该证书。
* `server_certificate` - (必填, 变更时重建) - 需要上传的公钥证书。如果不使用阿里云托管的证书，则此参数是必填的。
* `private_key` - (必填, 变更时重建) - 对应于`server_certificate`中指定的公钥证书的私钥。如果不使用阿里云托管的证书，则此参数是必填的。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 服务器证书(SSL Certificate)的ID。
* `name` - 服务器证书的名称。
* `server_certificate_name` - 创建时指定的服务器证书的名称。