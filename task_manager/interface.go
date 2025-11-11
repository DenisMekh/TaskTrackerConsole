package task_manager

import "github.com/TaskTrackerCLI/structures"

type TaskManagerInterface interface {
	AddTask(name, description string) int
	UpdateTask(id int, values map[string]string) bool
	DeleteTask(id int) int
	MarkTaskAsDone(id int) bool
	MarkTaskAsInProgress(id int) bool
	MarkTaskAsTodo(id int) bool
	ListAllTasks() []structures.Task
	ListDoneTasks() []structures.Task
	ListInProgressTasks() []structures.Task
	ListTodoTasks() []structures.Task
}
