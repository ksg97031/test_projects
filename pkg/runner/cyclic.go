package runner

import (
	"Yi/pkg/db"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/thoas/go-funk"
)

/**
  @author: yhy
  @since: 2022/12/7
  @desc: loop execution
**/

func Cyclic() {
	for {
		// todo is not elegant enough, in case there are too many items to monitor and the day is not over yet.
		// Wait 24 hours before cycling
		if !Option.RunNow {
			time.Sleep(24 * 60 * time.Minute)
			// Stop retrying the check when you start checking for updates.
			IsRetry = false
		}

		Option.RunNow = false

		// Updating the rule base
		UpdateRule()

		count := 0
		today := time.Now().Format("2006-01-02") + "/"
		DirNames = DirName{
			ZipDir:    Pwd + "/db/zip/" + today,
			ResDir:    Pwd + "/db/results/" + today,
			DbDir:     Pwd + "/db/database/" + today,
			GithubDir: Pwd + "/github/" + today,
		}
		os.MkdirAll(DirNames.ZipDir, 0755)
		os.MkdirAll(DirNames.ResDir, 0755)
		os.MkdirAll(DirNames.DbDir, 0755)
		os.MkdirAll(DirNames.GithubDir, 0755)

		var projects []db.Project
		globalDBTmp := db.GlobalDB.Model(&db.Project{})
		globalDBTmp.Order("id asc").Find(&projects)

		var wg sync.WaitGroup
		limit := make(chan bool, Option.Thread)

		for _, p := range projects {
			if p.DBPath == "" || !funk.Contains(Languages, p.Language) {
				continue
			}
			wg.Add(1)
			limit <- true
			go func(project db.Project) {
				defer func() {
					<-limit
					wg.Done()
				}()

				// Explanation The previous run failed, try again.
				if project.Count == 0 {
					Exec(project, nil)
				} else {
					// It's only when it's updated that it goes and generates a database
					update, dbPath, pushedAt := CheckUpdate(project)

					if !update {
						return
					}

					count++
					project.DBPath = dbPath
					project.PushedAt = pushedAt

					db.UpdateProject(project.Id, project)
					Exec(project, nil)
				}
			}(p)

		}

		wg.Wait()
		close(limit)

		// After running them all, start retrying the items that went wrong
		IsRetry = true

		record := db.Record{
			Color: "primary",
			Title: "A new round of scanning",
			Msg:   fmt.Sprintf("A new round of scanning has been completed, with a total of %d items scanned.", count),
		}
		db.AddRecord(record)
	}

}
