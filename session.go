package goKLC

type SessionInterface interface {
	Set(key string, value interface{})
	Get(key string, defaultValue interface{}) interface{}
	Delete()
}

type Session struct {
	Id     string
	Key    string
	Value  interface{}
	cookie string
}

func NewSession(r *Request) *Session {
	cookie := NewCookie()
	cookieName := _config.Get("SessionName", "goKLCSession").(string)
	cookie.Get(r, cookieName)

	return &Session{Id: _app.GetSessionKey(), cookie: cookie.Value}
}

func (s *Session) Set(key string, value interface{}) {
	s.Key = key
	s.Value = value

	collection := _sessionCollector.GetCollection(s.cookie)
	collection.Set(s)
}

func (s *Session) Get(key string, defaultValue interface{}) interface{} {
	collection := _sessionCollector.GetCollection(s.cookie)
	session := collection.Get(key)

	if session == nil || session.Value == nil {

		return defaultValue
	}

	return session.Value
}

func (s *Session) Delete() {
	collection := _sessionCollector.GetCollection(s.cookie)
	collection.Delete(s.Key)
}
