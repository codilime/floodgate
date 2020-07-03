/*
 * Spinnaker API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"os"
)

type File struct {
	Absolute      bool      `json:"absolute,omitempty"`
	AbsoluteFile  **os.File `json:"absoluteFile,omitempty"`
	AbsolutePath  string    `json:"absolutePath,omitempty"`
	CanonicalFile **os.File `json:"canonicalFile,omitempty"`
	CanonicalPath string    `json:"canonicalPath,omitempty"`
	Directory     bool      `json:"directory,omitempty"`
	Executable    bool      `json:"executable,omitempty"`
	File          bool      `json:"file,omitempty"`
	FreeSpace     int64     `json:"freeSpace,omitempty"`
	Hidden        bool      `json:"hidden,omitempty"`
	LastModified  int64     `json:"lastModified,omitempty"`
	Name          string    `json:"name,omitempty"`
	Parent        string    `json:"parent,omitempty"`
	ParentFile    **os.File `json:"parentFile,omitempty"`
	Path          string    `json:"path,omitempty"`
	Readable      bool      `json:"readable,omitempty"`
	TotalSpace    int64     `json:"totalSpace,omitempty"`
	UsableSpace   int64     `json:"usableSpace,omitempty"`
	Writable      bool      `json:"writable,omitempty"`
}
