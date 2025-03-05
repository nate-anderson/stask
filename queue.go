package stask

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Task struct {
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
}

func (t *Task) Sanitize() {
	t.Description = strings.ReplaceAll(t.Description, "\n", "")
}

type Queue struct {
	Tasks []Task `json:"tasks"`
}

func (q *Queue) Push(task Task) (int, error) {
	if task.Description == "" {
		return 0, errors.New("empty task")
	}
	task.Sanitize()
	q.Tasks = append(q.Tasks, task)
	return len(q.Tasks), nil
}

func (q *Queue) Pop() (*Task, error) {
	if len(q.Tasks) == 0 {
		return nil, errors.New("no tasks in queue, congratulations")
	}

	first := q.Tasks[0]
	q.Tasks = q.Tasks[1:]
	return &first, nil
}

func (q *Queue) Remove(i int) error {
	if len(q.Tasks) == 0 {
		return errors.New("no tasks in queue")
	}

	if len(q.Tasks) < i {
		return fmt.Errorf("there are only %d tasks", q.Count())
	}

	q.Tasks = append(q.Tasks[:i], q.Tasks[i+1:]...)
	return nil
}

func (q Queue) Count() int {
	return len(q.Tasks)
}

func ReadQueue(reader io.Reader) (*Queue, error) {
	decoder := json.NewDecoder(reader)
	q := Queue{}
	err := decoder.Decode(&q)
	return &q, err
}

func (q Queue) Print(out io.Writer) {
	for i, t := range q.Tasks {
		fmt.Fprintf(out, " \u001b[32m%d\u001b[0m: %s (created \u001b[36m%s\u001b[0m)\n", i, t.Description, t.Created.Format(time.RFC3339))
	}
	fmt.Println()
}

func (q Queue) Write(path string) error {
	data, err := json.Marshal(q)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0755)
}
