package res

type Vars struct {
	Conf       map[string]string
	TwitchMods []string
	Constants  *Const
}

type Const struct {
	CommandQueueName string
	DbFile           string
	PluginDbFile     string
}

func GetConst() *Const {
	return &Const{
		CommandQueueName: "commands",
	}
}

func (v *Vars) IsTwitchMod(nick string) bool {
	for _, value := range v.TwitchMods {
		if value == nick {
			return true
		}
	}
	return false
}
