package goKLC

type Response interface {
	GetHeaders() map[string]string
	HasHeaders(name string) bool
	GetHeader(name string) string
	WithHeader(name string, value string)
	WithBody(value string)
	GetBody() string
	GetStatusCode() int
	WithStatusCode(status int)
	WithCookie(cookie interface{})
	GetCookies() map[string]interface{}
	GetCookie(name string) interface{}
}
