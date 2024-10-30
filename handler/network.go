package handler

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"real-image-solution-2016/model"
	"real-image-solution-2016/util"
)

type NetworkConnector interface {
	Add(id *int)
	List()
	CheckPermission(q string)
	CreateSubDistributorNetwork()
	GetObject() *network
}

type network struct {
	CountryState map[string][]string
	StateCity    map[string][]string
	Current      model.Distributor
	Distributors []model.Distributor
}

func NewNetwork(path string) NetworkConnector {
	d := network{}
	d.LoadCities(path)

	return &d
}

func (d *network) GetObject() *network {
	return d
}

// Loads the cities from the csv file
func (d *network) LoadCities(filename string) (bool, error) {
	csvFile, err := os.Open(filename)
	if err != nil {
		return false, fmt.Errorf("error while reading, %+v", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(bufio.NewReader(csvFile))
	if _, err := reader.Read(); err != nil {
		panic(err)
	}

	d.CountryState = make(map[string][]string)
	d.StateCity = make(map[string][]string)

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		country := row[5]
		province := row[4]
		city := row[3]

		if _, ok := d.CountryState[country]; !ok {
			d.CountryState[country] = make([]string, 0)
		}

		if !util.Contains(d.CountryState[country], province) {
			d.CountryState[country] = append(d.CountryState[country], province)
		}

		if _, ok := d.StateCity[province]; !ok {
			d.StateCity[province] = make([]string, 0)
		}

		if !util.Contains(d.StateCity[province], city) {
			d.StateCity[province] = append(d.StateCity[province], city)
		}
	}

	return true, nil
}

func (d *network) List() {
	fmt.Println("-> Distributor List: ")
	for _, distributor := range d.Distributors {
		fmt.Printf("-> %d) '%s' has permission to access: \n", distributor.ID, distributor.Name)
		fmt.Printf("Permitted Places: %d, %s\n", len(distributor.PermittedPlaces), strings.Join(distributor.PermittedPlaces, ", "))

		exists := "NO"
		if distributor.HasSubDistributor {
			exists = "YES"
		}
		fmt.Printf("Sub Distributor: %s\n", exists)

		if distributor.Parent != "" {
			fmt.Printf("Parent: %s\n", distributor.Parent)
		} else {
			fmt.Println("Parent: NONE")
		}

		if len(distributor.Child) > 0 {
			fmt.Printf("Children: %s\n", strings.Join(distributor.Child, ", "))
		} else {
			fmt.Println("Children: NONE")
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func (d *network) Add(id *int) {
	var name string
	for {
		readName(&name)
		if name != "" {
			break
		}
		fmt.Println("Distributor name cannot be empty")
	}

	*id++
	d.Current.ID = *id
	d.Current.Name = name
	d.Distributors = append(d.Distributors, d.Current)
	fmt.Printf("-> Now add permissions for %s\n", d.Current.Name)
	for {
		var permission string
		fmt.Println("Enter permission(INCLUDE/EXCLUDE): REGION or press q for Main menu | Ex: INCLUDE: INDIA or EXCLUDE: KARNATAKA-INDIA")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		permission = scanner.Text()
		if permission == "q" {
			break
		}
		if permission == "" {
			fmt.Println("Permissions cannot be empty")
			fmt.Println("")
			continue
		}
		if !(strings.Contains(permission, "INCLUDE") || strings.Contains(permission, "EXCLUDE")) {
			fmt.Println("invalid input")
			fmt.Println("")
			continue
		}

		data := strings.Split(permission, ":")
		prefix := strings.TrimSpace(data[0])
		sufix := strings.TrimSpace(data[1])

		switch prefix {
		case "INCLUDE":
			d.include(sufix, *id-1)
		case "EXCLUDE":
			d.exclude(sufix, *id-1)
		default:
			fmt.Println("invalid input")
		}
	}

	fmt.Println("-> Distributor added successfully")
	fmt.Printf("DDDD: %+v", d.Distributors)
}

func (n *network) CheckPermission(q string) {
	for {
		if q == "" {
			fmt.Println("-> Enter Distributor Name: or press q for Main menu")
			var name string
			fmt.Scanln(&name)
			if name == "q" {
				break
			}

			exists := false
			dist := model.Distributor{}
			for _, d := range n.Distributors {
				if d.Name == name {
					exists = true
					dist = d
					break
				}
			}
			if !exists {
				fmt.Println("Distributor not found")
				break
			}
			n.Current = dist

			fmt.Println("-> Enter your query to check permission: ")
			fmt.Scanln(&q)
		}

		valid := n.verify(q)
		if valid {
			fmt.Println("")
			fmt.Println("YES")
			fmt.Println("")
		} else {
			fmt.Println("")
			fmt.Println("NO")
			fmt.Println("")
		}
	}
}

func (d *network) CreateSubDistributorNetwork() {
	for {
		fmt.Println("-> Your are currently in the Distributor Network Creation Mode")
		fmt.Println("-> Create Sub Distributor")
		fmt.Println("-> Enter Name of the Sub Distributor: or press 'q' to return to main menu")

		var name string
		fmt.Scanln(&name)
		if name == "q" {
			break
		}

		d.Current = model.Distributor{
			ID:                len(d.Distributors),
			Name:              name,
			HasSubDistributor: true,
			PermittedPlaces:   d.Distributors[len(d.Distributors)-1].PermittedPlaces,
		}

		fmt.Println("")
		fmt.Println("~~~ Create a network between distributors ~~~")
		fmt.Println("-> Enter the name of the distributors you want to connect in the following FORMAT: ChildDistributor<-ParentDistributor | ex : Dist2<-Dist1")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		rawSubNetwork := scanner.Text()
		if rawSubNetwork == "q" {
			break
		}

		// Split the input into parent and child distributors
		var parentName, childName string
		if strings.Contains(rawSubNetwork, "<-") {
			data := strings.Split(rawSubNetwork, "<-")
			parentName = strings.TrimSpace(data[1])
			childName = strings.TrimSpace(data[0])
		} else {
			fmt.Println("Invalid Format")
			return
		}

		d.Current.Parent = parentName
		d.Distributors = append(d.Distributors, d.Current)

		// Check if the parent distributor exists
		d.Distributors[len(d.Distributors)-1].Child = append(d.Distributors[len(d.Distributors)-1].Child, childName)
		d.Distributors[len(d.Distributors)-1].HasSubDistributor = true

		fmt.Println("")
		fmt.Printf("Added network connection between parent -> %s and child -> %s\n", parentName, childName)
		fmt.Println("DDD:", d.Distributors)
	}
}

func (n *network) verify(query string) bool {
	places := strings.Split(query, "-")
	for _, include := range n.Current.PermittedPlaces {
		if include == places[0] {
			return true
		}
	}

	return false
}

func (n *network) include(include string, id int) {
	for _, distributor := range n.Distributors {
		if distributor.ID == id {
			n.Current = distributor
		}
	}

	includeSlice := strings.Split(include, "-")

	switch len(includeSlice) {
	case 1:
		{
			n.Current.PermittedPlaces = append(n.Current.PermittedPlaces, includeSlice[0])

			// storing the state in the distributor include
			n.Current.PermittedPlaces = append(n.Current.PermittedPlaces, n.CountryState[includeSlice[0]]...)

			for _, state := range n.CountryState[includeSlice[0]] {
				n.Current.PermittedPlaces = append(n.Current.PermittedPlaces, n.StateCity[state]...)
			}
		}
	case 2:
		{
			n.Current.PermittedPlaces = append(n.Current.PermittedPlaces, includeSlice[0])

			for _, state := range n.CountryState[includeSlice[1]] {
				n.Current.PermittedPlaces = append(n.Current.PermittedPlaces, n.StateCity[state]...)
			}
		}
	case 3:
		{
			n.Current.PermittedPlaces = append(n.Current.PermittedPlaces, includeSlice[0])
		}
	default:
		fmt.Println("Invalid Input, Please try again")
	}

	n.Distributors[id].PermittedPlaces = n.Current.PermittedPlaces
}

func (d *network) exclude(exclude string, id int) {
	for _, distributor := range d.Distributors {
		if distributor.ID == id {
			d.Current = distributor
		}
	}

	excludeSlice := strings.Split(exclude, "-")
	switch len(excludeSlice) {
	case 1:
		{
			for i, value := range d.Current.PermittedPlaces {
				if value == excludeSlice[0] {
					for _, state := range d.CountryState[value] {
						for _, city := range d.StateCity[state] {
							for j, value := range d.Current.PermittedPlaces {
								if j >= 0 && j < len(d.Current.PermittedPlaces) {
									if value == city {
										d.Current.PermittedPlaces = append(d.Current.PermittedPlaces[:j], d.Current.PermittedPlaces[j+1:]...)
									}
								}
							}
						}
						for j, value := range d.Current.PermittedPlaces {
							if j >= 0 && j < len(d.Current.PermittedPlaces) {
								if value == state {
									d.Current.PermittedPlaces = append(d.Current.PermittedPlaces[:j], d.Current.PermittedPlaces[j+1:]...)
								}
							}
						}
					}
					if i >= 0 && i < len(d.Current.PermittedPlaces) {
						d.Current.PermittedPlaces = append(d.Current.PermittedPlaces[:i], d.Current.PermittedPlaces[i+1:]...)
					}
				}
			}
		}
	case 2:
		{
			for i, value := range d.Current.PermittedPlaces {
				if value == excludeSlice[0] {
					d.Current.PermittedPlaces = append(d.Current.PermittedPlaces[:i], d.Current.PermittedPlaces[i+1:]...)
				}
			}

			for _, city := range d.StateCity[excludeSlice[1]] {
				for i, value := range d.Current.PermittedPlaces {
					if value == city {
						d.Current.PermittedPlaces = append(d.Current.PermittedPlaces[:i], d.Current.PermittedPlaces[i+1:]...)
					}
				}
			}
		}
	case 3:
		{
			// we need to remove the city from the distributor include
			for i, value := range d.Current.PermittedPlaces {
				if value == excludeSlice[0] {
					d.Current.PermittedPlaces = append(d.Current.PermittedPlaces[:i], d.Current.PermittedPlaces[i+1:]...)
				}
			}
		}
	default:
		fmt.Println("Invalid input, please try again!")
	}

	d.Distributors[id].PermittedPlaces = d.Current.PermittedPlaces
}

func readName(s *string) {
	fmt.Println("")
	fmt.Println("-> Enter Distributor Name: ")
	fmt.Scanln(s)
}
