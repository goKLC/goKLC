package goKLC

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
)

type Request struct {
	Request *http.Request
	form    map[string][]string
}

func NewRequest(req *http.Request) *Request {
	r := &Request{}
	r.Request = req

	ct := req.Header.Get("Content-Type")
	ct, _, _ = mime.ParseMediaType(ct)

	switch ct {
	case "multipart/form-data":
		//todo memory config
		req.ParseMultipartForm(1024 * 1024 * 5)
		r.form = req.MultipartForm.Value
		break
	default:
		req.ParseForm()
		r.form = req.Form
		break
	}

	return r
}

func (r *Request) Input(key string) []string {

	ct := r.Request.Header.Get("Content-Type")
	ct, _, _ = mime.ParseMediaType(ct)

	return r.form[key]
}

func (r *Request) Json(dataModel *interface{}) {
	body, err := ioutil.ReadAll(r.Request.Body)

	if err != nil {
		fmt.Println(err.Error())
	}

	err = json.Unmarshal(body, &dataModel)
}
