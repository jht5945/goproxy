package main

// code from : https://gist.github.com/fabrizioc1/4327250
import (
    "fmt"
	"io"
	"log"
	"net/http"
)

type HttpConnection struct {
	Request  *http.Request
	Response *http.Response
}

type HttpConnectionChannel chan *HttpConnection

var connChannel = make(HttpConnectionChannel)

func PrintHTTP(conn *HttpConnection) {
	fmt.Println("[INFO] HTTP Request-------------------------------------------------------")
	fmt.Printf("%v %v\n", conn.Request.Method, conn.Request.RequestURI)
	for k, v := range conn.Request.Header {
		fmt.Println(k, ":", v)
	}
	fmt.Println("[INFO] HTTP Response------------------------------------------------------")
	fmt.Printf("HTTP/1.1 %v\n", conn.Response.Status)
	for k, v := range conn.Response.Header {
		fmt.Println(k, ":", v)
	}
	fmt.Println(conn.Response.Body)
	fmt.Println("[INFO] HTTP Ended---------------------------------------------------------")
}

type Proxy struct {
}

func NewProxy() *Proxy { return &Proxy{} }

func (p *Proxy) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error
	var req *http.Request
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	log.Printf("%v %v", r.Method, r.RequestURI)
	targetURI := "https://www.baidu.com" + r.RequestURI
	req, err = http.NewRequest(r.Method, targetURI, r.Body)
	for name, values := range r.Header {
		log.Printf("%v=%v", name, values)
		for _, value := range values {
			req.Header.Set(name, value)
		}
	}
	//rep.Host = r.Host
	resp, err = client.Do(req)
	r.Body.Close()

	// combined for GET/POST
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	conn := &HttpConnection{r, resp}

	for k, v := range resp.Header {
		wr.Header().Set(k, v[0])
	}
	wr.WriteHeader(resp.StatusCode)
	io.Copy(wr, resp.Body)
	resp.Body.Close()

	PrintHTTP(conn)
}

func main() {
	proxy := NewProxy()
	fmt.Println("==============================")
	fmt.Println("Listen: 12345")
	err := http.ListenAndServe(":12345", proxy)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
	// http.ListenAndServeTLS(":443", cert, key, proxy)
}