package suite

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"starter-boilerplate/internal"
	"starter-boilerplate/internal/shared/config"
	sharedjwt "starter-boilerplate/internal/shared/jwt"
	pkgjwt "starter-boilerplate/pkg/jwt"
	"starter-boilerplate/pkg/testcontainer"

	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

// FunctionalSuite is a reusable test suite that boots the full application
// with real Postgres and Redis containers. Embed it in a local test suite
// to get lifecycle management and HTTP helpers for free.
type FunctionalSuite struct {
	suite.Suite
	CM         *testcontainer.ContainerManager
	JWTManager *pkgjwt.Manager
	BaseURL    string
	cancel     context.CancelFunc

	// Configuration — set before suite.Run; defaults applied in SetupSuite.
	PgDatabase   string // default: "testdb"
	PgUsername   string // default: "testuser"
	PgPassword   string // default: "testpass"
	PgPort       string // default: "15432"
	RedisPort    string // default: "16379"
	AMQPPort     string // default: "15672"
	TestPassword string // default: "P@ssw0rd123"
	FixtureDir   string // required, e.g. "tests/functional/testdata/fixtures"
}

func defaultVal(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func (s *FunctionalSuite) SetupSuite() {
	s.PgDatabase = defaultVal(s.PgDatabase, "testdb")
	s.PgUsername = defaultVal(s.PgUsername, "testuser")
	s.PgPassword = defaultVal(s.PgPassword, "testpass")
	s.PgPort = defaultVal(s.PgPort, "15432")
	s.RedisPort = defaultVal(s.RedisPort, "16379")
	s.AMQPPort = defaultVal(s.AMQPPort, "15673")
	s.TestPassword = defaultVal(s.TestPassword, "P@ssw0rd123")

	s.Require().NotEmpty(s.FixtureDir, "FixtureDir must be set before running the suite")

	ctx := context.Background()

	hash, err := bcrypt.GenerateFromPassword([]byte(s.TestPassword), bcrypt.MinCost)
	s.Require().NoError(err, "generate bcrypt hash")

	s.CM, err = testcontainer.Setup(ctx,
		&testcontainer.PgContainer{
			Database: s.PgDatabase,
			Username: s.PgUsername,
			Password: s.PgPassword,
			HostPort: s.PgPort,
			Fixtures: map[string]testcontainer.FixtureSet{
				"default": {
					Dir:          s.FixtureDir,
					TemplateData: map[string]interface{}{"PasswordHash": string(hash)},
				},
			},
		},
		&testcontainer.RedisContainer{
			HostPort: s.RedisPort,
		},
		&testcontainer.AMQPContainer{
			HostPort: s.AMQPPort,
		},
	)
	s.Require().NoError(err, "setup containers")

	os.Setenv("APP_ENV", "test")
	cfg := config.SetupConfig()
	s.JWTManager = sharedjwt.NewJWTManager(cfg.JWT)

	appCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	application := internal.InitializeApp(appCtx)
	s.BaseURL = application.BaseURL()
	go func() { _ = application.Run(appCtx) }()

	select {
	case <-application.Ready():
	case err := <-application.StartErr():
		s.Require().NoError(err, "server start")
	case <-time.After(10 * time.Second):
		s.Require().Fail("server timeout")
	}
}

func (s *FunctionalSuite) TearDownSuite() {
	s.cancel()
	s.CM.Close()
	s.CM.Terminate(context.Background())
}

func (s *FunctionalSuite) SetupTest() {
	s.CM.LoadFixtures(s.T(), "default")
}

// --- Exported helpers ---

func (s *FunctionalSuite) DoRequest(method, path, body string, headers map[string]string) *http.Response {
	s.T().Helper()
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, s.BaseURL+path, bodyReader)
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	return resp
}

func (s *FunctionalSuite) DoAuthRequest(method, path, token, body string) *http.Response {
	s.T().Helper()
	return s.DoRequest(method, path, body, map[string]string{
		"Authorization": "Bearer " + token,
	})
}

func (s *FunctionalSuite) ReadJSON(resp *http.Response, target any) {
	s.T().Helper()
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)
	err = json.Unmarshal(data, target)
	s.Require().NoError(err, "failed to unmarshal: %s", string(data))
}

func (s *FunctionalSuite) IssueAccessToken(userID, role string) string {
	s.T().Helper()
	token, err := s.JWTManager.GenerateAccessToken(userID, role)
	s.Require().NoError(err)
	return token
}

func (s *FunctionalSuite) IssueRefreshToken(userID, role string) string {
	s.T().Helper()
	token, err := s.JWTManager.GenerateRefreshToken(userID, role)
	s.Require().NoError(err)
	return token
}
