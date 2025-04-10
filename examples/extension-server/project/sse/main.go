package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	messageChan chan string
}

type Broker struct {
	clients map[*Client]bool
	mu      sync.Mutex
}

func NewBroker() *Broker {
	return &Broker{
		clients: make(map[*Client]bool),
	}
}

func (b *Broker) AddClient(client *Client) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[client] = true
	log.Printf("Client added; total clients: %d", len(b.clients))
}

func (b *Broker) RemoveClient(client *Client) {
	b.mu.Lock()
	defer b.mu.Unlock()
	close(client.messageChan)
	delete(b.clients, client)
	log.Printf("Client removed; total clients: %d", len(b.clients))
}

var timeBroker = NewBroker()
var stockBroker = NewBroker()

func main() {
	// Handler for time updates
	http.HandleFunc("/time-stream", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		client := &Client{messageChan: make(chan string, 10)}
		timeBroker.AddClient(client)
		defer timeBroker.RemoveClient(client)

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "data: Client connected to time stream\n\n")
		flusher.Flush()

		timeZones := map[string]*time.Location{
			"UTC":     time.UTC,
			"NewYork": time.FixedZone("America/New_York", -4*3600),
			"Tokyo":   time.FixedZone("Asia/Tokyo", 9*3600),
			"London":  time.FixedZone("Europe/London", 1*3600),
		}

		for {
			select {
			case <-r.Context().Done():
				log.Println("Time stream client disconnected")
				return
			default:
				now := time.Now()
				for name, loc := range timeZones {
					timeStr := now.In(loc).Format(time.RFC3339)
					fmt.Fprintf(w, "data: Time in %s: %s\n\n", name, timeStr)
					flusher.Flush()
				}
				time.Sleep(2 * time.Second)
			}
		}
	})

	// Handler for stock updates
	http.HandleFunc("/stock-stream", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		client := &Client{messageChan: make(chan string, 10)}
		stockBroker.AddClient(client)
		defer stockBroker.RemoveClient(client)

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "data: Client connected to stock stream\n\n")
		flusher.Flush()

		stocks := []string{"AAPL", "GOOGL", "TSLA"}
		getStockPrice := func(stock string) float64 {
			return 100.0 + float64(time.Now().Second()) // Simulated price
		}

		for {
			select {
			case <-r.Context().Done():
				log.Println("Stock stream client disconnected")
				return
			default:
				for _, stock := range stocks {
					price := getStockPrice(stock)
					fmt.Fprintf(w, "data: %s Stock Price: %.2f\n\n", stock, price)
					flusher.Flush()
				}
				time.Sleep(2 * time.Second)
			}
		}
	})

	// Start server
	log.Println("Starting server on :5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal(err)
	}
}
