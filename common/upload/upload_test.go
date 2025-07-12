package upload

import (
	"math/rand"
	"os"
	"testing"

	"mira/config"
)

func TestUpload_SaveToLocal(t *testing.T) {
	// Mock config for testing
	config.Data = &config.Config{}
	config.Data.Ruoyi.Name = "test-project"
	config.Data.Ruoyi.Domain = "localhost"
	config.Data.Ruoyi.UploadPath = "uploads/"

	// Create a temporary directory for uploads
	tempDir := "uploads/20230101/"
	if err := os.MkdirAll(tempDir, 0o755); err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll("uploads")

	upload := New(
		SetSavePath(tempDir),
		SetUrlPath(tempDir),
	)

	fileContent := []byte("hello world")
	file := &File{
		FileName:    "test.txt",
		FileSize:    len(fileContent),
		FileType:    "txt",
		FileContent: fileContent,
	}

	result, err := upload.SetFile(file).Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if result.OriginalName != "test.txt" {
		t.Errorf("OriginalName = %s; want test.txt", result.OriginalName)
	}

	if _, err := os.Stat(tempDir + result.FileName); os.IsNotExist(err) {
		t.Errorf("File not saved to %s", tempDir+result.FileName)
	}
}

func TestUpload_SaveToLocal_WithRandomName(t *testing.T) {
	// Mock config for testing
	config.Data = &config.Config{}
	config.Data.Ruoyi.Name = "test-project"
	config.Data.Ruoyi.Domain = "localhost"
	config.Data.Ruoyi.UploadPath = "uploads/"

	// Create a temporary directory for uploads
	tempDir := "uploads/20230101/"
	if err := os.MkdirAll(tempDir, 0o755); err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll("uploads")

	upload := New(
		SetSavePath(tempDir),
		SetUrlPath(tempDir),
		SetRandomName(true),
	)
	upload.rand = rand.New(rand.NewSource(1)) // Seed for predictable random name

	fileContent := []byte("hello world")
	file := &File{
		FileName:    "test.txt",
		FileSize:    len(fileContent),
		FileType:    "txt",
		FileContent: fileContent,
	}

	result, err := upload.SetFile(file).Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Generate the expected random name
	r := rand.New(rand.NewSource(1))
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var randomName string
	for i := 0; i < 64; i++ {
		randomChar := chars[r.Intn(len(chars))]
		randomName = randomName + string(randomChar)
	}
	expectedFileName := randomName + ".txt"

	if result.FileName != expectedFileName {
		t.Errorf("FileName = %s; want %s", result.FileName, expectedFileName)
	}

	if _, err := os.Stat(tempDir + result.FileName); os.IsNotExist(err) {
		t.Errorf("File not saved to %s", tempDir+result.FileName)
	}
}
