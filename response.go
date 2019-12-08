package goKLC

import "net/http"

type Response struct {
	status  int
	content string
	cookies []*Cookie
}

func NewResponse() *Response {

	return &Response{}
}

func (r *Response) Ok(content string) *Response {
	r.status = http.StatusOK
	r.content = content

	return r
}

func (r *Response) Error(content string) *Response {
	r.status = http.StatusBadRequest
	r.content = content

	return r
}

func (r *Response) AddCookie(c *Cookie) *Response {
	r.cookies = append(r.cookies, c)

	return r
}
