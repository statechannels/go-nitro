package http

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	urlUtil "net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type clientHttpTransport struct {
	logger           *slog.Logger
	notificationChan chan []byte
	clientWebsocket  *websocket.Conn
	url              string
	wg               *sync.WaitGroup
}

// NewHttpTransportAsClient creates a transport that can be used to send http requests and a websocket connection for receiving notifications
// Initialization will block for 10 retries until the server endpoint is ready
func NewHttpTransportAsClient(url string, retryTimeout time.Duration) (*clientHttpTransport, error) {
	err := blockUntilHttpServerIsReady(url, retryTimeout)
	if err != nil {
		return nil, err
	}

	subscribeUrl, err := urlUtil.JoinPath("wss://", url, "subscribe")
	if err != nil {
		return nil, err
	}
	dialer := websocket.Dialer{
		// todo: InsecureSkipVerify needs to be removed in favor of trusting the self signed certificate
		TLSClientConfig: &tls.Config{ServerName: "statechannels.org", InsecureSkipVerify: true},
	}
	conn, _, err := dialer.Dial(subscribeUrl, nil)
	if err != nil {
		return nil, err
	}

	t := &clientHttpTransport{notificationChan: make(chan []byte, 10), clientWebsocket: conn, url: url, wg: &sync.WaitGroup{}, logger: slog.Default()}

	t.wg.Add(1)
	go t.readMessages()

	return t, nil
}

func (t *clientHttpTransport) Request(data []byte) ([]byte, error) {
	requestUrl, err := httpUrl(t.url)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				ServerName: "statechannels.org",
				// todo: InsecureSkipVerify needs to be removed in favor of trusting the self signed certificate
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := httpClient.Post(requestUrl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (t *clientHttpTransport) Subscribe() (<-chan []byte, error) {
	return t.notificationChan, nil
}

func (t *clientHttpTransport) Close() error {
	// This will also cause the go-routine to unblock waiting on `ReadMessage` and thus serves as a signal to exit
	err := t.clientWebsocket.Close()
	if err != nil {
		return err
	}
	t.wg.Wait()

	close(t.notificationChan)
	return nil
}

func (t *clientHttpTransport) readMessages() {
	t.logger.Debug("Starting to read websocket messages")
	for {
		_, data, err := t.clientWebsocket.ReadMessage()
		if err != nil {
			t.logger.Info("Websocket read error", "error", err)
			t.wg.Done()
			return
		}
		t.logger.Debug("Websocket received message", "data", string(data))

		t.notificationChan <- data
	}
}

// httpUrl joins the http prefix with the server url
func httpUrl(url string) (string, error) {
	httpUrl, err := urlUtil.JoinPath("https://", url)
	if err != nil {
		return "", err
	}
	return httpUrl, nil
}

// blockUntilHttpServerIsReady pings the health endpoint until the server is ready
func blockUntilHttpServerIsReady(url string, retryTimeout time.Duration) error {
	waitForServer := func(iteration int) {
		time.Sleep(retryTimeout * time.Duration(math.Pow(2, float64(iteration))))
	}

	httpUrl, err := httpUrl(url)
	if err != nil {
		return err
	}
	healthUrl, err := urlUtil.JoinPath(httpUrl, "health")
	if err != nil {
		return err
	}
	numAttempts := 10
	for i := 0; i < numAttempts; i++ {
		resp, err := http.Get(healthUrl)
		if err != nil {
			waitForServer(i)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return nil
		}
		waitForServer(i)
	}
	return fmt.Errorf("http server %v not ready after %d attempts", healthUrl, numAttempts)
}
