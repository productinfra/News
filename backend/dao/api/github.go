package api

import (
	"backend/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

// GetGithubTrendingAll retrieves all GitHub trending repositories across all languages
func GetGithubTrendingAll(p *models.ParamGithubTrending) (data *models.GithubTrending, err error) {
	url := "https://api.github.com/search/repositories?q=stars:%253E1&sort=stars&order=desc&page=" +
		fmt.Sprintf("%d", p.Page) + "&per_page=" + fmt.Sprintf("%d", p.Size)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		zap.L().Error("http.NewRequest failed", zap.Error(err))
		return
	}
	res, err := client.Do(req)
	if err != nil {
		zap.L().Error("client.Do failed", zap.Error(err))
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("ioutil.ReadAll failed", zap.Error(err))
		return
	}
	fmt.Println(string(body))
	var githubTrendingAll models.GithubTrending
	err = json.Unmarshal(body, &githubTrendingAll)
	if err != nil {
		zap.L().Error("json.Unmarshal failed", zap.Error(err))
		return
	}
	return &githubTrendingAll, nil
}
