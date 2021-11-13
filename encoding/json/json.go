package json

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Decode(i interface{}, r io.Reader) error {
	err := json.NewDecoder(r).Decode(i)
	if err != nil {
		return err
	}

	return validate.Struct(i)
}

func Encode(i interface{}, w io.Writer) error {
	return json.NewEncoder(w).Encode(i)
}
