package filesystem

type ConfigStruct struct {
	Webhooks struct {
		Version string `toml:"version"`
		Images  string `toml:"images"`
		Text    string `toml:"text"`
	} `toml:"webhooks"`
}
