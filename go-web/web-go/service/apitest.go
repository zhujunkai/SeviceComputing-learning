package service

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/unrolled/render"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Age      string `json:"Age"`
	Email    string `json:"email"`
}

type MY_Error struct {
	T_error string
}

var users []User

func init() {

}

func RegisterHandler(webRoot string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			var user User
			user.Username = req.FormValue("username")
			user.Password = req.FormValue("password")
			user.Age = req.FormValue("age")
			user.Email = req.FormValue("email")
			users = append(users, user)
		}
		http.FileServer(http.Dir(webRoot+"/assets/")).ServeHTTP(w, req)
	}
}

func submitHandler(webRoot string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			req.ParseForm()
			todo := req.FormValue("todo")
			login_username := req.FormValue("username")
			login_password := req.FormValue("password")
			if todo == "1" {
				if req.Method == "POST" {
					var user User
					user.Username = req.FormValue("username")
					user.Password = req.FormValue("password")
					user.Age = req.FormValue("age")
					user.Email = req.FormValue("email")
					users = append(users, user)
				}
				http.FileServer(http.Dir(webRoot+"/assets/")).ServeHTTP(w, req)
			} else if todo == "2" {
				t_name := login_username
				t_password := login_password
				var temp User
				var a int
				a = 0
				for _, it := range users {
					if it.Username == t_name && it.Password == t_password {
						temp = it
						a = 1
					}
				}
				if a == 1 {
					t := template.Must(template.ParseFiles("assets/info.html"))
					t.Execute(w, map[string]string{
						"username": login_username,
						"password": login_password,
						"age":      temp.Age,
						"email":    temp.Email,
					})
					return
				} else {
					w.Header().Set("Content-Type", "application/json")
					my_error := MY_Error{
						T_error: "登入失败,账号密码错误",
					}
					//将user转换为json格式
					json, _ := json.Marshal(my_error)
					w.Write(json)
					return
				}
			}
		}
		http.FileServer(http.Dir(webRoot+"/assets/")).ServeHTTP(w, req)
	}
}

// 找不到文件
func NotImplementedHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "501 Not Implemented", 501)
	}
}
func apiTestHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct {
			ID      string `json:"id"`
			Content string `json:"content"`
		}{ID: "8675309", Content: "Hello from Go!"})
	}
}
