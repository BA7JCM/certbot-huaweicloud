# certbot 华为云自动续期

certbot-huaweicloud

## 参考

[华为云开发者 Go 软件开发工具包（Go SDK）](https://console.huaweicloud.com/apiexplorer/#/sdkcenter/DNS?lang=Go)
[云解析服务 DNS API概览](https://support.huaweicloud.com/api-dns/zh-cn_topic_0132421999.html)
[certbot 验证前和验证后钩子](https://eff-certbot.readthedocs.io/en/stable/using.html#pre-and-post-validation-hooks)

## 安装SDK

```bash
# 安装华为云 Go SDK 库
go get github.com/huaweicloud/huaweicloud-sdk-go-v3
```

## 参数配置

将.env.example重命名为.env，并修改配置参数

- `CLOUD_SDK_AK`和`CLOUD_SDK_SK`
    华为云访问密钥，参见：[访问密钥](https://support.huaweicloud.com/usermanual-ca/zh-cn_topic_0046606340.html)

- `REGION`
    在华为云控制台查看所在区域名称。
    在[地区与终端节点 云解析服务](https://console.huaweicloud.com/apiexplorer/#/endpoint/DNS)查看对应区域变量名。

- `DOMAIN`
    你要申请证书的域名。

- `DESCRIPTION`
    创建recordSet时的描述信息，非必填。

- `CERTBOT_DOMAIN`和`CERTBOT_VALIDATION`
    DEBUG时才填，验证sdk是否正常。
    正式使用时，由certbot 自动生成，用于验证域名所有权。

- `MODE`
    DEBUG | DEV
    DEBUG时，将使用配置文件中的参数，用于验证SDK是否正常。
    正式编译时要改成DEV。
