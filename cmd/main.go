package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Deloff/test/internal/handlers"
	"github.com/Deloff/test/internal/pool"
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
		chSystem := make(chan os.Signal, 1)
		signal.Notify(chSystem, os.Interrupt, syscall.SIGTERM)
		<-chSystem
		cancel()
	}()

	return ctx
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Info("panic!", slog.Any("recover", r))
		}
	}()

	ctx := handleShutdown()

	poolSize := uint8(10)
	result := pool.NewResult(make(map[int]*pool.Task), make([]error, 50))
	tc := pool.NewTaskHandler(make(chan pool.Task, poolSize), result, poolSize)

	go tc.Create(ctx)

	tc.RunHandlers(ctx, &handlers.JobHandler{})

	<-ctx.Done()

	result.PrintErrors()
	result.PrintDone()
}
