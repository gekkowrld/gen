package src

import (
	"log"
	"os"

	"github.com/charmbracelet/huh"
)

func AskLicense(ig LicenseInput) LicenseInput {
	var li LicenseInput
	li = ig

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your license").
				Options(
					listr()...,
				).
				Value(&li.License),
			huh.NewInput().
				Title("Author(s)").
				Value(&li.Author).Placeholder(li.Output),
			huh.NewInput().
				Title("Where will the license be saved?").
				Value(&li.Output).Placeholder(li.Output),
			huh.NewInput().
				Title("Project Year").
				Value(&li.Year).Placeholder(li.Year),
			huh.NewInput().
				Title("Project Name").
				Value(&li.Project).Placeholder(li.Project),
		),
	)
	accessibleMode := os.Getenv("ACCESSIBLE") != ""
	form.WithAccessible(accessibleMode)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	return li
}

type lstr struct {
	name string
	cmd  string
}

func listr() []huh.Option[string] {
	var ltr []huh.Option[string]

	for lk := range AllLicense() {
		lii, _ := Metadata(lk)
		option := huh.Option[string]{
			Key:   lii.Title,
			Value: lk,
		}
		ltr = append(ltr, option)
	}

	return ltr
}
