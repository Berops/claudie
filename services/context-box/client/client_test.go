package cbox

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/Berops/platform/proto/pb"
	"github.com/Berops/platform/urls"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func ClientConnection() pb.ContextBoxServiceClient {
	cc, err := grpc.Dial(urls.ContextBoxURL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to server: %v", err)
	}

	// Creating the client
	c := pb.NewContextBoxServiceClient(cc)
	return c
}

func TestGetConfigScheduler(t *testing.T) {
	c := ClientConnection()

	res, err := GetConfigScheduler(c)
	require.NoError(t, err)
	t.Log("Config name", res.GetConfig().GetName())
}

func TestGetConfigBuilder(t *testing.T) {
	c := ClientConnection()

	res, err := GetConfigBuilder(c)
	require.NoError(t, err)
	t.Log("Config name", res.GetConfig().GetName())
}

func TestGetAllConfigs(t *testing.T) {
	c := ClientConnection()

	res, err := GetAllConfigs(c)
	require.NoError(t, err)
	for _, c := range res.GetConfigs() {
		t.Log(c.GetId(), c.GetName(), c.GetDesiredState(), c.CurrentState)
	}
}

func TestSaveConfigFrontEnd(t *testing.T) {
	c := ClientConnection()

	manifest, errR := ioutil.ReadFile("./manifest.yaml") //this is manifest from this test file
	if errR != nil {
		log.Fatalln(errR)
	}

	_, err := SaveConfigFrontEnd(c, &pb.SaveConfigRequest{
		Config: &pb.Config{
			Id:       "6126737f4f9bcdabaa336da4",
			Name:     "TestDeleteConfig Samo",
			Manifest: string(manifest),
		},
	})
	require.NoError(t, err)
}

func TestSaveConfigScheduler(t *testing.T) {
	c := ClientConnection()

	manifest, errR := ioutil.ReadFile("./manifest.yaml") //this is manifest from this test file
	if errR != nil {
		log.Fatalln(errR)
	}

	err := SaveConfigScheduler(c, &pb.SaveConfigRequest{
		Config: &pb.Config{
			Name:     "TestDeleteNodeSamo",
			Manifest: string(manifest),
		},
	})
	require.NoError(t, err)
}

func TestDeleteConfig(t *testing.T) {
	c := ClientConnection()
	err := DeleteConfig(c, "6126737f4f9bcdabaa336da4")
	require.NoError(t, err)
}

func TestPrintConfig(t *testing.T) {
	c := ClientConnection()
	_, err := PrintConfig(c, "6126737f4f9bcdabaa336da4")
	if err != nil {
		log.Fatalln("Config not found:", err)
	}
	require.NoError(t, err)
}
