package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type bazRequest struct {
	X string
}

type bazResponse struct {
	Y string
}

func handleBaz(req interface{}) (resp interface{}, err error) {
	req2 := req.(bazRequest)
	resp = bazResponse{Y: req2.X}
	return resp, nil
}

type bodyParser func(r *http.Request) (interface{}, error)
type myHandler func(req interface{}) (interface{}, error)

func wrapper(method string, bodyParser bodyParser, handler myHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//log.Printf("Method %s, params: %+v", method, body)

		body, err := bodyParser(r)
		if err != nil {
			log.Printf("body parse error: %s", err)
			fmt.Fprintf(w, "error_body_parse\n")
			return
		}

		// Run handler
		resp, err := handler(body)

		// Log response
		log.Printf("Response: %+v", resp)

		// Marshal response
		json.NewEncoder(w).Encode(resp)
	}
}

// ServeHTTP runs HTTP server
func ServeHTTP() {
	log.Printf("Starting serving HTTP")

	//addHandler("Foo", fooBody{}, handleFoo)
	//addHandler("Bar", barBody{}, handleBar)
	http.HandleFunc("/Baz", wrapper(
		"Baz",
		func(r *http.Request) (interface{}, error) {
			req := bazRequest{}
			err := json.NewDecoder(r.Body).Decode(&req)
			return req, err
		},
		handleBaz,
	))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
