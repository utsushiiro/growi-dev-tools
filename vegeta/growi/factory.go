package growi

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	vegeta "github.com/tsenart/vegeta/lib"
)

// PagesListResponse is a response structure of pages.list api.
type PagesListResponse struct {
	Pages []Page `json:"pages"`
}

// Page is a partial struct of PagesListResponse.
type Page struct {
	ID   string `json:"_id"`
	Path string `json:"path"`
}

// PagesGetResponse is a response structure of pages.get api.
type PagesGetResponse struct {
	Page PageDetail `json:"page"`
}

// PageDetail is a partial struct of PagesGetResponse.
type PageDetail struct {
	ID       string `json:"_id"`
	Path     string `json:"path"`
	Revision struct {
		ID   string `json:"_id"`
		Body string `json:"body"`
	} `json:"revision"`
}

// PagesUpdateRequest is a request structure of pages.update api.
type PagesUpdateRequest struct {
	PageID     string `json:"page_id"`
	RevisionID string `json:"revision_id"`
	Body       string `json:"body"`
}

// TargeterFactory : TargeterFactory
type TargeterFactory interface {
	NewRandomPageAccessTargeter() (vegeta.Targeter, error)
	NewRandomPageUpdateTargeter() (vegeta.Targeter, error)
}

type targeterFactory struct {
	host   string
	token  string
	cookie http.Cookie
}

// NewGrowiTargeterFactory creates growi.TargeterFactory instance.
func NewGrowiTargeterFactory() (TargeterFactory, error) {
	// set up global rand seed
	seed, err := strconv.ParseInt(os.Getenv("SEED"), 10, 64)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	rand.Seed(seed)

	return &targeterFactory{
		host:  os.Getenv("GROWI_URL"),
		token: os.Getenv("ADMIN_API_TOKEN"),
		cookie: http.Cookie{
			Name:  "connect.sid",
			Value: os.Getenv("ADMIN_SESSION_COOKIE"),
		},
	}, nil
}

// NewRandomPageAccessTargeter creates vegeta.Targeter that randomly access a page.
func (tf *targeterFactory) NewRandomPageAccessTargeter() (vegeta.Targeter, error) {
	pageList, err := tf.getPageList()
	if err != nil {
		return nil, err
	}

	return func(t *vegeta.Target) error {
		t.Method = "GET"
		t.URL = tf.host + pageList[rand.Intn(len(pageList))].Path
		t.Header = http.Header{}
		t.Header.Add("Cookie", tf.cookie.String())
		return nil
	}, nil
}

// NewRandomPageUpdateTargeter creates vegeta.Targeter that randomly get and update a page.
// Sometimes, this fails because of revision_id conflict
func (tf *targeterFactory) NewRandomPageUpdateTargeter() (vegeta.Targeter, error) {
	pageList, err := tf.getPageList()
	if err != nil {
		return nil, err
	}

	return func(t *vegeta.Target) error {
		targetPage := pageList[rand.Intn(len(pageList))]

		// get current page info for the latest revision_id
		currentTargetPage, err := tf.getPage(targetPage.ID)
		if err != nil {
			return nil
		}

		// update page
		t.Method = "POST"
		t.URL = fmt.Sprintf("%s/_api/pages.update?access_token=%s", tf.host, tf.token)
		t.Header = http.Header{}
		t.Header.Add("Content-Type", "application/json")
		t.Body, err = json.Marshal(&PagesUpdateRequest{
			PageID:     targetPage.ID,
			RevisionID: currentTargetPage.Revision.ID,
			Body:       currentTargetPage.Revision.Body + "\nUpdate!",
		})
		return nil
	}, nil
}

func (tf *targeterFactory) getPageList() ([]Page, error) {
	endpoint := fmt.Sprintf("%s/_api/pages.list?path=/&limit=1000&access_token=%s", tf.host, tf.token)
	response, err := http.Get(endpoint)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var pagesListResponse PagesListResponse
	if err := json.NewDecoder(response.Body).Decode(&pagesListResponse); err != nil {
		return nil, err
	}

	return pagesListResponse.Pages, nil
}

func (tf *targeterFactory) getPage(pageID string) (*PageDetail, error) {
	endpoint := fmt.Sprintf("%s/_api/pages.get?page_id=%s&access_token=%s", tf.host, pageID, tf.token)
	response, err := http.Get(endpoint)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var pagesGetResponse PagesGetResponse
	if err := json.NewDecoder(response.Body).Decode(&pagesGetResponse); err != nil {
		return nil, err
	}

	return &pagesGetResponse.Page, nil
}
