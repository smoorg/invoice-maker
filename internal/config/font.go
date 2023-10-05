package config

type FontCfg struct {
	Family   string `yaml:"family"`
	Style    string `yaml:"style"`
	Filepath string `yaml:"filepath"`
}

func (f *FontCfg) SetFamily(fontFamily string) {
	f.Family = fontFamily
}

func (f *FontCfg) SetStyle(fontStyle string) {
	f.Style = fontStyle
}

func (f *FontCfg) SetPath(fontPath string) {
	f.Filepath = fontPath
}
