package step_9_api_fixtures

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed fixtures
var Fixtures embed.FS

type FixtureLoader struct {
	t           *testing.T
	currentPath fs.FS
}

func NewFixtureLoader(t *testing.T, fixturePath fs.FS) *FixtureLoader {
	return &FixtureLoader{
		t:           t,
		currentPath: fixturePath,
	}
}

func (l *FixtureLoader) LoadString(path string) string {
	file, err := l.currentPath.Open(path)
	require.NoError(l.t, err)

	defer file.Close()

	data, err := io.ReadAll(file)
	require.NoError(l.t, err)

	return string(data)
}

func (l *FixtureLoader) LoadTemplate(path string, data any) string {
	tempData := l.LoadString(path)

	temp, err := template.New(path).Parse(tempData)
	require.NoError(l.t, err)

	buf := bytes.Buffer{}

	err = temp.Execute(&buf, data)
	require.NoError(l.t, err)

	return buf.String()
}
