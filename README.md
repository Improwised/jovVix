# Quizz App

Description of the quizz app.

## Table of Contents

- [Quick installation and run](#installation)
    - [Prerequisite](#prerequisite)
    - [Run the Project](#run-the-project)
- [Setup development environment](#setup-development-environment)
- [Commands](#commands)
    1. [Create](#create)
        - [Admin](#admin)

## Quick installation and run

Follow these steps to install and run the project:

### Prerequisite

Ensure that Docker is installed on your system.

### Run the Project

To run this project, execute the following command from the root of the project:


```
    docker-compose up
```

## Setup development environment

1. Ensure that you have [Go](https://go.dev/) and [Node.js](https://nodejs.org/en) installed on your system. Add Go and Node.js to your system's PATH for global access.

2. A database is required; the current code is written to work with [PostgreSQL](https://www.postgresql.org/).

3. Run the following commands to install necessary packages. (Assuming you are in the root of the project)

~~~
cd ./api && go get 
cd ./app && npm install 
~~~

4. Change .env.local to .env and update necessary configurations. '.env' file is already in .gitignore. If you need to add some new key-value pairs, you need to update .env.example as well.

5. For hot-reload functionality in Golang, we will use the [air](https://github.com/cosmtrek/air) package.
   - Install it as mentioned in the repository and add the air library to your system's PATH for global access. Add the path of the bin ($GOPATH/bin or $HOME/go/bin) folder if it doesn't exist.
   - Run the following command to start the server and enable hot reload:

~~~
cd api && air
~~~

6. We use Nuxt3 in our project, its development server provides hot reload by default.
    - Run the following command to start nuxt3 server

~~~
cd app && npm run dev
~~~

6. [Hooks](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks) are triggers that are executed based on specific events. We will use these hooks in our project for various checks like spell-checking, linting, and testing.
   - We will use the [Python virtual environment](https://docs.python.org/3/library/venv.html) to create a virtual environment.

~~~
python -m venv venv
source ./venv/bin/activate
~~~
   > Note: You can ignore this step if you want to install it directly on your system. If you create a virtual environment, ensure to add it to your .gitignore file..

   - We will use the Python library [pre-commit](https://pre-commit.com/) to handle these hooks. Hooks are associated with Git, and the pre-commit hook is activated when you have already staged your changes with `git add ...` and are about to commit. You can run it manually by typing `pre-commit` into the terminal.

~~~
source ./venv/bin/activate
pip install pre-commit

git add README.md
pre-commit
~~~

   - Our Node project is configured with Node.js and [husky](https://typicode.github.io/husky/get-started.html). In that case, all configurations are written within the package.json file itself..


## Commands

- Commands are quick way to do some occasional or manual task.
- You can get more details from [/api/README.md](https://github.com/Improwised/quizz-app/blob/develop/api/README.md) file 
- commands are made with [cobra library](https://pkg.go.dev/github.com/spf13/cobra) which written in go.
- change .env.example to .env if you haven't already.
- to run command, first go to **/api** folder (which have go.mod file, so you can run `go run file.go`).
- Use `make up` command to run database migration before running this command (once when you write new migration, make sure database is running).

### Create

- This command is used for the creation of a sub-command.
- It's made to reduce UI integration in some cases and directly create entities.


**command**: `go run . create [subcommand]`
```
# examples
go run . create admin [sub-command]

# sub-commands
- This command needs one subcommand
```

Here is the list of sub-commands which supported by **create** command

#### Admin
- Create admin user

**command**: `go run . create admin <username> <email> <first-name> <last-name> <-f>`
explain: 
```
# examples
go run . create admin adminxyz123 adminxyz@gmail.com admin xyz -f

# arguments
- username: required, string, unique* (see #flags section)
- email: required, string, unique
- first-name: required, string
- last-name: required, string

# flags
- force: bool, default-true
    - It will make your given username unique by modifying last some characters. 
```