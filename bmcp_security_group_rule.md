# alibabacloudstack_bmcp_security_group

> 裸金属算力平台安全组

## 资源Id

{sgrId}

## 创建资源

### 创建安全组规则
### api三元组信息
EasyAI_2023-11-01_AddSecurityGroupRule
### 请求参数示例
```json
{"Direction":"ingress","AccessPolicy":"ACCEPT","IpProtocol":"tcp","PortRange":"1/200","Priority":1,"CidrIp":"192.168.0.0/24","SgId":"sg-7bj7sq55wjnqhd7yccve","IpVersion":"ipv4","Description":"test",}
```

### 响应参数示例
```json
{
        "eagleEyeTraceId": "0aec110d17750381540333549e0026",
        "asapiSuccess": true,
        "responseVersion": "",
        "data": {
            "data": {
                "sgId": "sg-7bj7sq55wjnqhd7yccve",
                "sgrId": [
                    "sgr-kv0wfdfgzq6fbqtlbo61"
                ]
            }
        },
        "errorTitle": "",
        "requestId": "95C93B7F-4928-49A7-84AA-A7B5C52B177C",
        "success": true,
        "errorMessage": "",
        "errorCode": "",
        "status": 0
    }
```

## 读取资源

### api三元组信息

EasyAI_2023-11-01_ListSecurityGroupRule

### 请求参数示例
```json
{"Direction":"ingress","SgId":"sg-7bj7sq55wjnqhd7yccve"}
```

### 响应参数示例

```json
{
        "eagleEyeTraceId": "0aec110d17750375497667263e0025",
        "asapiSuccess": true,
        "responseVersion": "",
        "data": {
            "total": 1,
            "pageNumber": 1,
            "data": [
                {
                    "portRange": "1/200",
                    "SgId": "sg-7bj7sq55wjnqhd7yccve",
                    "description": "dsadsad",
                    "updateTime": "2026-04-01T16:48:53+08:00",
                    "ipProtocol": "tcp",
                    "priority": 1,
                    "relatedGroupId": "",
                    "ipVersion": "ipv4",
                    "deleted": false,
                    "createTime": "2026-04-01T16:48:53+08:00",
                    "ID": 14,
                    "accessPolicy": "ACCEPT",
                    "sgrId": "sgr-i2hzhpsem05dp43y8k0y",
                    "cidrIp": "192.168.0.0/24",
                    "direction": "ingress"
                }
            ],
            "pageSize": 10
        },
        "errorTitle": "",
        "requestId": "3B347C5F-7700-46B6-97F3-8E83C81C3773",
        "success": true,
        "errorMessage": "",
        "errorCode": "",
        "status": 0
    }
```

## 修改资源
### api三元组信息

EasyAI_2023-11-01_ModifySecurityGroupRule

### 请求参数示例
```json
{"Direction":"ingress","AccessPolicy":"ACCEPT","IpProtocol":"tcp","PortRange":"2/200","Priority":1,"CidrIp":"192.168.0.0/24","SgId":"sg-7bj7sq55wjnqhd7yccve","IpVersion":"ipv4","SgrId":"sgr-i2hzhpsem05dp43y8k0y","Description":"dsadsad"}
```

### 响应参数示例
不消费响应数据。

## 删除资源

### api三元组信息
EasyAI_2023-11-01_DeleteSecurityGroup
### 请求参数示例

```json
{"SgId":"sg-7bj7sq55wjnqhd7yccve","SgrId":"sgr-i2hzhpsem05dp43y8k0y"}
```

### 响应参数示例
不消费响应数据。

## 检索资源
### api三元组信息
EasyAI_2023-11-01_ListSecurityGroupRule
检索逻辑可以参考《读取资源》章节