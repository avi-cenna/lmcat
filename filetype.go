package main

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"strings"
	"unicode/utf8"
)

var nonTextFileExtensions mapset.Set[string]
var textFileExtensions mapset.Set[string]

func IsLikelyTextFile(filepath string) bool {
	ext := path.Ext(filepath)
	ext, _ = strings.CutPrefix(ext, ".")

	switch {
	case textFileExtensions.Contains(ext):
		log.Debug().Str("filepath", filepath).Msg("Detected text file")
		return true
	case nonTextFileExtensions.Contains(ext):
		log.Debug().Str("filepath", filepath).Msg("Skipping non-text file")
		return false
	default:
		// Fallback to checking the file content
		isBinary, err := isBinaryFile(filepath)
		if err != nil {
			log.Debug().Str("filepath", filepath).Msg("Failed to determine if file is binary")
			return false
		}
		return !isBinary
	}
}

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
		"bmp",
		"cr2",
		"dwg",
		"gif",
		"heif",
		"ico",
		"jp2",
		"jpg",
		"jpeg",
		"jxr",
		"png",
		"psd",
		"svg",
		"tif",
		"webp",

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

	textFileExtensions = mapset.NewSet(
		// Source Code
		"c",
		"cpp",
		"h",
		"hpp",
		"cs",
		"go",
		"java",
		"js",
		"ts",
		"jsx",
		"tsx",
		"py",
		"rb",
		"rs",
		"swift",
		"php",
		"html",
		"htm",
		"css",
		"scss",
		"xml",
		"json",
		"yml",
		"yaml",
		"toml",
		"ini",
		"sh",
		"bat",
		"sql",
		"md",
		"rst",
		"tex",
		"bib",

		// Plain Text
		"txt",
		"log",
		"csv",
		"tsv",

		// Configuration Files
		"conf",
		"cfg",
		"env",
		"properties",

		// Scripts
		"pl",
		"pm",
		"awk",
		"sed",
		"lua",
		"r",

		// Data Formats
		"sgml",
		"srt",
		"vtt",
		"tsv",
	)
}

// isBinaryFile checks if a file is binary based on its first few bytes
func isBinaryFile(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Err(err).Str("path", path).Msg("Failed to close file")
		}
	}(file)

	// Read the first 512 bytes
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && n == 0 {
		return false, err
	}

	// Check if the buffer contains non-UTF-8 characters
	for _, b := range buffer[:n] {
		if b > 0x7F || !utf8.ValidRune(rune(b)) {
			return true, nil // It's binary
		}
	}

	return false, nil // It's likely text
}
