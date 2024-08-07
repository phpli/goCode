package sarama

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
)

var addrs = []string{"localhost:9094"}

func TestSyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(addrs, cfg)

	assert.NoError(t, err)
	//message :=
	//messages := []*sarama.ProducerMessage{&sarama.ProducerMessage{
	//	Topic: "test-topic",
	//	Value: sarama.StringEncoder("hello, 这是一个消息"),
	//	Headers: []sarama.RecordHeader{
	//		{
	//			Key:   []byte("trace_id"),
	//			Value: []byte("123456"),
	//		},
	//	},
	//	Metadata: "则是metadata",
	//}}
	//err = producer.SendMessages(messages)

	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "test_topic",
		Key:   sarama.StringEncoder("oid-123"), //可以在业务内有序，都会分发到一个分区上
		Value: sarama.StringEncoder("hello, 这是一个消息2"),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("trace_id"),
				Value: []byte("123456"),
			},
		},
		Metadata: "则是metadata",
	})
	assert.NoError(t, err)
}

func TestAsyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Errors = true
	producer, err := sarama.NewAsyncProducer(addrs, cfg)
	assert.NoError(t, err)
	msgCh := producer.Input()
	go func() {
		for {

			msg := &sarama.ProducerMessage{
				Topic: "test_topic",
				Key:   sarama.StringEncoder("oid-123"), //可以在业务内有序，都会分发到一个分区上
				Value: sarama.StringEncoder("hello, 这是一个异步消息"),
				Headers: []sarama.RecordHeader{
					{
						Key:   []byte("trace_id"),
						Value: []byte("123456"),
					},
				},
				Metadata: "则是metadata",
			}
			select {
			case msgCh <- msg:
			default:
			}
		}
	}()
	errChan := producer.Errors()
	successChan := producer.Successes()
	for {
		select {
		case err := <-errChan:
			t.Log("发送出问题", err.Err)
		case <-successChan:
			t.Log("发送成功")
		}
		//time.Sleep(time)
	}

}
