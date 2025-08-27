Agentic SDLC API - Starter (Complete)

This starter provides a minimal API server wired to Ollama via langchaingo (pinned to v0.1.13).
It is a scaffold to extend into the configuration-driven workflow system described earlier.

Requirements:
- Go 1.22+
- Ollama running locally (ollama serve) and a model pulled, e.g. `ollama pull llama3`
- Set OLLAMA_MODEL env var or pass -model flag

Run:
    go run ./cmd/server -model llama3 -addr :8080

Endpoints (minimal):
- POST /api/v1/workflows          -> create workflow (JSON body)
- GET  /api/v1/workflows         -> list workflows
- POST /api/v1/workflows/{id}/executions -> start execution

This is an in-memory prototype (no DB). Extend to persist storage and secure auth for production.

----------------------------------------------------

v1 - this will not work
curl -X POST http://localhost:8080/api/v1/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Shopping App",
    "agents": ["Product", "Design", "Code"]
  }'

{"id":"f4499b93-4966-45ed-83c2-c267a72023ad","name":"New Shopping App","spec":""}

curl http://localhost:8080/api/v1/workflows

[{"id":"f4499b93-4966-45ed-83c2-c267a72023ad","name":"New Shopping App","spec":""}]


curl -X POST http://localhost:8080/api/v1/workflows/f4499b93-4966-45ed-83c2-c267a72023ad/executions

---------------------------
v1  - has been modified from previous update


curl -X POST http://localhost:8080/api/v1/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Shopping App",
    "agents": ["Code"],
    "spec": "Build a simple task management web app with backend in Go and frontend in React."
  }'

{
    "id":"89856e9c-5855-44b1-bf82-295476957d90",
    "name":"New Shopping App","spec":"Build a simple task management web app with backend in Go and frontend in React.",
    "agents":["Product","Design","Code"]
}

 curl -X POST http://localhost:8080/api/v1/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Shopping App",
    "agents": ["Code"],
    "spec": "Build a simple task management web app with backend in Go and frontend in React."
  }'

{"id":"93345d2f-e2fa-4b48-9a68-9fa71e0f5a6c","name":"New Shopping App","spec":"Build a simple task management web app with backend in Go and frontend in React.","agents":["Code"]}


curl -X POST http://localhost:8080/api/v1/workflows/1211e0e7-8f4a-41d9-af6f-3191dcfa1e82/executions


[{"id":"9c1ec936-661e-4439-9bd8-69bfd5a1a3c9","name":"New Shopping App","spec":"Build a simple task management web app with backend in Go and frontend in React.","agents":["Product","Design","Code","Test","Deploy"]}]


curl -X POST http://localhost:8080/api/v1/workflows/9c1ec936-661e-4439-9bd8-69bfd5a1a3c9/executions


-------------------------------------------------

v2 - with db -pgql 


it is working
run the init.sql
then with go run seeder.go


