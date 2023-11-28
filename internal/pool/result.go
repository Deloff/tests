package pool

import "sync"

type Result struct {
	done   map[int]*Task
	mu     sync.Mutex
	errors []error
}

func NewResult(done map[int]*Task, errors []error) *Result {
	return &Result{done: done, errors: errors}
}

func (r *Result) AddError(err error) {
	r.mu.Lock()
	r.errors = append(
		r.errors,
		err,
	)
	r.mu.Unlock()

}

func (r *Result) AddDone(job *Task) {
	r.mu.Lock()
	r.done[job.Id] = job
	r.mu.Unlock()
}

func (r *Result) Errors() []error {
	return r.errors
}

func (r *Result) Done() map[int]*Task {
	return r.done
}
