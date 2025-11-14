package task_manager

import (
	"os"
	"testing"
)

// TestAddTask проверяет успешное добавление и обработку ошибок для пустых полей.
func TestAddTask(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "task-test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFilePath := tmpFile.Name()
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}
	defer func() {
		if rErr := os.Remove(tmpFilePath); rErr != nil {
			t.Logf("Failed to remove temp file %s: %v", tmpFilePath, rErr)
		}
	}()

	tm, err := NewTaskManager(tmpFilePath)
	if err != nil {
		t.Fatalf("Failed to create TaskManager: %v", err)
	}

	tests := []struct {
		name      string
		inputName string
		inputDesc string
		wantErr   bool
	}{
		{
			name:      "Success: Valid Task",
			inputName: "Test",
			inputDesc: "Test description",
			wantErr:   false,
		},
		{
			name:      "Failure: Empty Name",
			inputName: "",
			inputDesc: "Should fail validation",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			id, err := tm.AddTask(tt.inputName, tt.inputDesc)

			if (err != nil) != tt.wantErr {
				t.Fatalf("AddTask() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if id <= 0 {
					t.Errorf("AddTask returned invalid ID: %d, want > 0", id)
				}

				task, ok := tm.Tasks[id]
				if !ok || task.TaskName != tt.inputName || task.TaskDescription != tt.inputDesc {
					t.Errorf("Task not correctly stored or retrieved from map")
				}
			}
		})
	}
}

// TestDeleteTask - проверят успешное удаление и обработку ошибок для удаления несуществующих тасков
func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantOK  bool
		wantErr bool
	}{
		{
			name:    "Success: Valid Task",
			id:      1,
			wantOK:  true,
			wantErr: false,
		},
		{
			name:    "Failure: No task",
			id:      999,
			wantOK:  false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, _ := os.CreateTemp("", "task-test-*.json")
			tmpFilePath := tmpFile.Name()
			err := tmpFile.Close()
			if err != nil {
				return
			}
			defer func(name string) {
				err := os.Remove(name)
				if err != nil {
					t.Logf("Failed to remove temp file %s: %v", name, err)
					return
				}
			}(tmpFilePath)

			tm, err := NewTaskManager(tmpFilePath)
			if err != nil {
				t.Fatalf("Failed to create TaskManager: %v", err)
			}

			if tt.wantOK {
				_, err := tm.AddTask("test", "test")
				if err != nil {
					t.Fatalf("Failed to setup task for deletion: %v", err)
				}
			}
			ok, err := tm.DeleteTask(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTask() error = %v, wantErr %v", err, tt.wantErr)
			}
			if ok != tt.wantOK {
				t.Errorf("DeleteTask() ok = %v, wantOK %v", ok, tt.wantOK)
			}
			if ok {
				if _, exists := tm.Tasks[tt.id]; exists {
					t.Errorf("DeleteTask() deleted task with id %d", tt.id)
				}
			}
		})
	}
}

