# Agenda

Agenda is a command line meeting manager.

## Collaborators

15331344 薛明淇
15331348 颜泽鑫
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

#### Global argument

```
-h, --help
-d, --debug
```
You can get helpful information whenever you are using Cobar. And you can get detailed log information when using -d.

#### Main Command
We implement the following command for you to use.

```
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
```

#### register
You can register an account by using `register` command, and we define the following requisite argument. All of them are type of `string`.

```
Flags:
  -m, --mail string       email.
  -p, --password string   Help message for username
  -t, --phone string      Phone
  -u, --user string       Username
```

#### Login
You need to login before you use the most of command in `Agenda`. Once you login, you don't have to login next time, we will keep your state. But you have to log out by yourself in order to protect your privacy. 

```
Flags:
  -p, --password string   Input password
  -u, --user string       Input username
```

#### Logout
No argument is needed here. When you have nothing else to do, you had better log out in order to protect your privacy.

#### delete
We are sorry when you use `delete`. Once you use this command, all of the information about you will be erased, and there is no way for us and you to get it back. So we recommend you not to use it.

#### List
No argument is needed here. You can get all users' information except password so that you can invite others to join your meetings.


#### Create Meetings
After you have login, you can create a meeting by using `createMeetings` command. All of the arguments are required. You should invite at least one participator to join your meeting. After all, you can't have a meeting only with yourself. :P

```
Flags:
  -e, --end string             Input end time as the format of (yyyy-mm-dd).
  -p, --participators string   Input participator name.
  -s, --start string           Input start time as the format of (yyyy-mm-dd).
  -t, --title string           Input title name.
```

#### List Meetings
You can using `listMeetings` to get all of the meetings' information during the specific interval(start - end).

```
Flags:
  -e, --end string     Input the end time.(yyyy-mm-dd)
  -s, --start string   Input the start time.(yyyy-mm-dd)
```

#### Change Participators
Whenever you need to invite or remove the participators of your meetings, you can using `changeParticipators`. But you have to ensure that he or she are available to attend your meeting.

```
Flags:
  -y, --delete                 If true, delete participators, otherwise append participators.
  -p, --participators string   Input the participators.
  -t, --title string           Input the title name.
```

#### Cancel
You may make some mistakes when creating a meeting, then you can just cancel it.

```
Flags:
  -t, --title string   Input the title.
```

#### Quit
You can refuse to attend a meeting by using `quit`, then you can be free to go out for a play.

```
Flags:
  -t, --title string   Input the title.
```

#### clear Meetings
You may feel busy and boring for attending meetings. At that moment, you can use `clear` to cancel all of the meeting which you attended and quit all of the meetings which you are invited. That's you are free!!!


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

