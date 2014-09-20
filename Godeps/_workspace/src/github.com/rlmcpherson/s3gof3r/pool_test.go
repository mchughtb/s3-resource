package s3gof3r

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"time"
)

func TestBP(t *testing.T) {

	// send log output to buffer
	lf := *bytes.NewBuffer(nil)
	SetLogger(&lf, "", log.LstdFlags, true)

	bp := newBufferPool(mb)
	bp.timeout = 1 * time.Millisecond
	b := <-bp.get
	if cap(b.Bytes()) != int(mb+100*kb) {
		t.Errorf("Expected buffer capacity: %d. Actual: %d", kb, b.Len())
	}
	bp.give <- b
	if bp.makes != 2 {
		t.Errorf("Expected makes: %d. Actual: %d", 2, bp.makes)
	}

	b = <-bp.get
	bp.give <- b
	time.Sleep(2 * time.Millisecond)
	if bp.makes != 3 {
		t.Errorf("Expected makes: %d. Actual: %d", 3, bp.makes)
	}
	close(bp.quit)
	expLog := "3 buffers of 1 MB allocated"
	time.Sleep(1 * time.Millisecond) // wait for log
	ls := lf.String()
	if !strings.Contains(ls, expLog) {
		t.Errorf("BP debug logging on quit: \nExpected: %s\nActual: %s",
			expLog, ls)
	}

}
