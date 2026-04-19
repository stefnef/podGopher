package user

import (
	"podGopher/adapter/outbound/repository/postgres/postgresTestSetup"
	repositoryShow "podGopher/adapter/outbound/repository/postgres/show"
	"podGopher/core/domain/model"
	"podGopher/core/domain/role"
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

func Test_should_not_assign_user_if_show_does_not_exist(t *testing.T) {
	db := postgresTestSetup.StartTestcontainersPostgres(t, "../postgresTestSetup/")
	defer postgresTestSetup.Teardown(t, db)

	nonExistingShowId := uuid.NewString()

	repository := NewPostgresUserRepository(db)
	user := &model.User{
		Id:       uuid.NewString(),
		Username: "some-username",
		ShowRoles: []domainRole.ShowRole{
			{
				ShowId: nonExistingShowId,
				Role:   domainRole.EDITOR,
			},
		},
	}
	err := repository.SaveUser(user)
	assert.NotNil(t, err)
}

func Test_should_save_a_new_user(t *testing.T) {
	db := postgresTestSetup.StartTestcontainersPostgres(t, "../postgresTestSetup/")
	defer postgresTestSetup.Teardown(t, db)

	repository := NewPostgresUserRepository(db)
	username := "Some username"
	user := &model.User{
		Id:        uuid.NewString(),
		Username:  username,
		ShowRoles: []domainRole.ShowRole{},
	}

	t.Run("should return false if user with username does not exist", func(t *testing.T) {
		exists := repository.ExistsByUsername(username)
		assert.False(t, exists)
	})

	t.Run("should save a user", func(t *testing.T) {
		err := repository.SaveUser(user)
		assert.Nil(t, err)
	})

	t.Run("should return true if user with username exists", func(t *testing.T) {
		exists := repository.ExistsByUsername(username)
		assert.True(t, exists)
	})

	t.Run("should return false if user with username does not exists", func(t *testing.T) {
		exists := repository.ExistsByUsername("some-other-user")
		assert.False(t, exists)
	})

	t.Run("should query a user", func(t *testing.T) {
		var id, username string
		var isAdmin bool
		err := db.QueryRow("SELECT * FROM podcast_user WHERE id = $1", user.Id).
			Scan(&id, &username, &isAdmin)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, user.Id, id)
		assert.Equal(t, user.Username, username)
		assert.False(t, isAdmin)
	})

	t.Run("should save an admin user", func(t *testing.T) {
		err := repository.SaveUser(&model.User{
			Id:       uuid.NewString(),
			Username: "admin-user",
			IsAdmin:  true,
		})
		assert.Nil(t, err)
		var adminUser *model.User
		adminUser, err = repository.GetUserByUsernameOrNil("admin-user")
		assert.Nil(t, err)
		assert.True(t, adminUser.IsAdmin)
	})
}

