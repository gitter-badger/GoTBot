package context

import "github.com/asdine/storm"

var Conf map[string]string
var DB *storm.DB
var Users storm.Node

const CommandQueueName = "commands"
const UserBucketName = "Users"
const DbFile = "gotbot.db"