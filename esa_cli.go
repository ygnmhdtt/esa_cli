package esa_cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

type Client_V1 struct {
	URL           *url.URL
	TeamName      string
	HTTPClient    *http.Client
	Authorization string
	Logger        *log.Logger
}

type Team struct {
	Name        string `json:"name"`
	Privacy     string `json:"privacy"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	URL         string `json:"url"`
}

type TeamStats struct {
	Members            int `json:"members"`
	Posts              int `json:"posts"`
	PostsWip           int `json:"posts_wip"`
	PostsShipped       int `json:"posts_shipped"`
	Comments           int `json:"comments"`
	Stars              int `json:"stars"`
	DailyActiveUsers   int `json:"daily_active_users"`
	WeeklyActiveUsers  int `json:"weekly_active_users"`
	MonthlyActiveUsers int `json:"monthly_active_users"`
}

type TeamMembers struct {
	Members []struct {
		Name       string `json:"name"`
		ScreenName string `json:"screen_name"`
		Icon       string `json:"icon"`
		Email      string `json:"email"`
		PostsCount int    `json:"posts_count"`
	} `json:"members"`
	PrevPage   interface{} `json:"prev_page"`
	NextPage   interface{} `json:"next_page"`
	TotalCount int         `json:"total_count"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	MaxPerPage int         `json:"max_per_page"`
}

type CreatedBy struct {
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
	Icon       string `json:"icon"`
}

type UpdatedBy struct {
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
	Icon       string `json:"icon"`
}

type PostGet struct {
	Number          int       `json:"number"`
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Wip             bool      `json:"wip"`
	BodyMd          string    `json:"body_md"`
	BodyHTML        string    `json:"body_html"`
	CreatedAt       time.Time `json:"created_at"`
	Message         string    `json:"message"`
	URL             string    `json:"url"`
	UpdatedAt       time.Time `json:"updated_at"`
	Tags            []string  `json:"tags"`
	Category        string    `json:"category"`
	RevisionNumber  int       `json:"revision_number"`
	CreatedBy       CreatedBy `json:"created_by"`
	UpdatedBy       UpdatedBy `json:"updated_by"`
	Kind            string    `json:"kind"`
	CommentsCount   int       `json:"comments_count"`
	TasksCount      int       `json:"tasks_count"`
	DoneTasksCount  int       `json:"done_tasks_count"`
	StargazersCount int       `json:"stargazers_count"`
	WatchersCount   int       `json:"watchers_count"`
	Star            bool      `json:"star"`
	Watch           bool      `json:"watch"`
}

type Posts struct {
	Posts      []PostGet   `json:"posts"`
	PrevPage   interface{} `json:"prev_page"`
	NextPage   int         `json:"next_page"`
	TotalCount int         `json:"total_count"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	MaxPerPage int         `json:"max_per_page"`
}

type PostCreate struct {
	Post struct {
		Name     string   `json:"name"`
		BodyMd   string   `json:"body_md"`
		Tags     []string `json:"tags"`
		Category string   `json:"category"`
		Wip      bool     `json:"wip"`
		Message  string   `json:"message"`
	} `json:"post"`
}

type Categories struct {
	Categories []Category
}

type Category struct {
	Name     string     `json:"name"`
	Post     bool       `json:"post,omitempty"`
	Count    int        `json:"count,omitempty"`
	Children []Category `json:"children,omitempty"`
}

func (c Category) Tree() []string {
	if c.Post {
		return nil
	}
	var trees []string

	if len(c.Children) == 0 {
		trees = append(trees, c.Name+"/")
		return trees
	}
	for _, child := range c.Children {
		trees = append(trees, c.Append(child.Tree())...)
	}
	return trees
}

func (c Category) Append(trees []string) []string {
	var res []string
	for _, tree := range trees {
		res = append(res, c.Name+"/"+tree)
	}
	return res
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

func NewClient(auth string, teamName string) *Client_V1 {
	client := new(Client_V1)
	if os.Getenv("TEST") == "1" {
		u, _ := url.Parse(os.Getenv("TEST_URL"))
		client.URL = u
	} else {
		u, _ := url.Parse("https://api.esa.io/v1")
		client.URL = u
	}
	client.TeamName = teamName
	client.HTTPClient = &http.Client{}
	client.Authorization = auth
	client.Logger = log.New(os.Stdout, "logger: ", log.Lshortfile)
	return client
}

func (c *Client_V1) newRequest(method string, spath string, body io.Reader) (*http.Request, error) {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.Authorization)
	return req, nil
}

func (c *Client_V1) GetTeam() (*Team, error) {
	spath := fmt.Sprintf("/teams/%v", c.TeamName)
	req, _ := c.newRequest("GET", spath, nil)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var team Team
	if err := decodeBody(res, &team); err != nil {
		return nil, err
	}
	return &team, err
}

func (c *Client_V1) GetTeamStats() (*TeamStats, error) {
	spath := fmt.Sprintf("/teams/%v/stats", c.TeamName)
	req, _ := c.newRequest("GET", spath, nil)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var teamStats TeamStats
	if err := decodeBody(res, &teamStats); err != nil {
		return nil, err
	}
	return &teamStats, err
}

func (c *Client_V1) GetTeamMembers(page int) (*TeamMembers, error) {
	spath := fmt.Sprintf("/teams/%v/members?page=%v", c.TeamName, page)
	req, _ := c.newRequest("GET", spath, nil)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var teamMembers TeamMembers
	if err := decodeBody(res, &teamMembers); err != nil {
		return nil, err
	}
	return &teamMembers, err
}

func (c *Client_V1) GetPosts(page int, q ...string) (*Posts, error) {
	spath := fmt.Sprintf("/teams/%v/posts?q=%v&page=%v", c.TeamName, q, page)
	req, _ := c.newRequest("GET", spath, nil)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var posts Posts
	if err := decodeBody(res, &posts); err != nil {
		return nil, err
	}
	return &posts, err
}

func (c *Client_V1) GetPost(id int) (*PostGet, error) {
	spath := fmt.Sprintf("/teams/%v/posts/%v", c.TeamName, id)
	req, _ := c.newRequest("GET", spath, nil)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var post PostGet
	if err := decodeBody(res, &post); err != nil {
		return nil, err
	}
	return &post, err
}

func (c *Client_V1) CreatePost(post *PostCreate) (*PostGet, error) {
	if post.Post.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	spath := fmt.Sprintf("/teams/$v/posts", c.TeamName)
	j, _ := json.Marshal(&post)
	req, _ := c.newRequest("POST", spath, bytes.NewBuffer([]byte(j)))
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		b, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf(string(b))
	}
	var cpost PostGet
	if err := decodeBody(res, &cpost); err != nil {
		return nil, err
	}
	return &cpost, nil
}

func (c *Client_V1) GetCategories() (*Categories, error) {
	spath := fmt.Sprintf("/teams/%v/categories", c.TeamName)
	req, _ := c.newRequest("GET", spath, nil)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var categories Categories
	if err := decodeBody(res, &categories); err != nil {
		return nil, err
	}
	return &categories, nil
}
