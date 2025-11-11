package task_manager

import (
	"time"

	"github.com/TaskTrackerCLI/structures"
)

type TaskManager struct {
	Tasks  map[int]structures.Task
	nextId int
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		Tasks:  make(map[int]structures.Task),
		nextId: 1,
	}

}

// AddTask Метод добавления нового таска с данными(имя, описание)
func (taskManager *TaskManager) AddTask(name, description string) int {
	newTask := structures.Task{
		TaskId:          taskManager.nextId,
		TaskName:        name,
		TaskDescription: description,
		TaskStatus:      "TODO",
		TaskCreatedAt:   time.Now().Format(time.RFC3339),
	}
	taskManager.Tasks[taskManager.nextId] = newTask
	taskManager.nextId++
	return taskManager.nextId - 1
}

// DeleteTask Метод удаления таска с id
func (taskManager *TaskManager) DeleteTask(id int) int {
	_, ok := taskManager.Tasks[id]
	if !ok {
		return 0
	}
	delete(taskManager.Tasks, id)
	taskManager.nextId--
	return id

}

// UpdateTask - Метод обновления данных(имя, описание) у таски с id
func (taskManager *TaskManager) UpdateTask(id int, values map[string]string) bool {
	task, ok := taskManager.Tasks[id]
	if !ok {
		return false // Возвращаем false, если таска не была найдена
	}
	// Обновлять имя только если оно было передано
	if name, exists := values["task_name"]; exists {
		task.TaskName = name
	}
	// Обновить описание только если оно было передано
	if description, exists := values["task_description"]; exists {
		task.TaskDescription = description
	}
	// Меняем время обновления таска
	task.TaskUpdatedAt = time.Now().Format(time.RFC3339)
	taskManager.Tasks[id] = task
	return true
}

func (taskManager *TaskManager) taskStatusHelper(id int, newStatus string) bool {
	task, ok := taskManager.Tasks[id]
	if !ok {
		return false
	}
	task.TaskStatus = newStatus
	task.TaskUpdatedAt = time.Now().Format(time.RFC3339)
	taskManager.Tasks[id] = task
	return true
}

// MarkTaskAsDone - Метод для установки статуса "DONE"
func (taskManager *TaskManager) MarkTaskAsDone(id int) bool {
	return taskManager.taskStatusHelper(id, "DONE")
}

// MarkTaskAsInProgress - Метод для установки статуса "IN_PROGRESS"
func (taskManager *TaskManager) MarkTaskAsInProgress(id int) bool {
	return taskManager.taskStatusHelper(id, "IN_PROGRESS")
}

// MarkTaskAsTodo - Метод для установки статуса "TODO"
func (taskManager *TaskManager) MarkTaskAsTodo(id int) bool {
	return taskManager.taskStatusHelper(id, "TODO")
}

func (taskManager *TaskManager) filterTaskByStatus(status string) []structures.Task {
	result := make([]structures.Task, 0, len(taskManager.Tasks))

	filterALL := status == "ALL"
	for _, task := range taskManager.Tasks {
		if filterALL || task.TaskStatus == status {
			result = append(result, task)
		}
	}
	return result
}

func (taskManager *TaskManager) ListDoneTasks() []structures.Task {
	return taskManager.filterTaskByStatus("DONE")
}

func (taskManager *TaskManager) ListTodoTasks() []structures.Task {
	return taskManager.filterTaskByStatus("TODO")
}

func (taskManager *TaskManager) ListInProgressTasks() []structures.Task {
	return taskManager.filterTaskByStatus("IN_PROGRESS")
}

func (taskManager *TaskManager) ListAllTasks() []structures.Task {
	return taskManager.filterTaskByStatus("ALL")
}
