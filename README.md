# Gosh (Go Shell)

Gosh 是一个使用 Go 语言编写的轻量级、模块化的类 Unix Shell 实现。它采用了现代化的架构设计，将解析（Parsing）、执行（Execution）和内建命令（Builtins）完全解耦，是一个学习 Go 语言系统编程、进程管理和 I/O 管道机制的绝佳示例。

## 功能特性 (Features)

### 核心功能

* **多级管道支持 (`|`)**：支持无限级管道操作，父子进程资源自动管理。
* 例如：`ls -l | grep ".go" | wc -l`


* **I/O 重定向**：
* 输入重定向：`wc -l < input.txt`
* 输出覆盖：`echo "hello" > output.txt`
* 输出追加：`echo "world" >> output.txt`
* 标准错误重定向：`ls not_exist 2> error.log` 或 `2>>` 追加模式。



### 内建命令 (Built-ins)

* `cd`：切换工作目录（支持 `~` 展开为用户主目录）。
* `pwd`：显示当前工作路径。
* `echo`：回显文本。
* `history`：命令历史管理。
* `history`：显示所有历史。
* `history 10`：显示最近 10 条。
* `history -w [file]`：将历史写入文件。
* `history -r [file]`：从文件加载历史。


* `type`：检测命令类型（显示是内建命令还是外部可执行路径）。
* `exit`：保存历史并安全退出 Shell。

### 交互体验

* **行编辑**：集成了 `readline` 库，支持光标移动、删除等行编辑操作。
* **自动补全**：按下 `Tab` 键可自动补全内建命令、PATH 中的外部命令以及文件路径。
* **持久化历史**：通过环境变量 `HISTFILE` 配置，或默认自动保存历史记录，重启 Shell 后依然可用。

##  架构设计 (Architecture)

项目遵循 Clean Architecture 原则，目录结构清晰，模块职责单一：

```text
├── cmd/
│   └── myshell/
│       ├── main.go      # 程序入口：负责 Session 初始化和 REPL 主循环
│       └── readline.go  # Readline 配置与自动补全逻辑
├── internal/
│   ├── builtins/        # 内建命令层：包含 cd, echo, exit 等具体实现
│   ├── executor/        # 执行引擎层：负责管道连接、文件句柄管理、进程启动
│   ├── parser/          # 解析层：词法分析与语法解析，生成 Pipeline 结构体
│   └── session/         # 状态层：管理全局状态（历史记录、环境变量、IO流）

```

## 具体实现(Concrete Accomplishment)


##  安装与运行 (Installation & Usage)

### 前置条件

* Go 1.18 或更高版本
* Linux / macOS 环境（Windows 下部分路径处理可能略有不同）

### 快速运行

你可以直接使用 Go 命令运行源码：

```bash
# 1. 下载依赖
go mod tidy

# 2. 运行 Shell
go run cmd/myshell/main.go

```

### 编译构建

生成可执行文件以便日常使用：

```bash
# 编译
go build -o gosh cmd/myshell/main.go

# 启动
./gosh

```

## 使用示例 (Examples)

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
# 查看最近 5 条历史
$ history 5

# 手动保存历史到指定文件
$ history -w my_history_backup.txt

```

##  贡献 (Contributing)

欢迎提交 Pull Requests 来改进代码或增加新功能。目前的待办事项（TODO）包括：

* 支持环境变量设置（`export`）。
* 支持逻辑运算符（`&&`, `||`）。
* 支持后台运行作业（`&`）。

##  License

MIT License