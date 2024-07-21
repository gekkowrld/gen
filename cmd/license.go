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
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "You must supply the license name\n\n")
			cmd.Help()
			os.Exit(1)
		}

		var linput src.LicenseInput
		linput.Author = cmd.Flag("author").Value.String()
		linput.Year = cmd.Flag("year").Value.String()
		linput.License = args[0]
		license, err := src.License(linput)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		fmt.Printf(license)
	},
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
	licenseCmd.PersistentFlags().StringP("output", "o", "LICENSE", "The path to the output file")
}
