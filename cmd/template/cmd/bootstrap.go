package cmd

import (
	"github.com/kamushadenes/template/internal/bootstrap"

	"github.com/spf13/cobra"
)

var name, language, outputDir, externalTemplatesDir, githubRepo, githubUsername, version string

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstrap a new project",
	Long:  `Bootstrap a new project using template files for different languages or from a GitHub repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var inputs bootstrap.Inputs

		inputs.Project.Name = name

		inputs.Language.Name = language
		inputs.Language.Version = version

		inputs.TemplateSource.Type = bootstrap.BootstrapTypeEmbedded

		inputs.GitHub.Username = githubUsername

		inputs.OutputDir = outputDir

		if githubRepo != "" {
			inputs.TemplateSource.Type = bootstrap.BootstrapTypeGitHub
			inputs.TemplateSource.Source = githubRepo
		}

		if externalTemplatesDir != "" {
			inputs.TemplateSource.Type = bootstrap.BootstrapTypeExternalDir
			inputs.TemplateSource.Source = externalTemplatesDir
		}

		return bootstrap.Bootstrap(&inputs)
	},
}

func init() {
	bootstrapCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the project")
	bootstrapCmd.Flags().StringVarP(&language, "language", "l", "", "Programming language of the project")
	bootstrapCmd.Flags().StringVarP(&version, "version", "v", "", "Programming language version")

	bootstrapCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory for the bootstrapped project")
	bootstrapCmd.Flags().StringVarP(&externalTemplatesDir, "templates-dir", "d", "", "Path to an external templates directory")
	bootstrapCmd.Flags().StringVarP(&githubRepo, "github-repo", "g", "", "GitHub repository URL for templates")
	bootstrapCmd.Flags().StringVarP(&githubUsername, "github-username", "u", "", "GitHub username for go mod")
}
