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

func handlePOST(w http.ResponseWriter, r *http.Request, c Config) {
  reqURL := fmt.Sprintf("%s%s", c.RemoteHost, r.RequestURI)
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
