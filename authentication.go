package goKLC

type Auth struct {
	Model ModelInterface
}

func (a *Auth) LoginById(id int, request *Request) {
	user := a.getUser(id)

	session := NewSession(request)
	session.Key = "auth"
	session.Value = user
	session.Set()
}

func (a *Auth) getUser(id int) ModelInterface {
	model := a.Model
	_app.DB().Where("id = ?", id).First(model)

	return model
}

func (a *Auth) Check(request *Request) bool {
	user := a.User(request)

	if user == nil {
		return false
	}

	return true
}

func (a *Auth) User(request *Request) ModelInterface {
	session := NewSession(request)
	user := session.Get("auth", nil)

	return user
}
