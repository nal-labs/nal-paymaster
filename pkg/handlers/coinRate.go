package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// func GetCoinRate() {

// 	url := "https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd&include_last_updated_at=true"

// 	req, _ := http.NewRequest("GET", url, nil)

// 	req.Header.Add("accept", "application/json")
// 	req.Header.Add("x-cg-demo-api-key", "FujYAV")

// 	res, _ := http.DefaultClient.Do(req)

// 	defer res.Body.Close()
// 	body, _ := io.ReadAll(res.Body)

// 	fmt.Println(string(body))
// }

type CoinGeckoResponse struct {
	ETH struct {
		USD           float64 `json:"usd"`
		LastUpdatedAt int64   `json:"last_updated_at"`
	} `json:"ethereum"`
}

func GetCoinRate() int64 {

	url := "https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd&include_last_updated_at=true"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return 0
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("x-cg-demo-api-key", "CG-p2uWegSdhu2MZkLcRMFujYAV")

	// res, err := http.DefaultClient.Do(req)

	// tr := &http.Transport{
	// 	MaxIdleConns: 100,
	// }
	client := &http.Client{
		// Transport: tr,
		Timeout: 3 * time.Second, // 超时加在这里，是每次调用的超时
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return 0
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return 0
	}

	// fmt.Println(string(body))
	// 解析JSON响应
	var coinResponse CoinGeckoResponse
	err = json.Unmarshal(body, &coinResponse)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return 0
	}

	// 打印解析后的数据
	fmt.Printf("Ethereum price in USD: %.2f\n", coinResponse.ETH.USD)
	fmt.Printf("Last updated at: %d\n", coinResponse.ETH.LastUpdatedAt)
	return int64(coinResponse.ETH.USD * 1.2)
}
