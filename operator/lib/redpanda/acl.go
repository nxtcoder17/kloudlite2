package redpanda

import (
	"fmt"
	exec2 "k8s.io/utils/exec"
	fn "operators.kloudlite.io/lib/functions"
)

type AdminClient interface {
	CreateUser(username, password string) error
	DeleteUser(username string) error
	UserExists(username string) error
	CreateTopic(topicName string, partitionCount int) error
	DeleteTopic(topicName string) error
	TopicExists(topicName string) (bool, error)
	AllowUserOnTopics(username string, topicNames ...string) error
}

type adminCli struct {
	kafkaBrokers  string
	adminEndpoint string
	saslAuthFlags string
}

func exitCode(err error) int {
	if exitErr, ok := err.(exec2.ExitError); ok {
		return exitErr.ExitStatus()
	}
	return 17
}

func (admin adminCli) UserExists(username string) error {
	err, _, _ := fn.Exec(fmt.Sprintf("rpk acl user list --brokers %s %s | grep -i %s", admin.kafkaBrokers, admin.saslAuthFlags, username))
	if err != nil {
		return err
	}
	return nil
}

func (admin adminCli) TopicExists(topicName string) (bool, error) {
	err, _, _ := fn.Exec(
		fmt.Sprintf(
			"rpk topic describe %s --brokers %s %s", topicName, admin.kafkaBrokers, admin.saslAuthFlags,
		),
	)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (admin adminCli) CreateUser(username, password string) error {
	err, _, _ := fn.Exec(fmt.Sprintf("rpk acl user create %s -p %s --api-urls %s", username, password, admin.adminEndpoint))
	if err != nil {
		return err
	}
	return nil
}

func (admin adminCli) DeleteUser(username string) error {
	err, _, _ := fn.Exec(fmt.Sprintf("rpk acl user delete %s --api-urls %s", username, admin.saslAuthFlags))
	if err != nil {
		return err
	}
	return nil
}

func (admin adminCli) CreateTopic(topicName string, partitionCount int) error {
	err, _, stderr := fn.Exec(
		fmt.Sprintf(
			"rpk topic create %s -p %d --brokers %s %s",
			topicName,
			partitionCount,
			admin.kafkaBrokers,
			admin.saslAuthFlags,
		),
	)
	if err != nil {
		return err
	}
	fmt.Println(stderr.String())
	return nil
}

func (admin adminCli) DeleteTopic(topicName string) error {
	err, _, _ := fn.Exec(fmt.Sprintf("rpk topic delete %s --brokers %s %s", topicName, admin.kafkaBrokers, admin.saslAuthFlags))
	if err != nil {
		return err
	}
	return nil
}

func (admin adminCli) AllowUserOnTopics(username string, topicNames ...string) error {
	topicFlags := ""
	for i := range topicNames {
		topicFlags += " --topic " + topicNames[i]
	}

	err, _, _ := fn.Exec(
		fmt.Sprintf(
			"rpk acl create --allow-principal %s --operation all %s --brokers %s %s", username, topicFlags, admin.kafkaBrokers,
			admin.saslAuthFlags,
		),
	)
	if err != nil {
		return err
	}
	return nil
}

func NewAdminClient(username, password, kafkaBrokers, adminEndpoint string) AdminClient {
	return &adminCli{
		kafkaBrokers:  kafkaBrokers,
		adminEndpoint: adminEndpoint,
		saslAuthFlags: fmt.Sprintf("--user %s --password %s --sasl-mechanism SCRAM-SHA-256", username, password),
	}
}
