# Quizz App

Description of the quizz app.

## Table of Contents

- [Installation](#installation)
    - [Prerequisite](#prerequisite)
    - [Run the Project](#run-the-project)
- [Commands](#commands)
    1. [Create](#create)
        1. [Admin](#admin)

## Installation

Follow these steps to install and run the project:

### Prerequisite

Ensure that Docker is installed on your system.

### Run the Project

To run this project, execute the following command from the root of the project:


```
    docker-compose up
```

## Commands

- Commands are quick way to do some occasional or manual task.
- You can get more details from [/api/README.md](https://github.com/Improwised/quizz-app/blob/develop/api/README.md) file 
- Use `make up` command to run database migration before running this command (no need to do this time)

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