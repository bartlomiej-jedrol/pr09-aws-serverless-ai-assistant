# pr09-aws-serverless-ai-assistant

- This project is a serverless AI assistant application built using AWS Lambda. It integrates with various services to provide functionalities such as link shortening, request routing, OpenAI model interaction, Slack message processing, and API authorization.

Key Functionalities:
- Link Shortening: Utilizes the Dub API to shorten URLs. Implemented in the link_shortener Lambda function.
- Request Routing: Routes API Gateway requests based on user intentions, determined using OpenAI's GPT model. Implemented in the router_lambda Lambda function.
- OpenAI Integration: Provides tools to interact with OpenAI models for generating AI-driven responses.
- Slack Message Processing: Parses and processes Slack messages to extract and handle relevant data.
- AWS Utilities: Includes utility functions for AWS service interactions, such as environment variable management and configuration loading.
- Authorization: Implements a custom authorizer for API Gateway to secure access using AWS Secrets Manager.
This modular architecture allows for scalable and maintainable serverless applications, capable of handling diverse tasks efficiently.