// TestUpdateTask проверяет обновление существующих задач, частичные обновления и ошибки.
func TestUpdateTask(t *testing.T) {
	tests := []struct {
		name          string
		taskID        int
		setupName     string
		setupDesc     string
		inputArgs     map[string]string
		wantOK        bool
		wantErr       bool
		wantFinalName string
		wantFinalDesc string
	}{
		{
			name:          "Success: Update Name Only",
			taskID:        1,
			setupName:     "Original Name",
			setupDesc:     "Original Description",
			inputArgs:     map[string]string{"task_name": "New Updated Name"},
			wantOK:        true,
			wantErr:       false,
			wantFinalName: "New Updated Name",
			wantFinalDesc: "Original Description",
		},
		{
			name:          "Success: Update Description Only",
			taskID:        1,
			setupName:     "Original Name",
			setupDesc:     "Original Description",
			inputArgs:     map[string]string{"task_description": "New Updated Description"},
			wantOK:        true,
			wantErr:       false,
			wantFinalName: "Original Name",
			wantFinalDesc: "New Updated Description",
		},
		{
			name:          "Success: Full Update",
			taskID:        1,
			setupName:     "Old Task",
			setupDesc:     "Old Description",
			inputArgs:     map[string]string{"task_name": "Finalized Task", "task_description": "Final Description"},
			wantOK:        true,
			wantErr:       false,
			wantFinalName: "Finalized Task",
			wantFinalDesc: "Final Description",
		},
		{
			name:      "Failure: Non-existent ID",
			taskID:    999,
			setupName: "", setupDesc: "",
			inputArgs:     map[string]string{"task_name": "Irrelevant"},
			wantOK:        false,
			wantErr:       false,
			wantFinalName: "", wantFinalDesc: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "task-update-*.json")
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
			tmpFilePath := tmpFile.Name()
			if err := tmpFile.Close(); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
			defer func() {
				err := os.Remove(tmpFilePath)
				if err != nil {
					return
				}
			}()

			tm, err := NewTaskManager(tmpFilePath)
			if err != nil {
				t.Fatalf("Failed to create TaskManager: %v", err)
			}

			if tt.wantOK {
				_, err := tm.AddTask(tt.setupName, tt.setupDesc)
				if err != nil {
					t.Fatalf("Failed to setup task for update: %v", err)
				}
			}

			ok, err := tm.UpdateTask(tt.taskID, tt.inputArgs)

			if (err != nil) != tt.wantErr {
				t.Fatalf("UpdateTask() error = %v, wantErr %v", err, tt.wantErr)
			}

			if ok != tt.wantOK {
				t.Fatalf("UpdateTask() returned ok=%v, want %v", ok, tt.wantOK)
			}

			if ok {
				updatedTask, exists := tm.Tasks[tt.taskID]
				if !exists {
					t.Fatalf("Task unexpectedly disappeared after successful update")
				}

				if updatedTask.TaskName != tt.wantFinalName {
					t.Errorf("TaskName mismatch: got %s, want %s", updatedTask.TaskName, tt.wantFinalName)
				}

				if updatedTask.TaskDescription != tt.wantFinalDesc {
					t.Errorf("TaskDescription mismatch: got %s, want %s", updatedTask.TaskDescription, tt.wantFinalDesc)
				}
			}
		})
	}
}

// TestMarkTaskStatus - проверяет, изменение статуса у определенной задачи и ошибки.
func TestMarkTaskStatus(t *testing.T) {
	tests := []struct {
		name            string
		taskID          int
		setupName       string
		setupDesc       string
		inputStatus     string
		wantOK          bool
		wantErr         bool
		wantFinalName   string
		wantFinalDesc   string
		wantFinalStatus string
	}{

		{
			name:            "Success: Mark Status as DONE",
			taskID:          1,
			setupName:       "Original Name",
			setupDesc:       "Original Description",
			inputStatus:     "DONE",
			wantOK:          true,
			wantErr:         false,
			wantFinalName:   "Original Name",
			wantFinalDesc:   "Original Description",
			wantFinalStatus: "DONE",
		},

		{
			name:            "Success: Mark Status as TODO",
			taskID:          1,
			setupName:       "Original Name",
			setupDesc:       "Original Description",
			inputStatus:     "TODO",
			wantOK:          true,
			wantErr:         false,
			wantFinalName:   "Original Name",
			wantFinalDesc:   "Original Description",
			wantFinalStatus: "TODO",
		},

		{
			name:            "Success: Mark Status as In-Progress",
			taskID:          1,
			setupName:       "Original Name",
			setupDesc:       "Original Description",
			inputStatus:     "IN_PROGRESS",
			wantOK:          true,
			wantErr:         false,
			wantFinalName:   "Original Name",
			wantFinalDesc:   "Original Description",
			wantFinalStatus: "IN_PROGRESS",
		},

		{
			name:            "Failure: Non-existent Task",
			taskID:          999,
			setupName:       "N/A",
			setupDesc:       "N/A",
			inputStatus:     "DONE",
			wantOK:          false,
			wantErr:         false,
			wantFinalName:   "",
			wantFinalDesc:   "",
			wantFinalStatus: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tmpFile, err := os.CreateTemp("", "task-update-*.json")
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
			tmpFilePath := tmpFile.Name()
			if err := tmpFile.Close(); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
			defer func() {

				if rErr := os.Remove(tmpFilePath); rErr != nil {
					t.Logf("Failed to remove temp file %s: %v", tmpFilePath, rErr)
				}
			}()

			tm, err := NewTaskManager(tmpFilePath)
			if err != nil {
				t.Fatalf("Failed to create TaskManager: %v", err)
			}

			if tt.wantOK {
				_, err := tm.AddTask(tt.setupName, tt.setupDesc)
				if err != nil {
					t.Fatalf("Failed to setup task for update: %v", err)
				}
			}

			var ok bool
			var markErr error

			switch tt.inputStatus {
			case "DONE":
				ok, markErr = tm.MarkTaskAsDone(tt.taskID)
			case "IN_PROGRESS":
				ok, markErr = tm.MarkTaskAsInProgress(tt.taskID)
			case "TODO":
				ok, markErr = tm.MarkTaskAsTodo(tt.taskID)
			default:
				t.Fatalf("Test setup error: Invalid input status %s", tt.inputStatus)
			}

			if (markErr != nil) != tt.wantErr {
				t.Fatalf("MarkTask() error = %v, wantErr %v", markErr, tt.wantErr)
			}

			if ok != tt.wantOK {
				t.Fatalf("MarkTask() returned ok=%v, want %v", ok, tt.wantOK)
			}

			if ok {
				updatedTask, exists := tm.Tasks[tt.taskID]
				if !exists {
					t.Fatalf("Task unexpectedly disappeared after successful update status")
				}

				if updatedTask.TaskName != tt.wantFinalName {
					t.Errorf("TaskName mismatch: got %s, want %s", updatedTask.TaskName, tt.wantFinalName)
				}

				if updatedTask.TaskStatus != tt.wantFinalStatus {
					t.Errorf("TaskStatus mismatch: got %s, want %s", updatedTask.TaskStatus, tt.wantFinalStatus)
				}
			}
		})
	}
}

