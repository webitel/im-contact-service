package integration

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"testing"

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

type BotStoreTestSuite struct {
	suite.Suite

	postgresContainer *testhelpers.PostgresContainer
	repo store.BotStore
	ctx context.Context
}

func TestBotStoreTestSuite(t *testing.T) {
	suite.Run(t, new(BotStoreTestSuite))
}

func newBot(domainId int, opts ...func(*model.WebitelBot)) *model.WebitelBot {
	b := &model.WebitelBot{
		BaseModel: model.BaseModel{
			DomainId: domainId,
		},
		FlowId: 100,
		DisplayName: "Antonio Banderas",
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func (suite *BotStoreTestSuite) SetupSuite() {
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

	botStore := postgres.NewBotStore(db)

	suite.repo = botStore
}

func (suite *BotStoreTestSuite) SetupTest() {
	truncateCmd := []string{
		"psql",
		"-U", "opensips",
		"-d", "webitel",
		"-c", "TRUNCATE TABLE im_contact.bot CASCADE;",
	}

	exitCode, _, err := suite.postgresContainer.Exec(suite.ctx, truncateCmd)
	if err != nil {
		log.Fatalf("failed to execute truncate command: %v", err)
	}

	if exitCode != 0 {
		log.Fatalf("truncate command failed with exit code: %d", exitCode)
	}
} 

func (suite *BotStoreTestSuite) TestCreate_Success() {
    bot := newBot(1)
    
    createdBot, err := suite.repo.Create(suite.ctx, bot)
    
    suite.NoError(err)
    suite.NotNil(createdBot)
    suite.NotEqual(uuid.Nil, createdBot.Id)
    suite.NotZero(createdBot.CreatedAt)
    suite.NotZero(createdBot.UpdatedAt)
    suite.Equal(bot.DomainId, createdBot.DomainId)
    suite.Equal(bot.FlowId, createdBot.FlowId)
    suite.Equal(bot.DisplayName, createdBot.DisplayName)
}

func (suite *BotStoreTestSuite) TestCreate_MultipleBots() {
    bot1 := newBot(1, func(b *model.WebitelBot) { b.DisplayName = "Bot One" }, func(wb *model.WebitelBot) {wb.FlowId = 99})
    bot2 := newBot(1, func(b *model.WebitelBot) { b.DisplayName = "Bot Two" }, func(wb *model.WebitelBot) {wb.FlowId = 100})
    
    createdBot1, err1 := suite.repo.Create(suite.ctx, bot1)
    createdBot2, err2 := suite.repo.Create(suite.ctx, bot2)
    
    suite.NoError(err1)
    suite.NoError(err2)
    suite.NotEqual(createdBot1.Id, createdBot2.Id)
}

func (suite *BotStoreTestSuite) TestSearch_ByDomainId() {
    bot1 := newBot(1, func(b *model.WebitelBot) { b.DisplayName = "Domain 1 Bot" }, func(wb *model.WebitelBot) {wb.FlowId = 99})
    bot2 := newBot(2, func(b *model.WebitelBot) { b.DisplayName = "Domain 2 Bot" }, func(wb *model.WebitelBot) {wb.FlowId = 101})
    bot3 := newBot(1, func(b *model.WebitelBot) { b.DisplayName = "Another Domain 1 Bot" }, func(wb *model.WebitelBot) {wb.FlowId = 100})
    
    suite.repo.Create(suite.ctx, bot1)
    suite.repo.Create(suite.ctx, bot2)
    suite.repo.Create(suite.ctx, bot3)
    
    filter := &dto.SearchBotRequest{
		BaseFilter: dto.BaseFilter{
			DomainId: 1,
			Page: 1,
			Size: 10,
			Sort: "+created_at",
		},
    }
    
    bots, err := suite.repo.Search(suite.ctx, filter)
    
    suite.NoError(err)
    suite.Len(bots, 2)
    for _, bot := range bots {
        suite.Equal(1, bot.DomainId)
    }
}

func (suite *BotStoreTestSuite) TestSearch_ByQ() {
    bot1 := newBot(1, func(b *model.WebitelBot) { b.DisplayName = "Antonio Banderas" }, func(wb *model.WebitelBot) {wb.FlowId = rand.Int()})
    bot2 := newBot(1, func(b *model.WebitelBot) { b.DisplayName = "Tom Cruise" }, func(wb *model.WebitelBot) {wb.FlowId = rand.Int()})
    bot3 := newBot(1, func(b *model.WebitelBot) { b.DisplayName = "Antonio Lopez" }, func(wb *model.WebitelBot) {wb.FlowId = rand.Int()})
    
    suite.repo.Create(suite.ctx, bot1)
    suite.repo.Create(suite.ctx, bot2)
    suite.repo.Create(suite.ctx, bot3)
    
    q := "%antonio%"
    filter := &dto.SearchBotRequest{
		BaseFilter: dto.BaseFilter{
			DomainId: 1,
			Page: 1,
			Q: q,
			Size: 10,
			Sort: "+display_name",
		},
    }
    
    bots, err := suite.repo.Search(suite.ctx, filter)
    
    suite.NoError(err)
    suite.Len(bots, 2)
    suite.Contains(bots[0].DisplayName, "Antonio")
    suite.Contains(bots[1].DisplayName, "Antonio")
}

func (suite *BotStoreTestSuite) TestSearch_ByFlowIds() {
    bot1 := newBot(1, func(b *model.WebitelBot) { b.FlowId = 100 })
    bot2 := newBot(1, func(b *model.WebitelBot) { b.FlowId = 200 })
    bot3 := newBot(1, func(b *model.WebitelBot) { b.FlowId = 300 })
    
    createdBot1, _ := suite.repo.Create(suite.ctx, bot1)
    createdBot2, _ := suite.repo.Create(suite.ctx, bot2)
    suite.repo.Create(suite.ctx, bot3)
    
    flowIds := []int64{int64(createdBot1.FlowId), int64(createdBot2.FlowId)}
    filter := &dto.SearchBotRequest{
		BaseFilter: dto.BaseFilter{
			DomainId: 1,
			Page: 1,
			Size: 10,
			Sort: "+id",
		},
        FlowIds:  flowIds,
    }
    
    bots, err := suite.repo.Search(suite.ctx, filter)
    
    suite.NoError(err)
    suite.Len(bots, 2)
}

func (suite *BotStoreTestSuite) TestSearch_ByDisplayNames() {
    bot1 := newBot(1, func(b *model.WebitelBot) { b.DisplayName = "Alpha Bot" }, func(wb *model.WebitelBot) {wb.FlowId = rand.Int()})
    bot2 := newBot(1, func(b *model.WebitelBot) { b.DisplayName = "Beta Bot" }, func(wb *model.WebitelBot) {wb.FlowId = rand.Int()})
    bot3 := newBot(1, func(b *model.WebitelBot) { b.DisplayName = "Gamma Bot" }, func(wb *model.WebitelBot) {wb.FlowId = rand.Int()})
    
    suite.repo.Create(suite.ctx, bot1)
    suite.repo.Create(suite.ctx, bot2)
    suite.repo.Create(suite.ctx, bot3)
    
    displayNames := []string{"%alpha%", "%beta%"}
    filter := &dto.SearchBotRequest{
		BaseFilter: dto.BaseFilter{
			DomainId: 1,
			Page: 1,
			Size: 10,
			Sort: "+display_name",
		},
        DisplayNames: displayNames,
    }
    
    bots, err := suite.repo.Search(suite.ctx, filter)
    
    suite.NoError(err)
    suite.Len(bots, 2)
}

func (suite *BotStoreTestSuite) TestSearch_ByIds() {
    bot1 := newBot(1, func(b *model.WebitelBot) { b.DisplayName = "Bot 1" }, func(wb *model.WebitelBot) {wb.FlowId = rand.Int()})
    bot2 := newBot(1, func(b *model.WebitelBot) { b.DisplayName = "Bot 2" }, func(wb *model.WebitelBot) {wb.FlowId = rand.Int()})
    bot3 := newBot(1, func(b *model.WebitelBot) { b.DisplayName = "Bot 3" }, func(wb *model.WebitelBot) {wb.FlowId = rand.Int()})
    
    createdBot1, _ := suite.repo.Create(suite.ctx, bot1)
    createdBot2, _ := suite.repo.Create(suite.ctx, bot2)
    suite.repo.Create(suite.ctx, bot3)
    
    ids := []uuid.UUID{createdBot1.Id, createdBot2.Id}
    filter := &dto.SearchBotRequest{
		BaseFilter: dto.BaseFilter{
			DomainId: 1,
			Page: 1,
			Size: 10,
			Sort: "+created_at",
		},
        Ids:      ids,
    }
    
    bots, err := suite.repo.Search(suite.ctx, filter)
    
    suite.NoError(err)
    suite.Len(bots, 2)
    suite.ElementsMatch(ids, []uuid.UUID{bots[0].Id, bots[1].Id})
}

func (suite *BotStoreTestSuite) TestSearch_Pagination() {
    for i := 1; i <= 5; i++ {
        bot := newBot(1, func(b *model.WebitelBot) {
            b.DisplayName = fmt.Sprintf("Bot %d", i)
        }, func(wb *model.WebitelBot) {wb.FlowId = rand.Int()})
        suite.repo.Create(suite.ctx, bot)
    }
    
    filter := &dto.SearchBotRequest{
		BaseFilter: dto.BaseFilter{
			DomainId: 1,
			Page: 1,
			Size: 2,
			Sort: "+display_name",
		},
    }
    
    page1, err := suite.repo.Search(suite.ctx, filter)
    suite.NoError(err)
    suite.Len(page1, 3)
    
    filter.Page = 2
    page2, err := suite.repo.Search(suite.ctx, filter)
    suite.NoError(err)
    suite.Len(page2, 3)
    
    suite.NotEqual(page1[0].Id, page2[0].Id)
}

func (suite *BotStoreTestSuite) TestSearch_EmptyResult() {
    filter := &dto.SearchBotRequest{
		BaseFilter: dto.BaseFilter{
			DomainId: 999,
			Page: 1,
			Size: 10,
			Sort: "+id",
		},
    }
    
    bots, err := suite.repo.Search(suite.ctx, filter)
    
    suite.NoError(err)
    suite.Empty(bots)
}

func (suite *BotStoreTestSuite) TestSearch_Sorting() {
    bot1 := newBot(1, func(b *model.WebitelBot) { b.DisplayName = "Charlie" }, func(wb *model.WebitelBot) {wb.FlowId = 2})
    bot2 := newBot(1, func(b *model.WebitelBot) { b.DisplayName = "Alice" }, func(wb *model.WebitelBot) {wb.FlowId = 3})
    bot3 := newBot(1, func(b *model.WebitelBot) { b.DisplayName = "Bob" }, func(wb *model.WebitelBot) {wb.FlowId = 4})
    
    suite.repo.Create(suite.ctx, bot1)
    suite.repo.Create(suite.ctx, bot2)
    suite.repo.Create(suite.ctx, bot3)
    
    filter := &dto.SearchBotRequest{
		BaseFilter: dto.BaseFilter{
			DomainId: 1,
			Page: 1,
			Size: 10,
			Sort: "+display_name",
		},
    }
    
    bots, err := suite.repo.Search(suite.ctx, filter)
    
    suite.NoError(err)
    suite.Len(bots, 3)
    suite.Equal("Alice", bots[0].DisplayName)
    suite.Equal("Bob", bots[1].DisplayName)
    suite.Equal("Charlie", bots[2].DisplayName)
}

func (suite *BotStoreTestSuite) TestUpdate_Success() {
    bot := newBot(1, func(b *model.WebitelBot) {
        b.DisplayName = "Original Name"
        b.FlowId = 100
    })
    createdBot, _ := suite.repo.Create(suite.ctx, bot)
    
    newFlowId := 200
    newDisplayName := "Updated Name"
    updateCmd := &dto.UpdateBotCommand{
        Id:          createdBot.Id,
        DomainId:    createdBot.DomainId,
        FlowId:      &newFlowId,
        DisplayName: &newDisplayName,
    }
    
    updatedBot, err := suite.repo.Update(suite.ctx, updateCmd)
    
    suite.NoError(err)
    suite.Equal(createdBot.Id, updatedBot.Id)
    suite.Equal(newFlowId, updatedBot.FlowId)
    suite.Equal(newDisplayName, updatedBot.DisplayName)
    suite.True(updatedBot.UpdatedAt.After(createdBot.UpdatedAt))
}

func (suite *BotStoreTestSuite) TestUpdate_PartialUpdate_OnlyFlowId() {
    bot := newBot(1, func(b *model.WebitelBot) {
        b.DisplayName = "Original Name"
        b.FlowId = 100
    })
    createdBot, _ := suite.repo.Create(suite.ctx, bot)
    
    newFlowId := 200
    updateCmd := &dto.UpdateBotCommand{
        Id:       createdBot.Id,
        DomainId: createdBot.DomainId,
        FlowId:   &newFlowId,
    }
    
    updatedBot, err := suite.repo.Update(suite.ctx, updateCmd)
    
    suite.NoError(err)
    suite.Equal(newFlowId, updatedBot.FlowId)
    suite.Equal(createdBot.DisplayName, updatedBot.DisplayName)
}

func (suite *BotStoreTestSuite) TestUpdate_PartialUpdate_OnlyDisplayName() {
    bot := newBot(1, func(b *model.WebitelBot) {
        b.DisplayName = "Original Name"
        b.FlowId = 100
    })
    createdBot, _ := suite.repo.Create(suite.ctx, bot)
    
    newDisplayName := "Updated Name"
    updateCmd := &dto.UpdateBotCommand{
        Id:          createdBot.Id,
        DomainId:    createdBot.DomainId,
        DisplayName: &newDisplayName,
    }
    
    updatedBot, err := suite.repo.Update(suite.ctx, updateCmd)
    
    suite.NoError(err)
    suite.Equal(newDisplayName, updatedBot.DisplayName)
    suite.Equal(createdBot.FlowId, updatedBot.FlowId)
}

func (suite *BotStoreTestSuite) TestUpdate_NonExistentBot() {
    newFlowId := 200
    updateCmd := &dto.UpdateBotCommand{
        Id:       uuid.New(),
        DomainId: 1,
        FlowId:   &newFlowId,
    }
    
    _, err := suite.repo.Update(suite.ctx, updateCmd)
    
    suite.Error(err)
}

func (suite *BotStoreTestSuite) TestUpdate_WrongDomain() {
    bot := newBot(1)
    createdBot, _ := suite.repo.Create(suite.ctx, bot)
    
    newFlowId := 200
    updateCmd := &dto.UpdateBotCommand{
        Id:       createdBot.Id,
        DomainId: 999,
        FlowId:   &newFlowId,
    }
    
    _, err := suite.repo.Update(suite.ctx, updateCmd)
    
    suite.Error(err)
}

func (suite *BotStoreTestSuite) TestDelete_ById() {
    bot := newBot(1)
    createdBot, _ := suite.repo.Create(suite.ctx, bot)
    
    deleteCmd := &dto.DeleteBotCommand{
        Id:       &createdBot.Id,
        DomainId: createdBot.DomainId,
    }
    
    err := suite.repo.Delete(suite.ctx, deleteCmd)
    
    suite.NoError(err)
    
    filter := &dto.SearchBotRequest{
		BaseFilter: dto.BaseFilter{
			DomainId: createdBot.DomainId,
			Page: 1,
			Size: 10,
			Sort: "+id",
		},
        Ids:      []uuid.UUID{createdBot.Id},
    }
    bots, _ := suite.repo.Search(suite.ctx, filter)
    suite.Empty(bots)
}

func (suite *BotStoreTestSuite) TestDelete_ByFlowId() {
    bot1 := newBot(1, func(b *model.WebitelBot) { b.FlowId = 100 })
    bot2 := newBot(1, func(b *model.WebitelBot) { b.FlowId = 101 })
    bot3 := newBot(1, func(b *model.WebitelBot) { b.FlowId = 200 })
    
    suite.repo.Create(suite.ctx, bot1)
    suite.repo.Create(suite.ctx, bot2)
    suite.repo.Create(suite.ctx, bot3)
    
    flowId := 100
    deleteCmd := &dto.DeleteBotCommand{
        FlowId:   &flowId,
        DomainId: 1,
    }
    
    err := suite.repo.Delete(suite.ctx, deleteCmd)
    
    suite.NoError(err)
    
    filter := &dto.SearchBotRequest{
		BaseFilter: dto.BaseFilter{
			DomainId: 1,
			Page: 1,
			Size: 10,
			Sort: "+id",
		},
    }
    bots, _ := suite.repo.Search(suite.ctx, filter)
    suite.Len(bots, 2)
}

func (suite *BotStoreTestSuite) TestDelete_NonExistentBot() {
    id := uuid.New()
    deleteCmd := &dto.DeleteBotCommand{
        Id:       &id,
        DomainId: 1,
    }
    
    err := suite.repo.Delete(suite.ctx, deleteCmd)
    
    suite.NoError(err)
}

func (suite *BotStoreTestSuite) TestDelete_WrongDomain() {
    bot := newBot(1)
    createdBot, _ := suite.repo.Create(suite.ctx, bot)
    
    deleteCmd := &dto.DeleteBotCommand{
        Id:       &createdBot.Id,
        DomainId: 999,
    }
    
    err := suite.repo.Delete(suite.ctx, deleteCmd)
    
    suite.NoError(err)
    
    filter := &dto.SearchBotRequest{
		BaseFilter: dto.BaseFilter{
			DomainId: 1,
			Page: 1,
			Size: 10,
			Sort: "+id",
		},
        Ids:      []uuid.UUID{createdBot.Id},
    }
    bots, _ := suite.repo.Search(suite.ctx, filter)
    suite.Len(bots, 1)
}

func (suite *BotStoreTestSuite) TestComplexScenario_CreateSearchUpdateDelete() {
    bot1 := newBot(1, func(b *model.WebitelBot) {
        b.DisplayName = "Test Bot 1"
        b.FlowId = 100
    })
    bot2 := newBot(1, func(b *model.WebitelBot) {
        b.DisplayName = "Test Bot 2"
        b.FlowId = 200
    })
    
    createdBot1, err := suite.repo.Create(suite.ctx, bot1)
    suite.NoError(err)
    createdBot2, err := suite.repo.Create(suite.ctx, bot2)
    suite.NoError(err)
    
    filter := &dto.SearchBotRequest{
		BaseFilter: dto.BaseFilter{
			DomainId: 1,
			Page: 1,
			Size: 10,
			Sort: "+display_name",
		},
    }
    bots, err := suite.repo.Search(suite.ctx, filter)
    suite.NoError(err)
    suite.Len(bots, 2)
    
    newName := "Updated Bot 1"
    updateCmd := &dto.UpdateBotCommand{
        Id:          createdBot1.Id,
        DomainId:    1,
        DisplayName: &newName,
    }
    updatedBot, err := suite.repo.Update(suite.ctx, updateCmd)
    suite.NoError(err)
    suite.Equal(newName, updatedBot.DisplayName)
    
    deleteCmd := &dto.DeleteBotCommand{
        Id:       &createdBot2.Id,
        DomainId: 1,
    }
    err = suite.repo.Delete(suite.ctx, deleteCmd)
    suite.NoError(err)
    
    finalBots, err := suite.repo.Search(suite.ctx, filter)
    suite.NoError(err)
    suite.Len(finalBots, 1)
    suite.Equal(createdBot1.Id, finalBots[0].Id)
    suite.Equal(newName, finalBots[0].DisplayName)
}

func (suite *BotStoreTestSuite) TearDownSuite() {
    if suite.postgresContainer != nil {
        suite.postgresContainer.Terminate(suite.ctx)
    }
}

