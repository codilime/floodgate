package resourcemanager

import fl "github.com/codilime/floodgate/parser/fileloader"

// Options store ResourceManager options
type Options struct {
	fileLoaders []fl.FileLoader
}

// Option ResourceManager option
type Option func(op *Options)

// FileLoaders set file loaders
func FileLoaders(fileLoaders ...fl.FileLoader) Option {
	return func(op *Options) {
		op.fileLoaders = fileLoaders
	}
}
