package goKLC

import "net/http"

type Controller struct{}
type ControllerFunc func(*http.Request, RouteParams) string