curl -X POST http://localhost:8080/api/v2/workflows/55528e2c-05ec-4116-afd3-0ce9581c47d8/executions   -H "Content-Type: application/json"   -d '{"input":"I want an e-commerce backend"}'
{"results":{"22695631-c485-46c0-831b-0f7c1d83d190":"\nThere are several options for building an e-commerce backend, depending on your specific needs and requirements. Here are a few popular options:\n\n1. Magento: Magento is a popular open-source e-commerce platform that allows you to build customized online stores with advanced features like product catalogs, shopping carts, and payment gateways. It also has a large community of developers who contribute to its codebase and provide support for various integrations.\n2. Shopify: Shopify is a cloud-based e-commerce platform that allows you to build customized online stores with a range of features like product catalogs, shopping carts, and payment gateways. It also has a user-friendly interface and a large community of developers who contribute to its codebase and provide support for various integrations.\n3. WooCommerce: WooCommerce is a popular open-source e-commerce plugin for WordPress that allows you to build customized online stores with advanced features like product catalogs, shopping carts, and payment gateways. It also has a large community of developers who contribute to its codebase and provide support for various integrations.\n4. Spree Commerce: Spree Commerce is an open-source e-commerce platform that allows you to build customized online stores with advanced features like product catalogs, shopping carts, and payment gateways. It also has a large community of developers who contribute to its codebase and provide support for various integrations.\n5. Vend: Vend is an open-source e-commerce platform that allows you to build customized online stores with advanced features like product catalogs, shopping carts, and payment gateways. It also has a large community of developers who contribute to its codebase and provide support for various integrations.\n\nUltimately, the best option for your e-commerce backend will depend on your specific needs and requirements. You may want to consider factors like ease of use, scalability, customization options, and integration with other systems.","841219b1-d443-444b-995e-217d2a64a5da":"\nThere are several options for building an e-commerce backend, depending on your specific needs and requirements. Here are a few popular options:\n\n1. Magento: Magento is a popular open-source e-commerce platform that allows you to build customized online stores with advanced features like product catalogs, shopping carts, and payment gateways. It also has a large community of developers who contribute to its codebase and provide support for various integrations.\n2. Shopify: Shopify is a cloud-based e-commerce platform that allows you to build customized online stores with a range of features like product catalogs, shopping carts, and payment gateways. It also has a user-friendly interface and a large community of developers who contribute to its codebase and provide support for various integrations.\n3. WooCommerce: WooCommerce is a popular open-source e-commerce plugin for WordPress that allows you to build customized online stores with advanced features like product catalogs, shopping carts, and payment gateways. It also has a large community of developers who contribute to its codebase and provide support for various integrations.\n4. Spree Commerce: Spree Commerce is an open-source e-commerce platform that allows you to build customized online stores with advanced features like product catalogs, shopping carts, and payment gateways. It also has a large community of developers who contribute to its codebase and provide support for various integrations.\n5. Vend: Vend is an open-source e-commerce platform that allows you to build customized online stores with advanced features like product catalogs, shopping carts, and payment gateways. It also has a large community of developers who contribute to its codebase and provide support for various integrations.\n\nUltimately, the best option for your e-commerce backend will depend on your specific needs and requirements. You may want to consider factors like ease of use, scalability, customization options, and integration with other systems.","f433d8f4-e90c-445d-bbb5-0cc545883c8f":"\nThere are several options for building an e-commerce backend, depending on your specific needs and requirements. Here are a few popular options:\n\n1. Magento: Magento is a popular open-source e-commerce platform that allows you to build customized online stores with advanced features like product catalogs, shopping carts, and payment gateways. It also has a large community of developers who contribute to its codebase and provide support for various integrations.\n2. Shopify: Shopify is a cloud-based e-commerce platform that allows you to build customized online stores with a range of features like product catalogs, shopping carts, and payment gateways. It also has a user-friendly interface and a large community of developers who contribute to its codebase and provide support for various integrations.\n3. WooCommerce: WooCommerce is a popular open-source e-commerce plugin for WordPress that allows you to build customized online stores with advanced features like product catalogs, shopping carts, and payment gateways. It also has a large community of developers who contribute to its codebase and provide support for various integrations.\n4. Spree Commerce: Spree Commerce is an open-source e-commerce platform that allows you to build customized online stores with advanced features like product catalogs, shopping carts, and payment gateways. It also has a large community of developers who contribute to its codebase and provide support for various integrations.\n5. Vend: Vend is an open-source e-commerce platform that allows you to build customized online stores with advanced features like product catalogs, shopping carts, and payment gateways. It also has a large community of developers who contribute to its codebase and provide support for various integrations.\n\nUltimately, the best option for your e-commerce backend will depend on your specific needs and requirements. You may want to consider factors like ease of use, scalability, customization options, and integration with other systems."}}


-----------------
go run ./cmd/server/main.go -backend=ollama -model=phi -addr=:8080

curl -X POST http://localhost:8080/api/v2/workflows/cd4e4333-6f9d-4669-b369-4a88b2fd4354/executions   -H "Content-Type: application/json"   -d '{"input":"I want an e-commerce backend"}'