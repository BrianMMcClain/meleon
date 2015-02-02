package main

import (
  "net/http"
  "fmt"
  "io/ioutil"
  "os"
  "encoding/json"
) 

type Config struct {
  RemoteHost string
  LogBody bool
}

func main() {
  c := readConfig()

  // We are an invisible proxy, so we will handle all requests the same way:
  // Read the request from the client, record it, and pass it along as intended
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    proxyRequest(w, r, c)
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

func proxyRequest(w http.ResponseWriter, r *http.Request, c Config) {
  recordRequest(r, c)
  switch r.Method {
    case "GET":
      handleGET(w, r, c)
    case "POST":
      handlePOST(w, r, c)
  }
}

// Record the method, headers, address and potentially the
// body of the reqest
func recordRequest(r *http.Request, c Config) {
  fmt.Println("--------------------------------------------------------------------")
  fmt.Printf("> %s %s\n", r.Method, r.RequestURI)
  for k, _ := range r.Header {
    fmt.Printf("> %s: %s\n", k, r.Header[k])
  }

  if c.LogBody {
    switch r.Method {
    case "POST":
      recordRequestBody(r)
    }
  }

  fmt.Println()
}

// Read the request body and record it
func recordRequestBody(r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    panic(fmt.Sprintf("%v", err))
    return
  }

  fmt.Printf("\n%v\n", string(body))
}

func recordResponse(r *http.Response, body string, c Config) {
  fmt.Printf("< %s\n", r.Status)
  for k, _ := range r.Header {
    fmt.Printf("< %s: %s\n", k, r.Header[k])
  }

  if c.LogBody {
    fmt.Printf("\n%s\n", body)
  }
  fmt.Println()
}

// GET requests are extremely straight-forward. We will receive the request,
// do all processing that we need to on our side, and pass the request along
// without modification
func handleGET(w http.ResponseWriter, r *http.Request, c Config) {
  reqURL := fmt.Sprintf("%s%s", c.RemoteHost, r.RequestURI)
  resp, err := http.Get(reqURL)
  if err != nil {
    panic(fmt.Sprintf("%v", err))
    return
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    panic(fmt.Sprintf("%v", err))
    return
  }

  recordResponse(resp, string(body), c)

  // Copy headers
  cRespHeader := w.Header()
  for k, _ := range resp.Header {
    cRespHeader[k] = resp.Header[k]
  }

  // Set status code
  w.WriteHeader(resp.StatusCode)

  // Send the response back to the client
  w.Write(body)
}

// POST requests only have the extra step of ensuring that we read the POST
// body from the request and pass it along as intended, unmodified
func handlePOST(w http.ResponseWriter, r *http.Request, c Config) {
  reqURL := fmt.Sprintf("%s%s", c.RemoteHost, r.RequestURI)
  // TODO: Read the request body and send it in the request
  resp, err := http.Post(reqURL, "", nil)
  if err != nil {
    panic(fmt.Sprintf("%v", err))
    return
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    panic(fmt.Sprintf("%v", err))
    return
  }

  recordResponse(resp, string(body), c)

  // Copy headers
  cRespHeader := w.Header()
  for k, _ := range resp.Header {
    cRespHeader[k] = resp.Header[k]
  }

  // Set status code
  w.WriteHeader(resp.StatusCode)

  // Send the response back to the client
  w.Write(body)
}
