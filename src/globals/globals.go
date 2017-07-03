package globals

import "github.com/thoj/go-ircevent"

var Conf map[string]string
var Connection *irc.Connection
const CommandQueueName = "commands"
const UserbucketName = "Users"