package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type ViewData struct {
	TextData string
}

var data ViewData
var access string
var refresh string
var user string

func main() {
	fmt.Println("Client run")
	tmpl, err := template.ParseFiles("form.html")
	errExc(err)
	router := mux.NewRouter()
	data = ViewData{TextData: ""}
	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data = ViewData{TextData: ""}
		errExc(tmpl.Execute(w, data))
	}))

	router.Handle("/receive", http.HandlerFunc(receiving)).Methods("Get")
	router.Handle("/refresh", http.HandlerFunc(refreshing)).Methods("Get")
	router.Handle("/delete", http.HandlerFunc(deleting)).Methods("Get")
	router.Handle("/clear", http.HandlerFunc(clearing)).Methods("Get")

	port := os.Getenv("PORT")
	fmt.Println("Listen on port: " + port)
	errExc(http.ListenAndServe(":"+port, router))
}

func errExc(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func receiving(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("user")
	tmpl, err := template.ParseFiles("form.html")
	errExc(err)
	if username == "" {
		data.TextData = "Пожалуйста введите имя пользователя и " +
			"нажмите кнопку \"Получить\" для получения токенов"
		errExc(tmpl.Execute(w, data))
		return
	}
	response, err := http.Get("https://witmanstasjwd.herokuapp.com/receive?user=" + username)
	errExc(err)
	body, err := ioutil.ReadAll(response.Body)
	errExc(err)
	tokens := strings.Fields(string(body))
	access = tokens[1]
	refresh = tokens[3]
	user = tokens[5]
	data.TextData = "Токены были получены"
	errExc(tmpl.Execute(w, data))
}
func refreshing(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("form.html")
	errExc(err)
	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://witmanstasjwd.herokuapp.com/refresh", nil)
	errExc(err)
	request.Header.Add("Authorization", "Bearer "+access)
	request.Header.Add("refresh", refresh)
	resp, err := client.Do(request)
	errExc(err)
	body, err := ioutil.ReadAll(resp.Body)
	errExc(err)
	data.TextData = string(body)
	if strings.Contains(data.TextData, "токен обновления") || strings.Contains(data.TextData, "token") {

		errExc(tmpl.Execute(w, data))
	} else {
		tokens := strings.Fields(data.TextData)
		access = tokens[1]
		refresh = tokens[3]
		data.TextData = "Токены успешно обновлены"
		errExc(tmpl.Execute(w, data))
	}
}

func deleting(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("form.html")
	errExc(err)
	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://witmanstasjwd.herokuapp.com/delete", nil)
	errExc(err)
	request.Header.Add("Authorization", "Bearer "+access)
	request.Header.Add("refresh", refresh)
	resp, err := client.Do(request)
	errExc(err)
	body, err := ioutil.ReadAll(resp.Body)
	errExc(err)
	data.TextData = string(body)
	access, refresh = "", ""
	errExc(tmpl.Execute(w, data))

}
func clearing(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("user")
	tmpl, err := template.ParseFiles("form.html")
	errExc(err)
	if username == "" {
		data.TextData = "Пожалуйста введите имя пользователя и " +
			"нажмите кнопку \"Очистить\" для удаления всех записей пользователя"
		errExc(tmpl.Execute(w, data))
		return
	}
	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://witmanstasjwd.herokuapp.com/clear?user="+username, nil)
	errExc(err)
	request.Header.Add("Authorization", "Bearer "+access)
	resp, err := client.Do(request)
	errExc(err)
	body, err := ioutil.ReadAll(resp.Body)
	errExc(err)
	data.TextData = string(body)
	if username == user {
		access, refresh = "", ""
	}
	errExc(tmpl.Execute(w, data))
}
