package doc

import (
	"fmt"

	. "goa.design/model/dsl"
	"goa.design/model/expr"
)

var _ = Design("Manga Scraper API", "Go microservice project using Domain Driven Design!", func() {
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

	Person("User", "Interacts with Service", func() {
		External()

		Tag(stylePerson)

		Uses(softwareSystem, "Reads and writes data using", "HTTPS/JSON", Synchronous, func() {
			Tag("Relationship", "Synchronous")
		})

		Uses(fmt.Sprintf("%s/%s", softwareSystem, containerRESTfulAPI), "Reads and writes data using", "HTTPS/JSON", Synchronous, func() {
			Tag("Relationship", "Synchronous")
		})
	})

	System = SoftwareSystem(softwareSystem, "Allows users to interact with Manga data", func() {
		URL("https://github.com/tigaron/manga-scraper")

		MySQL = Container("MySQL", "Stores Manga records", "MySQL 8.x", func() {
			Tag(styleDatabase)
			Tag(styleContainer)
		})

		Redis = Container("Redis", "Stores cached Manga records", "Redis 6.x", func() {
			Tag(styleDatabase)
			Tag(styleContainer)
		})

		ElasticSearch = Container("ElasticSearch", "Stores searchable Manga records", "OpenSearch 2.x", func() {
			Tag(styleDatabase)
			Tag(styleContainer)
		})

		Kafka = Container("Kafka", "Streams Scrape Request events", "Kafka 2.x", func() {
			Tag(styleDatabase)
			Tag(styleContainer)
		})

		RESTfulAPI = Container(containerRESTfulAPI, "RESTful API", "Go 1.22", func() {
			Uses(Redis, "Reads from and Writes to", "Redis", Synchronous, func() {})
			Uses(MySQL, "Reads from and Writes to", "MySQL", Synchronous, func() {})
			Uses(ElasticSearch, "Reads from", "HTTPS", Synchronous, func() {})
			Uses(Kafka, "Produces", "Kafka", Asynchronous, func() {})

			Component(componentElasticsearch, "interacts with ElasticSearch", "Go Package", func() {
				Uses(ElasticSearch, "Uses", Synchronous, func() {
					Tag("Relationship", "Synchronous")
				})

				Tag(styleComponent)
			})

			Component(componentMySQL, "interacts with MySQL", "Go Package", func() {
				Uses(MySQL, "Uses", Synchronous, func() {
					Tag("Relationship", "Synchronous")
				})

				Tag(styleComponent)
			})

			Component(componentRedis, "interacts with Redis", "Go Package", func() {
				Uses(Redis, "Uses", Synchronous, func() {
					Tag("Relationship", "Synchronous")
				})

				Tag(styleComponent)
			})

			Component(componentKafka, "interacts with Kafka", "Go Package", func() {
				Uses(Kafka, "Uses", Asynchronous, func() {
					Tag("Relationship", "Asynchronous")
				})

				Tag(styleComponent)
			})

			Component("internal.service", "interacts with all datastores", "Go Package", func() {
				Uses(componentElasticsearch, "Reads records from", Synchronous, func() {
					Tag("Relationship", "Synchronous")
				})

				Uses(componentKafka, "Produce events to", Synchronous, func() {
					Tag("Relationship", "Synchronous")
				})

				Uses(componentMySQL, "Uses", Synchronous, func() {
					Tag("Relationship", "Synchronous")
				})

				Uses(componentRedis, "Uses", Synchronous, func() {
					Tag("Relationship", "Synchronous")
				})

				Tag(styleComponent)
			})

			Component("internal.rest", "defines HTTP handlers", "Go Package", func() {
				Uses("internal.service", "Uses", Synchronous, func() {
					Tag("Relationship", "Synchronous")
				})

				Tag(styleComponent)
			})

			Tag(styleContainer)
		})

		Tag(styleSoftwareSystem)
	})

	Views(func() {
		SystemContextView(System, "Manga Scraper System", func() {
			AddDefault()

			EnterpriseBoundaryVisible()
		})

		ContainerView(softwareSystem, "Containers", "Container diagram for the Manga Scraper System", func() {
			AddDefault()

			SystemBoundariesVisible()
		})

		ComponentView(RESTfulAPI, "RESTful API", "Component diagram for the REST Server", func() {
			AddDefault()

			ContainerBoundariesVisible()
		})

		Styles(func() {
			ElementStyle(styleSoftwareSystem, func() {
				Background("#1168bd")
				Color("#ffffff")
			})

			ElementStyle(stylePerson, func() {
				Background("#08427b")
				Color("#ffffff")
				Shape(ShapePerson)
			})

			ElementStyle(styleComponent, func() {
				Background("#85bbf0")
				Color("#000000")
			})

			ElementStyle(styleContainer, func() {
				Background("#438dd5")
				Color("#ffffff")
			})

			ElementStyle(styleDatabase, func() {
				Shape(ShapeCylinder)
			})
		})
	})
})
