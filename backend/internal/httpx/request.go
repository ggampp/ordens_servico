package httpx

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// DecodeAndValidate decodes a JSON body into dst and runs struct validation.
func DecodeAndValidate(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return NewBadRequest("invalid JSON body: " + err.Error())
	}
	if err := validate.Struct(dst); err != nil {
		var verrs validator.ValidationErrors
		if ok := asValidation(err, &verrs); ok {
			msgs := make([]string, 0, len(verrs))
			for _, fe := range verrs {
				msgs = append(msgs, fe.Field()+" failed '"+fe.Tag()+"'")
			}
			return NewBadRequest("validation error: " + strings.Join(msgs, ", "))
		}
		return NewBadRequest("validation error")
	}
	return nil
}

func asValidation(err error, target *validator.ValidationErrors) bool {
	if v, ok := err.(validator.ValidationErrors); ok {
		*target = v
		return true
	}
	return false
}

// QueryInt reads an integer query parameter with a fallback default.
func QueryInt(r *http.Request, key string, def int) int {
	if v := r.URL.Query().Get(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

// QueryFloatPtr reads an optional float query parameter.
func QueryFloatPtr(r *http.Request, key string) *float64 {
	if v := r.URL.Query().Get(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return &f
		}
	}
	return nil
}
