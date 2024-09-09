package main

import (
	"errors"
	"fmt"
)

const (
	ArrayPrefix = '*'
	BulkPrefix  = '$'
)

func parseCommand(buf []byte) ([]string, error) {
	i := 0
	if i >= len(buf) {
		return nil, fmt.Errorf("invalid command")
	}
	if buf[i] != ArrayPrefix {
		return nil, fmt.Errorf("invalid command format")
	}
	i++
	var length int
	for i < len(buf) && buf[i] >= '0' && buf[i] <= '9' {
		length = length*10 + int(buf[i]-'0')
		i++
	}
	i, err := expect(buf, i, "\r\n")
	if err != nil {
		return nil, err
	}

	var args []string
	var arg string
	for j := 0; j < length; j++ {
		arg, i, err = parseBulkString(buf, i)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	return args, nil
}

func parseBulkString(buf []byte, i int) (string, int, error) {
	if i >= len(buf) || buf[i] != BulkPrefix {
		return "", i, errors.New("expecting $")
	}
	i++
	var length int
	for i < len(buf) && buf[i] >= '0' && buf[i] <= '9' {
		length = length*10 + int(buf[i]-'0')
		i++
	}
	i, err := expect(buf, i, "\r\n")
	if err != nil {
		return "", i, err
	}
	if i+length > len(buf) {
		return "", i, fmt.Errorf("buffer too short for the expected length")
	}
	bulkString := string(buf[i : i+length])
	i += length

	i, err = expect(buf, i, "\r\n")
	if err != nil {
		return "", i, err
	}
	return bulkString, i, nil
}

func expect(buf []byte, i int, exp string) (int, error) {
	if i+len(exp) <= len(buf) && string(buf[i:i+len(exp)]) == exp {
		return i + len(exp), nil
	}
	return i, errors.New("Expecting " + exp)
}
