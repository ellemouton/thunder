package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ellemouton/thunder/db"
)

func TestCreate(t *testing.T) {
	dbc := db.ConnectForTesting(t)
	ctx := context.Background()

	_, err := Create(ctx, dbc, "title", "the dash", "blah blah blah.....")
	require.NoError(t, err)
}

func TestLookupInfo(t *testing.T) {
	dbc := db.ConnectForTesting(t)
	ctx := context.Background()

	id, err := Create(ctx, dbc, "title", "the dash", "blah blah blah.....")
	require.NoError(t, err)

	info, err := LookupInfo(ctx, dbc, id)
	require.NoError(t, err)

	require.Equal(t, info.Name, "title")
	require.Equal(t, info.Description, "the dash")
}

func TestLookupContent(t *testing.T) {
	dbc := db.ConnectForTesting(t)
	ctx := context.Background()

	id, err := Create(ctx, dbc, "title", "the dash", "blah blah blah")
	require.NoError(t, err)

	info, err := LookupInfo(ctx, dbc, id)
	require.NoError(t, err)

	content, err := LookupContent(ctx, dbc, info.ContentID)
	require.NoError(t, err)
	require.Equal(t, content.Text, "blah blah blah")
}

func TestListAllInfos(t *testing.T) {
	dbc := db.ConnectForTesting(t)
	ctx := context.Background()

	id1, err := Create(ctx, dbc, "title", "the dash", "blah blah blah")
	require.NoError(t, err)

	id2, err := Create(ctx, dbc, "title", "the dash", "blah blah blah")
	require.NoError(t, err)

	id3, err := Create(ctx, dbc, "title", "the dash", "blah blah blah")
	require.NoError(t, err)

	infos, err := ListAllInfoRev(ctx, dbc)
	require.NoError(t, err)
	require.Equal(t, len(infos), 3)
	require.Equal(t, id3, infos[0].ID)
	require.Equal(t, id2, infos[1].ID)
	require.Equal(t, id1, infos[2].ID)
}

func TestUpdateBlog(t *testing.T) {
	dbc := db.ConnectForTesting(t)
	ctx := context.Background()

	id, err := Create(ctx, dbc, "title", "the dash", "blah blah blah")
	require.NoError(t, err)

	err = UpdateBlog(ctx, dbc, id, "change", "meh", "mwahahahah")
	require.NoError(t, err)

	info, err := LookupInfo(ctx, dbc, id)
	require.NoError(t, err)
	require.Equal(t, info.Name, "change")
	require.Equal(t, info.Description, "meh")

	content, err := LookupContent(ctx, dbc, info.ContentID)
	require.NoError(t, err)
	require.Equal(t, content.Text, "mwahahahah")
}
