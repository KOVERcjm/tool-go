package repository

type CacheConfig struct {
	Address       string `default:"localhost:6379" envconfig:"CACHE_ADDRESS"`
	Password      string `default:"" envconfig:"CACHE_PASSWORD"`
	Database      int    `default:"0" envconfig:"CACHE_DATABASE"`
	EnableTLS     bool   `default:"false" envconfig:"CACHE_ENABLE_TLS"`
	ScanBatchSize int64  `default:"10000" envconfig:"CACHE_SCAN_BATCH_SIZE"`
}
