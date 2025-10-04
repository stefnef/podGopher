package repository

import (
	"podGopher/core/port/outbound"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_should_implement_port(t *testing.T) {
	repository := NewInMemoryShowRepository()

	assert.NotNil(t, repository)
	assert.Implements(t, (*outbound.SaveShowPort)(nil), repository)
}

func Test_should_save_a_show(t *testing.T) {
	repository := NewInMemoryShowRepository()

	err := repository.SaveShow("Some title")
	assert.Nil(t, err)
}

func Test_should_return_false_if_show_does_not_exist(t *testing.T) {
	repository := NewInMemoryShowRepository()

	exists := repository.ExistsByTitle("Some title")
	assert.False(t, exists)
}

func Test_should_return_true_if_show_exists(t *testing.T) {
	repository := NewInMemoryShowRepository()

	_ = repository.SaveShow("Some title")
	exists := repository.ExistsByTitle("Some title")

	assert.True(t, exists)
}
