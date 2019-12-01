package goKLC

type Config struct {
	key   string
	value interface{}
}

type configCollector map[string]Config

func newConfigCollector() *configCollector {

	return &configCollector{}
}

func NewConfig() Config {

	return Config{}
}

func (c Config) Set(key string, value interface{}) {
	c.key = key
	c.value = value

	_configCollector.Set(key, c)
}

func (c Config) Get(key string, defaultValue interface{}) interface{} {
	value, found := _configCollector.Get(key)

	if !found {

		return defaultValue
	}

	return value.(Config).value
}

func (c Config) SetFromMap(config map[string]interface{}) {
	for key, value := range config {
		newConfig := NewConfig()
		newConfig.Set(key, value)
	}
}

func (cc *configCollector) Set(key string, config Config) {

	(*cc)[key] = config
}

func (cc *configCollector) Get(key string) (interface{}, bool) {
	value, found := (*cc)[key]

	return value, found
}
