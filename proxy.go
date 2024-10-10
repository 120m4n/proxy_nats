package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/nats-io/nats.go"
)

type natsConnection struct {
	nc *nats.Conn
	topic string
}

//Funcion para crearun proxy hacia un backend especifico
func newProxy(target string, natsConn *natsConnection) http.HandlerFunc{
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	return func(w http.ResponseWriter, r *http.Request){
		//Log de la solicitud
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)


		// Check if the request method is POST
		if r.Method == http.MethodPost {
			// Read the request body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusInternalServerError)
				return
			}
			// Restore the request body for the proxy
			r.Body = io.NopCloser(strings.NewReader(string(body)))
		
			// Send a copy of the request body as a NATS message
			if err := natsConn.nc.Publish(natsConn.topic, body); err != nil {
				log.Printf("Failed to publish message to NATS: %v", err)
			}
		}

		//modificar la solicitud antes de enviarla al backend si es necesario
		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = url.Host

		// Reenviar la solicitud al backend
		proxy.ServeHTTP(w, r)
	}
}

func main(){
	// Connect to NATS server
	natsUrl := os.Getenv("NATS_URL")
	if natsUrl == "" {
		natsUrl = nats.DefaultURL
	}
	
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	natsConn := &natsConnection{
		nc: nc,
		topic: "example.topic",
	}
	// Definir diferentes backends para diferentes rutas
	backend1 := "http://localhost:8081"
	backend2 := "http://localhost:8082"

	// Ruta pra el backend 1
	http.HandleFunc("/mock_1/", func(w http.ResponseWriter, r *http.Request){
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/mock_1")
		newProxy(backend1, natsConn)(w, r)
	})
	http.HandleFunc("/mock_2/", func(w http.ResponseWriter, r *http.Request){
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/mock_2")
		newProxy(backend2, natsConn)(w, r)
	})	
	
	// Iniciar el servidor en el puerto 8080
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)

}