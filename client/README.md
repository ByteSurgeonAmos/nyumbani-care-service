# NyumbaniCare Frontend

This repository contains the frontend code for NyumbaniCare, a healthcare platform that provides accessible healthcare services at home.

## Features

- User authentication (login, registration)
- Test kits browsing and ordering
- Health education resources
- Medical records management
- Telehealth consultations
- Symptom checking
- Admin dashboard for managing orders, tests, and users

## Technology Stack

- SvelteKit: Modern JavaScript framework
- TypeScript: For type safety
- Tailwind CSS: For styling
- Axios: For API calls
- Lucide: For icons
- Svelte French Toast: For notifications

## Developing

Once you've cloned the project and installed dependencies with `npm install`, start a development server:

```bash
npm run dev

# or start the server and open the app in a new browser tab
npm run dev -- --open
```

## Building

To create a production version of your app:

```bash
npm run build
```

You can preview the production build with `npm run preview`.

> To deploy your app, you may need to install an [adapter](https://svelte.dev/docs/kit/adapters) for your target environment.
