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
