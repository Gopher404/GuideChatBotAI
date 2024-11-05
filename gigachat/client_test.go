package gigachat

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"log"
	"main/config"
	"main/domain"
	"os"
	"path/filepath"
	"testing"
)

func getConfig() *config.Config {
	currentFileDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	projectDir, _ := filepath.Split(currentFileDir)

	cfg, err := config.Read(filepath.Join(projectDir, "config", "config.json"))
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

func TestRequest(t *testing.T) {

	messages := []domain.AIMessage{
		{Role: "user", Content: "Привет"},
		{Role: "assistant", Content: "Привет!"},
		{Role: "user", Content: "Отправь любой json"},
	}

	cfg := getConfig()

	c, err := NewClient(&cfg.GigaChat)
	require.NoError(t, err)

	resp, err := c.Request(messages, 0.4)
	require.NoError(t, err)
	fmt.Println(resp)

}
