package prompt

import "io"

//IOWriter wrap io.writer to ConsoleWriter
type IOWriter struct {
	VT100Writer
	out io.Writer
}

// Flush to flush buffer
func (w *IOWriter) Flush() error {
	_, err := w.out.Write(w.buffer)
	if err != nil {
		return err
	}
	w.buffer = []byte{}
	return nil
}

// NewIOWriter returns ConsoleWriter object to write to stdout.
// This generates win32 control sequences.
func NewIOWriter(writer io.Writer) ConsoleWriter {
	return &IOWriter{
		out: writer,
	}
}
