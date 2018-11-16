package service

import (
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {

	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	webRoot := os.Getenv("WEBROOT")
	if len(webRoot) == 0 {
		if root, err := os.Getwd(); err != nil {
			panic("Could not retrive working directory")
		} else {
			webRoot = root
			//fmt.Println(root)
		}
	}
	//实现对静态文件和js的引用
	mx.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(webRoot+"/assets/"))))
	//unknown的错误显示
	mx.HandleFunc("/unknown", NotImplementedHandler(formatter)).Methods("GET")
	mx.HandleFunc("/api/test", apiTestHandler(formatter)).Methods("GET")
	mx.HandleFunc("/register.html", RegisterHandler(webRoot)).Methods("POST")
	mx.HandleFunc("/", submitHandler(webRoot)).Methods("POST")
	mx.PathPrefix("/").Handler(http.FileServer(http.Dir(webRoot + "/assets/")))

}
