package languageworker

import (
	"fmt"
	"song-lyrics-indexer/fileworker"
	tikaclient "song-lyrics-indexer/tika-client"
)

type LanguageWorker struct {
	Client tikaclient.Client
}

func NewLanguageWorker(tikaClientUrl string) LanguageWorker {
	client := tikaclient.NewClient(tikaClientUrl)
	return LanguageWorker{Client: client}
}

func (languageWorker *LanguageWorker) GetFileLanguages(entries <-chan fileworker.FileWorkerEntry) <-chan fileworker.FileWorkerEntry {
	result := make(chan fileworker.FileWorkerEntry)

	go func() {
		for {
			entry, ok := <-entries
			if !ok {
				close(result)
			}

			language, err := languageWorker.Client.DetectLanguage(entry.FileContent)

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
