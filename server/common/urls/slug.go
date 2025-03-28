package urls

import (
	"bytes"
	"net/url"
	"strings"
	"unicode"
)

func NormalizeVietnamese(input string) string {
	input = strings.ToLower(input)
	normalizationMap := map[rune]rune{
		'á': 'a', 'à': 'a', 'ả': 'a', 'ã': 'a', 'ạ': 'a', 'ă': 'a', 'ằ': 'a', 'ắ': 'a', 'ẳ': 'a', 'ẵ': 'a', 'ặ': 'a',
		'â': 'a', 'ầ': 'a', 'ấ': 'a', 'ẩ': 'a', 'ẫ': 'a', 'ậ': 'a',
		'é': 'e', 'è': 'e', 'ẻ': 'e', 'ẽ': 'e', 'ẹ': 'e', 'ê': 'e', 'ề': 'e', 'ế': 'e', 'ể': 'e', 'ễ': 'e', 'ệ': 'e',
		'í': 'i', 'ì': 'i', 'ỉ': 'i', 'ĩ': 'i', 'ị': 'i',
		'ó': 'o', 'ò': 'o', 'ỏ': 'o', 'õ': 'o', 'ọ': 'o', 'ô': 'o', 'ồ': 'o', 'ố': 'o', 'ổ': 'o', 'ỗ': 'o', 'ộ': 'o',
		'ơ': 'o', 'ờ': 'o', 'ớ': 'o', 'ở': 'o', 'ỡ': 'o', 'ợ': 'o',
		'ú': 'u', 'ù': 'u', 'ủ': 'u', 'ũ': 'u', 'ụ': 'u', 'ư': 'u', 'ừ': 'u', 'ứ': 'u', 'ử': 'u', 'ữ': 'u', 'ự': 'u',
		'ý': 'y', 'ỳ': 'y', 'ỷ': 'y', 'ỹ': 'y', 'ỵ': 'y',
		'đ': 'd',
	}
	var result []rune
	for _, r := range input {
		if replacement, exists := normalizationMap[r]; exists {
			result = append(result, replacement)
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func GenerateSlug(title string) string {
	normalized := NormalizeVietnamese(title)
	var buf bytes.Buffer
	for _, r := range normalized {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' {
			buf.WriteRune(r)
		} else {
			buf.WriteRune('-')
		}
	}
	cleanStr := buf.String()
	for strings.Contains(cleanStr, "--") {
		cleanStr = strings.ReplaceAll(cleanStr, "--", "-")
	}
	cleanStr = strings.Trim(cleanStr, "-")
	return url.PathEscape(cleanStr)
}
