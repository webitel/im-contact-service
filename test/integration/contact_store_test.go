package integration

import (
	"context"
	"log"
	"log/slog"
	"sync"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/webitel/im-contact-service/cmd/migrate"
	"github.com/webitel/im-contact-service/config"
	"github.com/webitel/im-contact-service/infra/db/pg"
	"github.com/webitel/im-contact-service/internal/domain/model"
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

func TestContactStoreTestSuite(t *testing.T) {
	suite.Run(t, new(ContactStoreTestSuite))
}

func newContact(domain int, opts ...func(*model.Contact)) *model.Contact {
	c := &model.Contact{
		BaseModel: model.BaseModel{
			DomainId: domain,
		},
		IssuerId:      uuid.New().String(),
		ApplicationId: uuid.New().String(),
		SubjectId:     uuid.New().String(),
		Type:          "webitel",
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

	mig := migrate.NewMigrator(&config.Config{Postgres: config.PostgresConfig{DSN: pgContainer.ConnectionString}}, slog.Default())
	if err := mig.Run(suite.ctx); err != nil {
		log.Fatal(err)
	}

	db, err := pg.New(suite.ctx, slog.Default(), pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	contactStore := postgres.NewContactStore(db)

	suite.repo = contactStore
}

func (suite *ContactStoreTestSuite) SetupTest() {
	truncateCmd := []string{
		"psql",
		"-U", "opensips",
		"-d", "webitel",
		"-c", "TRUNCATE TABLE im_contact.contact CASCADE;",
	}

	exitCode, _, err := suite.postgresContainer.Exec(suite.ctx, truncateCmd)
	if err != nil {
		log.Fatalf("failed to execute truncate command: %v", err)
	}

	if exitCode != 0 {
		log.Fatalf("truncate command failed with exit code: %d", exitCode)
	}
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

func (suite *ContactStoreTestSuite) TestCreate_DuplicateSubject() {
	c1 := newContact(1)
	c2 := newContact(1, func(c *model.Contact) {
		c.SubjectId = c1.SubjectId
	}, func(c *model.Contact) {
		c.IssuerId = c1.IssuerId
	})

	_, err := suite.repo.Create(suite.ctx, c1)
	suite.Require().NoError(err)

	_, err = suite.repo.Create(suite.ctx, c2)
	suite.Error(err)
}

func (suite *ContactStoreTestSuite) TestCreate_WithMetadata() {
	contact := newContact(1, func(c *model.Contact) {
		c.Metadata = map[string]string{
			"lang":   "uk",
			"tz":     "Europe/Kyiv",
			"source": "landing",
		}
	})

	created, err := suite.repo.Create(suite.ctx, contact)
	suite.Require().NoError(err)

	suite.NotNil(created.Metadata)
	suite.Equal("uk", created.Metadata["lang"])
	suite.Equal("Europe/Kyiv", created.Metadata["tz"])
	suite.Equal("landing", created.Metadata["source"])
}

func (suite *ContactStoreTestSuite) TestCreate_WithNilMetadata() {
	contact := newContact(1, func(c *model.Contact) {
		c.Metadata = nil
	})

	created, err := suite.repo.Create(suite.ctx, contact)
	suite.Require().NoError(err)

	suite.True(len(created.Metadata) == 0)
}

//#endregion

// #region Search
func (suite *ContactStoreTestSuite) TestSearch_NoFilters() {
	_, err := suite.repo.Create(suite.ctx, newContact(1))
	suite.Require().NoError(err)
	_, err = suite.repo.Create(suite.ctx, newContact(1))
	suite.Require().NoError(err)

	res, err := suite.repo.Search(suite.ctx, &dto.ContactSearchFilter{
		Page:     1,
		Size:     10,
		DomainId: 1,
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

	q := "Ang"

	res, err := suite.repo.Search(suite.ctx, &dto.ContactSearchFilter{
		Page:     1,
		Size:     10,
		Q:        &q,
		Fields:   []string{"name", "id"},
		DomainId: 1,
	})

	suite.Require().NoError(err)
	suite.Len(res, 1)
	suite.Equal("Angelina Jolie", res[0].Name)
}

func (suite *ContactStoreTestSuite) TestSearch_ByApplication() {
	appID := uuid.New()

	_, err := suite.repo.Create(suite.ctx, newContact(1, func(c *model.Contact) {
		c.ApplicationId = appID.String()
	}))
	suite.Require().NoError(err)

	_, err = suite.repo.Create(suite.ctx, newContact(1))
	suite.Require().NoError(err)

	res, err := suite.repo.Search(suite.ctx, &dto.ContactSearchFilter{
		Page:     1,
		Size:     10,
		Apps:     []string{appID.String()},
		DomainId: 1,
	})

	suite.Require().NoError(err)
	suite.Len(res, 1)
}

func (suite *ContactStoreTestSuite) TestSearch_ByIssuer() {
	issuer := uuid.New()

	_, err := suite.repo.Create(suite.ctx, newContact(1, func(c *model.Contact) {
		c.IssuerId = issuer.String()
	}))
	suite.Require().NoError(err)

	res, err := suite.repo.Search(suite.ctx, &dto.ContactSearchFilter{
		Page:     1,
		Size:     10,
		Issuers:  []string{issuer.String()},
		DomainId: 1,
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
		Page:     2,
		Size:     10,
		DomainId: 1,
	})

	suite.Require().NoError(err)
	suite.Len(res, 11)
}

func (suite *ContactStoreTestSuite) TestSearch_EmptyResult() {
	q := "does-not-exist"

	res, err := suite.repo.Search(suite.ctx, &dto.ContactSearchFilter{
		Page:     1,
		Size:     10,
		Q:        &q,
		DomainId: 1,
	})

	suite.Require().NoError(err)
	suite.Empty(res)
}

//#endregion

// #region Update
func (suite *ContactStoreTestSuite) TestUpdate_HappyPath() {
	created, _ := suite.repo.Create(suite.ctx, newContact(1))

	var (
		updatedName     string = "Angelina Jolie"
		updatedUsername string = "a.jolie@webitel.com"
	)

	cmd := &dto.UpdateContactCommand{
		Id:       created.Id,
		DomainId: created.DomainId,
		Name:     &updatedName,
		Username: &updatedUsername,
	}

	updated, err := suite.repo.Update(suite.ctx, cmd)
	suite.Require().NoError(err)

	suite.Equal("Angelina Jolie", updated.Name)
	suite.Greater(updated.UpdatedAt, created.UpdatedAt)
}

func (suite *ContactStoreTestSuite) TestUpdate_Metadata_Clear() {
	created, _ := suite.repo.Create(suite.ctx, newContact(1, func(c *model.Contact) {
		c.Metadata = map[string]string{
			"lang": "en",
		}
	}))

	empty := map[string]string{}

	cmd := &dto.UpdateContactCommand{
		Id:       created.Id,
		DomainId: created.DomainId,
		Metadata: empty,
	}

	updated, err := suite.repo.Update(suite.ctx, cmd)
	suite.Require().NoError(err)

	suite.Empty(updated.Metadata)
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

	command := &dto.DeleteContactCommand{
		Id:       created.Id,
		DomainId: created.DomainId,
	}
	err := suite.repo.Delete(suite.ctx, command)
	suite.Require().NoError(err)

	res, _ := suite.repo.Search(suite.ctx, &dto.ContactSearchFilter{
		Page:     1,
		Size:     10,
		DomainId: 1,
	})

	suite.Empty(res)
}

func (suite *ContactStoreTestSuite) TestClearByDomain() {
	for range 10 {
		suite.repo.Create(suite.ctx, newContact(1))
	}

	for range 5 {
		suite.repo.Create(suite.ctx, newContact(2))
	}

	err := suite.repo.ClearByDomain(suite.ctx, 1)
	suite.Require().NoError(err)

	res, _ := suite.repo.Search(suite.ctx, &dto.ContactSearchFilter{
		Page:     1,
		Size:     6,
		DomainId: 2,
	})

	suite.Len(res, 5)
}

//#endregion

// #region Upsert

func (suite *ContactStoreTestSuite) TestUpsert_Insert_HappyPath() {
	contact := newContact(1)

	result, isInsert, err := suite.repo.Upsert(suite.ctx, contact)

	suite.Require().NoError(err)
	suite.True(isInsert, "Expected insert operation")
	suite.NotNil(result)
	suite.NotEqual(uuid.Nil, result.Id)
	suite.NotZero(result.CreatedAt)
	suite.NotZero(result.UpdatedAt)
	suite.Equal(contact.Name, result.Name)
	suite.Equal(contact.Username, result.Username)
	suite.Equal(contact.IssuerId, result.IssuerId)
	suite.Equal(contact.SubjectId, result.SubjectId)
}

func (suite *ContactStoreTestSuite) TestUpsert_Update_HappyPath() {
	contact := newContact(1, func(c *model.Contact) {
		c.Name = "Original Name"
		c.Username = "original@example.com"
	})

	first, isInsert, err := suite.repo.Upsert(suite.ctx, contact)
	suite.Require().NoError(err)
	suite.True(isInsert)

	updatedContact := newContact(1, func(c *model.Contact) {
		c.IssuerId = contact.IssuerId
		c.SubjectId = contact.SubjectId
		c.Name = "Updated Name"
		c.Username = "updated@example.com"
	})

	second, isInsert, err := suite.repo.Upsert(suite.ctx, updatedContact)
	suite.Require().NoError(err)
	suite.False(isInsert, "Expected update operation")
	suite.Equal(first.Id, second.Id, "ID should remain the same")
	suite.Equal("Updated Name", second.Name)
	suite.Equal("updated@example.com", second.Username)
	suite.Equal(first.CreatedAt, second.CreatedAt, "CreatedAt should not change")
	suite.Greater(second.UpdatedAt, first.UpdatedAt, "UpdatedAt should be newer")
}

func (suite *ContactStoreTestSuite) TestUpsert_Update_Metadata() {
	contact := newContact(1, func(c *model.Contact) {
		c.Metadata = map[string]string{
			"lang": "en",
			"tz":   "UTC",
		}
	})

	first, _, err := suite.repo.Upsert(suite.ctx, contact)
	suite.Require().NoError(err)
	suite.Equal("en", first.Metadata["lang"])

	updated := newContact(1, func(c *model.Contact) {
		c.IssuerId = contact.IssuerId
		c.SubjectId = contact.SubjectId
		c.Metadata = map[string]string{
			"lang":   "uk",
			"tz":     "Europe/Kyiv",
			"source": "landing",
		}
	})

	second, isInsert, err := suite.repo.Upsert(suite.ctx, updated)
	suite.Require().NoError(err)
	suite.False(isInsert)
	suite.Equal("uk", second.Metadata["lang"])
	suite.Equal("Europe/Kyiv", second.Metadata["tz"])
	suite.Equal("landing", second.Metadata["source"])
}

func (suite *ContactStoreTestSuite) TestUpsert_Update_MetadataClear() {
	contact := newContact(1, func(c *model.Contact) {
		c.Metadata = map[string]string{
			"lang": "en",
		}
	})

	first, _, err := suite.repo.Upsert(suite.ctx, contact)
	suite.Require().NoError(err)
	suite.NotEmpty(first.Metadata)

	updated := newContact(1, func(c *model.Contact) {
		c.IssuerId = contact.IssuerId
		c.SubjectId = contact.SubjectId
		c.Metadata = nil
	})

	second, isInsert, err := suite.repo.Upsert(suite.ctx, updated)
	suite.Require().NoError(err)
	suite.False(isInsert)
	suite.Empty(second.Metadata)
}

func (suite *ContactStoreTestSuite) TestUpsert_Insert_WithNilMetadata() {
	contact := newContact(1, func(c *model.Contact) {
		c.Metadata = nil
	})

	result, isInsert, err := suite.repo.Upsert(suite.ctx, contact)

	suite.Require().NoError(err)
	suite.True(isInsert)
	suite.True(len(result.Metadata) == 0)
}

func (suite *ContactStoreTestSuite) TestUpsert_Insert_MissingUsername() {
	contact := newContact(1, func(c *model.Contact) {
		c.Username = ""
	})

	_, _, err := suite.repo.Upsert(suite.ctx, contact)
	suite.Error(err, "Should fail on empty username")
}

func (suite *ContactStoreTestSuite) TestUpsert_MultipleContactsDifferentDomains() {
	contact1 := newContact(1, func(c *model.Contact) {
		c.IssuerId = "same-issuer"
		c.SubjectId = "same-subject"
		c.Name = "Domain 1 Contact"
	})

	contact2 := newContact(2, func(c *model.Contact) {
		c.IssuerId = "same-issuer"
		c.SubjectId = "same-subject"
		c.Name = "Domain 2 Contact"
	})

	result1, isInsert1, err := suite.repo.Upsert(suite.ctx, contact1)
	suite.Require().NoError(err)
	suite.True(isInsert1)

	result2, isInsert2, err := suite.repo.Upsert(suite.ctx, contact2)
	suite.Require().NoError(err)
	suite.True(isInsert2)

	suite.NotEqual(result1.Id, result2.Id, "Different domains should have different IDs")
	suite.Equal("Domain 1 Contact", result1.Name)
	suite.Equal("Domain 2 Contact", result2.Name)
}

func (suite *ContactStoreTestSuite) TestUpsert_ConsecutiveUpdates() {
	contact := newContact(1)

	v1, isInsert, err := suite.repo.Upsert(suite.ctx, contact)
	suite.Require().NoError(err)
	suite.True(isInsert)

	contact.Name = "Name v2"
	v2, isInsert, err := suite.repo.Upsert(suite.ctx, contact)
	suite.Require().NoError(err)
	suite.False(isInsert)
	suite.Equal("Name v2", v2.Name)
	suite.Greater(v2.UpdatedAt, v1.UpdatedAt)

	contact.Name = "Name v3"
	v3, isInsert, err := suite.repo.Upsert(suite.ctx, contact)
	suite.Require().NoError(err)
	suite.False(isInsert)
	suite.Equal("Name v3", v3.Name)
	suite.Greater(v3.UpdatedAt, v2.UpdatedAt)

	suite.Equal(v1.Id, v2.Id)
	suite.Equal(v1.Id, v3.Id)
	suite.Equal(v1.CreatedAt, v2.CreatedAt)
	suite.Equal(v1.CreatedAt, v3.CreatedAt)
}

func (suite *ContactStoreTestSuite) TestUpsert_ConcurrentOperations() {
	contact := newContact(1)

	var wg sync.WaitGroup
	results := make([]bool, 5)
	errors := make([]error, 5)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_, isInsert, err := suite.repo.Upsert(suite.ctx, contact)
			results[idx] = isInsert
			errors[idx] = err
		}(i)
	}

	wg.Wait()

	for i, err := range errors {
		suite.NoError(err, "Operation %d should not error", i)
	}

	insertCount := 0
	for _, isInsert := range results {
		if isInsert {
			insertCount++
		}
	}
	suite.Equal(1, insertCount, "Exactly one operation should be insert")
}

//#endregion