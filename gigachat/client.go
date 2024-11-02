package gigachat

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"main/config"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	masterToken   string
	accessToken   string
	expiresAt     time.Time
	timeout       time.Duration
	GigaChatModel string
}

type UpdateTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`

	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func NewClient(cfg *config.GigaChatConfig) (*Client, error) {
	client := &Client{
		masterToken:   cfg.Token,
		timeout:       time.Duration(cfg.Timeout) * time.Second,
		GigaChatModel: cfg.Model,
	}

	if err := client.updateToken(); err != nil {
		return nil, err
	}

	go func() {
		for {
			time.Sleep(client.expiresAt.Sub(time.Now()))
			log.Println("generate token")
			if err := client.updateToken(); err != nil {
				log.Printf("Error updating token: %s", err.Error())
				continue
			}

		}
	}()

	return client, nil
}

func newHttpClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func (c *Client) updateToken() error {
	data := url.Values{}
	data.Set("scope", "GIGACHAT_API_CORP")

	req, err := http.NewRequest("POST", "https://ngw.devices.sberbank.ru:9443/api/v2/oauth", bytes.NewBufferString(data.Encode()))

	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+c.masterToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Charset", "utf-8")
	req.Header.Add("RqUID", uuid.New().String())

	res, err := newHttpClient(c.timeout).Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	tokenResponse := &UpdateTokenResponse{}
	if err := json.Unmarshal(body, tokenResponse); err != nil {
		return err
	}

	if tokenResponse.Error.Code != 0 {
		return errors.New(tokenResponse.Error.Message)
	}

	c.accessToken = tokenResponse.AccessToken
	c.expiresAt = time.UnixMilli(tokenResponse.ExpiresAt)

	return nil
}

func (c *Client) Request(message string) (string, error) {
	jsonData := fmt.Sprintf(`{"model":"%s","messages":[{"role":"user","content":"%s"}],"stream":false,"repetition_penalty":1}`, c.GigaChatModel, message)

	fmt.Println(jsonData)

	req, err := http.NewRequest("POST", "https://gigachat.devices.sberbank.ru/api/v1/chat/completions", bytes.NewBuffer([]byte(jsonData)))

	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.accessToken)

	resp, err := newHttpClient(c.timeout).Do(req)
	if err != nil {
		return "", err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status code %d, resp: %s", resp.StatusCode, string(b))
	}

	respMessage := parseMessageContent(string(b))

	return respMessage, nil
}

func parseMessageContent(jsonResp string) string {
	var cont string

	key := "content"
	for i := len(key); i < len(jsonResp); i++ {
		if jsonResp[i-len(key):i] == key {

			fmt.Println("find", jsonResp[i-len(key):i])

			var idx1, idx2 int

			for j := i + 1; j < len(jsonResp); j++ {
				if string(jsonResp[j]) == "\"" {
					if idx1 == 0 {
						idx1 = j + 1
					} else {
						idx2 = j
						break
					}
				}
			}

			cont = jsonResp[idx1:idx2]

			break
		}
	}
	return cont
}
