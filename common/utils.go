package common

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func StringToInt(str string) int {
	number, _ := strconv.Atoi(str)
	return number
}

func Abort(statusCode int, message string) {
	errorData := make(ErrorMap)
	AbortWithData(statusCode, message, errorData)
}

func AbortWithData(statusCode int, message string, errorData ErrorMap) {
	panic(&HttpJsonError{StatusCode: statusCode, Message: message, ErrorData: errorData})
}

func Response(statusCode int, message string, data interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	responseData := map[string]interface{}{
		"code":    statusCode,
		"message": message,
	}

	isError := false
	if !(statusCode >= 200 && statusCode < 300) {
		isError = true
	}

	if data == nil {
		data = struct{}{}
	}

	if isError {
		responseData["errors"] = data
	} else {
		responseData["data"] = data
	}

	json.NewEncoder(w).Encode(responseData)
}
