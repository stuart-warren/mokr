package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func pipeCommands(commands ...*exec.Cmd) ([]byte, error) {
	for i, command := range commands[:len(commands)-1] {
		out, err := command.StdoutPipe()
		if err != nil {
			return nil, err
		}
		command.Start()
		commands[i+1].Stdin = out
	}
	final, err := commands[len(commands)-1].Output()
	if err != nil {
		return nil, err
	}
	return final, nil
}

func main() {
	output, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	sha := strings.TrimSpace(string(output))
	if err != nil {
		log.Fatal("Ensure git is installed and you are in a git repo.", err)
	}
	output, _ = exec.Command("git", "--no-pager", "show", "-s", "--format='%an <%ae>'", sha).Output()
	committer := strings.TrimSpace(string(output))
	fmt.Println(committer)
	gcmd3 := exec.Command("git", "archive", "HEAD")
	dcmd := exec.Command("docker", "run", "-i", "-a", "stdin", "deis/slugbuilder")
	output, err = pipeCommands(gcmd3, dcmd)
	if err != nil {
		log.Fatalf("%v", err)
	}
	cont := strings.TrimSpace(string(output))
	dcmd2 := exec.Command("docker", "logs", "-f", cont)
	dcmd2.Stdout = os.Stdout
	dcmd2.Run()
	fmt.Println("Extracting slug...")
	slugfile := fmt.Sprintf("./slug-%s.tgz", sha)
	dcmd3 := exec.Command("docker", "cp", cont+":/tmp/slug.tgz", slugfile)
	dcmd3.Run()
	// TODO:
	// create Dockerfile
	// - FROM deis/slugrunner
	// - MAINTAINER committer
	// - ADD slug.tgz /app
	// - Dockerfile.include
	// - ENTRYPOINT ["/runner/init"]
	// docker build
	// rm slug if not debug
	// push
	// update marathon
}
