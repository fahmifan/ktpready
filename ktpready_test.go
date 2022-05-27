package ktpready

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKTPReady_Check(t *testing.T) {
	nameChecker := NewNameChecker(2, 60)
	dirtyWords := `
alpha
beta
charlie
delta
`

	bannedWords := `
ban1
ban2
ban3
`
	err := nameChecker.LoadDirtyWords(strings.NewReader(dirtyWords))
	require.NoError(t, err)

	err = nameChecker.LoadBannedWords(strings.NewReader(bannedWords))
	require.NoError(t, err)

	t.Run("failed: contains dirty word", func(t *testing.T) {
		err := nameChecker.Check("john alpha")
		require.Error(t, err)
	})

	t.Run("failed: contains banned word", func(t *testing.T) {
		err := nameChecker.Check("john ban2")
		require.Error(t, err)
	})

	t.Run("failed: too short", func(t *testing.T) {
		err = nameChecker.Check("john")
		require.Error(t, err)
	})

	t.Run("failed: too long", func(t *testing.T) {
		longName := strings.Repeat("a", nameChecker.MaxChar+1)
		err = nameChecker.Check(longName)
		require.Error(t, err)
	})

	t.Run("success: max long", func(t *testing.T) {
		longName := "john " + strings.Repeat("b", nameChecker.MaxChar-4)
		err = nameChecker.Check(longName)
		require.NoError(t, err)
	})

	t.Run("success: no bad names", func(t *testing.T) {
		err = nameChecker.Check("john doe")
		require.NoError(t, err)
	})
}
