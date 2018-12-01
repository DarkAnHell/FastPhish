package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"github.com/DarkAnHell/FastPhish/pkg/db"
	"github.com/DarkAnHell/FastPhish/api"
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

	conn *grpc.ClientConn
	analconn *grpc.ClientConn

	dscli api.DB_GetDomainsScoreClient
	aclient api.Analyzer_AnalyzeClient
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
	if err := dscli.Send(domain); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "could not send domain to backend: %v", err)
		return
	}

	resp, err := dscli.Recv()
	if err == db.ErrDBNotFound {
		// send to analyze
		if err := aclient.Send(domain); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "could not send domain to analysis module: %v", err)
			return
		}

		resp, err := aclient.Recv()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "could not handle analysis result: %v", err)
			return
		}


		var phishing bool
		// TODO: config.
		if resp.Domain.Score > uint32(70) {
			phishing = true
		}
		userResp := &response{
			Phising: phishing,
			Score: resp.Domain.Score,
		}

		b, err := json.Marshal(userResp)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "could not convert response to JSON: %v", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, b)
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "could not handle backend response: %v", err)
		return
	}

	var phishing bool
	if resp.Domain.Score > uint32(70) {
		phishing = true
	}
	userResp := &response{
		Phising: phishing,
		Score: resp.Domain.Score,
	}
	b, err := json.Marshal(userResp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "could not convert response to JSON: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, b)
}

func main() {
	// Create the client TLS credentials
	var err error
	creds, err = credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	// TODO: Config
	conn, err = grpc.Dial("localhost:50000", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// TODO: Config.
	analconn, err = grpc.Dial("localhost:1338", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer analconn.Close()

	client := api.NewDBClient(conn)
	dscli, err = client.GetDomainsScore(context.Background())
	if err != nil {
		log.Fatalf("could not create DomainsScoreClient: %v", err)
	}
	defer func() {
		if err := dscli.CloseSend(); err != nil {
			log.Fatalf("could not close connection: %v", err)
		}
	}()

	ac := api.NewAnalyzerClient(analconn)
	aclient, err = ac.Analyze(context.Background())
	if err != nil {
		log.Fatalf("could not create analyzer: %v", err)
	}
	defer func() {
		if err := aclient.CloseSend(); err != nil {
			log.Fatalf("could not close connection: %v", err)
		}
	}()

	http.HandleFunc("/", analyze)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
