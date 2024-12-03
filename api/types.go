package api

type User struct {
	ID int		//User Id
	Name string	//User Name
	Trial bool		//Is Trial User
	Get, Add, Update, Delete int64 //Limits
}

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