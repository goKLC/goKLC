package goKLC

type sessionCollector map[string]*sessionCollection
type sessionCollection map[string]*Session

func newSessionCollector() sessionCollector {

	return sessionCollector{}
}

func newSessionCollection() *sessionCollection {

	return &sessionCollection{}
}

func (sc sessionCollector) GetCollection(key string) *sessionCollection {
	collection, ok := sc[key]

	if ok {

		return collection
	}

	return sc.SetCollection(key)
}

func (sc sessionCollector) SetCollection(key string) *sessionCollection {
	collection := newSessionCollection()
	(sc)[key] = collection

	return collection
}

func (sc sessionCollection) Get(key string) *Session {
	session, ok := (sc)[key]

	if ok {

		return session
	}

	return nil
}

func (sc sessionCollection) Set(session *Session) {

	(sc)[session.Key] = session
}

func (sc sessionCollection) Delete(key string) {
	delete(sc, key)
}
