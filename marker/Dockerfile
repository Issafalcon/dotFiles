FROM python:3.11-slim

RUN apt-get update && apt-get install -y \
    libgl1 \
    libglib2.0-0 \
    libsm6 \
    libxext6 \
    libxrender-dev \
    libgomp1 \
    poppler-utils \
    tesseract-ocr \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Install PyTorch CPU only (much smaller than CUDA build)
RUN pip install --no-cache-dir torch torchvision --index-url https://download.pytorch.org/whl/cpu

# Install marker with full extras (EPUB, DOCX, etc.)
RUN pip install --no-cache-dir "marker-pdf[full]"

WORKDIR /workspace

# All HuggingFace/transformers model cache paths unified under /models.
# Mount a host directory or named volume here to persist between runs:
#   docker run -v marker-models:/models ...
# HF_HOME=/models -> hub cache lands at /models/hub/
# TRANSFORMERS_CACHE and HF_HUB_CACHE must point to that same subdirectory.
ENV HF_HOME=/models
ENV HF_HUB_CACHE=/models/hub
ENV TRANSFORMERS_CACHE=/models/hub
ENV TORCH_HOME=/models/torch

ENTRYPOINT ["marker_single"]
