package main

import (
	"io"

	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"

	"github.com/zperf/nfs/v3/client/nfs"
	"github.com/zperf/nfs/v3/client/nfs/rpc"
	"github.com/zperf/nfs/v3/logging"
)

func main() {
	_ = logging.New()
	log := zap.S()

	//m, err := nfs.DialMountWithPort("192.168.126.15", 3456)
	m, err := nfs.DialMountWithPort("127.0.0.1", 3456)
	if err != nil {
		log.Fatalf("dial failed, %v", err)
	}
	defer m.Close()

	t, err := m.Mount("/", rpc.AuthNull)
	if err != nil {
		log.Fatalf("Get mount failed, %v", err)
	}
	defer t.Close()

	dirs, err := t.ReadDirPlus(".")
	if err != nil {
		log.Fatalf("ReadDirPlus: %v", err)
	}

	for _, dir := range dirs {
		log.Infof("%+v", dir.Name())
	}

	bar := progressbar.DefaultBytes(-1, "Reading")
	for {
		// read the test-file1 forever
		fh, err := t.Open("test-file1")
		if err != nil {
			log.Fatalf("open failed, %v", err)
		}
		defer fh.Close()

		_, err = io.Copy(bar, fh)
		if err != nil {
			log.Warn("copy failed", err)
		}
	}
}
