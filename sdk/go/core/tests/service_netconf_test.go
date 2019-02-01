package test

import (
	"fmt"
	"strconv"
	"strings"
	ysanity "github.com/CiscoDevNet/ydk-go/ydk/models/ydktest/sanity"
	ocInterfaces "github.com/CiscoDevNet/ydk-go/ydk/models/ydktest/openconfig_interfaces"
	ocBgp        "github.com/CiscoDevNet/ydk-go/ydk/models/ydktest/openconfig_bgp"
	encoding     "github.com/CiscoDevNet/ydk-go/ydk/types/encoding_format"
	"github.com/CiscoDevNet/ydk-go/ydk"
	"github.com/CiscoDevNet/ydk-go/ydk/providers"
	"github.com/CiscoDevNet/ydk-go/ydk/services"
	"github.com/CiscoDevNet/ydk-go/ydk/types"
	"github.com/CiscoDevNet/ydk-go/ydk/types/datastore"
	"github.com/stretchr/testify/suite"
	"testing"
)

type NetconfServiceTestSuite struct {
	suite.Suite
	Provider 	providers.NetconfServiceProvider
	Crud 		services.CrudService
	NS 			services.NetconfService
}

func (suite *NetconfServiceTestSuite) SetupSuite() {
	suite.Crud = services.CrudService{}
	suite.NS = services.NetconfService{}
	suite.Provider = providers.NetconfServiceProvider{
		Address:  "127.0.0.1",
		Username: "admin",
		Password: "admin",
		Port:     12022}
	suite.Provider.Connect()
}

func (suite *NetconfServiceTestSuite) TearDownSuite() {
	suite.Provider.Disconnect()
}

func (suite *NetconfServiceTestSuite) BeforeTest(suiteName, testName string) {
	suite.Crud.Delete(&suite.Provider, &ysanity.Runner{})
	suite.NS.DiscardChanges(&suite.Provider)
	fmt.Printf("%v: %v ...\n", suiteName, testName)
}

func (suite *NetconfServiceTestSuite) TestEditCommitGet() {
	runner := ysanity.Runner{}
	runner.YdktestSanityOne.Number = 1
	runner.YdktestSanityOne.Name = "runner:one:name"

	// Edit Config
	op := suite.NS.EditConfig(&suite.Provider, datastore.Candidate, &runner, "", "", "")
	suite.Equal(op, true)

	// Get Config
	readEntity := suite.NS.GetConfig(&suite.Provider, datastore.Candidate, &ysanity.Runner{})
	suite.Equal(types.EntityEqual(readEntity, &runner), true)

	// Commit
	op = suite.NS.Commit(&suite.Provider, false, -1, -1, -1)
	suite.Equal(op, true)

	// Get
	readEntity = suite.NS.Get(&suite.Provider, &ysanity.Runner{})
	suite.Equal(types.EntityEqual(readEntity, &runner), true)
}

func (suite *NetconfServiceTestSuite) TestLockUnlock() {
	op := suite.NS.Lock(&suite.Provider, datastore.Running)
	suite.Equal(op, true)

	op = suite.NS.Unlock(&suite.Provider, datastore.Running)
	suite.Equal(op, true)
}

func (suite *NetconfServiceTestSuite) TestLockUnlockFail() {
	op := suite.NS.Lock(&suite.Provider, datastore.Candidate)
	suite.Equal(op, true)

	funcDidPanic, panicValue := didPanic(func() { suite.NS.Unlock(&suite.Provider, datastore.Running) })
	suite.Equal(funcDidPanic, true)
	suite.Regexp("YServiceProviderError:", panicValue)
	errMsg := `<rpc-error>
    <error-type>application</error-type>
    <error-tag>operation-failed</error-tag>
    <error-severity>error</error-severity>
  </rpc-error>`
	suite.Regexp(errMsg, panicValue)
}

func (suite *NetconfServiceTestSuite) TestValidate() {
	op := suite.NS.Validate(&suite.Provider, datastore.Candidate, nil, "")
	suite.Equal(op, true)

	runner := ysanity.Runner{}
	runner.YdktestSanityOne.Number = 1
	runner.YdktestSanityOne.Name = "runner:one:name"
	op = suite.NS.Validate(&suite.Provider, datastore.NotSet, &runner, "")
	suite.Equal(op, true)
}

