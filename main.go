package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Ttype represents a task
type Ttype struct {
	id         int
	cT         string // Creation time
	fT         string // Finish time
	taskRESULT []byte
}

func main() {
	taskCreator := func(taskCh chan Ttype) {
		for {
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // Condition for erroneous tasks
				ft = "Some error occurred"
			}
			taskCh <- Ttype{cT: ft, id: int(time.Now().Unix())} // Send the task for processing
		}
	}

	taskCh := make(chan Ttype, 10)

	go taskCreator(taskCh)

	taskWorker := func(task Ttype, wg *sync.WaitGroup) {
		defer wg.Done()

		tt, _ := time.Parse(time.RFC3339, task.cT)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			task.taskRESULT = []byte("task has been succeeded")
		} else {
			task.taskRESULT = []byte("something went wrong")
		}
		task.fT = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		if string(task.taskRESULT[14:]) == "succeeded" {
			fmt.Printf("Task id %d succeeded\n", task.id)
		} else {
			fmt.Printf("Task id %d failed: %s\n", task.id, task.taskRESULT)
		}
	}

	var wg sync.WaitGroup
	done := make(chan struct{})

	go func() {
		for task := range taskCh {
			wg.Add(1)
			go taskWorker(task, &wg)
		}
	}()

	// Stop processing tasks after 3 seconds and close the channel
	go func() {
		time.Sleep(time.Second * 3)
		close(taskCh)
	}()

	wg.Wait()
	close(done)
}
