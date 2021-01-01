package deploy

// Diff describes how to get from one folder state to another
type Diff struct {
	Add    []File
	Rename []Rename
	Delete []Delete
}

// Rename represents a rename action
type Rename struct {
	From string
	To   string
}

// Delete represents a delete action
type Delete struct {
	Key string
}

// GenerateDiff generates a diffmanifest between two manifests
func GenerateDiff(oldManifest, newManifest Manifest) Diff {
	newFiles := []File{}
	renamedFiles := []Rename{}
	deletedFiles := []Delete{}

	// Slice containing renamed files ORIGINAL name
	renamedFilesSlice := []string{}

	// Loop through every file in the new manifest
	for _, file := range newManifest.Files {
		fileFound := false

		// Check if hash exists in old manifest
		for _, oldFile := range oldManifest.Files {
			// File exists in old manifest but file name is different
			if file.Checksum == oldFile.Checksum {
				fileFound = true

				if file.Key != oldFile.Key {
					// For it to be considered a rename, ensure the original key is absent from the new manifest
					keyAbsent := true

					for _, newFile := range newManifest.Files {
						if oldFile.Key == newFile.Key {
							keyAbsent = false
						}
					}

					if keyAbsent {
						renamedFiles = append(renamedFiles, Rename{
							From: oldFile.Key,
							To:   file.Key,
						})

						renamedFilesSlice = append(renamedFilesSlice, oldFile.Key)
					}
				}
			}
		}

		// If file doesn't exist at all in old manifest, it must be new
		if !fileFound {
			newFiles = append(newFiles, file)
		}
	}

	// Check for deleted files
	for _, oldFile := range oldManifest.Files {
		fileFound := false

		for _, newFile := range newManifest.Files {
			if oldFile.Key == newFile.Key {
				fileFound = true
			}
		}

		for _, renamedFile := range renamedFilesSlice {
			if renamedFile == oldFile.Key {
				fileFound = true
			}
		}

		if !fileFound {
			deletedFiles = append(deletedFiles, Delete{
				Key: oldFile.Key,
			})
		}
	}

	return Diff{
		Add:    newFiles,
		Rename: renamedFiles,
		Delete: deletedFiles,
	}
}
