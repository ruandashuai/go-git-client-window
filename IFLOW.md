# Go Git Client - 项目说明

## 项目概述

Go Git Client 是一个基于 Wails 框架构建的现代化 Git 客户端桌面应用程序。该项目采用 Go 作为后端，HTML/CSS/JavaScript 作为前端，支持桌面应用和 Web 访问两种模式，类似于 VSCode Web 版的体验。

### 主要特性

- **双模式运行**：既可作为桌面应用使用，也可通过 HTTP 服务器在浏览器中访问
- **Git 操作支持**：提供完整的 Git 基础操作，包括初始化、克隆、状态查看、添加文件和提交
- **浏览器集成**：支持在 Chrome 浏览器中打开任意 URL
- **RESTful API**：提供 HTTP API 接口，支持通过 API 调用执行 Git 操作
- **现代化 UI**：采用渐变色设计和卡片式布局，提供良好的用户体验

### 技术栈

- **后端**：Go 1.25
- **前端框架**：Wails v2.11.0
- **HTTP 服务器**：Go 标准库 `net/http`
- **Git 操作**：通过 `os/exec` 调用系统 Git 命令
- **UI 设计**：原生 HTML/CSS/JavaScript，无额外前端框架依赖

## 项目结构

```
go-git-client-window/
├── main.go                 # 主程序文件，包含应用逻辑和 HTTP 服务器
├── go.mod                  # Go 模块依赖配置
├── go.sum                  # Go 依赖锁定文件
├── wails.json              # Wails 项目配置文件
├── build/                  # 构建输出目录
│   ├── appicon.png        # 应用图标
│   ├── bin/               # 可执行文件
│   └── windows/           # Windows 特定资源
└── frontend/              # 前端资源目录
    ├── package.json       # 前端包配置（当前无构建步骤）
    ├── dist/              # 前端构建输出
    └── src/
        ├── index.html     # 主页面
        └── wailsjs/       # Wails 自动生成的绑定代码
```

## 构建和运行

### 环境要求

- Go 1.25 或更高版本
- Wails v2.11.0
- Git（系统需安装 Git 命令行工具）
- Windows/macOS/Linux 操作系统

### 安装依赖

```bash
# 安装 Wails CLI（如果尚未安装）
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 安装前端依赖（如有需要）
cd frontend
npm install
```

### 开发模式

```bash
# 在项目根目录运行
wails dev
```

这将启动开发服务器，自动编译并运行应用，支持热重载。

### 构建生产版本

```bash
# 构建 Windows 版本
wails build

# 构建特定平台版本
wails build -platform windows/amd64
```

构建完成后，可执行文件位于 `build/bin/` 目录。

### 运行应用

1. **桌面模式**：直接运行编译后的可执行文件
2. **Web 模式**：
   - 启动应用后，点击界面上的"启动服务器"按钮
   - 默认端口为 8080
   - 在浏览器中访问 `http://localhost:8080`

## 开发约定

### 代码风格

- **Go 代码**：遵循 Go 官方代码规范，使用 `gofmt` 格式化
- **函数命名**：使用 PascalCase 导出函数，camelCase 私有函数
- **注释**：导出函数必须包含注释说明，使用中文注释
- **错误处理**：所有错误必须处理，不忽略任何错误

### Git 操作方法

所有 Git 操作方法都遵循以下模式：

```go
func (a *App) GitOperationName(params ...string) (string, error) {
    // 参数验证
    if param == "" {
        return "", fmt.Errorf("参数不能为空")
    }
    
    // 执行 Git 命令
    cmd := exec.Command("git", "command", args...)
    cmd.Dir = path
    
    // 获取输出
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("操作失败: %s", string(output))
    }
    
    return string(output), nil
}
```

### HTTP API 规范

所有 API 端点遵循以下规范：

- **请求方法**：使用 POST 方法进行数据操作
- **请求格式**：JSON 格式
- **响应格式**：统一的 JSON 响应结构

```json
{
  "success": true,
  "message": "操作成功",
  "data": "返回数据"
}
```

### API 端点列表

| 端点 | 方法 | 描述 |
|------|------|------|
| `/api/server-info` | GET | 获取服务器信息 |
| `/api/init` | POST | 初始化 Git 仓库 |
| `/api/clone` | POST | 克隆 Git 仓库 |
| `/api/status` | POST | 获取 Git 状态 |
| `/api/add` | POST | 添加文件到暂存区 |
| `/api/commit` | POST | 提交更改 |

### 前端开发

- **HTML 结构**：使用语义化标签，保持清晰的 DOM 结构
- **CSS 样式**：使用内联样式（当前实现），保持样式与组件在一起
- **JavaScript**：使用 async/await 处理异步操作，通过 Wails 绑定调用 Go 方法
- **用户反馈**：所有操作必须有视觉反馈（成功/错误提示）

### 测试

当前项目未包含自动化测试。建议在添加新功能时：

1. 测试桌面应用模式下的功能
2. 测试 Web 访问模式下的功能
3. 测试 HTTP API 的正确性
4. 测试跨平台兼容性（Windows/macOS/Linux）

## 关键文件说明

### main.go

主程序文件，包含：

- **App 结构体**：应用的核心结构，管理上下文和 HTTP 服务器
- **HTTP 服务器**：提供静态文件服务和 RESTful API
- **Git 操作方法**：封装所有 Git 命令行操作
- **浏览器集成**：支持在 Chrome 中打开 URL
- **生命周期钩子**：startup、domReady、shutdown

### wails.json

Wails 项目配置文件，定义：

- 应用名称和版本
- 前端资源目录
- 构建和开发命令
- 应用元数据

### frontend/src/index.html

主页面文件，包含：

- UI 布局和样式
- HTTP 服务器控制界面
- Git 操作界面
- JavaScript 交互逻辑

## 使用示例

### 通过 API 调用 Git 操作

```bash
# 初始化仓库
curl -X POST http://localhost:8080/api/init \
  -H "Content-Type: application/json" \
  -d '{"path":"D:\\workspace\\my-project"}'

# 查看状态
curl -X POST http://localhost:8080/api/status \
  -H "Content-Type: application/json" \
  -d '{"path":"D:\\workspace\\my-project"}'

# 克隆仓库
curl -X POST http://localhost:8080/api/clone \
  -H "Content-Type: application/json" \
  -d '{"repoUrl":"https://github.com/user/repo.git","targetPath":"D:\\workspace\\repo"}'
```

### 在浏览器中打开

```bash
# 使用 Go 方法
window.go.main.App.OpenInBrowser("https://github.com")

# 或直接在浏览器地址栏输入
http://localhost:8080
```

## 注意事项

1. **Git 依赖**：系统必须安装 Git，否则所有 Git 操作将失败
2. **端口占用**：HTTP 服务器默认使用 8080 端口，如被占用需修改端口
3. **路径格式**：Windows 路径需要使用双反斜杠或正斜杠
4. **浏览器兼容**：Chrome 浏览器集成功能仅支持 Windows、macOS 和 Linux
5. **权限要求**：某些 Git 操作可能需要管理员权限

## 扩展建议

未来可以考虑添加以下功能：

- 更多 Git 操作（分支管理、合并、拉取、推送等）
- 文件浏览器和编辑器集成
- Git 历史记录可视化
- 冲突解决工具
- 多仓库管理
- 用户配置和偏好设置
- 主题切换（深色/浅色模式）
- 键盘快捷键支持
- 国际化支持

## 许可证

Copyright (C) 2026 Git Client