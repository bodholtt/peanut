# database

The database package handles all interaction between the 
backend and the Postgres database.

If you want to use an existing Postgres instance, please be sure
to create a user and a database for this application.

```
database:
  host: <database-url>
  port: <database-port>
  username: <database-user>
  password: <database-password
  database_name: <database-name>
  params: "?sslmode=disable"
```