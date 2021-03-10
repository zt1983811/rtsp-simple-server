package sourcerpicamera

import (
	"sync"
	"fmt"
	"time"

	"github.com/aler9/gortsplib"

	"github.com/aler9/rtsp-simple-server/internal/logger"
	"github.com/aler9/rtsp-simple-server/internal/stats"
)

// #cgo LDFLAGS: -ldl
// #include <dlfcn.h>
import "C"

const (
	retryPause = 5 * time.Second
)

// Parent is implemented by path.Path.
type Parent interface {
	Log(logger.Level, string, ...interface{})
	OnSourceSetReady(gortsplib.Tracks)
	OnSourceSetNotReady()
	OnFrame(int, gortsplib.StreamType, []byte)
}

// Source is a RPI Camera source.
type Source struct {
	wg          *sync.WaitGroup
	stats       *stats.Stats
	parent      Parent

	// in
	terminate chan struct{}
}

// New allocates a Source.
func New(wg *sync.WaitGroup,
	stats *stats.Stats,
	parent Parent) *Source {
	s := &Source{
		wg:          wg,
		stats:       stats,
		parent:      parent,
		terminate:   make(chan struct{}),
	}

	s.log(logger.Info, "started")

	s.wg.Add(1)
	go s.run()
	return s
}

// Close closes a Source.
func (s *Source) Close() {
	s.log(logger.Info, "stopped")
	close(s.terminate)
}

// IsSource implements path.source.
func (s *Source) IsSource() {}

// IsSourceExternal implements path.sourceExternal.
func (s *Source) IsSourceExternal() {}

func (s *Source) log(level logger.Level, format string, args ...interface{}) {
	s.parent.Log(level, "[rpi-camera source] "+format, args...)
}

func (s *Source) run() {
	defer s.wg.Done()

	for {
		ok := func() bool {
			ok := s.runInner()
			if !ok {
				return false
			}

			t := time.NewTimer(retryPause)
			defer t.Stop()

			select {
			case <-t.C:
				return true
			case <-s.terminate:
				return false
			}
		}()
		if !ok {
			break
		}
	}
}

func (s *Source) runInner() bool {
	s.log(logger.Info, "connecting")

	handle := C.dlopen(C.CString("/opt/vc/lib/libmmal_components.so"), C.RTLD_LAZY)

	fmt.Println(handle)

	s.log(logger.Info, "TODO")

	return true
}
