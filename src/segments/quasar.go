package segments

import (
	"encoding/json"

	"github.com/jandedobbeleer/oh-my-posh/src/platform"
	"github.com/jandedobbeleer/oh-my-posh/src/properties"
)

type Package struct {
	Version string `json:"version"`
	Dev     bool   `json:"dev"`
}

type Quasar struct {
	language

	HasVite bool
	Vite    *Package
	AppVite *Package
}

func (q *Quasar) Enabled() bool {
	return q.language.Enabled()
}

func (q *Quasar) Template() string {
	return " \uea6a {{.Full}}{{ if .HasVite }} \ueb29 {{ .Vite.Version }}{{ end }} "
}

func (q *Quasar) Init(props properties.Properties, env platform.Environment) {
	q.language = language{
		env:        env,
		props:      props,
		extensions: []string{"quasar.config", "quasar.config.js"},
		commands: []*cmd{
			{
				executable: "quasar",
				args:       []string{"--version"},
				regex:      `(?P<version>((?P<major>[0-9]+).(?P<minor>[0-9]+).(?P<patch>[0-9]+)))`,
			},
		},
		loadContext:        q.loadContext,
		inContext:          q.inContext,
		versionURLTemplate: "https://github.com/quasarframework/quasar/releases/tag/quasar-v{{ .Full }}",
	}
}

func (q *Quasar) loadContext() {
	if !q.language.env.HasFiles("package-lock.json") {
		return
	}

	content := q.language.env.FileContent("package-lock.json")

	var objmap map[string]json.RawMessage
	if err := json.Unmarshal([]byte(content), &objmap); err != nil {
		return
	}

	var dependencies map[string]*Package
	if err := json.Unmarshal(objmap["dependencies"], &dependencies); err != nil {
		return
	}

	if p, ok := dependencies["vite"]; ok {
		q.HasVite = true
		q.Vite = p
	}

	if p, ok := dependencies["@quasar/app-vite"]; ok {
		q.AppVite = p
	}
}

func (q *Quasar) inContext() bool {
	return q.HasVite
}
