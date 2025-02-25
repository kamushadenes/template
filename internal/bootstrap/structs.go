package bootstrap

type Inputs struct {
	GitHub         InputGitHub    `yaml:"github"`
	Project        InputProject   `yaml:"project"`
	Language       InputLanguage  `yaml:"language"`
	TemplateSource TemplateSource `yaml:"template-source"`
	OutputDir      string         `yaml:"output-dir"`
}

type InputGitHub struct {
	Username string `yaml:"username"`
}

type InputLanguage struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type InputProject struct {
	Name string `yaml:"name"`
}

type TemplateSource struct {
	Type   BootstrapType `yaml:"type"`
	Source string        `yaml:"source"`
}

type Commands struct {
	PreFiles  []string `yaml:"pre-files"`
	PostFiles []string `yaml:"post-files"`
}
