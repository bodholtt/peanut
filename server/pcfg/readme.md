# pcfg

The pcfg package holds the configuration object used by the other packages.

It simply takes the config file, parses it into a struct, 
and creates a public object that can be accessed by the other modules.

The config is stored in `config.yml` in the root directory of the server.

See below for a sample config.yml:
```
server:
  port: 8080 # port where server should listen on
images:
  root_location: "./static" # directory where images are stored
  delete_image_files: false # upon post deletion, remove files if true
database:
  host: "localhost" # locally reachable address of postgres database
  port: 5432 # port
  username: "postgres" # user with access to the database
  password: "password" # ^
  database_name: "imageboard" # name of the database to be used
  params: "?sslmode=disable" # 
client:
  host: "http://localhost:4321" # URL of client, including port if applicable
```