package lib

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func CloudFlare(APIKey string) {
	url := "https://api.cloudflare.com/client/v4"
	req, err := http.NewRequest(http.MethodGet, url+"/zones?per_page=50", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+APIKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", body)
}
