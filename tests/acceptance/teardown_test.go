//nolint:varnamelen
package readme_test

import (
	"log"
	"os"

	"github.com/liveoaklabs/readme-api-go-client/readme"
)

type teardownArgs struct {
	client  *readme.Client
	options readme.RequestOptions
}

func teardown() {
	token := os.Getenv("README_API_KEY")
	client, err := readme.NewClient(token)
	if err != nil {
		panic(err)
	}

	td := teardownArgs{
		client: client,
	}

	td.deleteSpecs(client)
	td.deleteCategories(client)
	td.deleteVersions(client)
	td.deleteCustomPages(client)
	td.deleteDocs(client)
}

func (t teardownArgs) deleteSpecs(client *readme.Client) {
	specs, resp, err := client.APISpecification.GetAll(t.options)
	if err != nil {
		log.Printf("Error getting specs (%d): %s", resp.HTTPResponse.StatusCode, err)
	}

	if len(specs) == 0 {
		return
	}

	for _, spec := range specs {
		log.Printf("Deleting spec: %s", spec.ID)
		_, resp, err := client.APISpecification.Delete(spec.ID)
		if err != nil {
			log.Printf("Error deleting spec (%d): %s", resp.HTTPResponse.StatusCode, err)
		}
	}
}

func (t teardownArgs) deleteCategories(client *readme.Client) {
	categories, resp, err := client.Category.GetAll(t.options)
	if err != nil {
		log.Printf("Error getting categories (%d): %s", resp.HTTPResponse.StatusCode, err)
	}

	if len(categories) == 0 {
		return
	}

	for _, category := range categories {
		log.Printf("Deleting category: %s", category.Slug)
		_, resp, err := client.Category.Delete(category.Slug, t.options)
		if err != nil {
			log.Printf("Error deleting category (%d): %s", resp.HTTPResponse.StatusCode, err)
		}
	}
}

func (t teardownArgs) deleteVersions(client *readme.Client) {
	versions, resp, err := client.Version.GetAll()
	if err != nil {
		log.Printf("Error getting versions (%d): %s", resp.HTTPResponse.StatusCode, err)
	}

	// Keep default 1.0.0
	for i, version := range versions {
		if version.Version == "1.0" || version.Version == "1.0.0" {
			versions = append(versions[:i], versions[i+1:]...)
		}
	}

	if len(versions) == 0 {
		return
	}

	for _, version := range versions {
		if version.Version == "1.0" || version.Version == "1.0.0" {
			continue
		}

		_, resp, err := client.Version.Delete(version.VersionClean)
		if err != nil {
			log.Printf("Error deleting version (%d): %s", resp.HTTPResponse.StatusCode, err)
		}
	}
}

func (t teardownArgs) deleteCustomPages(client *readme.Client) {
	pages, resp, err := client.CustomPage.GetAll(t.options)
	if err != nil {
		log.Printf("Error getting custom pages (%d): %s", resp.HTTPResponse.StatusCode, err)
	}

	if len(pages) == 0 {
		return
	}

	for _, page := range pages {
		_, resp, err := client.CustomPage.Delete(page.Slug)
		if err != nil {
			log.Printf("Error deleting custom page (%d): %s", resp.HTTPResponse.StatusCode, err)
		}
	}
}

func (t teardownArgs) deleteDocs(client *readme.Client) {
	docs, resp, err := client.Doc.Search("*", readme.RequestOptions{})
	if err != nil {
		log.Printf("Error getting docs (%d): %s", resp.HTTPResponse.StatusCode, err)
	}

	if len(docs) == 0 {
		return
	}

	for _, doc := range docs {
		_, resp, err := client.Doc.Delete(doc.Slug)
		if err != nil {
			log.Printf("Error deleting doc (%d): %s", resp.HTTPResponse.StatusCode, err)
		}
	}
}
