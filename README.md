# WASA
### This is a repository for the WASAtext project
[*Web And Software Architecture*](https://gamificationlab.uniroma1.it/en/wasa/) is a course I have been attending at Sapienza University of Rome during the academic year 2025/2026 as part of my Erasmus+ exchange program.  

The main goal of the course is to create a web application project called *WASAtext*, which is a basic communication platform that allows users to send messages, create channels, and interact with each other. 

The project is developed stage by stage according to the course syllabus: 
1. define APIs using the OpenAPI standard
2. design and develop the server side (“backend”) in Go
3. design and develop the client side (“frontend”) in JavaScript
4. create a Docker container image for deployment

The structure of the repository is imposed by the course requirements and is based on [fantastic coffee decaffeinated](https://github.com/sapienzaapps/fantastic-coffee-decaffeinated).

----

## 1. OpenAPI

All API endpoints and schemas are defined in the *OpenAPI* specification.  
For a convenient graphical representation of all resources, please open:

[Graphical API Documentation](doc/api-docs.html) built using the [*Redoc CLI*](https://redocly.com/docs/cli)


## 2. Backend and Database
The backend is implemented in Go and is located in the `service` directory, with the main entry point in `cmd/` directory.

The backend uses a SQLite database to store user data, and the database is created using migrations in `service/database/migrations/`.

### here will be a proper graphical architecture diagram soon...
----

The backend exposes the API endpoints defined in the OpenAPI specification and handles requests from the frontend.

According to the Go ideology (as far as I get it), the methods responsible for handling API requests are located in the `service/api` directory, while the database interactions are handled in the `service/database` directory. This part of the code is implemented using interfaces to allow for easy swapping of database implementations if needed.

According to the course requirements, I use bearer tokens for user authentication, which are generated using the `github.com/golang-jwt/jwt/v5` package. 

napisac jaki jestem super i ze robie to nawet na tranzakcjach (czasami jak trzeba)

napisac ze na moje male potrzeby zapisuje zdjecia w bazie jako bloby ale wiem ze to zle i w realnych projektach sie tak nie robi

napisac ze obrabiam zdjecia tak aby byly zapisywane jako kwadraty 200x200px latwe do odczytu i wyslania potel do frontu i ze uzywam do tego biblioteki imaging