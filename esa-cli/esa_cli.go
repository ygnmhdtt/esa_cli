package esa_cli

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

type Client struct {
	URLV1         *url.URL
	HTTPClient    *http.Client
	Authorization string
	Logger        *log.Logger
}

type Post struct {
	Number         int       `json:"number"`
	Name           string    `json:"name"`
	FullName       string    `json:"full_name"`
	Wip            bool      `json:"wip"`
	BodyMd         string    `json:"body_md"`
	BodyHTML       string    `json:"body_html"`
	CreatedAt      time.Time `json:"created_at"`
	Message        string    `json:"message"`
	URL            string    `json:"url"`
	UpdatedAt      time.Time `json:"updated_at"`
	Tags           []string  `json:"tags"`
	Category       string    `json:"category"`
	RevisionNumber int       `json:"revision_number"`
	CreatedBy      struct {
		Name       string `json:"name"`
		ScreenName string `json:"screen_name"`
		Icon       string `json:"icon"`
	} `json:"created_by"`
	UpdatedBy struct {
		Name       string `json:"name"`
		ScreenName string `json:"screen_name"`
		Icon       string `json:"icon"`
	} `json:"updated_by"`
	Kind            string `json:"kind"`
	CommentsCount   int    `json:"comments_count"`
	TasksCount      int    `json:"tasks_count"`
	DoneTasksCount  int    `json:"done_tasks_count"`
	StargazersCount int    `json:"stargazers_count"`
	WatchersCount   int    `json:"watchers_count"`
	Star            bool   `json:"star"`
	Watch           bool   `json:"watch"`
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

func NewClient(auth string) *Client {
	client := new(Client)
	u, _ := url.Parse("https://api.esa.io/v1")
	client.URLV1 = u
	client.HTTPClient = &http.Client{}
	client.Authorization = auth
	client.Logger = log.New(os.Stdout, "logger: ", log.Lshortfile)
	return client
}

func (c *Client) newRequest(method string, spath string, body io.Reader) (*http.Request, error) {
	u := *c.URLV1
	u.Path = path.Join(c.URLV1.Path, spath)
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.Authorization)
	return req, nil
}

func (c *Client) GetPost(id string) (*Post, error) {
	spath := fmt.Sprintf("/teams/mmmcorp/posts/%v", id)
	req, _ := c.newRequest("GET", spath, nil)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var post Post
	if err := decodeBody(res, &post); err != nil {
		return nil, err
	}
	return &post, err
}
