package infra

import (
	"bytes"
	"errors"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func TestResultRecorder_completed_always_logSuccessAndFailureCount(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	sut := ResultRecorder{}
	sut.Completed(&messaging.BatchResponse{
		SuccessCount: 1,
		FailureCount: 2,
		Responses:    []*messaging.SendResponse{},
	})

	readByte()

	loggedMessage := buf.String()

	assert.True(t, strings.Contains(loggedMessage, "success count: 1"))
	assert.True(t, strings.Contains(loggedMessage, "failure count: 2"))
}

func TestResultRecorder_completed_someMessageFailed_logThem(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	responses := []*messaging.SendResponse{
		{
			Success:   false,
			MessageID: "MSG1",
			Error:     errors.New(""),
		},
		{
			Success:   true,
			MessageID: "MSG2",
			Error:     nil,
		},
		{
			Success:   false,
			MessageID: "MSG3",
			Error:     errors.New(""),
		},
	}

	sut := ResultRecorder{}
	sut.Completed(&messaging.BatchResponse{
		FailureCount: 1,
		Responses:    responses,
	})

	readByte()
	loggedMessage := buf.String()

	assert.True(t, strings.Contains(loggedMessage, "messageId: MSG1 failed"))
	assert.True(t, strings.Contains(loggedMessage, "messageId: MSG3 failed"))
	assert.False(t, strings.Contains(loggedMessage, "MSG2"))

}

func readByte() {
	err := io.EOF
	if err != nil {
		fmt.Println("ERROR")
		log.Print("Couldn't read first byte")
		return
	}
}
