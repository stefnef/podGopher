package user

import (
	"podGopher/adapter/outbound/repository/postgres/postgresTestSetup"
	repositoryShow "podGopher/adapter/outbound/repository/postgres/show"
	"podGopher/core/domain/model"
	forSaveUser "podGopher/core/port/outbound/user"
	"sort"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_user_repository_should_implement_port(t *testing.T) {
	repository := NewPostgresUserRepository(nil)

	assert.NotNil(t, repository)
	assert.Implements(t, (*forSaveUser.SaveUserPort)(nil), repository)
}

func Test_should_not_save_user_if_show_does_not_exist(t *testing.T) {
	db := postgresTestSetup.StartTestcontainersPostgres(t, "../postgresTestSetup/")
	defer postgresTestSetup.Teardown(t, db)

	nonExistingShowId := uuid.NewString()

	repository := NewPostgresUserRepository(db)
	user := &model.User{
		Id:       uuid.NewString(),
		Username: "some-username",
		ShowRoles: []model.ShowRole{
			{
				ShowId: nonExistingShowId,
				Role:   model.EDITOR,
			},
		},
	}
	err := repository.SaveUser(user)
	assert.NotNil(t, err)
}

func Test_should_save_a_user(t *testing.T) {
	db := postgresTestSetup.StartTestcontainersPostgres(t, "../postgresTestSetup/")
	defer postgresTestSetup.Teardown(t, db)

	editorShowUuid := uuid.NewString()
	producerShowUuid := uuid.NewString()

	showRepository := repositoryShow.NewPostgresShowRepository(db)
	repository := NewPostgresUserRepository(db)
	username := "Some username"
	user := &model.User{
		Id:       uuid.NewString(),
		Username: username,
		ShowRoles: []model.ShowRole{
			{ShowId: producerShowUuid, Role: model.PRODUCER},
			{ShowId: editorShowUuid, Role: model.EDITOR},
		},
	}

	if err := showRepository.SaveShow(&model.Show{Id: editorShowUuid, Title: "editor-show", Slug: "editor-slug"}); err != nil {
		t.Fatal(err)
	}
	if err := showRepository.SaveShow(&model.Show{Id: producerShowUuid, Title: "producer-show", Slug: "producer-slug"}); err != nil {
		t.Fatal(err)
	}

	t.Run("should return false if user with username does not exist", func(t *testing.T) {
		exists := repository.ExistsByUsername(producerShowUuid, username)
		assert.False(t, exists)

		exists = repository.ExistsByUsername(editorShowUuid, username)
		assert.False(t, exists)
	})

	t.Run("should save a user", func(t *testing.T) {
		err := repository.SaveUser(user)
		assert.Nil(t, err)
	})

	t.Run("should return true if user with username exists", func(t *testing.T) {
		exists := repository.ExistsByUsername(producerShowUuid, username)
		assert.True(t, exists)

		exists = repository.ExistsByUsername(editorShowUuid, username)
		assert.True(t, exists)
	})

	t.Run("should return false if user with username does not exists", func(t *testing.T) {
		exists := repository.ExistsByUsername(editorShowUuid, "some-other-user")
		assert.False(t, exists)

		exists = repository.ExistsByUsername("i-do-not-exist", username)
		assert.False(t, exists)
	})

	t.Run("should query a user", func(t *testing.T) {
		var id, username string
		err := db.QueryRow("SELECT * FROM podcast_user WHERE id = $1", user.Id).
			Scan(&id, &username)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, user.Id, id)
		assert.Equal(t, user.Username, username)
	})

	t.Run("should query show roles", func(t *testing.T) {
		var role string
		err := db.QueryRow("SELECT role FROM user_roles WHERE user_id = $1 and show_id = $2", user.Id, editorShowUuid).
			Scan(&role)
		assert.Nil(t, err)
		assert.Equal(t, model.EDITOR.Name(), role)

		err = db.QueryRow("SELECT role FROM user_roles WHERE user_id = $1 and show_id = $2", user.Id, producerShowUuid).
			Scan(&role)
		assert.Nil(t, err)
		assert.Equal(t, model.PRODUCER.Name(), role)
	})
}

