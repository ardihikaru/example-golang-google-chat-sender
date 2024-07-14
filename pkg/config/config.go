package config

type General struct {
	Market    string `mapstructure:"market"`
	BuildMode string `mapstructure:"buildMode"`
}

type Log struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type Google struct {
	ServiceAccount string   `mapstructure:"serviceAccount"`
	Spaces         []string `mapstructure:"spaces"`
}

type Scheduler struct {
	Schedules []string `mapstructure:"schedules"`
}
