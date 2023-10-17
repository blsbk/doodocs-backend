package delivery

import (
	"archive/zip"
	"bytes"
	"doodocs-challenge/internal/models"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockArchiveUsecase struct{}

func (m *MockArchiveUsecase) GetArchiveInfo(fileHeader *multipart.FileHeader) models.ArchiveInfo {
	// Implement a mock version of GetArchiveInfo
	return models.ArchiveInfo{}
}

func (m *MockArchiveUsecase) GetFileInfo(zipReader *zip.Reader, res models.ArchiveInfo) models.ArchiveInfo {
	// Implement a mock version of GetFileInfo
	return models.ArchiveInfo{}
}

func (m *MockArchiveUsecase) SendFileToEmails(fileHeader *multipart.FileHeader, emails []string) error {
	// Implement a mock version of SendFileToEmails
	return nil
}

func (m *MockArchiveUsecase) GetFilePath(string) string {
	// Implement a mock version of GetFilePath
	return ""
}

// MockFile is a custom implementation of multipart.File.
type MockFile struct {
	Content io.Reader
}

func (m *MockFile) Read(p []byte) (n int, err error) {
	return m.Content.Read(p)
}

func (m *MockFile) Close() error {
	// You can implement a Close method if needed.
	return nil
}

func TestMockFile(t *testing.T) {
	// Create a mock file content as a byte slice
	mockFileContent := []byte("This is a mock file content.")

	// Create a buffer and write the mock content to it
	buffer := bytes.NewBuffer(mockFileContent)

	// Create a MockFile instance that wraps the buffer
	mockFile := &MockFile{Content: buffer}

	// Read from the MockFile
	readData := make([]byte, len(mockFileContent))
	n, err := mockFile.Read(readData)

	assert.NoError(t, err)
	assert.Equal(t, len(mockFileContent), n)
	assert.Equal(t, mockFileContent, readData)
}

func TestShowArchiveInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockArchiveUsecase := &MockArchiveUsecase{}
	archiveHandler := ArchiveHandler{AUsecase: mockArchiveUsecase}

	r.POST("/api/archive/information", archiveHandler.ShowArchiveInfo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/archive/information", nil)
	req.Header.Set("Content-Type", "multipart/form-data")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}

func TestCreateZipArchive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockArchiveUsecase := &MockArchiveUsecase{}
	archiveHandler := ArchiveHandler{AUsecase: mockArchiveUsecase}

	r.POST("/api/archive/files", archiveHandler.CreateZipArchive)

	// Create a test request with a file in a multipart form
	// You may need to create a mock file to simulate the form data

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/archive/files", nil)
	req.Header.Set("Content-Type", "multipart/form-data")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}

func TestSendFileByEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockArchiveUsecase := &MockArchiveUsecase{}
	archiveHandler := ArchiveHandler{AUsecase: mockArchiveUsecase}

	path := "/Users/blsbk/Desktop/go projects/doodocs/testfile.txt"

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", path)
	assert.NoError(t, err)
	sample, err := os.Open(path)
	assert.NoError(t, err)

	_, err = io.Copy(part, sample)
	assert.NoError(t, err)
	writer.WriteField("emails", "bagdat365@gmail.com,bagdatbilisbek@gmail.com")
	assert.NoError(t, writer.Close())

	r.POST("/api/mail/file", archiveHandler.SendFileByEmail)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "api/mail/file", body)
	req.Header.Set("Content-Type", "multipart/form-data")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
