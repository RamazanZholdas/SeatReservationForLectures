package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Teachers struct {
	Techers []Teacher `json:"teachers"`
}

type Teacher struct {
	Name   string   `json:"name"`
	Age    string   `json:"age"`
	Bio    string   `json:"bio"`
	Places []Places `json:"places"`
}

type Places struct {
	Id      int    `json:"id"`
	Level   string `json:"level"`
	Ordered bool   `json:"ordered"`
}

type CongratsTicket struct {
	Name        string
	Number      int
	Level       string
	AlreadyTrue bool
}

const (
	port          string = ":8080"
	filename      string = "jsonDb/teachers.json"
	filenameMySql string = "jsonDb/teachersMySql.json"
	filenameGo    string = "jsonDb/teachersGo.json"
)

var (
	allData             Teachers
	allDataMySql        Teachers
	allDataGo           Teachers
	chosenPythonTeacher string
	chosenMySqlTeacher  string
	chosenGoTeacher     string
	index               int
	indexMySql          int
	indexGo             int
	AlreadyTrue         bool = false
	AlreadyTrueMySql    bool = false
	AlreadyTrueGo       bool = false
)

func init() {
	createJsonDB(filename, "python")
	createJsonDB(filenameMySql, "mySql")
	createJsonDB(filenameGo, "go")
	allData = readJsonDb(filename)
	allDataMySql = readJsonDb(filenameMySql)
	allDataGo = readJsonDb(filenameGo)
}

func main() {
	http.HandleFunc("/main/", mainPage)
	http.HandleFunc("/pythonPage/", pythonPage)
	http.HandleFunc("/pythonPage/step2/", choosePlacePython)
	http.HandleFunc("/pythonPage/step2/pythonCongrats/", pythonCongrats)
	http.HandleFunc("/mySqlPage/", mySqlPage)
	http.HandleFunc("/mySqlPage/step2/", choosePlaceMySql)
	http.HandleFunc("/mySqlPage/step2/mySqlCongrats/", mySqlCongrats)
	http.HandleFunc("/goPage/", goPage)
	http.HandleFunc("/goPage/step2/", choosePlaceGo)
	http.HandleFunc("/goPage/step2/goCongrats/", goCongrats)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Can't launch server\nError:%v", err)
	}
}

func goCongrats(rw http.ResponseWriter, rq *http.Request) {
	tm, err := template.ParseFiles("goPages/goCongratsPage.html")
	cannotRunTemplate("goPages/goCongratsPage.html", err)
	ticket := CongratsTicket{
		chosenGoTeacher,
		allDataGo.Techers[indexGo].Places[indexGo].Id,
		allDataGo.Techers[indexGo].Places[indexGo].Level,
		AlreadyTrueGo,
	}
	tm.Execute(rw, ticket)
}

func mySqlCongrats(rw http.ResponseWriter, rq *http.Request) {
	tm, err := template.ParseFiles("mySqlPages/mySqlCongratsPage.html")
	cannotRunTemplate("mySqlPages/mySqlCongratsPage.html", err)
	ticket := CongratsTicket{
		chosenMySqlTeacher,
		allDataMySql.Techers[indexMySql].Places[indexMySql].Id,
		allDataMySql.Techers[indexMySql].Places[indexMySql].Level,
		AlreadyTrueMySql,
	}
	tm.Execute(rw, ticket)
}

func pythonCongrats(rw http.ResponseWriter, rq *http.Request) {
	tm, err := template.ParseFiles("pythonPages/pythonCongratsPage.html")
	cannotRunTemplate("pythonPages/pythonCongratsPage.html", err)
	ticket := CongratsTicket{
		chosenPythonTeacher,
		allData.Techers[index].Places[index].Id,
		allData.Techers[index].Places[index].Level,
		AlreadyTrue,
	}
	tm.Execute(rw, ticket)
}

func choosePlaceGo(rw http.ResponseWriter, rq *http.Request) {
	if chosenGoTeacher == "" {
		rw.Write([]byte("Go back u didn't choose teacher"))
		return
	}
	for i := range allDataGo.Techers {
		if allDataGo.Techers[i].Name == chosenGoTeacher {
			tm, err := template.ParseFiles("goPages/goPageS2.html")
			cannotRunTemplate("goPages/goPageS2.html", err)
			tm.Execute(rw, allDataGo.Techers[i].Places)
			indexGo = i
			break
		}
	}
	if len(rq.FormValue("smth")) != 0 {
		for i := range allDataGo.Techers[indexGo].Places {
			str := rq.FormValue("smth")
			num, _ := strconv.Atoi(str[len(str)-1:])
			if allDataGo.Techers[indexGo].Places[i].Id == num {
				if allDataGo.Techers[indexGo].Places[i].Ordered {
					AlreadyTrueGo = true
					return
				}
				allDataGo.Techers[indexGo].Places[i].Ordered = true
				AlreadyTrueGo = false
				saveDataJson(allDataGo, filenameGo)
				break
			}
		}
	}
}

