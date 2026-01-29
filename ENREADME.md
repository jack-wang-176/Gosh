
# 🐚 Gosh (Go Shell)

> **A high-performance interactive Shell written from scratch in Go.**
> A system-level interactive command-line interpreter implemented from scratch, deeply simulating core Bash mechanisms.

![Go Version](https://img.shields.io/badge/go-1.18+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

<p align="center">
  <strong>English Version</strong> | <a href="./README.md">中文版本</a>
</p>
**Gosh** is a lightweight, modular Unix-like Shell implementation. It adopts a modern **decoupled architecture design**, completely separating Parsing, Execution, and Builtins.

This project deeply practices the Unix philosophy, implementing efficient **multi-stage pipelines** and **file descriptor redirection** through **Go concurrency**. It is an excellent example for learning operating system process management, IPC communication, and lexical analysis.

---

## ✨ Features

### 🚀 Core Mechanisms
* **Multi-stage Pipelines (`|`)**: Supports infinite pipeline operations, automatically managing I/O resources and lifecycles of parent-child processes.
    * *Example:* `ls -l | grep ".go" | wc -l`
* **Advanced I/O Redirection**:
    * Input Redirection: `wc -l < input.txt`
    * Output Overwrite: `echo "hello" > output.txt`
    * Output Append: `echo "world" >> output.txt`
    * **Error Stream Redirection**: `ls not_exist 2> error.log` (Supports `2>>` append).

### 🛠️ Built-ins
* `cd`: Change directory (supports intelligent expansion of `~` to the user's home directory).
* `pwd`: Display current working directory.
* `echo`: Echo text, supports parameter concatenation.
* `history`: Powerful history management (supports `-w` write / `-r` read).
* `type`: Command type detection (distinguishes between Built-in and External).
* `exit`: Safely exit and persist history.

### 💻 Interaction Experience
* **Line Editing Enhancement**: Integrated `readline` library, supporting cursor movement, deletion, and history backtracking.
* **Smart Completion**: Pressing the `Tab` key automatically completes commands or file paths.
* **State Persistence**: Saves history via the `HISTFILE` environment variable or default strategy, preventing loss upon restart.

---

## 🧠 Implementation Details

This project performs deep module partitioning and logical encapsulation in its underlying implementation:

### 1. State Management (Session)
* **Global Control**: Maintains Shell global variables and standardizes I/O streams (Stdin/Stdout/Stderr), primarily for use by built-in commands in the `builtins` directory.
* **Convenient Initialization**: Provides a `NewSession` factory method for one-click environment initialization.

### 2. Parsing Layer (Parse)
* **Pipeline Splitting**: The `Parse` function first physically splits the input string using the `|` symbol as a delimiter, then encapsulates it into the **Pipeline** struct for storage.
* **Tokenization (Lexer)**: `parseToken` handles finer-grained parsing, processing escape logic for single quotes (`'`) and double quotes (`"`), further disassembling the string separated by `|` into concrete Command units. At the low level, a Command is converted into a `[]string` slice.
* **Redirection Parsing**: `parseSingleCommand` performs deep parsing on concrete commands, capable of identifying and processing **5 types of I/O stream redirection** (`<`, `>`, `>>`, `2>`, `2>>`).
    * **Key Logic**: At this stage, the parser **removes** I/O redirection identifiers and filenames, purifying the character fragments maintained within the Command into purely "Instruction Name + Arguments", while encapsulating the redirected filenames and modes separately in the Command struct.

### 3. Execution Engine (Executor)
* **Channel Construction (Execute)**: The main function of this is to build I/O channels for commands.
    * It iterates through `command.args`, first handling pipeline redirection logic.
    * It maintains temporary I/O objects internally, constructing file redirections and opening corresponding files based on flags set by the Parse layer.
    * If the current command is not the last in the Pipeline, it builds an `os.Pipe` channel for connection.
* **Command Scheduling (startSingleCommand)**:
    * Used to specifically invoke or execute command operations.
    * **Concurrency Safety**: For built-in commands, an extra **Goroutine** is constructed for execution to ensure the main function continues running and maintains memory safety, passing errors (`err`) through channels.
    * **Unified Interface**: Adopts the **`WaitFunc` interface** form to standardize the exit waiting logic of the main function (whether waiting for internal goroutines or external processes).

### 4. Built-in Commands (Builtins)
* **Function Mapping**: Designs internal command operations, mapping command strings to processing functions via a `CodFunc` map, ensuring concrete function interfaces are not exposed outside the package.
* **Unified Interface**: Uses `CommandFunc` as the unified implementation interface for all built-in command functions, guaranteeing feasibility of direct function invocation via mapping loops.
* **Extensibility**: To extend command-line opcodes, simply add corresponding logic in the concrete implementation functions.

### 5. Entry & Interaction (Main & Readline)
* **Myshell**: The entry point of the entire program, initializing `Session` and `Readline`.
* **Real-time Interaction**: Implements real-time input reading, capturing user input and entering the REPL loop.

---

## 🏗️ Directory Structure

```text
├── cmd/
│   └── myshell/
│       ├── main.go      # Program entry: Responsible for Session initialization and REPL main loop
│       └── readline.go  # Readline configuration and auto-completion logic
├── internal/
│   ├── builtins/        # Builtin command layer: CodFunc mapping and concrete implementation
│   ├── executor/        # Execution engine layer: Execute channel construction and startSingleCommand scheduling
│   ├── parser/          # Parsing layer: Pipeline splitting and redirection symbol stripping
│   └── session/         # State layer: Global Session management



---

## 📦 Installation & Usage

### Prerequisites

* Go 1.18 or higher
* Linux / macOS environment (Path handling on Windows may differ slightly)

### Quick Run

```bash
# 1. Download dependencies
go mod tidy

# 2. Run Shell
go run cmd/myshell/main.go

```

### Build

```bash
# Compile
go build -o gosh cmd/myshell/main.go

# Start
./gosh

```

---

## 📝 Examples

**1. Basic Commands and Pipelines**

```bash
$ ls -l | grep "main.go"

```

**2. File Writing and Reading**

```bash
# Overwrite
$ echo "Hello Gosh" > hello.txt
# Append
$ echo "Another line" >> hello.txt
# Input Redirection
$ cat < hello.txt

```

**3. Error Log Processing**

```bash
# Redirect standard error output to file
$ ls /file_not_exist 2> error.log

```

**4. History Operations**

```bash
# Manually save history to specified file
$ history -w my_history_backup.txt

```

---

## 🤝 Contributing

Pull Requests are welcome. To-Do List (TODO):

* Support environment variable setting (`export`).
* Support logical operators (`&&`, `||`).
* Support background jobs (`&`).

## 📄 License

MIT License

```

```