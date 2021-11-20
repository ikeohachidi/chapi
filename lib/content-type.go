package lib

import (
	"mime"
	"path/filepath"
)

func DetectContentType(fileName string) string {
	switch filepath.Ext(fileName) {
	case ".js":
		return "application/javascript"
	case ".css":
		return "text/css"
	case ".svg":
		return "image/svg+xml"
	case ".eot":
		return "font/eot"
	case ".ttf":
		return "font/ttf"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	default:
		return mime.TypeByExtension(fileName)
	}
}
