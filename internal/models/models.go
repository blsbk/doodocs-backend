package models

import (
	"archive/zip"
	"mime/multipart"
)

type ArchiveUsecases interface {
	GetArchiveInfo(*multipart.FileHeader) ArchiveInfo
	GetFileInfo(*zip.Reader, ArchiveInfo) ArchiveInfo
	SendFileToEmails(*multipart.FileHeader, []string) error
	GetFilePath(string) string
}

type ArchiveInfo struct {
	Filename    string     `json:"filename"`
	ArchiveSize float64    `json:"archive_size"`
	TotalSize   float64    `json:"total_size"`
	TotalFiles  float64    `json:"total_files"`
	Files       []FileInfo `json:"files"`
}

type FileInfo struct {
	FilePath string  `json:"file_path"`
	Size     float64 `json:"size"`
	MimeType string  `json:"mimetype"`
}
