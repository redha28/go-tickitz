package utils

import (
	"fmt"
	"gotickitz/pkg"
	"log"
	"mime/multipart"
	"os"
	fp "path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func (u *Utils) FileHandling(ctx *gin.Context, file *multipart.FileHeader, payload *pkg.Payload, oldFilename string) (filename, filepath string, err error) {
	ext := fp.Ext(file.Filename)
	allowedExt := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
	}
	if !allowedExt[ext] {
		return "", "", fmt.Errorf("ekstensi file tidak diizinkan")
	}

	// Buat nama file baru dengan timestamp
	filename = fmt.Sprintf("%d_%s_profile_image%s", time.Now().UnixNano(), payload.Id, ext)
	filepath = fp.Join("public", "img", filename)

	// Simpan file baru
	if err = ctx.SaveUploadedFile(file, filepath); err != nil {
		return "", "", err
	}

	// Hapus file lama jika ada
	if oldFilename != "" {
		oldPath := fp.Join("public", "img", oldFilename)
		if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
			log.Println("[WARNING] Gagal hapus file lama:", err)
		}
	}

	return filename, filepath, nil
}
