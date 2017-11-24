package amqpClient

//Config struct : describes the structure of the .yml file configuration
type Config struct {
	RabbitMq RabbitMq
}

//RabbitMq struct : describes the rabbitmq configuration (connection, loggig settings)
type RabbitMq struct {
	Connection struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Uri      string `yaml:"uri"`
		Exchange string `yaml:"exchange"`
		Key string `yaml:"key"`
	}
	Logs struct {
		Logfile string `yaml:"logfile"`
	}
}
