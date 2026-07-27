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

# All model cache paths unified under /models so a single mounted volume persists
# everything across container runs: docker run -v marker-models:/models ...
#
# HuggingFace (used by marker internals)
ENV HF_HOME=/models
ENV HF_HUB_CACHE=/models/hub
ENV TRANSFORMERS_CACHE=/models/hub
ENV TORCH_HOME=/models/torch
# Surya (OCR/layout engine used by marker) writes to its own cache dir by default
# (/root/.cache/datalab/models). Override it to land inside the volume.
ENV MODEL_CACHE_DIR=/models/datalab

# Reduce CPU memory pressure
ENV MALLOC_TRIM_THRESHOLD_=100000
ENV PYTORCH_NO_CUDA_MEMORY_CACHING=1
# Use float16 weights where possible to halve model memory usage
ENV TORCH_DTYPE=float16

ENTRYPOINT ["marker_single"]
