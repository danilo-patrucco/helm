package support

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

const outFilePath = "ignorer-inputs.jsonl"

type DumpMessage struct {
	Kind     string `json:"kind"`
	Index    int    `json:"index,omitempty"`
	Severity int    `json:"severity"`
	Path     string `json:"path"`
	ErrText  string `json:"err_text"`
}

type DumpError struct {
	Kind    string `json:"kind"`
	Index   int    `json:"index,omitempty"`
	ErrText string `json:"err_text"`
}

// DumpInputsAsJsonLines takes slices of Messages and errors and dumps them
// to outFilePath as a series of single-line JSON payloads.
func DumpInputsAsJsonLines(messages []Message, errors []error) {
	f, closer := openOutputFile(outFilePath)
	defer closer()
	dumpMessages(f, messages...)
	dumpErrors(f, errors...)
	f.Sync()
}

func dumpErrors(w io.Writer, errors ...error) {
	for i, err := range errors {
		de := DumpError{Kind: "error", Index: i, ErrText: err.Error()}
		jsonBytes, err := json.Marshal(de)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s\n", jsonBytes)
	}

}

func dumpMessages(w io.Writer, messages ...Message) {
	for i, message := range messages {
		dm := DumpMessage{
			Kind:     "message",
			Index:    i,
			Severity: message.Severity,
			Path:     message.Path,
			ErrText:  message.Err.Error(),
		}
		jsonBytes, err := json.Marshal(dm)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s\n", jsonBytes)
	}
}

// - opens an output file
// - creates it if it doesn't exist already
// - will append to the file rather than truncating
func openOutputFile(p string) (*os.File, func()) {
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	closeFile := func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
	closeFn := func() { closeFile(f) }

	return f, closeFn
}
