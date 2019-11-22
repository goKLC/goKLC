package goKLC

type Middleware struct {
	Handle    func()
	Terminate func()
}
