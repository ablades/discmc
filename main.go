package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/bwmarrin/discordgo"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Error reading config file, %s", err))
	}
}

func createSession() *session.Session {
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-2"),
		Credentials: credentials.NewStaticCredentials(viper.GetString("ACCESS_ID"), viper.GetString("ACCESS_KEY"), ""),
	})
	if err != nil {
		log.Error(err)
	}

	return session
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + viper.GetString("BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	
	guildID := m.GuildID
	member, err := s.GuildMember(guildID, m.Author.ID)
	if err != nil {
		panic(err)
	}

	if len(member.Roles) == 0 || member.Roles[0] != "783617119325126666"{
		s.ChannelMessageSend(m.ChannelID, "You're not a BIG BOI cant play mc with us")
		return
	}
	

	// If the message is "ping" reply with "Pong!"
	if m.Content == "start mc" {
		startInstance(viper.GetString("INSTANCE_ID"))
		s.ChannelMessageSend(m.ChannelID, "Starting MC Server! IP is: 3.22.45.58")
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "stop mc" {
		stopInstance(viper.GetString("INSTANCE_ID"))
		s.ChannelMessageSend(m.ChannelID, "Stopping MC Server!")
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

func startInstance(instance string) {
	svc := ec2.New(createSession())

	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(instance),
		},
	}

	result, err := svc.StartInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func stopInstance(instance string) {
	svc := ec2.New(createSession())

	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(instance),
		},
	}

	result, err := svc.StopInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}
