package task_manager

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
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
	if strings.TrimSpace(name) == "" {
		return 0, fmt.Errorf("task name must not be empty")
	}
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
		return false, nil
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
		return false, nil
	}

	if name, exists := values["task_name"]; exists {
		task.TaskName = name
	}

	if description, exists := values["task_description"]; exists {
		task.TaskDescription = description
	}

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
func (taskManager *TaskManager) MarkTaskAsDone(id int) (bool, error) {
	if !taskManager.taskStatusHelper(id, "DONE") {
		return false, nil
	}
	if err := taskManager.SaveTasks(); err != nil {
		return false, fmt.Errorf("failed to save task status change: %w", err)
	}
	return true, nil
}

// MarkTaskAsInProgress - Метод для установки статуса "IN_PROGRESS"
func (taskManager *TaskManager) MarkTaskAsInProgress(id int) (bool, error) {
	if !taskManager.taskStatusHelper(id, "IN_PROGRESS") {
		return false, nil
	}
	if err := taskManager.SaveTasks(); err != nil {
		return false, fmt.Errorf("failed to save task status change: %w", err)
	}
	return true, nil
}

// MarkTaskAsTodo - Метод для установки статуса Toдo
func (taskManager *TaskManager) MarkTaskAsTodo(id int) (bool, error) {
	if !taskManager.taskStatusHelper(id, "TODO") {
		return false, nil
	}
	if err := taskManager.SaveTasks(); err != nil {
		return false, fmt.Errorf("failed to save task status change: %w", err)
	}
	return true, nil
}

func (taskManager *TaskManager) filterTaskByStatus(status string) []structures.Task {
	result := make([]structures.Task, 0, len(taskManager.Tasks))

	filterALL := status == "ALL"
	for _, task := range taskManager.Tasks {
		if filterALL || task.TaskStatus == status {
			result = append(result, task)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].TaskId < result[j].TaskId
	})
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

// SearchTasks - Метод, позволяющий находить нужные таски по подстрокам
func (taskManager *TaskManager) SearchTasks(query string) []structures.Task {
	tasks := make([]structures.Task, 0, len(taskManager.Tasks))
	for _, task := range taskManager.Tasks {
		if strings.Contains(strings.ToLower(task.TaskName), strings.ToLower(query)) || strings.Contains(strings.ToLower(task.TaskDescription), strings.ToLower(query)) {
			tasks = append(tasks, task)
		}
	}
	if len(tasks) == 0 {
		return tasks
	}
	return tasks

}

// CleanDoneTasks - очищает таски со статусом DONE
func (taskManager *TaskManager) CleanDoneTasks() (int, error) {
	var idsToDelete []int
	for _, task := range taskManager.Tasks {
		if task.TaskStatus == "DONE" {
			idsToDelete = append(idsToDelete, task.TaskId)
		}
	}

	count := len(idsToDelete)
	if count == 0 {
		return 0, nil
	}

	for _, id := range idsToDelete {
		delete(taskManager.Tasks, id)
	}

	if err := taskManager.SaveTasks(); err != nil {
		return 0, fmt.Errorf("failed to save tasks: %w", err)
	}
	return count, nil
}
