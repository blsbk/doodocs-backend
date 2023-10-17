package usecase

import (
	"archive/zip"
	"doodocs-challenge/internal/models"
	"mime"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetArchiveInfo(t *testing.T) {
	// Create a mock file header for testing
	fileHeader := &multipart.FileHeader{
		Filename: "test.zip",
		Size:     1024,
	}

	archiveUsecase := NewArchiveUsecase()
	archiveInfo := archiveUsecase.GetArchiveInfo(fileHeader)

	assert.Equal(t, "test.zip", archiveInfo.Filename)
	assert.Equal(t, float64(1024), archiveInfo.ArchiveSize)
}

func TestGetFileInfo(t *testing.T) {
	// Create a mock zip.Reader for testing
	mockZipFile1 := &zip.File{
		FileHeader: zip.FileHeader{
			Name:               "file1.txt",
			UncompressedSize64: 100,
		},
	}
	mockZipFile2 := &zip.File{
		FileHeader: zip.FileHeader{
			Name:               "file2.jpg",
			UncompressedSize64: 200,
		},
	}
	mockZipReader := &zip.Reader{
		File: []*zip.File{mockZipFile1, mockZipFile2},
	}

	archiveUsecase := NewArchiveUsecase()
	archiveInfo := models.ArchiveInfo{}
	updatedArchiveInfo := archiveUsecase.GetFileInfo(mockZipReader, archiveInfo)

	// Assert that the method updated the ArchiveInfo correctly
	assert.Len(t, updatedArchiveInfo.Files, 2)
	assert.Equal(t, "file1.txt", updatedArchiveInfo.Files[0].FilePath)
	assert.Equal(t, "file2.jpg", updatedArchiveInfo.Files[1].FilePath)
	assert.Equal(t, float64(100), updatedArchiveInfo.Files[0].Size)
	assert.Equal(t, float64(200), updatedArchiveInfo.Files[1].Size)
	assert.Equal(t, mime.TypeByExtension(".txt"), updatedArchiveInfo.Files[0].MimeType)
	assert.Equal(t, mime.TypeByExtension(".jpg"), updatedArchiveInfo.Files[1].MimeType)
	assert.Equal(t, float64(300), updatedArchiveInfo.TotalSize)
	assert.Equal(t, 2, updatedArchiveInfo.TotalFiles)
}

func TestGetFilePath(t *testing.T) {
	archiveUsecase := NewArchiveUsecase()
	filePath := archiveUsecase.GetFilePath("dir/subdir/file.txt")
	assert.Equal(t, "subdir/file.txt", filePath)
}

func TestSendFileToEmails(t *testing.T) {
	// Create a mock file header and emails for testing
	fileHeader := &multipart.FileHeader{
		Filename: "test.txt",
		Size:     100,
	}
	emails := []string{"recipient1@example.com", "recipient2@example.com"}

	archiveUsecase := NewArchiveUsecase()
	err := archiveUsecase.SendFileToEmails(fileHeader, emails)

	assert.NoError(t, err)

}
