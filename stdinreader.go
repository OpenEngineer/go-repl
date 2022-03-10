package terminal

import (
	"bufio"
	"os"
	"sync"
	"time"
)

const MACHINE_INTERVAL = time.Millisecond

// _StdinReader collects inputs
type _StdinReader struct {
	reader   *bufio.Reader
	lastTime time.Time
	buffer   []byte
	lock     *sync.Mutex

	bytes chan []byte
}

func newStdinReader() *_StdinReader {
	return &_StdinReader{
		reader:   nil,
		lastTime: time.Time{},
		buffer:   make([]byte, 0),
		lock:     &sync.Mutex{},

		bytes: make(chan []byte),
	}
}

func (r *_StdinReader) start() {
	go func() {
		for {
			<-time.After(MACHINE_INTERVAL)

			r.lock.Lock()

			if len(r.buffer) > 0 {
				if time.Now().After(r.lastTime.Add(MACHINE_INTERVAL)) {
					msg := r.buffer

					r.buffer = make([]byte, 0)

					r.bytes <- msg
				}
			}

			r.lock.Unlock()
		}
	}()
}

func (r *_StdinReader) read() {
	if r.reader != nil {
		return
	}

	r.reader = bufio.NewReader(os.Stdin)
	r.lastTime = time.Now()

	go func() {
		for {
			b, err := r.reader.ReadByte()
			if err != nil {
				panic(err)
			}

			stopNow := false
			if b == 13 && time.Now().After(r.lastTime.Add(MACHINE_INTERVAL)) {
				// it is unlikely that a carriage return followed by some text is pasted into the terminal, so we can use this as a queu to quit
				stopNow = true
			}

			r.lastTime = time.Now()

			r.lock.Lock()

			r.buffer = append(r.buffer, b)

			r.lock.Unlock()

			if stopNow {
				r.reader = nil
				return
			}
		}
	}()
}
