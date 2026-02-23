package amqp

type AMQPConfig struct {
	Standalone bool   `yaml:"standalone"`
	URL        string `yaml:"url" validate:"required_unless=Standalone true"`
}
