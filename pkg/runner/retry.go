package runner

import (
	"Yi/pkg/db"
	"Yi/pkg/logging"
	"sync"
	"time"
)

/**
  @author: yhy
  @since: 2023/2/16
  @desc: Database generates errors for retry
**/

var IsRetry bool

// Retry Todo is not elegant
func Retry() {
	var wg sync.WaitGroup
	limit := make(chan bool, Option.Thread)

	for {
		if IsRetry { // After the runtime, you will not need an coroutine, one by one, one by one
			for _, perr := range RetryProject {
				if !IsRetry {
					break
				}
				wg.Add(1)
				limit <- true
				delete(RetryProject, perr.Url)

				go func(p ProError) {
					defer func() {
						<-limit
						wg.Done()
					}()
					logging.Logger.Printf("Project (%S) Review", p.Url)

					_, project := db.Exist(p.Url)

					if p.Code == 1 {
						// Get from github
						_, dbPath, res := GetRepos(p.Url)
						project.DBPath = dbPath
						project.Language = res.Language
						project.PushedAt = res.PushedAt
						project.DefaultBranch = res.DefaultBranch
					} else if p.Code == 2 { // Manually
						project.DBPath = CreateDb(p.Url, project.Language)
					}

					db.UpdateProject(project.Id, project)
					if project.DBPath != "" {
						Exec(project, nil)
					}
				}(perr)

			}
		}
		time.Sleep(time.Minute)
	}
}
