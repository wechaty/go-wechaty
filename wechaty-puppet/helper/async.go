package helper

import (
	"math"
	"runtime"
	"sync"
)

// DefaultWorkerNum default number of goroutines is twice the number of GOMAXPROCS
var DefaultWorkerNum = runtime.GOMAXPROCS(0) * 2

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
		tasks        []Task
		maxWorkerNum int
	}

	// Task task func
	Task func() (interface{}, error)
)

func NewAsync(maxWorkerNum int) IAsync {
	if maxWorkerNum <= 0 {
		maxWorkerNum = DefaultWorkerNum
	}
	return &async{
		maxWorkerNum: maxWorkerNum,
	}
}

func (a *async) AddTask(task Task) {
	a.tasks = append(a.tasks, task)

}

func (a *async) Result() []AsyncResult {
	taskChan := make(chan Task)
	resultChan := make(chan AsyncResult)

	taskNum := len(a.tasks)
	wg := sync.WaitGroup{}
	wg.Add(taskNum)

	workerNum := int(math.Min(float64(taskNum), float64(a.maxWorkerNum)))
	for i := 0; i < workerNum; i++ {
		go func() {
			for task := range taskChan {
				result := AsyncResult{}
				result.Value, result.Err = task()
				resultChan <- result
				wg.Done()
			}
		}()
	}

	go func() {
		for _, v := range a.tasks {
			taskChan <- v
		}
		wg.Wait()
		close(resultChan)
		close(taskChan)
		a.tasks = make([]Task, 0)
	}()

	result := make([]AsyncResult, taskNum)
	idx := 0
	for v := range resultChan {
		result[idx] = v
		idx += 1
	}
	return result
}
