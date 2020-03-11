package msg_q

import (
	"github.com/go-redis/redis"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type TaskQ struct {
	client *redis.Client

	addTaskSha         string
	selectTaskSha      string
	pushTaskResultSha  string
	redoExpiredTaskSha string
	delTaskSha         string
}

type Data struct {
	model string
}

type Result struct {
	progress int
	image    string
}

func CreateTaskQ() *TaskQ {
	var q TaskQ
	q.client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	pong, err := q.client.Ping().Result()
	if err != nil {
		log.Panicf("ping failed: %v", err)
	}
	log.Println(pong)

	q.addTaskSha = q.loadScript("lua/add-task.lua")
	q.selectTaskSha = q.loadScript("lua/select-task.lua")
	q.pushTaskResultSha = q.loadScript("lua/push-task-result.lua")
	q.redoExpiredTaskSha = q.loadScript("lua/redo-expired-task.lua")
	q.delTaskSha = q.loadScript("lua/del-task.lua")

	return &q
}

func (q *TaskQ) Close() {
	_ = q.client.Close()
}

func (q *TaskQ) loadScript(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		log.Panic(err)
	}
	defer func() { _ = file.Close() }()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		log.Panicf("read file(%v) failed: %v", filename, err)
	}

	sha, err := q.client.ScriptLoad(string(buf)).Result()
	if err != nil {
		log.Panicf("script(%v) load failed: %v", filename, err)
	}

	log.Printf("load script %v succeed, sha is: %v", filename, sha)
	return sha
}

func (q *TaskQ) AddTask(taskId string, data string) {
	_, err := q.client.EvalSha(q.addTaskSha, []string{}, taskId, data).Result()
	if err != nil {
		log.Panicf("add task(%v, %v) failed: %v", taskId, data, err)
	}
}

func (q *TaskQ) DelTask(taskId string) {
	res, err := q.client.EvalSha(q.delTaskSha, []string{}, taskId).Result()
	if err != nil {
		log.Panic(err)
	}

	if _, ok := res.(int64); !ok {
		log.Panic(res)
	}
}

func (q *TaskQ) UpdateTaskStat(taskId string, stat string) {
	res, err := q.client.HSet("task-stat-hset", taskId, stat).Result()
	if err != nil {
		log.Panicf("update task stat failed: %v", err)
	}

	log.Printf("udpate task stat res(%v)", res)
}

func (q *TaskQ) getTaskStat(taskId string) string {
	res, err := q.client.HGet("task-stat-hset", taskId).Result()
	if err != nil {
		log.Panicf("get task stat failed: %v", err)
	}

	return res
}

func (q *TaskQ) PushTaskResult(taskId string, tracerId string, result string) bool {
	res, err := q.client.EvalSha(q.pushTaskResultSha, []string{}, taskId, tracerId, result).Result()
	if err != nil {
		log.Panicf("push task(%v) result(%v) failed: %v", taskId, result, err)
	}

	if n, ok := res.(int64); ok {
		return n == 1
	}

	return false
}

func (q *TaskQ) PullTaskResult(taskId string) string {
	res, err := q.client.HGet("task-result-hset", taskId).Result()
	if err != nil {
		log.Println(err)
	}

	return res
}

func (q *TaskQ) RedoExpiredTask(timeout time.Duration) {
	res, err := q.client.EvalSha(q.redoExpiredTaskSha, []string{}, timeout.Seconds()).Result()
	if err != nil {
		log.Println(err)
	}

	// res must be a int64 number
	if _, ok := res.(int64); !ok {
		log.Panic(res)
	}
}

func (q *TaskQ) SelectTask(tracerId string) string {
	res, err := q.client.EvalSha(q.selectTaskSha, []string{}, tracerId).Result()
	if err != nil {
		log.Panicf("select task failed: %v", err)
	}

	// convert to string
	if s, ok := res.(string); ok {
		return s
	}

	return ""
}
