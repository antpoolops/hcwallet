// +build darwin

package omnilib

// #include <stdio.h>
// #include <stdlib.h>
// #include "./omniproxy.h"
// #cgo CFLAGS: -I./
//#cgo LDFLAGS:-L./ -lomnicored -lbitcoin_server -lbitcoin_common -lunivalue -lbitcoin_util -lbitcoin_wallet  -lbitcoin_consensus -lbitcoin_crypto -lleveldb -lmemenv -lsecp256k1 /usr/local/lib/libboost_system.a /usr/local/lib/libboost_filesystem.a /usr/local/lib/libboost_program_options.a /usr/local/lib/libboost_thread-mt.a /usr/local/lib/libboost_chrono.a /usr/local/lib/libboost_iostreams.a /usr/local/lib/libdb_cxx.a /usr/local/Cellar/openssl/1.0.2p/lib/libssl.a /usr/local/Cellar/openssl/1.0.2p/lib/libcrypto.a  /usr/local/lib/libevent_pthreads.a /usr/local/lib/libevent.a -lm -lz -ldl -lstdc++
import "C"
import (
	"unsafe"
	"fmt"

	"sync"
	"time"
)

var mutexOmni sync.Mutex

func JsonCmdReqHcToOm(strReq string) string {
	mutexOmni.Lock()
	defer mutexOmni.Unlock()
	strRsp := C.GoString(C.CJsonCmdReq(C.CString(strReq)))
	return strRsp
}

func LoadLibAndInit() {
	C.CLoadLibAndInit()
}

func OmniStart(strArgs string, strArgs1 string) {
	C.COmniStart(C.CString(strArgs), C.CString(strArgs1))
}

var ChanReqOmToHc=make(chan string )
var ChanRspOmToHc=make(chan string )

// callback to LegacyRPC.Server
//var PtrLegacyRPCServer *Server=nil

//export JsonCmdReqOmToHc
func JsonCmdReqOmToHc(pcReq *C.char) *C.char {
	strReq:=C.GoString(pcReq)
	fmt.Println("Go JsonCmdReqOmToHc strReq=",strReq)
	ChanReqOmToHc<-strReq
	strRsp:=<-ChanRspOmToHc
	fmt.Println("Go JsonCmdReqOmToHc strRsp=",strRsp)
	cs := C.CString(strRsp)

	defer func(){
		go func() {
			time.Sleep(time.Microsecond*200)
			C.free(unsafe.Pointer(cs))
		}()
	}()

	return cs
}
