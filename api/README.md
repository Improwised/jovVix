# Golang BoilerPlater

### **Index**
- [Introduction](#introduction)

- [Multiple Database System Support](#multiple-database-system-support)

- [Kick Start Commands](#kick-start-commands)

- [Migrations](#migrations)

- [Kratos Integration](#kratos-integration)

- [Messaging queue](#messaging-queue)

- [Code Walk-through](#code-walk-through)
    - [Config](#config)
    - [Command](#command)
    - [Route](#route)
    - [Middleware](#middleware)
    - [Model](#model)
    - [Controller](#controller)
    - [Utils](#utils)

- [Testcases](#testcases)
---
### **Introduction**
- This template is build on [gofiber](https://github.com/gofiber/fiber).
- To sync go packages you need to run following command.

    ### `go mod vendor`
- Copy `.env.example` as `.env` for environment variables.
- To add new package you can run following command.

    ### `go get {package_url}`
- You need to run `go mod vendor` after running above command to sync your pacakage with vendor directory.
- Following command is use to remove unused package from your `go.mod` and `go.sum` files.

    ### `go mod tidy`
- Please make sure you set proper database configuration to connect.
---
### **Multiple Database System Support**

- We are supporting 3 types of database at this time `postgres`, `mysql` & `sqlite3`

- It allows us to switch to database without changing too much stuff (but it requires to change function that are varies based on system.)

- [goqu](https://github.com/doug-martin/goqu) package allows us to do that.

- we have environment variable `DB_DIALECT` where you can set to database system type.
---
### **Kick Start Commands**
- This commands are define inside `Makefile`.

**Note:** Make sure you change this file if you are going to use this boilerplate with different folder structure.

for ex: `start` commands run `go run app.go api` if your app.go is somewhere else then you need to update this file for paths to work properly.

- **make start-api-dev** : This will start your app with `nodemon` which watch for changes in your directory and reload go.

- **make start-api** : To start api, it basically runs `go run app.go`

- **make migrate file_name={MIGRATION_FILE_NAME}** : This command will create migration `up` and `down` both respectively, you need to pass `file_name={migration_name}`

- **make build app_name={BINARY_NAME}**  : It will build binary of your project.

- **make install app_name={BINARY_NAME}**: It will generate optimized binary with `-s` and `-w` ldflags.

- **make test** : To run testcase for entire project, `.env.testing` env need to use while writing testcase, you need to load

- **make test-wo-cache**: To run testcases with without cache.

- **make swagger-gen** : To genrate swagger docs, This command will download binary if not exist.

- **make migrate-up** : To run `Up` migrations.
---
## Migration

Migrations are like **version control for your database**, allowing your team to define and share the application's database schema definition. If you have ever had to tell a teammate to manually add a column to their local database schema after pulling in your changes from source control, you've faced the problem that database migrations solve.

## Kratos Integration
Ory Kratos provides the user identity management service and different flows for user management (signup/sign in, forgot password, reset password, etc.). For more, you can see the official [documentation](https://www.ory.sh/docs/kratos/ory-kratos-intro).

Ory Kratos doesn't provide UI, You have to specify the endpoints for different UI pages inside the configuration and Kratos will use them. There are some other services that you can use for demo UIs. for example `kratos-selfservice-ui-node`.

**Note:** ory Kratos is an optional integration to the boilerplate, if you want to use it you need to follow the below steps.

- Inside the ```.env``` you'll have to set the ```KRATOS_ENABLED``` for enabling the kratos integration.
- According to your config requirements, you will need to change the corresponding files inside the ```/pkg/kratos``` folder.
- Then after that for all the endpoints you want Kratos authentication you'll have to add ```middlewares.Authenticated```. After that, you can add your own handle and write business logic over there using user details.
- For more details, you can see the [documentation](./pkg/kratos/readme.md) section.

## Execution

1. Run ```docker-compose up``` to spin up the database and admire. In the case of a Kratos Enabled user this command ```docker-compose --profile kratos up```.
2. Open ```localhost:8080```,  select **system** to ```PostgreSQL``` and put username and password.
3. Build image using ```docker build -t jovvix .```
4. Run ```docker run jovvix``` to run the container.

**Another Way:**
- For starting All services including kratos, runinng databse migrations and starting up golang server all together, run the script `local.sh` inside `/pkg/kratos`.
    ```bash
    cd pkg/kratos
    ./local.sh
    ```
**Note:** Use this for local development environment.

### **Migrations**
- **CREATE :** To create migrations i have discribed details in [Kick Start Commands](#kick-start-commands) section.
- **RUN :** To run migration there is two command `make migration-up` && `make migration-down`.
- Migration needs `-- +migrate Up` and `-- +migrate Down` respectively in starting of files, this is required because we are using [sql-migrate](https://github.com/rubenv/sql-migrate) package

### **Messaging Queue**
- We are using [watermill](https://watermill.io/) package for messaging queue.
- Watermill is a Golang library for working efficiently with message streams. It is intended for building event-driven applications.


- #### Multiple Message Queue Broker Support
    - We are supporting 5 types of message queue broker at this time `rabbitmq`, `redis`, `googleCloud`,`sql(postgres,mysql)` & `kafka`
    - It allows us to switch to message queue broker without changing too much stuff.
    - Watermill package allows us to do that.
    - We have environment variable `MQ_DIALECT` where you can set to message queue broker type.
    - Need to change env accoeding to message queue broker.
- #### Creating An Worker
    - All of the workers for your application are stored in the `cli/workers` directory.
    - To create a new job add a new file into the `cli/workers` directory.
- #### Class Structure
    - Workers class are very simple, consisting of a single method `Handle`. `Handle` executes when a message is received.
    - The `Handle()` method should return an `error` if the job fails.
        ```go
        type WelcomeMail struct {
            FirstName string
            LastName  string
            Email     string
            Roles     string
        }
        // Handle executes the job.
        func (w WelcomeMail) Handle() error {
            return nil
        }
        ```
- #### Register Worker
    - After creating the struct, you need to register it in `cli/workers/worker_handler.go`, so that it can be called correctly.
    - To register a new worker add struct to `RegisterWorkerStruct` function.
        ```go
        func RegisterWorkerStruct() []interface{} {
            return []interface{}{
                WelcomeMail{},
                // ...
            }
        }
        ```
 - #### Command to run worker
    ```go
    go run app.go worker --retry-delay 400 --retry-count 3 --topic user
    // --retry-delay 400 --retry-count 3 are optional
    // --retry-delay 400 means it will retry after 400ms
    // --retry-count 3 means it will retry 3 times

    ```

- #### Publish Message
    - The `InitPubliser` function initializes a `WatermillPubliser` based on the provided configuration.
        ```go
        pub, err := watermill.InitPubliser(cfg)
        if err != nil {
            // Handle error
        }
        ```
    - The `Publish` method on `WatermillPubliser` is used to publish a message to a specific topic(queue name). The message is encoded using the Go `encoding/gob` package before being sent.
        ```go
        // Worker struct must be registered before publishing
        err := pub.Publish(topic, workerStruct)
        if err != nil {
            // Handle error
        }
        ```
- #### Dead Letter Queue
    - The `dead letter queue`, also known as the `poison queue` in watermill, is a designated destination for messages that have failed to undergo processing by a consumer.
    - The name of this queue is specified in the `DEAD_LETTER_QUEUE` environment variable, we are storing failed job into database.
    - Command to run dead letter queue
        ```go
        go run app.go dead-letter-queue
        ```

---
### **Code Walk-through**
- #### Config:
    - We are using [envconfig](https://github.com/kelseyhightower/envconfig) which binds env values to struct elements.
    ```go
        type Config struct {
            AppName string `envconfig:"APP_NAME"`
        }
    ```
    - In above example `APP_NAME` is bind with `AppName` element of struct so further when it is used it will have value from `env`.
    - To load env variables from `.env` we are using [gotdotenv](https://github.com/joho/godotenv), which basically initialized when app starts [main.go](app.go).
    - We have also create a seprate method `LoadTestEnv()` in [config](main.go) which helps to load envs from `.env.testing`
- #### Command:
    - We are using [cobra](https://github.com/spf13/cobra) for commands, you need create seprate file for your command.
    - To register your command add in [cli/main.go](cli/main.go)
    - After registering your command you can run that by following syntax `go run app.go {your_command_name}`
- #### Route:
    - If you are building api then you need to add routes in [routes/main.go](routes/main.go).
    - In that we have `Setup()` which will be initialized while starting app so all routes that are required need to register their.
    - Currently we are using seprate method to register routes, that sepration is based on controller. e.g If your controller have multiple or single we will create separte function to register it.
    ```go
        func setupAuthController(v1 fiber.Router, goqu *goqu.Database) error {
	        authController, err := controller.NewAuthController(goqu)
            if err != nil {
                return err
            }
            v1.Post("/login", authController.DoAuth)
            return nil
        }
    ```
    - Above function will be created inside [routes/main.go](routes/main.go), and that function will be called in `Setup()`.
    - Why we do this ? because every model does required database object that will perform database operation. We will be using same connection for all the database operation.
    - `controller.NewAuthController(goqu)` This line at starting of function calls controller method which calls initialize method of model.
    - So when we routes will register we will initialize method of controller, that will initialize model struct with database object.
- #### Middleware:
    - Currently we have two middleware [http_logger](middleware/http_logger.go) & [jwt](middleware/jwt.go)
    - [http_logger](middleware/http_logger.go): This middleware enusres that it logs every incoming request to logger.
    - [jwt](middleware/jwt.go): JWT validates the user authentication, in routes we need to set which routes need authentication, so middleware will check for auth on that request.
- #### Service:
    - Also known as the Business Layer, contains the functionality that compounds the core of the application, thus becoming highly reusable for controllers, workers, jobs, and CLI.
- #### Model:
    - Model refer to database tables, currently we have followed structure that one table have one model.
    - When you create a model it should inside of `models` folder.
    - Create a model struct with initialize method.
    ```go
    type UserModel struct {
	    db *goqu.Database
    }

    func InitUserModel(goqu *goqu.Database) (UserModel, error) {
        return UserModel{
            db: goqu,
        }, nil
    }
    ```
    - Above code snippet represent model struct and its init method that will be call in controller method to initialize database object.
    - When you create new method after initialization that function should be object of `UserModel` in current example. If you didn't do that then function will not able to access to database object.
    - Define struct that will represent your table columns.
    ```go
    type User struct {
        ID        string `json:"id"`
        FirstName string `json:"first_name" db:"first_name" validate:"required"`
        LastName  string `json:"last_name" db:"last_name" validate:"required"`
        Email     string `json:"email" db:"email" validate:"required"`
        Password  string `json:"password,omitempty" db:"password" validate:"required"`
    }
    ```
    - Above example represent fields that table have (we have added limited fields here)
    - `json` - in struct represent key name when it will be returned in json format.
    - `db` - represents name of database column, it's supported by `goqu` package.
    - `validate` - if you want to validate fields that are required you can set `validate:"required"`
    - `omitempty` - this will be used with along side of json, if value will be empty then that element will be trimmed from response of json.
- #### Controller:
    - As descirbed in [route](#route) controller init method should be call in route.
    - Each controller must have their struct which contains model object, that will be use to call function of models.
- #### Utils:
    - We have define some common methods inside `utils` like json response.
    - We need common function which handles json response that's why we have created [json_response.go](utils/json_response.go)
    - Similarly we have created different files for different use in utils.
---
### **Testcases**
- Test file should be available within the same package.
- It is advisible to create one `TestMain()` for one package.
- `TestMain()` should initialize everything you need in testcase.
- After that you can call `m.Run()` which executes all testcases that are available in testcases, Also `m.Run()` return  exit status code.
- After `m.RUN()` we must delete data that was inserted.
- When you run testcase with `make test` it will also display `code coverage` which helps you to determine percentage of code covered in your testcases.

### Generators
We are using `mockery` for generate mock functions for interface. We can generate mock using below command.
```shell
$> go generate ./...
```
