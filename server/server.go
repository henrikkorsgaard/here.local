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
	publicPath   string
	templatePath string
	ipRouter     *mux.Router
	hostRouter   *mux.Router
	globalRouter *mux.Router
	configs      *configuration.Configuration
)

func Run() {
	devMode := viper.GetBool("dev") // retreive value from viper
	configs = configuration.GetConfiguration()
	var ip string
	var host string
	var mode string
	var exPath string

	if devMode {
		ip = "127.0.0.1"
		host = "localhost"
		mode = "CONFIG"
	} else {

		ip = configs.GetIP()
		host = configs.GetLocation()
		mode = configs.GetMode()

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
	r.HandleFunc("/config", BasicAuth(configHandler)).Methods("GET", "POST")
}

func setupSlave() {

	ipRouter.HandleFunc("/config", BasicAuth(configHandler)).Methods("GET", "POST")
	hostRouter.HandleFunc("/config", BasicAuth(configHandler)).Methods("GET", "POST")

}

func setupMaster() {

	ipRouter.HandleFunc("/config", BasicAuth(configHandler)).Methods("GET", "POST")
	hostRouter.HandleFunc("/config", BasicAuth(configHandler)).Methods("GET", "POST")
	globalRouter.HandleFunc("/config", BasicAuth(configHandler)).Methods("GET", "POST")

	//run config server
	//run api server
	//run context server
}

//BasicAuth handles authentication to the config server
//It will fetch the user/pass information from the configuration file on every query
//to ensure that any updates are reflected without a server restart
func BasicAuth(handler http.HandlerFunc) http.HandlerFunc {

	username := configs.GetAuthenticationLogin()
	password := configs.GetAuthenticationPassword()

	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter username and password"`)
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
