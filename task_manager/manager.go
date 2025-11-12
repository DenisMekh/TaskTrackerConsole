package task_manager

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/TaskTrackerCLI/structures"
)

type TaskManager struct {
	Tasks    map[int]structures.Task
	nextId   int
	FilePath string
}

func NewTaskManager(filePath string) (*TaskManager, error) {
	taskManager := &TaskManager{
		Tasks:    make(map[int]structures.Task),
		FilePath: filePath,
		nextId:   1}
	if err := taskManager.LoadTasks(); err != nil {
		return taskManager, err
	}
	return taskManager, nil

}

// AddTask Метод добавления нового таска с данными(имя, описание)
func (taskManager *TaskManager) AddTask(name, description string) (int, error) {
	newTask := structures.Task{
		TaskId:          taskManager.nextId,
		TaskName:        name,
		TaskDescription: description,
		TaskStatus:      "TODO",
		TaskCreatedAt:   time.Now().Format(time.RFC3339),
	}
	taskManager.Tasks[taskManager.nextId] = newTask
	taskManager.nextId++
	err := taskManager.SaveTasks()
	if err != nil {
		return 0, err
	}
	return taskManager.nextId - 1, nil
}

// DeleteTask Метод удаления таска с id
func (taskManager *TaskManager) DeleteTask(id int) (bool, error) {
	_, ok := taskManager.Tasks[id]
	if !ok {
		return false, fmt.Errorf("task Not Found")
	}
	delete(taskManager.Tasks, id)
	err := taskManager.SaveTasks()
	if err != nil {
		return false, err
	}
	return ok, nil

}

// UpdateTask - Метод обновления данных(имя, описание) у таски с id
func (taskManager *TaskManager) UpdateTask(id int, values map[string]string) (bool, error) {
	task, ok := taskManager.Tasks[id]
	if !ok {
		return false, fmt.Errorf("task not found") // Возвращаем false, если таска не была найдена
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
	err := taskManager.SaveTasks()
	if err != nil {
		return false, err
	}
	return true, nil
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
func (taskManager *TaskManager) MarkTaskAsDone(id int) error {
	if !taskManager.taskStatusHelper(id, "DONE") {
		return fmt.Errorf("task with id %d not found", id)
	}
	if err := taskManager.SaveTasks(); err != nil {
		return fmt.Errorf("failed to save task status change: %w", err)
	}
	return nil
}

// MarkTaskAsInProgress - Метод для установки статуса "IN_PROGRESS"
func (taskManager *TaskManager) MarkTaskAsInProgress(id int) error {
	if !taskManager.taskStatusHelper(id, "IN_PROGRESS") {
		return fmt.Errorf("task with id %d not found", id)
	}
	if err := taskManager.SaveTasks(); err != nil {
		return fmt.Errorf("failed to save task status change: %w", err)
	}
	return nil
}

// MarkTaskAsTodo - Метод для установки статуса Toдo
func (taskManager *TaskManager) MarkTaskAsTodo(id int) error {
	if !taskManager.taskStatusHelper(id, "TODO") {
		return fmt.Errorf("task with id %d not found", id)
	}
	if err := taskManager.SaveTasks(); err != nil {
		return fmt.Errorf("failed to save task status change: %w", err)
	}
	return nil
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

// SaveTasks - метод для сохранения созданных, обновленных, удаленных тасков в json
func (taskManager *TaskManager) SaveTasks() error {
	tasks, err := json.Marshal(taskManager.Tasks)
	if err != nil {
		return fmt.Errorf("failed to write tasks file: %w", err)
	}
	err = os.WriteFile(taskManager.FilePath, tasks, 0644)
	if err != nil {
		return fmt.Errorf("failed to write tasks file: %w", err)
	}
	return nil
}

// LoadTasks - метод для загрузки тасков из json файла
func (taskManager *TaskManager) LoadTasks() error {
	fileContent, err := os.ReadFile(taskManager.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			taskManager.nextId = 1
			return nil
		}
		return fmt.Errorf("failed to read tasks: %w", err)
	}

	if len(fileContent) == 0 {
		taskManager.nextId = 1
		return nil
	}

	err = json.Unmarshal(fileContent, &taskManager.Tasks)
	if err != nil {
		return fmt.Errorf("failed to unmarshal tasks: %w", err)
	}
	maxID := 0
	for _, task := range taskManager.Tasks {
		if task.TaskId > maxID {
			maxID = task.TaskId
		}
	}
	taskManager.nextId = maxID + 1
	return nil
}
