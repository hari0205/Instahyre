# Submission for my version of Instahyre assignment

### Steps
- Copy the code and paste it in your machine
- Install the dependencies   
```
go get {Insert package name here}
```

or 

```
go mod tidy 
```
- Run the following command
```bash
go run main.go
```
- If required, use 
```
go build 
```

to build the api as an executable/binary file. For the sake of saving space, I have not done it . Depending on OS and CPU architecture , you might want to run ` GOOS` and ` GOARCH` flag during build


### Instructions to use API
- The endpoints required are present in main.go file. Please use the endpoints in Postman/Insomnia
- Use signup endpoint for creating a new user
- Use login endpoint for logging in and authorizing the user
- To use other endpoint, authorization is required.
- The server will throw an error if the user is not authorised and trying to access protected routes
- After login, use the any/all endpoint as per requirement. 

