package runner

import (
	"Yi/pkg/db"
	"Yi/pkg/logging"
	"Yi/pkg/utils"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

/**
  @author: yhy
  @since: 2022/10/13
  @desc: //TODO
**/

func Analyze(database string, name string, language string, qls []string) map[string]string {
	if language == "Python" {
		qls = QLFiles.PythonQL
	}

	if len(qls) == 0 {
		logging.Logger.Debugln("len(qls) = 0")
		return nil
	}

	res := make(map[string]string)
	filePath := DirNames.ResDir + name
	os.MkdirAll(filePath, 0755)

	logging.Logger.Infof("[[%s:%s]] analyze start ...", name, database)
	for i, ql := range qls {
		fileName := fmt.Sprintf("%s/%d.json", filePath, time.Now().Unix())
		cmd := exec.Command("codeql", "database", "analyze", "--rerun", database, Option.Path+ql, "--format=sarif-latest", "-o", fileName)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout // standard output
		cmd.Stderr = &stderr // standard error
		err := cmd.Run()
		_, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		if err != nil {
			logging.Logger.Errorf("Analyze cmd.Run() failed with %s --  %s, %s %s", err, errStr, database, name)
			continue
		}

		lines := utils.LoadFile(fileName)

		if len(lines) == 0 {
			continue
		}

		var result string

		for _, line := range lines {
			result += line
		}
		res[fileName] = result

		ProgressBar[name] = float32(i+1) / float32(len(qls)) * 100
	}

	logging.Logger.Infof("[[%s:%s]] analysis completed.", name, database)
	record := db.Record{
		Project: name,
		Url:     name,
		Color:   "success",
		Title:   name,
		Msg:     fmt.Sprintf("%s end of analysis", name),
	}
	ProgressBar[name] = 100
	db.AddRecord(record)
	return res
}

// CreateDb pulls repository, creates database locally
func CreateDb(gurl, languages string) string {
	dbName := utils.GetName(gurl)
	err := GitClone(gurl, dbName)

	if err != nil {
		logging.Logger.Errorln("create db err:", err)
		return ""
	}

	// The todo batch run just jerked around, causing some projects to fail to generate a database "There's no CodeQL extractor named 'Go' installed."
	cmd := exec.Command("codeql", "database", "create", DirNames.DbDir+dbName, "-s", DirNames.GithubDir+dbName, "--language="+strings.ToLower(languages), "--overwrite")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // standard output
	cmd.Stderr = &stderr // standard error
	err = cmd.Run()
	out, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if err != nil {
		logging.Logger.Errorf("CreateDb cmd.Run() failed with %s\n %s --  %s\n", err, out, errStr)
		return ""
	}

	// It's strange that some of the generated databases are not in the project directory, but in the second level directory
	dbPath := filepath.Dir(path.Join(utils.CodeqlDb(DirNames.DbDir+dbName), "*"))
	logging.Logger.Debugln(gurl, " CreateDb success")
	return dbPath
}

// UpdateRule pulls the official repository every day to update the rules
func UpdateRule() {
	if Option.Path != "" {
		_, err := utils.RunGitCommand(Option.Path, "git", "pull")
		record := db.Record{
			Project: "CodeQL Rules",
			Url:     "CodeQL Rules",
			Color:   "success",
			Title:   "CodeQL Rules",
			Msg:     "CodeQL Rules Successful update",
		}

		if err != nil {
			record.Color = "danger"
			record.Msg = fmt.Sprintf("CodeQL Rules update failure, %s", err.Error())
		}

		db.AddRecord(record)
	}
}
