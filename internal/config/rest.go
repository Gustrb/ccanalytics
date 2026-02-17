package config

type RestConfig struct {
	Addr string `envconfig:"REST_ADDR" default:":8080"`
}

var Rest RestConfig

func init() {
	if err := Config(&Rest); err != nil {
		panic(err)
	}
}