func (suite *NetconfServiceTestSuite) TestDiscardChanges(){
	runner := ysanity.Runner{}
	runner.Two.Number = 2
	runner.Two.Name = "runner:two:name"

	// EditConfig
	op := suite.NS.EditConfig(&suite.Provider, datastore.Candidate, &runner, "", "", "")
	suite.Equal(op, true)

	// DiscardChanges
	op = suite.NS.DiscardChanges(&suite.Provider)
	suite.Equal(op, true)

	// GetConfig
	readEntity := suite.NS.GetConfig(&suite.Provider, datastore.Candidate, &ysanity.Runner{})
	suite.Nil(readEntity)
}

func (suite *NetconfServiceTestSuite) TestConfirmedCommit() {
	runner := ysanity.Runner{}
	runner.Two.Number = 2
	runner.Two.Name = "runner:two:name"

	op := suite.NS.EditConfig(&suite.Provider, datastore.Candidate, &runner, "", "", "")
	suite.Equal(op, true)

	op = suite.NS.Commit(&suite.Provider, true, 120, -1, -1)
	suite.Equal(op, true)

	readEntity := suite.NS.Get(&suite.Provider, &ysanity.Runner{})
	suite.Equal(types.EntityEqual(readEntity, &runner), true)
}

// skip
func (suite *NetconfServiceTestSuite) TestCancelCommit() {
	suite.T().Skip("No message id in cancel commit payload")

	runner := ysanity.Runner{}
	runner.Two.Number = 2
	runner.Two.Name = "runner:two:name"

	op := suite.NS.EditConfig(&suite.Provider, datastore.Candidate, &runner, "", "", "")
	suite.Equal(op, true)

	op = suite.NS.Commit(&suite.Provider, true, 120, -1, -1)
	suite.Equal(op, true)

	op = suite.NS.CancelCommit(&suite.Provider, -1)
	suite.Equal(op, true)

	readEntity := suite.NS.Get(&suite.Provider, &ysanity.Runner{})
	suite.Equal(types.EntityEqual(readEntity, &runner), true)
}

// TODO: Problems found with path functions:
//			readEntity changes with getFilter
//			EntityEqual(readEntity, &runner) returns true, even after the above
func (suite *NetconfServiceTestSuite) TestCopyConfig() {
	runner := ysanity.Runner{}
	runner.Two.Number = 2
	runner.Two.Name = "runner:two:name"

	var op bool
	var readEntity types.Entity

	// Modify Candidate via CopyConfig from runner
	op = suite.NS.CopyConfig(
		&suite.Provider, datastore.Candidate, datastore.NotSet, &runner, "")
	suite.Equal(op, true)

	readEntity = suite.NS.GetConfig(&suite.Provider, datastore.Candidate, &ysanity.Runner{})
	suite.Equal(types.EntityEqual(readEntity, &runner), true)

	// Modify Candidate via CopyConfig from runner
	runner.Two.Name = fmt.Sprintf("%s_modified", runner.Two.Name)

	op = suite.NS.CopyConfig(
		&suite.Provider, datastore.Candidate, datastore.NotSet, &runner, "")
	suite.Equal(op, true)

	readEntity = suite.NS.GetConfig(&suite.Provider, datastore.Candidate, &ysanity.Runner{})
	suite.Equal(types.EntityEqual(readEntity, &runner), true)

	// Modify Candidate via CopyConfig from Running. That will remove config.
	op = suite.NS.CopyConfig(&suite.Provider, datastore.Candidate, datastore.Running, nil, "")
	suite.Equal(op, true)

	readEntity = suite.NS.GetConfig(&suite.Provider, datastore.Candidate, &ysanity.Runner{})
	suite.Nil(readEntity)

	// DiscardChanges
	op = suite.NS.DiscardChanges(&suite.Provider)
	suite.Equal(op, true)
}

// skip
func (suite *NetconfServiceTestSuite) TestDeleteConfig() {
	suite.T().Skip("startup not enabled in ConfD")

	runner := ysanity.Runner{}
	runner.Two.Number = 2
	runner.Two.Name = "runner:two:name"

	op := suite.NS.CopyConfig(
		&suite.Provider, datastore.Startup, datastore.NotSet, &runner, "")
	suite.Equal(op, true)

	op = suite.NS.DeleteConfig(&suite.Provider, datastore.Startup, "")
	suite.Equal(op, true)
}

