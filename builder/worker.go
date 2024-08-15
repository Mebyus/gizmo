package builder

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/mebyus/gizmo/gencpp"
	"github.com/mebyus/gizmo/parser"
	"github.com/mebyus/gizmo/source"
)

type BuildTask struct {
	dep *DepEntry

	// order number of this task, used to produce deterministic
	// build results when gathering outputs from multiple workers
	order int

	// scan this unit for main function
	scanMain bool
}

type BuildTaskResult struct {
	task *BuildTask

	// output generated by build
	genout []byte

	// Additional object file built from unit asm files
	//
	// If there are no asm files in unit, then this field remains empty
	obj string

	// If this field is not empty that means unit contains entry point symbol
	entry string

	err error
}

const poolSize = 8

type Pool struct {
	output BuildOutput

	wg sync.WaitGroup

	tasks []*BuildTask

	cfg   *Config
	cache *Cache
}

func NewPool(cfg *Config, cache *Cache, tasknum int) *Pool {
	return &Pool{
		tasks: make([]*BuildTask, 0, tasknum),
		cfg:   cfg,
		cache: cache,
	}
}

func (p *Pool) AddTask(task *BuildTask) {
	order := len(p.tasks)
	task.order = order
	p.tasks = append(p.tasks, task)
}

func (p *Pool) Start() {
	tasks := make(chan *BuildTask, 16)
	results := make(chan *BuildTaskResult, 16)
	stop := make(chan struct{})

	wctl := WorkerControls{
		wg: &p.wg,

		tap:  tasks,
		sink: results,
		stop: stop,
	}

	for i := 0; i < poolSize; i++ {
		p.wg.Add(1)
		go SpawnWorker(p.cfg, p.cache, wctl)
	}

	sctl := StaplerControls{
		wg: &p.wg,

		tap:   results,
		stop:  stop,
		parts: len(p.tasks),
	}

	p.wg.Add(1)
	go SpawnStapler(p.cfg, &p.output, sctl)

	dctl := DispatcherControls{
		wg: &p.wg,

		sink: tasks,
		stop: stop,
	}

	p.wg.Add(1)
	go SpawnDispatcher(p.tasks, dctl)
}

func (p *Pool) WaitDone() {
	p.wg.Wait()
}

func (p *Pool) WaitOutput() BuildOutput {
	p.WaitDone()
	return p.output
}

type DispatcherControls struct {
	wg *sync.WaitGroup

	// outgoing build tasks
	sink chan<- *BuildTask

	// signals dispatcher that it should stop early
	stop <-chan struct{}
}

// SpawnDispatcher meant to be run in a separate goroutine. Fans out supplied tasks
// through sink channel to workers
func SpawnDispatcher(tasks []*BuildTask, ctl DispatcherControls) {
	defer ctl.wg.Done()

	for _, task := range tasks {
		select {
		case ctl.sink <- task:
		case <-ctl.stop:
			return
		}
	}
}

type BuildOutput struct {
	err error

	// entry point symbol link name
	entry string

	code PartsBuffer

	// list of paths to already built object files (produced from asm code)
	objs []string
}

type StaplerControls struct {
	wg *sync.WaitGroup

	// incoming build results from workers
	tap <-chan *BuildTaskResult

	// send signal to other goroutines in pool to stop early
	stop chan<- struct{}

	// how many parts stapler expects to receive from tap
	parts int
}

func SpawnStapler(cfg *Config, output *BuildOutput, ctl StaplerControls) {
	defer ctl.wg.Done()

	// collection of build results from received from workers
	// index in this slice directly corresponds to BuildTask.order
	results := make([]*BuildTaskResult, ctl.parts)

	// stapler tip index inside results slice
	tip := 0

	for i := 0; i < ctl.parts; i++ {
		result := <-ctl.tap
		if result.err != nil {
			// signal other goroutines in pool to stop early
			close(ctl.stop)
			output.err = fmt.Errorf("build unit \"%s\": %w", result.task.dep.UID(), result.err)
			return
		}
		if result.entry != "" {
			if debug {
				fmt.Println("entrypoint:", result.task.dep.Path.String(), result.entry)
			}
			if output.entry != "" {
				panic("two or more units with entry points")
			}
			output.entry = result.entry
		}

		results[result.task.order] = result

		// staple available parts together
		for tip < ctl.parts && results[tip] != nil {
			output.code.Add(results[tip].genout)
			if results[tip].obj != "" {
				output.objs = append(output.objs, results[tip].obj)
			}

			tip += 1
		}
	}

	if tip != ctl.parts {
		panic(fmt.Sprintf("%d parts were stapled, expected exactly %d", tip, ctl.parts))
	}

	close(ctl.stop)

	if output.entry != "" {
		var buf bytes.Buffer
		// gencpp.EntryPoint(&buf, output.entry, "coven_start", "coven::os::exit")
		output.code.Add(buf.Bytes())
	}
}

