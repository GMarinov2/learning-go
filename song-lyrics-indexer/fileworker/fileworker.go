package fileworker

import (
	"fmt"
	"os"
	"path"
	"sync"
)

type FileWorker struct {
	SourceFolder      string
	DestinationFolder string
}

type FileWorkerEntry struct {
	FileName     string
	FileContent  string
	FileLanguage string
}

func (fileWorker *FileWorker) ListAllFiles() (<-chan string, error) {
	result := make(chan string)

	files, error := os.ReadDir(fileWorker.SourceFolder)

	if error != nil {
		return result, error
	}

	go func() {
		defer close(result)

		for _, file := range files {
			result <- file.Name()
		}
	}()

	return result, nil
}

func (fileWorker *FileWorker) ReadAllFiles(fileNames <-chan string) (<-chan FileWorkerEntry, error) {
	result := make(chan FileWorkerEntry)
	go func() {
		for {
			fileName, ok := <-fileNames
			if !ok {
				close(result)
				return
			}

			content, err := fileWorker.readFile(fileName)

			if err != nil {
				fmt.Printf("Error reading file %v", fileName)
			}

			entry := FileWorkerEntry{
				FileName:    fileName,
				FileContent: content,
			}

			result <- entry
		}
	}()

	return result, nil
}

func (fileWorker *FileWorker) WriteAllFiles(fileWorkerEntries <-chan FileWorkerEntry) {
	wg := sync.WaitGroup{}
	filesByLanguage := make(map[string]chan FileWorkerEntry)

	_ = os.Mkdir(fileWorker.DestinationFolder, 0777)

	for {
		file, ok := <-fileWorkerEntries
		if !ok {
			for name, chanel := range filesByLanguage {
				fmt.Printf("Closing channel %v \n", name)
				close(chanel)
			}
			break
		}
		fmt.Println(file.FileLanguage)
		if filesByLanguage[file.FileLanguage] != nil {
			filesByLanguage[file.FileLanguage] <- file
		} else {
			wg.Add(1)
			_ = os.Mkdir(path.Join(fileWorker.DestinationFolder, file.FileLanguage), 0777)
			filesByLanguage[file.FileLanguage] = make(chan FileWorkerEntry)
			go func() {
				fileWorker.writeFiles(filesByLanguage[file.FileLanguage], &wg)
			}()
			filesByLanguage[file.FileLanguage] <- file
		}
	}

	wg.Wait()
}

func (fileWorker *FileWorker) writeFiles(files <-chan FileWorkerEntry, wg *sync.WaitGroup) {
	for {
		file, ok := <-files
		if !ok {
			wg.Done()
			return
		}

		filePath := path.Join(fileWorker.DestinationFolder, file.FileLanguage, file.FileName)
		newFile, err := os.Create(filePath)
		if err != nil {
			fmt.Println(err)
		}

		defer newFile.Close()

		err = os.WriteFile(filePath, []byte(file.FileContent), 0664)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (fileWorker *FileWorker) readFile(fileName string) (string, error) {
	filePath := path.Join(fileWorker.SourceFolder, fileName)
	content, err := os.ReadFile(filePath)

	if err != nil {
		return "", err
	}

	return string(content), nil
}
