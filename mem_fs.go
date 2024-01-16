package golden

type MemFs struct {
	files map[string][]byte
}

func NewMemFs() *MemFs {
	return &MemFs{
		files: make(map[string][]byte),
	}
}

func (fs *MemFs) Exists(name string) (bool, error) {
	_, ok := fs.files[name]
	if ok {
		return true, nil
	}
	return false, nil
}

func (fs *MemFs) WriteFile(name string, data []byte) error {
	fs.files[name] = data
	return nil
}

func (fs *MemFs) ReadFile(name string) ([]byte, error) {
	content, ok := fs.files[name]
	if ok {
		return content, nil
	}
	return []byte{}, SnapshotNotFound
}
