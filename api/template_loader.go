package main

import "html/template"

const templateFileShortList = "templates/short-list.html"
const templateFileNavigation = "templates/navigation.html"
const templateFileNames = "templates/name-suggestor.html"

func getShortListTemplate() (*template.Template, error) {
	return template.ParseFiles(templateFileShortList, templateFileNavigation)
}

func getNameTemplate() (*template.Template, error) {
	return template.ParseFiles(templateFileNames, templateFileNavigation)
}
