package figfont

import (
	"fmt"
	"strings"
)

// FIGFont contains all data of FIG Font to able to print message
type FIGFont struct {
	Name           string           `json:"name"`
	Hardblank      string           `json:"hardblank"`
	Height         int              `json:"height"`
	Baseline       int              `json:"baseline"`
	PrintDirection int              `json:"printDirection"`
	Letters        map[int][]string `json:"letters"`
}

func (font FIGFont) Print(phrase string) (string, error) {
	var printed []string

	if font.PrintDirection != 0 {
		phrase = strReverse(phrase)
	}

	for row := 0; row < font.Height; row++ {
		var printedRow string
		for _, letter := range phrase {
			data, ok := font.Letters[int(letter)]
			if !ok {
				return "", fmt.Errorf("unknown letter '%v'", letter)
			}

			printedRow = printedRow + strings.Replace(data[row], font.Hardblank, " ", -1)
		}
		printed = append(printed, printedRow)
	}

	return strings.Join(printed, "\n"), nil
}

func strReverse(str string) string {
	reversed := []rune(str)
	for left, right := 0, len(reversed)-1; left < right; left, right = left+1, right-1 {
		reversed[left], reversed[right] = reversed[right], reversed[left]
	}
	return string(reversed)
}
