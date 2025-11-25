package kafkaclient

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
)

const (
	OffsetFirst = "first"
	OffsetLast  = "last"
)

type KqConf struct {
	ServerName  string
	Brokers     []string
	Group       string
	Topic       string
	CaFile      string `json:",optional"`
	Offset      string `json:",options=first|last,default=last"`
	Conns       int    `json:",default=1"`
	Consumers   int    `json:",default=8"`
	Processors  int    `json:",default=8"`
	MinBytes    int    `json:",default=10240"`    // 10K
	MaxBytes    int    `json:",default=10485760"` // 10M
	Username    string `json:",optional"`
	Password    string `json:",optional"`
	ForceCommit bool   `json:",default=true"`
}

func CreateTopics(kafkaURL string, topicConfigs []kafka.TopicConfig) error {
	conn, err := kafka.Dial("tcp", kafkaURL)
	if err != nil {
		return fmt.Errorf("failed to dial kafka: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get controller: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return fmt.Errorf("failed to dial controller: %w", err)
	}
	defer controllerConn.Close()

	if err = controllerConn.CreateTopics(topicConfigs...); err != nil {
		return fmt.Errorf("failed to create topics: %w", err)
	}

	return nil
}

func ListTopic(kafkaURL string) (map[string]struct{}, error) {
	conn, err := kafka.Dial("tcp", kafkaURL)
	if err != nil {
		return nil, fmt.Errorf("failed to dial kafka: %w", err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return nil, fmt.Errorf("failed to read partitions: %w", err)
	}

	topics := make(map[string]struct{})
	for _, p := range partitions {
		topics[p.Topic] = struct{}{}
	}
	
	return topics, nil
}
