package mod

type GitBlob struct {
	header   string
	blobData string
}

func (gob *GitBlob) serialize() string {
	return gob.blobData
}

func (gob *GitBlob) deserialize(data string) {
	gob.blobData = data
}
