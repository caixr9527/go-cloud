package config

type redisConfig struct {
	Enable bool `yaml:"enable" mapstructure:"enable"`
	// The network type, either tcp or unix.
	// Default is tcp.
	Network string `yaml:"network" mapstructure:"network"`
	// host:port address.
	Addr string `yaml:"addr" mapstructure:"addr"`

	// ClientName will execute the `CLIENT SETNAME ClientName` command for each conn.
	ClientName string `yaml:"clientName" mapstructure:"clientName"`

	// Protocol 2 or 3. Use the version to negotiate RESP version with redis-server.
	// Default is 3.
	Protocol int `yaml:"protocol" mapstructure:"protocol"`
	// Use the specified Username to authenticate the current connection
	// with one of the connections defined in the ACL list when connecting
	// to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
	Username string `yaml:"username" mapstructure:"username"`
	// Optional password. Must match the password specified in the
	// requirepass server configuration option (if connecting to a Redis 5.0 instance, or lower),
	// or the User Password when connecting to a Redis 6.0 instance, or greater,
	// that is using the Redis ACL system.
	Password string `yaml:"password" mapstructure:"password"`

	// Database to be selected after connecting to the server.
	DB int `yaml:"DB" mapstructure:"DB"`

	// Maximum number of retries before giving up.
	// Default is 3 retries; -1 (not 0) disables retries.
	MaxRetries int `yaml:"maxRetries" mapstructure:"maxRetries"`
	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	MinRetryBackoff int64 `yaml:"minRetryBackoff" mapstructure:"minRetryBackoff"`
	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	MaxRetryBackoff int64 `yaml:"maxRetryBackoff" mapstructure:"maxRetryBackoff"`

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout int64 `yaml:"dialTimeout" mapstructure:"dialTimeout"`
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Supported values:
	//   - `0` - default timeout (3 seconds).
	//   - `-1` - no timeout (block indefinitely).
	//   - `-2` - disables SetReadDeadline calls completely.
	ReadTimeout int64 `yaml:"readTimeout" mapstructure:"readTimeout"`
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.  Supported values:
	//   - `0` - default timeout (3 seconds).
	//   - `-1` - no timeout (block indefinitely).
	//   - `-2` - disables SetWriteDeadline calls completely.
	WriteTimeout int64 `yaml:"writeTimeout" mapstructure:"writeTimeout"`
	// ContextTimeoutEnabled controls whether the client respects context timeouts and deadlines.
	// See https://redis.uptrace.dev/guide/go-redis-debugging.html#timeouts
	ContextTimeoutEnabled bool `yaml:"contextTimeoutEnabled" mapstructure:"contextTimeoutEnabled"`

	// Type of connection pool.
	// true for FIFO pool, false for LIFO pool.
	// Note that FIFO has slightly higher overhead compared to LIFO,
	// but it helps closing idle connections faster reducing the pool size.
	PoolFIFO bool `yaml:"poolFIFO" mapstructure:"poolFIFO"`
	// Base number of socket connections.
	// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	// If there is not enough connections in the pool, new connections will be allocated in excess of PoolSize,
	// you can limit it through MaxActiveConns
	PoolSize int `yaml:"poolSize" mapstructure:"poolSize"`
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout int64 `yaml:"poolTimeout" mapstructure:"poolTimeout"`
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	// Default is 0. the idle connections are not closed by default.
	MinIdleConns int `yaml:"minIdleConns" mapstructure:"minIdleConns"`
	// Maximum number of idle connections.
	// Default is 0. the idle connections are not closed by default.
	MaxIdleConns int `yaml:"maxIdleConns" mapstructure:"maxIdleConns"`
	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActiveConns int `yaml:"maxActiveConns" mapstructure:"maxActiveConns"`
	// ConnMaxIdleTime is the maximum amount of time a connection may be idle.
	// Should be less than server's timeout.
	//
	// Expired connections may be closed lazily before reuse.
	// If d <= 0, connections are not closed due to a connection's idle time.
	//
	// Default is 30 minutes. -1 disables idle timeout check.
	ConnMaxIdleTime int64 `yaml:"connMaxIdleTime" mapstructure:"connMaxIdleTime"`
	// ConnMaxLifetime is the maximum amount of time a connection may be reused.
	//
	// Expired connections may be closed lazily before reuse.
	// If <= 0, connections are not closed due to a connection's age.
	//
	// Default is to not close idle connections.
	ConnMaxLifetime int64 `yaml:"connMaxLifetime" mapstructure:"connMaxLifetime"`

	// Disable set-lib on connect. Default is false.
	DisableIndentity bool `yaml:"disableIndentity" mapstructure:"disableIndentity"`

	// Add suffix to client name. Default is empty.
	IdentitySuffix string `yaml:"identitySuffix" mapstructure:"identitySuffix"`
}