func choosePlaceMySql(rw http.ResponseWriter, rq *http.Request) {
	if chosenMySqlTeacher == "" {
		rw.Write([]byte("Go back u didn't choose teacher"))
		return
	}
	for i := range allDataMySql.Techers {
		if allDataMySql.Techers[i].Name == chosenMySqlTeacher {
			tm, err := template.ParseFiles("mySqlPages/mySqlPageS2.html")
			cannotRunTemplate("mySqlPages/mySqlPageS2.html", err)
			tm.Execute(rw, allDataMySql.Techers[i].Places)
			indexMySql = i
			break
		}
	}
	if len(rq.FormValue("smth")) != 0 {
		for i := range allDataMySql.Techers[indexMySql].Places {
			str := rq.FormValue("smth")
			num, _ := strconv.Atoi(str[len(str)-1:])
			if allDataMySql.Techers[indexMySql].Places[i].Id == num {
				if allDataMySql.Techers[indexMySql].Places[i].Ordered {
					AlreadyTrueMySql = true
					return
				}
				allDataMySql.Techers[indexMySql].Places[i].Ordered = true
				AlreadyTrueMySql = false
				saveDataJson(allDataMySql, filenameMySql)
				break
			}
		}
	}
}

func choosePlacePython(rw http.ResponseWriter, rq *http.Request) {
	if chosenPythonTeacher == "" {
		rw.Write([]byte("Go back u didn't choose teacher"))
		return
	}
	for i := range allData.Techers {
		if allData.Techers[i].Name == chosenPythonTeacher {
			tm, err := template.ParseFiles("pythonPages/pythonPageS2.html")
			cannotRunTemplate("pythonPages/pythonPageS2.html", err)
			tm.Execute(rw, allData.Techers[i].Places)
			index = i
			break
		}
	}
	if len(rq.FormValue("smth")) != 0 {
		for i := range allData.Techers[index].Places {
			str := rq.FormValue("smth")
			num, _ := strconv.Atoi(str[len(str)-1:])
			if allData.Techers[index].Places[i].Id == num {
				if allData.Techers[index].Places[i].Ordered {
					AlreadyTrue = true
					return
				}
				allData.Techers[index].Places[i].Ordered = true
				AlreadyTrue = false
				saveDataJson(allData, filename)
				break
			}
		}
	}
}

func goPage(rw http.ResponseWriter, rq *http.Request) {
	tm, err := template.ParseFiles("goPages/goPage.html")
	cannotRunTemplate("goPages/goPage.html", err)
	tm.Execute(rw, allDataGo.Techers)
	chosenGoTeacher = rq.FormValue("smth")
}

func mySqlPage(rw http.ResponseWriter, rq *http.Request) {
	tm, err := template.ParseFiles("mySqlPages/mySqlPage.html")
	cannotRunTemplate("mySqlPages/mySqlPage.html", err)
	tm.Execute(rw, allDataMySql.Techers)
	chosenMySqlTeacher = rq.FormValue("smth")
}

func pythonPage(rw http.ResponseWriter, rq *http.Request) {
	tm, err := template.ParseFiles("pythonPages/pythonPage.html")
	cannotRunTemplate("pythonPages/pythonPage.html", err)
	tm.Execute(rw, allData.Techers)
	chosenPythonTeacher = rq.FormValue("smth")
}

func mainPage(rw http.ResponseWriter, rq *http.Request) {
	tm, err := template.ParseFiles("mainPage.html")
	cannotRunTemplate("mainPage.html", err)
	tm.Execute(rw, nil)
}

