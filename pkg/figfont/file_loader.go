package figfont

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

//           flf2a$ 6 5 20 15 3 0 143 229    NOTE: The first five characters in
//             |  | | | |  |  | |  |   |     the entire file must be "flf2a".
//            /  /  | | |  |  | |  |   \
//   Signature  /  /  | |  |  | |   \   Codetag_Count
//     Hardblank  /  /  |  |  |  \   Full_Layout*
//          Height  /   |  |   \  Print_Direction
//          Baseline   /    \   Comment_Lines
//           Max_Length      Old_Layout*

const fontSignature = "flf2"
const endmarkChars = "@#%$"

// FileLoader is set of data to load font from file
type FileLoader struct {
	font         FIGFont
	commentLines int
	buf          *bufio.Scanner
}

// NewFileLoader creates FIG font loader from file
func NewFileLoader(r io.Reader) (FileLoader, error) {
	loader := FileLoader{font: FIGFont{}, buf: bufio.NewScanner(r)}
	loader.font.Letters = make(map[int][]string)

	return loader, nil
}

// Parse read font data and fill all necessary attributes to use font in future
func (loader *FileLoader) Parse() (FIGFont, error) {
	err := loader.parseHeaders()
	if err != nil {
		return loader.font, err
	}

	err = loader.parseLetters()
	if err != nil {
		return loader.font, err
	}

	return loader.font, nil
}

func (loader *FileLoader) parseHeaders() error {
	loader.buf.Scan()
	header := loader.buf.Text()
	if header[:len(fontSignature)] != fontSignature {
		return errors.New("bad font signature")
	}

	fields := strings.Fields(header)
	loader.font.Hardblank = fields[0][len(fields[0])-1:]

	var err error
	loader.font.Height, err = strconv.Atoi(fields[1])
	if err != nil {
		return err
	}
	loader.font.Baseline, err = strconv.Atoi(fields[2])
	if err != nil {
		return err
	}
	loader.commentLines, err = strconv.Atoi(fields[5])
	if err != nil {
		return err
	}
	if len(fields) > 6 {
		loader.font.PrintDirection, err = strconv.Atoi(fields[6])
		if err != nil {
			return err
		}
	}

	return nil
}

func (loader *FileLoader) parseLetters() error {
	// skip comment lines
	for i := 0; i < loader.commentLines; i++ {
		loader.buf.Scan()
	}

	var lineEndChar string
	var letter []string
	charCode := 32

	for loader.buf.Scan() {
		line := loader.buf.Text()

		if len(line) == 0 {
			continue
		}

		if strings.ContainsAny(line[len(line)-1:], endmarkChars) {
			// Letter
			if len(lineEndChar) == 0 {
				// start of letter
				lineEndChar = line[len(line)-1:]
			}
			if (len(line) > 1 && line[len(line)-2:] == strings.Repeat(lineEndChar, 2)) || len(letter) >= loader.font.Height-1 {
				endLength := 2
				if loader.font.Height == 1 {
					endLength = 1
				}
				letter = append(letter, line[:len(line)-endLength])
				loader.font.Letters[charCode] = letter
				letter = []string{}
				lineEndChar = ""
				charCode++
			} else {
				letter = append(letter, line[:len(line)-1])
			}
		} else {
			// extended char code
			var err error
			charCode, err = parseCharCode(line)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func parseCharCode(line string) (int, error) {
	substrings := strings.SplitN(line, " ", 2)

	base := 10
	charCodeStr := substrings[0]
	if charCodeStr[:2] == "0x" {
		base = 16
		charCodeStr = charCodeStr[2:]
	} else if charCodeStr[:1] == "0" {
		base = 8
		charCodeStr = charCodeStr[1:]
	}

	charCode, err := strconv.ParseInt(charCodeStr, base, 16)
	if err != nil {
		return 0, err
	}

	return int(charCode), nil
}
