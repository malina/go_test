package main

import (
	"encoding/json"
	"net/http"
	"reflect"
  "math"
	"./model"
	"github.com/gorilla/mux"
  "strconv"
)

var Work = map[string]int32{"admin": 1000000}
var Auth = map[string]string{}


func main() {
  model.GormInit()
  defer model.GormClose()
	r := mux.NewRouter()

	r.HandleFunc("/", mainPage)
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/login/pass", changePass).Methods("POST")
	r.HandleFunc("/do", doWork).Methods("POST")

	http.Handle("/", r)
  http.ListenAndServe(":8080", nil)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<!DOCTYPE html>
		<html>
		<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<meta name="theme-color" content="#375EAB">

			<title>main page</title>
		</head>
		<body>
			Page body and some more content
		</body>
		</html>`))
}

type LoginParams struct {
  Login string
  Pass string
  NewPass string
}

func login(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  var params LoginParams
  err := decoder.Decode(&params)

  if err != nil {
    panic(err)
  }

  login := params.Login
  pass := params.Pass

	if Auth[login] == pass {
		w.WriteHeader(http.StatusOK)
	}

	user := &model.User{}
	err = user.Get(login, pass)
	if err == nil {
		Auth[login] = pass
		Work[login] = user.WorkNumber

    w.WriteHeader(http.StatusOK)
	}

	w.WriteHeader(http.StatusBadRequest)
}

func changePass(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  var params LoginParams
  err := decoder.Decode(&params)

  if err != nil {
    panic(err)
  }

  login := params.Login
  pass := params.Pass
	newPass := params.NewPass

	if Auth[login] != pass {
		w.WriteHeader(http.StatusBadRequest)
	}

  user := &model.User{}
  err = user.UpdatePass(login, pass, newPass)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
  w.WriteHeader(http.StatusOK)
}

type DTO struct {
	BigNumber int32
	Text      string
}

func doWork(w http.ResponseWriter, r *http.Request) {
	var value DTO

	login := r.FormValue("login")

	if Work[login] <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.Unmarshal([]byte(r.FormValue("value")), &value)

	v := reflect.ValueOf(value)
  structType := reflect.TypeOf(value)

  result := map[string]string{}
	for i := 0; i < v.NumField(); i++ {
    result[structType.Field(i).Name] = reverse(v.Field(i))
	}

  json, err := json.Marshal(result)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(json)
}

func reverse(val reflect.Value) string {
	switch val.Kind().String() {
	case "int64":
		fallthrough
	case "int32":
		result := math.MaxInt32-val.Interface().(int32)
		return strconv.FormatInt(int64(result), 10)
	case "string":
		var result string

		for i := len(val.Interface().(string)) - 1; i >= 0; i-- {
			result += string(val.Interface().(string)[i])
		}
		return result
	}
	return ""
}
