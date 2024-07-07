package vm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type DebugServer struct {
	vm Machine

	hs http.Server

	// root directory where static files
	// with html, js and other web stuff if located
	dir string
}

func NewDebugServer(addr string, dir string) *DebugServer {
	s := &DebugServer{
		hs: http.Server{
			Addr: addr,
		},
		dir: dir,
	}
	s.hs.Handler = s
	return s
}

func (s *DebugServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%-10s%s\n", r.Method, r.RequestURI)

	p := r.URL.Path
	if strings.HasPrefix(p, "/vm/") {
		s.serve(w, r)
		return
	}

	path := filepath.Join(s.dir, p)
	http.ServeFile(w, r, path)
}

// serve vm command or query
func (s *DebugServer) serve(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/vm/")

	switch name {
	case "state":
		switch r.Method {
		case http.MethodGet:
			s.renderState(w)
			return
		}
	case "step":
		switch r.Method {
		case http.MethodPost:
			s.renderStep(w)
			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write(fmt.Appendf(nil, "unknown command: %s (%s)", r.Method, name))
}

// RenderStateObject is used to marshal vm state into json encoding.
type RenderStateObject struct {
	Registers RenderStateRegisters `json:"regs"`

	Exit *RenderStateExit `json:"exit,omitempty"`
}

type RenderStateExit struct {
	Error  string `json:"error,omitempty"`
	Status string `json:"status"`
	Clock  string `json:"clock"`
}

// All registers are rendered as fixed width hex integers.
type RenderStateRegisters struct {
	IP RenderStateRegister `json:"ip"`

	// Always contains 64 elemetns.
	R []RenderStateRegister `json:"r"`
}

type RenderStateRegister [8]string

func formatByte(b byte) string {
	if b < 0x10 {
		return "0" + strconv.FormatUint(uint64(b), 16)
	}
	return strconv.FormatUint(uint64(b), 16)
}

func formatRegister(r *RenderStateRegister, v uint64) {
	for i := 0; i < 8; i += 1 {
		b := byte(v)
		r[7-i] = formatByte(b)
		v = v >> 8
	}
}

func formatRegisters(vv []uint64) []RenderStateRegister {
	if len(vv) == 0 {
		return nil
	}

	ss := make([]RenderStateRegister, len(vv))
	for i, v := range vv {
		r := &ss[i]
		formatRegister(r, v)
	}
	return ss
}

func (s *DebugServer) state() *RenderStateObject {
	var obj RenderStateObject
	obj.Registers.R = formatRegisters(s.vm.r[:])
	formatRegister(&obj.Registers.IP, s.vm.ip)
	return &obj
}

func (s *DebugServer) renderState(w http.ResponseWriter) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(s.state())
	if err != nil {
		fmt.Printf("encode state error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(fmt.Appendf(nil, "encode state error: %v", err))
		return
	}
}

func (s *DebugServer) renderStep(w http.ResponseWriter) {

}

func (s *DebugServer) ListenAndServe() error {
	fmt.Printf("serving web ui from %s\n", s.dir)
	fmt.Printf("open your browser at %s\n", s.hs.Addr)
	return s.hs.ListenAndServe()
}
