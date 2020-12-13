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

	// Set the file name of the configurations file
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Error reading config file, %s", err))
	}
}

func createSession() *session.Session {
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-2"),
		Credentials: credentials.NewStaticCredentials(viper.GetString("aws.accessid"), viper.GetString("aws.accesskey"), ""),
	})
	if err != nil {
		log.Error(err)
	}

	return session
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + viper.GetString("bot.token"))
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
	// If the message is "ping" reply with "Pong!"
	if m.Content == "start mc" {
		startInstance(viper.GetString("aws.instanceid"))
		s.ChannelMessageSend(m.ChannelID, "Starting MC Server! IP is: 3.22.45.58")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "stop mc" {
		stopInstance(viper.GetString("aws.instanceid"))
		s.ChannelMessageSend(m.ChannelID, "Stopping MC Server!")
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

// Access Key ID:
// AKIAJL4VUZZM7HZTM35A
// Secret Access Key:
// 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
