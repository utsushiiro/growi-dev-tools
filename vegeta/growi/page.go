package growi

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	vegeta "github.com/tsenart/vegeta/lib"
)

// PageList : PageList
type PageList struct {
	Pages []struct {
		Path string `json:"path"`
	} `json:"pages"`
}

// NewRandomPageAccessTargeter : NewRandomPageAccessTargeter
func NewRandomPageAccessTargeter() (vegeta.Targeter, error) {
	host := os.Getenv("GROWI_URL")
	cookie := http.Cookie{
		Name:  "connect.sid",
		Value: os.Getenv("ADMIN_SESSION_COOKIE"),
	}

	pageList, err := getPageList()
	if err != nil {
		return nil, err
	}

	seed, err := strconv.ParseInt(os.Getenv("SEED"), 10, 64)
	if err != nil {
		return nil, err
	}
	rand.Seed(seed)

	return func(t *vegeta.Target) error {
		t.Method = "GET"
		t.URL = host + pageList.Pages[rand.Intn(len(pageList.Pages))].Path
		t.Header = http.Header{}
		t.Header.Add("Cookie", cookie.String())
		return nil
	}, nil
}

func getPageList() (*PageList, error) {
	host := os.Getenv("GROWI_URL")
	token := os.Getenv("ADMIN_API_TOKEN")
	endpoint := fmt.Sprintf("%s/_api/pages.list?path=/&limit=1000&access_token=%s", host, token)

	response, err := http.Get(endpoint)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var pageList PageList
	if err := json.NewDecoder(response.Body).Decode(&pageList); err != nil {
		return nil, err
	}

	return &pageList, nil
}
