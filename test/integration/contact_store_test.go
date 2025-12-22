package integration

import (
	"context"
	"log"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/webitel/im-contact-service/internal/model"
	"github.com/webitel/im-contact-service/internal/service/dto"
	"github.com/webitel/im-contact-service/internal/store"
	"github.com/webitel/im-contact-service/internal/store/postgres"
	testhelpers "github.com/webitel/im-contact-service/test/integration/test_helpers"
)

type ContactStoreTestSuite struct {
	suite.Suite
	postgresContainer *testhelpers.PostgresContainer
	repo              store.ContactStore
	ctx               context.Context
}

func newContact(domain int, opts ...func(*model.Contact)) *model.Contact {
	c := &model.Contact{
		BaseModel: model.BaseModel{
			DomainId: domain,
		},
		IssuerId:      uuid.New(),
		ApplicationId: uuid.New(),
		Type:          model.User,
		Name:          "Antonio Banderas",
		Username:      "a.banderas@webitel.com",
		Metadata: map[string]string{
			"lang": "en",
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (suite *ContactStoreTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.NewPostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	suite.postgresContainer = pgContainer
	//TODO: add correct setup and migrations run
	contactStore := postgres.NewContactStore()

	suite.repo = contactStore
}

func (suite *ContactStoreTestSuite) TearDownSuite() {
	if err := suite.postgresContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

//#region Create

func (suite *ContactStoreTestSuite) TestCreate_HappyPath() {
	contact := newContact(1)

	created, err := suite.repo.Create(suite.ctx, contact)
	suite.Require().NoError(err)
	suite.NotNil(created)

	suite.NotEqual(uuid.Nil, created.Id)
	suite.NotZero(created.CreatedAt)
	suite.NotZero(created.UpdatedAt)
	suite.Equal(contact.Name, created.Name)
	suite.Equal(contact.Username, created.Username)
}

func (suite *ContactStoreTestSuite) TestCreate_MissingUsername() {
	contact := newContact(1, func(c *model.Contact) {
		c.Username = ""
	})

	_, err := suite.repo.Create(suite.ctx, contact)
	suite.Error(err)
}

func (suite *ContactStoreTestSuite) TestCreate_InvalidDomain() {
	contact := newContact(999999)

	_, err := suite.repo.Create(suite.ctx, contact)
	suite.Error(err)
}

func (suite *ContactStoreTestSuite) TestCreate_DuplicateUsername() {
	c1 := newContact(1)
	c2 := newContact(1, func(c *model.Contact) {
		c.Username = c1.Username
	})

	_, err := suite.repo.Create(suite.ctx, c1)
	suite.Require().NoError(err)

	_, err = suite.repo.Create(suite.ctx, c2)
	suite.Error(err)
}

//#endregion

// #region Search
func (suite *ContactStoreTestSuite) TestSearch_NoFilters() {
	_, err := suite.repo.Create(suite.ctx, newContact(1))
	suite.Require().NoError(err)
	_, err = suite.repo.Create(suite.ctx, newContact(1))
	suite.Require().NoError(err)

	res, err := suite.repo.Search(suite.ctx, &dto.ContactSearchFilter{
		Page: 1,
		Size: 10,
	})

	suite.Require().NoError(err)
	suite.Len(res, 2)
}

func (suite *ContactStoreTestSuite) TestSearch_Q_ByName() {
	_, err := suite.repo.Create(suite.ctx, newContact(1, func(c *model.Contact) {
		c.Name = "Angelina Jolie"
	}))
	suite.Require().NoError(err)

	_, err = suite.repo.Create(suite.ctx, newContact(1, func(c *model.Contact) {
		c.Name = "Bob"
	}))
	suite.Require().NoError(err)

	res, err := suite.repo.Search(suite.ctx, &dto.ContactSearchFilter{
		Page:   1,
		Size:   10,
		Q:      "Ang",
		Fields: []string{"name"},
	})

	suite.Require().NoError(err)
	suite.Len(res, 1)
	suite.Equal("Angelina Jolie", res[0].Name)
}

func (suite *ContactStoreTestSuite) TestSearch_ByApplication() {
	appID := uuid.New()

	_, err := suite.repo.Create(suite.ctx, newContact(1, func(c *model.Contact) {
		c.ApplicationId = appID
	}))
	suite.Require().NoError(err)

	_, err = suite.repo.Create(suite.ctx, newContact(1))
	suite.Require().NoError(err)

	res, err := suite.repo.Search(suite.ctx, &dto.ContactSearchFilter{
		Page: 1,
		Size: 10,
		Apps: []string{appID.String()},
	})

	suite.Require().NoError(err)
	suite.Len(res, 1)
}

func (suite *ContactStoreTestSuite) TestSearch_ByIssuer() {
	issuer := uuid.New()

	_, err := suite.repo.Create(suite.ctx, newContact(1, func(c *model.Contact) {
		c.IssuerId = issuer
	}))
	suite.Require().NoError(err)

	res, err := suite.repo.Search(suite.ctx, &dto.ContactSearchFilter{
		Page:    1,
		Size:    10,
		Issuers: []string{issuer.String()},
	})

	suite.Require().NoError(err)
	suite.Len(res, 1)
}

func (suite *ContactStoreTestSuite) TestSearch_Pagination() {
	for range 25 {
		_, _ = suite.repo.Create(suite.ctx, newContact(1, func(c *model.Contact) {
			c.Username = faker.Username()
		}))
	}

	res, err := suite.repo.Search(suite.ctx, &dto.ContactSearchFilter{
		Page: 2,
		Size: 10,
	})

	suite.Require().NoError(err)
	suite.Len(res, 11)
}

func (suite *ContactStoreTestSuite) TestSearch_EmptyResult() {
	res, err := suite.repo.Search(suite.ctx, &dto.ContactSearchFilter{
		Page: 1,
		Size: 10,
		Q:    "does-not-exist",
	})

	suite.Require().NoError(err)
	suite.Empty(res)
}

//#endregion

// #region Update
func (suite *ContactStoreTestSuite) TestUpdate_HappyPath() {
	created, _ := suite.repo.Create(suite.ctx, newContact(1))

	cmd := &dto.UpdateContactCommand{
		Id:       created.Id,
		Name:     "Angelina Jolie",
		Username: "a.jolie@webitel.com",
	}

	updated, err := suite.repo.Update(suite.ctx, cmd)
	suite.Require().NoError(err)

	suite.Equal("Angelina Jolie", updated.Name)
	suite.Greater(updated.UpdatedAt, created.UpdatedAt)
}

func (suite *ContactStoreTestSuite) TestUpdate_NotFound() {
	_, err := suite.repo.Update(suite.ctx, &dto.UpdateContactCommand{
		Id: uuid.New(),
	})

	suite.Error(err)
}

//#endregion

// #region Delete
func (suite *ContactStoreTestSuite) TestDelete_HappyPath() {
	created, _ := suite.repo.Create(suite.ctx, newContact(1))

	err := suite.repo.Delete(suite.ctx, created.Id)
	suite.Require().NoError(err)

	res, _ := suite.repo.Search(suite.ctx, &dto.ContactSearchFilter{
		Page: 1,
		Size: 10,
	})

	suite.Empty(res)
}

//#endregion
