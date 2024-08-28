package toolkit

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const randomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Tools struct {
	MaxFileSize      int32
	AllowedFileTypes []string
}

func (t *Tools) RandomString(length int) string {
	s, r := make([]rune, length), []rune(randomStringSource)
	randomStringSourceLen := len(r)
	for i := range s {
		prime, err := rand.Prime(rand.Reader, randomStringSourceLen)
		if err != nil {
			return ""
		}
		x, y := prime.Uint64(), uint64(randomStringSourceLen)
		s[i] = r[x%y]
	}
	return string(s)
}

type UploadedFile struct {
	NewFileName      string
	OriginalFileName string
	Filesize         int64
}

func (t *Tools) UploadFile(request *http.Request, uploadDir string, rename ...bool) (*UploadedFile, error) {
	renameFile := true

	if len(rename) > 0 {
		renameFile = rename[0]
	}

	files, err := t.UploadFiles(request, uploadDir, renameFile)
	if err != nil {
		return nil, err
	}
	return files[0], nil

}

func (t *Tools) UploadFiles(request *http.Request, uploadDir string, rename ...bool) ([]*UploadedFile, error) {
	renameFile := true
	if len(rename) > 0 {
		renameFile = rename[0]
	}

	var uploadedFiles []*UploadedFile

	if t.MaxFileSize == 0 {
		t.MaxFileSize = 1024 * 1024 * 1024
	}

	err := request.ParseMultipartForm(int64(t.MaxFileSize))
	if err != nil {
		return uploadedFiles, errors.New("file Size is too big")
	}

	for _, fHeader := range request.MultipartForm.File {
		for _, hdr := range fHeader {
			uploadedFiles, err = func(uploadedFiles []*UploadedFile) ([]*UploadedFile, error) {
				var uploadedFile UploadedFile
				file, err := hdr.Open()
				if err != nil {
					return nil, err
				}
				defer file.Close()

				buff := make([]byte, 512)
				_, err = file.Read(buff)
				if err != nil {
					return nil, err
				}

				allowed := false

				fileType := http.DetectContentType(buff)

				if len(t.AllowedFileTypes) > 0 {
					for _, typeOfAllowedFile := range t.AllowedFileTypes {
						if strings.EqualFold(typeOfAllowedFile, fileType) {
							allowed = true
						}
					}
				} else {
					allowed = true
				}
				if !allowed {
					return nil, errors.New("file type is not allowed")
				}
				_, err = file.Seek(0, 0)
				if err != nil {
					return nil, err
				}
				uploadedFile.OriginalFileName = hdr.Filename
				if renameFile {
					uploadedFile.NewFileName = fmt.Sprintf("%s%s", t.RandomString(25), filepath.Ext(hdr.Filename))
				} else {
					uploadedFile.NewFileName = hdr.Filename
				}

				var outfile *os.File
				defer outfile.Close()

				if outfile, err = os.Create(filepath.Join(uploadDir, uploadedFile.NewFileName)); err != nil {
					return nil, err
				} else {
					fileSize, err := io.Copy(outfile, file)
					if err != nil {
						return nil, err
					}
					uploadedFile.Filesize = fileSize
				}
				uploadedFiles = append(uploadedFiles, &uploadedFile)
				return uploadedFiles, nil

			}(uploadedFiles)

			if err != nil {
				return uploadedFiles, err
			}
		}
	}
	return uploadedFiles, nil
}
