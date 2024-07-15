package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.File("youtube-downloader/index.html")
	})

	e.POST("/download", func(c echo.Context) error {
		videoURL := c.FormValue("videoUrl")
		thumbnailURL := getThumbnailURL(videoURL)
		if thumbnailURL == "" {
			return c.String(http.StatusBadRequest, "Invalid YouTube URL")
		}

		resp, err := http.Get(thumbnailURL)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to download thumbnail")
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatal("Error closing resp.Body", err)
				return
			}
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			return c.String(http.StatusInternalServerError, "Failed to download thumbnail")
		}

		c.Response().Header().Set("Content-Type", "image/jpeg")
		_, err = io.Copy(c.Response().Writer, resp.Body)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to write thumbnail")
		}
		return nil
	})

	err := e.Start(":8080")
	if err != nil {
		log.Fatal("Failed to start server!", err)
		return
	}
}

func getThumbnailURL(videoURL string) string {
	videoID := extractVideoID(videoURL)
	if videoID == "" {
		return ""
	}
	return fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", videoID)
}

func extractVideoID(videoURL string) string {
	if strings.Contains(videoURL, "youtu.be/") {
		parts := strings.Split(videoURL, "youtu.be/")
		if len(parts) > 1 {
			return parts[1]
		}
	} else if strings.Contains(videoURL, "watch?v=") {
		parts := strings.Split(videoURL, "watch?v=")
		if len(parts) > 1 {
			return parts[1]
		}
	}
	return ""
}
