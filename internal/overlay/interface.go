package overlay

type Overlay interface {
	Load(imagePath string) error
}
