package contentful

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleSpacesService_Get() {
	cma := NewCMA("cmaClient-token")

	space, err := cma.Spaces.Get("space-id")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(space.Name)
}

func ExampleSpacesService_List() {
	cma := NewCMA("cmaClient-token")
	collection, err := cma.Spaces.List().Next()
	if err != nil {
		log.Fatal(err)
	}

	spaces := collection.ToSpace()
	for _, space := range spaces {
		fmt.Println(space.Sys.ID, space.Name)
	}
}

func ExampleSpacesService_Upsert_create() {
	cma := NewCMA("cmaClient-token")

	space := &Space{
		Name:          "space-name",
		DefaultLocale: "en-US",
	}

	err := cma.Spaces.Upsert(space)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleSpacesService_Upsert_update() {
	cma := NewCMA("cmaClient-token")

	space, err := cma.Spaces.Get("space-id")
	if err != nil {
		log.Fatal(err)
	}

	space.Name = "modified"
	err = cma.Spaces.Upsert(space)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleSpacesService_Delete() {
	cma := NewCMA("cmaClient-token")

	space, err := cma.Spaces.Get("space-id")
	if err != nil {
		log.Fatal(err)
	}

	err = cma.Spaces.Delete(space)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleSpacesService_Delete_all() {
	cma := NewCMA("cmaClient-token")

	collection, err := cma.Spaces.List().Next()
	if err != nil {
		log.Fatal(err)
	}

	for _, space := range collection.ToSpace() {
		err := cma.Spaces.Delete(space)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestSpacesServiceList(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("spaces.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cmaClient client
	cmaClient = NewCMA(CMAToken)
	cmaClient.BaseURL = server.URL

	collection, err := cmaClient.Spaces.List().Next()
	assertions.Nil(err)

	spaces := collection.ToSpace()
	assertions.Equal(2, len(spaces))
	assertions.Equal("id1", spaces[0].Sys.ID)
	assertions.Equal("id2", spaces[1].Sys.ID)
}

func TestSpacesServiceList_Pagination(t *testing.T) {
	var err error
	assertions := assert.New(t)

	requestCount := 1
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces")
		checkHeaders(r, assertions)

		w.WriteHeader(200)
		query := r.URL.Query()
		if requestCount == 1 {
			assertions.Equal(query.Get("order"), "-sys.createdAt")
			assertions.Equal(query.Get("skip"), "")
			_, _ = fmt.Fprintln(w, readTestData("spaces.json"))
		} else {
			assertions.Equal(query.Get("order"), "-sys.createdAt")
			assertions.Equal(query.Get("skip"), "100")
			_, _ = fmt.Fprintln(w, readTestData("spaces-page-2.json"))
		}
		requestCount++
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cmaClient client
	cmaClient = NewCMA(CMAToken)
	cmaClient.BaseURL = server.URL

	collection, err := cmaClient.Spaces.List().Next()
	assertions.Nil(err)

	nextPage, err := collection.Next()
	assertions.Nil(err)
	assertions.IsType(&Collection{}, nextPage)
}

func TestSpacesServiceGet(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID)

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("space-1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cmaClient client
	cmaClient = NewCMA(CMAToken)
	cmaClient.BaseURL = server.URL

	space, err := cmaClient.Spaces.Get(spaceID)
	assertions.Nil(err)
	assertions.Equal("id1", space.Sys.ID)
}

func TestSpacesService_Get_2(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID)

		checkHeaders(r, assertions)

		w.WriteHeader(400)
		_, _ = fmt.Fprintln(w, readTestData("space-1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cmaClient client
	cmaClient = NewCMA(CMAToken)
	cmaClient.BaseURL = server.URL

	_, err = cmaClient.Spaces.Get(spaceID)
	assertions.NotNil(err)
}

func TestSpaceSaveForCreate(t *testing.T) {
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "POST")
		assertions.Equal(r.RequestURI, "/spaces")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("new space", payload["name"])
		assertions.Equal("en", payload["defaultLocale"])

		w.WriteHeader(201)
		_, _ = fmt.Fprintln(w, readTestData("spaces-newspace.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cmaClient client
	cmaClient = NewCMA(CMAToken)
	cmaClient.BaseURL = server.URL

	space := &Space{
		Name:          "new space",
		DefaultLocale: "en",
	}

	err := cmaClient.Spaces.Upsert(space)
	assertions.Nil(err)
	assertions.Equal("newspace", space.Sys.ID)
	assertions.Equal("new space", space.Name)
	assertions.Equal("en", space.DefaultLocale)
}

func TestSpaceSaveForUpdate(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/newspace")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("changed-space-name", payload["name"])
		assertions.Equal("de", payload["defaultLocale"])

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("spaces-newspace-updated.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cmaClient client
	cmaClient = NewCMA(CMAToken)
	cmaClient.BaseURL = server.URL

	space, err := spaceFromTestData("spaces-newspace.json")
	assertions.Nil(err)

	space.Name = "changed-space-name"
	space.DefaultLocale = "de"

	err = cmaClient.Spaces.Upsert(space)
	assertions.Nil(err)
	assertions.Equal("changed-space-name", space.Name)
	assertions.Equal("de", space.DefaultLocale)
	assertions.Equal(2, space.Sys.Version)
}

func TestSpaceDelete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "DELETE")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID)
		checkHeaders(r, assertions)

		w.WriteHeader(200)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cmaClient client
	cmaClient = NewCMA(CMAToken)
	cmaClient.BaseURL = server.URL

	space, err := spaceFromTestData("spaces-" + spaceID + ".json")
	assertions.Nil(err)

	err = cmaClient.Spaces.Delete(space)
	assertions.Nil(err)
}
