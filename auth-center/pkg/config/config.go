package config

import (
	"flag"
	"time"

	"github.com/runtime-radar/runtime-radar/lib/config"
)

// Config represents system configuration.
type Config struct {
	NewDB                  bool          // forces recreation of DB
	PostgresAddr           string        // Postgres address in host[:port] format
	PostgresDB             string        // Postgres db name
	PostgresUser           string        // Postgres user
	PostgresPassword       string        // Postgres password
	PostgresSSLMode        bool          // Postgres SSL mode
	PostgresSSLCheckCert   bool          // Check postgres SSL cert
	LogLevel               string        // log level can be INFO, WARN, ERROR, FATAL, DEBUG or ALL
	ListenGRPCAddr         string        // address "[host]:port" that server should be listening on
	ListenHTTPAddr         string        // address "[host]:port" that server should be listening for health checks
	InstrumentationAddr    string        // address "[host]:port" that instrumentation server should be listening for health checks and metrics
	TLS                    bool          // is TLS enabled?
	TokenKey               string        // key for jwt token
	Auth                   bool          // is auth enabled?
	OwnCSURL               string        // URL of current CS (http(s)://host[:port]).
	GopsAddr               string        // gops listen address
	PathToPasswordSecList  string        // nested path to file that context bad passwords
	AccessTokenExpiration  time.Duration // how long access token will be alive in minute. By default, 15 minute
	RefreshTokenExpiration time.Duration // how long refresh token will be alive in minute. By default, 8640 minute
	AdminUsername          string        // predeclared admin username
	AdminPassword          string        // predeclared admin password
	AdminEmail             string        // predeclared admin email
}

// New reads config from environment and returns pointer to a new Config.
func New() *Config {
	c := &Config{}

	flag.BoolVar(&c.NewDB, "newDB", config.LookupEnvBool("NEW_DB", false), "Set this flag to recreate DB from scratch. ALL EXISTING DATA WILL BE LOST.")
	flag.StringVar(&c.PostgresAddr, "postgresAddr", config.LookupEnvString("POSTGRES_ADDR", "127.0.0.1:5432"), "Set PostgreSQL address as host:port, where port is optional (without TLS).")
	flag.StringVar(&c.PostgresDB, "postgresDB", config.LookupEnvString("POSTGRES_DB", "cs"), "Set PostgreSQL DB.")
	flag.StringVar(&c.PostgresUser, "postgresUser", config.LookupEnvString("POSTGRES_USER", "cs"), "Set PostgreSQL user.")
	flag.StringVar(&c.PostgresPassword, "postgresPassword", config.LookupEnvString("POSTGRES_PASSWORD", "cs"), "Set PostgreSQL password.")
	flag.BoolVar(&c.PostgresSSLMode, "postgresSSLMode", config.LookupEnvBool("POSTGRES_SSL_MODE", false), "Set to enable PostgreSQL SSL mode.")
	flag.BoolVar(&c.PostgresSSLCheckCert, "postgresSSLCheckCert", config.LookupEnvBool("POSTGRES_SSL_CHECK_CERT", false), "Set to check PostgreSQL SSL cert.")
	flag.StringVar(&c.LogLevel, "logLevel", config.LookupEnvString("LOG_LEVEL", "TRACE"), "Set log level (DEBUG, INFO, WARN, ERROR, FATAL, any other value means TRACE).")
	flag.StringVar(&c.ListenGRPCAddr, "listenGRPCAddr", config.LookupEnvString("LISTEN_GRPC_ADDR", ":8000"), `Address in form of "[host]:port" that gRPC server should be listening on.`)
	flag.StringVar(&c.ListenHTTPAddr, "listenHTTPAddr", config.LookupEnvString("LISTEN_HTTP_ADDR", ":9000"), `Address in form of "[host]:port" that HTTP server should be listening on.`)
	flag.BoolVar(&c.TLS, "tls", config.LookupEnvBool("TLS", false), "Set to enable TLS.")
	flag.StringVar(&c.TokenKey, "tokenKey", config.LookupEnvString("TOKEN_KEY", ""), "Hex encoded token key to verify jwt token. Supported key sizes are 16, 24 and 32 bytes.")
	flag.BoolVar(&c.Auth, "auth", config.LookupEnvBool("AUTH", false), "Set to enable JWT auth.")
	flag.StringVar(&c.GopsAddr, "listenGopsAddr", config.LookupEnvString("LISTEN_GOPS_ADDR", "127.0.0.1:7000"), `Address in form of "[host]:port" that gops agent should be listening on. It's not safe to listen to interfaces other than loopback in production.`)
	flag.StringVar(&c.OwnCSURL, "ownCSURL", config.LookupEnvString("OWN_CS_URL", ""), "URL of current CS (http(s)://host[:port]).")
	flag.StringVar(&c.InstrumentationAddr, "listenInstrumentationAddr", config.LookupEnvString("LISTEN_INSTRUMENTATION_ADDR", ":9090"), `Address in form of "[host]:port" that instrumentation HTTP server should be listening on.`)
	flag.StringVar(&c.PathToPasswordSecList, "pathToPasswordSecList", config.LookupEnvString("PATH_TO_PASSWORD_SEC_LIST", "1000000-password-seclists.txt"), `nested path to file that context bad passwords`)
	flag.DurationVar(&c.AccessTokenExpiration, "accessTokenExpiration", config.LookupEnvDuration("ACCESS_TOKEN_EXPIRATION", 15*time.Minute), "how long access token will be alive. By default 15 minute")
	flag.DurationVar(&c.RefreshTokenExpiration, "refreshTokenExpiration", config.LookupEnvDuration("REFRESH_TOKEN_EXPIRATION", 6*24*time.Hour), "how long access token will be alive. By default 6 days")
	flag.StringVar(&c.AdminUsername, "administratorUsername", config.LookupEnvString("ADMINISTRATOR_USERNAME", ""), `predeclared administrator username`)
	flag.StringVar(&c.AdminPassword, "administratorPassword", config.LookupEnvString("ADMINISTRATOR_PASSWORD", ""), `predeclared administrator password`)
	flag.StringVar(&c.AdminEmail, "administratorEmail", config.LookupEnvString("ADMINISTRATOR_EMAIL", ""), `predeclared administrator email`)

	flag.Parse()

	return c
}
