package goKLC

type Controller struct{}
type ControllerFunc func(*Request, RouteParams) *Response
