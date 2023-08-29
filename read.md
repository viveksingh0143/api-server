github.com/vamika-digital/wms-api-server
|-- cmd
|   |-- root.go
|   |-- server.go
|-- config
|   |-- config.go
|   |-- config.yml
|-- internal
|   |-- app
|       |-- user
|           |-- domain
|               |-- user.go
|           |-- repository
|               |-- user_repository.go
|           |-- usecase
|               |-- user_usecase.go
|               |-- user_usecase_interface.go
|       |-- interface
|           |-- rest
|               |-- user_handler.go
|           |-- middlewares
|               |-- authentication.go
|               |-- authorization.go
|-- pkg
|   |-- database
|       |-- mysql.go
|       |-- postgresql.go
|       |-- sqlite3.go
|   |-- server
|       |-- server.go
