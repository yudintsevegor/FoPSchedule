package main

import (
	"text/template"
)

var (
	head = `
	<!DOCTYPE html>
	<html>
	  <head>
	  <style>
	  </style>
		<!-- Required meta tags -->
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

		<!-- Bootstrap CSS -->
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta/css/bootstrap.min.css" integrity="sha384-/Y6pD6FV/Vv2HJnA6t+vslU6fwYXjCFtcEpHbNJ0lyAFsXTsjBbfaDjzALeQsN6M" crossorigin="anonymous">
	  </head>
	  <body>
	  <p> Вы зашли с помощью {{.Email}} </p>
 	<form action="http://localhost:8080/result" method="post" enctype="application/x-www-form-urlencoded">
   <p><select name="group">
	`

	end = `
   </select></p>
   <p><input type="submit" value="Отправить"></p>
	</body>
	</html>
	`

	endOpt = `</optgroup>`

	label = template.Must(template.New("").Parse(`
    	<optgroup label="{{.Course}} курс">
	`))

	option = template.Must(template.New("").Parse(`
    	<option value="{{.Group}}">{{.Group}}</option>
	`))
)
