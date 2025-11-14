# 📝 TaskTrackerCLI: A Simple Command-Line Task Manager by https://roadmap.sh/projects/task-tracker

**TaskTrackerCLI** is a minimalist and fast command-line tool for
managing your tasks directly from the terminal.\
The application allows you to easily add, view, update, and delete tasks
using a simple JSON file for local data storage.

------------------------------------------------------------------------

## ✨ Features

-   **CRUD operations:** Full support for creating, reading, updating,
    and deleting tasks.\
-   **Status Management:** Quickly change task status (`TODO`,
    `IN_PROGRESS`, `DONE`).\
-   **Automatic ID Assignment:** Tasks are automatically assigned unique
    identifiers.\
-   **Cleanup:** Bulk deletion of completed (`DONE`) tasks.\
-   **Local Storage:** All data is stored in a single local JSON file.

------------------------------------------------------------------------

## 📦 Installation

### 1. Requirements

You must have **Go 1.18+** installed.

### 2. Install via `go install`

You can install the application directly using:

``` bash
go install github.com/TaskTrackerCLI@latest
```

> **Note:** Make sure your `$GOPATH/bin` or `$HOME/go/bin` directory is
> included in your system `$PATH`.\
> After installation, the application will be available as the `task`
> command.

------------------------------------------------------------------------

## 🚀 Usage

All commands start with:

``` bash
task
```

### 1. Add a Task (`task add`)

Create a new task with a description:

``` bash
task add "Learn new Go testing patterns" "Read Clean Code and apply TDT"
```

### 2. List Tasks (`task list`)

Display all existing tasks:

``` bash
task list
```

### 3. Update Task Status (`task mark`)

Change the status of a task (use the task ID, e.g., `1`):

``` bash
# Mark task 1 as DONE
task mark 1 DONE

# Mark task 2 as IN_PROGRESS
task mark 2 IN_PROGRESS
```

Available statuses: `TODO`, `IN_PROGRESS`, `DONE`.

### 4. Delete a Task (`task delete`)

Delete a task by its ID:

``` bash
task delete 1
```

### 5. Clean Completed Tasks (`task clean`)

Remove all tasks with the status `DONE`:

``` bash
task clean
```

(You will be asked for confirmation, as this action is irreversible.)

