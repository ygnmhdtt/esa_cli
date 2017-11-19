package esa_cli

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("TEST", "1")
	code := m.Run()
	os.Exit(code)
}

func TestGetTeam(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := `{
  "name": "docs",
  "privacy": "open",
  "description": "esa.io official documents",
  "icon": "https://img.esa.io/uploads/production/teams/105/icon/thumb_m_0537ab827c4b0c18b60af6cdd94f239c.png",
  "url": "https://docs.esa.io/"
}`
		fmt.Fprintf(w, res)
	}))
	os.Setenv("TEST_URL", ts.URL)
	client := NewClient("", "docs")
	team, _ := client.GetTeam()
	if team.Name != "docs" {
		t.Fatalf("invalid data: %v", team)
	}
	ts.Close()
}

func TestGetTeamStats(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := `{
  "members": 20,
  "posts": 1959,
  "posts_wip": 59,
  "posts_shipped": 1900,
  "comments": 2695,
  "stars": 3115,
  "daily_active_users": 8,
  "weekly_active_users": 14,
  "monthly_active_users": 15
}`
		fmt.Fprintf(w, res)
	}))
	os.Setenv("TEST_URL", ts.URL)
	client := NewClient("", "docs")
	teamStats, _ := client.GetTeamStats()
	if teamStats.Members != 20 {
		t.Fatalf("invalid data: %v", teamStats)
	}
	ts.Close()
}

func TestGetTeamMembers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := `{
  "members": [
    {
      "name": "Atsuo Fukaya",
      "screen_name": "fukayatsu",
      "icon": "https://img.esa.io/uploads/production/users/1/icon/thumb_m_402685a258cf2a33c1d6c13a89adec92.png",
      "email": "fukayatsu@esa.io",
      "posts_count": 222
    },
    {
      "name": "TAEKO AKATSUKA",
      "screen_name": "taea",
      "icon": "https://img.esa.io/uploads/production/users/2/icon/thumb_m_2690997f07b7de3014a36d90827603d6.jpg",
      "email": "taea@esa.io",
      "posts_count": 111
    }
  ],
  "prev_page": null,
  "next_page": null,
  "total_count": 2,
  "page": 1,
  "per_page": 20,
  "max_per_page": 100
}
`
		fmt.Fprintf(w, res)
	}))
	os.Setenv("TEST_URL", ts.URL)
	client := NewClient("", "docs")
	teamMembers, _ := client.GetTeamMembers(1)
	if teamMembers.Members[0].Name != "Atsuo Fukaya" {
		t.Fatalf("invalid data: %v", teamMembers)
	}
	ts.Close()
}

func TestGetPosts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := `{
  "posts": [
    {
      "number": 1,
      "name": "hi!",
      "full_name": "日報/2015/05/09/hi! #api #dev",
      "wip": true,
      "body_md": "# Getting Started",
      "body_html": "<h1 id=\"1-0-0\" name=\"1-0-0\">\n<a class=\"anchor\" href=\"#1-0-0\"><i class=\"fa fa-link\"></i><span class=\"hidden\" data-text=\"Getting Started\"> &gt; Getting Started</span></a>Getting Started</h1>\n",
      "created_at": "2015-05-09T11:54:50+09:00",
      "message": "Add Getting Started section",
      "url": "https://docs.esa.io/posts/1",
      "updated_at": "2015-05-09T11:54:51+09:00",
      "tags": [
        "api",
        "dev"
      ],
      "category": "日報/2015/05/09",
      "revision_number": 1,
      "created_by": {
        "name": "Atsuo Fukaya",
        "screen_name": "fukayatsu",
        "icon": "http://img.esa.io/uploads/production/users/1/icon/thumb_m_402685a258cf2a33c1d6c13a89adec92.png"
      },
      "updated_by": {
        "name": "Atsuo Fukaya",
        "screen_name": "fukayatsu",
        "icon": "http://img.esa.io/uploads/production/users/1/icon/thumb_m_402685a258cf2a33c1d6c13a89adec92.png"
      }
    }
  ],
  "prev_page": null,
  "next_page": null,
  "total_count": 1,
  "page": 1,
  "per_page": 20,
  "max_per_page": 100
}`
		fmt.Fprintf(w, res)
	}))
	os.Setenv("TEST_URL", ts.URL)
	client := NewClient("", "docs")
	posts, _ := client.GetPosts(1)
	if posts.Posts[0].Number != 1 {
		t.Fatalf("invalid data: %v", posts)
	}
	ts.Close()
}

func TestGetPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := `{
  "number": 1,
  "name": "hi!",
  "full_name": "日報/2015/05/09/hi! #api #dev",
  "wip": true,
  "body_md": "# Getting Started",
  "body_html": "<h1 id=\"1-0-0\" name=\"1-0-0\">\n<a class=\"anchor\" href=\"#1-0-0\"><i class=\"fa fa-link\"></i><span class=\"hidden\" data-text=\"Getting Started\"> &gt; Getting Started</span></a>Getting Started</h1>\n",
  "created_at": "2015-05-09T11:54:50+09:00",
  "message": "Add Getting Started section",
  "url": "https://docs.esa.io/posts/1",
  "updated_at": "2015-05-09T11:54:51+09:00",
  "tags": [
    "api",
    "dev"
  ],
  "category": "日報/2015/05/09",
  "revision_number": 1,
  "created_by": {
    "name": "Atsuo Fukaya",
    "screen_name": "fukayatsu",
    "icon": "http://img.esa.io/uploads/production/users/1/icon/thumb_m_402685a258cf2a33c1d6c13a89adec92.png"
  },
  "updated_by": {
    "name": "Atsuo Fukaya",
    "screen_name": "fukayatsu",
    "icon": "http://img.esa.io/uploads/production/users/1/icon/thumb_m_402685a258cf2a33c1d6c13a89adec92.png"
  },
  "kind": "flow",
  "comments_count": 1,
  "tasks_count": 1,
  "done_tasks_count": 1,
  "stargazers_count": 1,
  "watchers_count": 1,
  "star": true,
  "watch": true
}`
		fmt.Fprintf(w, res)
	}))
	os.Setenv("TEST_URL", ts.URL)
	client := NewClient("", "docs")
	post, _ := client.GetPost(1)
	if post.Number != 1 {
		t.Fatalf("invalid data: %v", post)
	}
	ts.Close()
}

func TestCreatePost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		res := `{
  "number": 5,
  "name": "hi!",
  "full_name": "dev/2015/05/10/hi! #api #dev",
  "wip": false,
  "body_md": "# Getting Started\n",
  "body_html": "<h1 id=\"1-0-0\" name=\"1-0-0\">\n<a class=\"anchor\" href=\"#1-0-0\"><i class=\"fa fa-link\"></i><span class=\"hidden\" data-text=\"Getting Started\"> &gt; Getting Started</span></a>Getting Started</h1>\n",
  "created_at": "2015-05-09T12:12:37+09:00",
  "message": "Add Getting Started section",
  "url": "https://docs.esa.io/posts/5",
  "updated_at": "2015-05-09T12:12:37+09:00",
  "tags": [
    "api",
    "dev"
  ],
  "category": "dev/2015/05/10",
  "revision_number": 1,
  "created_by": {
    "name": "Atsuo Fukaya",
    "screen_name": "fukayatsu",
    "icon": "http://img.esa.io/uploads/production/users/1/icon/thumb_m_402685a258cf2a33c1d6c13a89adec92.png"
  },
  "updated_by": {
    "name": "Atsuo Fukaya",
    "screen_name": "fukayatsu",
    "icon": "http://img.esa.io/uploads/production/users/1/icon/thumb_m_402685a258cf2a33c1d6c13a89adec92.png"
  },
  "kind": "flow",
  "comments_count": 0,
  "tasks_count": 0,
  "done_tasks_count": 0,
  "stargazers_count": 0,
  "watchers_count": 1,
  "star": false,
  "watch": false
}`
		fmt.Fprintf(w, res)
	}))
	os.Setenv("TEST_URL", ts.URL)
	client := NewClient("", "docs")

	// empty name
	post := PostCreate{}
	_, err := client.CreatePost(&post)
	if err == nil {
		t.Fatalf("err was nil")
	}

	// valid
	post = PostCreate{}
	post.Post.Name = "test"
	_, err = client.CreatePost(&post)
	if err != nil {
		t.Fatalf("err occured: %v", err)
	}
	ts.Close()
}
