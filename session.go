package goKLC

import "time"

type SessionInterface interface {
	Set(key string, value interface{})
	Get(key string, defaultValue interface{}) interface{}
	Delete()
}

type Session struct {
	Id          string
	Key         string
	Value       interface{}
	cookie      string
	maxDuration time.Duration
	Duration    time.Duration
}

func NewSession(r *Request) *Session {
	cookie := NewCookie()
	cookieName := _config.Get("SessionName", "goKLCSession").(string)
	cookie.Get(r, cookieName)

	session := &Session{
		Id:          _app.GetSessionKey(),
		cookie:      cookie.Value,
		maxDuration: time.Second * time.Duration(cookie.Duration),
	}

	return session
}

func (s *Session) Set() {
	collection := _sessionCollector.GetCollection(s.cookie)
	collection.Set(s)

	if s.Duration > 0 {
		if s.Duration < s.maxDuration && s.Duration != 0 {
			time.AfterFunc(s.Duration, s.Delete)
		} else {
			time.AfterFunc(s.maxDuration, s.Delete)
		}
	}
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
