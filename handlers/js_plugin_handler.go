package handlers

import (
	"github.com/3stadt/GoTBot/structs"
	"os"
	"io/ioutil"
	"github.com/robertkrimen/otto"
	"github.com/thoj/go-ircevent"
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
	_, _ = vm.Run(string(jsData))
	return nil, nil
}