// Worker accepts unit build tasks and processes them
type Worker struct {
	cfg   *Config
	cache *Cache
}

func NewWorker(cfg *Config, cache *Cache) *Worker {
	return &Worker{
		cfg:   cfg,
		cache: cache,
	}
}

type WorkerControls struct {
	wg *sync.WaitGroup

	// incoming tasks which worker will process
	tap <-chan *BuildTask

	// outgoing build results
	sink chan<- *BuildTaskResult

	// signals worker that it should stop early
	stop <-chan struct{}
}

// SpawnWorker meant to be run in a separate goroutine. Loops through incoming
// tasks, processes them and sends results further
func SpawnWorker(cfg *Config, cache *Cache, ctl WorkerControls) {
	defer ctl.wg.Done()

	w := NewWorker(cfg, cache)
	for {
		select {
		case task, ok := <-ctl.tap:
			if !ok {
				return
			}

			result := w.Process(task)

			select {
			case ctl.sink <- result:
			case <-ctl.stop:
				return
			}
		case <-ctl.stop:
			return
		}
	}
}

func (w *Worker) Process(task *BuildTask) *BuildTaskResult {
	result, err := w.process(task)
	if err != nil {
		return &BuildTaskResult{
			err:  err,
			task: task,
		}
	}
	if result == nil {
		return &BuildTaskResult{task: task}
	}
	return result
}

func (w *Worker) process(task *BuildTask) (*BuildTaskResult, error) {
	if len(task.dep.BuildInfo.Files) == 0 {
		return nil, nil
	}

	p := task.dep.Path
	_, err := w.cache.InitUnitDir(p)
	if err != nil {
		return nil, fmt.Errorf("create unit \"%s\" cache directory: %w", p.String(), err)
	}

	names := task.dep.BuildInfo.Files
	files := make([]*source.File, 0, len(names))
	asmFiles := make([]*source.File, 0)
	var unitModTime time.Time

	for _, name := range names {
		name = filepath.Clean(name)
		if name == "" || name == "." {
			return nil, fmt.Errorf("empty file name")
		}

		dir := filepath.Dir(name)
		if !(dir == "" || dir == ".") {
			return nil, fmt.Errorf("file \"%s\" points to a separate directory", name)
		}

		src, err := w.cache.InfoSourceFile(task.dep.Path, name)
		if err != nil {
			return nil, err
		}
		if src.Kind == source.NoExt {
			return nil, fmt.Errorf("file \"%s\" has no extension", name)
		}
		if src.Kind == source.Unknown {
			return nil, fmt.Errorf("file \"%s\" has extension unknown to builder", name)
		}

		if src.Kind == source.ASM {
			asmFiles = append(asmFiles, src)
			continue
		}

		if src.ModTime.After(unitModTime) {
			unitModTime = src.ModTime
		}
		files = append(files, src)
	}

	asmObj, err := asmBuildObject(w.cache, task.dep, asmFiles)
	if err != nil {
		return nil, err
	}

	genout, ok := w.cache.LoadUnitGenout(task.dep.Path, unitModTime)
	if ok {
		var cachedUnitEntryPoint string
		if task.scanMain {
			cachedUnitEntryPoint = w.cache.LoadUnitEntryPoint(task.dep.Path)
		}
		return &BuildTaskResult{
			task:   task,
			genout: genout,
			obj:    asmObj,
			entry:  w.wrapEntryPoint(task, cachedUnitEntryPoint),
		}, nil
	}

	var buf bytes.Buffer
	var unitEntryPoint string
	var unitEntryPointFile string
	for _, file := range files {
		var err error

		switch file.Kind {
		case source.KU:
			var fileEntryPoint string
			fileEntryPoint, err = w.gizmoGenout(&buf, task, file.Name)
			if fileEntryPoint != "" {
				if unitEntryPoint != "" {
					return nil, fmt.Errorf("duplicate entry point \"%s\"", fileEntryPoint)
				}
				unitEntryPoint = fileEntryPoint
				unitEntryPointFile = file.Name
			}
		case source.CPP:
			err = w.cppGenout(&buf, task.dep, file.Name)
		case source.ASM:
			panic("unexpected source file kind: " + file.Kind.String())
		default:
			panic(fmt.Sprintf("unknown source file kind %d", file.Kind))
		}

		if err != nil {
			return nil, err
		}
	}

	if task.scanMain {
		if unitEntryPoint == "" {
			unitEntryPoint = w.cache.LoadUnitEntryPoint(task.dep.Path)
		}
	}

	result := &BuildTaskResult{
		task:   task,
		genout: buf.Bytes(),
		obj:    asmObj,
		entry:  w.wrapEntryPoint(task, unitEntryPoint),
	}
	w.cache.SaveUnitGenout(task.dep.Path, result.genout, unitEntryPoint, unitEntryPointFile)
	return result, nil
}

