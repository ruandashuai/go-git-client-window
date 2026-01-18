# 🚀 Go Git Client

**一个现代化的跨平台 Git 客户端桌面应用程序，基于 Wails 框架构建**

Go Git Client 是一款功能强大的 Git 客户端，结合了桌面应用的便捷性和 Web 访问的灵活性。它提供了一个直观的图形界面来执行常见的 Git 操作，让开发者能够更高效地管理他们的代码仓库。

## 💻 技术栈

- **后端**: [Go](https://golang.org/) 1.25
- **前端框架**: [Wails](https://wails.io/) v2.11.0
- **前端技术**: [Vue.js](https://vuejs.org/) 3.x, [Vite](https://vitejs.dev/)
- **HTTP 服务器**: Go 标准库 `net/http`
- **Git 操作**: 通过 `os/exec` 调用系统 Git 命令
- **UI 设计**: Vue 组件化架构，现代化响应式界面

## ✨ 核心功能

- 🖥️ **双模式运行**: 支持桌面应用和 Web 访问两种模式
- 🛠️ **完整 Git 操作**: 初始化、克隆、状态查看、添加文件、提交等
- 🔗 **浏览器集成**: 支持在浏览器中打开任意 URL
- 🌐 **RESTful API**: 提供 HTTP API 接口，支持远程操作
- 🎨 **现代化 UI**: 基于 Vue 的响应式用户界面
- 📊 **实时状态显示**: 显示 Git 仓库的当前状态和历史记录
- 🔄 **分支管理**: 查看和切换 Git 分支

## 🚀 快速开始

### 环境要求 (Prerequisites)

- [Go](https://golang.org/dl/) 1.25 或更高版本
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) v2.11.0
- [Git](https://git-scm.com/downloads) (系统需安装 Git 命令行工具)
- Node.js 和 npm (用于前端开发)
- Windows/macOS/Linux 操作系统

### 安装步骤 (Installation)

1. **安装 Wails CLI** (如果尚未安装):
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

2. **克隆项目**:
   ```bash
   git clone <repository-url>
   cd go-git-client-window
   ```

3. **安装前端依赖**:
   ```bash
   cd frontend
   npm install
   cd ..
   ```

### 运行命令 (Usage)

1. **开发模式**:
   ```bash
   # 在项目根目录运行
   wails dev
   ```
   这将启动开发服务器，自动编译并运行应用，支持热重载。

2. **构建生产版本**:
   ```bash
   # 构建默认平台版本
   wails build
   
   # 构建特定平台版本 (例如 Windows)
   wails build -platform windows/amd64
   ```
   
   构建完成后，可执行文件位于 `build/bin/` 目录。

3. **运行应用**:
   - **桌面模式**: 直接运行编译后的可执行文件
   - **Web 模式**: 启动应用后，在浏览器中访问提供的 Web 界面

## 📁 目录结构

```
go-git-client-window/
├── core/                    # 核心业务逻辑
│   └── git_service.go       # Git 服务层
├── git/                     # Git 操作封装
│   ├── operations.go        # Git 操作实现
│   └── operations_test.go   # Git 操作测试
├── internal/gitcmd/         # Git 命令执行层
│   ├── git_exec.go          # Git 命令执行
│   └── git_operations.go    # Git 命令操作
├── models/                  # 数据模型定义
│   └── git_models.go        # Git 相关数据模型
├── frontend/                # 前端源码目录
│   ├── src/
│   │   ├── components/      # Vue 组件
│   │   │   ├── BranchList.vue      # 分支列表组件
│   │   │   ├── CommitHistory.vue   # 提交历史组件
│   │   │   ├── CommitPanel.vue     # 提交面板组件
│   │   │   ├── FileStatusPanel.vue # 文件状态面板组件
│   │   │   ├── StatusBar.vue       # 状态栏组件
│   │   │   └── TopNavBar.vue       # 顶部导航栏组件
│   │   ├── styles/          # 样式文件
│   │   ├── App.vue          # 根组件
│   │   └── main.js          # 前端入口文件
│   ├── package.json         # 前端包配置
│   └── vite.config.js       # Vite 配置文件
├── main.go                  # 主程序入口
├── go.mod                   # Go 模块依赖配置
├── go.sum                   # Go 依赖锁定文件
├── wails.json               # Wails 项目配置文件
└── IFLOW.md                 # 项目详细说明文档
```

## 🔧 功能特性

- **现代化 UI**: 使用 Vue.js 构建的响应式用户界面，提供流畅的用户体验
- **实时状态监控**: 实时显示 Git 仓库的状态变化
- **分支管理**: 清晰展示所有本地和远程分支
- **提交历史**: 以可视化方式展示提交历史记录
- **文件状态**: 直观显示工作目录中的文件变更状态
- **API 接口**: 提供 RESTful API，支持外部系统集成

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来帮助改进这个项目！

## 📄 许可证

Copyright (C) 2026 Git Client