# Arch Linux 专用 Discord Canary 更新器

一个用于在 Arch Linux 系统上自动保持 Discord Canary 最新版本的简单工具。

[English](README.md) | [Русский](README_ru.md)

## 问题

使用官方 Discord Canary 版本（tar.gz）的 Arch Linux 用户面临一个不便的问题：应用程序每次启动时都可能需要下载更新。每天手动下载压缩包、解压并替换文件既繁琐又不便。

## 解决方案

这个用 Go 语言编写的简单脚本实现了自动化处理。它执行以下操作：

- 从 Discord 官方服务器下载最新的 discord-canary.tar.gz 压缩包。
- 解压压缩包。
- 用新版本替换旧版本。

## 安装

1. 从 [发布页面](https://github.com/Lazlo105/Discord-canary-updater-for-Arch-linux/releases/latest) 下载最新版本。
2. 授予执行权限：
   ```bash
   chmod +x discord-updater
   ```
3. 运行更新器：
   ```bash
   ./discord-updater
   ```

## 系统要求

- Go 1.20 或更高版本（用于从源代码构建）
- Make（可选，用于使用 Makefile 目标）

## 从源代码构建

1. 克隆仓库：
   ```bash
   git clone https://github.com/Lazlo105/Discord-canary-updater-for-Arch-linux.git
   cd Discord-canary-updater-for-Arch-linux
   ```

2. 构建项目：
   ```bash
   make build
   ```

3. 运行应用程序：
   ```bash
   ./bin/updater
   ```

## 开发指南

- 格式化代码：`make fmt`
- 运行测试：`make test`
- 运行代码检查：`make lint`
- 查看测试覆盖率：`make test-coverage`

## 许可证

本项目采用 MIT 许可证。详情请参见 [LICENSE](LICENSE) 文件。

## 贡献

欢迎贡献！请在提交拉取请求之前阅读我们的 [贡献指南](CONTRIBUTORS.md)。

## 安全

如果您发现安全漏洞，请参阅我们的 [安全政策](SECURITY.md)。
