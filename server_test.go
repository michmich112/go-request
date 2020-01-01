package request_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

func CreateTestServer(suite *RequestSuite) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		log := suite.Logger.Child("server", "handler")
		headers := map[string]string{}
		for key, values := range req.Header {
			headers[key] = strings.Join(values, ", ")
		}
		log.Record("headers", headers).Infof("Request: %s %s", req.Method, req.URL)

		switch req.Method {
		case http.MethodPost:
			switch req.URL.Path {
			case "/redirect":
				res.Header().Add("Location", "/")
				res.WriteHeader(http.StatusSeeOther)
				// res.WriteHeader(http.StatusFound)
				log.Infof("Redirecting to /")
			default:
				if _, err := res.Write([]byte("body")); err != nil {
					log.Errorf("Failed to Write response to %s %s, error: %s", req.Method, req.URL, err)
				}
			}
		case http.MethodGet:
			switch req.URL.Path {
			case "/":
				if _, err := res.Write([]byte("body")); err != nil {
					log.Errorf("Failed to Write response to %s %s, error: %s", req.Method, req.URL, err)
				}
			case "/audio.mp3":
				res.Header().Add("Content-Type", "application/octet-stream")
				if _, err := res.Write([]byte(`body`)); err != nil {
					log.Errorf("Failed to Write response to %s %s, error: %s", req.Method, req.URL, err)
				}
			case "/bad_jpg_type":
				res.Header().Add("Content-Type", "image/jpg")
				if _, err := res.Write([]byte(`body`)); err != nil {
					log.Errorf("Failed to Write response to %s %s, error: %s", req.Method, req.URL, err)
				}
			case "/bad_response":
				res.Header().Add("Content-Length", "1")
				if _, err := res.Write([]byte(``)); err != nil {
					log.Errorf("Failed to Write response to %s %s, error: %s", req.Method, req.URL, err)
				}
			case "/data":
				res.Header().Add("Content-Type", "application/octet-stream")
				if _, err := res.Write([]byte(`body`)); err != nil {
					log.Errorf("Failed to Write response to %s %s, error: %s", req.Method, req.URL, err)
				}
			case "/token":
				auth := req.Header.Get("Authorization")
				if strings.Compare(auth, "Bearer ThisIsAToken") != 0 {
					res.WriteHeader(http.StatusUnauthorized)
					return
				}
				if _, err := res.Write([]byte("body")); err != nil {
					log.Errorf("Failed to Write response to %s %s, error: %s", req.Method, req.URL, err)
				}
			case "/redirect":
				res.Header().Add("Location", "/")
				res.WriteHeader(http.StatusFound)
				log.Infof("Redirecting to /")
			case "/results":
				if _, err := res.Write([]byte(`{"code": 1234}`)); err != nil {
					log.Errorf("Failed to Write response to %s %s, error: %s", req.Method, req.URL, err)
				}
			case "/timeout":
				time.Sleep(5 * time.Second)
			default:
				res.WriteHeader(http.StatusNotFound)
				if _, err := res.Write([]byte("{}")); err != nil {
					log.Errorf("Failed to Write response to %s %s, error: %s", req.Method, req.URL, err)
				}
				return
			}
		default:
			res.WriteHeader(http.StatusMethodNotAllowed)
			if _, err := res.Write([]byte("{}")); err != nil {
				log.Errorf("Failed to Write response to %s %s, error: %s", req.Method, req.URL, err)
			}
			return
		}
	}))

}
