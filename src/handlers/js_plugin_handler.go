package handlers

import (
	"github.com/3stadt/GoTBot/src/structs"
	"os"
	"io/ioutil"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
	"github.com/thoj/go-ircevent"
	"encoding/json"
	"fmt"
	"github.com/3stadt/GoTBot/src/db"
)

func JsPluginHandler(filePath string, channel string, sender string, params string, connection *irc.Connection) (*structs.Message, error) {
	var err error
	var jsData []byte
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}
	if jsData, err = ioutil.ReadFile(filePath); err != nil {
		return nil, err
	}
	vm := otto.New()
	vm.Set("channel", channel)
	vm.Set("sender", sender)
	vm.Set("params", params)
	vm.Set("sendMessage", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) > 0 {
			msg := call.Argument(0)
			connection.Privmsg(channel, msg.String())
		}
		return otto.Value{}
	})

	vm.Set("getUser", func(call otto.FunctionCall) otto.Value {
		result, _ := vm.ToValue("")
		if len(call.ArgumentList) < 1 {
			return result
		}
		username, err := call.Argument(0).ToString()
		if err != nil {
			return result
		}
		result, _ = vm.ToValue(*getBoltUserAsJson(username))
		return result
	})

	_, err = vm.Run(string(jsData))
	if err != nil {
		fmt.Println("ERROR in javascript file " + filePath + ":")
		fmt.Println(err)
	}
	return nil, nil
}

func getBoltUserAsJson(username string) *string {
	emptyJson := "{}"
	userStruct, err := db.GetUser(username)
	if err != nil {
		return &emptyJson
	}
	jUser, err := json.Marshal(*userStruct)
	if err != nil {
		return &emptyJson
	}
	userdata := string(jUser)
	return &userdata
}
