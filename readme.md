# certbot 华为云自动续期

certbot-huaweicloud

**服务器迁移到1Panel。此项目不再维护。**

---

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

## 编译

```bash
go build -o ./bin/ ./...
```

可执行程序将生成在bin目录下。

## 使用

```bash
# 不编译，debug用
go run main.go -use [auth|cleanup]

# 编译后使用
/path/to/project/bin/certbot.exe -use [auth|cleanup]
```

use 参数为 auth 时，创建recordSet。

use 参数为 cleanup 时，删除recordSet。

如果是手动设置记录集的主机和值，需要在`.env`中设置`CERTBOT_DOMAIN`和`CERTBOT_VALIDATION`，并将`MODE`设置为DEBUG。

如果由certbot自动续期，则将`MODE`设置为DEV并编译。

## 自动续期

```bash
certbot renew --manual --preferred-challenges=dns \
--manual-auth-hook /path/to/project/bin/certbot -use auth \
--manual-cleanup-hook /path/to/project/bin/certbot -use cleanup \
--deploy-hook "sudo nginx -s reload"
```

- `/path/to/project/bin/certbot` 为你编译生成的程序路径。
- `sudo nginx -s reload` 用于重新加载配置nginx配置。

编辑生成的程序需要读取配置文件。cd到项目目录或者将配置文件`.env`放到程序同一目录再执行上述命令。
