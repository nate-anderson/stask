package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"stask"
	"time"
)

const taskJSONFile = ".stasks.json"

var (
	filePath string
)

func printHeader(tasks int, out io.Writer) {
	fmt.Fprintf(out, "\n\u001b[33mstask\u001b[0m: %d task(s)\n", tasks)
	fmt.Fprintln(out, "-------------------------------------------")
}

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	taskJSONPath := filepath.Join(user.HomeDir, taskJSONFile)
	flag.StringVar(&filePath, "f", taskJSONPath, "the location of the stask JSON file")
	removeFlag := flag.Int("r", -1, "remove the numbered task")
	addTaskFlag := flag.String("t", "", "add a task to the queue")
	popFlag := flag.Bool("p", false, "pop the top task off the queue")
	flag.Parse()

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	queue, err := stask.ReadQueue(file)
	if err != nil {
		log.Fatal(err)
	}

	if popFlag != nil && *popFlag {
		_, err := queue.Pop()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(os.Stdout, "popped task 0!\n")
	}

	if *removeFlag >= 0 {
		err := queue.Remove(*removeFlag)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(os.Stdout, "removed task %d!\n", *removeFlag)
	}

	if *addTaskFlag != "" {
		t := stask.Task{
			Description: *addTaskFlag,
			Created:     time.Now(),
		}
		if id, err := queue.Push(t); err != nil {
			log.Fatal(err)
		} else {
			fmt.Fprintf(os.Stdout, "created task %d\n", id)
		}
	}

	printHeader(queue.Count(), os.Stdout)
	queue.Print(os.Stdout)

	if err := queue.Write(filePath); err != nil {
		log.Fatal(err)
	}
}
