package imageutil

import "strings"

type ImageEncoding int

const (
	ENCODING_PNG ImageEncoding = iota
	ENCODING_JPEG
	ENCODING_GIF
	ENCODING_UNKNOWN
)

func EncodingExtension(encoding ImageEncoding) string {
	switch encoding {
	case ENCODING_PNG:
		return ".png"
	case ENCODING_JPEG:
		return ".jpg"
	case ENCODING_GIF:
		return ".gif"
	default:
		return ".unknown"
	}
}

func ExtensionEncoding(ext string) ImageEncoding {
	ext = strings.ToLower(ext)

	switch ext {
	case ".png":
		return ENCODING_PNG
	case ".jpg":
		return ENCODING_JPEG
	case ".gif":
		return ENCODING_GIF
	}

	return ENCODING_UNKNOWN
}

var EncodingList = []ImageEncoding{
	ENCODING_PNG,
	ENCODING_JPEG,
	ENCODING_GIF,
	ENCODING_UNKNOWN,
}
