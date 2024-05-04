package editor

import (
	"io/ioutil"
	"os"
	"os/exec"
)

const DefaultEditor = "nano"

// Edit the text using editor
func Edit(text string, tmpPattern string) (string, error) {
	if tmpPattern == "" {
		tmpPattern = "tpot_*.txt"
	}

	f, err := ioutil.TempFile("", tmpPattern)
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())
	_, err = f.Write([]byte(text))
	if err != nil {
		return "", err
	}

	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")
		if editor == "" {
			editor = DefaultEditor
		}
	}

	cmd := exec.Command(editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	readFile, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return "", err
	}
	return string(readFile), nil
}
