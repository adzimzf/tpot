package ui

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

// Confirm shows basic popup confirmation Yes or No
// return false if user select no
func Confirm(text string) (bool, error) {
	prompt := promptui.Prompt{
		Label: text + " [Y/y/N/n]",
		Validate: func(s string) error {
			if s == "y" || s == "Y" || s == "N" || s == "n" {
				return nil
			}
			return fmt.Errorf("invalid option")
		},
	}
	run, err := prompt.Run()
	if err != nil {
		return false, err
	}
	return run == "y" || run == "Y", nil
}
