package goKLC

import (
	"encoding/json"
	"io/ioutil"
	"mime"
	"net/http"
)

type Request struct {
	routeParams RouteParams
	Request     *http.Request
	form        map[string][]string
}

func NewRequest(req *http.Request, routeParams RouteParams) *Request {
	r := &Request{}
	r.Request = req
	r.routeParams = routeParams

	ct := req.Header.Get("Content-Type")
	ct, _, _ = mime.ParseMediaType(ct)

	switch ct {
	case "multipart/form-data":
		req.ParseMultipartForm(_config.Get("MaxFormMemory", 1024*1024*5).(int64))
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
		_app.Log().Error(err.Error(), nil)
	}

	err = json.Unmarshal(body, &dataModel)
}

func (r *Request) GetParameter(key string) interface{} {

	return r.routeParams[key]
}

func (r *Request) GetCookie(name string) Cookie {
	cookieRaw, err := r.Request.Cookie(name)

	if err != nil {

		return nil
	}

	cookie := _app.cookie
	cookie.Create(cookieRaw.Name, cookieRaw.Value, cookieRaw.MaxAge, cookieRaw.Path)

	return cookie
}
