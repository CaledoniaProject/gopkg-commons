package commons

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type GithubRepo struct {
	Description string `json:"description"`
	Stars       int    `json:"stars"`
}

func GetGithubRepo(user, repo string) (*GithubRepo, error) {
	var (
		githubRepo = &GithubRepo{}
	)

	if _, body, err := HttpRequest(&RequestOptions{
		URL: fmt.Sprintf("https://api.github.com/repos/%s/%s", url.PathEscape(user), url.PathEscape(repo)),
		Headers: map[string]string{
			"User-Agent": RandomUserAgent(),
			"Connection": "close",
		},
	}); err != nil {
		return nil, err
	} else if err := json.Unmarshal(body, githubRepo); err != nil {
		return nil, err
	} else {
		return githubRepo, nil
	}
}
