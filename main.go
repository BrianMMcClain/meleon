package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/brianmmcclain/meleon/transaction"
)

type Config struct {
	RemoteHost string
	LogBody    bool
}

func main() {
	c := readConfig()

	transactionLog := make([]transaction.Transaction, 0)

	// We are an invisible proxy, so we will handle all requests the same way:
	// Read the request from the client, record it, and pass it along as intended
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.Host = c.RemoteHost
		transaction := proxyRequest(w, r, c)
		fmt.Println(transaction)
		transactionLog = append(transactionLog, transaction)
		// b, _ := json.Marshal(transactionLog)
		// fmt.Println(string(b))
	})
	http.ListenAndServe(":9999", nil)
}

func readConfig() Config {
	file, err := os.Open("meleon.json")
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}

func proxyRequest(w http.ResponseWriter, r *http.Request, c Config) transaction.Transaction {
	switch r.Method {
	case "GET":
		resp, body := handleGET(w, r, c)
		return transaction.NewTransaction(r, resp, body)
		//return recordTransaction(r, resp)
	case "POST":
		resp, body := handlePOST(w, r, c)
		return transaction.NewTransaction(r, resp, body)
	}

	return transaction.NewTransaction(nil, nil, "")
}

// GET requests are extremely straight-forward. We will receive the request,
// do all processing that we need to on our side, and pass the request along
// without modification
func handleGET(w http.ResponseWriter, r *http.Request, c Config) (*http.Response, string) {
	reqURL := fmt.Sprintf("%s%s", c.RemoteHost, r.RequestURI)
	resp, err := http.Get(reqURL)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
		return nil, ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
		return nil, ""
	}

	// Set status code
	w.WriteHeader(resp.StatusCode)

	// Send the response back to the client
	w.Write(body)

	return resp, string(body)
}

// POST requests only have the extra step of ensuring that we read the POST
// body from the request and pass it along as intended, unmodified
func handlePOST(w http.ResponseWriter, r *http.Request, c Config) (*http.Response, string) {
	reqURL := fmt.Sprintf("%s%s", c.RemoteHost, r.RequestURI)
	// TODO: Read the request body and send it in the request
	resp, err := http.Post(reqURL, "", nil)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
		return nil, ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
		return nil, ""
	}

	// Set status code
	w.WriteHeader(resp.StatusCode)

	// Send the response back to the client
	w.Write(body)

	return resp, string(body)
}
