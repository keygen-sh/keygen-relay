package db

type Config struct {
	DatabaseFilePath string
	DatabasePragmas  map[string]string
}

func NewConfig() *Config {
	return &Config{DatabasePragmas: map[string]string{
		"journal_mode": "WAL",         // set the journal mode to Write-Ahead Logging for concurrency
		"synchronous":  "NORMAL",      // set synchronous mode to NORMAL to better balance performance and safety
		"busy_timeout": "5000",        // set busy timeout to 5 seconds to avoid lock-related errors
		"cache_size":   "-20000",      // set cache size to 20MB for faster data access
		"foreign_keys": "ON",          // enable foreign key constraints
		"auto_vacuum":  "INCREMENTAL", // enable auto vacuuming and set it to incremental mode for gradual space reclaiming
		"temp_store":   "MEMORY",      // store temporary tables and data in memory for better performance
		"mmap_size":    "2147483648",  // set the mmap_size to 2GB for faster read/write access using memory-mapped I/O
		"page_size":    "8192",        // set the page size to 8KB for balanced memory usage and performance
	}}
}
