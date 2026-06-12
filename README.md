<p align="center">
  <img src="app/assets/images/readme-logo1.png" alt="Jovvix Logo" width="200"/>
</p>

<h3 align="center">
<b>
Open-Source Quizzing Built for Live Quiz Experiences
</b>
</h3 >



# jovVix

- jovVix is a fun and interactive platform where users can enjoy playing quizzes while admins have the ability to create engaging and diverse quizzes.
- Designed to provide a fun and educational experience while ensuring smooth admin management, jovVix is the perfect solution for interactive quiz-based learning, competitions, or corporate events.

- App is an interactive, real-time platform that supports diverse
  question types like
- App provides instant feedback and live ranks and leaderboards.
  With features like
   - In-depth analytics
   - real dynamic avatars
   - customizable csv with different question types like: Image based, code based, survey based questions
<br><br>

## Live Demo & Website

- **Website**: Visit [https://jovvix.com](https://jovvix.com) to learn more about the platform, explore the **About** section, and access full **Documentation**.
- **Live Demo**: Try out the platform live at [https://app.jovvix.com](https://app.jovvix.com) to experience how the quiz system works — including joining a quiz, answering questions in real-time, and viewing live leaderboards.

## Key Attributes:
- Cloud-Based: The application can be deployed on cloud platforms, providing high availability and scalability for users and organizations.
- Environment Agnostic: Compatible with multiple environments (local, development, staging, and production) for seamless deployment across setups.
- Cross-Platform: Designed for a smooth experience across devices and operating systems.
- Modular Architecture: Built with a component-based architecture, enabling high code reusability and scalability.
- Real-Time Features: Uses WebSockets for real-time communication and user interaction.

### Development and Quality
- Linting: Utilizes ESLint to maintain code quality and consistency.
- Authentication: Ory Kratos manages secure user authentication and session management.
- Image Handling: Images for questions and options are stored inline as base64 in PostgreSQL, removing the need for any external object-storage service.
- SMTP Services: Configured for email services, supporting verification and password recovery.
- Testing: We have performed the unit testing in this application
- Containerization: Packaged with Docker for streamlined deployment and environment consistency.
- RESTful API: Offers a RESTful API for integration and easy access to app functionality.
- Environment Variables: Follows the 12-factor app principles with environment variable configuration.

### User Experience
- Intuitive UI/UX: Designed for a user-friendly and engaging interface.
- Gamification: Includes elements like points, ranks, and leaderboards to boost user engagement.

## Features

- **Real-time Interaction:** Supports real-time quizzes with instant feedback and live leaderboards, ensuring a highly engaging user experience.
- **Public Quizzes:** Curated quizzes can be published publicly and explored right from the homepage — no account required:
  - Anyone, including guests who haven't logged in, can start a public quiz instantly
  - The quiz session can be shared with others so they can join and participate together
  - The host of a public quiz can also play along as a participant
- **Multiple Question Types:** Allows admins to create quizzes with different question formats:
  - Multiple-choice questions
  - Survey-based questions
  - Code-based challenges
  - Image-based questions
- **In-Depth Analytics:** Provides detailed insights and analytics for both users and admins, including answer breakdowns and performance tracking.
- **Reports-Dashboard:** Provides an in depth metrices of user's response to admin and user analysis statistics
- **Dynamic Avatars:** Real-time avatar updates for users in various stages:
  - When joining and waiting for a quiz
  - During the quiz at each question
  - After the quiz, showcasing winners and individual performance
- **Customization Options:** Offers customizable quiz features for educators, workshop hosts, and event organizers to suit their needs.
- **Share quiz feature** Users can share their quiz to other users through email and can grant their preferred permission while sharing to the other one
- **Add Questions via CSV or UI:** Add questions to a quiz in two ways:
    - **CSV Uploads:** Enables efficient quiz creation through CSV file uploads, allowing users to preview, edit, and seamlessly create quizzes
    - **Multiple CSV Uploads:** Upload multiple CSV files to combine questions from different files into a single quiz
    - **Add Questions through UI:** You can now also add questions directly through the UI, creating and editing them one by one without preparing a CSV file
    - You can see the CSV formatting guidelines here : [csv-formatting-guide.md](docs/csv-formatting-guide.md)
- **Mobile-Friendly Design:** Fully responsive and works seamlessly on mobile and desktop devices.
- **Admin Tools:** Advanced admin panel for managing quizzes, participants, and results.
- **API Documentation with Swagger:** Provides a visual interface for exploring API endpoints and testing requests, Useful for developers working with the API
- **Open Source:** Fully open-source platform, allowing developers to contribute, customize, and extend the app.

## Table of Contents:

- [About jovVix](#jovvix)
  - [Quickstart](#quickstart)
- [Features](#features)
- [Getting started](#getting-started-for-local-setup-from-source)
  - [Local setup](#prerequisites-1)
- [Documentation](#architecture-overview)
  - [Architecture Overview](#architecture-overview)
  - [Overview of API documentation](#api-overview)
  - [Upgrading and changelog](#upgrading-and-changelog)
- [Contributing to jovVix](#contributing-to-jovvix)
- [Code Of Conduct](#code-of-conduct)
- [Develop]()
  - [Guide](#guide)
  - [Adding shadcn-vue Components](#adding-shadcn-vue-components)
  - [Collaborators](#collaborators)
  - [Dependencies](#dependencies)

## Prerequisites
- Docker latest version installed in your system

## Quickstart

> **Important:**
>  Ensure all the tools mentioned in the prerequisites are installed before proceeding.


- First clone the repository on your terminal using command given below
  git clone https://github.com/improwised/jovVix.git


- Navigate to the project directory using:
   cd jovVix


- Configure your env settings into .env.docker if you want to integrate any changes otherwise keep the default ones

- Then build and run the docker compose file in your environment by following command
  docker-compose up --build


- Your app is now running successfully and you can access it on the ip:port as :
  127.0.0.1:5000

- For verification of email go to mailpit localhost and port
  127.0.0.1:8025


## Getting Started For local setup from source

### Prerequisites

| Package | Version |
| --- | --- |
| [Node.js](https://nodejs.org/en/) | v18.0+ |
| [Nuxt](https://nuxt.com/) | v3.0.0 |
| [Go](https://golang.org/) | v1.21+ |


- jovVix is a fun and interactive platform where users can enjoy playing quizzes while admins have the ability to create engaging and diverse quizzes.


## Installation steps

> **Important:**
>  Ensure all the tools mentioned in the prerequisites are installed before proceeding.



- First clone the repository using
  git clone https://github.com/Improwised/jovVix.git



- Navigate to the project directory:
   cd jovVix



- Then Copy environment files of both app and api folders  using:
   cp api/.env.example api/.env
   cp app/.env.example app/.env



- Start backend services:
  cd api
  docker-compose up
  go run app.go migrate
  go run app.go api



> **Warning**
> Install all tools and technologies we have mentioned above in prerequisites


- Then Install frontend dependencies:
  cd ../app
  npm install

- Afterwards run the frontend development server:
  npm i
  npm run dev

- Then you have setup jovVix successfully in your local environment


- Designed to provide a fun and educational experience while ensuring smooth admin management, jovVix is the perfect solution for interactive quiz-based learning, competitions, or corporate events.

## Documentation

## Architecture Overview:

### System Architecture

- **Backend:** jovVix uses a Golang backend to handle the server-side logic.
- **Frontend:** Vue.js and nuxt framework is employed for a single-page application (SPA) architecture and component-based development.
- **Database:** PostgreSQL is used to handle concurrent requests and manage data efficiently.
- **Caching:** Valkey is utilized for caching and manipulating users' data and requests, improving performance.
- **Authentication:** The app leverages Ory Kratos for user authentication and uses SMTP services for password recovery and email verification flows.
- **Real-Time Communication:** WebSockets are implemented for handling multiple sessions and managing cookies effectively.
- **Base64:** We are storing images as base64 


 ### API Overview:
- locally you could start your API server by running the following command :
http://127.0.0.1:3000/api/v1/docs

- This would open the swagger documentation

- **Backend:** jovVix uses a Golang backend to handle the server-side logic.
- This would open the swagger documentation

 ### Upgrading and changelog:

- It's being managed in this file: CHANGELOG.md

## Contributing to jovVix:
- You can see the contribution guidelines here : [CONTRIBUTING.md](./CONTRIBUTING.md)

## Code of Conduct:
- This platforms also provides the [Code of Conduct](./CODE_OF_CONDUCT.md)

## Developer

## Guide:
If you're a developer looking to contribute or modify jovVix, here is a brief guide to get started: [Getting started](#getting-started-for-local-setup-from-source) OR, you could just `docker-compose up`

### Code Structure Overview:

- app/: Contains Vue.js components for the quiz interface and holds frontend pages and all
- api/: Contains Golang source code for the server-side logic.

## Adding shadcn-vue Components:

- This project is configured with [shadcn-vue](https://www.shadcn-vue.com/) via the `shadcn-nuxt` module (version `2.2.0`). The configuration lives in [`app/components.json`](app/components.json) and uses the `new-york` style with the `neutral` base color and the `lucide` icon library.
- To add a new component to the project, run the shadcn-vue CLI from inside the `app/` directory:
  cd app
  npx shadcn-vue@2.2.0 add <component-name>


- For example, to add the `button` component:
  npx shadcn-vue@2.2.0 add button


- The CLI reads `components.json` and installs the component into `app/components/ui/<component-name>/` using the aliases defined in that file (`@/components/ui`, `@/lib/utils`, etc.).
- Please pin to version `2.2.0` of the CLI when adding components so that generated code stays consistent with the rest of the project. If you need to upgrade the CLI/module, bump both `shadcn-nuxt` in [`app/package.json`](app/package.json) and the version used in the command above in the same PR.
- For the full list of available components and their usage, refer to the [shadcn-vue documentation](https://www.shadcn-vue.com/docs/components).

## Collaborators:
- Thanks to all the people who already contributed!

<a href="https://github.com/Improwised/jovVix/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=Improwised/jovVix" alt="contributors">
</a>

## Dependencies:

This project uses the following key dependencies:

- Frontend:

  - Vue.js: For building reactive user interfaces.
  - Nuxt 3: Framework for server-side rendering and static site generation.

- Backend:

  - Go (Golang): High-performance backend logic.
  - PostgreSQL: Relational database management.

- Other Tools:

  - Valkey: Used for caching and session management.
  - Ory Kratos: For user authentication.
  - Mailpit setup locally for the SMTP like server

# Design & Image Credits

- The **jovVix logo** is designed by [Nirav Raval](https://www.linkedin.com/in/nirav-raval-06732858).
- All other images and visual assets used in the platform are either custom-designed or sourced from [Freepik](https://www.freepik.com/) under their free license.

> Note: Assets from Freepik are used in accordance with their [license terms](https://www.freepikcompany.com/legal#nav-freepik-license).

# License:

jovVix is put under a dual-licensing scheme: In general all of the provided code is open source via [GNU AGPL 3.0](https://www.gnu.org/licenses/agpl-3.0.en.html), please see the [LICENSE](LICENSE.txt) file for more details.