package state

import (
	"encoding/json"
	"io/fs"
	"os"
)

type ConfikStateGlobal map[string]ConfikState

var State ConfikStateGlobal

type ConfikState struct {
	File     string `json:"file"`
	Contents string `json:"contents"`
	Restart  string `json:"restart"`
}

var StateFile = "/confik_state/state.json"

func Retrieve() *ConfikStateGlobal {
	stateBytes, err := os.ReadFile(StateFile)
	if err != nil {
		stateBytes = []byte("{}")
	}
	err = json.Unmarshal(stateBytes, &State)
	if err != nil {
		panic(err)
	}
	return &State
}

func Store() {
	stateWriteBytes, err := json.Marshal(&State)
	if err != nil {
		panic(err)
	}
	os.WriteFile(StateFile, stateWriteBytes, fs.ModePerm)
}
