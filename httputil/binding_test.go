package httputil

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindForm(t *testing.T) {
	// 示例用法
	type MyStruct struct {
		Name   string `form:"name"`
		Age    int    `form:"age"`
		Active bool   `form:"active"`
	}

	req := &http.Request{
		Form: url.Values{
			"name":   {"John Doe"},
			"age":    {"30"},
			"active": {"true"},
		},
	}

	var myStruct MyStruct
	err := bindForm(req, &myStruct)
	assert.NoError(t, err)

	assert.Equal(t, "John Doe", myStruct.Name)
	assert.Equal(t, 30, myStruct.Age)
	assert.Equal(t, true, myStruct.Active)
}

func TestBindJSON(t *testing.T) {
	// 示例用法
	type MyStruct struct {
		Name   string `form:"name" json:"name,omitempty"`
		Age    int    `form:"age" json:"age,omitempty"`
		Active bool   `form:"active" json:"active,omitempty"`
	}

	req := &http.Request{
		Body: ioutil.NopCloser(strings.NewReader(`{"name":"John Doe","age":30,"active":true}`)),
	}

	var myStruct MyStruct
	err := bindJSON(req, &myStruct)
	assert.NoError(t, err)

	assert.Equal(t, "John Doe", myStruct.Name)
	assert.Equal(t, 30, myStruct.Age)
	assert.Equal(t, true, myStruct.Active)
}

func TestBind(t *testing.T) {
	// 示例用法
	type MyStruct struct {
		Name   string `form:"name" json:"name,omitempty"`
		Age    int    `form:"age" json:"age,omitempty"`
		Active bool   `form:"active" json:"active,omitempty"`
	}

	req := &http.Request{
		Body: ioutil.NopCloser(strings.NewReader(`{"name":"John Doe","age":30,"active":true}`)),
	}

	var myStruct MyStruct
	err := Bind(req, &myStruct)
	assert.NoError(t, err)

	assert.Equal(t, "John Doe", myStruct.Name)
	assert.Equal(t, 30, myStruct.Age)
	assert.Equal(t, true, myStruct.Active)
}
