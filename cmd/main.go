package main

import (
	"context"
	handlers "github.com/Deloff/test/internal/handlers"
	"github.com/Deloff/test/internal/pool"
	"os"
	"os/signal"
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

	poolSize := uint8(10)
	result := pool.NewResult(make(map[int]*pool.Task), make([]error, 50))
	tc := pool.NewTaskHandler(make(chan pool.Task, poolSize), result, poolSize)

	go tc.Create(ctx)

	tc.RunHandlers(ctx, &handlers.JobHandler{})

	<-ctx.Done()

	println("Errors:")
	for _, r := range result.Errors() {
		if r == nil {
			continue
		}
		println(r.Error())
	}

	println("Done tasks:")
	for r := range result.Done() {
		println(r)
	}

}
