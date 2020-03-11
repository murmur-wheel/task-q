package msg_q

import (
	uuid "github.com/satori/go.uuid"
	"testing"
	"time"
)

func TestTaskQ_AddTask(t *testing.T) {
	q := CreateTaskQ()
	defer q.Close()

	q.client.FlushDB() // clear database

	taskId := uuid.NewV4().String()
	q.AddTask(taskId, "task-data")
	q.UpdateTaskStat(taskId, "CANCELED")
}

func TestTaskQ_UpdateTaskStat(t *testing.T) {
	q := CreateTaskQ()
	defer q.Close()

	q.client.FlushDB() // clear database

	taskId := uuid.NewV4().String()
	q.AddTask(taskId, "test task")

	f := func(stat string) {
		q.UpdateTaskStat(taskId, stat)
		temp := q.getTaskStat(taskId)
		if temp != stat {
			t.Errorf("task stat should be %v, be we got %v", stat, temp)
		}
	}

	f("PENDING")
	f("DOING")
	f("CANCELED")
	f("FINISHED")
}

func TestTaskQ_SelectTask(t *testing.T) {
	q := CreateTaskQ()
	defer q.Close()

	q.client.FlushDB() // clear database

	taskId := uuid.NewV4().String()
	tracerId := uuid.NewV4().String()

	// prepare tasks
	q.AddTask(taskId, "task-data")

	// select tasks
	selected1 := q.SelectTask(tracerId)
	selected2 := q.SelectTask(tracerId)

	if selected1 != taskId && selected2 != "" {
		t.Errorf("selected1(%v) must be %v, selected2(%v) should be \"\"", selected1, taskId, selected2)
	}
}

func TestTaskQ_PushTaskResult(t *testing.T) {
	q := CreateTaskQ()
	defer q.Close()

	q.client.FlushDB()

	task := uuid.NewV4().String()
	q.AddTask(task, "task-data")

	tracer1 := uuid.NewV4().String()
	selected := q.SelectTask(tracer1)

	if !q.PushTaskResult(selected, tracer1, "trace-result1") {
		t.Errorf("tracer1 push task result should succeed")
	}

	tracer2 := uuid.NewV4().String()
	if q.PushTaskResult(selected, tracer2, "trace-result2") {
		t.Errorf("tracer2 push task result should failed")
	}
}

func TestTaskQ_PullTaskResult(t *testing.T) {
	q := CreateTaskQ()
	defer q.Close()

	q.client.FlushDB() // reset database

	task1 := uuid.NewV4().String()
	tracer1 := uuid.NewV4().String()
	result1 := "result"

	q.AddTask(task1, "task1-data")
	task2 := q.SelectTask(tracer1)
	q.PushTaskResult(task2, tracer1, result1)
	result2 := q.PullTaskResult(task2)
	if result1 != result2 || task1 != task2 {
		t.Error(task1, task2, result1, result2)
	}
}

func TestTaskQ_RedoExpiredTask(t *testing.T) {
	q := CreateTaskQ()
	defer q.Close()

	q.client.FlushDB()

	task1 := uuid.NewV4().String()
	tracer1 := uuid.NewV4().String()

	q.AddTask(task1, "task1")
	task2 := q.SelectTask(tracer1)
	time.Sleep(time.Second * 2)

	q.RedoExpiredTask(time.Second)
	stat1 := q.getTaskStat(task1)
	stat2 := q.getTaskStat(task2)
	if stat1 != stat2 || task1 != task2 {
		t.Errorf("stat1: %v, stat2: %v\ntask1: %v, task2: %v",
			stat1, stat2, task1, task2)
	}
}
