package server

import (
	"crypto/subtle"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/henrikkorsgaard/here.local/configuration"
	"github.com/henrikkorsgaard/here.local/logging"
	"github.com/rs/cors"
)

var (
	publicPath        string
	templatePath      string
	basicAuthlogin    string
	basicAuthPassword string
	ipRouter          mux.Router
	hostRouter        mux.Router
	globalRouter      mux.Router
)

func Run() {

	basicAuthlogin = configuration.GetBasicAuthLogin()
	basicAuthPassword = configuration.GetBasicAuthPassword()

	ex, err := os.Executable()
	logging.Fatal(err)

	exPath := filepath.Dir(ex)
	publicPath = filepath.Join(exPath, "./html/public")

	r := mux.NewRouter()
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir(publicPath))))

	ip := configuration.GetTP()
	host := configuration.GetLocation()
	ipRouter = r.Host(ip).Subrouter()
	hostRouter = r.Host(host + ".local").Subrouter()
	globalRouter = r.Host("here.local").Subrouter()

	ipRouter.HandlerFunc("/favicon.ico", faviconHandler).Methods("GET")
	hostRouter.HandlerFunc("/favicon.ico", faviconHandler).Methods("GET")
	globalRouter.HandlerFunc("/favicon.ico", faviconHandler).Methods("GET")

	mode := configuration.GetMode()
	if mode == configuration.CONFIG_MODE {
		setupConfig(ipRouter, hostRouter)
	} else if mode == configuration.CONFIG_MODE {
		setupSlave()
	} else {
		setupMaster()
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://here.local"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	err = http.ListenAndServe(":1337", handler)
	logging.Fatal(err)

}

func setupConfig(ipRouter *mux.Router, hostRouter *mux.Router) {
	if basicAuthlogin != "" {
		ipRouter.HandlerFunc("/", BasicAuth(configHandler, basicAuthLogin, basicAuthPassword, "Please enter username and password")).Methods("GET", "POST")
		hostRouter.HandlerFunc("/", BasicAuth(configHandler, basicAuthLogin, basicAuthPassword, "Please enter username and password")).Methods("GET", "POST")
	} else {
		ipRouter.HandlerFunc("/", configHandler).Methods("GET", "POST")
		hostRouter.HandlerFunc("/", configHandler).Methods("GET", "POST")
	}
}

func setupSlave(ipRouter *mux.Router, hostRouter *mux.Router) {

	if basicAuthlogin != "" {
		ipRouter.HandlerFunc("/config", BasicAuth(configHandler, basicAuthLogin, basicAuthPassword, "Please enter username and password")).Methods("GET", "POST")
		hostRouter.HandlerFunc("/config", BasicAuth(configHandler, basicAuthLogin, basicAuthPassword, "Please enter username and password")).Methods("GET", "POST")
	} else {
		ipRouter.HandlerFunc("/config", configHandler).Methods("GET", "POST")
		hostRouter.HandlerFunc("/config", configHandler).Methods("GET", "POST")
	}

}

func setupMaster(ipRouter *mux.Router, hostRouter *mux.Router, globalRouter *mux.Router) {
	if basicAuthlogin != "" {
		ipRouter.HandlerFunc("/config", BasicAuth(configHandler, basicAuthLogin, basicAuthPassword, "Please enter username and password")).Methods("GET", "POST")
		hostRouter.HandlerFunc("/config", BasicAuth(configHandler, basicAuthLogin, basicAuthPassword, "Please enter username and password")).Methods("GET", "POST")
		globalRouter.HandlerFunc("/config", BasicAuth(configHandler, basicAuthLogin, basicAuthPassword, "Please enter username and password")).Methods("GET", "POST")
	} else {
		ipRouter.HandlerFunc("/config", configHandler).Methods("GET", "POST")
		hostRouter.HandlerFunc("/config", configHandler).Methods("GET", "POST")
		globalRouter.HandlerFunc("/config", configHandler).Methods("GET", "POST")
	}

	//run config server
	//run api server
	//run context server
}

//BasicAuth does something
func BasicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join(publicPath, "./images/favicon.ico"))
}
