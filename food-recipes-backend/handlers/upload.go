package handlers

import (
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	// Initialize Cloudinary client
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cloudinary config error"})
		return
	}

	// Try multi-file upload first: fields files|media|image|file (max 3)
	if form, err := c.MultipartForm(); err == nil && form != nil {
		files := form.File["files"]
		if len(files) == 0 { files = form.File["media"] }
		if len(files) == 0 { files = form.File["image"] }
		if len(files) == 0 { files = form.File["file"] }
		if len(files) > 0 {
			if len(files) > 3 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "You can upload up to 3 files only"})
				return
			}
			urls := make([]string, 0, len(files))
			items := make([]map[string]string, 0, len(files))
			for _, fh := range files {
				src, err := fh.Open()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
					return
				}
				up, err := cld.Upload.Upload(c, src, uploader.UploadParams{ResourceType: "auto"})
				_ = src.Close()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
					return
				}
				urls = append(urls, up.SecureURL)
				items = append(items, map[string]string{"url": up.SecureURL, "type": up.ResourceType})
			}
			if len(urls) == 1 {
				c.JSON(http.StatusOK, gin.H{"url": urls[0], "item": items[0]})
			} else {
				c.JSON(http.StatusOK, gin.H{"urls": urls, "items": items})
			}
			return
		}
	}

	// Fallback: single file under common keys
	file, err := c.FormFile("image")
	if err != nil {
		if file, err = c.FormFile("file"); err != nil {
			if file, err = c.FormFile("media"); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
				return
			}
		}
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	up, err := cld.Upload.Upload(c, src, uploader.UploadParams{ResourceType: "auto"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": up.SecureURL, "item": map[string]string{"url": up.SecureURL, "type": up.ResourceType}})
}
