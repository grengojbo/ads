package config

type Config struct {
	Port int `default:"7000" env:"PORT"`
	DB   struct {
		Name     string `default:"qor_example"`
		Adapter  string `default:"postgres"`
		User     string
		Password string
		Host     string `default:"localhost"`
		Port     uint   `default:"5432"`
		Debug    bool   `default:"false"`
	}
	Redis struct {
		Host     string `default:"localhost"`
		Port     uint   `default:"6379"`
		Protocol string `default:"tcp"`
		Password string
	}
	Session struct {
		Name    string `default:"sessionid"`
		Adapter string `default:"cookie"`
	}
	Secret string `default:"secret"`
	Limit  int    `default:"5"`
}
