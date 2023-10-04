package runner

import (
	"Yi/pkg/db"
	"Yi/pkg/logging"
	"Yi/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/corpix/uarand"
	jsoniter "github.com/json-iterator/go"
	"github.com/thoas/go-funk"
)

/**
  @author: yhy
  @since: 2022/12/7
  @desc: //TODO
**/

type GithubRes struct {
	Language      string `json:"language"`
	PushedAt      string `json:"pushed_at"`
	DefaultBranch string `json:"default_branch"`
}

// ProError Retry if the project database is fetched incorrectly.
type ProError struct {
	Url  string
	Code int
}

var RetryProject = make(map[string]ProError)

// GetRepos from github Download and build a good database
func GetRepos(url_tmp string) (error, string, GithubRes) {
	// https://github.com/prometheus/prometheus  -> https://api.github.com/repos/prometheus/prometheus
	guri := strings.ReplaceAll(url_tmp, "github.com", "api.github.com/repos")

	res := GetTimeBran(guri, url_tmp)
	// https://api.github.com/repos/grafana/grafana Here will only display the most language in the project, but it is not necessarily the main language of the project. For example

	var flag bool
	// repos The language in the middle is only the type of language with the most proportion. It is possible that the TypeScript written by GoScript is the most proportion.
	res.Language, flag = GetLanguage(guri, url_tmp) // todo Now I just adapt to GO, Java language, and try to adapt to the mainstream language in the later period. At present

	logging.Logger.Debugln(url_tmp, " language: ", res.Language)
	if flag {
		guri = fmt.Sprintf("%s/code-scanning/codeql/databases/%s", guri, res.Language)
	} else {
		return errors.New("language does not support"), "", res
	}

	err, dbPath, code := GetDb(guri, url_tmp, res.Language)
	logging.Logger.Debugln(url_tmp, " dbPath: ", dbPath)
	if code != 0 { // No corresponding database is generated
		logging.Logger.Debugln(url_tmp, " GetDb err: ", err)
		RetryProject[url_tmp] = ProError{
			Url:  url_tmp,
			Code: code,
		}
	}
	return err, dbPath, res
}

// Gettimebran gets the project update time and main branch https://api.github.com/repos/prometheus/prometheus
func GetTimeBran(guri, url_tmp string) GithubRes {
	req, _ := http.NewRequest("GET", guri, nil)
	req.Header.Set("Accept", "application/vnd.github.v3.text-match+json")
	if Option.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Option.Token))
	}
	req.Header.Set("User-Agent", uarand.GetRandom())
	res := GithubRes{}

	req.Close = true
	Option.Session.RateLimiter.Take()
	resp, err := Option.Session.Client.Do(req)

	if err != nil {
		logging.Logger.Errorln("GetRepos client.Do(req) err:", err)
		// Caused by network errors, you need to try it out
		RetryProject[url_tmp] = ProError{
			Url:  url_tmp,
			Code: 1,
		}
		res.Language = ""
		return res
	}
	defer resp.Body.Close()

	if resp.Body != nil {
		result, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(result, &res)
	}
	return res
}

// GetLanguage Obtain the code language of the project  https://api.github.com/repos/prometheus/prometheus/languages
func GetLanguage(guri, url_tmp string) (string, bool) {
	req, _ := http.NewRequest("GET", guri+"/languages", nil)
	req.Header.Set("Accept", "application/vnd.github.v3.text-match+json")
	if Option.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Option.Token))
	}
	req.Header.Set("User-Agent", uarand.GetRandom())
	req.Close = true
	Option.Session.RateLimiter.Take()
	resp, err := Option.Session.Client.Do(req)

	if err != nil {
		logging.Logger.Errorln("GetLanguage client.Do(req) err:", err)
		// Caused by network errors, you need to try it out
		RetryProject[url_tmp] = ProError{
			Url:  url_tmp,
			Code: 1,
		}
		return "", false
	}

	defer resp.Body.Close()
	var language string
	if resp.Body != nil {
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return "", false
		}

		var m float64 = 1

		// Remove HTML, TypeScript, JavaScript, CSS, SCSS, and later, the language of the project is used to use the most useful language
		for k, v := range jsoniter.Get(body).GetInterface().(map[string]interface{}) {
			if funk.Contains(k, "HTML") || funk.Contains(k, "TypeScript") || funk.Contains(k, "JavaScript") || funk.Contains(k, "CSS") || funk.Contains(k, "SCSS") {
				continue
			}

			switch v.(type) {
			case float64:
				if v.(float64) > m {
					language = k
					m = v.(float64)
				}
			default:
				language = k
				logging.Logger.Errorf("GetLanguage err %s %s", guri, body)
			}
		}
		if funk.Contains(Languages, language) {
			return language, true
		}
	}

	return language, false
}

