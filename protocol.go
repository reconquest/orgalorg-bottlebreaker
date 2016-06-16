package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/seletskiy/hierr"
)

type (
	syncMode int
)

const (
	syncModeAll syncMode = iota
	syncModeUnknown
)

type protocolMessageFilter struct {
	tag   string
	parts int
}

func receiveProtocolPrefix(reader *bufio.Reader) (string, error) {
	message, err := receiveExpectedMessage(
		reader,
		[]protocolMessageFilter{{"HELLO", 2}},
	)
	if err != nil {
		return "", err
	}

	return message[0], nil
}

func receiveNodesList(reader *bufio.Reader) ([]string, error) {
	nodes := []string{}

	for {
		message, err := receiveExpectedMessage(
			reader,
			[]protocolMessageFilter{{"NODE", 3}, {"START", 2}},
		)
		if err != nil {
			return nil, err
		}

		if message[1] == "START" {
			break
		}

		nodes = append(nodes, message[2])
	}

	return nodes, nil
}

func synchronize(
	reader *bufio.Reader,
	writer io.Writer,
	phase string,
	prefix string,
	nodes []string,
	mode syncMode,
) error {
	fmt.Fprintf(writer, "%s SYNC %s\n", prefix, phase)

	threshold := 0
	switch mode {
	default:
		threshold = len(nodes)
	}

	for i := 0; i < threshold; i++ {
		_, err := receiveExpectedMessage(
			reader,
			[]protocolMessageFilter{{"SYNC", 4}},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func receiveExpectedMessage(
	reader *bufio.Reader,
	filters []protocolMessageFilter,
) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, hierr.Errorf(
			err,
			`can't receive one of messages: %v`,
			filters,
		)
	}

	message := []string{}
	for _, filter := range filters {
		message = strings.SplitN(strings.TrimSpace(line), " ", filter.parts)
		if len(message) < 2 {
			return nil, hierr.Errorf(
				err,
				`invalid message received: %v`,
				message,
			)
		}

		if message[1] != filter.tag {
			continue
		}

		if len(message) < filter.parts {
			return nil, hierr.Errorf(
				err,
				`short %s message received: %d of %d parts`,
				filter.tag,
				len(message), filter.parts,
			)
		}

		return message, nil
	}

	return nil, hierr.Errorf(
		err,
		`unexpected message %s received, expected %v`,
		message[1], filters,
	)
}
