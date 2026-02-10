# Deployment Guide

This guide covers how to deploy the entire application:
1.  **ML Service** (Python/FastAPI) on **Render**
2.  **Backend** (Go) on **Render**
3.  **Frontend** (React) on **Netlify**

---

## 1. Prerequisites

- Ensure all your code is pushed to **GitHub**.
- Make sure `ml-service` folder is committed.
- Have your **Supabase** credentials ready (`DATABASE_URL`, `SUPABASE_URL`, `SUPABASE_KEY`).

---

## 2. Deploy ML Service (Render)

This service processes certificates. It needs to be deployed first so we can give its URL to the Backend.

1.  Log in to [Render.com](https://render.com).
2.  Click **New +** -> **Web Service**.
3.  Connect your GitHub repository.
4.  Standard Settings:
    - **Name**: `eduvault-ml-service` (or similar)
    - **Root Directory**: `ml-service`
    - **Runtime**: **Python 3**
    - **Build Command**: `pip install -r requirements.txt`
    - **Start Command**: `uvicorn api:app --host 0.0.0.0 --port $PORT`
5.  Click **Create Web Service**.
6.  **Wait** for the build to finish.
7.  **Copy the URL** provided by Render (e.g., `https://eduvault-ml.onrender.com`).

---

## 3. Deploy Backend (Render)

The Go backend handles API requests and database interactions.

1.  On Render dashboard, click **New +** -> **Web Service**.
2.  Connect the **Same GitHub repository**.
3.  Standard Settings:
    - **Name**: `eduvault-backend`
    - **Root Directory**: `backend-c`
    - **Runtime**: **Go**
    - **Build Command**: `go build -o server ./cmd/server`
    - **Start Command**: `./server`
4.  **Environment Variables** (Scroll down to "Environment Variables"):
    - Add `DATABASE_URL`: *(Value from your `backend-c/.env`)*
    - Add `ML_SERVICE_URL`: *(The URL you copied in Step 2, e.g., `https://eduvault-ml.onrender.com`)* **(Do not add /verify at the end, just the base URL)**
    - Add `APP_ENV`: `production`
    - Add `ALLOWED_EMAIL_DOMAIN`: `citchennai.net`
5.  Click **Create Web Service**.
6.  **Wait** for the build to finish and the service to go live.
7.  **Copy the Backend URL** (e.g., `https://eduvault-backend.onrender.com`).

---

## 4. Deploy Frontend (Netlify)

The React frontend.

1.  Log in to [Netlify](https://www.netlify.com).
2.  Click **Add new site** -> **Import from existing project**.
3.  Select **GitHub** and choose your repository.
4.  **Build Settings**:
    - **Base directory**: `frontend-c`
    - **Build command**: `npm run build`
    - **Publish directory**: `frontend-c/dist` (Netlify might auto-detect `dist`, ensure it knows it's inside `frontend-c`)
5.  **Environment Variables** (Click "Add environment variables" or go to Site Configuration > Environment variables):
    - Key: `VITE_API_BASE_URL`
    - Value: *(The Backend URL you copied in Step 3, e.g., `https://eduvault-backend.onrender.com`)*
    - Key: `VITE_SUPABASE_URL`
    - Value: *(From `frontend-c/.env`)*
    - Key: `VITE_SUPABASE_PUBLISHABLE_KEY`
    - Value: *(From `frontend-c/.env`)*
    - Key: `VITE_SUPABASE_PROJECT_ID`
    - Value: *(From `frontend-c/.env`)*
6.  Click **Deploy Site**.

---

## 5. Final Checks

1.  Open your **Netlify URL**.
2.  Log in (Database should work).
3.  Go to **Faculty Dashboard** -> Upload a Certificate.
4.  The Frontend calls Backend -> Backend calls ML Service -> ML Service validates -> DB updated.
5.  Check the status updates to **VERIFIED**!
