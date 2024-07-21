/*
Copyright © 2024 Gekko Wrld

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"codeberg.org/gekkowrld/gen/src"
	conf "github.com/gekkowrld/go-gitconfig"
	"github.com/spf13/cobra"
)

// licenseCmd represents the license command
var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Generate a license file",
	Long: `Generate a license file

e.g "gen license gpl3" to generate GNU GPL version 3`,
	Run: func(cmd *cobra.Command, args []string) {

		if cmd.Flag("all").Changed {
			fmt.Println("Available licenses:")
			var i int
			for lkeys := range src.LicenseMaps {
				i++
				meta, _ := src.Metadata(lkeys)
				fmt.Printf("%s:\n%s\n", meta.Title, meta.Description)
			}
			os.Exit(0)
		}

		// Display info
		info := cmd.Flag("info").Value.String()
		if !is_empty(info) {
			printInfo(info)
		}

		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "You must supply the license name\n\n")
			cmd.Help()
			os.Exit(1)
		}

		// Write about a file
		var linput src.LicenseInput
		linput.Author = cmd.Flag("author").Value.String()
		linput.Year = cmd.Flag("year").Value.String()
		linput.License = args[0]
		output := cmd.Flag("output").Value.String()
		license, err := src.License(linput)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		if output == "1" {
			fmt.Println(license)
		} else {
			err = src.FileWrite(license, output)
			if err != nil {
				log.Fatalf("%+v\n", err)
			}
		}

	},
}

func is_empty(value string) bool {
	if strings.TrimSpace(value) == "" {
		return true
	} else {
		return false
	}
}

func printInfo(license string) {
	meta, err := src.Metadata(license)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	if !is_empty(meta.Nickname) {
		fmt.Printf("%s (%s) - %s\n", meta.Title, meta.SpdxId, meta.Nickname)
	} else {
		fmt.Printf("%s (%s)\n", meta.Title, meta.SpdxId)
	}
	if !is_empty(meta.Description) {
		fmt.Println(meta.Description)
	}
	if !is_empty(meta.Note) {
		fmt.Println("Note:\n", meta.Note)
	}
	if !is_empty(arr_d(meta.Conditions)) {
		fmt.Println("Conditions: ")
		fmt.Println(arr_d(meta.Conditions))
	}
	if !is_empty(arr_d(meta.Limitations)) {
		fmt.Println("Limitations: ")
		fmt.Println(arr_d(meta.Limitations))
	}
	if !is_empty(arr_d(meta.Permissions)) {
		fmt.Println("Permissions: ")
		fmt.Println(arr_d(meta.Permissions))
	}

	os.Exit(0)
}

func arr_d(strs []string) string {
	var __str strings.Builder
	for _, str := range strs {
		if !is_empty(str) {
			__str.WriteString("- ")
			__str.WriteString(str)
			__str.WriteRune('\n')
		}
	}

	return __str.String()
}

func init() {
	rootCmd.AddCommand(licenseCmd)
	gname, _ := conf.GetValue(conf.OptionsPassed{ConfigKey: "user.name"})
	licenseCmd.PersistentFlags().StringP("author", "a", gname, "The authors name")
	cwd, _ := os.Getwd()
	cwd = filepath.Base(cwd)
	licenseCmd.PersistentFlags().StringP("project", "p", cwd, "The project name")
	year := time.Now().Year()
	licenseCmd.PersistentFlags().StringP("year", "y", strconv.Itoa(year), "The project year")
	licenseCmd.PersistentFlags().StringP("output", "o", "LICENSE", "The path to the output file. 1 for stdout")
	licenseCmd.PersistentFlags().BoolP("all", "A", false, "List all the licenses available")
	licenseCmd.PersistentFlags().StringP("info", "i", "", "List information about a license")
}
