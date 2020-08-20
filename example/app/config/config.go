package config

type Demo struct {
	Name string
}
type Configuration struct {
	Name string `value:"name"`
	A int `value:"a"`
}

func NewConfiguration() *Configuration {
	return &Configuration{}
}

func (c *Configuration) NewDemo() *Demo {
	return &Demo{c.Name}
}

func (c *Configuration) NewDemo1() (string,*Demo) {
	return "demo",&Demo{Name: "larry"}
}

