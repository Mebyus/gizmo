package vm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
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
	IP string `json:"ip"`

	// Always contains 64 elemetns.
	R []string `json:"r"`
}

func formatRegister(v uint64) string {
	return fmt.Sprintf("%016x", v)
}

func formatRegisters(vv []uint64) []string {
	if len(vv) == 0 {
		return nil
	}

	ss := make([]string, 0, len(vv))
	for _, v := range vv {
		ss = append(ss, formatRegister(v))
	}
	return ss
}

func (s *DebugServer) renderState(w http.ResponseWriter) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(&RenderStateObject{
		Registers: RenderStateRegisters{
			IP: formatRegister(s.vm.ip),
			R:  formatRegisters(s.vm.r[:]),
		},
	})
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
