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
	"main/domain"
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

type Request struct {
	Model             string             `json:"model"`
	Messages          []domain.AIMessage `json:"messages"`
	Stream            bool               `json:"stream"`
	RepetitionPenalty float32            `json:"repetition_penalty"`
	Temperature       float32            `json:"temperature"`
	TopP              float32            `json:"top_p"`
}

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (c *Client) Request(messages []domain.AIMessage, temperature float32) (string, error) {
	AIRequest := Request{
		Model:             "GigaChat-Pro",
		Messages:          messages,
		Stream:            false,
		RepetitionPenalty: 1,
		Temperature:       temperature,
		TopP:              0.4,
	}

	jsonAIRequest, err := json.Marshal(AIRequest)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://gigachat.devices.sberbank.ru/api/v1/chat/completions", bytes.NewBuffer(jsonAIRequest))

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

	respStruct := new(Response)

	if err := json.Unmarshal(b, respStruct); err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status code %d, error: %s", resp.StatusCode, respStruct.Message)
	}

	if len(respStruct.Choices) == 0 {
		return "", errors.New("no choices found")
	}

	return respStruct.Choices[0].Message.Content, nil
}

func parseMessageContent(jsonResp string) string {
	var cont string

	key := "content"
	for i := len(key); i < len(jsonResp); i++ {
		if jsonResp[i-len(key):i] == key {

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
