package loki

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/BugBuster401/logx/slogx"
)

type LokiClient struct {
	URL     string
	AppName string
	Client  *http.Client
	wg      sync.WaitGroup
}

func NewLokiClient(url, appName string) (*LokiClient, error) {
	lc := &LokiClient{
		URL:     url,
		AppName: appName,
		Client:  &http.Client{Timeout: 10 * time.Second},
	}

	if err := lc.ping(); err != nil {
		return nil, fmt.Errorf("loki is not reachable: %w", err)
	}

	return lc, nil
}

func (lc *LokiClient) ping() error {
	req, err := http.NewRequest(http.MethodGet, lc.URL+"/ready", nil)
	if err != nil {
		return err
	}

	resp, err := lc.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	return nil
}

type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type LokiPayload struct {
	Streams []LokiStream `json:"streams"`
}

func (lc *LokiClient) Send(ctx context.Context, entry slogx.LogEntry) error {
	labels := map[string]string{
		"app":   lc.AppName,
		"level": entry.Level,
	}

	data := make(map[string]any, len(entry.Fields)+2)
	data["message"] = entry.Message

	for k, v := range entry.Fields {
		data[k] = v
	}

	entryJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal log entry: %w", err)
	}

	stream := LokiStream{
		Stream: labels,
		Values: [][]string{
			{
				fmt.Sprintf("%d", entry.Timestamp.UnixNano()),
				string(entryJSON),
			},
		},
	}

	payload := LokiPayload{
		Streams: []LokiStream{stream},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal loki payload: %w", err)
	}

	url := fmt.Sprintf("%s/loki/api/v1/push", lc.URL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create loki request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	lc.wg.Add(1)
	go func() {
		defer lc.wg.Done()

		resp, err := lc.Client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "LokiClient error: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "LokiClient error: failed read response body: %s", err.Error())
				return
			}
			fmt.Fprintf(os.Stderr,
				"LokiClient error: unexpected status %d, response: %s\n",
				resp.StatusCode, strings.TrimSpace(string(body)),
			)
		}
	}()

	return nil
}

func (lc *LokiClient) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		lc.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
