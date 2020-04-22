package fileloader

// FileLoader is the interface that wraps LoadFile and SupportedFileExtensions methods
//
// LoadFile loads file from given path.
// It returns file content as map[string]interface{} and any error encountered that caused
// the file loading to stop early.
//
// SupportedFileExtensions returns a slice of strings representing supported file format extensions
type FileLoader interface {
	LoadFile(filePath string) ([]map[string]interface{}, error)
	SupportedFileExtensions() []string
}
