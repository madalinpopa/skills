package tests

import (
	"net/http"
	"testing"

	"{{MODULE_PATH}}/internal/modules/{{MODULE_NAME}}/adapters/rest/client"
	"{{MODULE_PATH}}/internal/testkit"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAndGet{{ENTITY_PASCAL}}(t *testing.T) {
	testkit.RequireIntegration(t)

	// Arrange
	api, err := client.NewClientWithResponses(serverURL)
	require.NoError(t, err)

	// Act
	created, err := api.Create{{ENTITY_PASCAL}}WithResponse(t.Context(), client.Create{{ENTITY_PASCAL}}JSONRequestBody{Name: "first"})
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, created.StatusCode())
	require.NotNil(t, created.JSON201)

	got, err := api.Get{{ENTITY_PASCAL}}WithResponse(t.Context(), created.JSON201.Uuid)
	require.NoError(t, err)

	// Assert
	require.Equal(t, http.StatusOK, got.StatusCode())
	require.NotNil(t, got.JSON200)
	assert.Equal(t, created.JSON201.Uuid, got.JSON200.Uuid)
	assert.Equal(t, "first", got.JSON200.Name)
}

func TestGet{{ENTITY_PASCAL}}_notFound(t *testing.T) {
	testkit.RequireIntegration(t)

	// Arrange
	api, err := client.NewClientWithResponses(serverURL)
	require.NoError(t, err)

	// Act
	got, err := api.Get{{ENTITY_PASCAL}}WithResponse(t.Context(), "does-not-exist")
	require.NoError(t, err)

	// Assert
	require.Equal(t, http.StatusNotFound, got.StatusCode())
	require.NotNil(t, got.JSON404)
	assert.Equal(t, "{{ENTITY}}-not-found", got.JSON404.Slug)
}

func TestCreate{{ENTITY_PASCAL}}_emptyName(t *testing.T) {
	testkit.RequireIntegration(t)

	// Arrange
	api, err := client.NewClientWithResponses(serverURL)
	require.NoError(t, err)

	// Act
	created, err := api.Create{{ENTITY_PASCAL}}WithResponse(t.Context(), client.Create{{ENTITY_PASCAL}}JSONRequestBody{Name: "  "})
	require.NoError(t, err)

	// Assert
	require.Equal(t, http.StatusBadRequest, created.StatusCode())
	require.NotNil(t, created.JSON400)
	assert.Equal(t, "{{ENTITY}}-name-empty", created.JSON400.Slug)
}
