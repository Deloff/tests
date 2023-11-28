package handlers

import (
	"github.com/Deloff/test/internal/pool"
	"time"
)

type JobHandler struct{}

func (jh JobHandler) Handle(job *pool.Task) {
	tt, _ := time.Parse(time.RFC3339, job.Created)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		job.Result = []byte("task has been successed")
	} else {
		job.Result = []byte("something went wrong")
	}
	job.Handled = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)
}
