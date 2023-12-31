package tracer

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestNewMessageStack(t *testing.T) {

	logger := logrus.New()

	f := &logrus.TextFormatter{
		DisableTimestamp: true,
		ForceColors:      true,
	}

	db, err := InitDatabase("rx.db", &Tracer{}, &Message{})
	if err != nil {
		// Handle the error
		fmt.Printf("Error initializing database: %v\n", err)
	}

	logger.SetFormatter(f)

	stack, err := NewStack(debug, logger)
	if err != nil {
		return
	}

	tracer := NewTracer("my-plugin", "", stack.Logger(), db)

	err = tracer.AddTracer("kkk", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	tracer.Info("info !!!!1")
	tracer.Error("error #####")

	tracer.SaveMessagesToDB(10)
	//
	tracers, _ := tracer.GetAllTracers()
	for _, trace := range tracers {
		for i2, message := range trace.Messages {
			fmt.Println(i2+1, message.UUID, message.Text)
		}
	}

}
