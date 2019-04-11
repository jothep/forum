package main

import(
  "encoding/json"
  "fmt"
  "os"
  "log"
  "net/http"
  "strings"
  "html/template"
  "forum/data"
  "errors"
)
type Configuration struct {
  Address string
  ReadTimeout int64
  WriteTimeout int64
  Static  string
}

var config Configuration
var logger *log.Logger

//print out interfaces ?
func p(a ...interface{}){
  fmt.Println(a)
}

func loadConfig(){
  jsonFile, err := os.Open("config.json")
  if err != nil{
    log.Fatalln("Cannot open config file", err)
  }
  fmt.Println("Successfully Opened config.json")
  defer jsonFile.Close()

  config = Configuration{}
  err = json.NewDecoder(jsonFile).Decode(&config)

  if err != nil {
    log.Fatalln("Cannot get configuration from file", err)
  }
}

func init(){
  loadConfig()
  file, err := os.OpenFile("jaskochat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
  if err != nil {
    log.Fatalln("Failed to open log file", err)
  }
  logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile )
  defer file.Close()
}

//redirect url to error message
func error_message(writer http.ResponseWriter, request *http.Request, msg string) {
	url := []string{"/err?msg=", msg}
	http.Redirect(writer, request, strings.Join(url, ""), 302)
}

//check if user has session by cookie

func session(writer http.ResponseWriter, request *http.Request) (sess data.Session, err error) {
	cookie, err := request.Cookie("_cookie")
	if err == nil {
		sess = data.Session{Uuid: cookie.Value}
		if ok, _ := sess.Check(); !ok {
			err = errors.New("Invalid session")
		}
	}
	return
}


//template parse and generate
func parseTemplateFiles(filenames ...string) (t *template.Template) {
	var files []string
	t = template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	t = template.Must(t.ParseFiles(files...))
	return
}

func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(writer, "layout", data)
}

// for logging
func info(args ...interface{}) {
	logger.SetPrefix("INFO ")
	logger.Println(args...)
}

func danger(args ...interface{}) {
	logger.SetPrefix("ERROR ")
	logger.Println(args...)
}

func warning(args ...interface{}) {
	logger.SetPrefix("WARNING ")
	logger.Println(args...)
}

// version
func version() string {
	return "0.1"
}
