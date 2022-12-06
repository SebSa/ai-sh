package main

// Path: main.go
// cli tool to ask questions and get answers from openAI api using http requests

// usage: ai ask [question], ai auth [api-key]

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	color "github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func main() {
	if len(os.Args) < 2 {
		color.Yellow("usage: ai ask [question], ai auth [api-key]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "ask":
		// ask question and send and answer to answer()
		resp := ask()
		a := answer(resp)
		// if answer is empty, ask again
		if a == "" {
			color.Red("No answer found, try again")
			main()
		}
		// if answer contains sorry, ask again
		if strings.Contains(a, "Sorry, Can't answer that.") {
			color.Red("No answer found, try again")
			main()
		}
		// present answer to user with selector()
		selector(a)

	case "auth":
		auth()
	default:
		color.Yellow("usage: ai ask [question], ai auth [api-key]")
		os.Exit(1)
	}
}

func auth() {
	if len(os.Args) < 3 {
		color.Yellow("usage: ai auth [api key]")
		os.Exit(1)
	}

	apiKey := os.Args[2]
	// save api key in file in tmp dir
	file, err := os.Create(os.TempDir() + "/openai-api-key")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	color.Green("api key saved")
	// print location of file
	color.White("api key location: " + os.TempDir() + "/openai-api-key")
}

// ask sends a request to openAI api and returns the response
func ask() *CreateCompletionResponse {
	if len(os.Args) < 3 {
		color.Yellow("usage: ai ask [question]")
		os.Exit(1)
	}

	// get api key from temp file
	file, err := os.Open(os.TempDir() + "/openai-api-key")
	if err != nil {
		fmt.Println("api key not found, run ai auth [api-key]")
		os.Exit(1)
	}
	defer file.Close()

	var apiKey string
	_, err = fmt.Fscanf(file, "%s", &apiKey)
	if err != nil {
		log.Fatal(err)
	}

	// get question from args
	question := strings.Join(os.Args[2:], " ")

	// create client
	httpCli := &http.Client{}
	client := New(httpCli, apiKey)

	// create completion request
	input := &CreateCompletionRequest{
		Model:            ModelTextDaVinci,
		Prompt:           []string{"Correctly answer the asked question. Return 'Sorry, Can't Answer that.' if the question isn't related to technology. \n\nQ - get into a docker container.\nA - `docker exec -it mongodb` \n\nQ - Check what's listening on a port.\nA - `lsof -i tcp:4000` \n\nQ - How to ssh into a server with a specific file.\nA - `ssh -i ~/.ssh/id_rsa@127.0.0.1` \n\nQ -  How to set relative line numbers in vim.\nA - `:set relativenumber`\n\nQ - How to create alias?\nA - `alias my_command='my_real_command'`\n\nQ - Tail docker logs.\nA - `docker logs -f mongodb`\n\nQ - Forward port in kubectl.\nA - `kubectl port-forward <pod_name> 8080:3000`\n\nQ - Check if a port is accessible.\nA - `nc -vz host port`\n\nQ - Reverse SSH Tunnel Syntax.\nA - `ssh -R <remote_port>:<local_host>:<local_port> <user>@<remote_host>`\n\nQ - Kill a process running on port 3000.\nA - `lsof -ti tcp:3000 | xargs kill`\n\nQ - Backup database from a mongodb container.\nA - `docker exec -it mongodb bash -c 'mongoexport --db mongodb --collection collections --outdir backup'`\n\nQ - SSH Tunnel Remote Host port into a local port.\nA - `ssh -L <local_port>:<remote_host>:<remote_port> <user>@<remote_host>`\n\nQ - Copy local file to S3.\nA - `aws s3 cp <local_file> s3://<bucket_name>/<remote_file>`\n\nQ - Copy S3 file to local.\nA - `aws s3 cp s3://<bucket_name>/<remote_file> <local_file>`\n\nQ - Recursively remove a folder.\nA - `rm -rf <folder_name>`\n\nQ - Copy a file from local to ssh server.\nA - ` scp /path/to/file user@server:/path/to/destination`\n\nQ - Curl syntax with port.\nA - `curl http://localhost:3000`\n\nQ - Download a file from a URL with curl.\nA - `curl -o <file_name> <URL>`\n\nQ - Git commit with message.\nA - `git commit -m 'my commit message'`\n\nQ - Give a user sudo permissions.\nA - `sudo usermod -aG sudo <user>`\n\nQ - Check what's running on a port?\nA - `lsof -i tcp:<port>`\n\nQ - View last 5 files from history\nA - `history | tail -5`\n\nQ - When was China founded?\nA - Sorry, Can't answer that. \n\nQ - Pass auth header with curl.\nA - `curl -H 'Authorization: Bearer <token>' <URL>`\n\nQ - Filter docker container with labels\nA - `docker ps --filter 'label=<KEY>'`\n\nQ - When was Abraham Lincon born?\nA - Sorry, Can't answer that.\n\nQ - Get into a running kubernetes pod\nA - `kubectl exec -it <pod_name> bash`\n\nQ - Capital city of Ukraine?\nA - Sorry, Can't answer that.\n\nQ - " + question + "A - "},
		MaxTokens:        64,
		Temperature:      0.5,
		FrequencyPenalty: 0.5,
		PresencePenalty:  0,
		Stop:             []string{"\""},
	}

	// make request
	ctx := context.Background()
	s := spinner.New(spinner.CharSets[69], 100*time.Millisecond)
	s.Color("blue", "fgHiWhite")
	s.Suffix = "  🤔  Thinking..."
	s.Start()
	resp, err := client.CreateCompletion(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	s.Stop()
	return resp
}

// Answer processes the answer returned by ask()
func answer(resp *CreateCompletionResponse) string {

	// catch index out of range error
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error:", r)
		}
	}()

	// get answer from response
	answer := resp.Choices[0].Text
	// remove newlines and A - from answer
	answer = strings.ReplaceAll(answer, "A - ", "")
	answer = strings.ReplaceAll(answer, "\n", "")

	// switch on answer
	switch {
	// if answer contains sorry, print it
	case strings.Contains(answer, "Sorry, Can't answer that."):
		return answer

	// if answer contains nothing
	case answer == "":
		return answer

	default:
		// regex to get answer from response /`(.*?)`/
		re := regexp.MustCompile("`(.*)`")
		answer := re.FindStringSubmatch(resp.Choices[0].Text)

		a := answer[1]

		return a
		// print in green
	}
}

// Selector presents the user with a list of options to choose from
func selector(a string) {
	color.Green("Answer:\t'%s'", a)
	// promptui selector to present option to copy to clipboard or run command in terminal
	prompt := promptui.Select{
		Label: "Select an option",
		Items: []string{"Run in terminal", "Try again", "Exit"},
	}

	// get user selection
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	// switch on user selection
	switch result {

	case "Run in terminal":
		// run answer in terminal
		cmd := exec.Command("bash", "-c", a)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

	case "Try again":
		// go back to main
		main()

	case "Exit":
		// exit program
		os.Exit(0)

	}
}
