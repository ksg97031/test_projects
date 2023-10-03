# 弈 (the game of go)
以有算无 (To have calculation is to have everything; without it, you have nothing)

> [lgtm](https://lgtm.com/) is closing down, so I made my own monitoring tool. It's also easy to automate batch scanning after writing your own rules, so you can pick up holes efficiently.

Every day, check whether the github project is updated, automatically get/generate database query, automatically run CodeQL rule query, efficiently pick up holes.

Default web page open on port 8888, username, password if not specified, the default username is yhy, the password is random, will be output to the console.

Note: Because go-sqlite3 is used, each platform needs to be compiled separately.

```go
. /Yi -token githubToken -pwd password -f 1.txt -user username -path /Users/yhy/CodeQL/codeql
```

Considering that there are a bit too many projects to monitor, the github token is required to prevent access from being restricted.

**-path** must be specified to refer to the top-level directory of codeql's various language rulebases.

! [image-20221213212521373](images/image-20221213212521373.png)

Other parameters

```go
-p proxy
-t Monitor a project while running
-f the project to monitor after running, one github project address per line url
-port web access port, default is port 8888
-thread The number of scanning threads, default is 5.

-t -f Specify one or none, add them slowly via the Add button in the web interface.
```

After running, it will automatically generate the relevant folders (downloads, generated databases, clone repositories) and ql rule configuration files in the current directory.


Note: Run the program on a machine with **Codeql** (add environment variable), **Git**, **Docker**, **Go** installed.

**Java**, **Maven**, **Gradle** (if you want to monitor the Java project, otherwise it will lead to database generation failure)

If you need other languages, after modifying the code, it is best to also install the language corresponding to the compilation tool. emmmm is there a docker for all languages?

It's also a good idea to use `root` for execution, because when you use `makefile` in a monitoring project, there may be some tools that are not available on your machine that cause the database to fail, such as.

``go
[2022-12-14 16:34:26] [build-stdout] INFO: yq was not found, installing it
[2022-12-14 16:34:30] [build-stderr] make: go: not enough permissions
[2022-12-14 16:34:30] [build-stderr] make: go: not enough permissions
```

# Security risk
When `codeql` generates a database, it executes a `makefile`-like build process under the project, and there is a security risk here.

So be sure to monitor **trusted** **trusted** **trusted** projects, **don't get a shell bounced**.

All damages caused by **Trusted** are not related to this project or its author***.

All damages caused by ** are not related to this project or its authors**.

**The project and its author are not responsible for any damages caused by ***the project and its author**.

# Function

! [image-20221213143603327](images/image-20221213143603327.png)

! [image-20221215162315622](images/image-20221215162315622.png)

- [x] Monitor projects daily for updates, and fetch/generate databases for Codeql scanning if they are updated
- [x] monitor config file for updates, add new ql rules to fetch from database for scanning
- [x] blacklist, some rules will be false alarms, look at the time to blacklist the results of the scan, the results will not be displayed in the interface when scanning again in the future


# TODO

- [ ] now only adapt Go, Java language, later try to adapt the mainstream language, you can also modify the project where there is "Go", "Java" to add their own other languages
- [ ] codeql create database specify --[no-]db-cluster will automatically create database in all languages, if you don't specify --language, you need to specify github token to automatically analyze --github-auth-stdin
- [ ] Generate databases for download
- [ ] Docker wraps the languages and compilation tools.
- [ ] Read local codeql databases for closed-source or private projects.

# Known issues

- [ x ] http request with occasional `EOF` Solution: limit github access rate.


# 🌟 Star

[! [Stargazers over time](https://starchart.cc/ZhuriLab/Yi.svg)](https://starchart.cc/ZhuriLab/Yi)

# 📄 Disclaimer

This tool is only for legally authorized enterprise security construction behavior, when using this tool for inspection, you should ensure that the behavior is in accordance with local laws and regulations, and has obtained sufficient authorization.

If you use this tool in the process of any illegal behavior or cause all the losses, you need to **self bear the corresponding consequences, this project and its author will not assume any legal and joint liability **.

Before using this tool, please be sure to carefully read and fully understand the contents of the terms, limitations, disclaimers or other provisions involving your significant rights and interests may be bolded, underlined and other forms of attention. Unless you have fully read, fully understand and accept all the terms of this Agreement, please do not use this tool. Your use or any other express or implied acceptance of this Agreement shall be deemed that you have read and agreed to be bound by this Agreement.