// Getdb download/generate database https://api.github.com/repos/prometheus/prometheus/code-sCanning/Codeql/databases/ {languages}
/*
0: Success
1: Network or file creation error
2: No corresponding database is generated on github
*/
func GetDb(guri, url, languages string) (error, string, int) {
	req, _ := http.NewRequest("GET", guri, nil)
	req.Header.Set("Accept", "application/zip")
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Content-Type", "application/octet-stream")
	if Option.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Option.Token))
	}
	req.Header.Set("User-Agent", uarand.GetRandom())
	req.Close = true

	Option.Session.RateLimiter.Take()
	resp, err := Option.Session.Client.Do(req)

	if err != nil {
		logging.Logger.Errorln(guri, "HttpRequest Do err: ", err)
		return err, "", 1
	}
	defer resp.Body.Close()

	name := utils.GetName(url)
	filePath := DirNames.ZipDir + name + ".zip"

	var dbPath string
	if resp != nil && resp.StatusCode == 200 {
		out, err := os.Create(filePath)
		defer out.Close()
		if err != nil {
			logging.Logger.Errorln("os.Create(filePath) err:", err)
			return err, "", 1
		}

		if _, err = io.Copy(out, resp.Body); err != nil {
			logging.Logger.Errorln(url, " HttpRequest io.Copy err: ", err)
			return err, "", 1
		}
	} else { // Explain that the project does not configure CodeQL scan (404) in GitHub, or the project owner is configured with access required (403)
		dbPath = CreateDb(url, languages)
		if dbPath == "" {
			return err, "", 2
		} else {
			return nil, dbPath, 0
		}
	}

	err = utils.DeCompress(filePath, DirNames.DbDir+name+"/")
	if err != nil {
		logging.Logger.Errorln("DeCompress err:", err)
		return err, "", 1
	}

	dbPath = filepath.Dir(path.Join(utils.CodeqlDb(DirNames.DbDir+name), "*"))
	logging.Logger.Debugln(url, " downloadDb success.")
	return nil, dbPath, 0
}

// CheckUpdate Check whether the item is updated
func CheckUpdate(project db.Project) (bool, string, string) {
	guri := strings.ReplaceAll(project.Url, "github.com", "api.github.com/repos")

	res := GetTimeBran(guri, project.Url)

	var (
		dbPath string
		code   int
	)

	if project.PushedAt < res.PushedAt { // Explain the update

		guri = fmt.Sprintf("%s/code-scanning/codeql/databases/%s", guri, project.Language)

		_, dbPath, code = GetDb(guri, project.Url, project.Language)

		if code != 0 { // No corresponding database is generated
			RetryProject[project.Url] = ProError{
				Url:  project.Url,
				Code: code,
			}
		} else {
			// Are you in the list of reviews?
			delete(RetryProject, project.Url)
		}
	}

	if dbPath != "" {
		logging.Logger.Debugln(project.Url, " update, start a new scan.", dbPath)
		record := db.Record{
			Project: project.Project,
			Url:     project.Url,
			Color:   "warning",
			Title:   project.Project + " renew",
			Msg:     fmt.Sprintf("%s Project update, Re -generate Codeql database", project.Url),
		}
		db.AddRecord(record)

		return true, dbPath, res.PushedAt
	}

	return false, "", ""
}
