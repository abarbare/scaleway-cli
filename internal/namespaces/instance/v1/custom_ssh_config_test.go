package instance

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/sshconfig"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

func Test_SSHConfigInstall(t *testing.T) {
	t.Run("Install config and create default", core.Test(&core.TestConfig{
		TmpHomeDir: true,
		Commands:   GetCommands(),
		BeforeFunc: createServer("Server"),
		Args:       []string{"scw", "instance", "ssh", "install-config"},
		Check: core.TestCheckCombine(
			core.TestCheckGoldenAndReplacePatterns(
				core.GoldenReplacement{
					Pattern:     regexp.MustCompile("generated to .*scaleway.config"),
					Replacement: "generated to /tmp/scw/.ssh/scaleway.config",
				},
			),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				server := ctx.Meta["Server"].(*instance.Server)

				configPath := sshconfig.ConfigFilePath(ctx.Meta["HOME"].(string))
				content, err := os.ReadFile(configPath)
				assert.Nil(t, err)
				assert.Contains(t, string(content), "Host "+server.Name)

				included, err := sshconfig.ConfigIsIncluded(ctx.Meta["HOME"].(string))
				assert.Nil(t, err)
				assert.True(t, included)
			},
		),
		AfterFunc: deleteServer("Server"),
	}))

	t.Run("Install config and include", core.Test(&core.TestConfig{
		TmpHomeDir: true,
		Commands:   GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			func(ctx *core.BeforeFuncCtx) error {
				homeDir := ctx.Meta["HOME"].(string)
				configPath := sshconfig.DefaultConfigFilePath(homeDir)
				err := os.Mkdir(filepath.Join(homeDir, ".ssh"), 0700)
				assert.Nil(t, err)
				err = os.WriteFile(configPath, []byte(`Host myhost`), 0600)
				assert.Nil(t, err)

				return nil
			},
			createServer("Server"),
		),
		Args: []string{"scw", "instance", "ssh", "install-config"},
		Check: core.TestCheckCombine(
			core.TestCheckGoldenAndReplacePatterns(
				core.GoldenReplacement{
					Pattern:     regexp.MustCompile("generated to .*scaleway.config"),
					Replacement: "generated to /tmp/scw/.ssh/scaleway.config",
				},
			),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				server := ctx.Meta["Server"].(*instance.Server)

				defaultConfigPath := sshconfig.DefaultConfigFilePath(ctx.Meta["HOME"].(string))
				content, err := os.ReadFile(defaultConfigPath)
				assert.Nil(t, err)
				assert.Contains(t, string(content), "Include scaleway.config")
				assert.Contains(t, string(content), "Host myhost")

				configPath := sshconfig.ConfigFilePath(ctx.Meta["HOME"].(string))
				content, err = os.ReadFile(configPath)
				assert.Nil(t, err)
				assert.Contains(t, string(content), "Host "+server.Name)

				included, err := sshconfig.ConfigIsIncluded(ctx.Meta["HOME"].(string))
				assert.Nil(t, err)
				assert.True(t, included)
			},
		),
		AfterFunc: deleteServer("Server"),
	}))
}
