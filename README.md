# chatwork-command

*this project is work in progress*

Now newest version is none.

# usage

```
// send my chat
chatwork-command send "test message"
// specific room
chatwork-command send -r room_name "test message"

// show mychat latest message
chatwork-command show
// show specific room latest message
chatwork-command show -r room_name
// show specific room latest 100 message
chatwork-command show -r room_name -n 100

// watch new messages with verbose log every 10 minutes
chatwork-command watch -v 10 
room    username    message
Room1   User1       Hello!
Room1   User2       HRU
Room2   UserX       It looks like rain.  
```
