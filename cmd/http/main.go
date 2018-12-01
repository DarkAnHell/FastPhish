package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/DarkAnHell/FastPhish/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type request struct {
	// Domain is the DNS name.
	Domain string `json:"domain"`
}

type response struct {
	// Phishing indicates if the domain is a phishing site.
	Phising bool `json:"phishing"`
	// Score indicates the confidence value.
	Score uint32 `json:"score"`
}

var (
	creds credentials.TransportCredentials

	conn     *grpc.ClientConn
	aclient api.API_QueryClient
)

func analyze(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	d := json.NewDecoder(req.Body)
	var r request
	err := d.Decode(&r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "could not handle request: %v", err)
		return
	}

	domain := &api.Domain{Name: r.Domain}
	if err := aclient.Send(domain); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "could not send domain to backend: %v", err)
		return
	}
	resp, err := aclient.Recv()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "could not send domain to backend: %v", err)
		return
	}
	log.Printf("Got response with status %v: %s with score: %v", resp.GetStatus().Status, resp.GetDomain().Name, resp.GetDomain().Score)

	var phishing bool
	if resp.Domain.Score > uint32(70) {
		phishing = true
	}
	userResp := &response{
		Phising: phishing,
		Score:   resp.Domain.Score,
	}
	b, err := json.Marshal(userResp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "could not convert response to JSON: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", b)
}

func main() {
	// Create the client TLS credentials
	var err error
	creds, err = credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	// TODO: Config
	conn, err = grpc.Dial("localhost:1337", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := api.NewAPIClient(conn)
	aclient, err = client.Query(context.Background())
	if err != nil {
		log.Fatalf("could not create DomainsScoreClient: %v", err)
	}
	defer func() {
		if err := aclient.CloseSend(); err != nil {
			log.Fatalf("could not close connection: %v", err)
		}
	}()

	http.HandleFunc("/", analyze)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
