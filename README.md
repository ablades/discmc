# DiscMC

## About
DiscMC is a discord bot that leverages the use of the Discord API and AWS to provide a way to allow me and my friends to start and stop a custom minecraft server hosted on AWS.


## Process
- Discord Bot waits for a command.

- When proper command is given the discord bot will utilize my AWS Credentials to start an EC2 instance on AWS.

- Instance will begin to start and run a SystemCTL service that intializes the mincraft server on instance start.

- After a couple of minutes the server is up and ready to play!


## Technologies Used

* [Golang](https://golang.org/)
* [Discord Go Wrapper](https://github.com/bwmarrin/discordgo)
* [AWS Go SDK](https://aws.amazon.com/sdk-for-go/)
* [Minecraft](https://minecraft.com)
