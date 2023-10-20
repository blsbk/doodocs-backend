package usecase

import (
	"archive/zip"
	"doodocs-challenge/internal/models"
	"io"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"

	"gopkg.in/gomail.v2"
)

type archiveUsecase struct {
}

func NewArchiveUsecase() models.ArchiveUsecases {
	return &archiveUsecase{}
}

func (a *archiveUsecase) GetArchiveInfo(fileHeader *multipart.FileHeader) models.ArchiveInfo {
	res := models.ArchiveInfo{
		Filename:    fileHeader.Filename,
		ArchiveSize: float64(fileHeader.Size),
	}
	return res
}

func (a *archiveUsecase) GetFileInfo(zipReader *zip.Reader, res models.ArchiveInfo) models.ArchiveInfo {
	// zipReader, err := zip.NewReader(file, archiveSize)
	// if err != nil {
	// 	return models.ArchiveInfo{}
	// }

	for _, zipFile := range zipReader.File {

		filename := filepath.Base(zipFile.Name)

		if !zipFile.FileInfo().IsDir() && !strings.HasPrefix(filename, ".") {
			res.Files = append(res.Files, models.FileInfo{
				FilePath: a.GetFilePath(zipFile.Name),
				Size:     float64(zipFile.UncompressedSize64),
				MimeType: mime.TypeByExtension(filepath.Ext(zipFile.Name)),
			})
			res.TotalSize += float64(zipFile.UncompressedSize64)
			res.TotalFiles++
		}

	}
	return res
}

func (a *archiveUsecase) GetFilePath(filename string) string {
	filePath := ""
	for i, v := range filename {
		if v == '/' {
			filePath = filename[i+1:]
			break
		}
	}
	return filePath
}

func (a *archiveUsecase) SendFileToEmails(fileHeader *multipart.FileHeader, emails []string) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	message := gomail.NewMessage()
	message.SetHeader("From", "bagdat365@gmail.com")
	message.SetHeader("Subject", "File Attachment")
	message.SetBody("text/plain", "Please find the attached file.")

	message.Attach(fileHeader.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := io.Copy(w, file)
		return err
	}))

	d := gomail.NewDialer("smtp.gmail.com", 587, "your-emails", "your-password")

	for _, recipient := range emails {
		message.SetHeader("To", recipient)
		if err := d.DialAndSend(message); err != nil {
			return err
		}
		if _, err := file.Seek(0, 0); err != nil {
			return err
		}
	}
	return nil
}
