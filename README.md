# Gosh (Go Shell)

Gosh 是一个使用 Go 语言编写的轻量级、模块化的类 Unix Shell 实现。它采用了现代化的架构设计，将解析（Parsing）、执行（Execution）和内建命令（Builtins）完全解耦，是一个学习 Go 语言系统编程和 Shell 工作原理的绝佳示例。

## 🚀 功能特性 (Features)

* **管道支持 (`|`)**：支持多级管道操作，如 `ls | grep go | wc -l`。
* **I/O 重定向**：
    * 输入重定向 (`<`)
    * 输出重定向 (`>`, `>>`)
    * 标准错误重定向 (`2>`, `2>>`)
* **内建命令 (Built-ins)**：
    * `cd`：切换目录（支持 `~` 展开）。
    * `pwd`：显示当前路径。
    * `echo`：回显文本。
    * `history`：查看、读取 (`-r`)、写入 (`-w`) 历史记录。
    * `type`：判断命令是内建还是外部程序。
    * `exit`：退出 Shell。
* **交互体验**：
    * 基于 `readline` 的行编辑功能。
    * 支持命令自动补全（Tab 键）。
    * 支持历史记录持久化（保存到文件）。
* **健壮的架构**：
    * **Parser**：独立的词法分析与语法解析器。
    * **Executor**：资源安全的执行引擎，自动管理文件句柄与管道关闭。
    * **Session**：统一管理会话状态（环境变量、历史记录等）。

## 🛠️ 架构设计 (Architecture)

项目采用分层架构，目录结构清晰：

```text
├── cmd/
│   └── myshell/       # 程序入口 (Main)
├── internal/
│   ├── builtins/      # 内建命令的具体实现 (cd, echo, etc.)
│   ├── executor/      # 执行引擎 (负责管道连接、进程启动)
│   ├── parser/        # 解析器 (负责将字符串解析为 Pipeline 结构体)
│   └── session/       # 会话状态管理

📦 安装与运行 (Installation & Usage)
前置条件
Go 1.18+

运行
可以直接使用 Go 命令运行项目：

Bash
# 下载依赖
go mod tidy

# 运行
go run .
编译
Bash
go build -o gosh main.go
./gosh
📝 使用示例 (Examples)
1. 基础命令与管道

Bash
$ ls -l | grep ".go"
2. 文件重定向

Bash
# 将 hello 写入文件
$ echo "hello world" > test.txt

# 追加内容
$ echo "another line" >> test.txt

# 读取文件内容并统计行数
$ wc -l < test.txt
3. 错误处理

Bash
# 将错误信息重定向到 error.log (假设 file_not_exist 不存在)
$ ls file_not_exist 2> error.log
4. 历史记录

Bash
$ history 5
$ history -w my_history.txt