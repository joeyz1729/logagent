package kafka

import (
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

var (
	client  sarama.SyncProducer
	msgChan chan *Message
	log     *logrus.Logger
)

type Message struct {
	Topic string
	Data  string
}

func Init(addrs []string, chanSize int) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		log.Errorf("producer closed, err: %s", err)
		return err
	}

	msgChan = make(chan *Message, chanSize)
	go sendKafka()

	return
}

func sendKafka() {

}

func SendLog(msg *Message) (err error) {
	return
}
