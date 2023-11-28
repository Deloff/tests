package pool

import (
	"context"
	"fmt"
	"time"
)

type Handler interface {
	Handle(job *Task)
}

type TaskHandler struct {
	ch       chan Task
	result   *Result
	handlers uint8
}

func NewTaskHandler(ch chan Task, res *Result, handlers uint8) *TaskHandler {
	return &TaskHandler{ch: ch, result: res, handlers: handlers}
}

func (th *TaskHandler) Create(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(th.ch)
			return
		default:
			created := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				created = "Some error occured"
			}
			th.ch <- Task{Created: created, Id: int(time.Now().Unix())} // передаем таск на выполнение
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
			val.Handled = time.Now().Format(time.RFC3339Nano)
			th.sortResult(val)
		}
	}
}

func (th *TaskHandler) RunHandlers(ctx context.Context, job Handler) {
	for i := uint8(0); i < th.handlers; i++ {
		go th.handle(ctx, job)
	}
}

func (th *TaskHandler) sortResult(t *Task) {
	if string(t.Result[14:]) == "successed" {
		th.result.AddDone(t)
		return
	}
	th.result.AddError(fmt.Errorf("Task id %d time %s, error %s", t.Id, t.Created, t.Result))
}
