package goKLC

type Cookie struct {
	request  *Request
	Name     string
	Value    string
	Duration int
}

func NewCookie() *Cookie {

	return &Cookie{}
}

func (c *Cookie) Set(r *Response) {
	r.AddCookie(c)
}

func (c *Cookie) Get(r *Request, name string) bool {
	c.request = r
	cookie, err := r.Request.Cookie(name)

	if err != nil {

		return false
	}

	c.Name = cookie.Name
	c.Value = cookie.Value
	c.Duration = cookie.MaxAge

	return true
}
