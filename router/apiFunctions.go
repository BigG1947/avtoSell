package router

import (
	"avtoSell/model"
	"encoding/json"
	"io/ioutil"
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
