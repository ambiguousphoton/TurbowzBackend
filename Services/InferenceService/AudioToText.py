import fastapi
import os
import whisper
import tempfile

app = fastapi.FastAPI()
model = whisper.load_model("base")

ffmpeg_path = os.path.join(os.path.dirname(__file__), "tools")
os.environ["PATH"] += os.pathsep + ffmpeg_path

@app.post("/audio-to-text/")
async def audio_to_text(file: fastapi.UploadFile = fastapi.File(...)):
    print("Received file:", file.filename)
    
    
    with tempfile.NamedTemporaryFile(delete=False, suffix=".mp3") as tmp:
        tmp.write(await file.read())
        tmp_path = tmp.name
    
    result = model.transcribe(tmp_path)
    os.remove(tmp_path)
    print(result["text"])
    return {"transcribe": result["text"]}  


 