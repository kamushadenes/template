package bootstrap

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v3"
)

//go:embed all:templates
var embeddedTemplates embed.FS

const (
	commandsFile = "commands.yml"
)

func Bootstrap(inputs *Inputs) error {
	var templateFS fs.FS

	switch inputs.TemplateSource.Type {
	case BootstrapTypeEmbedded:
		// Use embedded templates
		subDir, err := fs.Sub(embeddedTemplates, "templates")
		if err != nil {
			return fmt.Errorf("error accessing embedded templates: %w", err)
		}
		templateFS = subDir
	case BootstrapTypeExternalDir:
		if _, err := os.Stat(inputs.TemplateSource.Source); os.IsNotExist(err) {
			return fmt.Errorf("template directory does not exist: %w", err)
		}
		templateFS = os.DirFS(inputs.TemplateSource.Source)
	case BootstrapTypeGitHub:
		dir, err := fetchGitHub(inputs.TemplateSource.Source)
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
	}

	// Validate if the template exists
	if _, err := fs.Stat(templateFS, inputs.Language.Name); err != nil {
		return fmt.Errorf("template %s does not exist", inputs.Language.Name)
	}

	if _, err := fs.Stat(templateFS, filepath.Join(inputs.Language.Name, commandsFile)); err != nil {
		return fmt.Errorf("commands file %s does not exist", filepath.Join(inputs.Language.Name, commandsFile))
	}

	b, err := fs.ReadFile(templateFS, filepath.Join(inputs.Language.Name, commandsFile))
	if err != nil {
		return fmt.Errorf("error reading commands file: %w", err)
	}

	var commands Commands

	if err := yaml.Unmarshal(b, &commands); err != nil {
		return fmt.Errorf("error unmarshalling commands file: %w", err)
	}

	if _, err := os.Stat(inputs.OutputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(inputs.OutputDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	for _, command := range commands.PreFiles {
		tmpl, err := template.New("cmd").Parse(command)
		if err != nil {
			return fmt.Errorf("error parsing command: %w", err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, inputs); err != nil {
			return fmt.Errorf("error executing command: %w", err)
		}

		fmt.Println("Executing:", buf.String())

		cmd := exec.Command("sh", "-c", buf.String())
		cmd.Dir = inputs.OutputDir

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error running command: %w", err)
		}
	}

	// Copy files
	filesPath := filepath.Join(inputs.Language.Name, "files")

	if err := fs.WalkDir(templateFS, filesPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(filesPath, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(inputs.OutputDir, relPath)

		// Read file content
		content, err := fs.ReadFile(templateFS, path)
		if err != nil {
			return err
		}

		tmpl, err := template.New("file").Parse(string(content))
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, inputs); err != nil {
			return err
		}

		content = buf.Bytes()

		// Ensure the directory exists
		if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return err
		}

		// Write to destination
		return os.WriteFile(destPath, content, 0644)
	}); err != nil {
		return fmt.Errorf("error copying files: %w", err)
	}

	fmt.Println("Project bootstrapped successfully.")

	return nil
}
