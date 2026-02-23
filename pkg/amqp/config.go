package amqp

type AMQPConfig struct {
	Standalone bool       `yaml:"standalone"`
	URL        string     `yaml:"url" validate:"required_unless=Standalone true"`
	Pool       PoolConfig `yaml:"pool"`
}

type PoolConfig struct {
	InitialSize int `yaml:"initialSize"`
	MaxSize     int `yaml:"maxSize"`
}

func (c PoolConfig) withDefaults() PoolConfig {
	if c.InitialSize <= 0 {
		c.InitialSize = 1
	}
	if c.MaxSize <= 0 {
		c.MaxSize = 5
	}
	if c.MaxSize < c.InitialSize {
		c.MaxSize = c.InitialSize
	}
	return c
}
