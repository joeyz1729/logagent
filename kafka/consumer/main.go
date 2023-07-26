package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"sync"
)

func main() {
	// 创建consumer
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err: %v\n", err)
		return
	}
	fmt.Println("start consumer success")

	// 返回给定主题的所有分区列表
	partitionList, err := consumer.Partitions("web_log")
	if err != nil {
		fmt.Printf("fail to get list of partiion: err: %v\n", err)
		return
	}
	fmt.Println("get partition list success")
	fmt.Println(partitionList)

	// 等待接收message
	var wg sync.WaitGroup
	for partition := range partitionList {
		// 每个分区创建PartitionConsumer
		pc, err := consumer.ConsumePartition("web_log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d, err: %v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		wg.Add(1)
		// 创建goroutine接收msg
		go func(sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				fmt.Printf("Partition: %d Offset: %d  Value: %s\n", msg.Partition, msg.Offset, msg.Value)
			}
		}(pc)
	}
	wg.Wait()
}
