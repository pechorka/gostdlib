package env

import (
	"bytes"
	"os"

	"github.com/pechorka/gostdlib/pkg/errs"
	"github.com/pechorka/gostdlib/pkg/stringx"
)

func ExportDotEnv() error {
	file, err := os.ReadFile(".env")
	if err != nil {
		return errs.Wrap(err, "failed to read .env file")
	}
	if err := exportDotEnv(file); err != nil {
		return errs.Wrap(err, "failed to export .env file")
	}

	return nil
}

func exportDotEnv(file []byte) error {
	lines := bytes.Split(file, []byte("\n"))
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 || bytes.HasPrefix(line, []byte("#")) {
			continue
		}
		name, value, ok := bytes.Cut(line, []byte("="))
		if !ok {
			return errs.Newf("invalid line: %s", line)
		}
		value = removeCommentsAndSpaces(value)
		err := os.Setenv(stringx.FromBytes(name), stringx.FromBytes(value))
		if err != nil {
			return errs.Wrap(err, "failed to set environment variable")
		}
	}

	return nil
}

func removeCommentsAndSpaces(value []byte) []byte {
	value = bytes.TrimSpace(value)
	if len(value) == 0 {
		return nil
	}

	if len(value) > 1 && value[0] == '\'' || value[0] == '"' {
		quoteEnd := bytes.LastIndexByte(value, value[0])
		if quoteEnd != -1 {
			return value[1:quoteEnd]
		}
	}

	i := bytes.IndexByte(value, '#')
	if i != -1 {
		return bytes.TrimRight(value[:i], " ")
	}
	return value
}
