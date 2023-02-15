package repository

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type Config struct {
	DBConfig
}
type DBConfig struct {
	DBHost                      string `default:"localhost" envconfig:"DB_HOST"`
	DBUser                      string `default:"root" envconfig:"DB_USER"`
	DBPassword                  string `default:"" envconfig:"DB_PASSWORD"`
	DBName                      string `envconfig:"DB_NAME"`
	DBLocation                  string `default:"Local" envconfig:"DB_LOCATION"`
	DBEnableTLS                 bool   `default:"false" envconfig:"DB_ENABLE_TLS"`
	DBCaCertPEM                 string `default:"/tmp/mysql-ca-cert.pem" envconfig:"DB_CA_CERT_PEM"`
	DBCustomTLSKey              string `default:"custom" envconfig:"DB_CUSTOM_TLS_KEY"`
	DBConnPoolSize              int    `default:"5" envconfig:"DB_CONN_POOL_SIZE"`
	DBLogLevel                  string `default:"info" envconfig:"DB_LOG_LEVEL"`
	DBSlowThresholdMS           int    `default:"100" envconfig:"DB_SLOW_THRESHOLD_MS"`
	DBDefaultPageSize           int    `default:"10" envconfig:"DB_DEFAULT_PAGE_SIZE"`
	DBIgnoreRecordNotFoundError bool   `default:"false" envconfig:"DB_IGNORE_RECORD_NOT_FOUND_ERROR"`
}

var mysqlTLSOnce sync.Once

func (c Config) MySQL() *mysql.Config {
	var err error
	config := &mysql.Config{
		Addr:                 c.DBHost,
		User:                 c.DBUser,
		Passwd:               c.DBPassword,
		DBName:               c.DBName,
		Net:                  "tcp",
		Loc:                  time.Local,
		ParseTime:            true,
		Params:               map[string]string{"multiStatements": "true", "charset": "utf8mb4"},
		AllowNativePasswords: true,
		Collation:            "utf8mb4_unicode_ci",
	}
	if c.DBLocation != "" {
		config.Loc, err = time.LoadLocation(c.DBLocation)
		if err != nil {
			panic(err)
		}
	}
	if c.DBEnableTLS {
		mysqlTLSOnce.Do(func() {
			var pem []byte
			pem, err = os.ReadFile(c.DBCaCertPEM)
			if err != nil {
				err = errors.Wrap(err, "failed to read ca cert pem")
				return
			}
			rootCertPool := x509.NewCertPool()
			if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
				err = errors.Wrap(err, "failed to append ca cert pem")
				return
			}
			if err = mysql.RegisterTLSConfig(c.DBCustomTLSKey, &tls.Config{RootCAs: rootCertPool}); err != nil {
				err = errors.Wrap(err, "failed to register tls config")
				return
			}
			config.TLSConfig = c.DBCustomTLSKey
		})
	}
	return config
}
