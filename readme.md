# Schedule of Faculty of Physics && Google-Calendar.
This is a little project for shedule of Faculty of Physics(http://ras.phys.msu.ru/table/4/1.htm) and Google-Calendar(https://developers.google.com/calendar) interaction. All my app is written with golang. For storage data i used MySQL.


# Content of repository.
* branch: master (main branch with web-application)
    - `/Parsing` (files for shedule of faculty of physics parsing)
	- `/Schedule` (files for console-representation of shedule)
	- `/debug` (there are not important files)
	- `/WebApp` (main directory with web-application)
	- `TODO.md`
	- `readme.md`
* branch: console (with console-application)

# Opportunities.
The program allows you to upload the schedule of any group(in scope the Faculty of Physics) to your google-calendar using your google-account.

# Packages.
During development, i used following packages:
* For parsing [goqury](https://godoc.org/github.com/fzipp/goquery)
* For google-authorization [oauth2](https://godoc.org/golang.org/x/oauth2)
* [GoogleAPI](https://godoc.org/google.golang.org/api/calendar/v3)
* Standart libs of golang

# How to launch it on local machine?
Firstly, you need to install mysql DB. Secondly, using `/Parsing`, launch `sh Parsing.sh`. Thirdly, using `/WebApp`, launch `sh WebApp.sh`, after you can use `localhost:8080` for using application. If you want to check correctness of my parsing, you can use `/Schedule`, launch `Schedule.sh`for console-representation.

# Using the application with web-connection.
[Click it](https://fopschedule.herokuapp.com).
For deployment, i used [Heroku](https://heroku.com) with [ClearDB](https://www.cleardb.com).

# Essential Requirements.
* google-account

# P.S.
This is a very raw, but working version of my application. I need to upgrade front-end of app(sorry, but i'm only backend developer...). Also, i need to add some features. It's comming soon...
