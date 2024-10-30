package main

import (
	"fmt"
	"os"

	"real-image-solution-2016/handler"
	"real-image-solution-2016/util"
)

func main() {
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~| Welcome to Real Image Challenge |~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	fmt.Println("			")
	var id int = 0
	network := handler.NewNetwork("./data/cities.csv")

	for {
		fmt.Println("~~~~~~~~ MAIN MENU ~~~~~~~~")
		util.GetMainMenu()
		var arg int
		fmt.Scanln(&arg)

		switch arg {
		case 1:
			fmt.Println("")
			fmt.Println("~~~~ LIST DISTRIBUTOR ~~~~")
			network.List()
		case 2:
			fmt.Println("")
			fmt.Println("~~~~ ADD A DISTRIBUTOR WITH REQUIRED PERMISSIONS ~~~~")
			network.Add(&id)
		case 3:
			fmt.Println("")
			fmt.Println("~~~~ CHECK DISTRIBUTOR PERMISSIONS ~~~~")
			network.CheckPermission("")
		case 4:
			fmt.Println("")
			fmt.Println("~~~~ CREATE THE NETWORK OF DISTRIBUTORS ~~~~")
			network.CreateSubDistributorNetwork()
		case 5:
			fmt.Println("")
			fmt.Println("~~~~ BACK TO THE MAIN MENU ~~~~")
			util.GetMainMenu()
		case 6:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid Choice, please try again")
		}
	}
}
