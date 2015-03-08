package main

import (
	"log"
	"os/exec"
	"strings"

	"github.com/mitchellh/multistep"
)

// stepSha1 is to fetch the commit hash of the current commit
type stepSha1 struct{}

func (s *stepSha1) Run(state multistep.StateBag) multistep.StepAction {
	output, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		log.Println("Ensure git is installed and you are in a valid repo")
		log.Println(err)
		return multistep.ActionHalt
	}
	sha1 := strings.TrimSpace(string(output))
	state.Put("sha1", sha1)
	imagename := state.Get("imagename").(string)
	state.Put("imagename", imagename+":"+sha1)
	log.Println("sha1:", sha1)
	return multistep.ActionContinue
}

func (s *stepSha1) Cleanup(state multistep.StateBag) {}

// stepAuthor is to fetch the author of the current commit
type stepAuthor struct{}

func (s *stepAuthor) Run(state multistep.StateBag) multistep.StepAction {
	sha1 := state.Get("sha1").(string)
	output, err := exec.Command("git", "--no-pager", "show", "-s", "--format=%an <%ae>", sha1).Output()
	if err != nil {
		log.Println("Ensure git is installed and you are in a valid repo")
		log.Println(err)
		return multistep.ActionHalt
	}
	author := strings.TrimSpace(string(output))
	state.Put("author", author)
	log.Println("author:", author)
	return multistep.ActionContinue
}

func (s *stepAuthor) Cleanup(state multistep.StateBag) {}

// stepBranch is to fetch the branchname of the current commit
type stepBranch struct{}

func (s *stepBranch) Run(state multistep.StateBag) multistep.StepAction {
	output, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		log.Println("Ensure git is installed and you are in a valid repo")
		log.Println(err)
		return multistep.ActionHalt
	}
	branch := strings.TrimSpace(string(output))
	state.Put("branch", branch)
	log.Println("branch:", branch)
	return multistep.ActionContinue
}

func (s *stepBranch) Cleanup(state multistep.StateBag) {}

// stepArchive is to fetch the branchname of the current commit
type stepArchive struct{}

func (s *stepArchive) Run(state multistep.StateBag) multistep.StepAction {
	cmd := exec.Command("git", "archive", "HEAD")
	outpipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("Ensure git is installed and you are in a valid repo")
		log.Println(err)
		return multistep.ActionHalt
	}
	state.Put("pipe", outpipe)
	cmd.Start()
	log.Println("archive: completed")
	return multistep.ActionContinue
}

func (s *stepArchive) Cleanup(state multistep.StateBag) {}

// getFile returns a file from the current commit
func getFile(filename string) (file []byte, err error) {
	file, err = exec.Command("git", "show", "HEAD:"+filename).Output()
	return file, err
}
