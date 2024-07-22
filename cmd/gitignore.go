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

	"codeberg.org/gekkowrld/gen/src"
	"github.com/spf13/cobra"
)

// gitignoreCmd represents the gitignore command
var gitignoreCmd = &cobra.Command{
	Use:   "gitignore",
	Short: "Create a .gitignore for different languages",
	Long:  `Create a .gitignore for different languages`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatalf("You require atleast one argument for the gitignore template")
		}
		var gitignores string
		of := cmd.Flag("output").Value.String()
		if of == "1" {
			gitignores = src.GitIgnore(src.GitInput{Ignores: args})
		} else {
			gitignores = src.GitIgnore(src.GitInput{Ignores: args, IsInput: true, Output: of})
		}
		fmt.Println(gitignores)
	},
}

func init() {
	rootCmd.AddCommand(gitignoreCmd)
	gitignoreCmd.PersistentFlags().BoolP("all", "A", false, "List all the available gitignore templates")
	gitignoreCmd.PersistentFlags().StringP("output", "o", ".gitignore", "The output file. use 1 for stdout")
}
