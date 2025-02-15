package limiting

import (
	"github.com/eolinker/apinto/drivers"
	"github.com/eolinker/apinto/resources"
	"github.com/eolinker/eosc"
)

type Config struct {
	Cache eosc.RequireId `json:"cache" skill:"github.com/eolinker/apinto/resources.resources.ICache" required:"false" label:"缓存位置"`
}

func Create(id, name string, cfg *Config, workers map[eosc.RequireId]eosc.IWorker) (eosc.IWorker, error) {

	return &Strategy{
		
		WorkerBase: drivers.Worker(id, name),
		buildProxy: resources.NewVectorBuilder(string(cfg.Cache)),
	}, nil
}
