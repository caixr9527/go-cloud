package config

type client struct {
	// timeout for requesting Nacos server, default value is 10000ms
	TimeoutMs uint64 `yaml:"timeoutMs" mapstructure:"timeoutMs"`
	// the namespaceId of Nacos
	NamespaceId string `yaml:"namespaceId" mapstructure:"namespaceId"`
	// the endpoint for ACM. https://help.aliyun.com/document_detail/130146.html
	Endpoint string `yaml:"endpoint" mapstructure:"endpoint"`
	// the regionId for ACM & KMS
	RegionId string `yaml:"regionId" mapstructure:"regionId"`
	// the AccessKey for ACM & KMS
	AccessKey string `yaml:"accessKey" mapstructure:"accessKey"`
	// the SecretKey for ACM & KMS
	SecretKey string `yaml:"secretKey" mapstructure:"secretKey"`
	// it's to open KMS, default is false. https://help.aliyun.com/product/28933.html
	// , to enable encrypt/decrypt, DataId should be start with "cipher-"
	OpenKMS bool `yaml:"openKMS" mapstructure:"openKMS"`
	// the directory for persist nacos service info,default value is current path
	CacheDir string `yaml:"cacheDir" mapstructure:"cacheDir"`
	// the number of goroutine for update nacos service info,default value is 20
	UpdateThreadNum int `yaml:"updateThreadNum" mapstructure:"cacheDir"`
	// not to load persistent nacos service info in CacheDir at start time
	NotLoadCacheAtStart bool `yaml:"notLoadCacheAtStart" mapstructure:"notLoadCacheAtStart"`
	// update cache when get empty service instance from server
	UpdateCacheWhenEmpty bool `yaml:"updateCacheWhenEmpty" mapstructure:"updateCacheWhenEmpty"`
	// the username for nacos auth
	Username string `yaml:"username" mapstructure:"username"`
	// the password for nacos auth
	Password string `yaml:"password" mapstructure:"password"`
	// the directory for log, default is current path
	LogDir string `yaml:"logDir" mapstructure:"logDir"`
	// the level of log, it's must be debug,info,warn,error, default value is info
	LogLevel string `yaml:"logLevel" mapstructure:"logLevel"`
}

type server struct {
	// the nacos server scheme,defaut=http,this is not required in 2.0
	Scheme string `yaml:"scheme" mapstructure:"scheme"`
	// the nacos server contextpath,defaut=/nacos,this is not required in 2.0
	ContextPath string `yaml:"contextPath" mapstructure:"contextPath"`
	// the nacos server address
	IpAddr string `yaml:"ipAddr" mapstructure:"ipAddr"`
	// nacos server port
	Port uint64 `yaml:"port" mapstructure:"port"`
	// nacos server grpc port, default=server port + 1000, this is not required
	GrpcPort uint64 `yaml:"grpcPort" mapstructure:"grpcPort"`
}

type discoverConfig struct {
	EnableDiscover bool     `yaml:"enableDiscover"  mapstructure:"enableDiscover"`
	EnableConfig   bool     `yaml:"enableConfig" mapstructure:"enableConfig"`
	Client         client   `yaml:"client" mapstructure:"client"`
	Server         []server `yaml:"server" mapstructure:"server"`
}
