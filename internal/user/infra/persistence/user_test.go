//go:build integration

package persistence

import (
	"context"
	"fmt"
	"os"
	"testing"

	"starter-boilerplate/internal/user/domain/model"
	"starter-boilerplate/internal/user/domain/repository"
	"starter-boilerplate/pkg/testcontainer"

	"github.com/stretchr/testify/suite"
)

type UserRepoSuite struct {
	suite.Suite
	pg   *testcontainer.PgContainer
	repo repository.UserRepository
}

func TestUserRepository(t *testing.T) {
	if err := os.Chdir("../../../.."); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	pg, err := testcontainer.SetupPgContainer(context.Background(), &testcontainer.PgContainer{
		Database: "testdb",
		Username: "testuser",
		Password: "testpass",
		HostPort: "25432",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "setup pg container: %v\n", err)
		os.Exit(1)
	}

	s := &UserRepoSuite{pg: pg, repo: NewUserRepository(pg.DB())}
	suite.Run(t, s)
}

func (s *UserRepoSuite) TearDownSuite() {
	s.pg.Close()
	s.pg.Terminate(context.Background())
}

func (s *UserRepoSuite) SetupTest() {
	s.Require().NoError(s.pg.Clean(context.Background()))
}

func newTestUser(id, email string) *model.User {
	return &model.User{
		ID:           id,
		Email:        email,
		PasswordHash: "hashed_password",
		Role:         model.RoleUser,
	}
}

func (s *UserRepoSuite) TestCreate_Success() {
	u := newTestUser("id-1", "alice@example.com")
	err := s.repo.Create(context.Background(), u)
	s.Assert().NoError(err)
}

func (s *UserRepoSuite) TestCreate_DuplicateEmail() {
	ctx := context.Background()

	u1 := newTestUser("id-1", "dup@example.com")
	s.Require().NoError(s.repo.Create(ctx, u1))

	u2 := newTestUser("id-2", "dup@example.com")
	err := s.repo.Create(ctx, u2)
	s.Assert().Error(err)
}

func (s *UserRepoSuite) TestFindByID_Success() {
	ctx := context.Background()

	u := newTestUser("id-1", "bob@example.com")
	s.Require().NoError(s.repo.Create(ctx, u))

	found, err := s.repo.FindByID(ctx, "id-1")
	s.Require().NoError(err)
	s.Require().NotNil(found)
	s.Assert().Equal(u.ID, found.ID)
	s.Assert().Equal(u.Email, found.Email)
	s.Assert().Equal(u.PasswordHash, found.PasswordHash)
	s.Assert().Equal(u.Role, found.Role)
}

func (s *UserRepoSuite) TestFindByID_NotFound() {
	found, err := s.repo.FindByID(context.Background(), "nonexistent")
	s.Assert().NoError(err)
	s.Assert().Nil(found)
}

func (s *UserRepoSuite) TestFindByEmail_Success() {
	ctx := context.Background()

	u := newTestUser("id-1", "carol@example.com")
	s.Require().NoError(s.repo.Create(ctx, u))

	found, err := s.repo.FindByEmail(ctx, "carol@example.com")
	s.Require().NoError(err)
	s.Require().NotNil(found)
	s.Assert().Equal(u.ID, found.ID)
	s.Assert().Equal(u.Email, found.Email)
	s.Assert().Equal(u.PasswordHash, found.PasswordHash)
	s.Assert().Equal(u.Role, found.Role)
}

func (s *UserRepoSuite) TestFindByEmail_NotFound() {
	found, err := s.repo.FindByEmail(context.Background(), "nobody@example.com")
	s.Assert().NoError(err)
	s.Assert().Nil(found)
}

func (s *UserRepoSuite) TestUpdate_Success() {
	ctx := context.Background()

	u := newTestUser("id-1", "dave@example.com")
	s.Require().NoError(s.repo.Create(ctx, u))

	u.Role = model.RoleAdmin
	s.Require().NoError(s.repo.Update(ctx, u))

	found, err := s.repo.FindByID(ctx, "id-1")
	s.Require().NoError(err)
	s.Require().NotNil(found)
	s.Assert().Equal(model.RoleAdmin, found.Role)
}
