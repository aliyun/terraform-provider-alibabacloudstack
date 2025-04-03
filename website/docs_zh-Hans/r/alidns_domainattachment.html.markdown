---
subcategory: "DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_domain_attachment"
sidebar_current: "docs-alibabacloudstack-resource-dns-domain-attachment"
description: |-
  域名绑定到DNS实例

---

# alibabacloudstack_dns_domain_attachment
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_alidns_domainattachment`

使用Provider配置的凭证在指定的资源集下将域名绑定到DNS实例的资源。


## 示例用法

```
resource "alibabacloudstack_dns_domain_attachment" "dns" {
  instance_id     = "dns-cn-mp91lyq9xxxx"
  domain_names    = ["test111.abc", "test222.abc"]
}
```

## 参数参考

以下参数是支持的：

* `instance_id` - (必填，变更时重建) DNS实例的ID。
* `domain_names` - (必填) 绑定到DNS实例的域名列表。

## 属性参考

以下属性是导出的：

* `id` - 此资源的ID。其值与`instance_id`相同。
* `domain_names` - 绑定到DNS实例的域名列表。
* `instance_id` - DNS实例的ID。