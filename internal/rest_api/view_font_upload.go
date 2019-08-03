package rest_api

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/quard/asciiwrite/internal/storage"
	"github.com/quard/asciiwrite/pkg/figfont"
	"github.com/thedevsaddam/govalidator"
)

var ErrUnableToParseFont = errors.New("unable to process font")

type fontUploadRequest struct {
	Name string `json:"name"`
	Font string `json:"font"`
}

func (srv RestAPIServer) FontUpload(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	var requestData fontUploadRequest
	rules := govalidator.MapData{
		"name": []string{"required", "alpha_space", "between:2,20", "font_not_exists"},
		"font": []string{"required"},
	}
	opts := govalidator.Options{
		Request: request,
		Rules:   rules,
		Data:    &requestData,
	}
	validator := govalidator.New(opts)
	validationError := validator.ValidateJSON()
	if len(validationError) > 0 {
		responseValidationErrors(response, validationError)
	} else {
		err := addNewFont(srv.storage, requestData.Name, requestData.Font)
		if err != nil {
			log.Printf("unable to upload new font: %v", err)
			responseBadRequest(response, request, err)
		} else {
			response.WriteHeader(http.StatusCreated)
		}
	}
}

func addNewFont(storage storage.FontStorage, fontName, fontData string) error {
	fontLoader, errLoader := figfont.NewFileLoader(strings.NewReader(fontData))
	if errLoader != nil {
		log.Printf("unable to create file font loader: %v", errLoader)
		return ErrUnableToParseFont
	}
	font, errParse := fontLoader.Parse()
	if errParse != nil {
		log.Printf("unable to parse font with file loader: %v", errParse)
		return ErrUnableToParseFont
	}
	font.Name = fontName

	err := storage.Add(font)

	return err
}
