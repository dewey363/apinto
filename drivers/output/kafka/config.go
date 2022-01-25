package kafka

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/eolinker/eosc"
	"strings"
	"time"
)

var (
	errTopic           = errors.New("topic can not be null. ")
	errAddress         = errors.New("address is invalid. ")
	errorFormatterType = errors.New("error formatter type")
	errorPartitionKey  = errors.New("partition key is invalid")
)

type Config struct {
	Config *Kafka `json:"config" yaml:"config"`
}
type Kafka struct {
	Topic         string               `json:"topic" yaml:"topic"`
	Address       string               `json:"address" yaml:"address"`
	Timeout       int                  `json:"timeout" yaml:"timeout"`
	Version       string               `json:"version" yaml:"version"`
	PartitionType string               `json:"partition_type" yaml:"partition_type"`
	Partition     int32                `json:"partition" yaml:"partition"`
	PartitionKey  string               `json:"partition_key" yaml:"partition_key"`
	Type          string               `json:"type" yaml:"type"`
	Formatter     eosc.FormatterConfig `json:"formatter" yaml:"formatter"`
}

type ProducerConfig struct {
	Address       []string             `json:"address" yaml:"address"`
	Topic         string               `json:"topic" yaml:"topic"`
	Partition     int32                `json:"partition" yaml:"partition"`
	PartitionKey  string               `json:"partition_key" yaml:"partition_key"`
	PartitionType string               `json:"partition_type" yaml:"partition_type"`
	Conf          *sarama.Config       `json:"conf" yaml:"conf"`
	Type          string               `json:"type" yaml:"type"`
	Formatter     eosc.FormatterConfig `json:"formatter" yaml:"formatter"`
}

func (c *Config) doCheck() (*ProducerConfig, error) {
	conf := c.Config
	if conf.Topic == "" {
		return nil, errTopic
	}
	if conf.Address == "" {
		return nil, errAddress
	}
	p := &ProducerConfig{}
	s := sarama.NewConfig()
	if conf.Version != "" {
		v, err := sarama.ParseKafkaVersion(conf.Version)
		if err != nil {
			return nil, err
		}
		s.Version = v
	}
	p.PartitionType = conf.PartitionType
	switch conf.PartitionType {
	case "robin":
		s.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	case "hash":
		// 通过hash获取
		if conf.PartitionKey == "" {
			// key为空则还是用随机
			s.Producer.Partitioner = sarama.NewRandomPartitioner
			p.PartitionType = "random"
		} else {
			if !strings.HasPrefix(conf.PartitionKey, "$") {
				return nil, errorPartitionKey
			}
			s.Producer.Partitioner = sarama.NewHashPartitioner
			p.PartitionKey = strings.TrimLeft(conf.PartitionKey, "$")
		}
	case "manual":
		// 手动指定分区
		s.Producer.Partitioner = sarama.NewManualPartitioner
		// 默认为0
		p.Partition = conf.Partition
	default:
		s.Producer.Partitioner = sarama.NewRandomPartitioner
		p.PartitionType = "random"
	}
	// 只监听错误
	s.Producer.Return.Errors = true
	s.Producer.Return.Successes = false
	s.Producer.RequiredAcks = sarama.WaitForLocal

	p.Address = strings.Split(conf.Address, ",")
	if len(p.Address) == 0 {
		return nil, errAddress
	}
	// 超时时间
	if conf.Timeout != 0 {
		s.Producer.Timeout = time.Duration(conf.Timeout) * time.Second
	}

	if conf.Type == "" {
		conf.Type = "line"
	}
	p.Type = conf.Type
	p.Formatter = conf.Formatter
	p.Conf = s
	return p, nil
}
