package transaction

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Transaction struct {
	Timestamp    time.Time
	Request      *http.Request
	Response     *http.Response
	ResponseBody string
}

func NewTransaction(req *http.Request, resp *http.Response, body string) Transaction {
	t := Transaction{time.Now(), req, resp, body}
	return t
}

func (t Transaction) String() string {
	// Request Method and URI
	outS := fmt.Sprintf("> %s %s\n", t.Request.Method, t.Request.RequestURI)
	for k, _ := range t.Request.Header {
		outS += fmt.Sprintf("> %s: %s\n", k, t.Request.Header[k])
	}

	switch t.Request.Method {
	case "POST":
		body, err := ioutil.ReadAll(t.Request.Body)
		if err != nil {
			panic(fmt.Sprintf("%v", err))
			return ""
		}

		outS += fmt.Sprintf("\n%v\n", string(body))
	}

	outS += "\n" // Newlinr for formatting

	// Response
	outS += fmt.Sprintf("< %s\n", t.Response.Status)
	for k, _ := range t.Response.Header {
		outS += fmt.Sprintf("< %s: %s\n", k, t.Response.Header[k])
	}

	outS += fmt.Sprintf("\n%v\n", t.ResponseBody)

	return outS
}
