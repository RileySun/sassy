package main

import(
	"fmt"

	"api"
)

func main() {
	api := api.NewAPI()
	//fmt.Println(api.GetVideosBy("model_id", 0))
	//_ = api.AddModel("Riley", "Programmer")
	//something wrong here.
	model := api.GetModelBy("id", 6)
	fmt.Println(model)
	//api.DeleteModel()
}