# Agenda

Agenda is a command line meeting manager.

## Collaborators

15331344 薛明淇
`TODO`

## Usage

```shell
$ go get -u github.com/VinaLx/service-computing-homework/agenda

$ $GOPATH/bin/agenda -h
Agenda is a meeting manager based on CLI using cobra library.
It supports different operation on meetings including register, create meeting, query and so on.
It's a cooperation homework assignment for service computing.

Usage:
  Agenda [command]

Available Commands:
  cancel              Cancel your own meeting by specifying title name.
  changeParticipators Change your own meetings' participators.
  clear               Clear all meetings you attended or created.
  createMeetings      Create meetings.
  delete              A brief description of your command
  help                Help about any command
  list
  listMeetingsCmd     List all of your own meetings during a time interval.
  login               Login
  logout              Logout
  quit                Quit meetings.
  register            Register user.

Flags:
  -d, --debug   display log message
  -h, --help    help for Agenda

Use "Agenda [command] --help" for more information about a command.
```

## Examples

`TODO`

## Implementation Details

The program this time is solely business logic and (almost) no efficiency or programming tricks involved, so language feature and trick doesn't really matters here. The difficulties come down to the logic and hierarchy design of the whole program.

First we provide the overall program structure and then point out some crucial details.

### Overall Structure

On a high level abstraction, the whole program can be separated into four part.

- the entity logic that abstract out the overall invariant of program (package `entity`)
- the program model that implements the basic operations of program (package `model`)
- the logic that pack all basic operations to satisfy all business demand (package `cmd`)
- the (CLI) user interface parsing command line arguments to the action and parameters (package `cli`)

There should be a seperate part for validation for user input parameters and other error handling, but in this case it's combined with the third part due to the time limit of homework.

### Command Line Interface Using Cobra

TODO

### Entity Logic

We hadle the constraints of meetings and users directly by enforcing the instances of `Users` and `Meetings` to be logically correct. That is, no duplicate users can ever be added into `Users`, and no two meetings whose participants would have overlapped busy time.

So that there's a strong limit on where contraint problem can happen. Also we won't be ever worried about whether users and meetings instances are in a valid state in other contexts, because they always are. This design greatly simplifies the logic implementation other parts program.

### Filesystem Management

Obviously there's a requirement of data persistence of meetings and user record, along with the login state. The management method is not complicated at all, but it still worth a mention.

Just at the program start, agenda checks whether its "home directory"(currently `~/.agenda`) exist, if not, it creates one. And program reads all user and meeting infomation on demand (listing all user do not care any meetings at all and therefore meetings are not loaded in that command). And if anything change during the operation, the new user and meeting infomation rewrites the original ones before program exits.

The serialization and deserialization, as request, are implement directly using the `json` standard library. Nothing strange here. The only thing is that we didn't handle the inconsistency caused by the corruption of storage file, if any unexpected inconsistencies are detected, the program should abort directly with error messages printed.

### Login State

We stored user password in the form of md5 hash, and login file contains the plain username and password. When checking the login state, the password would be checked again in case of the corruption of login file. If the password and hash don't match, the login file is deleted and some operation requiring login state is denied.

### Business Logic

After the implementation of other parts, implementing business logic is straightforward, just doing basic error handling and logic by invoking underlying program model.

## Last

I did have some fun here :P