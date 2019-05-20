package server

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/henrikkorsgaard/here.local/configuration"
	"github.com/henrikkorsgaard/here.local/logging"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

var (
	publicPath        string
	templatePath      string
	basicAuthLogin    string
	basicAuthPassword string
	ipRouter          *mux.Router
	hostRouter        *mux.Router
	globalRouter      *mux.Router
)

func Run() {
	devMode := viper.GetBool("dev") // retreive value from viper

	var ip string
	var host string
	var mode string
	var exPath string

	if devMode {
		basicAuthLogin = ""
		basicAuthPassword = ""
		ip = "127.0.0.1"
		host = "localhost"
		mode = "CONFIG"
	} else {
		basicAuthLogin = configuration.GetBasicAuthLogin()
		basicAuthPassword = configuration.GetBasicAuthPassword()
		ip = configuration.GetIP()
		host = configuration.GetLocation()
		mode = configuration.GetMode()

	}

	fmt.Println(ip, host, mode)
	ex, err := os.Executable()
	logging.Fatal(err)
	exPath = filepath.Dir(ex)
	publicPath = filepath.Join(exPath, "./html/public")
	templatePath = filepath.Join(exPath, "./html/templates")
	fmt.Println(templatePath)
	r := mux.NewRouter()
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir(publicPath))))

	setupConfigurationServer(r)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://here.local"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	err = http.ListenAndServe(":1337", handler)
	logging.Fatal(err)
}

func setupConfigurationServer(r *mux.Router) {

	if basicAuthLogin != "" {
		r.HandleFunc("/config", BasicAuth(configHandler, basicAuthLogin, basicAuthPassword, "Please enter username and password")).Methods("GET", "POST")
	} else {
		r.HandleFunc("/config", configHandler).Methods("GET", "POST")
	}
}

func setupSlave() {

	if basicAuthLogin != "" {
		ipRouter.HandleFunc("/config", BasicAuth(configHandler, basicAuthLogin, basicAuthPassword, "Please enter username and password")).Methods("GET", "POST")
		hostRouter.HandleFunc("/config", BasicAuth(configHandler, basicAuthLogin, basicAuthPassword, "Please enter username and password")).Methods("GET", "POST")
	} else {
		ipRouter.HandleFunc("/config", configHandler).Methods("GET", "POST")
		hostRouter.HandleFunc("/config", configHandler).Methods("GET", "POST")
	}

}

func setupMaster() {
	if basicAuthLogin != "" {
		ipRouter.HandleFunc("/config", BasicAuth(configHandler, basicAuthLogin, basicAuthPassword, "Please enter username and password")).Methods("GET", "POST")
		hostRouter.HandleFunc("/config", BasicAuth(configHandler, basicAuthLogin, basicAuthPassword, "Please enter username and password")).Methods("GET", "POST")
		globalRouter.HandleFunc("/config", BasicAuth(configHandler, basicAuthLogin, basicAuthPassword, "Please enter username and password")).Methods("GET", "POST")
	} else {
		ipRouter.HandleFunc("/config", configHandler).Methods("GET", "POST")
		hostRouter.HandleFunc("/config", configHandler).Methods("GET", "POST")
		globalRouter.HandleFunc("/config", configHandler).Methods("GET", "POST")
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
