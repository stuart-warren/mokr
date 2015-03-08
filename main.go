package main

import (
	"log"
	"os"

	"github.com/mitchellh/multistep"
)

type stepConfig struct{}

func (s *stepConfig) Run(state multistep.StateBag) multistep.StepAction {
	image := "test"
	if len(os.Args) > 1 {
		image = os.Args[1]
	}
	state.Put("imagename", image)
	state.Put("dockerfile", "Dockerfile.mokr")
	state.Put("runner", "deis/slugrunner")
	state.Put("builder", "deis/slugbuilder")
	return multistep.ActionContinue
}

func (s *stepConfig) Cleanup(state multistep.StateBag) {}

func buildSlugRunner() *multistep.BasicStateBag {
	state := new(multistep.BasicStateBag)
	steps := []multistep.Step{
		&stepConfig{},
		&stepSha1{},
		&stepAuthor{},
		&stepBranch{},
		&stepArchive{},
		&stepSlugbuilder{},
		&stepSlugExtract{},
		&stepCreateDockerfile{},
		&stepBuildImage{},
	}
	runner := &multistep.BasicRunner{Steps: steps}
	runner.Run(state)
	return state
}

func main() {
	state := buildSlugRunner()
	_, ok := state.GetOk("imagebuilt")
	if !ok {
		log.Fatal("Something went wrong :(")
	}
	log.Println("Great success!")
}
