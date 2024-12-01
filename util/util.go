package util

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"
)

type Task struct {
	ID          int
	Description string
	CreatedAt   string
	IsComplete  bool
}

type File struct {
	OpenFlag   int
	Permission fs.FileMode
}

var openModes = map[string]File{
	"readOnly":    {OpenFlag: os.O_RDONLY, Permission: 0},
	"writeOnly":   {OpenFlag: os.O_WRONLY, Permission: 0644},
	"readWrite":   {OpenFlag: os.O_RDWR, Permission: 0644},
	"writeAppend": {OpenFlag: os.O_WRONLY | os.O_APPEND, Permission: 0644},
	"create":      {OpenFlag: os.O_CREATE | os.O_WRONLY, Permission: 0644},
}

func loadFile(filepath string, openMode int, permission fs.FileMode) (*os.File, error) {
	file, err := os.OpenFile(filepath, openMode, permission)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading")
	}

	// Exclusive lock obtained on the file descriptor
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		_ = file.Close()
		return nil, err
	}

	return file, nil
}

func closeFile(f *os.File) error {
	syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	return f.Close()
}

func CreateInitialDatabase(taskFilename string, taskHeaders []string) {
	if _, err := os.Stat(taskFilename); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("")
		} else {
			fmt.Printf("Error checking file: %v\n", err)
			return
		}
	} else {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Printf("Do you want to reset the database (y/n): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "y" {
				fmt.Println("Resetting...")
				time.Sleep(200 * time.Millisecond)
				break
			} else if input == "n" {
				return
			} else {
				fmt.Println("Error: Wrong Input. Please enter 'y' or 'n'.")
				return
			}
		}
	}

	fileConfig := openModes["overWrite"]
	file, err := loadFile(taskFilename, fileConfig.OpenFlag, fileConfig.Permission)
	if err != nil {
		fmt.Println("Error loading file: ", taskFilename)
		return
	}
	defer closeFile(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write(taskHeaders)
}

func getLastTaskID(taskFilename string) (int, error) {
	fileConfig := openModes["readOnly"]
	file, err := loadFile(taskFilename, fileConfig.OpenFlag, fileConfig.Permission)
	if err != nil {
		return 0, fmt.Errorf("opening file [%s]: %s", taskFilename, err)
	}
	defer closeFile(file)

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("reading CSV: %v", err)
	}

	// If no rows or only header row exists
	if len(rows) <= 1 {
		return 0, nil // No IDs yet
	}

	lastRow := rows[len(rows)-1]
	lastID, err := strconv.Atoi(lastRow[0])
	if err != nil {
		return 0, fmt.Errorf("converting ID to integer: %v\n", err)
	}

	return lastID, nil
}

func AddTask(taskFilename, taskDescription string) error {
	lastTaskId, err := getLastTaskID(taskFilename)
	if err != nil {
		return fmt.Errorf("Error Getting TaskID: %s\n", err)
	}

	fileConfig := openModes["writeAppend"]
	file, err := loadFile(taskFilename, fileConfig.OpenFlag, fileConfig.Permission)
	if err != nil {
		return fmt.Errorf("Error loading file [%s]: [%s]\n", taskFilename, err)
	}
	defer closeFile(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	currentTime := time.Now().Format("2006-01-02T15:04:05Z")

	task := []string{
		fmt.Sprintf("%d", lastTaskId+1), // ID
		taskDescription,                 // Description
		currentTime,                     // CreatedAt
		"false",                         // IsComplete
	}
	writer.Write(task)

	fmt.Printf("Task Added [ID: %v]\n", task[0])
	return nil
}

func ListTasks(taskFilename string, allFlagPassed bool) error {
	fileConfig := openModes["readOnly"]
	file, err := loadFile(taskFilename, fileConfig.OpenFlag, fileConfig.Permission)
	if err != nil {
		return fmt.Errorf("opening file [%s]: %s\n", taskFilename, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("reading file [%s]: %s\n", taskFilename, err)
	}

	if len(rows) > 1 {
		rows = rows[1:]
		tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

		if allFlagPassed {
			fmt.Fprintln(tabWriter, "ID\tTask\tCreated\tDone")
			for _, row := range rows {
				fmt.Fprintf(tabWriter, "%s\t%s\t%s\t%s\n", row[0], row[1], row[2], row[3])
			}

		} else {
			fmt.Fprintln(tabWriter, "ID\tTask\tCreated")
			for _, row := range rows {
				if row[len(row)-1] != "true" {
					fmt.Fprintf(tabWriter, "%s\t%s\t%s\n", row[0], row[1], row[2])
				}
			}
		}
		tabWriter.Flush()
	}
	return nil
}

func MarksTaskAsComplete(taskFilename, taskId string) error {
	taskIdInt, err := strconv.ParseInt(taskId, 10, 64)
	if err != nil {
		return fmt.Errorf("parsing taskId [%s]: %s\n", taskId, err)
	}

	fileConfig := openModes["readWrite"]
	file, err := loadFile(taskFilename, fileConfig.OpenFlag, fileConfig.Permission)
	if err != nil {
		return fmt.Errorf("loading file [%s]: %s\n", taskFilename, err)
	}
	defer closeFile(file)

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("reading file [%s]: %s\n", taskFilename, err)
	}

	if len(rows) != 0 {
		for i, row := range rows {
			if row[0] != "ID" && row[0] == taskId {
				// Modifiy Slice Element
				if row[3] == "true" {
					fmt.Println("Task already Completed")
					return nil
				}
				row[3] = "true"
				rows[i] = row
				break
			}
		}
	}

	// Empty File
	file.Truncate(0)
	file.Seek(0, 0)

	writer := csv.NewWriter(file)
	err = writer.WriteAll(rows)
	if err != nil {
		return fmt.Errorf("writing to file [%s]: %s\n", taskFilename, err)
	}
	writer.Flush()
	fmt.Println("Task Completed: ", rows[taskIdInt][1])
	return nil
}

func DeleteTaskFromCSV(taskFilename, taskId string) error {
	taskIdInt, err := strconv.ParseInt(taskId, 10, 64)
	if err != nil {
		return fmt.Errorf("parsing taskId [%s]: %s\n", taskId, err)
	}

	// Load the file
	fileConfig := openModes["readWrite"]
	file, err := loadFile(taskFilename, fileConfig.OpenFlag, fileConfig.Permission)
	if err != nil {
		return fmt.Errorf("loading file [%s]: %s\n", taskFilename, err)
	}
	defer closeFile(file)

	// Read the content of the file
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("reading file [%s]: %s\n", taskFilename, err)
	}

	// Flag to check if taskId is found
	taskFound := false

	// Check for the task and delete it
	if len(rows) != 0 {
		for i, row := range rows {
			if row[0] != "ID" && row[0] == taskId {
				// Remove task row from slice
				rows = append(rows[:i], rows[i+1:]...)
				taskFound = true
				break
			}
		}
	}

	// taskId was not found
	if !taskFound {
		fmt.Println("Task already Deleted")
		return nil
	}

	// Empty the file
	file.Truncate(0)
	file.Seek(0, 0)

	// Write the updated rows back to the file
	writer := csv.NewWriter(file)
	err = writer.WriteAll(rows)
	if err != nil {
		return fmt.Errorf("writing to file [%s]: %s\n", taskFilename, err)
	}
	writer.Flush()

	fmt.Println("Task Deleted: ", rows[taskIdInt][1])
	return nil
}
