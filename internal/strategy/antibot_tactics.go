package strategy

import (
	"crypto/rand"
	"math/big"
	"net/http"
)

type AntibotTactics struct {
	userAgents []string
	secChUa    []string
}

func NewAntibotTactics() *AntibotTactics {
	return &AntibotTactics{
		userAgents: []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		},
		secChUa: []string{
			`"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`,
			`"Google Chrome";v="123", "Not:A-Brand";v="8"`,
		},
	}
}

// InjectEvasionHeaders applies authentic browser metadata and header ordering
func (at *AntibotTactics) InjectEvasionHeaders(req *http.Request) {
	uaIndex := getRandomInt(len(at.userAgents))
	secIndex := getRandomInt(len(at.secChUa))
	req.Header.Set("User-Agent", at.userAgents[uaIndex])
	req.Header.Set("Sec-Ch-Ua", at.secChUa[secIndex])
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Accept", "audio/webm,audio/ogg,audio/wav,audio/*;q=0.9,application/json;q=0.8,*/*;q=0.5")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
}

func getRandomInt(max int) int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0
	}
	return int(nBig.Int64())
}
