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
	Executable bool `json:"executable,omitempty"`
	LastModified int64 `json:"lastModified,omitempty"`
	Directory bool `json:"directory,omitempty"`
	Writable bool `json:"writable,omitempty"`
	TotalSpace int64 `json:"totalSpace,omitempty"`
	Readable bool `json:"readable,omitempty"`
	CanonicalFile **os.File `json:"canonicalFile,omitempty"`
	FreeSpace int64 `json:"freeSpace,omitempty"`
	File bool `json:"file,omitempty"`
	Path string `json:"path,omitempty"`
	UsableSpace int64 `json:"usableSpace,omitempty"`
	AbsolutePath string `json:"absolutePath,omitempty"`
	Parent string `json:"parent,omitempty"`
	Hidden bool `json:"hidden,omitempty"`
	ParentFile **os.File `json:"parentFile,omitempty"`
	Absolute bool `json:"absolute,omitempty"`
	AbsoluteFile **os.File `json:"absoluteFile,omitempty"`
	Name string `json:"name,omitempty"`
	CanonicalPath string `json:"canonicalPath,omitempty"`
}
