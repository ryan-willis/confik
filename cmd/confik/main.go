package main

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/ryan-willis/confik/pkg/client"
	"github.com/ryan-willis/confik/pkg/state"
)

func main() {
	var confikState *state.ConfikStateGlobal
	var docker *client.Client
	var err error
	docker, err = client.NewClient()
	confikState = state.Retrieve()
	if err != nil {
		panic(err)
	}
	// TODO: subscribe to docker container events
	// instead of polling every 5 seconds
	for {
		// fmt.Println("sleeping 5 seconds...")
		time.Sleep(5 * time.Second)
		// fmt.Println("checking...")
		syncContainerConfigs(confikState, docker)
	}
}

func syncContainerConfigs(confikState *state.ConfikStateGlobal, docker *client.Client) {
	containerList := docker.GetContainerList()
	hasChanges := false
	for _, container := range containerList {
		var file string
		var contents string
		var restart string
		needsRestart := false
		for labelName, labelValue := range container.Labels {
			if strings.HasPrefix(labelName, "confik") {
				if labelName == "confik.file" {
					file = labelValue
				}
				if labelName == "confik.contents" {
					contents = labelValue
				}
				if labelName == "confik.restart" {
					restart = labelValue
					needsRestart = strings.ToLower(labelValue) == "true" || labelValue == "1"
				}
			}
		}
		if len(file) > 0 {
			existingConfig, exists := (*confikState)[container.ID]
			if exists && strings.Compare(existingConfig.Contents, contents) == 0 && strings.Compare(existingConfig.File, file) == 0 && strings.Compare(existingConfig.Restart, restart) == 0 {
				fmt.Println("Found exact match, skipping " + container.ID)
				continue
			}
			hasChanges = true
			(*confikState)[container.ID] = state.ConfikState{
				File:     file,
				Contents: contents,
				Restart:  restart,
			}
			err := os.WriteFile(file, []byte(contents), fs.ModePerm)
			if err != nil {
				fmt.Printf("failed to create/write to file %s: %s\n", file, err.Error())
				continue
			}
			if needsRestart {
				fmt.Println("Restarting " + container.ID + "...")
				restarted := docker.Restart(container.ID)
				if restarted {
					fmt.Println("Restarted " + container.ID + ".")
				}
			}
		}
	}
	if hasChanges {
		state.Store()
	}
}
