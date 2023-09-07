package config_test

import (
	"bufio"
	"io"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

type DotEnv struct {
	KeyA          string `default:"Default_A" envconfig:"KEY_A"`
	KeyB          string `default:"Default_B" split_words:"true"`
	OverwrittenBy string `envconfig:"OVERWRITTEN_BY"`
}

var defaultEnvFileContent = []string{
	"KEY_A=VALUE_A\n",
	"KEY_B=VALUE_B\n",
	"OVERWRITTEN_BY=DOT_ENV\n",
}

var _ = BeforeSuite(func() {
	defaultEnvFile, err := os.Create(".env")
	Ω(err).Should(Succeed())
	Ω(prepareFile(defaultEnvFile, defaultEnvFileContent)).Should(Succeed())

	DeferCleanup(func() {
		Ω(defaultEnvFile.Close()).Should(Succeed())
	})
})

var _ = AfterSuite(func() {
	os.Clearenv()
})

func prepareFile(file io.Writer, lines []string) error {
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := writer.WriteString(line); err != nil {
			return err
		}
	}
	return writer.Flush()
}