// TestCleanDoneTasks - проверяет корректность удаления задач со статусом DONE.
func TestCleanDoneTasks(t *testing.T) {
	tests := []struct {
		name           string
		tasksToSetup   map[int]string
		wantCount      int
		tasksRemaining int
		wantErr        bool
	}{
		{
			name:           "Success: Deleting multiple tasks, keeping one",
			tasksToSetup:   map[int]string{1: "DONE", 2: "TODO", 3: "DONE", 4: "IN_PROGRESS"},
			wantCount:      2,
			tasksRemaining: 2,
			wantErr:        false,
		},
		{
			name:           "Success: Deleting all tasks",
			tasksToSetup:   map[int]string{1: "DONE", 2: "DONE"},
			wantCount:      2,
			tasksRemaining: 0,
			wantErr:        false,
		},
		{
			name:           "Success: Empty manager (count zero)",
			tasksToSetup:   map[int]string{},
			wantCount:      0,
			tasksRemaining: 0,
			wantErr:        false,
		},
		{
			name:           "Success: Nothing to delete",
			tasksToSetup:   map[int]string{1: "TODO", 2: "IN_PROGRESS"},
			wantCount:      0,
			tasksRemaining: 2,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tmpFile, err := os.CreateTemp("", "task-clean-*.json")
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
			tmpFilePath := tmpFile.Name()
			if err := tmpFile.Close(); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			t.Cleanup(func() {
				err := os.Remove(tmpFilePath)
				if err != nil {
					return
				}
			})

			tm, err := NewTaskManager(tmpFilePath)
			if err != nil {
				t.Fatalf("Failed to create TaskManager: %v", err)
			}

			for _, status := range tt.tasksToSetup {

				newID, err := tm.AddTask("Task "+status, "Status: "+status)
				if err != nil {
					t.Fatalf("Setup AddTask failed: %v", err)
				}

				if task, ok := tm.Tasks[newID]; ok {
					// Изменяем копию
					task.TaskStatus = status

					tm.Tasks[newID] = task
				}
			}

			count, err := tm.CleanDoneTasks()

			if (err != nil) != tt.wantErr {
				t.Fatalf("CleanDoneTasks() error = %v, wantErr %v", err, tt.wantErr)
			}

			if count != tt.wantCount {
				t.Errorf("CleanDoneTasks() returned count %d, want %d", count, tt.wantCount)
			}

			if len(tm.Tasks) != tt.tasksRemaining {
				t.Errorf("CleanDoneTasks() left %d tasks, want %d", len(tm.Tasks), tt.tasksRemaining)
			}

			for _, task := range tm.Tasks {
				if task.TaskStatus == "DONE" {
					t.Errorf("CleanDoneTasks failed: Task ID %d with status DONE was not deleted.", task.TaskId)
				}
			}
		})
	}
}
