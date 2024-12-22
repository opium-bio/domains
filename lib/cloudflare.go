package lib

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/opium-bio/config"
	"github.com/opium-bio/utils"
)

type CloudflareResponse struct {
	Success bool `json:"success"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

var cfg *config.Config

func init() {
	var err error
	cfg, err = config.LoadConfig("./config.toml")
	if err != nil {
		utils.Error("Error loading config", true)
	}
}

func CloudFlare() {
	url := "https://api.cloudflare.com/client/v4"
	req, err := http.NewRequest(http.MethodGet, url+"/zones?per_page=50", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Cloudflare.CFApiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	utils.Log("Successfully connected to Cloudflare!")
}

func AddDomain(domain string) error {
	url := "https://api.cloudflare.com/client/v4/zones"
	var str = []byte(`{
        "account": {
            "id": "` + cfg.Cloudflare.CFAccount + `"
        },
        "name": "` + domain + `",
        "type": "full"
    }`)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(str))
	if err != nil {
		utils.Error("Error creating Cloudflare request: "+err.Error(), true)
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Cloudflare.CFApiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.Error("Error executing Cloudflare request: "+err.Error(), true)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Error("Error reading Cloudflare response: "+err.Error(), true)
		return err
	}

	var cfResponse CloudflareResponse
	err = json.Unmarshal(body, &cfResponse)
	if err != nil {
		utils.Error("Error unmarshaling Cloudflare response: "+err.Error(), true)
		return err
	}

	if !cfResponse.Success {
		if len(cfResponse.Errors) > 0 {
			utils.Error("Cloudflare API error: "+cfResponse.Errors[0].Message, true)
		} else {
			utils.Error("Cloudflare API error: Unknown error occurred", true)
		}
		utils.Error("Cloudflare API error occurred", true)
	}

	utils.Log("Successfully added domain: " + domain)
	return nil
}
