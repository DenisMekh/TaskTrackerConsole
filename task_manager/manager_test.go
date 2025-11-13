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
