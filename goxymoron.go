package main

import "bytes"
import "io/ioutil"
import "log"
import "net/http"
import "runtime"
import "strconv"

type ResponseForwarder struct {
}

func Transform(body []byte) (output []byte) {
    return bytes.Replace(
        body,
        []byte(`http://watch.sling.com`),
        []byte(`http://hahaha.com`),
        -1)
}

func (this ResponseForwarder) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
    backend_url := *req.URL
    backend_url.Host = "localhost:23206"
    backend_url.Scheme = "http"
    backend_request, err := http.NewRequest(req.Method, backend_url.String(), nil)
    if err != nil {
        resp.WriteHeader(503)
        resp.Write([]byte(`Failed to build backend request`))
        return
    }

    backend_response, err := http.DefaultClient.Do(backend_request)
    var body []byte
    if err == nil {
        body, err = ioutil.ReadAll(backend_response.Body)
        backend_response.Body.Close()
    }

    if err != nil {
        resp.WriteHeader(503)
        resp.Write([]byte(`Backend request failed`))
        return
    }

    body = Transform(body)

    resp.WriteHeader(200)
    resp.Header().Set("Content-Length", strconv.Itoa(len(body)))
    resp.Write(body)
}

func main() {
    runtime.GOMAXPROCS(1)

    response_forwarder := ResponseForwarder{}

    server := http.Server{
        Addr: ":8081",
        Handler: response_forwarder,
    }

    err := server.ListenAndServe()
    log.Fatal(err)
}
