package vinna

import (
	"errors"
	"io"
	"time"

	"golang.org/x/net/context"

	"github.com/akresling/vinna/pb"
)

var (
	// ErrTopicAlreadyExists is the error if the topic already exists
	ErrTopicAlreadyExists = errors.New("topic already exists")
	// ErrTopicDoesNotExist is returned if the topic doesnt exist
	ErrTopicDoesNotExist = errors.New("topic does not exist")
)

var (
	topicSize   = 10000000
	delayedSize = 100000
)

// Vinna is the implementation of the vinna service
type Vinna struct {
	topics       map[string]chan *pb.Message
	blockedRetry chan *pb.Message
	exit         chan bool
}

// New will initialize and return a new vinna service structure
func New() Vinna {
	v := Vinna{
		topics:       make(map[string]chan *pb.Message),
		blockedRetry: make(chan *pb.Message, delayedSize),
		exit:         make(chan bool),
	}
	return v
}

// NewTopic will create a topic to pump data into
func (v Vinna) NewTopic(_ context.Context, tp *pb.Topic) (*pb.TopicCreated, error) {
	if _, ok := v.topics[tp.GetTopic()]; ok {
		return &pb.TopicCreated{Success: false}, ErrTopicAlreadyExists
	}
	topic := make(chan *pb.Message, topicSize)
	v.topics[tp.GetTopic()] = topic
	return &pb.TopicCreated{Success: true}, nil
}

// Add will add things into the vinna service
func (v Vinna) Add(_ context.Context, msg *pb.Message) (*pb.Success, error) {
	err := v.put(msg)
	if err != nil {
		return &pb.Success{Success: false}, err
	}
	return &pb.Success{Success: true}, nil
}

// Take will take things from the vinna service
func (v Vinna) Take(_ context.Context, req *pb.MessageRequest) (*pb.Message, error) {
	tName := req.GetTopic()
	msg, err := v.get(tName)
	if err != nil {
		return &pb.Message{}, err
	}
	return msg, nil
}

// Produce will open a stream for client side message sending
func (v Vinna) Produce(stream pb.Vinna_ProduceServer) error {
	var msgCount int32
	startTime := time.Now()
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			elapsedTime := int32(time.Now().Sub(startTime).Seconds())
			return stream.SendAndClose(&pb.ProducerSummary{
				MsgCount:    msgCount,
				ElapsedTime: elapsedTime,
			})
		}
		if err != nil {
			return err
		}
		msgCount++
		err = v.put(msg)
		if err != nil {
			return err
		}
	}
}

// Consume will be the service endpoint for client side consuming
func (v Vinna) Consume(cr *pb.ConsumeRequest, stream pb.Vinna_ConsumeServer) error {
	tName := cr.GetTopic()
	for {
		msg, err := v.get(tName)
		if err != nil {
			return err
		}
		stream.Send(msg)
	}
}

func (v Vinna) get(tName string) (*pb.Message, error) {
	topic, ok := v.topics[tName]
	if !ok {
		return &pb.Message{}, ErrTopicDoesNotExist
	}
	msg := <-topic
	return msg, nil
}

func (v Vinna) put(msg *pb.Message) error {
	tName := msg.GetTopic()
	topic, ok := v.topics[tName]
	if !ok {
		return ErrTopicDoesNotExist
	}
	topic <- msg
	return nil
}
