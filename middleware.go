package goKLC

type Middleware struct {
	Handle    func(request *Request)
	Terminate func()
}
