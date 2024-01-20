# WeeSVC Gorilla
Implementation of the WeeSVC application using [Go](https://golang.org/) and the [Gorilla Mux](https://www.gorillatoolkit.org/pkg/mux) web toolkit.

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
In order to build the CLI, you will need to have Go (1.21+) installed on your system.

The default target for the `Makefile` will perform several tasks: 
* organize imports using `goimports`
* format code using `gofmt`
* vet code for errors using `go vet`
* compile binary for the current platform

> [!NOTE]
> To initially build the project, you may need to run the `make setup` command to install the tools utilized for builds.

Once built using `make`, you can **migrate** the database scripts and run the application:
```shell script
bin/weesvc migrate; bin/weesvc serve
```
#### Docker
For those who do not have Go available, [Docker](https://hub.docker.com/) is an option to build the application and run the application within a container.
Using the `make build-docker` command will build the application within a Linux container, then copy the resulting binary into a slim docker image to utilize for execution.

Once the image is available, you can simply run the provided script which will open a browser to access the service at http://localhost:9092/api/hello .

```shell script
./docker-run.sh
```
> [!IMPORTANT]
> The `docker-run.sh` script is setup to **not** maintain state between executions.
> This means each time you start the container, you will be starting with a freshly created database.

## Using the Application
Update the `DatabaseURI` setting in your `config.yaml` for the absolute path to the base project directory, i.e. the path for the directory containing this README.

> [!TIP]
> Use the very cool [HTTPie](https://httpie.org/) application for testing locally from the command-line.

Execute a `GET` command to retrieve the available _places_ from the database.
```shell script
http GET :9092/api/places
```
```shell
HTTP/1.1 200 OK
Content-Length: 2
Content-Type: application/json
Date: Sat, 25 Jan 2020 05:33:57 GMT

[]
```
Add a _place_ into the database using a `POST` command.
```shell script
http POST :9092/api/places name=NISC description="NISC Lake St. Louis Office" latitude:=38.7839 longitude:=90.7878
```
```shell
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
http GET :9092/api/places/1
```
```shell
HTTP/1.1 200 OK
Content-Length: 217
Content-Type: application/json
Date: Sat, 25 Jan 2020 05:34:18 GMT

[
    {
        "created_at": "2020-01-24T23:34:08.491999-06:00",
        "description": "NISC Lake St. Louis Office",
        "id": 1,
        "latitude": 38.7839,
        "longitude": 90.7878,
        "name": "NISC",
        "updated_at": "2020-01-24T23:34:08.491999-06:00"
    }
]
```
Use the `PATCH` command to update a specific value.
For example, we'll update the `Description` as follows:
```shell script
http PATCH :9092/api/places/1 description="Lake St. Louis"
```
```shell
HTTP/1.1 200 OK
Content-Length: 203
Content-Type: application/json
Date: Sat, 25 Jan 2020 18:13:13 GMT

{
    "created_at": "2020-01-24T23:34:08.491999-06:00",
    "description": "Lake St. Louis",
    "id": 1,
    "latitude": 38.7839,
    "longitude": 90.7878,
    "name": "NISC",
    "updated_at": "2020-01-25T12:13:13.351201-06:00"
}
```
This returns the newly "patched" version of the _place_.
Next we'll remove the row using the `DELETE` method.
```shell script
http DELETE :9092/api/places/1
```
```shell
HTTP/1.1 200 OK
Content-Length: 21
Content-Type: application/json
Date: Sat, 25 Jan 2020 18:15:16 GMT

{
    "message": "removed"
}
```

## API Compliance
A core requirement for all _WeeSVC_ implementations is to implement the same API which are utilized for benchmark comparisons.
To ensure compliance with the required API, [k6](https://k6.io/) is utilized within the [Workbench](https://github.com/weesvc/workbench) project.

To be a valid service, the following command MUST pass at 100%:
```shell script
k6 run -e PORT=9092 https://raw.githubusercontent.com/weesvc/workbench/main/scripts/api-compliance.js
```
