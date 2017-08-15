package handlers

func (d *deps) Echo() error {
	d.connection.Privmsg(d.channel, "Thanks's for sending goSay with '"+d.params+"' on "+d.channel+", "+d.sender+"!")
	return nil
}
