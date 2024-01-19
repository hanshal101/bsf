package build

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"

	"github.com/buildsafedev/bsf/cmd/styles"
	nixcmd "github.com/buildsafedev/bsf/pkg/nix/cmd"
)

// BuildCmd represents the build command
var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "builds the project",
	Long: `builds the project based on instructions defined in bsf.hcl.
	Build occurs in a sandboxed environment where only current directory is available. 
	It is recommended to check in the files in version control system(ex: Git) before building.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat("bsf.hcl"); err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: Has the project been initialized?"))
			fmt.Println(styles.HintStyle.Render("hint: ", "run `bsf init` to initialize the project"))
			os.Exit(1)
		}
		fmt.Println(styles.HighlightStyle.Render("Building..."))

		_, err := nixcmd.Build()
		if err != nil {
			gotHash := isHashMismatchError(err.Error())
			if gotHash != "" {
				fmt.Println(fmt.Sprintf(styles.ErrorStyle.Render("Hash mismatch detected.")))
				fmt.Println(styles.HintStyle.Render("hint: 1) please insert the following hash in the build app/module section(ex: vendorHash, vendorSha256) of bsf.hcl` : "))
				err = clipboard.WriteAll(gotHash)
				copied := true
				if err != nil {
					fmt.Println(styles.ErrorStyle.Render("Failed to copy hash to clipboard: ", err.Error()))
					copied = false
				}
				if copied {
					fmt.Println(styles.HighlightStyle.Render(gotHash + " (copied to clipboard) " + "📋"))
				} else {
					fmt.Println(styles.HighlightStyle.Render(gotHash))
				}
				fmt.Println(styles.HintStyle.Render("2) run `bsf generate` to lock dependencies and generate Nix files"))
				os.Exit(1)
			}

			if isNoFileError(err.Error()) {
				fmt.Println(styles.ErrorStyle.Render(err.Error() + "\n Please ensure all necessary files are added/committed in your version control system"))
				fmt.Println(styles.HintStyle.Render("hint: run git add .  "))
				os.Exit(1)
			}

			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))

			os.Exit(1)
		}

		fmt.Println(styles.SucessStyle.Render("Build completed successfully, please check the result directory"))

	},
}

func isHashMismatchError(err string) string {
	if strings.Contains(err, "hash mismatch") {
		re := regexp.MustCompile(`got:\s+(sha256-.*)`)
		matches := re.FindStringSubmatch(err)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}

func isNoFileError(err string) bool {
	return strings.Contains(err, "No such file or directory")
}