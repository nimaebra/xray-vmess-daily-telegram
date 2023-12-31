package main

import (
	"fmt"
	"xray-telegram/service/builder"
	"xray-telegram/service/execute"
	"xray-telegram/service/subscribe"
	"xray-telegram/service/telegram"
)

func main() {

	executeInstance := execute.NewExecute()
	executeInstance.ExecuteCommand("./reinstall.sh")

	fmt.Println("read setting file...")

	builderInstance := builder.NewBuilder().
		SetServerIP().
		SetSettingsFile().
		SetPublicKeyAndPrivateKey().
		SetConfigurations().
		SetBlock().
		Save()

	sub := subscribe.NewSubscribe(builderInstance, executeInstance)

	tel := telegram.NewTelegramClient(builderInstance, sub)

	tel.SendVNstat()

}
