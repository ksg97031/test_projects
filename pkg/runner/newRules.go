package runner

import (
	"Yi/pkg/db"
	"Yi/pkg/utils"
	"sync"
)

/**
  @author: yhy
  @since: 2022/12/13
  @desc: //TODO
**/

// NewRules When the rules are updated, get the project from the database and run it with new rules
func NewRules(oldQls *QLFile, newQls *QLFile) {
	pythonQLs := utils.Difference(oldQls.PythonQL, newQls.PythonQL)

	globalDBTmp := db.GlobalDB.Model(&db.Project{})

	if len(pythonQLs) != 0 {
		var projects []db.Project
		globalDBTmp.Where("language = Python").Order("id asc").Find(&projects)
		scan(projects, pythonQLs)
	}
}

func scan(projects []db.Project, qls []string) {
	var wg sync.WaitGroup
	limit := make(chan bool, Option.Thread)

	for _, project := range projects {
		if project.DBPath == "" {
			continue
		}
		wg.Add(1)
		limit <- true
		Exec(project, qls)
		<-limit
		wg.Done()
	}

	wg.Wait()
	close(limit)
}