func (suite *NetconfServiceTestSuite) TestCloseSession() {
	op := suite.NS.CloseSession(&suite.Provider)
	suite.Equal(op, true)

	funcDidPanic, panicValue := didPanic(func() { suite.NS.Lock(&suite.Provider, datastore.Running) })
	suite.Equal(funcDidPanic, true)
	suite.Equal(panicValue, "YClientError: Could not send payload")

	suite.Provider.Connect()

	op = suite.NS.Lock(&suite.Provider, datastore.Running)
	op = suite.NS.Unlock(&suite.Provider, datastore.Running)
	suite.Equal(op, true)
}

// skip
func (suite *NetconfServiceTestSuite) TestKillSession() {
	suite.T().Skip("session-id not recognized")

	op := suite.NS.Lock(&suite.Provider, datastore.Candidate)
	suite.Equal(op, true)

	funcDidPanic, panicValue := didPanic(func() { suite.NS.Lock(&suite.Provider, datastore.Candidate) })
	suite.Equal(funcDidPanic, true)
	suite.Regexp("<session-id>", panicValue)

	sessionIDStr := strings.Split(panicValue, "<session-id>")[1]
	sessionIDStr = strings.Split(sessionIDStr, "</session-id>")[0]
	sessionID, err := strconv.Atoi(sessionIDStr)
	suite.Equal(err, nil)

	op = suite.NS.KillSession(&suite.Provider, sessionID)
	suite.Equal(op, true)
}

func (suite *NetconfServiceTestSuite) TestGetOCEntities() {

    // Build filter
    interfacesFilter := ocInterfaces.Interfaces{};
    bgpFilter := ocBgp.Bgp{};
    filterList := types.NewFilter(&interfacesFilter, &bgpFilter)

    // Read running config
    getConfigEntity := suite.NS.Get(&suite.Provider, filterList);
    suite.Equal( types.IsEntityCollection(getConfigEntity), true)

    // Get results
    getConfigEC := types.EntityToCollection(getConfigEntity)
    for _, entity := range getConfigEC.Entities() {
    	ydk.YLogDebug(fmt.Sprintf("Printing %s", GetEntityXMLString(entity)))
    }
}

func (suite *NetconfServiceTestSuite) TestGetEditCopyConfigSanity() {
	// Build configuration
	runner := ysanity.Runner{}
	runner.Two.Number = 2
	runner.Two.Name = "runner-two-name"
	
	native := ysanity.Native{}
	native.Version = "0.1.0"
	native.Hostname = "MyHost"
	
	configEC := types.NewConfig(&runner, &native)

	result := suite.NS.EditConfig(&suite.Provider, datastore.Candidate, configEC, "", "", "")
	suite.Equal(result, true)
	
    // Build filter
	runnerFilter := ysanity.Runner{}
	nativeFilter := ysanity.Native{}
    filterEC := types.NewFilter(&runnerFilter, &nativeFilter)

    // Read running config
    getConfigEntity := suite.NS.GetConfig(&suite.Provider, datastore.Candidate, filterEC);
    suite.Equal( types.IsEntityCollection(getConfigEntity), true)

    // Get results
    getConfigEC := types.EntityToCollection(getConfigEntity)
    for _, entity := range getConfigEC.Entities() {
    	ydk.YLogDebug(fmt.Sprintf("Printing %s", GetEntityXMLString(entity)))
    }

    // Discard changes
    result = suite.NS.DiscardChanges(&suite.Provider);
    suite.Equal(result, true)
}

func (suite *NetconfServiceTestSuite) TestSanityGetRunningConfig() {
	// Create empty EntityCollection
    filterEC := types.NewFilter()

    // Read running-config
    getEntity := suite.NS.GetConfig(&suite.Provider, datastore.Running, filterEC);
    suite.Equal( true, types.IsEntityCollection(getEntity))

    // Get results
    getEC := types.EntityToCollection(getEntity)
    for _, entity := range getEC.Entities() {
    	ydk.YLogDebug(fmt.Sprintf("Printing %s", GetEntityXMLString(entity)))
    }
}

func TestNetconfServiceTestSuite(t *testing.T) {
	suite.Run(t, new(NetconfServiceTestSuite))
}

func GetEntityXMLString(entity types.Entity) string {
    codec := services.CodecService{}
	provider := providers.CodecServiceProvider{}
	provider.Encoding = encoding.XML
	payload := codec.Encode(&provider, entity)
	out := fmt.Sprintf("Entity{%s}\n%s", types.EntityToString(entity), payload)
	return out
}
