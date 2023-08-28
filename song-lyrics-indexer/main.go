package main

import (
	"fmt"
	"song-lyrics-indexer/args"
	"song-lyrics-indexer/fileworker"
	tikaclient "song-lyrics-indexer/tika-client"

	"github.com/alexflint/go-arg"
)

func main() {
	args := args.Args{}

	arg.MustParse(&args)

	languageDetector := tikaclient.NewClient(args.Tika)

	worker := fileworker.FileWorker{
		SourceFolder:      args.Source,
		DestinationFolder: args.Destination,
	}

	fileNames, err := worker.ListAllFiles()
	if err != nil {
		return
	}

	fileEntries, err := worker.ReadAllFiles(fileNames)
	if err != nil {
		return
	}

	detectedLanguages := getFileLanguages(fileEntries, languageDetector)

	worker.WriteAllFiles(detectedLanguages)

}

func getFileLanguages(entries <-chan fileworker.FileWorkerEntry, languageDetector tikaclient.Client) <-chan fileworker.FileWorkerEntry {
	result := make(chan fileworker.FileWorkerEntry)

	go func() {
		for {
			entry, ok := <-entries
			if !ok {
				close(result)
			}

			language, err := languageDetector.DetectLanguage(entry.FileContent)

			if err != nil {
				fmt.Println(err)
				continue
			}

			resultEntry := fileworker.FileWorkerEntry{
				FileLanguage: language,
				FileName:     entry.FileName,
				FileContent:  entry.FileContent,
			}

			result <- resultEntry
		}
	}()

	return result
}
