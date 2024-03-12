package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

// BindResp parse the response body as JSON and store it in v.
func BindResp[V any](res *http.Response, v V) error {
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(v); err != nil {
		return err
	}

	return nil
}

// Bind parse the request body as JSON or form data and store it in v.
func Bind[V any](req *http.Request, v V) error {
	contentType := req.Header.Get("Content-Type")
	if contentType == "application/x-www-form-urlencoded" {
		return bindForm(req, v)
	}

	return bindJSON(req, v)
}

func bindJSON[V any](req *http.Request, v V) error {
	if err := json.NewDecoder(req.Body).Decode(v); err != nil {
		return err
	}

	return nil
}

func bindForm[V any](req *http.Request, v V) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	val := reflect.ValueOf(v)

	// make sure v is a pointer
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("Bind: v must be a pointer")
	}

	// get the value that v points to
	elem := val.Elem()

	// make sure v points to a struct
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("Bind: v must point to a struct")
	}

	// iterate over the fields of the struct
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := field.Type()
		fieldName := elem.Type().Field(i).Tag.Get("form")
		if fieldName == "" {
			fieldName = elem.Type().Field(i).Name
		}

		// check if the field can be set
		if !field.CanSet() {
			continue
		}

		fieldValue := req.FormValue(fieldName)

		// set the value of the field
		switch fieldType.Kind() {
		case reflect.String:
			field.SetString(fieldValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if num, err := strconv.ParseInt(fieldValue, 10, 64); err == nil {
				field.SetInt(num)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if num, err := strconv.ParseUint(fieldValue, 10, 64); err == nil {
				field.SetUint(num)
			}
		case reflect.Bool:
			if boolValue, err := strconv.ParseBool(fieldValue); err == nil {
				field.SetBool(boolValue)
			}
		default:
			// ignore other types
			continue
		}
	}

	return nil
}
