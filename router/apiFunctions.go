package router

import (
	"avtoSell/model"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func checkPhone(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	phone, err := ioutil.ReadAll(request.Body)
	responseMap := make(map[string]string)
	if err != nil {
		writer.WriteHeader(500)
		responseMap["error"] = err.Error()
		responseMap["ok"] = "false"
		response, err := json.Marshal(responseMap)
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}
		writer.Write(response)
		return
	}
	if model.CheckUserPhoneExist(string(phone), connection) {
		writer.WriteHeader(200)
		responseMap["ok"] = "true"
		response, err := json.Marshal(responseMap)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
			return
		}
		writer.Write(response)
		return
	} else {
		writer.WriteHeader(200)
		responseMap["ok"] = "false"
		response, err := json.Marshal(responseMap)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
			return
		}
		writer.Write(response)
		return
	}
}

func checkEmail(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	email, err := ioutil.ReadAll(request.Body)
	responseMap := make(map[string]string)
	if err != nil {
		writer.WriteHeader(500)
		responseMap["error"] = err.Error()
		responseMap["ok"] = "false"
		response, err := json.Marshal(responseMap)
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}
		writer.Write(response)
		return
	}
	if model.CheckUserEmailExist(string(email), connection) {
		writer.WriteHeader(200)
		responseMap["ok"] = "true"
		response, err := json.Marshal(responseMap)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
			return
		}
		writer.Write(response)
		return
	} else {
		writer.WriteHeader(200)
		responseMap["ok"] = "false"
		response, err := json.Marshal(responseMap)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
			return
		}
		writer.Write(response)
		return
	}
}

func checkLogin(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	login, err := ioutil.ReadAll(request.Body)
	responseMap := make(map[string]string)
	if err != nil {
		writer.WriteHeader(500)
		responseMap["error"] = err.Error()
		responseMap["ok"] = "false"
		response, err := json.Marshal(responseMap)
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}
		writer.Write(response)
		return
	}
	if model.CheckUserLoginExist(string(login), connection) {
		writer.WriteHeader(200)
		responseMap["ok"] = "true"
		response, err := json.Marshal(responseMap)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
			return
		}
		writer.Write(response)
		return
	} else {
		writer.WriteHeader(200)
		responseMap["ok"] = "false"
		response, err := json.Marshal(responseMap)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
			return
		}
		writer.Write(response)
		return
	}
}

func ApiNewsAll(writer http.ResponseWriter, request *http.Request) {
	var news model.NewsList
	if err := news.GetAll(connection); err != nil {
		log.Printf("Error in api news all routers: %s\n", err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(news)
	if err != nil {
		log.Printf("Error in api news all routers: %s\n", err)
		return
	}
	writer.WriteHeader(200)
	if _, err := writer.Write(jsonResponse); err != nil {
		log.Printf("%s\n", err)
	}
	return
}
