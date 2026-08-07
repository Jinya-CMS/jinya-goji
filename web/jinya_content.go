package web

import (
	"log/slog"
	"math"
	"net/http"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/h2non/bimg"
	"jinya.de/goji/database"
	"jinya.de/goji/storage"
)

func serveConvertedImage(targetPath string, targetWidth *int, targetFileType *bimg.ImageType, w http.ResponseWriter) {
	slog.Debug("Read image file into buffer", "path", targetPath)
	buffer, err := os.ReadFile(targetPath)
	if err != nil {
		slog.Warn("Error reading file", "path", targetPath, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newImage := bimg.NewImage(buffer)

	slog.Debug("Determine image size", "path", targetPath)
	size, err := newImage.Size()
	if err != nil {
		slog.Warn("Error getting image size", "path", targetPath, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := buffer

	if targetWidth != nil && *targetWidth < 7680 {
		slog.Debug("Calculate the image aspect ratio and calculate the target width", "path", targetPath)
		ratio := float64(size.Width) / float64(size.Height)
		targetHeight := int(math.Round(float64(*targetWidth) / ratio))

		slog.Debug("Calculated the target width, resizing now", "width", targetWidth, "height", targetHeight, "ratio", ratio)
		resized, err := newImage.Resize(*targetWidth, targetHeight)
		if err != nil {
			slog.Error("Failed to resize the image", "path", targetPath, "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result = resized
	}

	slog.Debug("Determine the image type of the resized image", "path", targetPath)
	imgType := bimg.DetermineImageType(result)
	if targetFileType != nil && imgType != *targetFileType {
		img := bimg.NewImage(result)
		result, err = img.Convert(*targetFileType)
		if err != nil {
			slog.Error("Failed to convert image, serve the original type", "targetFileType", bimg.ImageTypeName(*targetFileType), "error", err)
		} else {
			imgType = *targetFileType
		}
	}
	switch imgType {
	case bimg.JPEG:
		w.Header().Set("Content-Type", "image/jpeg")
		break
	case bimg.WEBP:
		w.Header().Set("Content-Type", "image/webp")
		break
	case bimg.PNG:
		w.Header().Set("Content-Type", "image/png")
		break
	case bimg.GIF:
		w.Header().Set("Content-Type", "image/gif")
		break
	default:
		http.Error(w, "Unsupported image type", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(result)))
	_, _ = w.Write(result)
}

func jinyaContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	if !database.FilePathExists("/jinya-content/" + hash) {
		http.NotFound(w, r)
	} else {
		http.ServeFile(w, r, path.Join(storage.GetStoragePath(), hash))
	}
}

func imagePhp(w http.ResponseWriter, r *http.Request) {
	width := r.URL.Query().Get("width")
	imageType := r.URL.Query().Get("type")
	queryId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(queryId)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var targetWidth *int
	parsedWidth, err := strconv.Atoi(width)
	if err != nil {
		targetWidth = nil
	} else {
		targetWidth = &parsedWidth
	}

	imageType = strings.ToLower(imageType)
	var targetType *bimg.ImageType
	if slices.Contains([]string{"webp", "png", "jpeg", "jpg"}, imageType) {
		switch imageType {
		case "webp":
			targetType = new(bimg.WEBP)
			break
		case "png":
			targetType = new(bimg.PNG)
			break
		case "jpg":
		case "jpeg":
			targetType = new(bimg.JPEG)
			break
		}
	}

	file, err := database.FindFileById(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	_, targetPath := path.Split(file.Path)
	targetPath = path.Join(storage.GetStoragePath(), targetPath)

	if strings.HasPrefix(file.Type, "image/") {
		serveConvertedImage(targetPath, targetWidth, targetType, w)
	} else {
		http.ServeFile(w, r, targetPath)
	}
}

func themeAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	themeName := vars["themeName"]
	assetName := vars["assetName"]

	asset, err := database.FindThemeAssetByThemeAndName(themeName, assetName)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, asset.Path)
}

func themeCache(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	theme := vars["themeName"]
	file := vars["file"]
	if file == "cache-status.yaml" {
		http.NotFound(w, r)
		return
	}

	fullPath := path.Join(os.Getenv("CACHE_PATH"), theme, file)
	http.ServeFile(w, r, fullPath)
}
