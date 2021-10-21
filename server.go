package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"

	"./args"
)

type Server struct {
	// Alumno : Materia : Calificación
	Students map[string]map[string]float64
	// Materia : Alumno : Calificación
	Subjects map[string]map[string]float64
}

// agregar la calificación de un alumno por materia
func (s *Server) AddNoteBySubject(data args.Args, resp *string) error {
	// materia : calificación
	sbj := make(map[string]float64)
	stdt := make(map[string]float64)
	sbj[data.Subject] = data.Note
	stdt[data.Name] = data.Note

	if _,err := s.Students[data.Name]; err {
		s.Students[data.Name][data.Subject] = data.Note
	} else {
		s.Students[data.Name] = sbj
	}

	if _,err := s.Subjects[data.Subject]; err {
		s.Subjects[data.Subject][data.Name] = data.Note
	} else {
		s.Subjects[data.Subject] = stdt
	}

	// relación estudiantes, materia y calificación
	fmt.Println("Alumnos")
	for name,_ := range s.Students {
		fmt.Println("  >",name)
		for subject,note := range s.Students[name] {
			fmt.Printf("    - %s : %f\n", subject, note)
		}
	}

	*resp = "Calificación guardada"
	return nil
}

// obtener el promedio del alumno
func (s *Server) GetStudentAVG(name string, reply *float64) error {
	var avg float64 = 0
	var size float64 = 0

	if _, err := s.Students[name]; !err {
		return errors.New("El Alumno no existe")
	}

	for _, note := range s.Students[name] {
		avg += note
		size++
	}
	*reply = (avg/size)
	return nil
}

// obtener el promedio de todos los alumnos
func (s *Server) AVGsByStudents(size float64, reply *float64) error {
	var avg float64 = 0
	var gAvg float64 = 0
	var ammount float64 = 0
	if len(s.Students) == 0 {
		return errors.New("No hay calificaciones por mostrar")
	}

	for name,_ := range s.Students {
		avg = 0
		size = 0
		for _,note := range s.Students[name] {
			avg += note
			size++
		}
		gAvg += avg/size
		ammount++
	}
	*reply = gAvg/ammount
	return nil
}

// obtener el promedio por materia
func (s *Server) AVGsBySubjects(name string, reply *float64) error {
	var avg float64 = 0
	var size float64 = 0
	
	if _,err := s.Subjects[name]; !err{
		return errors.New("No existe dicha materia")
	}

	for _, note := range s.Subjects[name] {
		avg += note
		size++
	}
	*reply = (avg/size)

	return nil
}

func server() {
	m1 := make(map[string]map[string]float64)
	m2 := make(map[string]map[string]float64)
	s := new(Server)
	s.Students = m1
	s.Subjects = m2
	rpc.Register(s)
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
	}
	defer ln.Close()
	
	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(c)
	}
}

func main() {
	go server()
	var input string	
	fmt.Scanln(&input)
}