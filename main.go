package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Ttype represents a meaninglessness of our life

type (
	Ttype struct {
		id         int
		cT         string // время создания
		fT         string // время выполнения
		taskRESULT []byte
	}
	TaskHandler struct {
		ch     chan Ttype
		result *Result
	}
	Handler interface {
		Handle(job *Ttype)
	}

	JobHandler struct{}

	Result struct {
		done   map[int]*Ttype
		mu     sync.RWMutex
		errors []error
	}
)

func NewTaskHandler(ch chan Ttype, res *Result) *TaskHandler {
	return &TaskHandler{ch: ch, result: res}
}

func (th *TaskHandler) create(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(th.ch)
			return
		default:
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				ft = "Some error occured"
			}
			th.ch <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
		}

	}
}

func (th *TaskHandler) handle(ctx context.Context, job Handler) {
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-th.ch:
			val := &data
			job.Handle(val)
			th.sortResult(val)
		}
	}
}

func (th *TaskHandler) sortResult(t *Ttype) {
	if string(t.taskRESULT[14:]) == "successed" {
		th.result.AddDone(t)
		return
	}
	th.result.AddError(fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT))
}

func (jh JobHandler) Handle(job *Ttype) {
	tt, _ := time.Parse(time.RFC3339, job.cT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		job.taskRESULT = []byte("task has been successed")
	} else {
		job.taskRESULT = []byte("something went wrong")
	}
	job.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)
}

func (r *Result) AddError(err error) {
	r.mu.Lock()
	r.errors = append(
		r.errors,
		err,
	)
	r.mu.Unlock()

}

func (r *Result) AddDone(job *Ttype) {
	r.mu.Lock()
	r.done[job.id] = job
	r.mu.Unlock()
}

func handleShutdown() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	go func() {
		chSystem := make(chan os.Signal)
		signal.Notify(chSystem, os.Interrupt, syscall.SIGTERM)
		<-chSystem
		cancel()

	}()

	return ctx
}

func main() {
	ctx := handleShutdown()

	handlers := 5

	result := &Result{errors: make([]error, 50), done: make(map[int]*Ttype)}
	tc := NewTaskHandler(make(chan Ttype, handlers), result)

	for i := 0; i < handlers; i++ {
		go tc.handle(ctx, &JobHandler{})
	}

	go tc.create(ctx)

	<-ctx.Done()

	println("Errors:")
	for _, r := range result.errors {
		if r == nil {
			continue
		}
		println(r.Error())
	}

	println("Done tasks:")
	for r := range result.done {
		println(r)
	}

}
