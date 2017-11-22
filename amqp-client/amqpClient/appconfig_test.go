package amqpClient

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

//TestNewAppConfiguration testing if the application is able to comprehend and decode the configuration and can create a
//usable config set
func TestNewAppConfiguration(t *testing.T) {
	// <setup code>
	file := newAppConfigSetup()
	defer newAppConfigTearDown(file)

	c, err := newAppConfiguration(file.Name())
	if err != nil {
		t.Fail()
	}

	if c.RabbitMq.Connection.Host != "localhost" {
		t.Fail()
	}

	if c.RabbitMq.Logs.Logfile != "./logs.txt" {
		t.Fail()
	}
}

func newAppConfigSetup() *os.File {
	// <this is just a test setup func>
	//copying dummy/ template content of the config.yml file
	content, err := ioutil.ReadFile("./examples/config.yml")
	tmpfile, err := ioutil.TempFile(".", "config.yml")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	return tmpfile
}

func newAppConfigTearDown(configFile *os.File) {
	// this is a test tear down function>
	os.Remove(configFile.Name()) // clean up
}
