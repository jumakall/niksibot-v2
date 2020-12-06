package player

type Sound struct {
	Name string
	File string
	Path string
}

func CreateSound(name string, filename string, path string) *Sound {
	return &Sound{
		Name: name,
		File: filename,
		Path: path,
	}
}

func (s Sound) PathToFile() string {
	return s.Path + "/" + s.File
}
