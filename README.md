# local-webserver

Goose Migration Workflow:

1. Write a new migration file, it must have (-- +goose Up) to add to the migration and (-- +goose Down) to Drop. The file must be ordered so that goose knows what to parse, for instance an ALTER TABLE should be after the CREATE TABLE file.

2. Run goose up in the files dir to apply it to the database (goose postgres "connection-string" up) Sometimes a reset is required to parse a new migration, run (goose postgres "connection_string" down-to 0) then run (goose up).

3. Update your SQL queries in sql/queries/ (calls sqlc to create a go function, call in this format 
(-- name: function_name :(one/many/exec)))

4. Run sqlc generate in the root directory to regenerate the Go database code

Go Request with JSON:

1. Create a requestBody Struct that holds the request parameters such as the body with JSON tags, 
    type requestBody Struct {
        Body string `json:"body"`
    }

2. Initialize an empty struct of requestBody to be filled in with the decoded values,
    reqBody := requestBody{}

3. Create New decoder for the Request Body with,
    decoder := json.NewDecoder(request.Body)

4. Decode the request.Body with a pointer to the struct, the pointer is added the decoded values are parsed to the original struct,
    err := decoder.Decode(&reqBody)

Go Response with JSON:

1. Encode the data into JSON and check for any problems with: 
    response, err := json.Marshal(data)

2. Set the headers, status codes and response to the ResponseWriter.
