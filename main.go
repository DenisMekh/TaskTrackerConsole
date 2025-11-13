package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/TaskTrackerCLI/structures"
	"github.com/TaskTrackerCLI/task_manager"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var tm *task_manager.TaskManager

var mainCmd = &cobra.Command{
	Use:   "TaskTracker",
	Short: "TaskTracker for track your tasks",
	Long:  "A little bit long description for TaskTracker",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to the TaskTracker CLI! Use --help for usage ")
	}}

var addCmd = &cobra.Command{
	Use:   "add [task_name] [task_description]",
	Short: "add a new task",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		taskName := args[0]
		taskDescription := args[1]
		fmt.Printf("Adding task: Name='%s', Description='%s'\n", taskName, taskDescription)
		value, err := tm.AddTask(taskName, taskDescription)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error adding task: %v\n", err)
			return
		}
		fmt.Printf("✅ Task added successfully! ID: %d\n", value)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update [task_id] [task_name] [task_description]",
	Short: "update a task",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Task ID must be an integer. %v\n", err)

			return
		}
		taskName := args[1]
		taskDescription := args[2]
		arguments := make(map[string]string)
		arguments["task_name"] = taskName // 💡
		arguments["task_description"] = taskDescription

		updated, err := tm.UpdateTask(id, arguments)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating task: %v\n", err)
			return
		}

		if !updated {
			fmt.Fprintf(os.Stderr, "Error: Task with ID %d not found.\n", id)
			return
		}

		fmt.Printf("🔄 Task ID %d updated successfully.\n", id)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [task_id]",
	Short: "delete a task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Task ID must be an integer. %v\n", err)
			return
		}
		value, err := tm.DeleteTask(taskID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting task: %v\n", err)

			return
		}

		if !value {
			fmt.Fprintf(os.Stderr, "Error: Task with ID %d not found.\n", taskID)

			return
		}
		fmt.Printf("🗑️ Task ID %d deleted successfully.\n", taskID)
	},
}

var markTaskCmd = &cobra.Command{
	Use:   "mark [task_id] [task_status]",
	Short: "mark a task",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		taskStatus := args[1]
		taskID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(err)
		}
		switch taskStatus {
		case "TODO":
			err := tm.MarkTaskAsTodo(taskID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error marking task: %v\n", err)

				return
			}
		case "IN_PROGRESS":
			err := tm.MarkTaskAsInProgress(taskID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error marking task: %v\n", err)

				return
			}
		case "DONE":
			err := tm.MarkTaskAsDone(taskID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error marking task: %v\n", err)
				return
			}
		default:
			fmt.Println("Invalid task status")
			return
		}
		fmt.Printf("🏷️ Task ID %d marked as %s successfully.\n", taskID, taskStatus)
	},
}

var listTasksCmd = &cobra.Command{
	Use:   "list [status]",
	Short: "list tasks with different status (e.g., done, todo, all)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		status := "ALL"
		if len(args) == 1 {
			status = args[0]
		}

		var tasks []structures.Task

		switch status {
		case "ALL":
			tasks = tm.ListAllTasks()
		case "DONE":
			tasks = tm.ListDoneTasks()
		case "TODO":
			tasks = tm.ListTodoTasks()
		case "IN_PROGRESS":
			tasks = tm.ListInProgressTasks()
		default:
			fmt.Fprintf(os.Stderr, "Error: Invalid status '%s'. Use all, done, todo, or in_progress.\n", status)
			return
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.Header("ID", "Name", "Description", "Status", "Created", "Updated")
		for _, task := range tasks {
			tableRow := []string{strconv.Itoa(task.TaskId), task.TaskName, task.TaskDescription, task.TaskStatus, task.TaskCreatedAt, task.TaskUpdatedAt}
			err := table.Append(tableRow)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error appending row: %v\n", err)
			}
		}
		err := table.Render()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error rendering table: %v\n", err)
		}
		fmt.Printf("Listing %s tasks (Total: %d):\n", status, len(tasks))

	},
}

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "search tasks by query",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		tasks := tm.SearchTasks(query)
		if len(tasks) == 0 {
			fmt.Printf("No tasks found matching query '%s'.\n", query)
			return
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.Header("ID", "Name", "Description", "Status", "Created", "Updated")
		for _, task := range tasks {
			tableRow := []string{strconv.Itoa(task.TaskId), task.TaskName, task.TaskDescription, task.TaskStatus, task.TaskCreatedAt, task.TaskUpdatedAt}
			err := table.Append(tableRow)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error appending row: %v\n", err)
			}
		}
		err := table.Render()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error rendering table: %v\n", err)
		}
	},
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "clean tasks with DONE status",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(tm.ListDoneTasks()) == 0 {
			fmt.Println("No 'DONE tasks to clean'")
			return
		}
		if !ConfirmAction("Are you sure you want to delete all DONE tasks? This action is irreversible.") {
			fmt.Println("Operation cancelled")
			return
		}
		count, err := tm.CleanDoneTasks()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error cleaning tasks: %v\n", err)
			return
		}
		fmt.Printf("Successfully deleted %d DONE tasks.\n", count)
	},
}

func execute() {
	if err := mainCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	mainCmd.AddCommand(addCmd)
	mainCmd.AddCommand(updateCmd)
	mainCmd.AddCommand(deleteCmd)
	mainCmd.AddCommand(markTaskCmd)
	mainCmd.AddCommand(listTasksCmd)
	mainCmd.AddCommand(searchCmd)
	mainCmd.AddCommand(cleanCmd)
}

func main() {
	var err error

	tm, err = task_manager.NewTaskManager("tasks.json")
	if err != nil {
		fmt.Println("Error creating task manager:", err)
		os.Exit(1)
	}
	execute()
}

func ConfirmAction(promt string) bool {
	fmt.Printf("%s? [y/N] ", promt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())
	return strings.ToLower(input) == "y"
}
