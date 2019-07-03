package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	//Initialisation
	//imported cli package https://github.com/urfave/cli to create the command line app
	var username, password, serverUrl, destination, inputTitle, description string
	app := cli.NewApp()
	app.Name = "CPR"
	app.Usage = "A tool creating pull request in bitbucket"
	app.EnableBashCompletion = true
	app.Author = "Zehui Zhang"
	app.Email = "peter.zhang@simpro.co"
	app.Version = "1.0"
	app.Compiled = time.Now()

	//since short name functionality is not compatible with input source
	//used all full name at this moment
	flags := []cli.Flag{
		cli.StringFlag{
			Name:   "load",
			EnvVar: "CPR_CONFIG_FILE_PATH",
			Usage:  "load .yaml config file from the path or load from environment variable in default",
		},
		altsrc.NewStringFlag(
			cli.StringFlag{
				Name:        "username",
				Usage:       "Bitbucket account username",
				Destination: &username,
			}),
		altsrc.NewStringFlag(
			cli.StringFlag{
				Name:        "password",
				Usage:       "Bitbucket account password",
				Destination: &password,
			}),
		altsrc.NewStringFlag(
			cli.StringFlag{
				Name:        "url",
				Usage:       "bitbucket server url ",
				Value:       "https://bitbucket.simprocloud.com",
				Destination: &serverUrl,
			}),
		altsrc.NewStringFlag(
			cli.StringFlag{
				Name:        "destBranch",
				Usage:       "PR destination branch",
				Destination: &destination,
			}),
		altsrc.NewStringFlag(
			cli.StringFlag{
				Name:        "title",
				Usage:       "PR title, branch name will be used if the title is not given",
				Destination: &inputTitle,
			}),
		cli.StringFlag{
			Name:        "description",
			Usage:       "PR description, could be empty",
			Destination: &description,
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "turn debug on, will turn on all arguments and flags value",
		},
		altsrc.NewStringSliceFlag(
			cli.StringSliceFlag{
				Name:  "reviewer",
				Usage: "PR reviewer `firstName.lastName`, could be multiple",
			}),
	}

	//Main execution part
	app.Action = func(c *cli.Context) error {
		var currentBranch string
		var PRTitlePointer *string

		if password == "" {
			log.Fatal("No Password Given, Exit")
		}

		// if title not given in cli then using current branch name instead for PR Title
		if inputTitle == "" {
			PRTitlePointer = &currentBranch
		} else {
			PRTitlePointer = &inputTitle
		}

		//get reviewers from flag
		reviewers := c.StringSlice("reviewer")

		//get username from git email if it is empty
		if username == "" {
			out, err := exec.Command("git", "config", "--global", "user.email").Output()
			logMessage := "No username found from [git config --global user.email] ;"
			if err != nil {
				log.Fatal(logMessage, err)
			}
			username = strings.Split(string(out), "@")[0]
		}
		if !strings.Contains(username, ".") {
			log.Fatal("username has no dot symbol inside, ", username)
		}

		//get project and repo name from git remote.origin.url
		out, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
		logMessage := "remoteUrl is not found from [git config --get remote.origin.url] ;"
		if err != nil {
			log.Fatal(logMessage, err)
		}
		remoteUrl := strings.TrimSuffix(string(out), "\n")
		if remoteUrl == "" {
			log.Fatal(logMessage, remoteUrl)
		}
		urlSplits := strings.Split(remoteUrl, "/")
		project := urlSplits[len(urlSplits)-2]
		repo := strings.TrimSuffix(urlSplits[len(urlSplits)-1], ".git")

		//get current branch name
		out, err = exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
		logMessage = "currentBranch name is not found from [git rev-parse --abbrev-ref HEAD] ;"
		if err != nil {
			log.Fatal(logMessage, err)
		}
		currentBranch = strings.TrimSuffix(string(out), "\n")
		if currentBranch == "" {
			log.Fatal(logMessage, currentBranch)
		}

		//Construct Bitbucket url with project and repo
		url := serverUrl + "/rest/api/1.0/projects/" + project + "/repos/" + repo + "/pull-requests"

		//Construct post json data
		type user struct {
			Name string `json:"name"`
		}
		type reviewerObject struct {
			User user `json:"user"`
		}
		var reviewersArray []reviewerObject
		for _, reviewer := range reviewers {
			var r reviewerObject
			r.User.Name = reviewer
			reviewersArray = append(reviewersArray, r)
		}
		values := map[string]interface{}{
			"title":       *PRTitlePointer,
			"description": description,
			"fromRef":     map[string]string{"id": "refs/heads/" + currentBranch},
			"toRef":       map[string]string{"id": "refs/heads/" + destination},
			"reviewers":   reviewersArray,
		}
		jsonValue, _ := json.Marshal(values)

		//if has --debug in cli, it will print out all parameters
		if c.Bool("debug") {
			fmt.Println("Post Json - ", string(jsonValue))
			fmt.Println("username - ", username)
			fmt.Println("password - ", password)
			fmt.Println("destBranch - ", destination)
			fmt.Println("title - ", *PRTitlePointer)
			fmt.Println("description - ", description)
			fmt.Println("reviewers - ", reviewers)
			fmt.Println("currentBranch - ", currentBranch)
			fmt.Println("repo - ", repo)
			fmt.Println("project - ", project)
			fmt.Println("url - ", url)
		}

		//Create new http request and post to Bitbucket API
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
		if err != nil {
			panic(err)
		}
		req.SetBasicAuth(username, password)
		req.Header.Set("Content-Type", "application/json")
		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		//parse result, print out success message or the errors
		fmt.Println("---------------------Return -----------------------")
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			var result map[string]interface{}
			err := json.Unmarshal(body, &result)
			if err != nil {
				panic(err)
			}
			link := result["links"].(map[string]interface{})["self"].([]interface{})[0].(map[string]interface{})["href"]
			fmt.Println("Succeed")
			fmt.Println("PR link -> ", link)
			fmt.Println("Title -> ", *PRTitlePointer)
		} else {
			fmt.Println("Argh! Broken")
			fmt.Printf("%v", string(body))
		}
		return nil
	}

	//Read config from yaml file registered on load flag before app starts action function
	app.Before = altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("load"))
	app.Flags = flags
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
