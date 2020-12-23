# ServerInformationSearch
 Web API that enables you to search for the information of the servers an specific domain is using to run its services.

## To run the app follow these steps
## In the Frontend folder

## Project setup
```
npm install
```

### Compiles and hot-reloads for development
```
npm run serve
```

### Compiles and minifies for production
```
npm run build
```

### Lints and fixes files
```
npm run lint
```

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).

## In the Backend folder

### You will have to change the CockroachDB credentials for personal ones and update the references in persistence.go file.

Aditionally you will have to create the table in the DB so it matches the queries done in this file. There are three columns in a single table, all Strings; dominio, grade and dateVisited

### Run the Go app
```
go run .
```
