package logic

import (
	"backend/dao/api"
	"backend/models"
)

// GetGithubTrending Fetches the Github trending projects
func GetGithubTrending(p *models.ParamGithubTrending) (data *models.GithubTrending, err error) {
	switch p.Language {
	case 0:
		data, err = api.GetGithubTrendingAll(p)
		//case 1:
		//	data, err = models.GetGithubTrendingGo(p.Since, p.Page, p.Size)
		//default:
		//	data, err = models.GetGithubTrendingAll(p.Since, p.Page, p.Size)
	}
	return
}
