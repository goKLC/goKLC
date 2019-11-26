package goKLC

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

type MiddlewareInterface interface {
	Handle(request *Request) *Response
	Terminate(request *Response)
}

func (m Middleware) Handle(request *Request) *Response {

	return nil
}

func (m Middleware) Terminate(response *Response) {

}
