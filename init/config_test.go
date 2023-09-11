package init_test

import (
	"os"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/kovercjm/tool-go/init"
)

var _ = Describe("Testing function InitFromEnv", func() {
	Context("for expected behaviors", func() {
		It("should load from .env by default if no file is designated", func() {
			var cfg DotEnv
			init.InitFromEnv(&cfg)

			Ω(cfg.KeyA).Should(Equal("VALUE_A"))
			Ω(cfg.KeyB).Should(Equal("VALUE_B"))
			Ω(cfg.OverwrittenBy).Should(Equal("DOT_ENV"))
		})
		It("should load from designated files", func() {
			Ω(os.MkdirAll("some_file_path/", 0755)).Should(Succeed())
			customEnvFile, err := os.Create("some_file_path/.custom.env")
			Ω(err).Should(Succeed())
			Ω(prepareFile(customEnvFile, []string{
				"KEY_C=VALUE_C\n",
				"OVERWRITTEN_BY=CUSTOM_DOT_ENV\n",
			})).Should(Succeed())

			var cfg struct {
				DotEnv
				KeyC string `envconfig:"KEY_C"`
			}
			init.InitFromEnv(&cfg, "some_file_path/.custom.env")

			Ω(cfg.KeyA).Should(Equal("VALUE_A"))
			Ω(cfg.KeyB).Should(Equal("VALUE_B"))
			Ω(cfg.KeyC).Should(Equal("VALUE_C"))
			Ω(cfg.OverwrittenBy).Should(Equal("DOT_ENV"))

			DeferCleanup(func() {
				Ω(os.RemoveAll("some_file_path")).Should(Succeed())
			})
		})
		It("should load config by prefix for key without struct tag 'envconfig'", func() {
			Ω(os.Setenv("ENV_PREFIX", "DEV")).Should(Succeed())
			Ω(os.Setenv("DEV_KEY_B", "DEV_VALUE_B")).Should(Succeed())
			var cfg DotEnv
			init.InitFromEnv(&cfg)

			Ω(cfg.KeyA).Should(Equal("VALUE_A"))
			Ω(cfg.KeyB).Should(Equal("DEV_VALUE_B"))
			Ω(cfg.OverwrittenBy).Should(Equal("DOT_ENV"))

			Ω(os.Setenv("ENV_PREFIX", "")).Should(Succeed())
		})
		It("should load by struct tag 'default' if the key is not been set", func() {
			os.Clearenv()
			_ = os.Remove(".env")

			var cfg DotEnv
			init.InitFromEnv(&cfg)

			Ω(cfg.KeyA).Should(Equal("Default_A"))
			Ω(cfg.KeyB).Should(Equal("Default_B"))
			Ω(cfg.OverwrittenBy).Should(Equal(""))
		})
	})

	Context("for unexpected behaviors", func() {
		It("should panic if .env file exists, but cannot be read", func() {
			patchLoad := ApplyFuncSeq(godotenv.Load, []OutputCell{
				{Values: Params{os.ErrNotExist}, Times: 1},
				{Values: Params{errors.New("some error")}, Times: 1},
			})

			Ω(func() {
				init.InitFromEnv(&DotEnv{})
			}).ShouldNot(Panic())
			Ω(func() {
				init.InitFromEnv(&DotEnv{})
			}).Should(Panic())

			DeferCleanup(func() {
				patchLoad.Reset()
			})
		})
		It("should panic if env files are designated and exists, but cannot be read", func() {
			patchLoad := ApplyFuncSeq(godotenv.Load, []OutputCell{
				{Values: Params{nil}, Times: 1},
				{Values: Params{os.ErrNotExist}, Times: 1},
				{Values: Params{nil}, Times: 1},
				{Values: Params{errors.New("some error")}, Times: 1},
			})

			Ω(func() {
				init.InitFromEnv(&DotEnv{}, ".custom.env")
			}).ShouldNot(Panic())
			Ω(func() {
				init.InitFromEnv(&DotEnv{}, ".custom.env")
			}).Should(Panic())

			DeferCleanup(func() {
				patchLoad.Reset()
			})
		})
		It("should panic if env files are designated and exists, but cannot be read", func() {
			patchLoad := ApplyFuncSeq(godotenv.Load, []OutputCell{
				{Values: Params{nil}, Times: 1},
				{Values: Params{nil}, Times: 1},
			})
			patchEnvConfig := ApplyFuncReturn(envconfig.Process, errors.New("some error"))

			Ω(func() {
				init.InitFromEnv(&DotEnv{}, ".custom.env")
			}).Should(Panic())

			DeferCleanup(func() {
				patchLoad.Reset()
				patchEnvConfig.Reset()
			})
		})
	})
})