// add global namespace prefix and unit default namespace to entry point name
func (w *Worker) wrapEntryPoint(task *BuildTask, entry string) string {
	if entry == "" {
		return ""
	}

	globalPrefix := w.cfg.BaseNamespace
	if globalPrefix != "" && !strings.HasSuffix(globalPrefix, "::") {
		globalPrefix += "::"
	}
	return globalPrefix + task.dep.BuildInfo.DefaultNamespace + "::" + entry
}

// returns path to built object
func asmBuildObject(cache *Cache, dep *DepEntry, srcs []*source.File) (string, error) {
	if len(srcs) == 0 {
		// if there no asm files in unit, we do not need to build anything
		return "", nil
	}

	args := make([]string, 0, 2+len(srcs))
	path := filepath.Join(cache.dir, "unit", formatHash(dep.Path.Hash()), genPartObjName(srcs))
	args = append(args, "-o", path)
	for _, src := range srcs {
		args = append(args, src.Path)
	}

	cmd := exec.Command("as", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return path, nil
}

func (w *Worker) cppGenout(buf *bytes.Buffer, dep *DepEntry, name string) error {
	p := dep.Path
	src, err := w.cache.LoadSourceFile(p, name)
	if err != nil {
		return err
	}
	buf.Write(src.Bytes)
	buf.WriteString("\n")
	return nil
}

func (w *Worker) debug(format string, args ...any) {
	fmt.Printf("[debug] worker  | "+format+"\n", args...)
}

func (w *Worker) gizmoGenout(buf *bytes.Buffer, task *BuildTask, name string) (string, error) {
	p := task.dep.Path
	src, err := w.cache.LoadSourceFile(p, name)
	if err != nil {
		return "", err
	}
	genout, ok := w.cache.LoadFileGenout(p, src)
	if ok {
		// implementation always writes all bytes without error
		buf.Write(genout)
		return "", nil
	}

	if debug {
		w.debug("compile \"%s\"", src.Path)
	}

	atom, err := parser.ParseSource(src)
	if err != nil {
		return "", err
	}

	var entryPoint string
	if task.scanMain {
		entryPoint = "main"
	}

	genConfig := gencpp.Config{
		Size:             len(src.Bytes),
		DefaultNamespace: task.dep.BuildInfo.DefaultNamespace,

		GlobalNamespacePrefix:  w.cfg.BaseNamespace,
		SourceLocationComments: w.cfg.BuildKind == BuildDebug,
	}
	start := buf.Len()
	err = gencpp.Gen(buf, &genConfig, atom)
	if err != nil {
		return "", err
	}
	genout = buf.Bytes()[start:]
	w.cache.SaveFileGenout(p, src, genout)
	w.cache.SaveUnitFileToMap(p, src)
	return entryPoint, nil
}
