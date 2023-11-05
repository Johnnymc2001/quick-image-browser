package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"golang.design/x/clipboard"
	// "github.com/srwiley/rasterx"
)

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func GetBase64FromFile(path string) string {
	// Read the entire file into a byte slice
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var base64Encoding string
	mimeType := http.DetectContentType(bytes)

	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	base64Encoding += toBase64(bytes)

	return base64Encoding
}

// Greet returns a greeting for the given name

func GetCurrentAppDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path := filepath.Dir(ex)
	return path
}

func GetImages(directory string) []ImageObj {

	// imagePath := a.GetCurrentAppDir() + "/images"
	imagePath := directory

	files, err := os.ReadDir(imagePath)
	if err != nil {
		var imageList []ImageObj
		return imageList
	} else {

		var imageList []ImageObj
		allowList := []string{"png", "jpg", "jpeg"}
		for _, file := range files {
			if !file.IsDir() {
				splits := strings.Split(file.Name(), ".")
				ext := splits[len(splits)-1]

				if slices.Contains(allowList, ext) {

					base64 := GetBase64FromFile(imagePath + "/" + file.Name())
					newImageObj := ImageObj{Path: imagePath + "/" + file.Name(), Name: file.Name(), Base64: base64}
					imageList = append(imageList, newImageObj)
				}
			}
		}

		return imageList
	}
}

func CopyImageToClipboard(filePath string) {
	f, err := os.Open(filePath)
	// splits := strings.Split(filePath, ".")
	// ext := splits[len(splits)-1]

	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		panic(err)
	}

	clipboard.Write(clipboard.FmtImage, buf.Bytes())

}