func Test_should_retrieve_a_user(t *testing.T) {
	db := postgresTestSetup.StartTestcontainersPostgres(t, "../postgresTestSetup/")

	defer postgresTestSetup.Teardown(t, db)

	editorShowUuid := uuid.NewString()
	producerShowUuid := uuid.NewString()
	showRepository := repositoryShow.NewPostgresShowRepository(db)

	repository := NewPostgresUserRepository(db)
	username := "Some username"
	user := &model.User{
		Id:       uuid.NewString(),
		Username: username,
		ShowRoles: []model.ShowRole{
			{ShowId: editorShowUuid, Role: model.EDITOR},
			{ShowId: producerShowUuid, Role: model.PRODUCER},
		},
	}

	err := showRepository.SaveShow(&model.Show{Id: editorShowUuid, Title: "editor-show", Slug: "editor-slug"})
	assert.Nil(t, err)

	err = showRepository.SaveShow(&model.Show{Id: producerShowUuid, Title: "producer-show", Slug: "producer-slug"})
	assert.Nil(t, err)

	err = repository.SaveUser(user)
	assert.Nil(t, err)

	t.Run("should return nil if user does not exist", func(t *testing.T) {
		foundUser, err := repository.GetUserOrNil(uuid.NewString())
		assert.Nil(t, err)
		assert.Nil(t, foundUser)
	})

	t.Run("should retrieve a user", func(t *testing.T) {
		foundUser, err := repository.GetUserOrNil(user.Id)
		assert.Nil(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user.Id, foundUser.Id)
		assert.Equal(t, user.Username, foundUser.Username)
		sort.Sort(model.ByRole(user.ShowRoles))
		sort.Sort(model.ByRole(foundUser.ShowRoles))
		assert.EqualValues(t, user.ShowRoles, foundUser.ShowRoles)
	})
}

func Test_should_update_a_user(t *testing.T) {
	db := postgresTestSetup.StartTestcontainersPostgres(t, "../postgresTestSetup/")
	defer postgresTestSetup.Teardown(t, db)

	firstShowUuid := uuid.NewString()
	secondShowUuid := uuid.NewString()

	showRepository := repositoryShow.NewPostgresShowRepository(db)
	repository := NewPostgresUserRepository(db)

	username := "old name"
	user := &model.User{
		Id:       uuid.NewString(),
		Username: username,
		ShowRoles: []model.ShowRole{
			{ShowId: firstShowUuid, Role: model.PRODUCER},
		},
	}

	userUpdate := &model.User{
		Id:       user.Id,
		Username: "new username",
		ShowRoles: []model.ShowRole{
			{ShowId: firstShowUuid, Role: model.EDITOR},
			{ShowId: secondShowUuid, Role: model.PRODUCER},
		},
	}

	if err := showRepository.SaveShow(&model.Show{Id: firstShowUuid, Title: "first-show", Slug: "first-slug"}); err != nil {
		t.Fatal(err)
	}
	if err := showRepository.SaveShow(&model.Show{Id: secondShowUuid, Title: "second-show", Slug: "second-slug"}); err != nil {
		t.Fatal(err)
	}

	if err := repository.SaveUser(user); err != nil {
		t.Fatal(err)
	}

	t.Run("should update all fields", func(t *testing.T) {
		err := repository.UpdateUser(userUpdate)
		assert.Nil(t, err)
	})

	t.Run("should query updated user", func(t *testing.T) {
		if fetchedUser, err := repository.GetUserOrNil(user.Id); err != nil {
			t.Fatal(err)
		} else {
			sort.Sort(model.ByRole(fetchedUser.ShowRoles))
			sort.Sort(model.ByRole(userUpdate.ShowRoles))
			assert.Equal(t, userUpdate, fetchedUser)
		}
	})

	t.Run("should remove showRoles", func(t *testing.T) {
		userUpdate.ShowRoles = []model.ShowRole{}
		err := repository.UpdateUser(userUpdate)
		assert.Nil(t, err)

		if fetchedUser, err := repository.GetUserOrNil(user.Id); err != nil {
			t.Fatal(err)
		} else {
			assert.Empty(t, fetchedUser.ShowRoles)
		}
	})
}
