package certs

import (
	"github.com/eolinker/apinto/drivers"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/setting"
)

const driverName = "cert"

func Register(register eosc.IExtenderDriverRegister) {
	register.RegisterExtenderDriver(driverName, newFactory())
	setting.RegisterSetting(driverName, controller)
}

func newFactory() eosc.IExtenderDriverFactory {
	return &factory{IExtenderDriverFactory: drivers.NewFactory[Config](Create)}
}

type factory struct {
	eosc.IExtenderDriverFactory
}

func (f *factory) Create(profession string, name string, label string, desc string, params map[string]interface{}) (eosc.IExtenderDriver, error) {
	controller.driver = name
	controller.profession = profession
	return f.IExtenderDriverFactory.Create(profession, name, label, desc, params)
}
