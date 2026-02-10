# Use Python 3.10
FROM python:3.10

# Set working directory
WORKDIR /code

# Install system dependencies for OpenCV and Poppler
RUN apt-get update && apt-get install -y \
    libgl1-mesa-glx \
    libglib2.0-0 \
    poppler-utils \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements
COPY requirements.txt .

# Install Python dependencies (including CPU-only torch to save space)
RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir -r requirements.txt && \
    pip install --no-cache-dir torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cpu

# Copy the rest of the application
COPY . .

# Create a non-root user (required by Hugging Face Spaces)
RUN useradd -m -u 1000 user
USER user
ENV HOME=/home/user \
    PATH=/home/user/.local/bin:$PATH

WORKDIR $HOME/app
COPY --chown=user . $HOME/app

# Expose port (Hugging Face uses 7860 by default)
EXPOSE 7860

# Start command
CMD ["uvicorn", "api:app", "--host", "0.0.0.0", "--port", "7860"]
