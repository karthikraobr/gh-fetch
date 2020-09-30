package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/karthikraobr/gh-fetch/internal/cache"
	"github.com/karthikraobr/gh-fetch/internal/gh"
	"github.com/karthikraobr/gh-fetch/internal/mock"
)

func TestHandler_repoHandler(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := []*gh.Repository{{
			ID:        1,
			CreatedAt: time.Now(),
			Name:      "blog",
			NodeID:    "1",
			Owner:     "me"}}
		fakeGh := mock.NewMockFetcher(ctrl)
		fakeGh.EXPECT().ListRepositories(gomock.Any(), gomock.Any(), gomock.Any()).Return(repo, nil)
		fakeStore := mock.NewMockDB(ctrl)
		fakeStore.EXPECT().CreateRepositories(gomock.Any()).Return(nil, nil)
		fakeHandler := New(fakeGh, &log.Logger{}, fakeStore, cache.New(1, 1))
		router := fakeHandler.SetUpRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/user/karthikraobr/repositories", nil)
		router.ServeHTTP(w, req)

		var result []*gh.Repository
		json.NewDecoder(w.Body).Decode(&result)
		if !(cmp.Equal(200, w.Code) && cmp.Equal(repo[0], result[0])) {
			t.Error("ok")
			t.Errorf("Code-want:%vgot:%v\n Result-want:%vgot:%v", 200, w.Code, repo[0], result[0])
		}
	})

	t.Run("missing-username", func(t *testing.T) {
		wantErr := "empty username"
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		fakeGh := mock.NewMockFetcher(ctrl)
		fakeStore := mock.NewMockDB(ctrl)
		fakeHandler := New(fakeGh, &log.Logger{}, fakeStore, cache.New(1, 1))
		router := fakeHandler.SetUpRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/user//repositories", nil)
		router.ServeHTTP(w, req)
		err := w.Body.String()
		if !(cmp.Equal(400, w.Code) && strings.Contains(err, wantErr)) {
			t.Error("missing-username")
			t.Errorf("Code-want:%vgot:%v\n Result-want:%v got:%v", 400, w.Code, wantErr, err)
		}
	})

	t.Run("gh-error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := []*gh.Repository{{
			ID:        1,
			CreatedAt: time.Now(),
			Name:      "blog",
			NodeID:    "1",
			Owner:     "me"}}
		fakeGh := mock.NewMockFetcher(ctrl)
		fakeGh.EXPECT().ListRepositories(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("network issue"))
		fakeStore := mock.NewMockDB(ctrl)
		fakeStore.EXPECT().GetRepositories(gomock.Any()).Return(repo, nil)
		fakeHandler := New(fakeGh, &log.Logger{}, fakeStore, cache.New(1, 1))
		router := fakeHandler.SetUpRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/user/karthikraobr/repositories", nil)
		router.ServeHTTP(w, req)
		var result []*gh.Repository
		json.NewDecoder(w.Body).Decode(&result)
		if !(cmp.Equal(200, w.Code) && cmp.Equal(repo[0], result[0])) {
			t.Errorf("gh-error failed")
			t.Errorf("Code-want:%vgot:%v\n Result-want:%v got:%v", 200, w.Code, repo[0], result[0])
		}
	})

	t.Run("db-get-error", func(t *testing.T) {
		dbErr := "db get error"
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		fakeGh := mock.NewMockFetcher(ctrl)
		fakeGh.EXPECT().ListRepositories(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("network issue"))
		fakeStore := mock.NewMockDB(ctrl)
		fakeStore.EXPECT().GetRepositories(gomock.Any()).Return(nil, errors.New(dbErr))
		fakeHandler := New(fakeGh, &log.Logger{}, fakeStore, cache.New(1, 1))
		router := fakeHandler.SetUpRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/user/karthikraobr/repositories", nil)
		router.ServeHTTP(w, req)
		err := w.Body.String()
		if !(cmp.Equal(500, w.Code) && strings.Contains(err, dbErr)) {
			t.Error("db-get-error failed")
			t.Errorf("Code-want:%vgot:%v\n Result-want:%v got:%v", 500, w.Code, dbErr, err)
		}
	})

}

func TestHandler_commitHandler(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		commits := []*gh.Commit{{
			Author:      "author",
			CommentsURL: "commentURL",
			NodeID:      "nodeid",
			SHA:         "sha",
		}}
		fakeGh := mock.NewMockFetcher(ctrl)
		fakeGh.EXPECT().ListCommits(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(commits, nil)
		fakeStore := mock.NewMockDB(ctrl)
		fakeHandler := New(fakeGh, &log.Logger{}, fakeStore, cache.New(1, 1))
		router := fakeHandler.SetUpRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/user/karthikraobr/repository/myrepo/commits", nil)
		router.ServeHTTP(w, req)

		var result []*gh.Commit
		json.NewDecoder(w.Body).Decode(&result)
		if !(cmp.Equal(200, w.Code) && cmp.Equal(commits[0], result[0])) {
			t.Error("ok")
			t.Errorf("Code-want:%vgot:%v\n Result-want:%v got:%v", 200, w.Code, commits[0], result[0])
		}
	})

	t.Run("missing-username", func(t *testing.T) {
		wantErr := "empty username"
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		fakeGh := mock.NewMockFetcher(ctrl)
		fakeStore := mock.NewMockDB(ctrl)
		fakeHandler := New(fakeGh, &log.Logger{}, fakeStore, cache.New(1, 1))
		router := fakeHandler.SetUpRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/user//repository/myrepo/commits", nil)
		router.ServeHTTP(w, req)
		err := w.Body.String()
		if !(cmp.Equal(400, w.Code) && strings.Contains(err, wantErr)) {
			t.Errorf("missing-username failed")
			t.Errorf("Code-want:%vgot:%v\n Result-want:%v got:%v", 400, w.Code, wantErr, err)
		}
	})

	t.Run("missing-reponame", func(t *testing.T) {
		wantErr := "empty repo name"
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		fakeGh := mock.NewMockFetcher(ctrl)
		fakeStore := mock.NewMockDB(ctrl)
		fakeHandler := New(fakeGh, &log.Logger{}, fakeStore, cache.New(1, 1))
		router := fakeHandler.SetUpRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/user/karthikraobr/repository//commits", nil)
		router.ServeHTTP(w, req)
		err := w.Body.String()
		if !(cmp.Equal(400, w.Code) && strings.Contains(err, wantErr)) {
			t.Errorf("ok failed")
			t.Errorf("Code-want:%vgot:%v\n Result-want:%v got:%v", 401, w.Code, wantErr, err)
		}
	})

}
