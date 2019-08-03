package figfont

import (
	"bufio"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

const correctFontSignature = "flf2a$ 6 5 20 15 3 0 143 229"

func TestFileLoaderCheckSignature(t *testing.T) {
	t.Run("correct signature", func(t *testing.T) {
		header := strings.NewReader(correctFontSignature)
		loader, err := NewFileLoader(header)
		assertNoError(t, err)

		err = loader.parseHeaders()
		assertNoError(t, err)
	})

	badSignatures := []string{
		"fIf2a$ 6 5 20 15 3 0 143 229",
		"flf3a$ 6 5 20 15 3 0 143 229",
		"fdjnsk$ 6 5 20 15 3 0 143 229",
	}
	for _, badSignature := range badSignatures {
		caseName := fmt.Sprintf("bad signature %s", badSignature)
		t.Run(caseName, func(t *testing.T) {
			header := strings.NewReader(badSignature)
			loader, err := NewFileLoader(header)
			assertNoError(t, err)

			err = loader.parseHeaders()
			assertError(t, err, "bad font signature")
		})
	}
}

func TestFileLoaderHeaderParse(t *testing.T) {
	testCases := []struct {
		header         string
		hardblank      string
		height         int
		baseline       int
		commentLines   int
		printDirection int
	}{
		{
			correctFontSignature,
			"$",
			6,
			5,
			3,
			0,
		},
		{
			"flf2a% 18 128 0 0 512 -1",
			"%",
			18,
			128,
			512,
			-1,
		},
		{
			"flf2a~ 312 5 54 32 5 1",
			"~",
			312,
			5,
			5,
			1,
		},
		{
			"flf2a1 8 6 7 3 3", // default print direction
			"1",
			8,
			6,
			3,
			0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.header, func(t *testing.T) {
			reader := strings.NewReader(testCase.header)
			loader, err := NewFileLoader(reader)
			assertNoError(t, err)

			err = loader.parseHeaders()
			assertNoError(t, err)
			assertStringEqual(t, testCase.hardblank, loader.font.Hardblank)
			assertIntEqual(t, testCase.height, loader.font.Height)
			assertIntEqual(t, testCase.baseline, loader.font.Baseline)
			assertIntEqual(
				t,
				testCase.printDirection,
				loader.font.PrintDirection,
			)
		})
	}

	errorTestCases := []string{
		"flf2a@ A 1 2 3 4 5",
		"flf2a@ 1 A 2 3 4 5",
		"flf2a@ 1 2 3 4 A 5",
		"flf2a@ 1 2 3 4 5 A",
	}
	for _, errorTestCase := range errorTestCases {
		t.Run(errorTestCase, func(t *testing.T) {
			reader := strings.NewReader(errorTestCase)
			loader, err := NewFileLoader(reader)
			assertNoError(t, err)

			err = loader.parseHeaders()
			assertError(t, err, "strconv.Atoi: parsing \"A\": invalid syntax")
		})
	}
}

func TestFileLoaderLetterParse(t *testing.T) {
	testCases := []struct {
		commentLines int
		height       int
		raw          string
		letters      map[int][]string
	}{
		{
			0,
			5,
			" #  @\n# # @\n### @\n# # @\n# # @@",
			map[int][]string{
				32: {" #  ", "# # ", "### ", "# # ", "# # "},
			},
		},
		{
			3,
			5,
			"c1\nc2\nc3\n _______ @\n|   _   |@\n|       |@\n|___|___|@\n         @@",
			map[int][]string{
				32: {" _______ ", "|   _   |", "|       |", "|___|___|", "         "},
			},
		},
		{
			0,
			5,
			" _______ @\n|   _   |@\n|       |@\n|___|___|@\n         @@\n ______ @\n|      |@\n|   ---|@\n|______|@\n        @@",
			map[int][]string{
				32: {" _______ ", "|   _   |", "|       |", "|___|___|", "         "},
				33: {" ______ ", "|      |", "|   ---|", "|______|", "        "},
			},
		},
		{
			0,
			5,
			" _______ @\n|   _   |@\n|       |@\n|___|___|@\n         @@\n ______ @\n|      |@\n|   ---|@\n|______|@\n        @@",
			map[int][]string{
				32: {" _______ ", "|   _   |", "|       |", "|___|___|", "         "},
				33: {" ______ ", "|      |", "|   ---|", "|______|", "        "},
			},
		},
		{
			0,
			4,
			"(  __)@\n ) _) @\n(____)@@\n0x0422  CYRILLIC CAPITAL LETTER TE\n ____ @\n(_  _)@\n  )(  @\n (__) @@",
			map[int][]string{
				32:   {"(  __)", " ) _) ", "(____)"},
				1058: {" ____ ", "(_  _)", "  )(  ", " (__) "},
			},
		},
	}

	for idx, testCase := range testCases {
		t.Run(fmt.Sprintf("test case: %d", idx), func(t *testing.T) {
			reader := strings.NewReader(testCase.raw)
			loader := FileLoader{FIGFont{}, testCase.commentLines, bufio.NewScanner(reader)}
			loader.font.Height = testCase.height
			loader.font.Letters = make(map[int][]string)
			err := loader.parseLetters()
			assertNoError(t, err)

			if !reflect.DeepEqual(loader.font.Letters, testCase.letters) {
				t.Errorf("letter parsing error:\nexpect: %v\n   got: %v", testCase.letters, loader.font.Letters)
			}
		})
	}
}

func TestParseCharCode(t *testing.T) {
	testCases := []struct {
		line  string
		value int
	}{
		{"0x46 comment", 70},
		{"0xE6", 230},
		{"032", 26},
		{"13", 13},
	}

	for _, testCase := range testCases {
		t.Run(testCase.line, func(t *testing.T) {
			charCode, err := parseCharCode(testCase.line)
			assertNoError(t, err)
			assertIntEqual(t, testCase.value, charCode)
		})
	}
}

func TestFileLoader(t *testing.T) {
	fontFile := `flf2a$ 6 4 6 -1 4
3x5 font by Richard Kirk (rak@crosfield.co.uk).
Ported to figlet, and slightly changed (without permission :-})
by Daniel Cabeza Gras (bardo@dia.fi.upm.es)

	@
	@
	@
	@
	@
	@@
	@
 #  @
 #  @
 #  @
	@
 #  @@
	@
# # @
# # @
	@
	@
	@@
	@
# # @
### @
# # @
### @
# # @@`

	loader, err := NewFileLoader(strings.NewReader(fontFile))
	assertNoError(t, err)
	font, err := loader.Parse()
	assertNoError(t, err)

	assertStringEqual(t, "$", font.Hardblank)
	assertIntEqual(t, 6, font.Height)
	assertIntEqual(t, 4, font.Baseline)
	assertIntEqual(t, 4, loader.commentLines)

	letters := map[int][]string{
		32: {"	", "	", "	", "	", "	", "	"},
		33: {"	", " #  ", " #  ", " #  ", "	", " #  "},
		34: {"	", "# # ", "# # ", "	", "	", "	"},
		35: {"	", "# # ", "### ", "# # ", "### ", "# # "},
	}

	for charCode, letter := range letters {
		if !reflect.DeepEqual(letter, font.Letters[charCode]) {
			t.Errorf(
				"letter not match:\nexcept: %v\n   got: %v",
				letter,
				loader.font.Letters[charCode],
			)
		}
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func assertError(t *testing.T, err error, msg string) {
	t.Helper()

	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != msg {
		t.Errorf("expect error with message '%s'. got '%s'", msg, err)
	}

}

func assertStringEqual(t *testing.T, want, got string) {
	t.Helper()

	if want != got {
		t.Errorf("want '%s' but got '%s'", want, got)
	}
}

func assertIntEqual(t *testing.T, want, got int) {
	t.Helper()

	if want != got {
		t.Errorf("want '%d' but got '%d'", want, got)
	}
}
