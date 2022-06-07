package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"time"
)

type envs struct {
	username string
	hostname string
}

type pipeline struct {
	ID       int    `json:"id"`
	Stage    string `json:"stage"`
	Status   string `json:"status"`
	Name     string `json:"name"`
	Pipeline struct {
		CreatedAt string `json:"created_at"`
		ID        string `json:"id"`
		Status    string `json:"status"`
		WebURL    string `json:"web_url"`
	} `json:"pipeline"`
	WebURL string `json:"web_url"`
}

var (
	counter   int
	branch    string
	project   string
	err       error
	resp      []byte
	sample    []string
	ripest    []string
	fulresp   []pipeline
	urlstring = "https://gitlab.com/api/v4/projects/"
	headers   = map[string]string{
		"PRIVATE-TOKEN": "",
		"Content-Type":  "application/json",
	}
	//for custom header
	env = envs{
		username: "USER",
		hostname: "HOSTNAME",
	}
)

func myUsage() {
	fmt.Printf("Usage:\n  %s [OPTIONS] project stage1 stage2 ...\nOptions:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Printf("Example:\n  %s -b master /backend/tabl-suite migrate deploy \n", os.Args[0])
	fmt.Printf(" OR\n  %s --b master /frontend \"run tests\" deploy \n", os.Args[0])
	fmt.Printf(" OR\n  %s -b=master /backend/suite Db_migrate deploy \n\n", os.Args[0])
}
func breakError(err error, message string) {
	if message != "" {
		fmt.Printf(message)
	}
	fmt.Println(err)
	os.Exit(1)
}

func checkStage(stg string) bool {
	for i := 0; i < len(ripest); i++ {
		if stg == ripest[i] {
			return true
		}
	}
	return false
}

func main() {

	flag.Usage = myUsage
	//flags
	flag.StringVar(&branch, "b", "develop", "Source branch")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Printf("\n Required parameters are missing \n")
		flag.Usage()
		os.Exit(1)
	}

	args := ToLow(flag.Args()) //Подразумевает что проект тоже будет с маленькой буквы
	for i := 1; i <= len(args)-1; i++ {
		if "build" == args[i] {
			fmt.Println(` Warning: stage "build" usually starts automatically when the pipeline start`)
		}
	}

	project = url.PathEscape(args[0])

	_, err = GetReq(nil, urlstring+project, headers)
	if err != nil {
		breakError(err, "\n Possibly a mistake in the name of the project\n")
	}

	stendname, err := CustomEnv(env.hostname, env.username)
	if err != nil {
		breakError(err, "")
	}
	//Create new branch
	if branch != "develop" {
		_, err = GetReq(nil, urlstring+project+"/repository/branches/"+url.PathEscape(branch), headers)
		if err != nil {
			breakError(err, "\n Maybe the wrong branch?\n")
		}
	}
	newbranch := "cli/" + stendname + "-" + RandString(strconv.Itoa(1 + rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000-2)))[:6]
	_, err = PostReq(nil, urlstring+project+"/repository/branches/"+"?branch="+newbranch+"&ref="+branch, headers)
	if err != nil {
		breakError(err, "")
	}
	//Run pipline
	data := map[string]interface{}{
		"ref": newbranch,
		"variables": [1]map[string]string{{
			"key":   "CUSTOM_ENV",
			"value": stendname,
		}},
	}
	resp, err = PostReq(data, urlstring+project+"/pipeline", headers)
	if err != nil {
		breakError(err, "")
	}

	var respid pipeline
	json.Unmarshal(resp, &respid)
	if respid.ID != 0 { //Если в ответе айди не пустое
		fmt.Printf(" Pipeline was created: %s\n", respid.WebURL) //Пишем ссылку

		for { //И начинаем бесконечный цикл
			time.Sleep(10 * time.Second)
			//Запрашиваем джобы
			resp, err = GetReq(nil, urlstring+project+"/pipelines/"+strconv.Itoa(respid.ID)+"/jobs", headers)
			if err != nil {
				breakError(err, "")
			}
			json.Unmarshal(resp, &fulresp)
			if len(sample) < len(fulresp) {
				for i := 0; i < len(fulresp); i++ {
					gdv := fulresp[i]
					//var gdv []string
					stages := gdv.Stage
					sample = append(sample, stages)
				}
			}

			if len(sample) == len(fulresp) {
				if Compar(args[1:], sample) != true {
					_, err = DelReq(nil, urlstring+project+"/repository/branches/"+url.PathEscape(newbranch), headers)
					if err != nil {
						breakError(err, "")
					}
					err := fmt.Errorf("\n Apparently you are mistaken in the name of the stages.\n %s", respid.WebURL)
					breakError(err, "")
				}
			}

			for i := 0; i < len(fulresp); i++ { //Проверка статутсов
				abc := fulresp[i]
				if abc.Status == "failed" || abc.Status == "canceled" {
					_, err = DelReq(nil, urlstring+project+"/repository/branches/"+url.PathEscape(newbranch), headers)
					if err != nil {
						breakError(err, "")
					}
					err = fmt.Errorf("\n Something went wrong while doing the job.\n Status received: %s\n More details:    %s", abc.Status, abc.Pipeline.WebURL)
					breakError(err, "")
				}
				if abc.Status == "pending" || abc.Status == "running" {
					continue
				}
				if abc.Status == "success" && checkStage(abc.Stage) != true && len(ripest) == len(args)-1 {
					fmt.Printf(" %s: success \n", abc.Name)
					_, err = DelReq(nil, urlstring+project+"/repository/branches/"+url.PathEscape(newbranch), headers)
					if err != nil {
						breakError(err, "")
					}
					os.Exit(0)
				}
				if abc.Status == "success" && checkStage(abc.Stage) != true {
					fmt.Printf(" %s: success \n", abc.Name)
					ripest = append(ripest, abc.Stage)
					continue
				}

				if len(args) != 1 && checkStage(abc.Stage) != true && len(ripest) > 0 {
					for i := 0; i < len(args); i++ {
						if len(ripest) > counter && counter < len(args) && abc.Stage == args[counter+1] {
							counter++
							_, err = PostReq(nil, urlstring+project+"/jobs/"+strconv.Itoa(abc.ID)+"/play", headers)
							if err != nil {
								breakError(err, "")
							} else {
								continue
							}
						}
					}
				} else if len(args) == 1 && checkStage(abc.Stage) == true {
					_, err = DelReq(nil, urlstring+project+"/repository/branches/"+url.PathEscape(newbranch), headers)
					if err != nil {
						breakError(err, "")
					}
					os.Exit(0)
				} else {
					continue
				}
			}
		}
	}
} //main
