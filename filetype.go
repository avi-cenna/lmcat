package main

import (
	mapset "github.com/deckarep/golang-set/v2"
	"path"
	"strings"
)

func IsLikelyNonTextFile(filepath string) bool {
	ext := path.Ext(filepath)
	ext, _ = strings.CutPrefix(ext, ".")
	return nonTextFileExtensions.Contains(ext)
}

var nonTextFileExtensions mapset.Set[string]

func init() {
	nonTextFileExtensions = mapset.NewSet(
		// Applications
		"wasm",
		"dex",
		"dey",

		// Archives
		"epub",
		"zip",
		"tar",
		"rar",
		"gz",
		"bz2",
		"7z",
		"xz",
		"zst",
		"pdf",
		"exe",
		"swf",
		"rtf",
		"eot",
		"ps",
		"sqlite",
		"nes",
		"crx",
		"cab",
		"deb",
		"ar",
		"Z",
		"lz",
		"rpm",
		"elf",
		"dcm",
		"iso",
		"macho",

		// Audio
		"mid",
		"mp3",
		"m4a",
		"ogg",
		"flac",
		"wav",
		"amr",
		"aac",
		"aiff",

		// Documents
		"doc",
		"docx",
		"xls",
		"xlsx",
		"ppt",
		"pptx",

		// Fonts
		"woff",
		"woff2",
		"ttf",
		"otf",

		// Images
		"jpg",
		"jp2",
		"png",
		"gif",
		"webp",
		"cr2",
		"tif",
		"bmp",
		"jxr",
		"psd",
		"ico",
		"heif",
		"dwg",

		// Videos
		"mp4",
		"m4v",
		"mkv",
		"webm",
		"mov",
		"avi",
		"wmv",
		"mpg",
		"flv",
		"3gp",
	)
}
