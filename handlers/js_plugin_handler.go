package handlers

import (
	"github.com/3stadt/GoTBot/structs"
	"os"
	"io/ioutil"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
	"github.com/thoj/go-ircevent"
	"github.com/3stadt/GoTBot/bolt"
	"encoding/json"
	"fmt"
"github.com/mitchellh/mapstructure"
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

	vm.Set("getUser", func(username string) string {
		return *getBoltUserAsJson(username)
	})

	vm.Set("updateUser", func(userdata interface{}) bool {
		if err != nil {
			fmt.Println(err)
			return false
		}
		userdataMap := userdata.(map[string]interface{})
		updateBoltUserFromJson(userdataMap)
		return true
	})
	_, _ = vm.Run(string(jsData))
	return nil, nil
}

func updateBoltUserFromJson(userdata map[string]interface{}) error {
	user := structs.User{}
	err := mapstructure.Decode(userdata, &user)
	if err != nil {
		fmt.Println(userdata)
		fmt.Println(err)
		return err
	}
	fmt.Println(user)
	return nil
}

func getBoltUserAsJson(username string) *string {
	userStruct := bolt.GetUser(username)
	jUser, err := json.Marshal(*userStruct)
	if err != nil {
		emptyJson := "{}"
		return &emptyJson
	}
	userdata := string(jUser)
	return &userdata
}
