package utils

import (
	"regexp"
	"strings"
)

func Latinizer(str string) string {

	var m = map[rune]string{
		'а': "a",
		'б': "b",
		'в': "v",
		'г': "g",
		'д': "d",
		'е': "e",
		'ё': "yo",
		'ж': "zh",
		'з': "z",
		'и': "i",
		'й': "j",
		'к': "k",
		'л': "l",
		'м': "m",
		'н': "n",
		'о': "o",
		'п': "p",
		'р': "r",
		'с': "s",
		'т': "t",
		'у': "u",
		'ф': "f",
		'х': "h",
		'ц': "c",
		'ч': "ch",
		'ш': "sh",
		'щ': "sch",
		'ъ': "'",
		'ы': "y",
		'ь': "",
		'э': "e",
		'ю': "ju",
		'я': "ja",
	}

	tr := make([]byte, 0, len(str))
	for _, r := range []rune(strings.ToLower(str)) {
		if v, ok := m[r]; ok {
			tr = append(tr, []byte(v)...)
		} else {
			tr = append(tr, []byte(string(r))...)
		}
	}

	re := regexp.MustCompile(`[^\w ]`)

	res := re.ReplaceAllString(string(tr), "")

	re = regexp.MustCompile(`[ ]`)
	res = re.ReplaceAllString(res, "-")

	return res
}
