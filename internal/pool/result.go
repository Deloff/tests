package pool

import (
	"fmt"
	"sync"
)

type Result struct {
	done   map[int]*Task
	mu     sync.Mutex
	errors []error
}

func NewResult(done map[int]*Task, errors []error) *Result {
	return &Result{done: done, errors: errors}
}

// AddError - add error task.
func (r *Result) AddError(err error) {
	r.mu.Lock()
	r.errors = append(
		r.errors,
		err,
	)
	r.mu.Unlock()
}

// AddDone - add done task.
func (r *Result) AddDone(job *Task) {
	r.mu.Lock()
	r.done[job.ID] = job
	r.mu.Unlock()
}

// Errors - return errors tasks
func (r *Result) Errors() []error {
	return r.errors
}

// Done - return done tasks.
func (r *Result) Done() map[int]*Task {
	return r.done
}

// PrintErrors - print errors tasks.
func (r *Result) PrintErrors() {
	fmt.Println("Errors:")

	for _, res := range r.Errors() {
		if res == nil {
			continue
		}
		fmt.Println(res.Error())
	}
}

// PrintDone - print done tasks.
func (r *Result) PrintDone() {
	fmt.Println("Done tasks:")

	for res := range r.Done() {
		fmt.Println(res)
	}
}
