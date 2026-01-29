# 🐚 Gosh (Go Shell)

> **A high-performance interactive Shell written from scratch in Go.**
> 一个从零实现的系统级交互式命令行解释器，深度模拟了 Bash 核心机制。


![Go Version](https://img.shields.io/badge/go-1.18+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

<p align="center">
  <a href="./ENREADME.md">English Version</a> | <strong>中文版本</strong>
</p>

**Gosh** 是一个轻量级、模块化的类 Unix Shell 实现。它采用了现代化的**解耦架构设计**，将解析（Parsing）、执行（Execution）和内建命令（Builtins）完全分离。

本项目深度实践了 Unix 哲学，通过 **Go 并发编程** 实现了高效的**多级管道**与**文件描述符重定向**，是学习操作系统进程管理、IPC 通信及词法分析的绝佳示例。

---

## ✨ 功能特性 (Features)

### 🚀 核心机制
* **多级管道 (`|`)**：支持无限级管道操作，自动管理父子进程的 I/O 资源与生命周期。
    * *Example:* `ls -l | grep ".go" | wc -l`
* **高级 I/O 重定向**：
    * 输入重定向：`wc -l < input.txt`
    * 输出覆盖：`echo "hello" > output.txt`
    * 输出追加：`echo "world" >> output.txt`
    * **错误流重定向**：`ls not_exist 2> error.log` (支持 `2>>` 追加)。

### 🛠️ 内建命令 (Built-ins)
* `cd`：切换目录（支持 `~` 智能展开为用户主目录）。
* `pwd`：显示当前工作路径。
* `echo`：回显文本，支持参数拼接。
* `history`：强大的历史记录管理（支持 `-w` 写入 / `-r` 读取）。
* `type`：命令类型检测（区分 Built-in 与 External）。
* `exit`：安全退出并持久化历史。

### 💻 交互体验
* **行编辑增强**：集成 `readline` 库，支持光标移动、删除、历史回溯。
* **智能补全**：按下 `Tab` 键可自动补全命令或文件路径。
* **状态持久化**：通过 `HISTFILE` 环境变量或默认策略保存历史，重启不丢失。

---

## 🧠 具体实现 (Implementation Details)

本项目在底层实现上进行了深度的模块划分与逻辑封装：

### 1. 状态管理 (Session)
* **全局控制**：维护 Shell 的全局变量并规范 I/O 流（Stdin/Stdout/Stderr），主要供 `builtins` 目录下的内建命令使用。
* **便捷初始化**：提供 `NewSession` 工厂方法，实现环境的一键初始化。

### 2. 解析层 (Parse)
* **管道切分**：`Parse` 函数首先将输入字符串以 `|` 符号为界进行物理切分，随后将其封装进 **Pipeline** 结构体中储存。
* **Token 化 (Lexer)**：`parseToken` 负责更细粒度的解析，它处理单引号 (`'`) 和双引号 (`"`) 的转义逻辑，将以 `|` 分隔的字符串进一步拆解为具体的 Command 单元。在底层，一个 Command 被转换为 `[]string` 切片。
* **重定向解析**：`parseSingleCommand` 对具体命令进行深度解析，能够识别并处理 **5种 I/O 流重定向**（`<`, `>`, `>>`, `2>`, `2>>`）。
    * **关键逻辑**：在此阶段，解析器会**去除** I/O 重定向的标识符和文件名，使得 Command 内部维护的字符片段被净化为纯粹的“指令名 + 参数”，而将重定向的文件名和模式单独封装在 Command 结构体中。

### 3. 执行引擎 (Executor)
* **通道构建 (Execute)**：此函数的主要作用是构建命令的 I/O 通道。
    * 它遍历 `command.args`，首先处理管道的重定向逻辑。
    * 在函数内部维护临时 I/O 对象，根据 Parse 层打下的标识位构建文件重定向并打开相应文件。
    * 如果当前命令不是 Pipeline 的最后一个，则构建 `os.Pipe` 通道进行连接。
* **命令调度 (startSingleCommand)**：
    * 用来具体调用或执行命令操作。
    * **并发安全**：对于内建命令，额外构建 **协程 (Goroutine)** 运行，以保证主函数的持续运行和内存安全，并通过管道传递错误 (`err`)。
    * **接口统一**：采用 **`WaitFunc` 接口** 形式来统一规范主函数的退出等待逻辑（无论是等待内建协程还是外部进程）。

### 4. 内建命令 (Builtins)
* **函数映射**：设计内部命令操作，通过 `CodFunc` map 进行命令字符串和处理函数的映射调用，保证了具体函数接口不暴露在包外部。
* **统一接口**：使用 `CommandFunc` 作为所有内建命令函数的统一实现接口，保证了可以通过映射循环直接调用函数的可行性。
* **扩展性**：对于命令行操作码的扩充，只需在具体的实现函数中添加对应逻辑即可。

### 5. 入口与交互 (Main & Readline)
* **Myshell**：整个程序的入口，初始化 `Session` 和 `Readline`。
* **实时交互**：实现了实时读取输入符，捕获用户输入并进入 REPL 循环。

---

## 🏗️ 目录结构 (Directory Structure)

```text
├── cmd/
│   └── myshell/
│       ├── main.go      # 程序入口：负责 Session 初始化和 REPL 主循环
│       └── readline.go  # Readline 配置与自动补全逻辑
├── internal/
│   ├── builtins/        # 内建命令层：CodFunc 映射与具体实现
│   ├── executor/        # 执行引擎层：Execute 通道构建与 startSingleCommand 调度
│   ├── parser/          # 解析层：Pipeline 切分与重定向符号剥离
│   └── session/         # 状态层：Global Session 管理
```
---
## 📦 安装与运行 (Installation & Usage)

### 前置条件

* Go 1.18 或更高版本
* Linux / macOS 环境（Windows 下路径处理可能略有差异）

### 快速运行

```bash
# 1. 下载依赖
go mod tidy

# 2. 运行 Shell
go run cmd/myshell/main.go

```

### 编译构建

```bash
# 编译
go build -o gosh cmd/myshell/main.go

# 启动
./gosh

```

---

## 📝 使用示例 (Examples)

**1. 基础命令与管道**

```bash
$ ls -l | grep "main.go"

```

**2. 文件写入与读取**

```bash
# 覆盖写入
$ echo "Hello Gosh" > hello.txt
# 追加写入
$ echo "Another line" >> hello.txt
# 输入重定向
$ cat < hello.txt

```

**3. 错误日志处理**

```bash
# 将标准错误输出重定向到文件
$ ls /file_not_exist 2> error.log

```

**4. 历史记录操作**

```bash
# 手动保存历史到指定文件
$ history -w my_history_backup.txt

```

---

## 🤝 贡献 (Contributing)

欢迎提交 Pull Requests。待办事项（TODO）：

* 支持环境变量设置（`export`）。
* 支持逻辑运算符（`&&`, `||`）。
* 支持后台运行作业（`&`）。

## 📄 License

MIT License

```

### 修改说明：
1.  **恢复了 Parse 细节**：明确写回了“去除 I/O 重定向标识符”、“Command 内部维护纯粹指令名”的逻辑。
2.  **恢复了 Executor 细节**：明确写回了 `WaitFunc` 接口、内建命令使用协程（Goroutine）运行、以及 `Execute` 中维护临时 I/O 对象的逻辑。
3.  **恢复了 Builtins 细节**：明确写回了 `CodFunc` 和 `CommandFunc` 的映射与接口设计。
4.  **保留了格式优化**：虽然内容恢复了，但我还是帮您把排版（缩进、代码块、小标题）弄得更整齐了，方便别人阅读。

```