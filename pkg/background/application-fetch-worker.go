package background

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	log "github.com/sirupsen/logrus"
	"github.com/wufe/polo/pkg/models"
	"github.com/wufe/polo/pkg/storage"
	"github.com/wufe/polo/pkg/versioning"
)

type ApplicationFetchWorker struct {
	sessionStorage *storage.Session
	mediator       *Mediator
}

func NewApplicationFetchWorker(sessionStorage *storage.Session, mediator *Mediator) *ApplicationFetchWorker {
	worker := &ApplicationFetchWorker{
		sessionStorage: sessionStorage,
		mediator:       mediator,
	}

	worker.startAcceptingFetchRequests()

	return worker
}

func (w *ApplicationFetchWorker) startAcceptingFetchRequests() {
	go func() {
		for {
			applicationFetchReq := <-w.mediator.ApplicationFetch.RequestChan
			w.FetchApplicationRemote(applicationFetchReq.Application, applicationFetchReq.WatchObjects)
			w.mediator.ApplicationFetch.ResponseChan <- nil
		}
	}()
}

func (w *ApplicationFetchWorker) FetchApplicationRemote(application *models.Application, watchObjects bool) {

	// TODO: Handle all these errors

	auth, err := application.GetAuth()
	if err != nil {
		return
	}

	var appName string
	var baseFolder string
	var watch models.Watch

	application.WithRLock(func(a *models.Application) {
		appName = a.Name
		baseFolder = a.BaseFolder
		watch = a.Watch
	})

	objectsToHashMap := make(map[string]string)
	hashToObjectsMap := make(map[string]*models.RemoteObject)
	appBranches := make(map[string]*models.Branch)
	appTags := []string{}
	appCommits := []string{}
	appCommitMap := make(map[string]*object.Commit)

	checkObjectExists := func(hashToObjectsMap map[string]*models.RemoteObject) func(hash string) {
		return func(hash string) {
			if _, exists := hashToObjectsMap[hash]; !exists {
				hashToObjectsMap[hash] = &models.RemoteObject{
					Branches: []string{},
					Tags:     []string{},
				}
			}
		}
	}(hashToObjectsMap)

	registerHash, watchResults := w.registerObjectHash(objectsToHashMap, watch)

	gitClient := versioning.GetGitClient(application, auth)

	// Open repository
	repo, err := git.PlainOpen(baseFolder)
	defaultApplicationErrorLog(application, err)
	if err != nil {
		return
	}

	// Fetch
	err = gitClient.FetchAll(baseFolder)
	defaultApplicationErrorLog(application, err, git.NoErrAlreadyUpToDate)

	// Branches
	branches, err := repo.Branches()
	defaultApplicationErrorLog(application, err)
	refPrefix := "refs/heads/"
	branches.ForEach(func(ref *plumbing.Reference) error {
		refName := ref.Name().String()
		refHash := ref.Hash().String()

		if !strings.HasPrefix(refName, refPrefix) {
			return nil
		}
		branchName := refName[len(refPrefix):]

		registerHash(branchName, refHash)
		registerHash(fmt.Sprintf("origin/%s", branchName), refHash)
		registerHash(ref.Name().String(), refHash)

		checkObjectExists(refHash)

		hashToObjectsMap[refHash].Branches = appendWithoutDup(hashToObjectsMap[refHash].Branches, branchName)

		commit, err := repo.CommitObject(ref.Hash())
		if err != nil {
			return err
		}

		appBranches[branchName] = &models.Branch{
			Name:    branchName,
			Hash:    refHash,
			Author:  commit.Author.Email,
			Date:    commit.Author.When,
			Message: commit.Message,
		}

		return nil
	})
	defaultApplicationErrorLog(application, err)

	// Tags
	tags, err := repo.Tags()
	if err != nil {
		return
	}
	defaultApplicationErrorLog(application, err)

	tagPrefix := "refs/tags/"
	err = tags.ForEach(func(ref *plumbing.Reference) error {
		refName := ref.Name().String()
		refHash := ref.Hash().String()

		if !strings.HasPrefix(refName, tagPrefix) {
			return nil
		}

		tagName := refName[len(tagPrefix):]
		registerHash(tagName, refHash)

		appTags = appendWithoutDup(appTags, tagName)
		registerHash(refName, refHash)
		checkObjectExists(refHash)

		hashToObjectsMap[refHash].Tags = appendWithoutDup(hashToObjectsMap[refHash].Tags, tagName)

		return nil
	})

	// Log
	// TODO: Configure "since"
	since := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Now().UTC()
	logs, err := repo.Log(&git.LogOptions{All: true, Since: &since, Until: &until, Order: git.LogOrderCommitterTime})
	defaultApplicationErrorLog(application, err)
	if err != nil {
		return
	}

	err = logs.ForEach(func(commit *object.Commit) error {
		commitHash := commit.Hash.String()
		registerHash(commitHash, commitHash)
		appCommits = append(appCommits, commitHash)
		appCommitMap[commitHash] = commit
		return nil
	})

	log.Infof("[APP:%s] Found %d commits", appName, len(appCommits))

	application.WithLock(func(a *models.Application) {
		a.ObjectsToHashMap = objectsToHashMap
		a.HashToObjectsMap = hashToObjectsMap
		a.Branches = appBranches
		a.Tags = appTags
		a.Commits = appCommits
		a.CommitMap = appCommitMap
	})

	if !watchObjects {
		return
	}

	for ref, hash := range *watchResults {
		sessions := w.sessionStorage.GetAllAliveSessions()
		var foundSession *models.Session
		for _, session := range sessions {

			sessionAppName := session.ApplicationName
			sessionCheckout := session.Checkout

			if sessionAppName == appName && (sessionCheckout == ref) {
				foundSession = session
			}
		}
		buildSession := requestSessionBuilder(application, ref)
		if foundSession != nil {
			sessionCommitID := foundSession.CommitID
			if sessionCommitID != hash {
				log.Infof("[APP:%s][WATCH] Detected new commit on %s", appName, ref)
				w.mediator.DestroySession.Enqueue(foundSession, func(s *models.Session) {
					buildSession(w.mediator)
				})
			}
		} else {
			log.Infof("[APP:%s][WATCH] Auto-start on %s", appName, ref)
			buildSession(w.mediator)
		}
	}
}

func (w *ApplicationFetchWorker) registerObjectHash(objectsToHashMap map[string]string, watch models.Watch) (func(refName string, hash string), *map[string]string) {
	watchResults := make(map[string]string)
	watchedHashes := make(map[string]bool)
	return func(refName, hash string) {
		objectsToHashMap[refName] = hash
		if watch.Contains(refName) {
			if _, ok := watchedHashes[hash]; !ok {
				watchResults[refName] = hash
				watchedHashes[hash] = true
			}
		}
	}, &watchResults
}

func requestSessionBuilder(a *models.Application, ref string) func(*Mediator) {
	return func(mediator *Mediator) {
		mediator.BuildSession.Enqueue(ref, a, nil)
	}
}
