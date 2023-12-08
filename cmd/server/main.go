package main

import (
	"net"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/willscott/go-nfs"
	nfshelper "github.com/willscott/go-nfs/helpers"
	"go.uber.org/zap"

	"github.com/zperf/nfs/v3/logging"
)

func main() {
	logger := logging.New()
	nfs.SetLogger(&logging.MyLogger{Z: logger})

	listener, err := net.Listen("tcp", ":3456")
	if err != nil {
		zap.S().Fatalf("listen failed, %v", err)
	}
	zap.S().Infof("Server running at %s", listener.Addr())

	local := osfs.New("testdir")
	handler := nfshelper.NewNullAuthHandler(local)
	cache := nfshelper.NewCachingHandler(handler, 8192)
	err = nfs.Serve(listener, cache)
	if err != nil {
		zap.S().Infof("Serving NFS, %v", err)
	}
}
