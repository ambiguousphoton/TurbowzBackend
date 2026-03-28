import fastapi
import os
import time
import tempfile
from fastapi import Request
from faster_whisper import WhisperModel

# -------------------- APP --------------------
app = fastapi.FastAPI()

# -------------------- CPU THREADING --------------------
os.environ["OMP_NUM_THREADS"] = "8"
os.environ["MKL_NUM_THREADS"] = "8"

# -------------------- MODEL (TINY + CPU) --------------------
model = WhisperModel(
    "tiny",
    device="cpu",
    compute_type="int8",
    cpu_threads=8,
    num_workers=2
)

# -------------------- FFMPEG PATH --------------------
ffmpeg_path = os.path.join(os.path.dirname(__file__), "tools")
os.environ["PATH"] += os.pathsep + ffmpeg_path


# -------------------- MIDDLEWARE --------------------
@app.middleware("http")
async def latency_middleware(request: Request, call_next):
    start = time.perf_counter()

    response = await call_next(request)

    latency_ms = int((time.perf_counter() - start) * 1000)

    response.headers["X-Response-Time-ms"] = str(latency_ms)

    print(
        f"method={request.method} "
        f"path={request.url.path} "
        f"status={response.status_code} "
        f"latency_ms={latency_ms}"
    )

    return response


# -------------------- ENDPOINT --------------------
@app.post("/audio-to-text/")
async def audio_to_text(file: fastapi.UploadFile = fastapi.File(...)):
    print("Received file:", file.filename)

    with tempfile.NamedTemporaryFile(delete=False, suffix=".mp3") as tmp:
        tmp.write(await file.read())
        tmp_path = tmp.name

    segments, info = model.transcribe(
        tmp_path,
        beam_size=1,        # FAST
        vad_filter=True,    # skip silence
        chunk_length=30,    # stable for long audio
        language="en" ,     # skip language detection
        without_timestamps=True,
    )

    text = "".join(segment.text for segment in segments)

    os.remove(tmp_path)

    return {"transcribe": text}
