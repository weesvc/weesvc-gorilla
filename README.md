# WeeSVC Gorilla
Implementation of the WeeSVC application using [Go](https://golang.org/) and the [Gorilla Mux](https://www.gorillatoolkit.org/pkg/mux) 
web toolkit.

## Ingredients
The following external libraries were *directly* utilized in this project.

| Package     | Link                                   | Description                                                             |
| ---         | ---                                    | ---                                                                     |
| Go          | https://golang.org/                    | Well...it's obvious isn't it?!                                          |
| Gorilla Mux | https://www.gorillatoolkit.org/pkg/mux | Web Toolkit providing HTTP server and routing                           |
| GORM        | https://gorm.io/                       | Database ORM                                                            |
| SQLite      | https://www.sqlite.org/index.html      | The lightweight database                                                |
| Cobra       | https://github.com/spf13/cobra         | Command-line library                                                    |
| Viper       | https://github.com/spf13/viper         | Awesome configuration library for settings                              |
| Logrus      | https://github.com/sirupsen/logrus     | Logging abstraction for the Go standard library                         |
| UUID        | https://github.com/google/uuid         | Implementation for generation of universally unique identifiers (UUIDs) |

## Build
Builds are performed using the `Makefile` provided in the project root.  

#### CLI
In order to build the CLI, you will need to have Go (1.17+) installed on your system.

The default target for the `Makefile` will perform several tasks: 
* organize imports using `goimports`
* format code using `gofmt`
* perform linting using `golint`
* vet code for errors using `go vet`
* compile binary for the current platform

:point_up: NOTE: To initially build the project, you may need to run the `make setup` command to install the tools utilized for builds.

Once built, you can **migrate** the database scripts and run the application:
```shell script
$ bin/weesvc migrate; bin/weesvc serve
```
#### Docker
For those who do not have Go available, [Docker](https://hub.docker.com/) is an option to build the application and run 
the application within a container.  Using the `make build-docker` command will build the application within a Linux
container, then copy the resulting binary into a slim docker image to utilize for execution.

Once the the image is available, you can simply run the provided script which will open a browser to access the service
at http://localhost:9092/api/hello .

```shell script
$ ./docker-run.sh
```
:point_up: NOTE: the `docker-run.sh` script is setup to **not** maintain state between executions.  This means each
time you start the container, you will be starting with a freshly created database.

## Using the Application
Update the `DatabaseURI` setting in your `config.yaml` for the absolute path to the base project directory, 
i.e. the path for the directory containing this README.

:point_up: TIP: Use the very cool [HTTPie](https://httpie.org/) application for testing locally from the command-line.  

Execute a `GET` command to retrieve the available _places_ from the database.
```shell script
$ http GET localhost:9092/api/places

HTTP/1.1 200 OK
Content-Length: 2
Content-Type: application/json
Date: Sat, 25 Jan 2020 05:33:57 GMT

[]
```
Add a _place_ into the database using a `POST` command.
```shell script
$ http POST localhost:9092/api/places name=NISC desc="NISC Lake St. Louis Office" lat:=38.7839 lon:=90.7878

HTTP/1.1 200 OK
Content-Length: 8
Content-Type: application/json
Date: Sat, 25 Jan 2020 05:34:08 GMT

{
    "id": 1
}
```
Run the `GET` command again to retrieve _places_ which now include your newly added _place_!
```shell script
$ http GET localhost:9092/api/places/1

HTTP/1.1 200 OK
Content-Length: 217
Content-Type: application/json
Date: Sat, 25 Jan 2020 05:34:18 GMT

[
    {
        "CreatedAt": "2020-01-24T23:34:08.491999-06:00",
        "DeletedAt": null,
        "Description": "NISC Lake St. Louis Office",
        "ID": 1,
        "Latitude": 38.7839,
        "Longitude": 90.7878,
        "Name": "NISC",
        "UpdatedAt": "2020-01-24T23:34:08.491999-06:00"
    }
]
```
Use the `PATCH` command to update a specific value.  For example we'll update the `Description` as follows:
```shell script
$ http PATCH localhost:9092/api/places/1 desc="Lake St. Louis"

HTTP/1.1 200 OK
Content-Length: 203
Content-Type: application/json
Date: Sat, 25 Jan 2020 18:13:13 GMT

{
    "CreatedAt": "2020-01-24T23:34:08.491999-06:00",
    "DeletedAt": null,
    "Description": "Lake St. Louis",
    "ID": 1,
    "Latitude": 38.7839,
    "Longitude": 90.7878,
    "Name": "NISC",
    "UpdatedAt": "2020-01-25T12:13:13.351201-06:00"
}
```
This returns the newly "patched" version of the _place_.  Next we'll remove the row using the `DELETE` method.
```shell script
$ http DELETE localhost:9092/api/places/1

HTTP/1.1 200 OK
Content-Length: 21
Content-Type: application/json
Date: Sat, 25 Jan 2020 18:15:16 GMT

{
    "message": "removed"
}
```
