package kafka

import (
	"fmt"
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
	Key   string
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
	for {
		select {
		case msg := <-msgChan:
			producerMessage := &sarama.ProducerMessage{
				Topic: msg.Topic,
				Value: sarama.StringEncoder(msg.Data),
			}
			pid, offset, err := client.SendMessage(producerMessage)
			if err != nil {
				logrus.Warning("send msg failed, err: ", err)
				return
			}
			logrus.Infof("send msg to kafka success, pid: %v offset: %v", pid, offset)
		}
	}

}

func SendLog(msg *Message) (err error) {
	select {
	case msgChan <- msg:
	default:
		err = fmt.Errorf("msgChan is full")
	}
	return
}
