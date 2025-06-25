package api

import "time"

type FileObject struct {
	Name       string    `json:"name"`
	Mode       string    `json:"mode"`
	ModeBits   string    `json:"mode_bits"`
	Size       int64     `json:"size"`
	IsFile     bool      `json:"is_file"`
	IsSymlink  bool      `json:"is_symlink"`
	MimeType   string    `json:"mimes"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

// SignedURL represents a response containing a temporary, signed URL.
type SignedURL struct {
	URL string `json:"url"`
}

// RenameFile represents a single file rename operation.
type RenameFile struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// RenameFilesOptions defines the request body for renaming one or more files/folders.
type RenameFilesOptions struct {
	// The root directory for the rename operations.
	Root string `json:"root"`
	// A list of files to rename.
	Files []RenameFile `json:"files"`
}

// CopyFileOptions defines the request body for copying a file.
type CopyFileOptions struct {
	// The location to copy the file to.
	Location string `json:"location"`
}

// CompressFilesOptions defines the request body for compressing files.
type CompressFilesOptions struct {
	// The root directory where the files are located.
	Root string `json:"root"`
	// A list of file and folder paths to include in the archive.
	Files []string `json:"files"`
}

// DecompressFileOptions defines the request body for decompressing an archive.
type DecompressFileOptions struct {
	// The root directory where the archive is located.
	Root string `json:"root"`
	// The path to the archive file to decompress.
	File string `json:"file"`
}

// DeleteFilesOptions defines the request body for deleting files.
type DeleteFilesOptions struct {
	// The root directory where the files are located.
	Root string `json:"root"`
	// A list of file and folder paths to delete.
	Files []string `json:"files"`
}

// CreateFolderOptions defines the request body for creating a new folder.
type CreateFolderOptions struct {
	// The root directory where the new folder should be created.
	Root string `json:"root"`
	// The name of the new folder.
	Name string `json:"name"`
}

// fileObjectResponse and signedURLResponse are helpers for unmarshaling.
type FileObjectResponse struct {
	Object     string      `json:"object"`
	Attributes *FileObject `json:"attributes"`
}

type SignedURLResponse struct {
	Object     string     `json:"object"`
	Attributes *SignedURL `json:"attributes"`
}

// fileListResponse is a special helper as this endpoint is not paginated like others.
type FileListResponse struct {
	Object string                `json:"object"`
	Data   []*FileObjectResponse `json:"data"`
}
