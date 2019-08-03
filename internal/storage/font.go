package storage

import (
	"errors"

	"github.com/quard/asciiwrite/pkg/figfont"
)

type FontStorage interface {
	Add(font figfont.FIGFont) error
	Get(name string) (figfont.FIGFont, error)
	IsExist(name string) (bool, error)
	Names() ([]string, error)
}

var ErrFontNotFound = errors.New("font not found")
