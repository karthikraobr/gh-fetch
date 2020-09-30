package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v32/github"
	"github.com/karthikraobr/gh-fetch/internal/cache"
	"github.com/karthikraobr/gh-fetch/internal/gh"
	"github.com/karthikraobr/gh-fetch/internal/store"
)

type Handler struct {
	log    *log.Logger
	client gh.Fetcher
	cache  *cache.TTLCache
	store  store.DB
}

// New initializes the handler struct
func New(client gh.Fetcher, log *log.Logger, store store.DB, cache *cache.TTLCache) *Handler {
	return &Handler{
		log:    log,
		client: client,
		store:  store,
		cache:  cache,
	}
}

// SetUpRouter assigns the handlers to the URLs
func (h *Handler) SetUpRouter() *gin.Engine {
	r := gin.Default()
	r.Use(ErrorHandler())
	r.GET("/user/:username/repositories", h.HandleRepositories())
	r.GET("/user/:username/repository/:repository/commits", h.HandleCommits())
	r.GET("/user/:username/top20", h.HandleTop20())
	return r
}

// HandleRepositories fetches the public gh repositories
func (h *Handler) HandleRepositories() func(c *gin.Context) {
	return h.repoHandler
}

func (h *Handler) repoHandler(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	perPage, err := strconv.Atoi(c.DefaultQuery("perpage", "20"))
	if err != nil {
		perPage = 20
	}
	username := c.Param("username")
	if username == "" {
		c.Error(NewHttpError(http.StatusBadRequest, errors.New("empty username")))
		return
	}
	// Enable cache pagination. This would make cache management difficult. Due to lack of time this is being skipped rn.
	// Due to this results will be inconsistent and not follow page and perpage requests.
	cacheVal := h.cache.Get(username)
	if cacheVal != nil {
		if val, ok := cacheVal.([]*gh.Repository); ok {
			c.JSON(http.StatusOK, val)
			return
		}
	}
	opt := github.RepositoryListOptions{Type: "public", ListOptions: github.ListOptions{Page: page, PerPage: perPage}}
	repos, err := h.client.ListRepositories(c, username, &opt)
	if err != nil {
		repos, err := h.store.GetRepositories(username)
		if err != nil {
			c.Error(NewHttpError(http.StatusInternalServerError, err))
			return
		}
		c.JSON(http.StatusOK, repos)
		return
	}
	h.cache.Put(username, repos)
	c.JSON(http.StatusOK, repos)
	if _, err := h.store.CreateRepositories(repos); err != nil {
		h.log.Println("error in creating rows", err.Error())
		return
	}

}

//HandleCommits fetches the commits of a gh repository.
func (h *Handler) HandleCommits() func(c *gin.Context) {
	return h.commitHandler
}

func (h *Handler) commitHandler(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	perPage, err := strconv.Atoi(c.DefaultQuery("perpage", "20"))
	if err != nil {
		perPage = 20
	}
	username := c.Param("username")
	if username == "" {
		c.Error(NewHttpError(http.StatusBadRequest, errors.New("empty username")))
		return
	}
	repo := c.Param("repository")
	if repo == "" {
		c.Error(NewHttpError(http.StatusBadRequest, errors.New("empty repo name")))
		return
	}
	cKey := username + "/" + repo
	commits := h.cache.Get(cKey)
	if commits != nil {
		c.JSON(http.StatusOK, commits)
		return
	}
	opt := github.CommitsListOptions{ListOptions: github.ListOptions{Page: page, PerPage: perPage}}
	commits, err = h.client.ListCommits(c, username, repo, &opt)
	if err != nil {
		c.Error(NewHttpError(http.StatusInternalServerError, err))
		return
	}
	h.cache.Put(cKey, commits)
	c.JSON(http.StatusOK, commits)
}

//HandleTop20 fetches the top 20 recently accessed repositories.
func (h *Handler) HandleTop20() func(c *gin.Context) {
	return h.top20Handler
}

func (h *Handler) top20Handler(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.Error(NewHttpError(http.StatusBadRequest, errors.New("empty username")))
		return
	}

	repos, err := h.store.GetRepositoriesOrderedBy(username, 20, "desc", "last_access")
	if err != nil {
		c.Error(NewHttpError(http.StatusInternalServerError, err))
		return
	}
	c.JSON(http.StatusOK, repos)
	return
}