func saveDataJson(teachers Teachers, filenameToSaveData string) {
	file, err := os.Open(filenameToSaveData)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	byteSlice, err := json.MarshalIndent(teachers, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filenameToSaveData, byteSlice, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func readJsonDb(filename string) Teachers {
	var slice Teachers
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	byteSlice, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(byteSlice, &slice)
	if err != nil {
		log.Fatal(err)
	}

	return slice
}

func createJsonDB(fileName string, pLang string) {
	var slice Teachers
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if pLang == "python" {
		slice.addSomeData()
	} else if pLang == "mySql" {
		slice.addSomeDataMySql()
	} else {
		slice.addSomeDataGo()
	}

	byteSlice, err := json.MarshalIndent(slice, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(file.Name(), byteSlice, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func cannotRunTemplate(page string, err error) {
	if err != nil {
		log.Fatalf("Cannot run template called %s\nError:%v", page, err)
	}
	return
}

func (teacher *Teachers) addSomeDataGo() {
	teacher.Techers = append(teacher.Techers,
		Teacher{
			Name: "Richard Hendricks",
			Age:  "23",
			Bio:  "Some cv information",
			Places: []Places{
				{
					Id:      1,
					Level:   "Vip",
					Ordered: false,
				},
				{
					Id:      2,
					Level:   "Normal",
					Ordered: false,
				},
				{
					Id:      3,
					Level:   "Normal",
					Ordered: false,
				},
			},
		},
		Teacher{
			Name: "Bertram Gilfoyle",
			Age:  "26",
			Bio:  "Programming experience is over 5 years",
			Places: []Places{
				{
					Id:      1,
					Level:   "Vip",
					Ordered: false,
				},
				{
					Id:      2,
					Level:   "Normal",
					Ordered: false,
				},
				{
					Id:      3,
					Level:   "Normal",
					Ordered: false,
				},
			},
		},
		Teacher{
			Name: "Dinesh Chugtai",
			Age:  "25",
			Bio:  "Just a cool guy",
			Places: []Places{
				{
					Id:      1,
					Level:   "Vip",
					Ordered: false,
				},
				{
					Id:      2,
					Level:   "Normal",
					Ordered: false,
				},
				{
					Id:      3,
					Level:   "Normal",
					Ordered: false,
				},
			},
		},
	)
}

func (teacher *Teachers) addSomeDataMySql() {
	teacher.Techers = append(teacher.Techers,
		Teacher{
			Name: "John Dorian",
			Age:  "24",
			Bio:  "Some cv information",
			Places: []Places{
				{
					Id:      1,
					Level:   "Vip",
					Ordered: false,
				},
				{
					Id:      2,
					Level:   "Normal",
					Ordered: false,
				},
				{
					Id:      3,
					Level:   "Normal",
					Ordered: false,
				},
			},
		},
		Teacher{
			Name: "Christopher Duncan Turk",
			Age:  "26",
			Bio:  "Programming experience is over 5 years",
			Places: []Places{
				{
					Id:      1,
					Level:   "Vip",
					Ordered: false,
				},
				{
					Id:      2,
					Level:   "Normal",
					Ordered: false,
				},
				{
					Id:      3,
					Level:   "Normal",
					Ordered: false,
				},
			},
		},
		Teacher{
			Name: "Percival Ulysses Willis Cox",
			Age:  "35",
			Bio:  "Just a cool guy",
			Places: []Places{
				{
					Id:      1,
					Level:   "Vip",
					Ordered: false,
				},
				{
					Id:      2,
					Level:   "Normal",
					Ordered: false,
				},
				{
					Id:      3,
					Level:   "Normal",
					Ordered: false,
				},
			},
		},
	)
}

func (teacher *Teachers) addSomeData() {
	teacher.Techers = append(teacher.Techers,
		Teacher{
			Name: "Nanami",
			Age:  "30",
			Bio:  "Just a nice guy",
			Places: []Places{
				{
					Id:      1,
					Level:   "Vip",
					Ordered: true,
				},
				{
					Id:      2,
					Level:   "Normal",
					Ordered: false,
				},
				{
					Id:      3,
					Level:   "Normal",
					Ordered: false,
				},
			},
		},
		Teacher{
			Name: "Gojo",
			Age:  "26",
			Bio:  "Just a strong guy",
			Places: []Places{
				{
					Id:      1,
					Level:   "Vip",
					Ordered: false,
				},
				{
					Id:      2,
					Level:   "Normal",
					Ordered: false,
				},
				{
					Id:      3,
					Level:   "Normal",
					Ordered: false,
				},
			},
		},
		Teacher{
			Name: "Aoi Todou",
			Age:  "22",
			Bio:  "Just a cool guy",
			Places: []Places{
				{
					Id:      1,
					Level:   "Vip",
					Ordered: false,
				},
				{
					Id:      2,
					Level:   "Normal",
					Ordered: false,
				},
				{
					Id:      3,
					Level:   "Normal",
					Ordered: false,
				},
			},
		},
	)
}

/*
tm, err := template.ParseFiles("pythonPage.html")
		cannotRunTemplate("pythonPage.html", err)
		err = tm.Execute(rw, file.Places)
		cannotRunTemplate("pythonPage.html", err)

func firstHandler(rw http.ResponseWriter, rq *http.Request) {
	if rq.Method == "GET" {
		t, _ := template.ParseFiles("index.html")
		t.Execute(rw, nil)
	} else {
		rq.ParseForm()
		fmt.Println(rq.Form["body"])
		fmt.Printf("%T", rq.Form["body"])
		http.Redirect(rw, rq, "/redirect", http.StatusSeeOther)
	}
}
*/
