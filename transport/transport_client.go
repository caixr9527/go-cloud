package transport

type Fallback interface {
}
type Configuration interface {
}

type Client struct {
	Protocol      string
	Name          string
	Url           string
	Path          string
	Fallback      Fallback
	Configuration Configuration
}

type Service interface {
	Client() Client
	Register()
}
