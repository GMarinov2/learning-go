package main

import (
	"song-lyrics-indexer/args"
	"song-lyrics-indexer/fileworker"
	"song-lyrics-indexer/languageworker"

	"github.com/alexflint/go-arg"
)

func main() {
	args := args.Args{}

	arg.MustParse(&args)
	languageWorker := languageworker.NewLanguageWorker(args.Tika)

	fileWorker := fileworker.FileWorker{
		SourceFolder:      args.Source,
		DestinationFolder: args.Destination,
	}

	fileNames, err := fileWorker.ListAllFiles()
	if err != nil {
		return
	}

	fileEntries, err := fileWorker.ReadAllFiles(fileNames)
	if err != nil {
		return
	}

	detectedLanguages := languageWorker.GetFileLanguages(fileEntries)

	fileWorker.WriteAllFiles(detectedLanguages)
}
