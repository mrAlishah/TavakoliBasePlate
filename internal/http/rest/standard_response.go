package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
)

var statues = [...]string{
	"Ok",
	"Fail",
}

const (
	OK ResponseStatus = iota
	Fail
)

type ResponseStatus uint8

func (rs ResponseStatus) String() string {
	return statues[rs]
}

func (rs ResponseStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(rs.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (rs *ResponseStatus) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*rs = 0
	for i, v := range statues {
		if v == s {
			*rs = ResponseStatus(i)
			break
		}
	}

	return nil
}

var (
	InternalServerError = GetFailResponse("Internal Server Error", nil)
	NotImplemented      = GetFailResponse("Not Implemented", nil)
	NotFound            = GetFailResponse("Not Found!", nil)
	NotAcceptable       = GetFailResponse("Not Acceptable!", nil)
	EmptyQuery          = GetFailResponse("Query should not be empty string", nil)
)

type StandardResponse struct {
	Status  ResponseStatus    `json:"status"`
	Message string            `json:"message"`
	Data    map[string]string `json:"data,omitempty"`
}

func GetFailResponse(msg string, data map[string]string) StandardResponse {
	return StandardResponse{
		Status:  Fail,
		Message: msg,
		Data:    data,
	}
}

func GetSuccessResponse(msg string) StandardResponse {
	return StandardResponse{
		Status:  OK,
		Message: msg,
	}
}

func GetFailResponseFromValidationErrors(errs validator.ValidationErrors) StandardResponse {
	return StandardResponse{
		Status:  Fail,
		Message: "Params Validation Error",
		Data:    SimplifyError(errs),
	}
}

func GetFailResponseFromErrors(errors []error) StandardResponse {
	buf := bytes.NewBufferString("")
	for _, err := range errors {
		buf.WriteString(err.Error())
		buf.WriteString("\n")
	}

	return GetFailResponse(buf.String(), nil)
}

func SimplifyError(verr validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)

	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs[f.Field()] = err
	}

	return errs
}
