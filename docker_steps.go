package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mitchellh/multistep"
)

// stepSlugbuilder is to take the git archive (stdout) and build a slug
type stepSlugbuilder struct{}

func (s *stepSlugbuilder) Run(state multistep.StateBag) multistep.StepAction {
	inpipe := state.Get("pipe").(io.ReadCloser)
	builder := state.Get("builder").(string)
	cmd := exec.Command("docker", "run", "-i", "-a", "stdin", builder)
	cmd.Stdin = inpipe
	output, err := cmd.Output()
	if err != nil {
		log.Println("Ensure docker is installed and daemon is running")
		log.Println(err)
		return multistep.ActionHalt
	}
	containerid := strings.TrimSpace(string(output))
	state.Put("containerid", containerid)
	log.Println("containerid:", containerid)
	err = getLogs(containerid)
	if err != nil {
		log.Println(err)
		return multistep.ActionHalt
	}
	state.Put("slugbuilt", true)
	return multistep.ActionContinue
}

func (s *stepSlugbuilder) Cleanup(state multistep.StateBag) {}

// stepSlugExtract is to extroct the slug from the container
type stepSlugExtract struct{}

func (s *stepSlugExtract) Run(state multistep.StateBag) multistep.StepAction {
	containerid := state.Get("containerid").(string)
	sha1 := state.Get("sha1").(string)
	branch := state.Get("branch").(string)
	slugfile := filepath.Join(".", "slugs", branch, sha1) + string(os.PathSeparator)
	cmd := exec.Command("docker", "cp", containerid+":/tmp/slug.tgz", slugfile)
	err := cmd.Run()
	if err != nil {
		log.Println("Ensure docker is installed and daemon is running")
		log.Println(err)
		return multistep.ActionHalt
	}
	state.Put("slugfile", slugfile)
	log.Println("Slug extracted into:", slugfile)
	return multistep.ActionContinue
}

func (s *stepSlugExtract) Cleanup(state multistep.StateBag) {}

// stepCreateDockerfile is to create a Dockerfile from a template
type stepCreateDockerfile struct{}

func (s *stepCreateDockerfile) Run(state multistep.StateBag) multistep.StepAction {

	type dockerfile struct {
		Author   string
		Sha1     string
		Branch   string
		Slugfile string
		Include  string
		Runner   string
	}
	author := state.Get("author").(string)
	sha1 := state.Get("sha1").(string)
	branch := state.Get("branch").(string)
	slugfile := state.Get("slugfile").(string)
	include := ""
	runner := state.Get("runner").(string)
	file := state.Get("dockerfile").(string)

	content, err := getFile("Dockerfile.include")
	if err != nil {
		log.Println("Dockerfile.include not found in git repo, ignoring")
	} else {
		log.Println("Dockerfile.include found")
		include = string(content)
	}

	d := dockerfile{author, sha1, branch, slugfile, include, runner}
	templ := template.New("mokrfile")
	templ.Parse(dockerfileMokr)

	buf := new(bytes.Buffer)
	err = templ.Execute(buf, d)
	if err != nil {
		log.Println("Failed to execute Dockerfile template")
		log.Println(err)
		return multistep.ActionHalt
	}

	bytes, _ := ioutil.ReadAll(buf)
	if err := ioutil.WriteFile(file, bytes, 0644); err != nil {
		log.Println("Failed to write Dockerfile")
		log.Println(err)
		return multistep.ActionHalt
	}
	log.Println("Dockerfile created:", file)
	return multistep.ActionContinue
}

func (s *stepCreateDockerfile) Cleanup(state multistep.StateBag) {}

// stepBuildImage is to build an image from a Dockerfile
type stepBuildImage struct{}

func (s *stepBuildImage) Run(state multistep.StateBag) multistep.StepAction {
	imagename := state.Get("imagename").(string)
	dockerfile := state.Get("dockerfile").(string)
	cmd := exec.Command("docker", "build", "-t", imagename, "-f", dockerfile, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return multistep.ActionHalt
	}
	state.Put("imagebuilt", true)
	log.Println("image built:", imagename)
	log.Println("For a web app try:")
	log.Println("    docker run -ti -p 8080:8080 -e PORT=8080", imagename, "start web")
	return multistep.ActionContinue
}

func (s *stepBuildImage) Cleanup(state multistep.StateBag) {}

func getLogs(containerid string) error {
	cmd := exec.Command("docker", "logs", "-f", containerid)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}
