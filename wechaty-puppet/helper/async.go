package helper

import (
	"runtime"
	"sync"
)

// DefaultWorkerNum default number of goroutines is twice the number of CPUs
var DefaultWorkerNum = runtime.NumCPU() * 2

type (
	// IAsync interface
	IAsync interface {
		AddTask(task Task)
		Result() []AsyncResult
	}

	// AsyncResult result struct
	AsyncResult struct {
		Value interface{}
		Err   error
	}

	async struct {
		wg         sync.WaitGroup
		tasks      []Task
		workerNum  int
	}

	// Task task func
	Task func() (interface{}, error)
)

func NewAsync(workerNum int) IAsync {
	if workerNum <= 0 {
		workerNum = DefaultWorkerNum
	}
	return &async{
		wg:         sync.WaitGroup{},
		workerNum:  workerNum,
	}
}

func (a *async) AddTask(task Task) {
	a.wg.Add(1)
	a.tasks = append(a.tasks, task)
}

func (a *async) Result() []AsyncResult {
	taskChan := make(chan Task)
	resultChan := make(chan AsyncResult)
	for i := 0; i < a.workerNum; i++ {
		go func() {
			for task := range taskChan {
				result := AsyncResult{}
				result.Value, result.Err = task()
				resultChan <- result
				a.wg.Done()
			}
		}()
	}

	go func() {
		for _,v := range a.tasks {
			taskChan <- v
		}
		a.wg.Wait()
		close(resultChan)
		close(taskChan)
	}()

	result := make([]AsyncResult, 0, len(a.tasks))
	for v := range resultChan {
		result = append(result, v)
	}
	return result
}
