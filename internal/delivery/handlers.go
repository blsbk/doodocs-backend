package delivery

import (
	"archive/zip"
	"bytes"
	"doodocs-challenge/internal/models"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ArchiveHandler struct {
	AUsecase models.ArchiveUsecases
}

func NewArchiveHandler(r *gin.Engine, us models.ArchiveUsecases) {
	handler := &ArchiveHandler{
		AUsecase: us,
	}

	r.POST("/api/archive/information", handler.ShowArchiveInfo)
	r.POST("/api/archive/files", handler.CreateZipArchive)
	r.POST("/api/mail/file", handler.SendFileByEmail)
}

func (h *ArchiveHandler) ShowArchiveInfo(context *gin.Context) {
	err := context.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form data"})
		return
	}

	fileHeader, err := context.FormFile("file")
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "wrong content-disposition name"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	defer file.Close()

	fileData := make([]byte, 512)
	_, err = file.Read(fileData)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error reading zip file"})
		return
	}

	if http.DetectContentType(fileData) != "application/zip" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "file is not zip-archive"})
		return
	}

	res := h.AUsecase.GetArchiveInfo(fileHeader)

	zipReader, err := zip.NewReader(file, int64(res.ArchiveSize))
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error reading files in the archive"})
		return
	}

	res = h.AUsecase.GetFileInfo(zipReader, res)

	context.IndentedJSON(http.StatusOK, res)
}

func (h *ArchiveHandler) CreateZipArchive(context *gin.Context) {

	form, err := context.MultipartForm()
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form data"})
		return
	}

	files := form.File["files[]"]
	if len(files) == 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No files provided"})
		return
	}

	zipBuffer := new(bytes.Buffer)

	zipWriter := zip.NewWriter(zipBuffer)

	for _, fileHeader := range files {
		if !isValidMIMEType(fileHeader, allowedFileMIMETypes) {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid MIME type for file: " + fileHeader.Filename})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file: " + fileHeader.Filename})
			return
		}

		defer file.Close()

		zipFile, err := zipWriter.Create(fileHeader.Filename)
		if err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file in ZIP archive"})
			return
		}

		_, err = io.Copy(zipFile, file)
		if err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy file data to ZIP archive"})
			return
		}
	}
	err = zipWriter.Close()
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "closing writer"})
		return
	}

	// os.WriteFile("archive.zip", zipBuffer.Bytes(), 0777)
	context.Header("Content-Type", "application/zip")
	context.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", "archive.zip"))
	context.Data(http.StatusOK, "application/zip", zipBuffer.Bytes())

}

var allowedFileMIMETypes = []string{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"application/xml",
	"image/jpeg",
	"image/png",
}

var allowedEmailMIMETypes = []string{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"application/pdf",
}

func isValidMIMEType(fileHeader *multipart.FileHeader, types []string) bool {
	for _, allowedType := range types {
		if strings.HasPrefix(fileHeader.Header.Get("Content-Type"), allowedType) {
			return true
		}
	}
	return false
}

func (h *ArchiveHandler) SendFileByEmail(context *gin.Context) {
	form, err := context.MultipartForm()
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form data"})
		return
	}

	filesHeader := form.File["file"]
	if len(filesHeader) == 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	uploadedFile := filesHeader[0]

	if !isValidMIMEType(uploadedFile, allowedEmailMIMETypes) {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid MIME type for file: "})
		return
	}

	emails := context.PostForm("emails")
	emailList := strings.Split(emails, ",")

	fmt.Println(emailList)

	if err := h.AUsecase.SendFileToEmails(uploadedFile, emailList); err != nil {
		fmt.Print(err)
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"message": "File successfuly sent to email addresses"})

}
