package storage

import (
	"context"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/quard/asciiwrite/pkg/figfont"
)

// FirebaseFontStorage is a font storage realisation with Firebase as database
type FirebaseFontStorage struct {
	db *db.Client
}

type firebaseFIGFont struct {
	Name           string     `json:"name"`
	Hardblank      string     `json:"hardblank"`
	Height         int        `json:"height"`
	Baseline       int        `json:"baseline"`
	PrintDirection int        `json:"printDirection"`
	Letters        [][]string `json:"letters"`
}

// NewFirebaseFontStorage connect to Firebase and return instance of font storage
func NewFirebaseFontStorage() (FirebaseFontStorage, error) {
	var err error

	storage := FirebaseFontStorage{}
	config := &firebase.Config{
		ProjectID:   "asciiwrite",
		DatabaseURL: "https://asciiwrite.firebaseio.com/",
	}
	app, err := firebase.NewApp(context.Background(), config)
	if err != nil {
		log.Fatalf("unable to connect to Firebase: %v", err)
	}

	storage.db, err = app.Database(context.Background())
	if err != nil {
		log.Fatalf("unable to initialize DB: %v", err)
	}

	return storage, nil
}

func (stor FirebaseFontStorage) Add(font figfont.FIGFont) error {
	fontRef := stor.db.NewRef("fonts")
	_, err := fontRef.Push(context.Background(), &font)

	return err
}

func (stor FirebaseFontStorage) Get(name string) (figfont.FIGFont, error) {
	var font figfont.FIGFont

	fontRef := stor.db.NewRef("fonts")
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	fonts, err := fontRef.OrderByChild("name").EqualTo(name).GetOrdered(ctx)
	if err != nil {
		return font, err
	}
	var fFont firebaseFIGFont
	if len(fonts) == 0 {
		return font, ErrFontNotFound
	}
	for _, fontData := range fonts {
		err := fontData.Unmarshal(&fFont)
		if err != nil {
			err := fontData.Unmarshal(&font) // https://stackoverflow.com/a/17777278
			if err != nil {
				return font, err
			}
		} else {
			font = fFont.FIGFont()
		}
	}

	return font, nil
}

func (stor FirebaseFontStorage) IsExist(name string) (bool, error) {
	fontRef := stor.db.NewRef("fonts")
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	fonts, err := fontRef.OrderByChild("name").EqualTo(name).GetOrdered(ctx)
	if err != nil {
		return false, err
	}

	return len(fonts) != 0, nil
}

func (stor FirebaseFontStorage) Names() ([]string, error) {
	var names []string
	fontRef := stor.db.NewRef("fonts")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	fonts, err := fontRef.OrderByChild("name").GetOrdered(ctx)
	if err != nil {
		return []string{}, err
	}
	type Font struct {
		Name string `json:"name"`
	}
	for _, fontData := range fonts {
		var font Font
		if err := fontData.Unmarshal(&font); err != nil {
			log.Printf("unable to unmarshal font '%s': %v", fontData.Key(), err)
		} else {
			names = append(names, font.Name)
		}
	}

	return names, nil
}

func (fFont firebaseFIGFont) FIGFont() (font figfont.FIGFont) {
	font.Name = fFont.Name
	font.Hardblank = fFont.Hardblank
	font.Height = fFont.Height
	font.Baseline = fFont.Baseline
	font.PrintDirection = fFont.PrintDirection
	font.Letters = make(map[int][]string)
	for num, letter := range fFont.Letters {
		if letter != nil {
			font.Letters[num] = letter
		}
	}

	return font
}
