package segments

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert"
)

func TestQuasar(t *testing.T) {
	cases := []struct {
		Case           string
		ExpectedString string
		Version        string
	}{
		{Case: "@quasar/cli v2.2.1", ExpectedString: "\uea6a 2.2.1", Version: "@quasar/cli v2.2.1"},
	}
	for _, tc := range cases {
		params := &mockedLanguageParams{
			cmd:           "quasar",
			versionParam:  "--version",
			versionOutput: tc.Version,
			extension:     "quasar.config",
		}
		env, props := getMockedLanguageEnv(params)
		env.On("HasFiles", "package-lock.json").Return(false)
		quasar := &Quasar{}
		quasar.Init(props, env)
		assert.True(t, quasar.Enabled(), fmt.Sprintf("Failed in case: %s", tc.Case))
		assert.Equal(t, tc.ExpectedString, renderTemplate(env, quasar.Template(), quasar), fmt.Sprintf("Failed in case: %s", tc.Case))
	}
}
