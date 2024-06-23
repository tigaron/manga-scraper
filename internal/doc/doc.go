package doc

import (
	"fmt"

	"goa.design/model/dsl"
	expr "goa.design/model/expr"
)

var _ = dsl.Design("Manga Scraper API", "Go microservice project using Domain Driven Design!", func() {
	var (
		System        *expr.SoftwareSystem
		ElasticSearch *expr.Container
		Kafka         *expr.Container
		MySQL         *expr.Container
		Redis         *expr.Container
		RESTfulAPI    *expr.Container
	)

	const (
		containerRESTfulAPI = "RESTful API"
		softwareSystem      = "Manga Scraper System"
		goBinary            = "Go Binary"
	)

	const (
		styleSoftwareSystem = "Software System"
		styleComponent      = "Component"
		styleContainer      = "Container"
		styleDatabase       = "Database"
		stylePerson         = "Person"
	)

	const (
		componentElasticsearch = "internal.elasticsearch"
		componentKafka         = "internal.kafka"
		componentMySQL         = "internal.database.prisma"
		componentRedis         = "internal.database.redis"
	)

	dsl.Person("User", "Interacts with Service", func() {
		dsl.External()

		dsl.Tag(stylePerson)

		dsl.Uses(softwareSystem, "Reads and writes data using", "HTTPS/JSON", dsl.Synchronous, func() {
			dsl.Tag("Relationship", "Synchronous")
		})

		dsl.Uses(fmt.Sprintf("%s/%s", softwareSystem, containerRESTfulAPI), "Reads and writes data using", "HTTPS/JSON", dsl.Synchronous, func() {
			dsl.Tag("Relationship", "Synchronous")
		})
	})

	System = dsl.SoftwareSystem(softwareSystem, "Allows users to interact with Manga data", func() {
		dsl.URL("https://github.com/tigaron/manga-scraper")

		MySQL = dsl.Container("MySQL", "Stores Manga records", "MySQL 8.x", func() {
			dsl.Tag(styleDatabase)
			dsl.Tag(styleContainer)
		})

		Redis = dsl.Container("Redis", "Stores cached Manga records", "Redis 6.x", func() {
			dsl.Tag(styleDatabase)
			dsl.Tag(styleContainer)
		})

		ElasticSearch = dsl.Container("ElasticSearch", "Stores searchable Manga records", "OpenSearch 2.x", func() {
			dsl.Tag(styleDatabase)
			dsl.Tag(styleContainer)
		})

		Kafka = dsl.Container("Kafka", "Streams Scrape Request events", "Kafka 2.x", func() {
			dsl.Tag(styleDatabase)
			dsl.Tag(styleContainer)
		})

		RESTfulAPI = dsl.Container(containerRESTfulAPI, "RESTful API", "Go 1.22", func() {
			dsl.Uses(Redis, "Reads from and Writes to", "Redis", dsl.Synchronous, func() {})
			dsl.Uses(MySQL, "Reads from and Writes to", "MySQL", dsl.Synchronous, func() {})
			dsl.Uses(ElasticSearch, "Reads from", "HTTPS", dsl.Synchronous, func() {})
			dsl.Uses(Kafka, "Produces", "Kafka", dsl.Asynchronous, func() {})

			dsl.Component(componentElasticsearch, "interacts with ElasticSearch", "Go Package", func() {
				dsl.Uses(ElasticSearch, "Uses", dsl.Synchronous, func() {
					dsl.Tag("Relationship", "Synchronous")
				})

				dsl.Tag(styleComponent)
			})

			dsl.Component(componentMySQL, "interacts with MySQL", "Go Package", func() {
				dsl.Uses(MySQL, "Uses", dsl.Synchronous, func() {
					dsl.Tag("Relationship", "Synchronous")
				})

				dsl.Tag(styleComponent)
			})

			dsl.Component(componentRedis, "interacts with Redis", "Go Package", func() {
				dsl.Uses(Redis, "Uses", dsl.Synchronous, func() {
					dsl.Tag("Relationship", "Synchronous")
				})

				dsl.Tag(styleComponent)
			})

			dsl.Component(componentKafka, "interacts with Kafka", "Go Package", func() {
				dsl.Uses(Kafka, "Uses", dsl.Asynchronous, func() {
					dsl.Tag("Relationship", "Asynchronous")
				})

				dsl.Tag(styleComponent)
			})

			dsl.Component("internal.service", "interacts with all datastores", "Go Package", func() {
				dsl.Uses(componentElasticsearch, "Reads records from", dsl.Synchronous, func() {
					dsl.Tag("Relationship", "Synchronous")
				})

				dsl.Uses(componentKafka, "Produce events to", dsl.Synchronous, func() {
					dsl.Tag("Relationship", "Synchronous")
				})

				dsl.Uses(componentMySQL, "Uses", dsl.Synchronous, func() {
					dsl.Tag("Relationship", "Synchronous")
				})

				dsl.Uses(componentRedis, "Uses", dsl.Synchronous, func() {
					dsl.Tag("Relationship", "Synchronous")
				})

				dsl.Tag(styleComponent)
			})

			dsl.Component("internal.rest", "defines HTTP handlers", "Go Package", func() {
				dsl.Uses("internal.service", "Uses", dsl.Synchronous, func() {
					dsl.Tag("Relationship", "Synchronous")
				})

				dsl.Tag(styleComponent)
			})

			dsl.Tag(styleContainer)
		})

		dsl.Tag(styleSoftwareSystem)
	})

	dsl.Views(func() {
		dsl.SystemContextView(System, "Manga Scraper System", func() {
			dsl.AddDefault()

			dsl.EnterpriseBoundaryVisible()
		})

		dsl.ContainerView(softwareSystem, "Containers", "Container diagram for the Manga Scraper System", func() {
			dsl.AddDefault()

			dsl.SystemBoundariesVisible()
		})

		dsl.ComponentView(RESTfulAPI, "RESTful API", "Component diagram for the REST Server", func() {
			dsl.AddDefault()

			dsl.ContainerBoundariesVisible()
		})

		dsl.Styles(func() {
			dsl.ElementStyle(styleSoftwareSystem, func() {
				dsl.Background("#1168bd")
				dsl.Color("#ffffff")
			})

			dsl.ElementStyle(stylePerson, func() {
				dsl.Background("#08427b")
				dsl.Color("#ffffff")
				dsl.Shape(dsl.ShapePerson)
			})

			dsl.ElementStyle(styleComponent, func() {
				dsl.Background("#85bbf0")
				dsl.Color("#000000")
			})

			dsl.ElementStyle(styleContainer, func() {
				dsl.Background("#438dd5")
				dsl.Color("#ffffff")
			})

			dsl.ElementStyle(styleDatabase, func() {
				dsl.Shape(dsl.ShapeCylinder)
			})
		})
	})
})
