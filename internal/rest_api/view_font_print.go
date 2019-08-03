package rest_api

import (
	"errors"
	"log"
	"net/http"

	"github.com/quard/asciiwrite/internal/storage"
	"github.com/thedevsaddam/govalidator"
)

var ErrUnableToPrint = errors.New("unable to print phrase")

type printRequest struct {
	Name   string `json:"name"`
	Phrase string `json:"phrase"`
}

func (srv RestAPIServer) Print(response http.ResponseWriter, request *http.Request) {
	var requestData printRequest

	rules := govalidator.MapData{
		"name":   []string{"required", "alpha_space", "between:2,20"},
		"phrase": []string{"required"},
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
		text, err := getPrintedPhrase(srv.storage, requestData.Name, requestData.Phrase)
		if err != nil {
			responseBadRequest(response, request, err)
		} else {
			response.Write([]byte(text))
		}
	}
}

func getPrintedPhrase(stor storage.FontStorage, fontName, phrase string) (string, error) {
	font, err := stor.Get(fontName)
	if err == storage.ErrFontNotFound {
		return "", err
	} else if err != nil {
		log.Printf("unable to retrieve font: %v", err)
		return "", errors.New("unable to retrieve font")
	}

	return font.Print(phrase)
}
