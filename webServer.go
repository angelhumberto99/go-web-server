package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/rpc"
	"strconv"

	"./args"
)

func loadFile(filename string) string {
	file, _ := ioutil.ReadFile(filename)
	return string(file)
}

func form(res http.ResponseWriter, req *http.Request) {
	html := "<p>Los siguientes datos han sido agregados exitosamente</p>"+
			"<p><strong>Alumno: </strong>%s</p>" +
			"<p><strong>Materia: </strong>%s</p>" +
			"<p><strong>Calificación: </strong>%s</p>"
	switch req.Method {
	case "POST":
		c, err := rpc.Dial("tcp", ":9999")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer c.Close()

		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		note,_ := strconv.ParseFloat(req.FormValue("note"), 64)
		obj := args.Args{Name: req.FormValue("student"), 
					   Subject: req.FormValue("subject"), 
					   Note: note}

		// invocar agregar calificación por rpc
		var reply string
		err = c.Call("Server.AddNoteBySubject", obj, &reply)
		if err != nil {
			fmt.Println(err)
		}

		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			loadFile("./html/response.html"),
			loadFile("./css/styles.css"),
			fmt.Sprintf(html, 
				req.FormValue("student"),
				req.FormValue("subject"),
				req.FormValue("note")),
		)
	case "GET":
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			loadFile("./html/index.html"),
			loadFile("./css/styles.css"),
		)
	}
}

func studentAVG(res http.ResponseWriter, req *http.Request) {
	html := "<p>El promedio de %s es %f</p>"

	switch req.Method {
	case "POST":
		c, err := rpc.Dial("tcp", ":9999")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer c.Close()

		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}

		// invocar obtener promedio por rpc
		var reply float64
		err = c.Call("Server.GetStudentAVG", req.FormValue("student"), &reply)
		if err != nil {
			fmt.Println(err)
		}
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			loadFile("./html/response.html"),
			loadFile("./css/styles.css"),
			fmt.Sprintf(html, req.FormValue("student"),reply),
		)
	case "GET":
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			loadFile("./html/studentAVG.html"),
			loadFile("./css/styles.css"),
		)
	}
}

func globalAVG(res http.ResponseWriter, req *http.Request) {
	html := "<p>El promedio global es %f</p>"

	switch req.Method {
	case "GET":
		c, err := rpc.Dial("tcp", ":9999")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer c.Close()

		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}

		var reply float64
		// invocar obtener promedio general por rpc
		err = c.Call("Server.AVGsByStudents", 0.0, &reply)
		if err != nil {
			fmt.Println(err)
		}
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			loadFile("./html/response.html"),
			loadFile("./css/styles.css"),
			fmt.Sprintf(html, reply),
		)
	}
}

func subjectAVG(res http.ResponseWriter, req *http.Request) {
	html := "<p>El promedio de la materia %s es %f</p>"

	switch req.Method {
	case "POST":
		c, err := rpc.Dial("tcp", ":9999")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer c.Close()

		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		var reply float64

		// invocar obtener promedio de materia por rpc
		err = c.Call("Server.AVGsBySubjects", req.FormValue("subject"), &reply)
		if err != nil {
			fmt.Println(err)
		}

		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			loadFile("./html/response.html"),
			loadFile("./css/styles.css"),
			fmt.Sprintf(html, req.FormValue("subject"), reply),
		)
	case "GET":
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			loadFile("./html/subjectAvg.html"),
			loadFile("./css/styles.css"),
		)
	}
}

func main() {
	http.HandleFunc("/", form)
	http.HandleFunc("/student-avg", studentAVG)
	http.HandleFunc("/global-avg", globalAVG)
	http.HandleFunc("/subject-avg", subjectAVG)
	fmt.Println("Servidor corriendo en el puerto 9000")
	http.ListenAndServe(":9000", nil)
}