# local-webserver

Goose Migration Workflow:
1. Write a new migration file
2. Run goose up to apply it to the database (goose postgres "connection-string" up)
3. Update your SQL queries in sql/queries/ (calls sqlc to create a go function, call in this format -- name: function_name :(one/many/exec))
4. Run sqlc generate in the root directory to regenerate the Go database code