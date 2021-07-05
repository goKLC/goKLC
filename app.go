package goKLC

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/jinzhu/gorm"
	"io"
	"net/http"
	"strings"
)

type App struct {
	key         string
	logger      Log
	response    Response
	cookie      Cookie
	assetPrefix string
	assetFolder string
}

var _app *App
var _routeTree *RouteNode
var _middlewareList *MiddlewareNode
var _routeNameList routeNameList
var _configCollector configCollector
var _config Config
var _sessionCollector sessionCollector
var _DB *gorm.DB
var _auth *Auth

func GetApp() *App {

	if _app == nil {
		_routeTree = NewRouteTree()
		_middlewareList = NewMiddlewareNode()
		_routeNameList = NewRouteNameList()
		_configCollector = newConfigCollector()
		_config = NewConfig()
		_sessionCollector = newSessionCollector()
		_app = &App{}
	}

	return _app
}

func (a *App) Run() {

	a.key = _config.Get("AppKey", "").(string)
	port := _config.Get("HttpPort", 8080)
	httpAddr := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(httpAddr, a)

	a.Log().Error(err.Error(), nil)

}

func (a *App) Config() Config {

	return NewConfig()
}

func (a *App) Route() Route {
	rg := NewRouteGroup()
	rg.prefix = ""

	return rg.Route()
}

func GetRoute(name string) string {

	return _routeNameList.Get(name)
}

func (a *App) Middleware(m MiddlewareInterface) {
	if _middlewareList == nil {
		_middlewareList.middleware = m
	} else {
		mn := NewMiddlewareNode()
		mn.middleware = m

		_middlewareList.AddChild(mn)
	}
}

func (a *App) GetSessionKey() string {
	key := make([]byte, 128)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {

		a.Log().Error(err.Error(), nil)
	}

	return base64.URLEncoding.EncodeToString(key)
}

func (a *App) SetLogger(logger Log) {
	a.logger = logger
}

func (a *App) Log() Log {

	return a.logger
}

func (a *App) SetResponse(response Response) {
	a.response = response
}

func (a *App) Response() Response {
	return a.response
}

func (a *App) SetCookie(cookie Cookie) {
	a.cookie = cookie
}

func (a *App) Cookie() Cookie {
	return a.cookie
}

func (a *App) GetDBURL(dbType DBType) string {

	return connectDB(dbType)
}

func (a *App) SetDB(db *gorm.DB) {

	_DB = db
}

func (a *App) DB() *gorm.DB {

	return _DB
}

func (a *App) Auth() *Auth {
	if _auth == nil {

		_auth = &Auth{}
	}

	return _auth
}

func (a *App) setAssetConf() {
	a.assetPrefix = fmt.Sprintf("/%s", _config.Get("AssetsPrefix", "assets"))
	a.assetFolder = fmt.Sprintf("./%s", _config.Get("AssetsFolder", "public"))
}

func (a *App) Assets(path string) string {
	if len(a.assetPrefix) == 0 {
		a.setAssetConf()
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return a.assetPrefix + path
}

func (a *App) GetRouteByName(name string) string {
	return fmt.Sprintf("%s/%s", _config.Get("AppDomain", ""), GetRoute(name))
}

func (a *App) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if len(a.assetPrefix) == 0 {
		a.setAssetConf()
	}

	if strings.HasPrefix(req.URL.Path, a.assetPrefix) {
		fileServer := http.StripPrefix(a.assetPrefix, http.FileServer(http.Dir(a.assetFolder)))
		fileServer.ServeHTTP(rw, req)

		return
	}

	route, ok, params := match(req)
	request := NewRequest(req, params)
	var response Response
	var middleware *MiddlewareNode

	if !ok {

		fmt.Fprintf(rw, "404", nil)
		return
	}

	response, middleware = _middlewareList.Handle(request)

	if response == nil {
		if route.middleware != nil {

			var rm *MiddlewareNode
			response, rm = route.middleware.Handle(request)

			if response == nil {
				response = route.controller(request)
			}

			rm.Terminate(response)
			rm = nil

		} else {
			response = route.controller(request)
		}
	}

	middleware.Terminate(response)

	if len(response.GetCookies()) > 0 {
		writeCookies(rw, response)
	}

	if len(response.GetHeaders()) > 0 {
		writeHeaders(rw, response)
	}

	rw.WriteHeader(response.GetStatusCode())
	rw.Write([]byte(response.GetBody()))

	request = nil
	response = nil
}

func writeHeaders(rw http.ResponseWriter, r Response) {
	headers := r.GetHeaders()

	for key, value := range headers {
		rw.Header().Add(key, value)
	}
}

func writeCookies(rw http.ResponseWriter, r Response) {
	for _, cookie := range r.GetCookies() {
		ck := cookie.(Cookie)

		c := http.Cookie{
			Name:   ck.GetName(),
			Value:  ck.GetValue(),
			MaxAge: ck.GetDuration(),
			Path:   ck.GetPath(),
		}

		http.SetCookie(rw, &c)
	}
}
