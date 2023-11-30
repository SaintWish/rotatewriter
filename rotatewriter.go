// Credits to https://stackoverflow.com/a/28797984 for
// the basis of RotateWriter struct.

// RotateWriter interfaces to io.Writer to automaticly rotate the
// log file based on max size and/or time
package rotatewriter

import (
	"os"
	"sync"
	"time"
	"strings"
	"path/filepath"
)

var (
	Megabyte int64 = 1024 * 1024 // Variable to use along with MaxSize to use megabytes.
	Kilobyte int64 = 1024 // Variable to use along with MaxSize to use kilobytes.
)

// The structure for RotateWriter, which should interface io.Writer
type RotateWriter struct {
	Dir string // the directory to put log files.
	Filename string // should be set to the actual filename and extension.
	ExpireTime time.Duration // how often the log should rotate.
	MaxSize int64 // max size a log file is allowed to be in bytes.

	mu sync.Mutex
	now time.Time
	fp *os.File
}

// Write satisfies the io.Writer interface.
func (w *RotateWriter) Write(output []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.fp == nil {
		if err := w.Resume(); err != nil {
			return 0, err
		}
	}

	fi, _ := os.Stat(w.Dir+w.Filename);
	if (w.MaxSize > 0) && (fi.Size() >= w.MaxSize) {
		if err := w.Rotate(); err != nil {
			return 0, err
		}
	}

	if (w.ExpireTime > 0) && (time.Now().After(w.now.Add(w.ExpireTime))) {
		if err := w.Rotate(); err != nil {
			return 0, err
		}
	}

	return w.fp.Write(output)
}

// Create a new log file does not exists, opens if log file does exists.
func (w *RotateWriter) Resume() error {
	var err error
	var filename = w.Dir+w.Filename

	w.fp, err = os.OpenFile(filename, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	fi, _ := os.Stat(filename);
	w.now = fi.ModTime()
	return nil
}

// Perform the actual act of rotating and reopening file.
func (w *RotateWriter) Rotate() error {
	var err error
	var filename = w.Dir+w.Filename

	// create the needed direactories if they don't exists.
	if err := os.MkdirAll(w.Dir, 0755); err != nil {
		return err
	}

	// Close existing file if open
	if w.fp != nil {
		if err := w.fp.Close(); err != nil {
			return err
		}
		w.fp = nil
	}

	// Rename dest file if it already exists
	if _, err := os.Stat(filename); err == nil {
		if err := w.renameFile(); err != nil {
			return err
		}
	}

	// Create a file.
	w.now = time.Now()
	w.fp, err = os.OpenFile(filename, os.O_WRONLY | os.O_CREATE, 0666)
	return err
}

// Rename the log file to include the current date. Uses RFC3339 time format.
func (w *RotateWriter) renameFile() error {
	var fn = w.Dir+w.Filename
	newfn := fn[:len(fn)-len(filepath.Ext(w.Filename))]+"-"+time.Now().Format(time.RFC3339)+filepath.Ext(w.Filename)

	return os.Rename(fn, strings.ReplaceAll(newfn, ":", "-"))
}