func Test_should_assign_a_user(t *testing.T) {
	db := postgresTestSetup.StartTestcontainersPostgres(t, "../postgresTestSetup/")
	defer postgresTestSetup.Teardown(t, db)

	editorShowUuid := uuid.NewString()
	producerShowUuid := uuid.NewString()

	showRepository := repositoryShow.NewPostgresShowRepository(db)
	repository := NewPostgresUserRepository(db)
	username := "Some username"
	user := &model.User{
		Id:        uuid.NewString(),
		Username:  username,
		ShowRoles: nil,
	}

	if err := showRepository.SaveShow(&model.Show{Id: editorShowUuid, Title: "editor-show", Slug: "editor-slug"}); err != nil {
		t.Fatal(err)
	}
	if err := showRepository.SaveShow(&model.Show{Id: producerShowUuid, Title: "producer-show", Slug: "producer-slug"}); err != nil {
		t.Fatal(err)
	}

	if err := repository.SaveUser(user); err != nil {
		t.Fatal(err)
	}

	t.Run("should return false if user with username does not exist", func(t *testing.T) {
		exists := repository.ExistsByShowIdAndByUserId(producerShowUuid, user.Id)
		assert.False(t, exists)

		exists = repository.ExistsByShowIdAndByUserId(editorShowUuid, user.Id)
		assert.False(t, exists)
	})

	t.Run("should assign a user", func(t *testing.T) {
		user.ShowRoles = []domainRole.ShowRole{
			{ShowId: producerShowUuid, Role: domainRole.PRODUCER},
			{ShowId: editorShowUuid, Role: domainRole.EDITOR},
		}
		err := repository.UpdateUser(user)
		assert.Nil(t, err)
	})

	t.Run("should return true if user with id exists", func(t *testing.T) {
		exists := repository.ExistsByShowIdAndByUserId(producerShowUuid, user.Id)
		assert.True(t, exists)

		exists = repository.ExistsByShowIdAndByUserId(editorShowUuid, user.Id)
		assert.True(t, exists)
	})

	t.Run("should return false if user with user id does not exists", func(t *testing.T) {
		exists := repository.ExistsByShowIdAndByUserId(editorShowUuid, "some-other-user")
		assert.False(t, exists)

		exists = repository.ExistsByShowIdAndByUserId("i-do-not-exist", user.Id)
		assert.False(t, exists)
	})

	t.Run("should query show roles", func(t *testing.T) {
		var role string
		err := db.QueryRow("SELECT role FROM user_roles WHERE user_id = $1 and show_id = $2", user.Id, editorShowUuid).
			Scan(&role)
		assert.Nil(t, err)
		assert.Equal(t, domainRole.EDITOR.Name(), role)

		err = db.QueryRow("SELECT role FROM user_roles WHERE user_id = $1 and show_id = $2", user.Id, producerShowUuid).
			Scan(&role)
		assert.Nil(t, err)
		assert.Equal(t, domainRole.PRODUCER.Name(), role)
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
		IsAdmin:  false,
		ShowRoles: []domainRole.ShowRole{
			{ShowId: editorShowUuid, Role: domainRole.EDITOR},
			{ShowId: producerShowUuid, Role: domainRole.PRODUCER},
		},
	}

	err := showRepository.SaveShow(&model.Show{Id: editorShowUuid, Title: "editor-show", Slug: "editor-slug"})
	assert.Nil(t, err)

	err = showRepository.SaveShow(&model.Show{Id: producerShowUuid, Title: "producer-show", Slug: "producer-slug"})
	assert.Nil(t, err)

	err = repository.SaveUser(user)
	assert.Nil(t, err)

	t.Run("should return nil if user does not exist", func(t *testing.T) {
		foundUser, err := repository.GetUserByIdOrNil(uuid.NewString())
		assert.Nil(t, err)
		assert.Nil(t, foundUser)
	})

	t.Run("should return nil if user does not exist on getByUsername", func(t *testing.T) {
		foundUser, err := repository.GetUserByUsernameOrNil(uuid.NewString())
		assert.Nil(t, err)
		assert.Nil(t, foundUser)
	})

	t.Run("should retrieve a user", func(t *testing.T) {
		foundUser, err := repository.GetUserByIdOrNil(user.Id)
		assert.Nil(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user.Id, foundUser.Id)
		assert.Equal(t, user.Username, foundUser.Username)
		assert.False(t, foundUser.IsAdmin)
		sort.Sort(domainRole.ByRole(user.ShowRoles))
		sort.Sort(domainRole.ByRole(foundUser.ShowRoles))
		assert.EqualValues(t, user.ShowRoles, foundUser.ShowRoles)
	})

	t.Run("should retrieve a user by username", func(t *testing.T) {
		foundUser, err := repository.GetUserByUsernameOrNil(user.Username)
		assert.Nil(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user.Id, foundUser.Id)
	})
}

func Test_should_update_a_user_assignment(t *testing.T) {
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
		ShowRoles: []domainRole.ShowRole{
			{ShowId: firstShowUuid, Role: domainRole.PRODUCER},
		},
	}

	userUpdate := &model.User{
		Id:       user.Id,
		Username: "new username",
		ShowRoles: []domainRole.ShowRole{
			{ShowId: firstShowUuid, Role: domainRole.EDITOR},
			{ShowId: secondShowUuid, Role: domainRole.PRODUCER},
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
		if fetchedUser, err := repository.GetUserByIdOrNil(user.Id); err != nil {
			t.Fatal(err)
		} else {
			sort.Sort(domainRole.ByRole(fetchedUser.ShowRoles))
			sort.Sort(domainRole.ByRole(userUpdate.ShowRoles))
			assert.Equal(t, userUpdate, fetchedUser)
		}
	})

	t.Run("should remove showRoles", func(t *testing.T) {
		userUpdate.ShowRoles = []domainRole.ShowRole{}
		err := repository.UpdateUser(userUpdate)
		assert.Nil(t, err)

		if fetchedUser, err := repository.GetUserByIdOrNil(user.Id); err != nil {
			t.Fatal(err)
		} else {
			assert.Empty(t, fetchedUser.ShowRoles)
		}
	})
}
