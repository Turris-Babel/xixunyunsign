# 习讯云签到工具

[自动签到看这里](scheduler.md)

[serve子命令签到api详解](docs/webapi.md)

![License](https://img.shields.io/github/license/theshdowaura/xixunyunsign.svg)
![Version](https://img.shields.io/github/v/release/theshdowaura/xixunyunsign.svg)
![Build Status](https://github.com/theshdowaura/xixunyunsign/actions/workflows/release.yml/badge.svg)
![Run Status](https://github.com/theshdowaura/xixunyunsign/actions/workflows/test.yml/badge.svg)
![Docker Status](https://github.com/theshdowaura/xixunyunsign/actions/workflows/docker-publish.yml/badge.svg)

## 快速上手
[示例操作流程](#示例操作流程)

## 项目简介

习讯云签到工具是一个使用 Go 语言和 Cobra 库编写的命令行工具，用于与习讯云平台的 API 进行交互，实现账户登录、查询签到信息和执行签到等功能。该工具支持多用户操作，使用 SQLite 数据库存储用户的 Token 及签到位置信息。

## 功能介绍

- **登录**：使用账号和密码登录到习讯云平台，并保存会话 Token。
- **查询签到信息**：查询当月的签到信息，包括连续签到天数、应签到位置的经纬度等。
- **执行签到**：根据指定的地址、经纬度等信息执行签到操作。

## 环境要求

- Go 1.16 或更高版本
- 已安装 `Git`（用于克隆项目代码）
- 已安装 `gcc` 或其他 C 编译器（用于编译 `github.com/mattn/go-sqlite3` 库）

## 安装步骤

1. **克隆项目代码**

   ```bash
   git clone https://github.com/theshdowaura/xixunyunsign.git
   ```

2. **进入项目目录**

   ```bash
   cd xixunyunsign
   ```

3. **下载依赖库**

   ```bash
   go mod tidy
   ```

## 使用方法

### 编译项目

在使用之前，您需要先编译项目：

```bash
go build -o xixunyunsign.exe
```

这将在当前目录下生成一个可执行文件 `xixunyunsign.exe`。

### 命令概览

该工具提供以下命令：

- `login`：登录到习讯云平台。
- `query`：查询签到信息。
- `sign`：执行签到。
- `search`:通过学校名查询id

您可以使用 `xixunyunsign.exe help` 查看所有可用命令。

---

### 登录

在执行其他操作之前，您需要先登录：

```bash
./xixunyunsign.exe login -a <账号> -p <密码> -i <学校id>
```

#### 参数说明

- `-a` 或 `--account`：您的登录账号。
- `-p` 或 `--password`：您的登录密码。
- `-i` :学校id(使用search子命令查询)
#### 示例

```bash
./xixunyunsign.exe login -a user_number -p yourpassword
```

登录成功后，程序会在当前目录下生成一个 `config.db` 数据库文件，保存您的账号和 Token 信息。

---

### 查询签到信息

登录成功后，您可以查询当月的签到信息，并获取应签到的位置经纬度：

```bash
./xixunyunsign.exe query -a <账号>
```

#### 参数说明

- `-a` 或 `--account`：您的登录账号。

#### 示例

```bash
./xixunyunsign.exe query -a user_number
```

执行成功后，程序会输出查询结果，并将应签到的经纬度信息保存到数据库中。

---

### 执行签到

您可以使用以下命令进行签到：

```bash
./xixunyunsign.exe sign -a <账号> --address <地址>
```

#### 参数说明

- `-a` 或 `--account`：您的登录账号。
- `--address`：签到地点的详细地址。
- `--latitude`：签到地点的纬度（可选，程序会自动从数据库中获取）。
- `--longitude`：签到地点的经度（可选，程序会自动从数据库中获取）。
- `--remark`：备注（可选，默认值为 `0`）。
- `--comment`：评论（可选）。
- `-p`:省份（可选）默认 为空
- `-c`:城市（可选）默认 为空

#### 示例

```bash
./xixunyunsign.exe sign -a user_number --address "江苏省南京市玄武区北京东路41号"
```

**注意：**

- 如果您在命令行中未提供 `--latitude` 和 `--longitude` 参数，程序会自动从数据库中获取该账号的应签到经纬度信息。
- 如果数据库中不存在经纬度信息，程序会提示错误。此时，您可以手动提供经纬度，或者先执行 `query` 命令获取并保存经纬度信息。

#### 手动提供经纬度示例

```bash
./xixunyunsign.exe sign -a user_number \
  --address "江苏省南京市玄武区北京东路41号" \
  --latitude "32.069759" \
  --longitude "118.802972"
```

---

## 多用户支持

本工具支持多用户操作，每个用户的信息（账号、Token、应签到的经纬度）都会存储在 SQLite 数据库中。

### 数据库文件位置

数据库文件 `config.db` 存储在程序运行的当前目录下。

### 使用方法

#### 登录

```bash
./xixunyunsign.exe login -a <账号> -p <密码>
```

#### 查询签到信息

```bash
./xixunyunsign.exe query -a <账号>
```

#### 执行签到

```bash
./xixunyunsign.exe sign -a <账号> --address <地址>
```

#### 查询学校ID

```bash
./xixunyunsign.exe search -s <学校名称>
```

### 注意事项

- **首次签到前**，请确保已经执行过 `query` 命令，或者手动提供经纬度信息。
- **确保数据库文件一致性**：在同一个目录下执行所有命令，以确保程序能够正确访问数据库文件。
- **Token 有效期**：如果在执行 `query` 或 `sign` 命令时提示“会话已失效”，请重新执行 `login` 命令进行登录。

---

## 项目结构

```
.
├── main.go        // 程序入口
├── cmd            // 命令定义
│   ├── login.go   // 登录命令
│   ├── query.go   // 查询签到信息命令
│   ├── sign.go    // 执行签到命令
|   ├── search_school_id.go //查询签到的学校id
|   └── practice-report.go //ai自动编写周报月报,定时提交图片···
└── utils          // 工具函数
    ├── database.go // 数据库操作
    └── config.go   // 配置文件读写（已废弃）
```

## 依赖库

- **Cobra**：用于创建命令行界面
- **SQLite**：用于数据存储
- **`github.com/mattn/go-sqlite3`**：Go 的 SQLite3 驱动

请确保在项目的 `go.mod` 文件中添加以下依赖：

```go
require (
    github.com/mattn/go-sqlite3 v1.14.15 // 最新版本
    github.com/spf13/cobra v1.5.0        // 最新版本
)
```

---

## 常见问题解答

### Q1：执行签到时提示“未提供经纬度信息”怎么办？

**A**：请先执行 `query` 命令获取并保存应签到的经纬度信息，或者在执行 `sign` 命令时手动提供 `--latitude` 和 `--longitude` 参数。

### Q2：程序提示“会话已失效，请重新登录”怎么办？

**A**：请重新执行 `login` 命令进行登录，然后再次执行相应的操作。

### Q3：如何在不同目录下运行程序并保持数据一致？

**A**：确保将 `config.db` 数据库文件复制到程序运行的目录下，或者修改代码中的数据库路径，将数据库文件放在固定的绝对路径。

---

## 注意事项

- **确保数据库文件一致性**

    - 数据库文件 `config.db` 会在您运行程序的当前目录下生成。
    - 请在同一个目录下执行所有命令，以确保程序能够正确访问数据库文件。

- **Token 有效期**

    - 如果在执行 `query` 或 `sign` 命令时提示“会话已失效”，请重新执行 `login` 命令进行登录。

- **经纬度信息**

    - 如果未提供经纬度参数且数据库中不存在，应签到的经纬度信息，程序会提示错误。
    - 建议在首次签到前先执行 `query` 命令，获取并保存应签到的经纬度信息。

- **多用户支持**

    - 程序支持多用户操作，所有命令都需要使用 `-a` 或 `--account` 参数指定账号。

- **依赖环境**

    - 请确保已安装 Go 语言环境和必要的依赖库。
    - Windows 用户需要安装 `gcc` 或其他 C 编译器，用于编译 `github.com/mattn/go-sqlite3` 库。

---

## 贡献指南

欢迎对本项目提出意见或建议，您可以通过以下方式参与贡献：

- 提交 Issue
- 提交 Pull Request

在提交之前，请确保您的代码符合项目的编码规范，并经过充分测试。

---

## 许可证

本项目采用 GPL3.0 许可证进行许可。有关更多信息，请参阅 [LICENSE](LICENSE) 文件。

---

## 联系方式

如有特殊问题或定制需求，请我喝一杯咖啡吧：

- 开发者邮箱：kuilanmin@gmail.com

一般性程序功能问题直接提交[issue](https://github.com/theshdowaura/xixunyunsign/issues)即可，勿扰邮箱。

软件使用上的问题打开[discussions](https://github.com/theshdowaura/xixunyunsign/discussions)提问即可

---

## 致谢

感谢您的使用和支持！

---

## 示例操作流程

以下是一个完整的操作示例，展示如何使用该工具登录、查询签到信息和执行签到。

```bash
# 编译项目
go build -o xixunyunsign.exe

# 登录账户
./xixunyunsign.exe login -a user_number -p yourpassword

# 查询签到信息，获取并保存应签到的经纬度
./xixunyunsign.exe query -a user_number

# 执行签到，使用保存的经纬度信息
./xixunyunsign.exe sign -a user_number --address "江苏省南京市玄武区北京东路41号" --address_name "南京市人民政府研究室"
```

执行成功后，程序会输出签到结果。


---

## 更新日志

### v1.0.0

- 初始版本，实现了登录、查询签到信息和执行签到的基本功能。
- 支持多用户操作，使用 SQLite 数据库存储用户信息。
- 实现了 RSA 加密功能，对经纬度信息进行加密传输。

### v1.1.0

- 实现查询各个学校的ID功能，支持模糊匹配并显示所有匹配结果。
- 程序优化了数据库的操作流程，提升性能和稳定性。

### v1.2.0

- 增加通知渠道方糖(Server酱)，通过微信公众号接收签到信息


### v1.2.x.alpha

- 实现周报月报，ai自动代写，图片上传

### v1.3.0.alpha2
- 新增web api调用，实现网页端一键在线多签。

---

## 开发计划

- [x] 添加自动签到功能，支持定时任务。
- [x] 优化错误处理，提供更友好的提示信息。
- [x] 提供更详细的日志记录，方便调试和问题排查。
- [x] 增加对其他 API 接口的支持，如请假申请,自动定时提交周报月报，提交实习总结等(v1.3.0)。

---

## 参考资料

- [Go 语言官方文档](https://golang.org/doc/)
- [Cobra 命令行库](https://github.com/spf13/cobra)
- [SQLite3 Go 驱动](https://github.com/mattn/go-sqlite3)
- [RSA 加密解密](https://pkg.go.dev/crypto/rsa)

---

## 查看数据库

请使用navicat连接本项目下的config.db,使用sqlite的方式

## 为我买一杯coffee

[![Buy Me A Coffee](https://img.shields.io/badge/Buy%20Me%20A%20Coffee-%F0%9F%8D%8B-yellow.svg)](https://www.buymeacoffee.com/theshdowaura)
