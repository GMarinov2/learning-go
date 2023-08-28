package main

import (
	"fmt"
	"os"
	"path"
	tikaclient "song-lyrics-indexer/tika-client"
	"sync"
)

func main() {
	languageDetector := tikaclient.NewClient("http://localhost:9998")

	filePaths, err := listFilesInDirectory("./data")
	if err != nil {
		fmt.Println(err)
	}
	fileContents := readAllFiles(filePaths)
	wg := sync.WaitGroup{}

	lyricsByLanguage := make(map[string]chan string)

	for {
		file, ok := <-fileContents
		if !ok {
			for name, chanel := range lyricsByLanguage {
				fmt.Printf("Closing channel %v \n", name)
				close(chanel)
			}
			break
		}

		response, err := languageDetector.DetectLanguage(file)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if lyricsByLanguage[response] != nil {
			lyricsByLanguage[response] <- file
		} else {
			lyricsByLanguage[response] = make(chan string)
			_ = os.Mkdir(fmt.Sprintf("./output/%v", response), 0777)
			go func() {
				wg.Add(1)
				writeFiles(response, lyricsByLanguage[response])
				wg.Done()
			}()
			lyricsByLanguage[response] <- file
		}
	}

	wg.Wait()
}

func listFilesInDirectory(basePath string) (<-chan string, error) {
	result := make(chan string)

	files, error := os.ReadDir(basePath)

	if error != nil {
		return result, error
	}

	go func() {
		defer close(result)

		for _, file := range files {
			result <- path.Join(basePath, file.Name())
		}
	}()

	return result, nil
}

func readFile(filePath string) string {
	content, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Println(err)
	}

	return string(content)
}

func readAllFiles(filePaths <-chan string) <-chan string {
	result := make(chan string)
	go func() {
		for {
			filePath, ok := <-filePaths
			if !ok {
				close(result)
				return
			}

			content := readFile(filePath)

			result <- content
		}
	}()

	return result
}

func writeFiles(language string, fileContents <-chan string) {
	index := 0
	for {
		text, ok := <-fileContents
		if !ok {
			return
		}

		fileName := fmt.Sprintf("./output/%v/%v.txt", language, index)
		fmt.Println(fileName)
		newFile, err := os.Create(fileName)
		if err != nil {
			fmt.Println(err)
		}

		defer newFile.Close()

		err = os.WriteFile(fileName, []byte(text), 0664)
		if err != nil {
			fmt.Println(err)
		}
		index++
	}
}
