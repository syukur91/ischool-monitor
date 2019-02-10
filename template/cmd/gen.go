// This program generates contributors.go. It can be invoked by running
// go generate
package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"text/template"
)

func main() {

	// Load template
	serviceTemplatePath, _ := filepath.Abs("service.tmpl")
	serviceTemplateContent, err := ioutil.ReadFile(serviceTemplatePath)
	if err != nil {
		log.Fatal(err)
	}
	serviceTemplate := template.Must(template.New("").Parse(string(serviceTemplateContent)))

	controllerTemplatePath, _ := filepath.Abs("controller.tmpl")
	controllerTemplateContent, err := ioutil.ReadFile(controllerTemplatePath)
	if err != nil {
		log.Fatal(err)
	}
	controllerTemplate := template.Must(template.New("").Parse(string(controllerTemplateContent)))

	schemaTemplatePath, _ := filepath.Abs("schema.tmpl")
	schemaTemplateContent, err := ioutil.ReadFile(schemaTemplatePath)
	if err != nil {
		log.Fatal(err)
	}
	schemaTemplate := template.Must(template.New("").Parse(string(schemaTemplateContent)))

	// models that we want to generate
	// WARNING !!. You might want to set skip=false if model already generated and modified
	models := []struct {
		Model          string
		ModelLowerCase string
		ControllerFile string
		ServiceFile    string
		SchemaFile     string
		Skip           bool
	}{
		{
			Skip:           true,
			Model:          "Mata_Pelajaran",
			ModelLowerCase: "mata_pelajaran",
			ControllerFile: "../api/controller/mata_pelajaran.go",
			ServiceFile:    "../service/mata_pelajaran.go",
			SchemaFile:     "../api/schema/mata_pelajaran.go",
		},
		{
			Skip:           true,
			Model:          "Kelas",
			ModelLowerCase: "kelas",
			ControllerFile: "../api/controller/kelas.go",
			ServiceFile:    "../service/kelas.go",
			SchemaFile:     "../api/schema/kelas.go",
		},
		{
			Skip:           true,
			Model:          "Wali_Kelas",
			ModelLowerCase: "wali_kelas",
			ControllerFile: "../api/controller/wali_kelas.go",
			ServiceFile:    "../service/wali_kelas.go",
			SchemaFile:     "../api/schema/wali_kelas.go",
		},
		{
			Skip:           true,
			Model:          "Siswa",
			ModelLowerCase: "siswa",
			ControllerFile: "../api/controller/siswa.go",
			ServiceFile:    "../service/siswa.go",
			SchemaFile:     "../api/schema/siswa.go",
		},
		{
			Skip:           true,
			Model:          "User",
			ModelLowerCase: "user",
			ControllerFile: "../api/controller/user.go",
			ServiceFile:    "../service/user.go",
			SchemaFile:     "../api/schema/user.go",
		},
	}

	// Create file
	for _, v := range models {

		if v.Skip {
			continue
		}

		log.Printf("Generating service and controller for %s", v.Model)
		{
			f, err := os.Create(v.ServiceFile)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			serviceTemplate.Execute(f, v)
		}

		{
			f, err := os.Create(v.ControllerFile)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			controllerTemplate.Execute(f, v)
		}

		{
			f, err := os.Create(v.SchemaFile)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			schemaTemplate.Execute(f, v)
		}
	}

}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
