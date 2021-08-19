package restapi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/hasujin/ston_common/log2paul"
)

var (
	log = log2paul.GetLogger("REST-API")
)

func SendHttpRequest(method, url, auth string, params map[string]string, body []byte) (int, []byte, error) {
	return sendHttpRequestInternal(method, url, auth, params, body, 0)
}

func SendHttpRequestWithTimeout(method, url, auth string, params map[string]string, body []byte, timeout time.Duration) (int, []byte, error) {
	return sendHttpRequestInternal(method, url, auth, params, body, timeout)
}

func sendHttpRequestInternal(method, url, auth string, params map[string]string, body []byte, timeout time.Duration) (int, []byte, error) {
	buff := bytes.NewBuffer(body)

	req, err := http.NewRequest(method, url, buff)
	if err != nil {
		log.Errorf("failed to create request. [%s] %s. param: %+v, body: %s, reason: %+v", method, url, string(body), err)
		return -1, nil, ErrCreateRequest
	}
	setAuthorization(req, auth)
	setQueryParameters(req, params)
	req.Header.Set("Content-Type", "application/json")

	var client *http.Client
	if timeout > 0 {
		client = &http.Client{
			Timeout: timeout,
		}
	} else {
		client = &http.Client{}
	}
	response, err := client.Do(req)
	if err != nil {
		log.Errorf("failed to send request. [%s] %s. body: %s, reason: %+v", method, url, string(body), err)
		if isTimeout(err) {
			return -1, nil, ErrRequestTimeout
		}
		return -1, nil, ErrSendRequest
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Errorf("failed to close response. [%s] %s. body: %s, reason: %+v", method, url, string(body), err)
		}
	}()
	if http.StatusNotFound == response.StatusCode {
		log.Errorf("requested resource could not be found. [%s] %s", method, url)
		return response.StatusCode, nil, ErrInvalidStatus
	}

	var responseBody []byte
	if responseBody, err = ioutil.ReadAll(response.Body); err != nil {
		log.Errorf("failed to read response. [%s] %s. body: %s, reason: %+v", method, url, string(body), err)
		return response.StatusCode, nil, ErrReadResponse
	}

	return response.StatusCode, responseBody, nil
}

const errStrTimeout = "context deadline exceeded"

func isTimeout(err error) bool {
	if strings.Contains(err.Error(), errStrTimeout) {
		return true
	}
	return false
}

func setAuthorization(req *http.Request, auth string) {
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
}

func setQueryParameters(req *http.Request, params map[string]string) {
	if nil == params {
		return
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
}

func WriteError(w http.ResponseWriter, statusCode int, e *Error) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Status: false,
		Error:  e,
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Errorf("failed to convert json obj to bytes. %+v: %v", response, err)
		return
	}
	_, err = w.Write(responseBytes)
	if err != nil {
		log.Errorf("failed to write message. %v", err)
	}
}

func WriteErrorCode(w http.ResponseWriter, code int) {
	WriteError(w, code, &Error{
		Code:    code,
		Message: http.StatusText(code),
	})
}

func WriteErrorCodeMsg(w http.ResponseWriter, code int, msg string) {
	WriteError(w, code, &Error{
		Code:    code,
		Message: msg,
	})
}

func WriteResponse(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Status: true,
		Data:   data,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Errorf("failed to encode data. %+v: %v", response, err)
		WriteError(w, http.StatusInternalServerError, &Error{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)})
	}
}

func GenBearerToken(token string) string {
	return "Bearer " + token
}

func GetBearerToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	splitToken := strings.Split(token, "Bearer ")
	if len(splitToken) < 2 {
		return ""
	}
	return splitToken[1]
}

func GenBasicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func DecodeRequestBody(w http.ResponseWriter, req *http.Request, data interface{}) error {
	err := json.NewDecoder(req.Body).Decode(data)
	if err != nil {
		log.Errorf("failed to decode data. %v: %v", req.Body, err)
		return err
	}
	return nil
}

func ToBytes(i interface{}) (bytes []byte, err error) {
	if bytes, err = json.Marshal(i); err != nil {
		log.Errorf("failed to convert to bytes. %+v", i)
		return nil, errors.New("malformed json")
	}
	return bytes, nil
}
