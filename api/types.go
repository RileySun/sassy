package api


type Model struct {
	ID int
	Name string
	Desc string
}

type Image struct {
	ID int
	Model_ID int
	Path string
	Desc string
}

type Video struct {
	ID int
	Model_ID int
	Path string
	Desc string
}