package goKLC

type Cookie struct {
	request  *Request
	Name     string
	Value    string
	Duration int
	Path     string
}

func NewCookie() *Cookie {

	return &Cookie{Path: "/"}
}

func (c *Cookie) Set(r *Response) {
	r.AddCookie(c)
}

func (c *Cookie) Get(r *Request, name string) (bool, *Cookie) {
	c.request = r
	cookie, err := r.Request.Cookie(name)

	if err != nil {

		return false, nil
	}

	c.Name = cookie.Name
	c.Value = cookie.Value
	c.Duration = cookie.MaxAge

	return true, c
}
